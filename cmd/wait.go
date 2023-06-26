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
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/erayan/k8s-wait-for-multi/pkg"
	"github.com/erayan/k8s-wait-for-multi/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	toolscache "k8s.io/client-go/tools/cache"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
)

var cancelFn context.CancelFunc
var timeoutCtx context.Context
var cc cache.Cache
var mu sync.Mutex

var waits *pkg.Waitables

func wait(cmd *cobra.Command, args []string) error {
	if *WaitForConfigFlags.PrintVersion {
		return printVersion(cmd, args)
	}

	if len(args) < 1 {
		return errors.New("command needs one or more arguments to wait for")
	}

	if KubernetesConfigFlags.Namespace == nil || *KubernetesConfigFlags.Namespace == "" {
		KubernetesConfigFlags.Namespace = pointer.String("default")
	}

	waits = pkg.NewWaitables(WaitForConfigFlags)

	waits.Start()

	timeout, err := cmd.Flags().GetDuration("timeout")

	if err != nil {
		return err
	}

	timeoutCtx, cancelFn = context.WithTimeout(context.Background(), timeout)
	defer cancelFn()

	namespaces := []string{}

	for _, arg := range args {
		arg_items := strings.Split(arg, ",")
		len := len(arg_items)
		err = nil
		if len == 1 || len == 2 {
			if !utils.StringInSlice(*KubernetesConfigFlags.Namespace, namespaces) {
				namespaces = append(namespaces, *KubernetesConfigFlags.Namespace)
			}
		} else if len == 3 {
			if !utils.StringInSlice(arg_items[0], namespaces) {
				namespaces = append(namespaces, arg_items[0])
			}
		} else {
			log.Printf("illegal argument '%s'", arg)
		}
	}

	log.Printf("Starting with namespaces: %v", namespaces)

	opts := cache.Options{
		Namespaces: namespaces,
		SyncPeriod: WaitForConfigFlags.SyncPeriod,
	}

	conf, err := KubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	cc, err = cache.New(conf, opts)
	if err != nil {
		return err
	}

	mu = sync.Mutex{}
	waits.WithCache(cc)

	illegals := false

	for _, arg := range args {
		arg_items := strings.Split(arg, ",")
		len := len(arg_items)
		err = nil
		if len == 1 {
			err = waits.AddItem("pod", *KubernetesConfigFlags.Namespace, arg_items[0])
			illegals = err != nil
		} else if len == 2 {
			err = waits.AddItem(arg_items[0], *KubernetesConfigFlags.Namespace, arg_items[1])
			illegals = err != nil
		} else if len == 3 {
			err = waits.AddItem(arg_items[1], arg_items[0], arg_items[2])
			illegals = err != nil
		} else {
			log.Printf("illegal argument '%s'", arg)
			illegals = true
		}
		if err != nil {
			log.Printf("illegal argument '%s': %s", arg, err.Error())
		}
	}

	if illegals {
		return errors.New("illegal argument provided")
	}

	if waits.TotalCount() < 1 {
		return errors.New("not enough arguments")
	}

	waits.PrintStatus()

	if waits.HasServices() {
		svc_informer, err := cc.GetInformerForKind(timeoutCtx, schema.FromAPIVersionAndKind("v1", "Service"))
		if err != nil {
			return err
		}

		svc_informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventAddService, obj.(*corev1.Service))
			},
			UpdateFunc: func(obj interface{}, newObj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventUpdateService, newObj.(*corev1.Service))
			},
			DeleteFunc: func(obj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventDeleteService, obj.(*corev1.Service))
			},
		})
	}

	if waits.HasServices() || waits.HasPods() {
		pod_informer, err := cc.GetInformerForKind(timeoutCtx, schema.FromAPIVersionAndKind("v1", "Pod"))
		if err != nil {
			return err
		}

		pod_informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventAddPod, obj.(*corev1.Pod))
			},
			UpdateFunc: func(obj interface{}, newObj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventUpdatePod, newObj.(*corev1.Pod))
			},
			DeleteFunc: func(obj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventDeletePod, obj.(*corev1.Pod))
			},
		})
	}

	if waits.HasJobs() {
		job_informer, err := cc.GetInformerForKind(timeoutCtx, schema.FromAPIVersionAndKind("batch/v1", "Job"))
		if err != nil {
			return err
		}

		job_informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventAddJob, obj.(*batchv1.Job))
			},
			UpdateFunc: func(obj interface{}, newObj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventUpdateJob, newObj.(*batchv1.Job))
			},
			DeleteFunc: func(obj interface{}) {
				handleEvent(timeoutCtx, waits.ProcessEventDeleteJob, obj.(*batchv1.Job))
			},
		})
	}

	err = cc.Start(timeoutCtx)
	if err != nil {
		return err
	}

	waits.Done()

	return nil
}

func processCompletion() {
	cancelFn()
}

func handleEvent[V *corev1.Pod | *corev1.Service | *batchv1.Job](ctx context.Context, f func(ctx context.Context, obj V) (bool, error), obj V) {
	mu.Lock()
	defer mu.Unlock()

	matches, err := f(ctx, obj)
	if err != nil {
		log.Fatal(err)
	}

	if matches {
		waits.PrintStatus()

		if waits.IsDone() {
			processCompletion()
		}
	}

}
