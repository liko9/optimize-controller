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

package controller

import (
	redskyapi "github.com/redskyops/redskyops-controller/redskyapi/experiments/v1alpha1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// These methods return a pointer to a Result struct and an error; this is useful for implementing the
// pattern where Reconcile methods are split up and check for a non-nil Result to determine if it is
// time to return.

// RequeueIfUnavailable will return a new result and the supplied error, adjusted for trial unavailable errors
func RequeueIfUnavailable(err error) (*ctrl.Result, error) {
	result := &ctrl.Result{}
	if rse, ok := err.(*redskyapi.Error); ok && rse.Type == redskyapi.ErrTrialUnavailable {
		result.RequeueAfter = rse.RetryAfter
		err = nil
	}
	return result, err
}

// RequeueConflict will return a new result and the supplied error, adjusted for Kubernetes conflict errors
func RequeueConflict(err error) (*ctrl.Result, error) {
	result := &ctrl.Result{}
	if apierrs.IsConflict(err) {
		ReconcileConflictErrors.WithLabelValues(guessController()).Inc()
		result.Requeue = true
		err = nil
	}
	return result, err
}
