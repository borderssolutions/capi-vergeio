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

package controllers

import (
	"context"

	infrav1 "github.com/mashalabbas/cluster-api-provider-verge/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	//import the pkg verge apis
	vergeapis "github.com/mashalabbas/cluster-api-provider-verge/pkg/verge"

	//need to add this to the repo: TODO
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// VergeMachineReconciler reconciles a VergeMachine object
type VergeMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vergemachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vergemachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vergemachines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VergeMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// // - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
// func (r *VergeMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
// 	_ = log.FromContext(ctx)

// 	// TODO(user): your logic here

// 	return ctrl.Result{}, nil
// }

// // SetupWithManager sets up the controller with the Manager.
// func (r *VergeMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
// 	return ctrl.NewControllerManagedBy(mgr).
// 		For(&infrastructurev1alpha1.VergeMachine{}).
// 		Complete(r)
// }

func (r *VergeMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// logger := log.FromContext(ctx)

	vergeMachine := &infrav1.VergeMachine{}
	if err := r.Get(ctx, req.NamespacedName, vergeMachine); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// If being deleted
	if !vergeMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, vergeMachine)
	}

	return r.reconcileNormal(ctx, vergeMachine)
}

func (r *VergeMachineReconciler) reconcileNormal(ctx context.Context, vergeMachine *infrav1.VergeMachine) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if vergeMachine.Spec.ProviderID != nil && *vergeMachine.Spec.ProviderID != "" {
		logger.Info("VergeMachine already has a ProviderID", "ProviderID", *vergeMachine.Spec.ProviderID)
		return ctrl.Result{}, nil
	}

	machine, err := vergeapis.GetOwnerMachine(ctx, r.Client, vergeMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}

	if !vergeMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	// Check if the machine is part of the control plane
	if clusterv1.IsControlPlaneMachine(machine) {
		logger.Info("Provisioning control plane node")
		// Add control-plane specific logic here
	} else {
		logger.Info("Provisioning worker node")
		// Add worker-specific logic here
	}

	providerID, err := vergeapis.CreateMachineInVergeInfra(machine, *machine.Spec.Bootstrap.Data)
	if err != nil {
		return ctrl.Result{}, err
	}

	vergeMachine.Spec.ProviderID = &providerID
	vergeMachine.Status.Ready = true

	if err := r.Status().Update(ctx, vergeMachine); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VergeMachineReconciler) reconcileDelete(ctx context.Context, vergeMachine *infrav1.VergeMachine) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	machine, err := vergeapis.GetOwnerMachine(ctx, r.Client, vergeMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("Deleting infrastructure machine for VergeMachine", "Machine", machine.Name)
	err = vergeapis.DeleteMachineInVergeInfra(machine)
	if err != nil {
		return ctrl.Result{}, err
	}

	// controllerutil.RemoveFinalizer(vergeMachine, infrav1.MachineFinalizer)
	if err := r.Update(ctx, vergeMachine); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VergeMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.VergeMachine{}).
		Owns(&clusterv1.Machine{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
