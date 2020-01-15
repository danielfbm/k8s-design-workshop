# ConfigMapReplica 

Sample controller for TDD study

The state of the code repository can be navigated using tags:

- step-1
- step-2 

etc.

## Step 1


Create base code 

- Install [kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)
- Change below `domain`, `license`, `owner`, and `repo` flags


```shell
kubebuilder init --domain example.com --license MIT  --repo github.com/danielfbm/k8s-design-workshop/controller
```


Create a resource `ConfigMapReplica`


```shell
kubebuilder create api --group replica --version v1alpha1 --kind ConfigMapReplica --namespaced=false --resource --controller --example 
```


## Step 2

Change the resource to achieve the business requirements:

Open `api/v1alpha1/configmapreplica_types.go` and update:

```golang
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
```