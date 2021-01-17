package poddelete

import (
	"context"
	"fmt"

	poddeletev1alpha1 "github.com/Srikrishnabh/pod-delete-operator/pkg/apis/poddelete/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_poddelete")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PodDelete Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePodDelete{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("poddelete-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource PodDelete
	err = c.Watch(&source.Kind{Type: &poddeletev1alpha1.PodDelete{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner PodDelete
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &poddeletev1alpha1.PodDelete{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePodDelete implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePodDelete{}

// ReconcilePodDelete reconciles a PodDelete object
type ReconcilePodDelete struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a PodDelete object and makes changes based on the state read
// and what is in the PodDelete.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePodDelete) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PodDelete")

	// Fetch the PodDelete instance
	instance := &poddeletev1alpha1.PodDelete{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//if err := deletePod(instance.Spec.Namespace, instance.Spec.PodName); err != nil {
	//	return reconcile.Result{RequeueAfter: time.Second*10}, err
	//}

	//Todo: ignore self pod deletion via env variable

	if err := r.deletePod(instance.Spec.Namespace, instance.Spec.PodName); err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info(fmt.Sprintf("reconcile for delete pod %s:%s success", instance.Spec.Namespace,																						instance.Spec.PodName))
	return reconcile.Result{}, nil
}


func (r *ReconcilePodDelete) deletePod(namespace, podName string) error {
	pod := &corev1.Pod{}
	err := r.client.Get(context.TODO(), client.ObjectKey{
		Namespace: namespace,
		Name: podName,
	}, pod)

	if err != nil {
		if errors.IsNotFound(err) { // ignore if pod not found
			logf.Log.Error(err, "error")
			logf.Log.Info(fmt.Sprintf("pod %s:%s not found", namespace, podName))
			return nil
		}
		return err
	}

	//Dont delete the pod that is in terminating state
	if !pod.DeletionTimestamp.IsZero(){
		logf.Log.Info(fmt.Sprintf("pod %s:%s is in terminating state", namespace, podName))
		return nil
	}

	err = r.client.Delete(context.TODO(), pod)
	if err != nil && !errors.IsNotFound(err){ // ignore if pod not found
		logf.Log.Error(err, fmt.Sprintf("failed to delete the pod %s:%s", namespace, podName))
		return err
	}

	return nil
}

//func deletePod(namespace, podName string) error {
//	k8sClient := k8sclientset.GetClient()
//
//	getOptions := metav1.GetOptions{}
//	pod, err := k8sClient.CoreV1().Pods(namespace).Get(podName, getOptions)
//	if err != nil {
//		return err
//	}
//
//	// Check if pod is already in terminating state
//	if !pod.DeletionTimestamp.IsZero(){
//		logf.Log.Info(fmt.Sprintf("pod %s:%s is in terminating state", namespace, podName))
//		return nil
//	}
//
//	deleteOptions := &metav1.DeleteOptions{}
//	err = k8sClient.CoreV1().Pods(namespace).Delete(podName, deleteOptions)
//	if err != nil && !errors.IsNotFound(err){ // ignore if pod not found
//		logf.Log.Error(err, fmt.Sprintf("failed to delete the pod %s:%s", namespace, podName))
//		return err
//	}
//
//	return nil
//}
