/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var log = logf.Log.WithName("awsmachinepool-resource")

// SetupWebhookWithManager will setup the webhooks for the AWSMachinePool
func (r *AWSMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-awsmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,versions=v1alpha3,name=validation.awsmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha3-awsmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,versions=v1alpha3,name=default.awsmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None

var _ webhook.Defaulter = &AWSMachinePool{}
var _ webhook.Validator = &AWSMachinePool{}

func (r *AWSMachinePool) validateDefaultCoolDown() []error {
	var allErrs []error
	if int(r.Spec.DefaultCoolDown.Duration.Seconds()) < 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.DefaultCoolDown"), "DefaultCoolDown must be greater than zero"))
	}
	return allErrs
}

// ValidateCreate will do any extra validation when creating a AWSMachinePool
func (r *AWSMachinePool) ValidateCreate() error {
	log.Info("AWSMachinePool validate create", "name", r.Name)

	var allErrs []error

	if errs := r.validateDefaultCoolDown(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	return kerrors.NewAggregate(allErrs)
}

// ValidateUpdate will do any extra validation when updating a AWSMachinePool
func (r *AWSMachinePool) ValidateUpdate(old runtime.Object) error {
	var allErrs []error
	if errs := r.validateDefaultCoolDown(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	return kerrors.NewAggregate(allErrs)
}

// ValidateDelete allows you to add any extra validation when deleting
func (r *AWSMachinePool) ValidateDelete() error {
	return nil
}

// Default will set default values for the AWSMachinePool
func (r *AWSMachinePool) Default() {
	if int(r.Spec.DefaultCoolDown.Duration.Seconds()) == 0 {
		log.Info("DefaultCoolDown is zero, setting 300 seconds as default")
		r.Spec.DefaultCoolDown.Duration = 300 * time.Second
	}
}
