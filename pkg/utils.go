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

import "github.com/xlab/treeprint"

const (
	TreeStatusDone    = "✅"
	TreeStatusIgnored = "☑️"
	TreeStatusNotDone = "❌"
	TreeStatusUnknown = "❔"
)

type TreeStatus string

func metaInSlice(a treeprint.MetaValue, list []treeprint.MetaValue) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getNodesStatus(item *treeprint.Node, collapseTree bool) treeprint.MetaValue {
	if len(item.Nodes) > 0 && item.Meta == TreeStatusUnknown {
		status := []treeprint.MetaValue{}

		for _, node := range item.Nodes {
			status = append(status, getNodesStatus(node, collapseTree))
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
