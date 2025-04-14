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
	corev1 "github.com/rajeshdeshpande02/platform-engineering-labs/k8s-operators/namespace-ripper/api/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// NamespaceTTLReconciler reconciles a NamespaceTTL object
type NamespaceTTLReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.pelabs.com,resources=namespacettls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.pelabs.com,resources=namespacettls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.pelabs.com,resources=namespacettls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NamespaceTTL object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *NamespaceTTLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)

	// Fetch the NamespaceTTL instance
	var namespaceTTL corev1.NamespaceTTL
	if err := r.Get(ctx, req.NamespacedName, &namespaceTTL); err != nil {
		log.Error(err, "unable to fetch NamespaceTTL")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Access the spec field (for example TTL or AdditionalField)
	ttlDuration, err := time.ParseDuration(namespaceTTL.Spec.TTL)
	if err != nil {
		log.Error(err, "Invalid TTL format")
		return ctrl.Result{}, err
	}

	namespaceList := &v1.NamespaceList{}
	err = r.Client.List(ctx, namespaceList)
	if err != nil {
		log.Error(err, "unable to list namespaces")
		return ctrl.Result{}, err
	}
	protectedNamespace := []string{"kube-system", "kube-public", "kube-node-lease", "namespace-ripper-system", "default"}
	for _, ns := range namespaceList.Items {

		currentTime := time.Now()
		creationTimeSeconds := ns.CreationTimestamp.Time
		elapsedTime := currentTime.Sub(creationTimeSeconds)
		if elapsedTime > ttlDuration {
			if notInList(ns.Name, protectedNamespace) {
				log.Info("Deleting Namespace", "Name", ns.Name)
				err := r.Client.Delete(ctx, &ns)
				if err != nil {
					log.Error(err, "unable to delete namespace")
					return ctrl.Result{}, err
				}
			}
		}
	}

	// Requeue the reconciliation after 15 seconds. This is imp bcz we dont have any trigger for reconcilation
	return ctrl.Result{
		RequeueAfter: 15 * time.Second,
	}, nil
}

func notInList(val string, list []string) bool {
	// Convert list to a map
	set := make(map[string]struct{}, len(list))
	for _, item := range list {
		set[item] = struct{}{}
	}

	// Check if the value is in the set
	_, found := set[val]
	return !found
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceTTLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.NamespaceTTL{}).
		Complete(r)
}
