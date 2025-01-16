/*
Copyright 2024 The AlaudaDevops Authors.

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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SecretDataChangedPredicate implements a default update predicate function on secret data change.
type SecretDataChangedPredicate struct {
	predicate.Funcs
}

// Update implements default UpdateEvent filter for validating generation change.
func (SecretDataChangedPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil {
		return false
	}
	oldObj := e.ObjectOld.(*corev1.Secret)

	if e.ObjectNew == nil {
		return false
	}
	newObj := e.ObjectNew.(*corev1.Secret)

	return !reflect.DeepEqual(oldObj.Data, newObj.Data)
}

// AnnotationChangedPredicate implements a predicate that checks for changes in specific annotations.
// It extends the default AnnotationChangedPredicate from controller-runtime and allows filtering
// on specific annotation keys.
type AnnotationChangedPredicate struct {
	// Keys is a list of annotation keys to watch for changes.
	// If empty, all annotation changes will be considered.
	Keys []string
	predicate.AnnotationChangedPredicate
}

// Create implements Predicate interface for creation events.
// It checks if any of the specified annotation keys have changed from nil to a value.
func (p AnnotationChangedPredicate) Create(e event.CreateEvent) bool {

	if len(p.Keys) == 0 {
		return p.AnnotationChangedPredicate.Create(e)
	}

	return valuesChangeInMap(p.Keys, nil, e.Object.GetAnnotations())
}

// Delete implements Predicate interface for deletion events.
// It checks if any of the specified annotation keys have changed from a value to nil.
func (p AnnotationChangedPredicate) Delete(e event.DeleteEvent) bool {

	if len(p.Keys) == 0 {
		return p.AnnotationChangedPredicate.Delete(e)
	}

	return valuesChangeInMap(p.Keys, e.Object.GetAnnotations(), nil)
}

// Generic implements Predicate interface for generic events.
// It checks if any of the specified annotation keys have changed.
func (p AnnotationChangedPredicate) Generic(e event.GenericEvent) bool {

	if len(p.Keys) == 0 {
		return p.AnnotationChangedPredicate.Generic(e)
	}

	return valuesChangeInMap(p.Keys, e.Object.GetAnnotations(), nil)
}

// Update implements Predicate interface for update events.
// It checks if any of the specified annotation keys have different values between old and new objects.
func (p AnnotationChangedPredicate) Update(e event.UpdateEvent) bool {

	if len(p.Keys) == 0 {
		return p.AnnotationChangedPredicate.Update(e)
	}

	if e.ObjectOld == nil {
		return false
	}
	if e.ObjectNew == nil {
		return false
	}

	return valuesChangeInMap(p.Keys, e.ObjectOld.GetAnnotations(), e.ObjectNew.GetAnnotations())
}

// valuesChangeInMap checks if any of the specified keys have different values in two maps.
// Returns true if there's a difference in values for any of the specified keys.
func valuesChangeInMap(keys []string, old, new map[string]string) bool {
	var getValue = func(key string, kv map[string]string) string {
		if len(kv) == 0 {
			return ""
		}
		return kv[key]
	}

	for _, key := range keys {
		if getValue(key, old) != getValue(key, new) {
			return true
		}
	}

	return false
}
