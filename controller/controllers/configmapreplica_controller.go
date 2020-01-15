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
