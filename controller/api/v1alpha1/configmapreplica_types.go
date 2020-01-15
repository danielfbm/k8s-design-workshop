package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigMapReplicaSpec defines the desired state of ConfigMapReplica
type ConfigMapReplicaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ConfigMapReplica. Edit ConfigMapReplica_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ConfigMapReplicaStatus defines the observed state of ConfigMapReplica
type ConfigMapReplicaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ConfigMapReplica is the Schema for the configmapreplicas API
type ConfigMapReplica struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigMapReplicaSpec   `json:"spec,omitempty"`
	Status ConfigMapReplicaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConfigMapReplicaList contains a list of ConfigMapReplica
type ConfigMapReplicaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigMapReplica `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigMapReplica{}, &ConfigMapReplicaList{})
}
