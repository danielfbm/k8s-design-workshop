package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapReplicaSpec defines the desired state of ConfigMapReplica
type ConfigMapReplicaSpec struct {
	// Template defines the data that should be replicated
	Template ConfigMapTemplate `json:"template"`

	// Selector as namespace selector rule to replicate configmaps to
	Selector map[string]string `json:"selector"`
}

// ConfigMapTemplate template data for all replicated ConfigMaps
type ConfigMapTemplate struct {
	// Labels to be given to replicated ConfigMap
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Data to be replicated
	Data map[string]string `json:"data,omitempty"`
}

// ConfigMapReplicaStatus defines the observed state of ConfigMapReplica
type ConfigMapReplicaStatus struct {
	// Status for each configmap
	// +optional
	ConfigMapStatuses []ConfigMapReplicaCopy `json:"configMapStatuses,omitempty"`
}

// ConfigMapReplicaCopy a condition for one Copy
type ConfigMapReplicaCopy struct {
	// Name for resource
	Name string `json:"name"`
	// Namespace of resource
	Namespace string `json:"namespace"`
	// Last time we probed the condition
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transitioned
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Ready returns true when a configmap is ready
	Ready bool `json:"ready"`
	// Reason for not being ready. CamelCase
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message detail for Reason
	// +optional
	Message string `json:"message,omitempty"`
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
