/*
Copyright 2023 The AlaudaDevops Authors.

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

package v1alpha1

// Common Annotations
const (
	// DisplayNameAnnotationKey display name for objects
	DisplayNameAnnotationKey = "cpaas.io/displayName"
	// CreatedTimeAnnotationKey creation time for objects
	CreatedTimeAnnotationKey = "cpaas.io/creationTime"
	// UpdatedTimeAnnotationKey update time for objects
	UpdatedTimeAnnotationKey = "cpaas.io/updateTime"
	// DeletedTimeAnnotationKey deletion time for objects
	DeletedTimeAnnotationKey = "cpaas.io/deletionTime"

	// NamespaceAnnotationKey namespace of objects
	NamespaceAnnotationKey = "cpaas.io/namespace"

	// CreatedByAnnotationKey annotation key to store resource creation username
	CreatedByAnnotationKey = "cpaas.io/createdBy"
	// UpdatedByAnnotationKey annotation key to store resource update username
	UpdatedByAnnotationKey = "cpaas.io/updatedBy"
	// DeletedByAnnotationKey annotation key to store resource update username
	DeletedByAnnotationKey = "cpaas.io/deletedBy"

	// UIDescriptorsAnnotationKey annotation for storing ui descriptors in resources
	UIDescriptorsAnnotationKey = "ui.cpaas.io/descriptors"
)
