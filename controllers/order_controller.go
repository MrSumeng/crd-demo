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

package controllers

import (
	"context"
	demov1 "crd-demo/api/v1"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

// OrderReconciler reconciles a Order object
type OrderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const MaxSpeedTime int = 60

//+kubebuilder:rbac:groups=demo.sumeng.com,resources=orders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.sumeng.com,resources=orders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.sumeng.com,resources=orders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Order object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *OrderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the order instance
	order := &demov1.Order{}
	err := r.Get(context.TODO(), req.NamespacedName, order)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	status := order.Status.DeepCopy()
	defer func() {
		err := r.updateScaleStatusInternal(order, *status)
		if err != nil {
			klog.Errorf("update order(%s/%s) status failed: %s",
				order.Namespace, order.Name, err.Error())
			return
		}
	}()

	//Simulate the time spent in each phase
	speedTime := rand.Int() % MaxSpeedTime
	switch status.Phase {
	case "":
		status.Phase = demov1.OrderNotAccepted
		status.Message = "Order not accepted"
	case demov1.OrderNotAccepted:
		status.Phase = demov1.OrderAccepted
		status.Message = "Order accepted"
		cond := NewOrderCondition(demov1.ConditionShop, corev1.ConditionFalse, status.Message, status.Message)
		SetOrderCondition(status, *cond)
	case demov1.OrderAccepted:
		status.Phase = demov1.OrderInMaking
		status.Message = "Order in making"
	case demov1.OrderInMaking:
		status.Phase = demov1.OrderMakeComplete
		status.Message = "Order make complete"
		cond := NewOrderCondition(demov1.ConditionShop, corev1.ConditionTrue, status.Message, status.Message)
		SetOrderCondition(status, *cond)
	case demov1.OrderMakeComplete:
		status.Phase = demov1.OrderWaiting
		status.Message = "Order wait delivery"
		cond := NewOrderCondition(demov1.ConditionDelivery, corev1.ConditionFalse, status.Message, status.Message)
		SetOrderCondition(status, *cond)
	case demov1.OrderWaiting:
		status.Phase = demov1.OrderDelivery
		status.Message = "Order delivery"
	case demov1.OrderDelivery:
		status.Phase = demov1.OrderFinish
		status.Message = "Order finished,customer has signed"
	case demov1.OrderFinish:
		cond := NewOrderCondition(demov1.ConditionDelivery, corev1.ConditionTrue, "Success", status.Message)
		SetOrderCondition(status, *cond)
		return ctrl.Result{}, nil
	}
	return ctrl.Result{RequeueAfter: time.Duration(speedTime) * time.Second}, nil
}

func (r *OrderReconciler) updateScaleStatusInternal(scale *demov1.Order, newStatus demov1.OrderStatus) error {
	if reflect.DeepEqual(scale.Status, newStatus) {
		return nil
	}
	clone := scale.DeepCopy()
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := r.Client.Get(context.TODO(),
			types.NamespacedName{Name: scale.Name, Namespace: scale.Namespace},
			clone); err != nil {
			klog.Errorf("error getting updated scale(%s/%s) from client",
				scale.Namespace, scale.Name)
			return err
		}
		clone.Status = newStatus
		if err := r.Client.Status().Update(context.TODO(), clone); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	oldBy, _ := json.Marshal(scale.Status)
	newBy, _ := json.Marshal(newStatus)
	klog.V(5).Infof("order(%s/%s) status from(%s) -> to(%s)", scale.Namespace, scale.Name, string(oldBy), string(newBy))
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//Order status changes do not trigger the reconcile process
	predicates := builder.WithPredicates(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldObject := e.ObjectOld.(*demov1.Order)
			newObject := e.ObjectNew.(*demov1.Order)
			if oldObject.Generation != newObject.Generation || newObject.DeletionTimestamp != nil {
				klog.V(3).Infof("Observed updated for order: %s/%s", newObject.Namespace, newObject.Name)
				return true
			}
			return false
		},
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.Order{}, predicates).
		Complete(r)
}
