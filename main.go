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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var inCluster bool
	var deleteFailedAfter int
	var deleteSuccessfulAfter int
	var namespaces []string
	var clientset *kubernetes.Clientset

	flag.BoolVar(&inCluster, "in-cluster", true, "in cluster config or ~/.kubeconfig")
	flag.IntVar(&deleteFailedAfter, "delete-failed-after", 10, "delete failed pods after x days")
	flag.IntVar(&deleteSuccessfulAfter, "delete-successful-after", 5, "delete succesful pods after x days")
	flag.StringSliceVar(&namespaces, "namespaces", []string{"default"}, "namespaces to delete pods from")

	flag.Parse()

	if inCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Println(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Println(err.Error())
		}
	}

	if !inCluster {
		kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Println(err.Error())
		}

		// create the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Println(err.Error())
		}

	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(fmt.Sprintf("Deleting pods that have failed after %d days and succeeded after %d days", deleteFailedAfter, deleteSuccessfulAfter))

	var wg sync.WaitGroup

	for _, v := range pods.Items {
		for k, v2 := range v.Labels {
			if k == "workflows.argoproj.io/completed" && v2 == "true" {
				if v.Status.Phase == "Succeeded" && v.CreationTimestamp.Time.Before(time.Now().AddDate(0, 0, -deleteFailedAfter)) {
					for _, v3 := range namespaces {
						if v.Namespace == v3 {

							wg.Add(1)

							go func() {
								defer wg.Done()
								deletePod(v.Namespace, v.Name, clientset)
							}()

							wg.Wait()
						}
					}
				}

				if v.Status.Phase == "Failed" && v.CreationTimestamp.Time.Before(time.Now().AddDate(0, 0, -deleteSuccessfulAfter)) {
					for _, v3 := range namespaces {
						if v.Namespace == v3 {
							wg.Add(1)

							go func() {
								defer wg.Done()
								deletePod(v.Namespace, v.Name, clientset)
							}()

							wg.Wait()
						}
					}
				}
			}
		}
	}

}

func deletePod(namespace string, podName string, clientset *kubernetes.Clientset) {

	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", true, "dry run deletion")

	if dryRun {
		m := fmt.Sprintf("Would delete pod %s in namespace %s", podName, namespace)
		log.Println(m)

		return
	}

	m := fmt.Sprintf("Deleting pod %s in namespace %s", podName, namespace)
	log.Println(m)

	err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}

	m = fmt.Sprintf("Deleted pod %s in namespace %s", podName, namespace)
	log.Println(m)
}
