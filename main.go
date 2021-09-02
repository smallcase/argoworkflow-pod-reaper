package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var inCluster bool
	var deleteFailedAfter int
	var deleteSuccessfulAfter int
	var namespaces []string
	var dryRun bool

	flag.BoolVar(&dryRun, "dry-run", true, "dry run deletion")
	flag.BoolVar(&inCluster, "in-cluster", true, "in cluster config or ~/.kubeconfig")
	flag.IntVar(&deleteFailedAfter, "delete-failed-after", 10, "delete failed pods after x days")
	flag.IntVar(&deleteSuccessfulAfter, "delete-successful-after", 5, "delete succesful pods after x days")
	flag.StringSliceVar(&namespaces, "namespaces", []string{"default"}, "namespaces to delete pods from")

	flag.Parse()

	var clientset *kubernetes.Clientset

	if inCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Println(err.Error())
			return
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if !inCluster {
		kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// create the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Println(err.Error())
			return
		}

	}

	Reap(deleteFailedAfter, deleteSuccessfulAfter, namespaces, clientset.CoreV1(), dryRun)
}

func Reap(deleteFailedAfter int, deleteSuccessfulAfter int, namespaces []string, clientset v1.CoreV1Interface, dryRun bool) {

	pods, err := clientset.Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(fmt.Sprintf("Deleting pods that failed %d days ago", deleteFailedAfter))
	log.Println(fmt.Sprintf("Deleting pods that succeeded %d days ago", deleteSuccessfulAfter))

	var wg sync.WaitGroup

	for _, v := range pods.Items {
		for k, v2 := range v.Labels {
			if k == "workflows.argoproj.io/completed" && v2 == "true" {
				if v.CreationTimestamp.Time.Before(time.Now().AddDate(0, 0, -deleteFailedAfter)) {
					for _, v3 := range namespaces {
						if v.Namespace == v3 {

							wg.Add(1)

							go func() {
								defer wg.Done()
								deletePod(v.Namespace, v.Name, clientset, dryRun)
							}()

							wg.Wait()
						}
					}
				}
			}
			if k == "workflows.argoproj.io/completed" && v2 == "false" {
				if v.CreationTimestamp.Time.Before(time.Now().AddDate(0, 0, -deleteSuccessfulAfter)) {
					for _, v3 := range namespaces {
						if v.Namespace == v3 {
							wg.Add(1)

							go func() {
								defer wg.Done()
								deletePod(v.Namespace, v.Name, clientset, dryRun)
							}()

							wg.Wait()
						}
					}
				}
			}
		}
	}

}

func deletePod(namespace string, podName string, clientset v1.CoreV1Interface, dryRun bool) {

	if dryRun {
		m := fmt.Sprintf("Would delete pod %s in namespace %s", podName, namespace)
		log.Println(m)

		return
	}

	m := fmt.Sprintf("Deleting pod %s in namespace %s", podName, namespace)
	log.Println(m)

	err := clientset.Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	m = fmt.Sprintf("Deleted pod %s in namespace %s", podName, namespace)
	log.Println(m)
}
