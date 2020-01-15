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

On terminal run

```shell
make
make install
```

the last command will fail, but will generate the necessary crd files

## Step 3

Now starting the TDD cycle we need to create the test cases. Create a `controlers/configmapreplica_controller_test.go` with the basic structure:


```golang
package controllers

import (
	"context"
	replicav1alpha1 "github.com/danielfbm/k8s-design-workshop/controller/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// logf "sigs.k8s.io/controller-runtime/pkg/log"
	mgr "sigs.k8s.io/controller-runtime/pkg/manager"
	// "time"
)

var _ = Describe("ConfigMapReplica.Reconcile", func() {

	var (
		// variable used in the test or configuration for tests
		input     *replicav1alpha1.ConfigMapReplica
		manager   ctrl.Manager
		opts      mgr.Options
		ctx       context.Context
		config    *rest.Config
		k8sclient client.Client
		err       error
		stop      chan struct{}
	)

	// Basic initialization
	BeforeEach(func() {
		// Basic initialization
		// cfg  and k8sClient variables declared on suite_test.go
		config = cfg
		k8sclient = k8sClient
		stop = make(chan struct{})
		ctx = context.TODO()

		// Create and start manager
		manager, err = ctrl.NewManager(config, opts)
		Expect(err).ToNot(HaveOccurred(), "building manager")
		go func() {
			Expect(manager.Start(stop)).ToNot(HaveOccurred(), "starting manager")
		}()
	})

	// TODO: add specific api calls
	JustBeforeEach(func() {
		Expect(k8sclient).ToNot(BeNil())
		Expect(ctx).ToNot(BeNil())
	})

	// TODO: add cleanup code
	AfterEach(func() {
		close(stop)
	})

	// not a test case, just to make sure it compiles
	It("TODO: implement real test case", func() {
		Expect(input).To(BeNil())
	})
})
```

If run the test cases it should succeed. To follow TDD we should implement a test and make sure it fails.

For this controller we need to have `Namespace`s and `ConfigMapReplica`s. Lets change the basic structure to support our use case:

```golang
package controllers

import (
	"context"
	replicav1alpha1 "github.com/danielfbm/k8s-design-workshop/controller/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	mgr "sigs.k8s.io/controller-runtime/pkg/manager"
	"time"
)

var _ = Describe("ConfigMapReplica.Reconcile", func() {

	var (
		// variable used in the test or configuration for tests
		input, result *replicav1alpha1.ConfigMapReplica
		// namespaces to create
		namespaces []*corev1.Namespace
		// number of configmaps to be expected
		expectedConfigmapNumber int
		manager                 ctrl.Manager
		controller              *ConfigMapReplicaReconciler
		opts                    mgr.Options
		ctx                     context.Context
		config                  *rest.Config
		k8sclient               client.Client
		err                     error
		stop                    chan struct{}
	)

	// Basic initialization
	BeforeEach(func() {
		// Basic initialization
		// cfg  and k8sClient variables declared on suite_test.go
		config = cfg
		k8sclient = k8sClient
		stop = make(chan struct{})
		ctx = context.TODO()
		namespaces = []*corev1.Namespace{}

		// Create and start manager
		manager, err = ctrl.NewManager(config, opts)
		Expect(err).ToNot(HaveOccurred(), "building manager")
		go func() {
			Expect(manager.Start(stop)).ToNot(HaveOccurred(), "starting manager")
		}()

		// Create and start controller
		controller = &ConfigMapReplicaReconciler{Log: logf.Log}
		Expect(controller.SetupWithManager(manager)).To(Succeed(), "starting controller")

		// this input data is invalid on purpose, it should be added using a specific
		// context and valid test case
		input = &replicav1alpha1.ConfigMapReplica{ObjectMeta: metav1.ObjectMeta{Name: "a"}}
	})

	JustBeforeEach(func() {
		// initialize namespaces
		// if necessary add all needed namespaces
		for _, ns := range namespaces {
			Expect(k8sclient.Create(ctx, ns)).To(Succeed(), "should create ns %s", ns.Name)
		}

		// initialize input
		Expect(k8sclient.Create(ctx, input)).To(Succeed(), "should create a configmapreplica %s", input)

		// wait for reconcile loop to finish
		// in this case we will check the status of ConfigMapReplica
		// but it can be any other way
		result = &replicav1alpha1.ConfigMapReplica{}
		objKey := client.ObjectKey{Name: input.Name}
		Eventually(func() int {
			err = k8sclient.Get(ctx, objKey, result)
			if err != nil {
				return -1
			}
			return len(result.Status.ConfigMapStatuses)
		}, 
		// This is the timeout time for this Eventually process
		// for more information check  http://onsi.github.io/gomega/
		time.Second,
		).Should(Equal(expectedConfigmapNumber), "should have %d configmaps", expectedConfigmapNumber)
	})


	// Basic cleanup code
	AfterEach(func() {
		k8sclient.Delete(ctx, input)
		k8sclient.DeleteAllOf(ctx, &corev1.ConfigMap{})
		k8sclient.DeleteAllOf(ctx, &corev1.Namespace{})
		close(stop)
	})

	// not a test case, just to make sure it compiles
	It("TODO: implement real test case", func() {
		Expect(input).To(BeNil())
	})
})

```

**Now the test cases will start to fail**. In the next step we will implement a test case and make sure the validation works

## Step 4

Add the first test case to make sure we can correctly validate our business logic

Delete:

```golang
    // not a test case, just to make sure it compiles
	It("TODO: implement real test case", func() {
		Expect(input).To(BeNil())
	})
```

Add the following test case:

```golang

	Context("one namespace with matching label", func() {
		BeforeEach(func() {
			// add this namespace to make sure it will be generated
			namespaces = append(namespaces, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "sample",
					Labels: map[string]string{"key": "value"},
				},
			})

			input = &replicav1alpha1.ConfigMapReplica{
				ObjectMeta: metav1.ObjectMeta{
					Name: "replica",
				},
				Spec: replicav1alpha1.ConfigMapReplicaSpec{
					Template: replicav1alpha1.ConfigMapTemplate{
						Labels: map[string]string{},
						Data: map[string]string{"data.yaml": "some value for configmap"},
					},
					Selector: map[string]string{"key": "value"},
				},
			}
			expectedConfigmapNumber = 1
		})

		It("should have one configmap", func() {
			list := &corev1.ConfigMapList{}
			Expect(k8sclient.List(ctx, list)).To(Succeed(), "listing configmaps")

			Expect(list).ToNot(BeNil(), "should have a configmap list")
			Expect(list.Items).To(HaveLen(1), "should have 1 configmap")
			Expect(result).ToNot(BeNil(), "crd should exist")
			Expect(result.Status.ConfigMapStatuses).To(HaveLen(1), "should have 1 configmapStatus")
		})
	})
```


Running should give the following error

```shell
    Expected
        <int>: 0
    to equal
        <int>: 1
```

Now we are ready to implement our first reconcile case


## Step 5

We can implement the reconciler logic to satisfy this test case

Inside the `controllers/configmapreplica_controller.go` file replace content with

```golang
package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/api/errors"

	replicav1alpha1 "github.com/danielfbm/k8s-design-workshop/controller/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ConfigMapReplicaReconciler reconciles a ConfigMapReplica object
type ConfigMapReplicaReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=replica.example.com,resources=configmapreplicas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=replica.example.com,resources=configmapreplicas/status,verbs=get;update;patch

func (r *ConfigMapReplicaReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	ctx := context.Background()
	log := r.Log.WithValues("configmapreplica", req.NamespacedName)

	configMapReplica := &replicav1alpha1.ConfigMapReplica{}
	if err = r.Get(ctx, req.NamespacedName, configMapReplica); err != nil {
		// not found error can be ignore, for all others we return
		// it means the object was delete before the reconcile loop started
		if errors.IsNotFound(err) {
			err = nil
		}
		return
	}

	// build selector from labels in spec
	selector := labels.SelectorFromSet(configMapReplica.Spec.Selector)
	namespaceList := &corev1.NamespaceList{}
	if err = r.List(ctx, namespaceList, &client.ListOptions{LabelSelector: selector}); err != nil {
		// log.Error("error listing namespace", "err", err)
		log.Error(err, "selector", selector)
		return
	}

	// base data for syncing
	baseConfigmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configMapReplica.Name,
			Labels: configMapReplica.Spec.Template.Labels,
		},
		Data: configMapReplica.Spec.Template.Data,
	}
	// making it editable
	configMapReplica = configMapReplica.DeepCopy()
	if configMapReplica.Status.ConfigMapStatuses == nil {
		configMapReplica.Status.ConfigMapStatuses = []replicav1alpha1.ConfigMapReplicaCopy{}
	}
	if err = controllerutil.SetControllerReference(configMapReplica, baseConfigmap, r.Scheme); err != nil {
		log.Error(err, "base", baseConfigmap, "owner", configMapReplica)
		return
	}

	for _, ns := range namespaceList.Items {
		clone := baseConfigmap.DeepCopy()
		clone.Namespace = ns.Name

		current := &corev1.ConfigMap{}
		key := types.NamespacedName{Namespace: ns.Name, Name: clone.Name}
		err = r.Get(ctx, key, current)
		switch {
			// no item, we can create
		case errors.IsNotFound(err):
			log.Info("will create configmap", "configmap", clone.ObjectMeta)
			err = r.Create(ctx, clone)
			configMapReplica.Status.ConfigMapStatuses = append(configMapReplica.Status.ConfigMapStatuses, replicav1alpha1.ConfigMapReplicaCopy{
				Name: clone.Name,
				Namespace: clone.Namespace,
				Ready: err == nil,
				LastTransitionTime: metav1.Now(),
				LastProbeTime: metav1.Now(),
			})

			// item exist. Should we update?
		case  err == nil:
			// TODO: add update

		}
	}

	err = r.Update(ctx, configMapReplica)
	log.Info("update?", "err", err)
	return
}

func (r *ConfigMapReplicaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.Scheme = mgr.GetScheme()
	return ctrl.NewControllerManagedBy(mgr).
		For(&replicav1alpha1.ConfigMapReplica{}).
		Complete(r)
}
```


Running this package tests should succeed. 

This completes one cycle of the TDD loop. 

## Next Steps:

As you may noticed there are several use cases that were not supported. It is up to you to add the last final touches to this controller and its test cases.