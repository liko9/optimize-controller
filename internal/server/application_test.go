/*
Copyright 2020 GramLabs, Inc.

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

package server

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	optimizev1beta2 "github.com/thestormforge/optimize-controller/v2/api/v1beta2"
	"github.com/thestormforge/optimize-go/pkg/api"
	applications "github.com/thestormforge/optimize-go/pkg/api/applications/v2"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestClusterExperimentToAPITemplate(t *testing.T) {
	one := intstr.FromInt(1)
	two := intstr.FromInt(2)
	three := intstr.FromString("three")

	exp := &optimizev1beta2.Experiment{
		Spec: optimizev1beta2.ExperimentSpec{
			Parameters: []optimizev1beta2.Parameter{
				{Name: "one", Min: 111, Max: 222, Baseline: &one},
				{Name: "two", Min: 1111, Max: 2222, Baseline: &two},
				{Name: "three", Min: 11111, Max: 22222, Baseline: &three},
				{Name: "test_case", Min: 1, Max: 1},
			},
			Metrics: []optimizev1beta2.Metric{
				{Name: "one", Minimize: true},
				{Name: "two", Minimize: false},
				{Name: "three", Minimize: true},
			},
		},
	}

	testCases := []struct {
		desc string
	}{
		{
			desc: "default",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.desc), func(t *testing.T) {
			template, err := ClusterExperimentToAPITemplate(exp)
			assert.NoError(t, err)
			assert.NotNil(t, template.Parameters)
			assert.NotNil(t, template.Metrics)

			// test_case is silently filters/dropped because min==max
			assert.Equal(t, 3, len(template.Parameters))
			assert.Equal(t, 3, len(template.Metrics))
		})
	}
}

func TestAPITemplateToClusterExperiment(t *testing.T) {
	templateOne := api.FromInt64(500)
	templateTwo := api.FromInt64(2000)
	templateThree := api.FromInt64(20000)
	expOne := intstr.FromString("500")
	expTwo := intstr.FromString("2000")
	expThree := intstr.FromString("20000")
	pTrue := true

	testCases := []struct {
		desc       string
		experiment *optimizev1beta2.Experiment
		template   *applications.Template
		expected   *optimizev1beta2.Experiment
	}{
		{
			desc: "params",
			experiment: &optimizev1beta2.Experiment{
				Spec: optimizev1beta2.ExperimentSpec{},
			},
			template: &applications.Template{
				Parameters: []applications.TemplateParameter{
					{
						Name: "one",
						Bounds: &applications.TemplateParameterBounds{
							Min: json.Number("1"),
							Max: json.Number("1000"),
						},
						Type:     "int",
						Baseline: &templateOne,
					},
					{
						Name: "two",
						Bounds: &applications.TemplateParameterBounds{
							Min: json.Number("1111"),
							Max: json.Number("2222"),
						},
						Type:     "int",
						Baseline: &templateTwo,
					},
					{
						Name: "three",
						Bounds: &applications.TemplateParameterBounds{
							Min: json.Number("11111"),
							Max: json.Number("22222"),
						},
						Type:     "int",
						Baseline: &templateThree,
					},
					// TODO this might be an edge case we need to handle
					// {Name: "test_case", Min: 1, Max: 1},
				},
			},
			expected: &optimizev1beta2.Experiment{
				Spec: optimizev1beta2.ExperimentSpec{
					Parameters: []optimizev1beta2.Parameter{
						{Name: "one", Min: 1, Max: 1000, Baseline: &expOne},
						{Name: "two", Min: 1111, Max: 2222, Baseline: &expTwo},
						{Name: "three", Min: 11111, Max: 22222, Baseline: &expThree},
					},
				},
			},
		},
		{
			desc: "metrics",
			experiment: &optimizev1beta2.Experiment{
				Spec: optimizev1beta2.ExperimentSpec{
					Metrics: []optimizev1beta2.Metric{
						{Name: "one", Minimize: true},
						{Name: "two", Minimize: false},
						{Name: "three", Minimize: true},
					},
				},
			},
			template: &applications.Template{
				Metrics: []applications.TemplateMetric{
					{
						Name:     "one",
						Minimize: true,
						Bounds: &applications.TemplateMetricBounds{
							Min: 1,
							Max: 10,
						},
					},
					{
						Name:     "two",
						Minimize: false,
						Bounds: &applications.TemplateMetricBounds{
							Min: 2,
							Max: 5,
						},
					},
					{
						Name:     "three",
						Minimize: false,
						Optimize: &pTrue,
					},
				},
			},
			expected: &optimizev1beta2.Experiment{
				Spec: optimizev1beta2.ExperimentSpec{
					Parameters: []optimizev1beta2.Parameter{},
					Metrics: []optimizev1beta2.Metric{
						{Name: "one", Minimize: true, Min: resource.NewQuantity(1, resource.DecimalSI), Max: resource.NewQuantity(10, resource.DecimalSI)},
						{Name: "two", Minimize: false, Min: resource.NewQuantity(2, resource.DecimalSI), Max: resource.NewQuantity(5, resource.DecimalSI)},
						{Name: "three", Minimize: false, Optimize: &pTrue},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.desc), func(t *testing.T) {
			err := APITemplateToClusterExperiment(tc.experiment, tc.template)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tc.experiment)
		})
	}
}
