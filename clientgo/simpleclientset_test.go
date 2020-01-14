package main

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

/*
NewSimpleClientset constructor can help with most of all the
standard positive test cases and other simple test cases:

1. empty set, create a resource, get the same resource
2. check an not existing resource, create and update a resource, etc.

But cannot satisfy negative test cases when reproducing specific failures like:
1. internal server errors
2. specific validation errors
3. specific actions returning error etc.
*/
func TestSimpleClientset(t *testing.T) {
	// Creating fake client set
	// if you need to populate some necessary data from any api group
	// it can be given to the cosntructor
	fakeclient := fake.NewSimpleClientset()
	// make sure it is kubernetes.Interface
	var _ kubernetes.Interface = fakeclient

	// GET Getting from an empty set will return an NotFound error
	configmap, err := fakeclient.CoreV1().ConfigMaps("default").Get("test", v1.GetOptions{})
	if err == nil || !errors.IsNotFound(err) {
		t.Errorf("this should return a 404 error: %v, %v", configmap, err)
	}

	// CREATE Using its API to create objects also works
	configmap = &corev1.ConfigMap{ObjectMeta: v1.ObjectMeta{Name: "test", Namespace: "default"}}
	configmap, err = fakeclient.CoreV1().ConfigMaps("default").Create(configmap)
	if err != nil {
		t.Errorf("should not return any error: %v, %v", configmap, err)
	}
	// Getting now should not return any error
	configmap, err = fakeclient.CoreV1().ConfigMaps("default").Get("test", v1.GetOptions{})
	if err != nil {
		t.Errorf("should not return any error: %v, %v", configmap, err)
	}

	// UPDATE Updates also works fine
	configmap = configmap.DeepCopy()
	configmap.Labels = map[string]string{"key":"value"}
	configmap, err = fakeclient.CoreV1().ConfigMaps("default").Update(configmap)
	if err != nil {
		t.Errorf("should not return any error: %v, %v", configmap, err)
	}
	// getting will have latest changes
	cm, err := fakeclient.CoreV1().ConfigMaps("default").Get("test", v1.GetOptions{})
	if err != nil || cm.Labels == nil || cm.Labels["key"] != "value" {
		t.Errorf("did not get updated version of the object")
	}

}
