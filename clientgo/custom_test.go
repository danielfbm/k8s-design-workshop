package main

/* Warning: First check simpleclientset_test.go */

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

/*
When facing cases where the NewSimpleClientset cannot help
is necessary to implement a few things to just get it right.

Below is the implementation of a very simple case:
1. Gets a configmap: it returns
2. Updates a configmap: an error is returned.

There are more complex use cases that can involve more complex operations like `watch`  please
refer to the source code

In any case both can be used together on complementing cases
*/
func TestClientset(t *testing.T) {
	// basic initialization, extrated from NewSimpleClientset code
	cs := &fake.Clientset{}

	// Defining behaviour for the Clientset using Reactors.
	// 1. Should get a object
	// for this case there are two methods:
	// 1.1 use a object tracker to activate specific behaviour
	scheme := runtime.NewScheme()
	corev1.SchemeBuilder.AddToScheme(scheme)
	codecs := serializer.NewCodecFactory(scheme)
	tracker := k8stesting.NewObjectTracker(scheme, codecs.UniversalDecoder())
	cs.AddReactor("get", "configmaps", k8stesting.ObjectReaction(tracker))

	expected := &corev1.ConfigMap{ObjectMeta: v1.ObjectMeta{Name: "test", Namespace: "default"}}
	tracker.Add(expected)
	cm, err := cs.CoreV1().ConfigMaps("default").Get("test", v1.GetOptions{})
	if err != nil {
		t.Errorf("should be able to get configmap and not return error, %v %v", cm, err)
	}

	// 1.2 implement the Reaction
	// using prepend to put it before the tracker implementation of the method above
	cs.PrependReactor("get", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		// this implementation is just an example on the complexity level
		// that it can support. for most use cases the behaviour should be simple
		// to not need much logic

		ns := action.GetNamespace()
		gvr := action.GetResource()
		switch act := action.(type) {
		case k8stesting.GetActionImpl:
			// return our object
			if ns == expected.Namespace {
				ret = expected
			} else {
				// name only available on specific type
				err = errors.NewNotFound(gvr.GroupResource(), act.GetName())
			}
		default:
			// returning some kind of error just in case
			err = errors.NewMethodNotSupported(gvr.GroupResource(), fmt.Sprintf("%s", act))
		}
		// set this to halt and return result from this Reactor
		handled = true
		return
	})
	// Make our updates always fails
	cs.PrependReactor("update", "configmaps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		err = errors.NewInternalError(fmt.Errorf("crash?"))
		// this flag will let us give a passthrough if not set
		handled = true
		return
	})

	// check that it works as expected
	_, err = cs.CoreV1().ConfigMaps("default").Get("test", v1.GetOptions{})
	t.Logf("%v", err)
	if err != nil {
		t.Errorf("should not fail")
	}

	_, err = cs.CoreV1().ConfigMaps("other").Get("test", v1.GetOptions{})
	t.Logf("%v", err)
	if err == nil {
		t.Errorf("should fail")
	}

	_, err = cs.CoreV1().ConfigMaps("default").Update(expected)
	t.Logf("%v", err)
	if err == nil || !errors.IsInternalError(err) {
		t.Errorf("is not the error defined: %v", err)
	}
}
