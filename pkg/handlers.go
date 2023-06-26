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
	"context"
	"log"

	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (w *Waitables) ProcessEventAddService(ctx context.Context, svc *corev1.Service) (bool, error) {
	if w.HasService(svc.ObjectMeta) {
		//log.Printf("Add %T %s %s", svc, svc.Namespace, svc.Name)
		pods, err := w.getPodsForSvc(ctx, svc)

		if err != nil {
			return true, err
		}

		w.SetServiceChildren(&svc.ObjectMeta, pods.Items)
		return true, nil
	}
	return false, nil
}

func (w *Waitables) ProcessEventUpdateService(ctx context.Context, svc *corev1.Service) (bool, error) {
	if w.HasService(svc.ObjectMeta) {
		//log.Printf("Update %T %s %s", svc, svc.Namespace, svc.Name)
		pods, err := w.getPodsForSvc(ctx, svc)

		if err != nil {
			return true, err
		}

		w.SetServiceChildren(&svc.ObjectMeta, pods.Items)
		return true, nil
	}
	return false, nil
}

func (w *Waitables) ProcessEventDeleteService(ctx context.Context, svc *corev1.Service) (bool, error) {
	if w.HasService(svc.ObjectMeta) {
		//log.Printf("Delete %T %s %s", svc, svc.Namespace, svc.Name)

		w.SetServiceChildren(&svc.ObjectMeta, nil)
		return true, nil
	}
	return false, nil
}

func (w *Waitables) ProcessOldPodEvents(ctx context.Context, pod *corev1.Pod) (bool, error) {
	if val, ok := w.LastPodEvents[pod.UID]; ok {
		//log.Printf("Running LastPodEvents for %s/%s of type %v", pod.Namespace, pod.Name, val.EventType)
		//defer delete(w.LastPodEvents, pod.UID)
		if val.EventType == EventTypeAdd {
			return w.ProcessEventAddPod(ctx, pod)
		} else if val.EventType == EventTypeUpdate {
			return w.ProcessEventUpdatePod(ctx, pod)
		} else if val.EventType == EventTypeDelete {
			return w.ProcessEventDeletePod(ctx, pod)
		}
	}
	return false, nil
}

func (w *Waitables) ProcessEventAddPod(ctx context.Context, pod *corev1.Pod) (bool, error) {
	// if w.HasPod(pod.ObjectMeta) {
	// 	log.Printf("Add %T %s %s", pod, pod.Namespace, pod.Name)
	// }

	if w.HasPodDirect(pod.ObjectMeta) {
		w.SetPodReadyFromPod(pod)
	}

	if podItems, ok := w.Services.GetPods(pod); ok {
		for _, podItem := range podItems {
			podItem.WithReadyFromPod(pod)
		}
	}

	w.LastPodEvents[pod.UID] = Event{EventType: EventTypeAdd, Pod: pod}

	return w.HasPod(pod.ObjectMeta), nil
}

func (w *Waitables) ProcessEventUpdatePod(ctx context.Context, pod *corev1.Pod) (bool, error) {
	// if w.HasPod(pod.ObjectMeta) {
	// 	log.Printf("Update %T %s %s", pod, pod.Namespace, pod.Name)
	// }

	if w.HasPodDirect(pod.ObjectMeta) {
		w.SetPodReadyFromPod(pod)
	}

	if podItems, ok := w.Services.GetPods(pod); ok {
		for _, podItem := range podItems {
			podItem.WithReadyFromPod(pod)
		}
	}

	w.LastPodEvents[pod.UID] = Event{EventType: EventTypeUpdate, Pod: pod}

	return w.HasPod(pod.ObjectMeta), nil
}

func (w *Waitables) ProcessEventDeletePod(ctx context.Context, pod *corev1.Pod) (bool, error) {
	// if w.HasPod(pod.ObjectMeta) {
	// 	log.Printf("Delete %T %s %s", pod, pod.Namespace, pod.Name)
	// }

	if w.HasPodDirect(pod.ObjectMeta) {
		w.UnsetPodReady(pod)
	}

	if podItems, ok := w.Services.GetPods(pod); ok {
		for _, podItem := range podItems {
			podItem.WithReady(false)
		}
		w.Services.DeletePod(&pod.ObjectMeta)
	}

	w.LastPodEvents[pod.UID] = Event{EventType: EventTypeDelete, Pod: pod}

	return w.HasPod(pod.ObjectMeta), nil
}

func (w *Waitables) ProcessEventAddJob(ctx context.Context, job *batchv1.Job) (bool, error) {
	if w.HasJob(job.ObjectMeta) {
		//log.Printf("Add %T %s %s", job, job.Namespace, job.Name)
		w.SetJobCompleteFromJob(job)
		return true, nil
	}
	return false, nil
}

func (w *Waitables) ProcessEventUpdateJob(ctx context.Context, job *batchv1.Job) (bool, error) {
	if w.HasJob(job.ObjectMeta) {
		//log.Printf("Update %T %s %s", job, job.Namespace, job.Name)
		w.SetJobCompleteFromJob(job)
		return true, nil
	}
	return false, nil
}

func (w *Waitables) ProcessEventDeleteJob(ctx context.Context, job *batchv1.Job) (bool, error) {
	if w.HasJob(job.ObjectMeta) {
		//log.Printf("Delete %T %s %s", job, job.Namespace, job.Name)
		w.UnsetJobComplete(job)
		return true, nil
	}
	return false, nil
}

func (w *Waitables) getPodsForSvc(context context.Context, svc *corev1.Service) (*corev1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := &client.ListOptions{Namespace: svc.Namespace, LabelSelector: set.AsSelector()}
	pods := &corev1.PodList{}
	err := w.List(context, pods, listOptions)
	if err != nil {
		return nil, err
	}
	return pods, err
}

func (w *Waitables) printRolloutStatus(pod *corev1.Pod) error {
	log.Printf("Pod %s is %v", pod.Name, pod.Status.Phase)
	return nil
}
