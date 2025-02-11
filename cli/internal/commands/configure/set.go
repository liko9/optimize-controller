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

package configure

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/thestormforge/optimize-controller/v2/cli/internal/commander"
	"github.com/thestormforge/optimize-go/pkg/config"
)

// SetOptions are the options for setting a configuration property to a new value
type SetOptions struct {
	// Config is the Optimize Configuration to view
	Config *config.OptimizeConfig

	// Key is the name of the property being set
	Key string
	// Value is the new value for the property
	Value string
}

// NewSetCommand creates a new command for setting a configuration property
func NewSetCommand(o *SetOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set NAME [VALUE]",
		Short: "Modify the configuration file",
		Long:  "Modify the Optimize Configuration file",

		Example: `# Add an environment variable to the controller
stormforge config set controller.default.env.FOOBAR example

# Set the controller memory
stormforge config set controller.default.resources.memory 512Mi`,

		Args: cobra.RangeArgs(1, 2),

		PreRun: func(cmd *cobra.Command, args []string) {
			o.Complete(args)
		},
		RunE: commander.WithoutArgsE(o.set),
	}

	return cmd
}

// Complete overwrites the options using from an argument slice
func (o *SetOptions) Complete(args []string) {
	if len(args) > 0 {
		o.Key = args[0]
	}
	if len(args) > 1 {
		o.Value = args[1]
	} else if strings.Contains(o.Key, "=") {
		s := strings.SplitN(o.Key, "=", 2)
		o.Key = s[0]
		o.Value = s[1]
	}
}

func (o *SetOptions) set() error {
	// Allow an explicitly empty key value to just save the configuration. This
	// is very much an edge case, e.g. a developer wants to explicitly save a
	// migrated configuration before making further changes.
	if o.Key == "" {
		return o.Config.Write()
	}

	if o.Value != "" {
		if err := o.Config.Update(config.SetProperty(o.Key, o.Value)); err != nil {
			return err
		}
	} else {
		if err := o.Config.Update(config.UnsetProperty(o.Key)); err != nil {
			return err
		}
	}

	if err := o.Config.Write(); err != nil {
		return err
	}

	return nil
}
