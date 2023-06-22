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

package items

import (
	"k8s.io/kubectl/pkg/util/podutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespacedPodCollection map[string]PodCollection

type PodCollection map[string]*PodItem

type PodItem struct {
	namespace string
	name      string
	ready     bool
}

func Pod(ns string, n string) *PodItem {
	return &PodItem{
		namespace: ns,
		name:      n,
		ready:     false,
	}
}

func (i *PodItem) WithReady(ready bool) *PodItem {
	i.ready = ready
	return i
}

func (i *PodItem) WithReadyFromPod(pod *corev1.Pod) *PodItem {
	i.ready = podutils.IsPodReady(pod) && podutils.IsPodAvailable(pod, 2, metav1.Now())
	return i
}

func (i *PodItem) GetName() string {
	return i.name
}

func (i *PodItem) GetNamespace() string {
	return i.namespace
}

func (i *PodItem) IsReady() bool {
	return i.ready
}

func (c NamespacedPodCollection) EnsureNamespace(ns string) {
	if _, ok := c[ns]; !ok {
		c[ns] = PodCollection{}
	}
}

func (c PodCollection) Contains(i ItemInterface) bool {
	return c.ContainsName(i.GetName())
}

func (c NamespacedPodCollection) Contains(i ItemInterface) bool {
	return c.ContainsNamespacedName(i.GetNamespace(), i.GetName())
}

func (c PodCollection) ContainsName(n string) bool {
	_, ok := c[n]
	return ok
}

func (c NamespacedPodCollection) ContainsNamespacedName(ns string, n string) bool {
	_, ok := c[ns][n]
	return ok
}

func (c NamespacedPodCollection) TotalCount() int {
	count := 0
	for _, items := range c {
		count += len(items)
	}
	return count
}

func (c NamespacedPodCollection) AreAllReady() bool {
	for _, items := range c {
		for _, item := range items {
			if !item.ready {
				return false
			}
		}
	}
	return true
}
