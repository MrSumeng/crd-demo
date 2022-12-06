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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OrderSpec defines the desired state of Order
type OrderSpec struct {
	// the information for the Shop
	// +kubebuilder:validation:Required
	Shop *ShopInfo `json:"shop"`

	// Commodities is a list of CommodityInfo
	// +kubebuilder:validation:Required
	Commodities []CommodityInfo `json:"commodity"`

	// TotalPrice is the total price of the Order
	// +kubebuilder:validation:Required
	TotalPrice int64 `json:"totalPrice"`

	// Remark of Order
	// +optional
	Remark string `json:"remark,omitempty"`
}

type ShopInfo struct {
	// Name of the shop
	Name string `json:"name"`
}

type CommodityInfo struct {
	// Name of the commodity
	Name string `json:"name"`

	// Price of the commodity
	Price int64 `json:"price"`

	// Quantity of commodity
	Quantity int64 `json:"quantity"`
}

// OrderStatus defines the observed state of Order
type OrderStatus struct {
	// Conditions a list of conditions an order can have.
	// +optional
	Conditions []OrderCondition `json:"conditions,omitempty"`
	// +optional
	Phase OrderPhase `json:"phase,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

type OrderCondition struct {
	// Type of order condition.
	Type OrderConditionType `json:"type"`
	// Phase of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human-readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type OrderConditionType string

const (
	ConditionShop     OrderConditionType = "Shop"
	ConditionDelivery OrderConditionType = "Delivery"
)

type OrderPhase string

const (
	OrderNotAccepted  OrderPhase = "未接单"
	OrderAccepted     OrderPhase = "已接单"
	OrderInMaking     OrderPhase = "制作中"
	OrderMakeComplete OrderPhase = "制作完成"
	OrderWaiting      OrderPhase = "待配送"
	OrderDelivery     OrderPhase = "配送中"
	OrderFinish       OrderPhase = "订单完成"
)

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.phase",description="The order status phase"
//+kubebuilder:printcolumn:name="MESSAGE",type="string",JSONPath=".status.message",description="The order status message"

// Order is the Schema for the orders API
type Order struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrderSpec   `json:"spec,omitempty"`
	Status OrderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OrderList contains a list of Order
type OrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Order `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Order{}, &OrderList{})
}
