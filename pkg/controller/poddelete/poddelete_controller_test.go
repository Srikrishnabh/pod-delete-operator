package poddelete

import (
	poddeletev1alpha1 "github.com/Srikrishnabh/pod-delete-operator/pkg/apis/poddelete/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestReconcilePodDelete_Reconcile(t *testing.T) {

	podDeleteInstance := &poddeletev1alpha1.PodDelete{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name: "delete-pod-test",
		},
		Spec:       poddeletev1alpha1.PodDeleteSpec{
			Namespace: "default",
			PodName: "nginx",
		},
		Status:     poddeletev1alpha1.PodDeleteStatus{},
	}

	podToDeleteInstance := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name: "nginx",
		},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}

	objs := []runtime.Object{podDeleteInstance, podToDeleteInstance}

	s := scheme.Scheme
	s.AddKnownTypes(poddeletev1alpha1.SchemeGroupVersion, podDeleteInstance)
	s.AddKnownTypes(corev1.SchemeGroupVersion, podToDeleteInstance)

	c1 := fake.NewFakeClientWithScheme(s, objs...)

	r := ReconcilePodDelete{client:c1, scheme:s}

	// posistive case to delete the pod
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "default",
			Name:      "nginx",
		},
	}

	res, err := r.Reconcile(req)

	assert.Nil(t, err, "reconcile error for pod delete valid")
	assert.Equal(t, false, res.Requeue)


	//negative case delete pod which dont exists
	reqPodNotFound := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "nginx",
			Name:      "nginx",
		},
	}

	res, err = r.Reconcile(reqPodNotFound)

	assert.Nil(t, err, "reconcile error for delete pod thats not found")
	assert.Equal(t, false, res.Requeue)

}
