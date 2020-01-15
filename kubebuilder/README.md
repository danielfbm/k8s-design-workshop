# kubebuilder unit test style

This is the full implementation step-by-step. Kubebuilder utilizes a more straight forward method of unit testing that relies on external dependencies. It follows the same principles of a Integration testing, but running everthing locally. This kind of testing is not common when writting unit tests because it depends heavily on the environment, and because of multiple components involved, it generally has a longer execution time compared to mocks


## Step 1

- Install [kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)
- Change below `domain`, `license`, `owner`, and `repo` flags

```shell
kubebuilder init --domain danielfbm.github.io --license MIT --owner "Daniel Morinigo" --repo github.com/danielfbm/k8s-design-workshop/kubebuilder
```


## Step 2

Create a new resource answering `Y` for all.

```
kubebuilder create api --group ship --version v1beta1 --kind Frigate
```



Open `api/v1beta1/frigate_types.go` and add a `Phase` to `FrigateStatus`:

```golang
// FrigateStatus defines the observed state of Frigate
type FrigateStatus struct {
    // Phase in which the Frigate is currently at
    // this comment will become a CRD definition
	Phase string `json:"phase,omitempty"`
}
```

Update definition for CRD

```shell
make
```


## Step 3

Create test file `controllers/frigate_controller_test.go`:

```golang
package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
	shipv1beta1 "github.com/danielfbm/k8s-design-workshop/controller/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	mgr "sigs.k8s.io/controller-runtime/pkg/manager"
)

/*
In TDD it is generally recommended to not be conserned with implementation
but focus on result, in this controller test case we can define our input (CRD instance)
and focus on the end result.

To simplify the business logic we will just add a Phase "Completed" to the CRD instance
*/
var _ = Describe("Reconcile", func() {

	var (
		// variable used in the test or configuration for tests
		frigate    *shipv1beta1.Frigate
		result     *shipv1beta1.Frigate
		controller *FrigateReconciler
		manager    ctrl.Manager

		opts mgr.Options
		ctx  context.Context
		
		config    *rest.Config
		k8sclient client.Client
		err       error
		stop      chan struct{}
	)

	// Ginkgo framework is based around a few blocks: 
	// Describe, Context, BeforeEach, JustBeforeEach, It, JustAfterEach, AfterEach
	// being that for each Describe/Context every time a function is declared it will be used for each It
	// BeforeEach is generally used for initialization
	// combined with a JustBeforeEach that can be used to run the specific test
	// leaving It to only run specific validations
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

		// Create controller
		controller = &FrigateReconciler{Log: logf.Log}
		err = controller.SetupWithManager(manager)
		Expect(err).ToNot(HaveOccurred(), "building controller")

		// Base data input (can be overwritten, example bellow)
		frigate = &shipv1beta1.Frigate{
			ObjectMeta: metav1.ObjectMeta{Name: "some", Namespace: "default"},
			Spec:       shipv1beta1.FrigateSpec{Foo: "foo"},
		}
	})

	// Here are the steps we take for every test case
	// for this case:
	// 1. create resource (resource data can be overwritten)
	// 2. wait for reconcile loop and keep result in result and err variables
	JustBeforeEach(func() {
		// create resource
		err = k8sclient.Create(ctx, frigate)
		Expect(err).To(BeNil(), "create frigate instance")

		objKey := client.ObjectKey{Namespace: frigate.Namespace, Name: frigate.Name}

		// wait for result
		// for this specific case we can validate the phase but
		// each controller might have a different way to validate
		// when does the reconcile loop finishes
		// For more on Eventually workings: http://onsi.github.io/gomega/
		result = &shipv1beta1.Frigate{}
		Eventually(func() string {
			err = k8sclient.Get(ctx, objKey, result)
			logf.Log.Info("got?", "result", result, "err", err)
			return result.Status.Phase
		}, time.Second).ShouldNot(BeEmpty())
	})

	// Some cleanup tasks between each test case
	AfterEach(func() {
		k8sclient.Delete(ctx, frigate)
		close(stop)
	})

	// This is the specific test case
	// here we will use the default data and variable set in BeforeEach
	// and can validate the result directly
	It("should have a Completed phase", func() {
		Expect(result).ToNot(BeNil(), "should have a result")
		Expect(result.Status.Phase).To(Equal("Completed"))
	})

	// How to reuse all the above code and add a new test case?
	// context can make it happen
	Context("new frigate instance with empty Foo", func() {
		// Adding this method will add a new BeforeEach to be executed
		// right after the one executed on top for this Context
		BeforeEach(func() {
			// lets say for the sake of simplicity that "another" Frigate 
			// should have a "Failure" phase
			frigate = &shipv1beta1.Frigate{
				ObjectMeta: metav1.ObjectMeta{Name: "another", Namespace: "default"},
				Spec:       shipv1beta1.FrigateSpec{Foo: ""},
			}
		})

		It("should have a Failure phase", func() {
			Expect(result).ToNot(BeNil(), "should have a result")
			Expect(result.Status.Phase).To(Equal("Failure"))
		})
	})
})
```

## Step 4

Change `frigate_controller.go` file for the following content:

```golang
package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/apimachinery/pkg/api/errors"

	shipv1beta1 "github.com/danielfbm/k8s-design-workshop/controller/api/v1beta1"
)

// FrigateReconciler reconciles a Frigate object
type FrigateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ship.danielfbm.github.io,resources=frigates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ship.danielfbm.github.io,resources=frigates/status,verbs=get;update;patch

func (r *FrigateReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	ctx := context.Background()
	log := r.Log.WithValues("frigate", req.NamespacedName)
	log.Info("got req", "req",req)

	frigate := &shipv1beta1.Frigate{}
	if err = r.Get(ctx, req.NamespacedName, frigate); err != nil {
		// not found error can be ignore, for all others we return
		// it means the object was delete before the reconcile loop started
		if errors.IsNotFound(err) {
			err = nil
		}
		return
	}


	frigateCopy := frigate.DeepCopy()
	// this logic is simple enough, the point being
	// how to write unit tests (check _test.go file)
	if req.Name == "another" {
		frigateCopy.Status.Phase = "Failure"
	} else {
		frigateCopy.Status.Phase = "Completed"
	}

	err = r.Update(ctx, frigateCopy)
	return
}

func (r *FrigateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.Scheme = mgr.GetScheme()

	return ctrl.NewControllerManagedBy(mgr).
		For(&shipv1beta1.Frigate{}).
		Complete(r)
}
```

run the tests