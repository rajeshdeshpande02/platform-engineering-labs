/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
	//"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	//pelabscorev1 "github.com/rajeshdeshpande02/platform-engineering-labs/k8s-operators/mini-reloade/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// AutoReloadConfigReconciler reconciles a AutoReloadConfig object
type AutoReloadConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.pelabs.com,resources=autoreloadconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.pelabs.com,resources=autoreloadconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.pelabs.com,resources=autoreloadconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutoReloadConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *AutoReloadConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling AutoReloadConfig", "name", req.Name, "namespace", req.Namespace)
	cm := &corev1.ConfigMap{}
	if err := r.Get(ctx, req.NamespacedName, cm); err != nil {
		log.Error(err, "unable to fetch ConfigMap")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Change detected in ConfigMap:", "name", cm.Name)
	err := r.restartDeployments(ctx, cm)

	if err != nil {
		log.Error(err, "Failed to restart deployments")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AutoReloadConfigReconciler) restartDeployments(ctx context.Context, cm *corev1.ConfigMap) error {

	var deployments appsv1.DeploymentList
	err := r.Client.List(ctx, &deployments, &client.ListOptions{
		Namespace: cm.Namespace,
	})
	if err != nil {
		return err
	}

	for _, deployemt := range deployments.Items {
		if isCMUsedByDeployment(cm, &deployemt) {
			deployemt.Spec.Template.Labels["restartedAt"] = fmt.Sprintf("%v", time.Now().Unix())
			err := r.Client.Update(ctx, &deployemt)
			if err != nil {
				return err
			}
			//	log.Info("Restarted deployment", "name", deployemt.Name)
		}
	}
	return nil
}

func isCMUsedByDeployment(cm *corev1.ConfigMap, dep *appsv1.Deployment) bool {

	for _, vol := range dep.Spec.Template.Spec.Volumes {
		if vol.ConfigMap != nil && vol.ConfigMap.Name == cm.Name {
			return true
		}
	}

	for _, container := range dep.Spec.Template.Spec.Containers {
		for _, env := range container.Env {
			if env.ValueFrom.ConfigMapKeyRef != nil && env.ValueFrom.ConfigMapKeyRef.Name == cm.Name {
				return true
			}
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoReloadConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(obj client.Object) bool {
			return obj.GetLabels()["mini-reloader.pelabs/enabled"] == "true"
		})).
		Complete(r)
}
