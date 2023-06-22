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

package kube

import (
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
	c      client.Client
	ww     client.WithWatch
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func Client(flags *genericclioptions.ConfigFlags, rbFlags *genericclioptions.ResourceBuilderFlags) (client.Client, error) {
	if c != nil {
		return c, nil
	}

	kubeconfig, err := flags.ToRESTConfig()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	kubeclient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if rbFlags != nil {
		if rbFlags.AllNamespaces == nil || !*rbFlags.AllNamespaces {
			if flags.Namespace != nil && len(*flags.Namespace) > 0 {
				log.Printf("Running with namespace %s.\n", *flags.Namespace)
				kubeclient = client.NewNamespacedClient(kubeclient, *flags.Namespace)
			} else {
				log.Println("Running with namespace default.")
				kubeclient = client.NewNamespacedClient(kubeclient, "default")
			}
		} else {
			log.Println("Running with all namespaces.")
			rbFlags.AllNamespaces = pointer.Bool(true)
		}
	}
	c = kubeclient
	return c, nil
}

func ClientWithWatch(flags *genericclioptions.ConfigFlags) (client.WithWatch, error) {
	if ww != nil {
		return ww, nil
	}

	kubeconfig, err := flags.ToRESTConfig()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	kubeclient, err := client.NewWithWatch(kubeconfig, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ww = kubeclient
	return ww, nil
}
