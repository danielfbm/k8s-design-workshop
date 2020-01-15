package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	replicav1alpha1 "github.com/danielfbm/k8s-design-workshop/controller/api/v1alpha1"
)

// ConfigMapReplicaReconciler reconciles a ConfigMapReplica object
type ConfigMapReplicaReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=replica.example.com,resources=configmapreplicas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=replica.example.com,resources=configmapreplicas/status,verbs=get;update;patch

func (r *ConfigMapReplicaReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("configmapreplica", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ConfigMapReplicaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&replicav1alpha1.ConfigMapReplica{}).
		Complete(r)
}
