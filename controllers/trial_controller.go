/*
Copyright 2019 GramLabs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	redskyv1alpha1 "github.com/redskyops/k8s-experiment/pkg/apis/redsky/v1alpha1"
	"github.com/redskyops/k8s-experiment/pkg/controller/metric"
	"github.com/redskyops/k8s-experiment/pkg/controller/template"
	redskytrial "github.com/redskyops/k8s-experiment/pkg/controller/trial"
	"github.com/redskyops/k8s-experiment/pkg/util"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// TrialReconciler reconciles a Trial object
type TrialReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *TrialReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&redskyv1alpha1.Trial{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=batch;extensions,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch;extensions,resources=jobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=list
// +kubebuilder:rbac:groups="",resources=services,verbs=list
// +kubebuilder:rbac:groups=redskyops.dev,resources=trials,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=redskyops.dev,resources=trials/status,verbs=get;update;patch

func (r *TrialReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("trial", req.NamespacedName)
	now := metav1.Now()

	// Fetch the Trial instance
	trial := &redskyv1alpha1.Trial{}
	if err := r.Get(ctx, req.NamespacedName, trial); err != nil {
		return util.IgnoreNotFound(err)
	}

	// Ahead of everything is the setup/teardown (contains finalization logic)
	if result, err := redskytrial.ManageSetup(r.Client, r.Scheme, ctx, &now, trial); result != nil {
		if err != nil {
			log.Error(err, "Setup task failed")
		}
		return *result, err
	}

	// If we are in a finished or deleted state there is nothing for us to do
	if redskytrial.IsTrialFinished(trial) || !trial.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	// Copy the patches over from the experiment
	if len(trial.Spec.PatchOperations) == 0 {
		e := &redskyv1alpha1.Experiment{}
		if err := r.Get(ctx, trial.ExperimentNamespacedName(), e); err != nil {
			return ctrl.Result{}, err
		}
		if err := checkAssignments(trial, e, log); err != nil {
			return ctrl.Result{}, err
		}
		if err := evaluatePatches(trial, e); err != nil {
			return ctrl.Result{}, err
		}
		if len(trial.Spec.PatchOperations) > 0 {
			// We know we have at least one patch to apply, use an unknown status until we start applying them
			redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialPatched, corev1.ConditionUnknown, "", "", &now)
			return r.forTrialUpdate(trial, ctx, log)
		}
	}

	// Check the "initializer" annotation, do not progress unless it is empty (don't requeue, wait for a change)
	if trial.HasInitializer() {
		return ctrl.Result{}, nil
	}

	// Apply the patches
	for i := range trial.Spec.PatchOperations {
		p := &trial.Spec.PatchOperations[i]
		if p.AttemptsRemaining == 0 {
			continue
		}

		u := unstructured.Unstructured{}
		u.SetName(p.TargetRef.Name)
		u.SetNamespace(p.TargetRef.Namespace)
		u.SetGroupVersionKind(p.TargetRef.GroupVersionKind())
		if err := r.Patch(ctx, &u, client.ConstantPatch(p.PatchType, p.Data)); err != nil {
			p.AttemptsRemaining = p.AttemptsRemaining - 1
			if p.AttemptsRemaining == 0 {
				// There are no remaining patch attempts remaining, fail the trial
				redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialFailed, corev1.ConditionTrue, "PatchFailed", err.Error(), &now)
			}
		} else {
			p.AttemptsRemaining = 0
			if p.Wait {
				// We successfully applied a patch that requires a wait, use an unknown status until we start waiting
				redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialStable, corev1.ConditionUnknown, "", "", &now)
			}
		}

		// We have started applying patches (success or fail), transition into a false status
		redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialPatched, corev1.ConditionFalse, "", "", &now)
		return r.forTrialUpdate(trial, ctx, log)
	}

	// If there is a patched condition that is not yet true, update the status
	if cc, ok := redskytrial.CheckCondition(&trial.Status, redskyv1alpha1.TrialPatched, corev1.ConditionTrue); ok && !cc {
		redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialPatched, corev1.ConditionTrue, "", "", &now)
		return r.forTrialUpdate(trial, ctx, log)
	}

	// Wait for a stable (ish) state
	for i := range trial.Spec.PatchOperations {
		p := &trial.Spec.PatchOperations[i]
		if !p.Wait {
			continue
		}

		var requeueAfter time.Duration
		if err := redskytrial.WaitForStableState(r.Client, ctx, log, p); err != nil {
			if serr, ok := err.(*redskytrial.StabilityError); ok && serr.RetryAfter > 0 {
				// Mark the trial as not stable and wait
				redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialStable, corev1.ConditionFalse, "Waiting", err.Error(), &now)
				requeueAfter = serr.RetryAfter
			} else {
				// No retry delay specified, fail the whole trial
				redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialFailed, corev1.ConditionTrue, "WaitFailed", err.Error(), &now)
			}
		} else {
			// We have successfully waited for one patch so we are no longer "unknown"
			redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialStable, corev1.ConditionFalse, "", "", &now)
			p.Wait = false
		}

		// Inject the retry delay if necessary
		rr, re := r.forTrialUpdate(trial, ctx, log)
		if re == nil && requeueAfter > 0 {
			rr.RequeueAfter = requeueAfter
		}
		return rr, re
	}

	// If there is a stable condition that is not yet true, update the status
	if cc, ok := redskytrial.CheckCondition(&trial.Status, redskyv1alpha1.TrialStable, corev1.ConditionTrue); ok && !cc {
		redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialStable, corev1.ConditionTrue, "", "", &now)
		return r.forTrialUpdate(trial, ctx, log)
	}

	// Find jobs labeled for this trial
	list := &batchv1.JobList{}
	matchingSelector, err := util.MatchingSelector(trial.GetJobSelector())
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := r.List(ctx, list, matchingSelector); err != nil {
		return ctrl.Result{}, err
	}

	// Update the trial run status using the job status
	needsJob := true
	for i := range list.Items {
		// Setup jobs always have "role=trialSetup" so ignore jobs with that label
		if list.Items[i].Labels[redskyv1alpha1.LabelTrialRole] != "trialSetup" {
			if applyJobStatus(trial, &list.Items[i], &now) {
				return r.forTrialUpdate(trial, ctx, log)
			}
			needsJob = false
		}
	}

	// Create a trial run job if needed
	if needsJob {
		job := createJob(trial)
		if err := controllerutil.SetControllerReference(trial, job, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		err = r.Create(ctx, job)
		return ctrl.Result{}, err
	}

	// The completion time will be non-nil as soon as the (a?) trial run job finishes
	if trial.Status.CompletionTime != nil {
		e := &redskyv1alpha1.Experiment{}
		if err = r.Get(ctx, trial.ExperimentNamespacedName(), e); err != nil {
			return ctrl.Result{}, err
		}

		// If we have metrics to collect, use an unknown status to fill the gap (e.g. TCP timeout) until the transition to false
		if len(e.Spec.Metrics) > 0 {
			if _, ok := redskytrial.CheckCondition(&trial.Status, redskyv1alpha1.TrialObserved, corev1.ConditionUnknown); !ok {
				redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialObserved, corev1.ConditionUnknown, "", "", &now)
				return r.forTrialUpdate(trial, ctx, log)
			}
		}

		// Look for metrics that have not been collected yet
		for _, m := range e.Spec.Metrics {
			v := findOrCreateValue(trial, m.Name)
			if v.AttemptsRemaining == 0 {
				continue
			}

			// Capture the metric
			var captureError error
			if target, err := getMetricTarget(r, ctx, &m); err != nil {
				captureError = err
			} else if value, stddev, err := metric.CaptureMetric(&m, trial, target); err != nil {
				if merr, ok := err.(*metric.CaptureError); ok && merr.RetryAfter > 0 {
					// Do not count retries against the remaining attempts
					return ctrl.Result{RequeueAfter: merr.RetryAfter}, nil
				}
				captureError = err
			} else {
				v.AttemptsRemaining = 0
				v.Value = strconv.FormatFloat(value, 'f', -1, 64)
				if stddev != 0 {
					v.Error = strconv.FormatFloat(stddev, 'f', -1, 64)
				}
			}

			// Handle any errors the occurred while collecting the value
			if captureError != nil && v.AttemptsRemaining > 0 {
				v.AttemptsRemaining = v.AttemptsRemaining - 1
				if v.AttemptsRemaining == 0 {
					redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialFailed, corev1.ConditionTrue, "MetricFailed", captureError.Error(), &now)
					if merr, ok := captureError.(*metric.CaptureError); ok {
						// Metric errors contain additional information which should be logged for debugging
						log.Error(err, "Metric collection failed", "address", merr.Address, "query", merr.Query, "completionTime", merr.CompletionTime)
					}
				}
			}

			// Set the observed condition to false since we have observed at least one, but possibly not all of, the metrics
			redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialObserved, corev1.ConditionFalse, "", "", &now)
			return r.forTrialUpdate(trial, ctx, log)
		}

		// If all of the metrics are collected, finish the observation
		if cc, ok := redskytrial.CheckCondition(&trial.Status, redskyv1alpha1.TrialObserved, corev1.ConditionTrue); ok && !cc {
			redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialObserved, corev1.ConditionTrue, "", "", &now)
		}

		// Mark the trial as completed
		redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialComplete, corev1.ConditionTrue, "", "", &now)
		return r.forTrialUpdate(trial, ctx, log)
	}

	// If nothing changed, check again
	return ctrl.Result{Requeue: true}, nil
}

// Returns from the reconcile loop after updating the supplied trial instance
func (r *TrialReconciler) forTrialUpdate(trial *redskyv1alpha1.Trial, ctx context.Context, log logr.Logger) (ctrl.Result, error) {
	// If we are going to be updating the trial, make sure the status is synchronized
	assignments := make([]string, len(trial.Spec.Assignments))
	for i := range trial.Spec.Assignments {
		assignments[i] = fmt.Sprintf("%s=%d", trial.Spec.Assignments[i].Name, trial.Spec.Assignments[i].Value)
	}
	trial.Status.Assignments = strings.Join(assignments, ", ")

	values := make([]string, len(trial.Spec.Values))
	for i := range trial.Spec.Values {
		if trial.Spec.Values[i].AttemptsRemaining == 0 {
			values[i] = fmt.Sprintf("%s=%s", trial.Spec.Values[i].Name, trial.Spec.Values[i].Value)
		}
	}
	trial.Status.Values = strings.Join(values, ", ")

	err := r.Update(ctx, trial)
	return util.IgnoreConflict(err)
}

func evaluatePatches(trial *redskyv1alpha1.Trial, e *redskyv1alpha1.Experiment) error {
	var err error
	te := template.NewTemplateEngine()
	for _, p := range e.Spec.Patches {
		po := redskyv1alpha1.PatchOperation{
			AttemptsRemaining: 3,
			Wait:              true,
		}

		// Evaluate the patch template
		po.Data, err = te.RenderPatch(&p, trial)
		if err != nil {
			return err
		}

		// If the patch is effectively null, we do not need to evaluate it
		if len(po.Data) == 0 || string(po.Data) == "null" {
			po.AttemptsRemaining = 0
		}

		// Determine the patch type
		switch p.Type {
		case redskyv1alpha1.PatchStrategic, "":
			po.PatchType = types.StrategicMergePatchType
		case redskyv1alpha1.PatchMerge:
			po.PatchType = types.MergePatchType
		case redskyv1alpha1.PatchJSON:
			po.PatchType = types.JSONPatchType
		default:
			return fmt.Errorf("unknown patch type: %s", p.Type)
		}

		// Attempt to populate the target reference
		if p.TargetRef != nil {
			p.TargetRef.DeepCopyInto(&po.TargetRef)
		} else if po.PatchType == types.StrategicMergePatchType {
			// TODO Allow strategic merge patches to specify the target reference
		}
		if po.TargetRef.Namespace == "" {
			po.TargetRef.Namespace = trial.Spec.TargetNamespace
		}
		if po.TargetRef.Namespace == "" {
			po.TargetRef.Namespace = trial.Namespace
		}

		trial.Spec.PatchOperations = append(trial.Spec.PatchOperations, po)
	}

	return nil
}

func checkAssignments(trial *redskyv1alpha1.Trial, experiment *redskyv1alpha1.Experiment, log logr.Logger) error {
	// Index the assignments
	assignments := make(map[string]int64, len(trial.Spec.Assignments))
	for _, a := range trial.Spec.Assignments {
		assignments[a.Name] = a.Value
	}

	// Verify against the parameter specifications
	var missing []string
	for _, p := range experiment.Spec.Parameters {
		if a, ok := assignments[p.Name]; ok {
			if a < p.Min || a > p.Max {
				log.Info("Assignment out of bounds", "trialName", trial.Name, "parameterName", p.Name, "assignment", a, "min", p.Min, "max", p.Max)
			}
		} else {
			missing = append(missing, p.Name)
		}
	}

	// Fail if there are missing assignments
	if len(missing) > 0 {
		return fmt.Errorf("trial %s is missing assignments for %s", trial.Name, strings.Join(missing, ", "))
	}
	return nil
}

func getMetricTarget(r client.Reader, ctx context.Context, m *redskyv1alpha1.Metric) (runtime.Object, error) {
	switch m.Type {
	case redskyv1alpha1.MetricLocal, "":
		// There is no target for local metrics
		return nil, nil
	case redskyv1alpha1.MetricPrometheus, redskyv1alpha1.MetricJSONPath:
		// Both Prometheus and JSONPath target a service
		target := &corev1.ServiceList{}
		if sel, err := util.MatchingSelector(m.Selector); err != nil {
			return nil, err
		} else if err := r.List(ctx, target, sel); err != nil {
			return nil, err
		}
		return target, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", m.Type)
	}
}

func findOrCreateValue(trial *redskyv1alpha1.Trial, name string) *redskyv1alpha1.Value {
	for i := range trial.Spec.Values {
		if trial.Spec.Values[i].Name == name {
			return &trial.Spec.Values[i]
		}
	}

	trial.Spec.Values = append(trial.Spec.Values, redskyv1alpha1.Value{Name: name, AttemptsRemaining: 3})
	return &trial.Spec.Values[len(trial.Spec.Values)-1]
}

func applyJobStatus(trial *redskyv1alpha1.Trial, job *batchv1.Job, time *metav1.Time) bool {
	var dirty bool

	if trial.Status.StartTime == nil {
		// Establish a start time if available
		trial.Status.StartTime = job.Status.StartTime.DeepCopy()
		dirty = dirty || job.Status.StartTime != nil
		if dirty && trial.Spec.StartTimeOffset != nil {
			*trial.Status.StartTime = metav1.NewTime(trial.Status.StartTime.Add(trial.Spec.StartTimeOffset.Duration))
		}
	} else if job.Status.StartTime != nil && trial.Status.StartTime.Before(job.Status.StartTime) {
		// Move the start time back
		trial.Status.StartTime = job.Status.StartTime.DeepCopy()
		dirty = true
	}

	if trial.Status.CompletionTime == nil {
		// Establish an end time if available
		trial.Status.CompletionTime = job.Status.CompletionTime.DeepCopy()
		dirty = dirty || job.Status.CompletionTime != nil
	} else if job.Status.CompletionTime != nil && trial.Status.CompletionTime.Before(job.Status.CompletionTime) {
		// Move the completion time back
		trial.Status.CompletionTime = job.Status.CompletionTime.DeepCopy()
		dirty = true
	}

	// Mark the trial as failed if the job itself failed
	for _, c := range job.Status.Conditions {
		// If activeDeadlineSeconds was used a workaround for having a sidecar, ignore the failure
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue && c.Reason != "DeadlineExceeded" {
			redskytrial.ApplyCondition(&trial.Status, redskyv1alpha1.TrialFailed, corev1.ConditionTrue, c.Reason, c.Message, time)
			dirty = true
		}
	}

	return dirty
}

func createJob(trial *redskyv1alpha1.Trial) *batchv1.Job {
	job := &batchv1.Job{}

	// Start with the job template
	if trial.Spec.Template != nil {
		trial.Spec.Template.ObjectMeta.DeepCopyInto(&job.ObjectMeta)
		trial.Spec.Template.Spec.DeepCopyInto(&job.Spec)
	}

	// Provide default metadata
	if job.Name == "" {
		job.Name = trial.Name
	}
	if job.Namespace == "" {
		job.Namespace = trial.Namespace
	}

	// Provide default labels
	if len(job.Labels) == 0 {
		job.Labels = trial.GetDefaultLabels()
	}

	// TODO Also add the "trial" label to the pod template?

	// The default restart policy for a pod is not acceptable in the context of a job
	if job.Spec.Template.Spec.RestartPolicy == "" {
		job.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
	}

	// Containers cannot be empty, inject a sleep by default
	if len(job.Spec.Template.Spec.Containers) == 0 {
		s := trial.Spec.ApproximateRuntime
		if s == nil || s.Duration == 0 {
			s = &metav1.Duration{Duration: 2 * time.Minute}
		}
		if trial.Spec.StartTimeOffset != nil {
			s = &metav1.Duration{Duration: s.Duration + trial.Spec.StartTimeOffset.Duration}
		}
		job.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name:    "default-trial-run",
				Image:   "busybox",
				Command: []string{"/bin/sh"},
				Args:    []string{"-c", fmt.Sprintf("echo 'Sleeping for %s...' && sleep %.0f && echo 'Done.'", s.Duration.String(), s.Seconds())},
			},
		}
	}

	return job
}
