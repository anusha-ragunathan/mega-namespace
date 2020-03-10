/*


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
	"fmt"
	"os/exec"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	megav1 "github.com/anusha-ragunathan/mega-namespace/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceTemplateReconciler reconciles a NamespaceTemplate object
type NamespaceTemplateReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	PostCreateHookMap map[string]bool
}

// creates additional resources. All additional resources listed in nstObj.Spec.AdditionalResources
// are created in the ns namespace
func (r *NamespaceTemplateReconciler) createadditionalresources(ctx context.Context, nsName string, nstObj megav1.NamespaceTemplate) error {

	// provision pods
	var p v1.Pod
	key := types.NamespacedName{Namespace: nsName, Name: nstObj.Spec.AddResources.Pod.Name}
	if err := r.Get(ctx, key, &p); err != nil {
		fmt.Printf("unable to get pod info due to %v\n", err)
		// assume that the error is "pod doesnt exist". in theory, err can be due to other issues as well
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nsName,
				Name:      nstObj.Spec.AddResources.Pod.Name,
			},
			Spec: nstObj.Spec.AddResources.Pod.Spec,
		}
		if err := r.Create(ctx, pod); err != nil {
			fmt.Printf("unable to create pod due to %v\n", err)
			return err
		}
	}

	// provision secrets
	var s v1.Secret
	key = types.NamespacedName{Namespace: nsName, Name: nstObj.Spec.AddResources.Secret.Name}
	if err := r.Get(ctx, key, &s); err != nil {
		fmt.Printf("unable to get secret info due to %v\n", err)
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nsName,
				Name:      nstObj.Spec.AddResources.Secret.Name,
			},
			Data:       nstObj.Spec.AddResources.Secret.Data,
			StringData: nstObj.Spec.AddResources.Secret.StringData,
		}
		if err := r.Create(ctx, secret); err != nil {
			fmt.Printf("unable to create secret due to %v\n", err)
			return err
		}
	}

	// provision LimitRange
	var lr v1.LimitRange
	key = types.NamespacedName{Namespace: nsName, Name: nstObj.Spec.AddResources.LimitRange.Name}
	if err := r.Get(ctx, key, &lr); err != nil {
		fmt.Printf("unable to get LimitRange info due to %v\n", err)
		limitrange := &v1.LimitRange{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: nsName,
				Name:      nstObj.Spec.AddResources.LimitRange.Name,
			},
			Spec: nstObj.Spec.AddResources.LimitRange.Spec,
		}
		if err := r.Create(ctx, limitrange); err != nil {
			fmt.Printf("unable to create Limitrange due to %v\n", err)
			return err
		}
	}

	return nil
}

// executePostCreateHook runs the postCreate lifecycle hooks and returns the duration it takes.
// add the "-n", namespace.Name suffix here.
func executePostCreateHook(namespace v1.Namespace, gracePeriod int64, nst megav1.NamespaceTemplate) int64 {
	fmt.Printf("Running postCreate hook for namespace %q\n", namespace.Name)

	nst.Spec.PostCreateHook.Command = append(nst.Spec.PostCreateHook.Command, "-n", namespace.Name)

	start := metav1.Now()
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer utilruntime.HandleCrash()

		fmt.Printf("running the postCreateHook %v in a go routine\n", nst.Spec.PostCreateHook.Command)
		n := len(nst.Spec.PostCreateHook.Command)

		cmd := exec.Command(nst.Spec.PostCreateHook.Command[0], nst.Spec.PostCreateHook.Command[1:n]...)
		if err := cmd.Run(); err != nil {
			fmt.Printf("executePostCreateHook failed to run the hook due to %v\n", err)
		}

	}()

	select {
	case <-time.After(time.Duration(gracePeriod) * time.Second):
		fmt.Printf("postCreate hook for namespace %q did not complete in %d seconds\n", namespace.Name, gracePeriod)
	case <-done:
		fmt.Printf("postCreate hook for namespace %q completed\n", namespace.Name)
	}

	return int64(metav1.Now().Sub(start.Time).Seconds())
}

// +kubebuilder:rbac:groups=mega.aragunathan.com,resources=namespacetemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mega.aragunathan.com,resources=namespacetemplates/status,verbs=get;update;patch

func (r *NamespaceTemplateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()
	log := r.Log.WithValues("namespacetemplate", req.NamespacedName)

	// 1. Load NamespaceTemplate obj by name
	var nstObj megav1.NamespaceTemplate
	if err := r.Get(ctx, req.NamespacedName, &nstObj); err != nil {
		log.Error(err, "unable to fetch NamespaceTemplate")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	actualObj := nstObj.DeepCopy()
	fmt.Printf("Options %v\n", actualObj.Spec.Options)
	fmt.Printf("PrecreateHook: %v\n", actualObj.Spec.PreCreateHook)
	fmt.Printf("PostcreateHook %v\n", actualObj.Spec.PostCreateHook)
	fmt.Printf("AdditionalResources %+v\n", actualObj.Spec.AddResources)

	// 2. List all namespaces that belong to this namespacetemplate
	var namespaces v1.NamespaceList
	if err := r.List(ctx, &namespaces, client.InNamespace(req.Namespace), client.MatchingLabels{"namespacetemplate": req.Name}); err != nil {
		log.Error(err, "unable to list namespaces matching the nst in this request")
		return ctrl.Result{}, err
	}

	for _, ns := range namespaces.Items {
		fmt.Printf("namespace matching this request: %v\n", ns.ObjectMeta.Name)
	}

	// 3. For all namespaces matching this nst, create additional resources
	for _, ns := range namespaces.Items {
		fmt.Printf("creating additionalresources for namespace %v\n", ns.ObjectMeta.Name)
		if err := r.createadditionalresources(ctx, ns.ObjectMeta.Name, nstObj); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 4. For all namespaces matching this nst, run the postcreatehook for the nst
	// Check r.PostCreateHookRan to run this just once.
	for _, ns := range namespaces.Items {
		fmt.Printf("running postCreateHooks for %v\n", ns.Name)
		_, ok := r.PostCreateHookMap[ns.Name]
		if !ok {
			fmt.Printf("no entry for %v in PostCreateHookMap\n", ns.Name)
			executePostCreateHook(ns, int64(time.Second*1), nstObj)
			r.PostCreateHookMap[ns.Name] = true
		} else {
			fmt.Printf("there's an entry for %v in PostCreateHookMap\n", ns.Name)
		}
	}

	// requeue so that we can continually monitor for any ns that match nst
	return ctrl.Result{Requeue: true}, nil
}

func (r *NamespaceTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&megav1.NamespaceTemplate{}).Complete(r)
}
