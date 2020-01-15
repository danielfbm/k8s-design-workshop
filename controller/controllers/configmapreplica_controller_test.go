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
})
