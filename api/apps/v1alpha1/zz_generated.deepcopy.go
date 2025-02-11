// +build !ignore_autogenerated

/*
Copyright 2021 GramLabs, Inc.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/thestormforge/konjure/pkg/konjure"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Application) DeepCopyInto(out *Application) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make(konjure.Resources, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Configuration != nil {
		in, out := &in.Configuration, &out.Configuration
		*out = make([]Parameter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(Ingress)
		**out = **in
	}
	if in.Scenarios != nil {
		in, out := &in.Scenarios, &out.Scenarios
		*out = make([]Scenario, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Objectives != nil {
		in, out := &in.Objectives, &out.Objectives
		*out = make([]Objective, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Application.
func (in *Application) DeepCopy() *Application {
	if in == nil {
		return nil
	}
	out := new(Application)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Application) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerResources) DeepCopyInto(out *ContainerResources) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]v1.ResourceName, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerResources.
func (in *ContainerResources) DeepCopy() *ContainerResources {
	if in == nil {
		return nil
	}
	out := new(ContainerResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomScenario) DeepCopyInto(out *CustomScenario) {
	*out = *in
	if in.PodTemplate != nil {
		in, out := &in.PodTemplate, &out.PodTemplate
		*out = new(v1.PodTemplateSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomScenario.
func (in *CustomScenario) DeepCopy() *CustomScenario {
	if in == nil {
		return nil
	}
	out := new(CustomScenario)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatadogGoal) DeepCopyInto(out *DatadogGoal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatadogGoal.
func (in *DatadogGoal) DeepCopy() *DatadogGoal {
	if in == nil {
		return nil
	}
	out := new(DatadogGoal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DurationGoal) DeepCopyInto(out *DurationGoal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DurationGoal.
func (in *DurationGoal) DeepCopy() *DurationGoal {
	if in == nil {
		return nil
	}
	out := new(DurationGoal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariable) DeepCopyInto(out *EnvironmentVariable) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariable.
func (in *EnvironmentVariable) DeepCopy() *EnvironmentVariable {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ErrorRateGoal) DeepCopyInto(out *ErrorRateGoal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ErrorRateGoal.
func (in *ErrorRateGoal) DeepCopy() *ErrorRateGoal {
	if in == nil {
		return nil
	}
	out := new(ErrorRateGoal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Goal) DeepCopyInto(out *Goal) {
	*out = *in
	if in.Max != nil {
		in, out := &in.Max, &out.Max
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Min != nil {
		in, out := &in.Min, &out.Min
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Optimize != nil {
		in, out := &in.Optimize, &out.Optimize
		*out = new(bool)
		**out = **in
	}
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = new(RequestsGoal)
		(*in).DeepCopyInto(*out)
	}
	if in.Latency != nil {
		in, out := &in.Latency, &out.Latency
		*out = new(LatencyGoal)
		**out = **in
	}
	if in.ErrorRate != nil {
		in, out := &in.ErrorRate, &out.ErrorRate
		*out = new(ErrorRateGoal)
		**out = **in
	}
	if in.Duration != nil {
		in, out := &in.Duration, &out.Duration
		*out = new(DurationGoal)
		**out = **in
	}
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(PrometheusGoal)
		**out = **in
	}
	if in.Datadog != nil {
		in, out := &in.Datadog, &out.Datadog
		*out = new(DatadogGoal)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Goal.
func (in *Goal) DeepCopy() *Goal {
	if in == nil {
		return nil
	}
	out := new(Goal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LatencyGoal) DeepCopyInto(out *LatencyGoal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LatencyGoal.
func (in *LatencyGoal) DeepCopy() *LatencyGoal {
	if in == nil {
		return nil
	}
	out := new(LatencyGoal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocustScenario) DeepCopyInto(out *LocustScenario) {
	*out = *in
	if in.Users != nil {
		in, out := &in.Users, &out.Users
		*out = new(int)
		**out = **in
	}
	if in.SpawnRate != nil {
		in, out := &in.SpawnRate, &out.SpawnRate
		*out = new(int)
		**out = **in
	}
	if in.RunTime != nil {
		in, out := &in.RunTime, &out.RunTime
		*out = new(metav1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocustScenario.
func (in *LocustScenario) DeepCopy() *LocustScenario {
	if in == nil {
		return nil
	}
	out := new(LocustScenario)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Objective) DeepCopyInto(out *Objective) {
	*out = *in
	if in.Goals != nil {
		in, out := &in.Goals, &out.Goals
		*out = make([]Goal, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Objective.
func (in *Objective) DeepCopy() *Objective {
	if in == nil {
		return nil
	}
	out := new(Objective)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parameter) DeepCopyInto(out *Parameter) {
	*out = *in
	if in.ContainerResources != nil {
		in, out := &in.ContainerResources, &out.ContainerResources
		*out = new(ContainerResources)
		(*in).DeepCopyInto(*out)
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(Replicas)
		**out = **in
	}
	if in.EnvironmentVariable != nil {
		in, out := &in.EnvironmentVariable, &out.EnvironmentVariable
		*out = new(EnvironmentVariable)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parameter.
func (in *Parameter) DeepCopy() *Parameter {
	if in == nil {
		return nil
	}
	out := new(Parameter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusGoal) DeepCopyInto(out *PrometheusGoal) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusGoal.
func (in *PrometheusGoal) DeepCopy() *PrometheusGoal {
	if in == nil {
		return nil
	}
	out := new(PrometheusGoal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Replicas) DeepCopyInto(out *Replicas) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Replicas.
func (in *Replicas) DeepCopy() *Replicas {
	if in == nil {
		return nil
	}
	out := new(Replicas)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequestsGoal) DeepCopyInto(out *RequestsGoal) {
	*out = *in
	if in.Weights != nil {
		in, out := &in.Weights, &out.Weights
		*out = make(v1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequestsGoal.
func (in *RequestsGoal) DeepCopy() *RequestsGoal {
	if in == nil {
		return nil
	}
	out := new(RequestsGoal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scenario) DeepCopyInto(out *Scenario) {
	*out = *in
	if in.StormForge != nil {
		in, out := &in.StormForge, &out.StormForge
		*out = new(StormForgeScenario)
		**out = **in
	}
	if in.Locust != nil {
		in, out := &in.Locust, &out.Locust
		*out = new(LocustScenario)
		(*in).DeepCopyInto(*out)
	}
	if in.Custom != nil {
		in, out := &in.Custom, &out.Custom
		*out = new(CustomScenario)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scenario.
func (in *Scenario) DeepCopy() *Scenario {
	if in == nil {
		return nil
	}
	out := new(Scenario)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StormForgeScenario) DeepCopyInto(out *StormForgeScenario) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StormForgeScenario.
func (in *StormForgeScenario) DeepCopy() *StormForgeScenario {
	if in == nil {
		return nil
	}
	out := new(StormForgeScenario)
	in.DeepCopyInto(out)
	return out
}
