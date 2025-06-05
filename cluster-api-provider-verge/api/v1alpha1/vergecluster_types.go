/*
Copyright 2025.

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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
const (
	// ClusterFinalizer allows cleaning up resources associated with
	// VergeCluster before removing it from the apiserver.
	ClusterFinalizer = "vergecluster.infrastructure.cluster.x-k8s.io"
)

// VergeClusterSpec defines the desired state of VergeCluster
type VergeClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

// VergeClusterStatus defines the observed state of VergeCluster
type VergeClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready indicates that the cluster is ready.
	// +optional
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VergeCluster is the Schema for the vergeclusters API
type VergeCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VergeClusterSpec   `json:"spec,omitempty"`
	Status VergeClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VergeClusterList contains a list of VergeCluster
type VergeClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VergeCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VergeCluster{}, &VergeClusterList{})
}
