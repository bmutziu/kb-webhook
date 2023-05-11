/*
Copyright 2023 bmutziu.

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
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/bmutziu/kb-webhook/api/v1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.bmutziu.me,resources=apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.bmutziu.me,resources=apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.bmutziu.me,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	instance := &appv1.App{}

	err := r.Client.Get(context.Background(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			log.Info(fmt.Sprintf(`Instance for app "%s" does not exist, should delete associated resources`, req.Name))
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, fmt.Sprintf(`Failed to retrieve instance "%s"`, req.Name))
		return ctrl.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		log.Info("Get deleted App, clean up subResources")
		return ctrl.Result{}, nil
	}

	log.Info(`Hello from your new app reconciler!`)

	labels := make(map[string]string)
	labels["app"] = instance.Name

	deploySpec := appsv1.DeploymentSpec{}
	deploySpec = instance.Spec.Deploy
	deploySpec.Selector = &metav1.LabelSelector{MatchLabels: labels}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-deploy",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: deploySpec,
	}

	scheme := runtime.Scheme{}
	if err := controllerutil.SetControllerReference(instance, deploy, &scheme); err != nil {
		log.Error(err, "Set DeployVersion CtlRef Error")
		return ctrl.Result{}, err
	}

	found := &appsv1.Deployment{}
	err = r.Client.Get(context.Background(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Old Deployment NotFound and Creating New One", "namespace", deploy.Namespace, "name", deploy.Name)
		if err = r.Client.Create(ctx, deploy); err != nil {
			log.Error(err, fmt.Sprintf(`Failed to create deployment for instance "%s"`, instance.Name))
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Get Deployment Info Error", "namespace", deploy.Namespace, "name", deploy.Name)
		return ctrl.Result{}, err
	} else if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		// Update the found object and write the result back if there are any changes
		found.Spec = deploy.Spec
		log.Info("Old Deployment Changed and Updating Deployment to Reconcile", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Client.Update(ctx, found)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.App{}).
		Complete(r)
}
