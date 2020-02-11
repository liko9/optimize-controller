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

package util

import (
	"github.com/redskyops/redskyops-controller/internal/config"
	redskyclient "github.com/redskyops/redskyops-controller/redskyapi"
	"github.com/spf13/pflag"
)

// Red Sky server specific configuration flags

type ServerFlags struct {
}

func NewServerFlags() *ServerFlags {
	return &ServerFlags{}
}

func (f *ServerFlags) AddFlags(*pflag.FlagSet) {}

func (f *ServerFlags) ToClientConfig() (redskyclient.Config, error) {
	cfg := &config.RedSkyConfig{}
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	return cfg, nil
}
