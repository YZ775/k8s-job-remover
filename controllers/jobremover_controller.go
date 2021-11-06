/*
Copyright 2021.

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
	"time"

	removerv1 "github.com/Onikle/job-remover/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// JobRemoverReconciler reconciles a JobRemover object
type JobRemoverReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=remover.onikle.com,resources=jobremovers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=remover.onikle.com,resources=jobremovers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=remover.onikle.com,resources=jobremovers/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JobRemover object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *JobRemoverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here
	var jobremover removerv1.JobRemover
	err := r.Get(ctx, req.NamespacedName, &jobremover)

	// if errors.IsNotFound(err) {
	// 	r.removeMetrics(jobremover)
	// 	return ctrl.Result{}, nil
	// }
	if err != nil {
		logger.Error(err, "unable to get JobRemover", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if !jobremover.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	ns := jobremover.Spec.Namespace
	var jobs batchv1.JobList
	err = r.List(ctx, &jobs, &client.ListOptions{Namespace: ns})

	if err != nil {
		return ctrl.Result{}, err
	}
	for _, job := range jobs.Items {
		fmt.Println(job.Name)
		time_elapsed := time.Since(job.Status.CompletionTime.Time)

		if job.Status.Succeeded == 1 && time_elapsed.Minutes() > float64(jobremover.Spec.TTL) {
			r.Delete(ctx, &job)
			logger.Info("DeletedJob", "Jobname", job.Name, "completed", job.Status.Succeeded, "CompletionTime", job.Status.CompletionTime, "Elapsed", time_elapsed)
		}

	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JobRemoverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&removerv1.JobRemover{}).
		Complete(r)
}
