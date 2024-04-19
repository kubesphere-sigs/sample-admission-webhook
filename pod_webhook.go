/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sample_admission_webhook

import (
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var log = logf.Log.WithName("sample_admission_webhook")

type Webhook struct{}

func (r *Webhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	defeater := admission.WithCustomDefaulter(scheme.Scheme, &corev1.Pod{}, r)
	validator := admission.WithCustomValidator(scheme.Scheme, &corev1.Pod{}, r)
	mgr.GetWebhookServer().Register("/mutate-sample-admission-webhook-v1-pod", defeater)
	mgr.GetWebhookServer().Register("/validate-sample-admission-webhook-v1-pod", validator)

	//ctrl.NewWebhookManagedBy(mgr).WithDefaulter(r).WithValidator(r).Complete()
	return nil
}

// +kubebuilder:webhook:path=/mutate-sample-admission-webhook-v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups=core,resources=pods,verbs=create;update,versions=v1,name=defaulter.sample-admission-webhook,admissionReviewVersions=v1
var _ webhook.CustomDefaulter = &Webhook{}

func (r *Webhook) Default(ctx context.Context, obj runtime.Object) error {
	pod := obj.(*corev1.Pod)

	log.Info("default", "name", pod.Name)
	return nil
}

// +kubebuilder:webhook:path=/validate-sample-admission-webhook-v1-pod,mutating=false,failurePolicy=fail,sideEffects=None,groups=core,resources=pods,verbs=create;update,versions=v1,name=validator.sample-admission-webhook,admissionReviewVersions=v1
var _ webhook.CustomValidator = &Webhook{}

func (r *Webhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	pod := obj.(*corev1.Pod)

	if !allContainerImageIsAllowed(pod) {
		return nil, errors.New("nginx:latest is not allowed")
	}

	log.Info("validate create", "name", pod.Name)
	return nil, nil
}

func allContainerImageIsAllowed(pod *corev1.Pod) bool {
	for _, container := range pod.Spec.Containers {
		if container.Image == "nginx:latest" {
			return false
		}
	}
	return true
}

func (r *Webhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	newPod := newObj.(*corev1.Pod)
	if !allContainerImageIsAllowed(newPod) {
		return nil, errors.New("nginx:latest is not allowed")
	}
	log.Info("validate update", "name", newPod.Name)

	return nil, nil
}

func (r *Webhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	pod := obj.(*corev1.Pod)

	log.Info("validate delete", "name", pod.Name)
	return nil, nil
}
