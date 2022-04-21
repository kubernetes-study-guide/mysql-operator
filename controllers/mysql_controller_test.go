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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	wordpressv1 "github.com/Sher-Chowdhury/mysql-operator/api/v1"
)

type mockKubeClient struct {
	getCallCounter      int
	failOnGetCallNumber int
	getResponse         error

	createCallCounter      int
	failOnCreateCallNumber int
	createResponse         error

	deleteCallCounter        int
	deleteOnCreateCallNumber int
	deleteResponse           error

	updateCallCounter        int
	updateOnCreateCallNumber int
	updateResponse           error

	patchCallCounter        int
	patchOnCreateCallNumber int
	patchResponse           error

	deleteAllOfresponseCallCounter        int
	deleteAllOfresponseOnCreateCallNumber int
	deleteAllOfResponse                   error
}

// Need to implement all of these methods - https://github.com/kubernetes-sigs/controller-runtime/blob/eb39b8eb28cfe920fa2450eb38f814fc9e8003e8/pkg/client/interfaces.go
func (m *mockKubeClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	// this counter system is used because the create method is called
	// multiple times in the
	m.createCallCounter = m.createCallCounter + 1
	if m.failOnCreateCallNumber == m.createCallCounter {
		return m.createResponse
	}
	return nil
}
func (m *mockKubeClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return m.deleteResponse
}

func (m *mockKubeClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return m.updateResponse
}

func (m *mockKubeClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return m.patchResponse
}

func (m *mockKubeClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return m.deleteAllOfResponse
}

func (m *mockKubeClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	// this counter system is used because the Get method is called
	// multiple times in the controller
	m.getCallCounter = m.getCallCounter + 1
	if m.failOnGetCallNumber == m.getCallCounter {
		return m.getResponse
	}
	return nil
}

func (m *mockKubeClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return nil
}

func (m *mockKubeClient) RESTMapper() meta.RESTMapper {
	return nil
}

func (m *mockKubeClient) Scheme() *runtime.Scheme {
	return nil
}

func (m *mockKubeClient) Status() client.StatusWriter {
	return nil
}

func TestReconcile(t *testing.T) {
	mysqlCrName := "test-name"
	mysqlCrNamespace := "test-namespace"

	// A mysql object with metadata and spec.
	mysqlCR := &wordpressv1.Mysql{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlCrName,
			Namespace: mysqlCrNamespace,
		},
		Spec:   wordpressv1.MysqlSpec{},
		Status: wordpressv1.MysqlStatus{},
	}

	mysqlDep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlCrName,
			Namespace: mysqlCrNamespace,
		},
		Spec:   appsv1.DeploymentSpec{},
		Status: appsv1.DeploymentStatus{},
	}

	scheme := runtime.NewScheme()
	wordpressv1.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      mysqlCrName,
			Namespace: mysqlCrNamespace,
		},
	}

	t.Run("Success Scenarios", func(t *testing.T) {
		t.Run("CR successfully reconciled", func(t *testing.T) {

			cr := mysqlCR
			dep := mysqlDep
			// Create a fake client to mock API calls.
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr, dep).Build()

			reconciler := &MysqlReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}
			_, err := reconciler.Reconcile(context.TODO(), req)
			require.NoError(t, err)

		})

		t.Run("errors.IsNotFound is encountered if cr missing", func(t *testing.T) {

			fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

			reconciler := &MysqlReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}
			_, err := reconciler.Reconcile(context.TODO(), req)
			require.NoError(t, err)

		})

		t.Run("errors.IsNotFound encountered if deployment is missing", func(t *testing.T) {

			cr := mysqlCR
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr).Build()

			reconciler := &MysqlReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}
			_, err := reconciler.Reconcile(context.TODO(), req)
			require.NoError(t, err)

		})

	})

	t.Run("Failure Scenarios", func(t *testing.T) {

		t.Run("expect error if kubeclient returns generic error when querying cluster for cr", func(t *testing.T) {

			// fakeclient can't return custom error so have to use mockclient instead
			mockKubeClient := &mockKubeClient{
				getCallCounter:      0,
				failOnGetCallNumber: 1,
				getResponse:         errors.NewBadRequest("Something went wrong"),
			}
			reconciler := &MysqlReconciler{
				Client: mockKubeClient,
				Scheme: scheme,
			}

			_, err := reconciler.Reconcile(context.TODO(), req)
			require.Error(t, err)

		})

		t.Run("expect error if kubeclient returns generic error when querying cluster for cr", func(t *testing.T) {

			// fakeclient can't return custom error so have to use mockclient instead
			mockKubeClient := &mockKubeClient{
				getCallCounter:      0,
				failOnGetCallNumber: 1,
				getResponse:         errors.NewBadRequest("Something went wrong"),
			}
			reconciler := &MysqlReconciler{
				Client: mockKubeClient,
				Scheme: scheme,
			}

			_, err := reconciler.Reconcile(context.TODO(), req)
			require.Error(t, err)

		})

		t.Run("expect error if kubeclient returns generic error when GETting deployment", func(t *testing.T) {

			// fakeclient can't return custom error so have to use mockclient instead
			mockKubeClient := &mockKubeClient{
				getCallCounter:      0,
				failOnGetCallNumber: 2,
				getResponse:         errors.NewBadRequest("Something went wrong"),
			}
			reconciler := &MysqlReconciler{
				Client: mockKubeClient,
				Scheme: scheme,
			}

			_, err := reconciler.Reconcile(context.TODO(), req)
			require.Error(t, err)

		})

		t.Run("expect error if kubeclient returns generic error when creating deployment", func(t *testing.T) {

			// fakeclient can't return custom error so have to use mockclient instead
			mockKubeClient := &mockKubeClient{
				getCallCounter:         0,
				failOnGetCallNumber:    2,
				getResponse:            errors.NewNotFound(schema.GroupResource{}, "Deployment"),
				createCallCounter:      0,
				failOnCreateCallNumber: 1,
				createResponse:         errors.NewBadRequest("Something went wrong"),
			}
			reconciler := &MysqlReconciler{
				Client: mockKubeClient,
				Scheme: scheme,
			}

			_, err := reconciler.Reconcile(context.TODO(), req)
			require.Error(t, err)

		})
	})
}

func TestSetupWithManager(t *testing.T) {
	scheme := runtime.NewScheme()
	wordpressv1.AddToScheme(scheme)

	config := &rest.Config{
		Host: "https://example.com:443",
	}

	options := ctrl.Options{
		Scheme: scheme,
	}

	mgr, err := ctrl.NewManager(config, options)
	require.NoError(t, err)

	err = (&MysqlReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
		Scheme: scheme,
	}).SetupWithManager(mgr)

	require.NoError(t, err)
}

func TestDeploymentForMysql(t *testing.T) {

	scheme := runtime.NewScheme()
	wordpressv1.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)

	// A mysql object with metadata and spec.
	mysqlCR := &wordpressv1.Mysql{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-Namespace",
		},
		Spec:   wordpressv1.MysqlSpec{},
		Status: wordpressv1.MysqlStatus{},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects().Build()

	reconciler := &MysqlReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	res := reconciler.deploymentForMysql(mysqlCR)
	require.IsType(t, &appsv1.Deployment{}, res)
	assert.Equal(t, res.Name, "test-name-mysql")
}
