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

package get

import (
	"context"
	"path"
	"sort"
	"strings"

	redsky "github.com/redskyops/k8s-experiment/pkg/api/redsky/v1alpha1"
	"github.com/redskyops/k8s-experiment/pkg/controller/experiment"
	cmdutil "github.com/redskyops/k8s-experiment/pkg/redskyctl/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	getExperimentListLong    = `Prints a list of experiments using a tabular format by default`
	getExperimentListExample = ``
)

// NewGetExperimentListCommand returns a new get experiment list command
func NewGetExperimentListCommand(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewGetOptions(ioStreams)

	printFlags := cmdutil.NewPrintFlags(&experimentTableMeta{})

	cmd := &cobra.Command{
		Use:     "experiments",
		Short:   "Display a list of experiments",
		Long:    getExperimentListLong,
		Example: getExperimentListExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args, printFlags))
			cmdutil.CheckErr(RunGetExperimentList(o))
		},
	}

	o.AddFlags(cmd)
	printFlags.AddFlags(cmd)

	return cmd
}

// RunGetExperimentList gets an experiment list for the given get options
func RunGetExperimentList(o *GetOptions) error {
	if o.RedSkyAPI != nil {
		if err := o.printIf(getRedSkyAPIExperimentList(*o.RedSkyAPI, o.ChunkSize)); err != nil {
			return err
		}
	}

	if o.RedSkyClientSet != nil {
		if err := o.printIf(getKubernetesExperimentList(o)); err != nil {
			return err
		}
	}

	return nil
}

func getRedSkyAPIExperimentList(api redsky.API, chunkSize int) (*redsky.ExperimentList, error) {
	l, err := api.GetAllExperiments(context.TODO(), &redsky.ExperimentListQuery{Limit: chunkSize})
	if err != nil {
		return nil, err
	}

	n := l
	for n.Next != "" {
		if n, err = api.GetAllExperimentsByPage(context.TODO(), n.Next); err != nil {
			return nil, err
		}
		l.Experiments = append(l.Experiments, n.Experiments...)
	}

	return filterAndSortExperiments(&l, "", "")
}

func getKubernetesExperimentList(o *GetOptions) (*redsky.ExperimentList, error) {
	clientset := o.RedSkyClientSet
	experiments := clientset.RedskyopsV1alpha1().Experiments(o.Namespace)
	opts := metav1.ListOptions{Limit: int64(o.ChunkSize)}
	l := &redsky.ExperimentList{}
	for {
		el, err := experiments.List(opts)
		if err != nil {
			return nil, err
		}

		err = experiment.ConvertExperimentList(el, l)
		if err != nil {
			return nil, err
		}

		if opts.Continue == "" {
			break
		}
	}
	return filterAndSortExperiments(l, "", "")
}

func filterAndSortExperiments(el *redsky.ExperimentList, selector, sortBy string) (*redsky.ExperimentList, error) {
	// Experiments do not have labels so anything but the empty selector will just nil out the list
	if sel, err := labels.Parse(selector); err != nil {
		return nil, err
	} else if !sel.Empty() {
		el.Experiments = nil
	} else if sortBy != "" {
		// TODO Do we really want to expose the whole experiment?
		sort.Slice(el.Experiments, sortByField(sortBy, func(i int) interface{} { return &el.Experiments[i] }))
	}
	return el, nil
}

type experimentTableMeta struct{}

func (*experimentTableMeta) IsListType(obj interface{}) bool {
	if _, ok := obj.(*redsky.ExperimentList); ok {
		return true
	}
	return false
}

func (*experimentTableMeta) ExtractList(obj interface{}) ([]interface{}, error) {
	switch o := obj.(type) {
	case *redsky.ExperimentList:
		list := make([]interface{}, len(o.Experiments))
		for i := range o.Experiments {
			list[i] = &o.Experiments[i]
		}
		return list, nil
	default:
		return []interface{}{obj}, nil
	}
}

func (*experimentTableMeta) ExtractValue(obj interface{}, column string) (string, error) {
	switch o := obj.(type) {
	case *redsky.ExperimentItem:
		switch column {
		case "name":
			return path.Base(o.Self), nil
		}
	}
	return "", nil
}

func (*experimentTableMeta) Allow(outputFormat string) bool {
	return outputFormat == "" || strings.ToLower(outputFormat) == "name"
}

func (*experimentTableMeta) Columns(outputFormat string) []string {
	return []string{"name"}
}

func (*experimentTableMeta) Header(outputFormat string, column string) string {
	return strings.ToUpper(column)
}
