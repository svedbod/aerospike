package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	var namespaceArg *string
	var podNameArg *string
	var labelSelectorArg *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	namespaceArg = flag.String("namespace", "aerospike", "(optional) namespace")
	podNameArg = flag.String("podName", "helloworldpod", "(optional) podName")
	labelSelectorArg = flag.String("labelSelector", "k8s-app=kube-dns", "(optional) labelSelector")

	flag.Parse()
	var namespace string = *namespaceArg
	var podName string = *podNameArg
	var labelSelector string = *labelSelectorArg

	var clientset kubernetes.Interface

    for {

		fmt.Println(`
	0. exit
	1. connect to the k8s cluster
	2. print out the namespaces on the cluster
	3. create a new namespace
	4. create a pod in that namespace that runs a simple hello-world container
	5. print out pod names and the namespace they are in for any pods that have a label of ‘k8s-app=kube-dns’
	6. delete the hello-world pod created from above
	7. extra credit - show how an client-go informer works
	8. run all`)

        fmt.Print("Please enter a number from 0 to 8 : ")
	    var num int
	    fmt.Scanf("%d", &num)
	    
	    switch num {
	    case 0:
	        os.Exit(0)
	    case 1:
			clientset = connect(kubeconfig)
	    case 2:
	        listNameSpaces(clientset)
	    case 3:
	        createNamespace(clientset, namespace)
	    case 4:
			createPod(clientset, namespace, podName)
	    case 5:
			printPodInfo(clientset, labelSelector)
	    case 6:
	        deletePod(clientset,        namespace, podName)
	    case 7:
	        testInformer(clientset, namespace, podName)
	    case 8:
	        fmt.Println("Running all 1-7:")
	        clientset = connect(kubeconfig)
	        listNameSpaces(clientset)
	        createNamespace(clientset, namespace)
	        createPod(clientset, namespace, podName)
	        printPodInfo(clientset, labelSelector)
	        testInformer(clientset, namespace, podName)
	       	deletePod(clientset, namespace, podName)

	    default:
	        fmt.Println("Please enter a valid action")
	    }
    }                                      
}