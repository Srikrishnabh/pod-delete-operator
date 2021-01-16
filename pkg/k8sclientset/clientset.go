package k8sclientset

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"sync"
)

var once sync.Once

var k8sClientSet *kubernetes.Clientset

func SetupClient(k8sConfig *rest.Config) {
	once.Do(func() {
		clientSet, err := kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			log.Fatal(err, "failed to create clientSet")
		}

		k8sClientSet = clientSet
	})
}

func GetClient() *kubernetes.Clientset {
	return k8sClientSet
}
