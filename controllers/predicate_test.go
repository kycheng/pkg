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
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestSecretDataChangedPredicate(t *testing.T) {
	var data = []struct {
		desc string
		old  map[string][]byte
		new  map[string][]byte

		expected bool
	}{
		{
			desc: "old is nil",
			old:  nil,
			new: map[string][]byte{
				"a": []byte("1"),
			},
			expected: true,
		},
		{
			desc: "old is not nil and no changes",
			old: map[string][]byte{
				"a": []byte("1"),
			},
			new: map[string][]byte{
				"a": []byte("1"),
			},
			expected: false,
		},
		{
			desc: "old is not nil and no changes 2",
			old: map[string][]byte{
				"b": []byte("0"),
				"a": []byte("1"),
			},
			new: map[string][]byte{
				"a": []byte("1"),
				"b": []byte("0"),
			},
			expected: false,
		},
		{
			desc: "old is not nil and changes",
			old: map[string][]byte{
				"b": []byte("0"),
				"a": []byte("1"),
			},
			new: map[string][]byte{
				"a": []byte("1"),
				"b": []byte("1"),
			},
			expected: true,
		},
	}

	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {
			g := NewGomegaWithT(t)
			e := event.UpdateEvent{
				ObjectOld: &corev1.Secret{Data: item.old},
				ObjectNew: &corev1.Secret{Data: item.new},
			}
			actual := SecretDataChangedPredicate{}.Update(e)

			g.Expect(actual).Should(BeEquivalentTo(item.expected))
		})

	}
}

func TestAnnotationChangedPredicate(t *testing.T) {
	tests := []struct {
		name           string
		keys           []string
		oldAnnotations map[string]string
		newAnnotations map[string]string
		eventType      string // "create", "update", "delete", "generic"
		expected       bool
	}{
		{
			name:           "create event - no keys specified",
			keys:           nil,
			oldAnnotations: nil,
			newAnnotations: map[string]string{"test": "value"},
			eventType:      "create",
			expected:       true,
		},
		{
			name:           "create event - no keys specified with nil",
			keys:           nil,
			oldAnnotations: nil,
			newAnnotations: nil,
			eventType:      "create",
			expected:       true,
		},
		{
			name:           "create event - specific key changed",
			keys:           []string{"test"},
			oldAnnotations: nil,
			newAnnotations: map[string]string{"test": "value"},
			eventType:      "create",
			expected:       true,
		},
		{
			name:           "create event - irrelevant key changed",
			keys:           []string{"test"},
			oldAnnotations: nil,
			newAnnotations: map[string]string{"other": "value"},
			eventType:      "create",
			expected:       false,
		},
		{
			name:           "update event - keys specified will nil",
			keys:           []string{"test"},
			oldAnnotations: nil,
			newAnnotations: nil,
			eventType:      "update",
			expected:       false,
		},
		{
			name:           "update event - specific key changed",
			keys:           []string{"test"},
			oldAnnotations: map[string]string{"test": "old"},
			newAnnotations: map[string]string{"test": "new"},
			eventType:      "update",
			expected:       true,
		},
		{
			name:           "update event - no change in specified key",
			keys:           []string{"test"},
			oldAnnotations: map[string]string{"test": "same"},
			newAnnotations: map[string]string{"test": "same"},
			eventType:      "update",
			expected:       false,
		},
		{
			name:           "delete event - specific key exists",
			keys:           []string{"test"},
			oldAnnotations: map[string]string{"test": "value"},
			newAnnotations: nil,
			eventType:      "delete",
			expected:       true,
		},
		{
			name:           "generic event - specific key exists",
			keys:           []string{"test"},
			oldAnnotations: map[string]string{"test": "value"},
			newAnnotations: nil,
			eventType:      "generic",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			pred := AnnotationChangedPredicate{
				Keys: tt.keys,
			}

			var result bool
			switch tt.eventType {
			case "create":
				obj := &corev1.Pod{}
				obj.SetAnnotations(tt.newAnnotations)
				result = pred.Create(event.CreateEvent{Object: obj})
			case "update":
				oldObj := &corev1.Pod{}
				newObj := &corev1.Pod{}
				oldObj.SetAnnotations(tt.oldAnnotations)
				newObj.SetAnnotations(tt.newAnnotations)
				result = pred.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})
			case "delete":
				obj := &corev1.Pod{}
				obj.SetAnnotations(tt.oldAnnotations)
				result = pred.Delete(event.DeleteEvent{Object: obj})
			case "generic":
				obj := &corev1.Pod{}
				obj.SetAnnotations(tt.oldAnnotations)
				result = pred.Generic(event.GenericEvent{Object: obj})
			}

			g.Expect(result).Should(Equal(tt.expected))
		})
	}
}
