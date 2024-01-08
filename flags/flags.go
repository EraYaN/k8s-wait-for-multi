/*
 *  Copyright 2023 The k8s-wait-for-multi authors.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  	http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package flags

import (
	"time"

	"github.com/spf13/pflag"

	utilpointer "k8s.io/utils/pointer"
)

type ConfigFlags struct {
	PrintVersion              *bool
	PrintTree                 *bool
	PrintCollapsedTree        *bool
	OnlyOnePerServiceRequired *bool

	Timeout    *time.Duration
	SyncPeriod *time.Duration
}

func NewConfigFlags() *ConfigFlags {
	return &ConfigFlags{
		PrintVersion:              utilpointer.Bool(false),
		PrintTree:                 utilpointer.Bool(true),
		PrintCollapsedTree:        utilpointer.Bool(true),
		OnlyOnePerServiceRequired: utilpointer.Bool(false),

		Timeout:    utilpointer.Duration(time.Duration(600 * time.Second)),
		SyncPeriod: utilpointer.Duration(time.Duration(90 * time.Second)),
	}
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.Timeout != nil {
		flags.DurationVarP(f.Timeout, "timeout", "t", *f.Timeout, "The length of time to wait before ending watch, zero means never. Any other values should contain a corresponding time unit (e.g. 1s, 2m, 3h)")
	}

	if f.SyncPeriod != nil {
		flags.DurationVar(f.SyncPeriod, "sync-period", *f.SyncPeriod, "The length of time to pass to the cache to initiate a sync. (e.g. 1s, 2m, 3h)")
	}

	if f.OnlyOnePerServiceRequired != nil {
		flags.BoolVar(f.OnlyOnePerServiceRequired, "only-one-per-service-required", *f.OnlyOnePerServiceRequired, "When true a service is ready when at least one pod is ready. When false all pods must be ready.")
	}

	if f.PrintVersion != nil {
		flags.BoolVarP(f.PrintVersion, "version", "v", *f.PrintVersion, "Display version info")
	}

	if f.PrintTree != nil {
		flags.BoolVar(f.PrintTree, "print-tree", *f.PrintTree, "Print the status as a tree")
	}

	if f.PrintCollapsedTree != nil {
		flags.BoolVar(f.PrintCollapsedTree, "print-collapsed-tree", *f.PrintCollapsedTree, "Collapse the status tree for done subtrees")
	}
}
