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

type NamespacedServiceCollection map[string]ServiceCollection

type ServiceCollection map[string]*ServiceItem

type ServiceItem struct {
	namespace string
	name      string
	children  PodCollection
}

func Service(ns string, n string) *ServiceItem {
	return &ServiceItem{
		namespace: ns,
		name:      n,
		children:  nil,
	}
}

func (i *ServiceItem) GetChildren() *PodCollection {
	return &i.children
}

func (i *ServiceItem) WithChildren(children PodCollection) *ServiceItem {
	i.children = children
	return i
}

func (i *ServiceItem) IsAvailable() bool {
	if i.children == nil || len(i.children) == 0 {
		return false
	}

	for _, item := range i.children {
		if !item.ready {
			return false
		}
	}
	return true
}

func (i *ServiceItem) IsAtLeastOneAvailable() bool {
	if i.children == nil || len(i.children) == 0 {
		return false
	}

	for _, item := range i.children {
		if item.ready {
			return true
		}
	}
	return false
}

func (i *ServiceItem) GetName() string {
	return i.name
}

func (i *ServiceItem) GetNamespace() string {
	return i.namespace
}

func (i *ServiceItem) GetPod(pod ItemInterface) (*PodItem, bool) {
	if i.children == nil {
		return nil, false
	}
	val, ok := i.children[pod.GetName()]
	return val, ok
}

func (c ServiceItem) DeletePod(i ItemInterface) {
	if c.namespace == i.GetNamespace() {
		delete(c.children, i.GetName())
	}
}

func (c ServiceCollection) DeletePod(i ItemInterface) {
	for _, svc := range c {
		svc.DeletePod(i)
	}
}

func (c NamespacedServiceCollection) DeletePod(i ItemInterface) {
	for ns, nssvcs := range c {
		if ns == i.GetNamespace() {
			nssvcs.DeletePod(i)
		}
	}
}

func (c NamespacedServiceCollection) EnsureNamespace(ns string) {
	if _, ok := c[ns]; !ok {
		c[ns] = ServiceCollection{}
	}
}

func (c NamespacedServiceCollection) Contains(i ItemInterface) bool {
	_, ok := c[i.GetNamespace()][i.GetName()]
	return ok
}

func (c NamespacedServiceCollection) ContainsPod(i ItemInterface) bool {
	for _, items := range c {
		for _, item := range items {
			if item.children.Contains(i) {
				return true
			}
		}
	}
	return false
}

func (c NamespacedServiceCollection) GetPods(i ItemInterface) ([]*PodItem, bool) {

	pods := []*PodItem{}

	for _, items := range c {
		for _, item := range items {
			if val, ok := item.GetPod(i); ok {
				pods = append(pods, val)
			}
		}
	}

	return pods, len(pods) > 0
}

func (c NamespacedServiceCollection) ContainsNamespacedName(ns string, n string) bool {
	_, ok := c[ns][n]
	return ok
}

func (c NamespacedServiceCollection) TotalCount() int {
	count := 0
	for _, items := range c {
		count += len(items)
	}
	return count
}

func (c NamespacedServiceCollection) AreAllAvailable(onlyOnePerServiceRequired bool) bool {
	for _, services := range c {
		for _, service := range services {
			if onlyOnePerServiceRequired {
				if !service.IsAtLeastOneAvailable() {
					return false
				}
			} else {
				if !service.IsAvailable() {
					return false
				}
			}
		}
	}
	return true
}
