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

package pkg

import (
	"fmt"
	"strings"

	"github.com/erayan/k8s-wait-for-multi/pkg/items"
	"github.com/xlab/treeprint"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type Waitables struct {
	cache.Cache

	UnprocessablePodEvents map[types.UID]Event

	Services items.NamespacedServiceCollection
	Pods     items.NamespacedPodCollection
	Jobs     items.NamespacedJobCollection
}

func (w *Waitables) AddItem(kind string, namespace string, name string) error {
	if kind == "pod" {
		w.addPod(namespace, name)
	} else if kind == "job" {
		w.addJob(namespace, name)
	} else if kind == "service" {
		w.addService(namespace, name)
	} else {
		return fmt.Errorf("unsupported kind '%s'", kind)
	}
	return nil
}

func (w *Waitables) addPod(namespace string, name string) *items.PodItem {
	w.Pods.EnsureNamespace(namespace)
	if !w.Pods.ContainsNamespacedName(namespace, name) {
		w.Pods[namespace][name] = items.Pod(namespace, name)
	}
	return w.Pods[namespace][name]
}

func (w *Waitables) addService(namespace string, name string) *items.ServiceItem {
	w.Services.EnsureNamespace(namespace)
	if !w.Services.ContainsNamespacedName(namespace, name) {
		w.Services[namespace][name] = items.Service(namespace, name)
	}
	return w.Services[namespace][name]
}

func (w *Waitables) addJob(namespace string, name string) *items.JobItem {
	w.Jobs.EnsureNamespace(namespace)
	if !w.Jobs.ContainsNamespacedName(namespace, name) {
		w.Jobs[namespace][name] = items.Job(namespace, name)
	}
	return w.Jobs[namespace][name]
}

func (w *Waitables) HasPod(meta metav1.ObjectMeta) bool {
	return w.Pods.Contains(&meta)
}

func (w *Waitables) HasService(meta metav1.ObjectMeta) bool {
	return w.Services.Contains(&meta)
}

func (w *Waitables) HasJob(meta metav1.ObjectMeta) bool {
	return w.Jobs.Contains(&meta)
}

func (w *Waitables) HasPods() bool {
	return w.Pods.TotalCount() > 0
}

func (w *Waitables) HasServices() bool {
	return w.Services.TotalCount() > 0
}

func (w *Waitables) HasJobs() bool {
	return w.Jobs.TotalCount() > 0
}

func (w *Waitables) IsDone() bool {
	s := w.Services.AreAllAvailable()
	p := w.Pods.AreAllReady()
	j := w.Jobs.AreAllComplete()
	return s && p && j
}

func (w *Waitables) GetStatusString() string {
	items := []string{}
	for ns, nsitems := range w.Services {
		for n, val := range nsitems {
			if !val.IsAvailable() {
				items = append(items, fmt.Sprintf("%s/service/%s", ns, n))
			}
		}
	}
	for ns, nsitems := range w.Pods {
		for n, val := range nsitems {
			if !val.IsReady() {
				items = append(items, fmt.Sprintf("%s/pod/%s", ns, n))
			}
		}
	}
	for ns, nsitems := range w.Jobs {
		for n, val := range nsitems {
			if !val.IsComplete() {
				items = append(items, fmt.Sprintf("%s/job/%s", ns, n))
			}
		}
	}
	return fmt.Sprintf("Waiting for: %s", strings.Join(items, ", "))
}

func (w *Waitables) GetStatusTreeString(collapseTree bool) string {

	tree := treeprint.NewWithRoot("wait status")

	namespace_branches := map[string]treeprint.Tree{}

	for _, ns := range w.GetAllNamespaces() {
		namespace_branches[ns] = tree.AddMetaBranch(TreeStatusUnknown, fmt.Sprintf("namespace/%s", ns))
	}

	for ns, nsitems := range w.Services {
		branch := namespace_branches[ns]
		for n, val := range nsitems {
			status := "Unavailable"
			if val.IsAvailable() {
				status = "Available"
			}
			svc_branch := branch.AddMetaBranch(TreeStatusUnknown, fmt.Sprintf("service/%s: %s", n, status))

			for podname, pod := range *val.GetChildren() {
				status := "NotReady"
				meta := TreeStatusNotDone
				if pod.IsReady() {
					status = "Ready"
					meta = TreeStatusDone
				}
				svc_branch.AddMetaNode(meta, fmt.Sprintf("pod/%s: %s", podname, status))
			}
		}
	}
	for ns, nsitems := range w.Pods {
		branch := namespace_branches[ns]
		for n, val := range nsitems {
			status := "NotReady"
			meta := TreeStatusNotDone
			if val.IsReady() {
				status = "Ready"
				meta = TreeStatusDone
			}
			branch.AddMetaNode(meta, fmt.Sprintf("pod/%s: %s", n, status))
		}
	}
	for ns, nsitems := range w.Jobs {
		branch := namespace_branches[ns]
		for n, val := range nsitems {
			status := "NotComplete"
			meta := TreeStatusNotDone
			if val.IsComplete() {
				status = "Complete"
				meta = TreeStatusDone
			}
			branch.AddMetaNode(meta, fmt.Sprintf("job/%s: %s", n, status))
		}
	}

	// if you need to iterate over the whole tree
	// call `VisitAll` from your top root node.
	tree.VisitAll(func(item *treeprint.Node) {
		GetNodesStatus(item, collapseTree)
	})

	return tree.String()
}

func GetNodesStatus(item *treeprint.Node, collapseTree bool) treeprint.MetaValue {
	if len(item.Nodes) > 0 && item.Meta == TreeStatusUnknown {
		status := []treeprint.MetaValue{}

		for _, node := range item.Nodes {
			status = append(status, GetNodesStatus(node, collapseTree))
		}
		if metaInSlice(TreeStatusUnknown, status) {
			item.Meta = TreeStatusUnknown
		} else if metaInSlice(TreeStatusNotDone, status) {
			item.Meta = TreeStatusNotDone
		} else {
			if collapseTree {
				item.Nodes = nil
			}
			item.Meta = TreeStatusDone
		}
		// branch nodes
	}

	return item.Meta
}

func (w *Waitables) SetPodReadyFromPod(pod *corev1.Pod) {
	w.Pods[pod.Namespace][pod.Name].WithReadyFromPod(pod)
}

func (w *Waitables) SetPodReady(pod *corev1.Pod) {
	w.Pods[pod.Namespace][pod.Name].WithReady(true)
}

func (w *Waitables) UnsetPodReady(pod *corev1.Pod) {
	w.Pods[pod.Namespace][pod.Name].WithReady(false)
}

func (w *Waitables) SetJobCompleteFromJob(job *batchv1.Job) {
	w.Jobs[job.Namespace][job.Name].WithCompleteFromJob(job)
}

func (w *Waitables) SetJobComplete(job *batchv1.Job) {
	w.Jobs[job.Namespace][job.Name].WithComplete(true)
}

func (w *Waitables) UnsetJobComplete(job *batchv1.Job) {
	w.Jobs[job.Namespace][job.Name].WithComplete(false)
}

func (w *Waitables) SetServiceChildren(meta *metav1.ObjectMeta, pods []corev1.Pod) {
	podItems := items.PodCollection{}
	for _, pod := range pods {
		podItems[pod.Name] = items.Pod(pod.Namespace, pod.Name).WithReadyFromPod(&pod)
	}
	w.Services[meta.Namespace][meta.Name].WithChildren(podItems)
}

func (w *Waitables) TotalCount() int {
	return w.Services.TotalCount() + w.Pods.TotalCount() + w.Jobs.TotalCount()
}

func (w *Waitables) GetAllNamespaces() []string {
	namespaces := []string{}

	for ns := range w.Services {
		if !stringInSlice(ns, namespaces) {
			namespaces = append(namespaces, ns)
		}
	}

	for ns := range w.Jobs {
		if !stringInSlice(ns, namespaces) {
			namespaces = append(namespaces, ns)
		}
	}
	for ns := range w.Pods {
		if !stringInSlice(ns, namespaces) {
			namespaces = append(namespaces, ns)
		}
	}

	return namespaces
}

func (w *Waitables) WithCache(c cache.Cache) *Waitables {
	w.Cache = c
	return w
}

func NewWaitables() *Waitables {
	return &Waitables{
		UnprocessablePodEvents: map[types.UID]Event{},
		Services:               items.NamespacedServiceCollection{},
		Pods:                   items.NamespacedPodCollection{},
		Jobs:                   items.NamespacedJobCollection{},
	}
}
