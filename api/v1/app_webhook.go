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

package v1

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var applog = logf.Log.WithName("app-resource")

func (r *App) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-app-bmutziu-me-v1-app,mutating=true,failurePolicy=fail,sideEffects=None,groups=app.bmutziu.me,resources=apps,verbs=create;update,versions=v1,name=mapp.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &App{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *App) Default() {
	applog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	var cns []core.Container
	cns = r.Spec.Deploy.Template.Spec.Containers

	container := core.Container{
		Name:  "sidecar-nginx",
		Image: "nginx:1.12.2",
	}

	cns = append(cns, container)
	r.Spec.Deploy.Template.Spec.Containers = cns
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-app-bmutziu-me-v1-app,mutating=false,failurePolicy=fail,sideEffects=None,groups=app.bmutziu.me,resources=apps,verbs=create;update,versions=v1,name=vapp.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &App{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *App) ValidateCreate() error {
	applog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *App) ValidateUpdate(old runtime.Object) error {
	applog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	a := old.(*App)
	if r.Spec.Deploy.Template.Spec.Containers[0].Image != a.Spec.Deploy.Template.Spec.Containers[0].Image {
		return field.Invalid(field.NewPath("spec").Child("deploy").Child("template").Child("spec").Child("containers").Index(0).Child("image"), r.Spec.Deploy.Template.Spec.Containers[0].Image, "For the 1st container image is immutable")
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *App) ValidateDelete() error {
	applog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
