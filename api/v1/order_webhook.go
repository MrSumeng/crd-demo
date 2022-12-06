/*
Copyright 2022 The CRD-Demo Authors.

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

package v1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var orderlog = logf.Log.WithName("order-resource")

func (in *Order) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-demo-sumeng-com-v1-order,mutating=true,failurePolicy=fail,sideEffects=None,groups=demo.sumeng.com,resources=orders,verbs=create;update,versions=v1,name=morder.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Order{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Order) Default() {
	// Set the default value.
	// However, we have noting to do in this crd resources.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-demo-sumeng-com-v1-order,mutating=false,failurePolicy=fail,sideEffects=None,groups=demo.sumeng.com,resources=orders,verbs=create;update,versions=v1,name=vorder.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Order{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Order) ValidateCreate() error {
	orderlog.Info("validate create", "name", in.Name)
	return in.Spec.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Order) ValidateUpdate(old runtime.Object) error {
	orderlog.Info("validate update", "name", in.Name)
	return in.Spec.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Order) ValidateDelete() error {
	return nil
}

func (in *OrderSpec) validate() error {
	if in.TotalPrice <= 0 {
		return fmt.Errorf("total price must be greater than 0")
	}

	var totalPrice int64 = 0
	for i := range in.Commodities {
		err := in.Commodities[i].validate()
		if err != nil {
			return err
		}
		totalPrice += in.Commodities[i].Price * in.Commodities[i].Quantity
	}

	if totalPrice != in.TotalPrice {
		return fmt.Errorf("the total price of the item is incorrect")
	}
	return nil
}

func (in *CommodityInfo) validate() error {
	if in.Quantity <= 0 {
		return fmt.Errorf("commodity %s quantity must be greater than 0", in.Name)
	}
	if in.Price <= 0 {
		return fmt.Errorf("commodity %s price must be greater than 0", in.Name)
	}
	return nil
}
