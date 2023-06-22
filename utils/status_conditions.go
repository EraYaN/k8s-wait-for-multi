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

package utils

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// Adapted from https://github.com/kubernetes/apimachinery/blob/master/pkg/api/meta/conditions.go

// IsJobStatusConditionTrue returns true when the conditionType is present and set to `metav1.ConditionTrue`
func IsJobStatusConditionTrue(conditions []batchv1.JobCondition, conditionType batchv1.JobConditionType) bool {
	return IsJobStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionTrue)
}

// IsJobStatusConditionFalse returns true when the conditionType is present and set to `metav1.ConditionFalse`
func IsJobStatusConditionFalse(conditions []batchv1.JobCondition, conditionType batchv1.JobConditionType) bool {
	return IsJobStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionFalse)
}

// IsJobStatusConditionPresentAndEqual returns true when conditionType is present and equal to status.
func IsJobStatusConditionPresentAndEqual(conditions []batchv1.JobCondition, conditionType batchv1.JobConditionType, status corev1.ConditionStatus) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition.Status == status
		}
	}
	return false
}

// IsPodStatusConditionTrue returns true when the conditionType is present and set to `metav1.ConditionTrue`
func IsPodStatusConditionTrue(conditions []corev1.PodCondition, conditionType corev1.PodConditionType) bool {
	return IsPodStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionTrue)
}

// IsPodStatusConditionFalse returns true when the conditionType is present and set to `metav1.ConditionFalse`
func IsPodStatusConditionFalse(conditions []corev1.PodCondition, conditionType corev1.PodConditionType) bool {
	return IsPodStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionFalse)
}

// IsPodStatusConditionPresentAndEqual returns true when conditionType is present and equal to status.
func IsPodStatusConditionPresentAndEqual(conditions []corev1.PodCondition, conditionType corev1.PodConditionType, status corev1.ConditionStatus) bool {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition.Status == status
		}
	}
	return false
}
