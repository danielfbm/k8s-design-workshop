package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
	log.Info("got req", "req", req)

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
