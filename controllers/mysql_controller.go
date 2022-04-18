/*
Copyright 2022.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	wordpressv1 "github.com/Sher-Chowdhury/mysql-operator/api/v1"
)

// MysqlReconciler reconciles a Mysql object
type MysqlReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=wordpress.codingbee.net,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=wordpress.codingbee.net,resources=mysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=wordpress.codingbee.net,resources=mysqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Mysql object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx).WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

	// Fetch the mysql isntance
	cr := &wordpressv1.Mysql{}
	err := r.Get(ctx, req.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			// Requested object is not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get mysql cr, going to requeue and retry this")
		return ctrl.Result{Requeue: true}, err
	}

	// Check if the deployment already exists, if not, then create the new deployment.
	// https://sdk.operatorframework.io/docs/building-operators/golang/references/client/#example-usage
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define and create a new deployment.
			dep := r.deploymentForMysql(cr)
			if err = r.Create(ctx, dep); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	reqLogger.Info("Reconcile completed successfully")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wordpressv1.Mysql{}).
		Complete(r)
}

// deploymentForMemcached returns a Deployment object for data from m.
func (r *MysqlReconciler) deploymentForMysql(cr *wordpressv1.Mysql) *appsv1.Deployment {
	// lbls := labelsForApp(cr.Name)

	var replicas int32 = 1


	mysqlEnvVars := cr.Spec.Environment

	containerEnvVars := []corev1.EnvVar{
        {
            Name:  "MYSQL_ROOT_PASSWORD",
            Value: mysqlEnvVars.MysqlRootPassword,
        },
        {
            Name:  "MYSQL_DATABASE",
            Value: mysqlEnvVars.MysqlDatabase,
        },
        {
            Name:  "MYSQL_USER",
            Value: mysqlEnvVars.MysqlUser,
        },
        {
            Name:  "MYSQL_PASSWORD",
            Value: mysqlEnvVars.MysqlPassword,
        },
	}

	deploymentObject := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-msyql",
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "docker.io/mysql:latest",
						Name:  "wordpressMysql",
                        Env: containerEnvVars,
					}},
				},
			},
		},
	}

	// This sets this mysql cr as the owner of this deployment object. 
    // so that if cr is deleted then this deployment should also get deleted (garbage collected)
	controllerutil.SetControllerReference(cr, deploymentObject, r.Scheme)
	return deploymentObject
}
