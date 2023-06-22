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
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kubectlget "k8s.io/kubectl/pkg/cmd/get"
)

var (
	KubernetesConfigFlags    *genericclioptions.ConfigFlags
	KubeResourceBuilderFlags *genericclioptions.ResourceBuilderFlags
	KubernetesPrintFlags     *genericclioptions.PrintFlags
	KubernetesGetPrintFlags  *kubectlget.PrintFlags
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-wait-for-multi NAMESPACE,KIND,NAME [ NAMESPACE,KIND,NAME ]",
	Short: "This is an implementation of k8s-wait-for that allows you to wait for multiple items in one process.",
	Long: `k8s-wait-for-multi
This is an implementation of k8s-wait-for that allows you to wait for multiple items in one process.
This uses informers to get the status updates for all the items that this application is waiting for.

You can omit the NAMESPACE and KIND, they default to the value of the --namespace flag and pod respectively. Supported string for KIND are service, job and pod.`,
	RunE:    wait,
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)

	KubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())

	rootCmd.Flags().BoolP("version", "v", false, "Display version info")

	rootCmd.Flags().Bool("no-collapse", false, "Do not collapse the status tree for done subtrees")

	rootCmd.Flags().Bool("no-tree", false, "Do not print the status as a tree")

	rootCmd.PersistentFlags().DurationP("timeout", "t", time.Duration(600*time.Second), "The length of time to wait before ending watch, zero means never. Any other values should contain a corresponding time unit (e.g. 1s, 2m, 3h). Default is 10 minutes")

}
