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
	"github.com/erayan/k8s-wait-for-multi/utils"

	batchv1 "k8s.io/api/batch/v1"
)

type NamespacedJobCollection map[string]JobCollection

type JobCollection map[string]*JobItem

type JobItem struct {
	namespace string
	name      string
	complete  bool
}

func Job(ns string, n string) *JobItem {
	return &JobItem{
		namespace: ns,
		name:      n,
		complete:  false,
	}
}

func (i *JobItem) WithComplete(complete bool) *JobItem {
	i.complete = complete
	return i
}

func (i *JobItem) WithCompleteFromJob(job *batchv1.Job) *JobItem {
	i.complete = utils.IsJobStatusConditionTrue(job.Status.Conditions, batchv1.JobComplete)
	return i
}

func (i *JobItem) GetName() string {
	return i.name
}

func (i *JobItem) GetNamespace() string {
	return i.namespace
}

func (i *JobItem) IsComplete() bool {
	return i.complete
}

func (c NamespacedJobCollection) EnsureNamespace(ns string) {
	if _, ok := c[ns]; !ok {
		c[ns] = JobCollection{}
	}
}

func (c NamespacedJobCollection) Contains(i ItemInterface) bool {
	_, ok := c[i.GetNamespace()][i.GetName()]
	return ok
}

func (c NamespacedJobCollection) ContainsNamespacedName(ns string, n string) bool {
	_, ok := c[ns][n]
	return ok
}

func (c NamespacedJobCollection) TotalCount() int {
	count := 0
	for _, items := range c {
		count += len(items)
	}
	return count
}

func (c NamespacedJobCollection) AreAllComplete() bool {
	for _, items := range c {
		for _, item := range items {
			if !item.complete {
				return false
			}
		}
	}
	return true
}
