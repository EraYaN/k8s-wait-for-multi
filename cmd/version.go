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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

// SetVersion set the application version for consumption in the output of the command.
func SetVersionInfo(v string, c string, d string) {
	version = v
	commit = c
	date = d
}

func printVersion(cmd *cobra.Command, args []string) error {
	_, err := fmt.Printf("k8s-wait-for-multi v%s %s %s\n", version, commit, date)
	if err != nil {
		return err
	}
	return nil
}
