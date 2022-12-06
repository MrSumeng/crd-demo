/*
Copyright 2022 The KubePort Authors.
*/

package controllers

import (
	demov1 "crd-demo/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewOrderCondition creates a new order condition.
func NewOrderCondition(condType demov1.OrderConditionType, status corev1.ConditionStatus, reason, message string) *demov1.OrderCondition {
	return &demov1.OrderCondition{
		Type:               condType,
		Status:             status,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func GetOrderCondition(status demov1.OrderStatus, condType demov1.OrderConditionType) *demov1.OrderCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func SetOrderCondition(status *demov1.OrderStatus, condition demov1.OrderCondition) bool {
	currentCond := GetOrderCondition(*status, condition.Type)
	if currentCond != nil && currentCond.Status == condition.Status && currentCond.Reason == condition.Reason &&
		currentCond.Message == condition.Message {
		return false
	}
	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := filterOrderCondition(status.Conditions, condition.Type)
	status.Conditions = append(newConditions, condition)
	return true
}

// filterOrderCondition returns a new slice of Order conditions without conditions with the provided type.
func filterOrderCondition(conditions []demov1.OrderCondition, condType demov1.OrderConditionType) []demov1.OrderCondition {
	var newConditions []demov1.OrderCondition
	for _, c := range conditions {
		if c.Type == condType {
			continue
		}
		newConditions = append(newConditions, c)
	}
	return newConditions
}
