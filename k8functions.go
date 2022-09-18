package main

import (
	"context"
	"fmt"
	"log"
    corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/informers"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

func connect(kubeconfig *string)*kubernetes.Clientset{
	fmt.Println("Connect");
	    // use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalln("Could not create kubernetes config")
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Could not to create a clientset")
		panic(err.Error())
	}
	return clientset
}

func getNamespaces(clientset kubernetes.Interface) []corev1.Namespace{
	if clientset == nil {
		fmt.Println("Please connect to a k8s cluster")
		return nil
	}
    nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
    if err != nil {
    	fmt.Println(err)
    }
    return nsList.Items 
}

func listNameSpaces(clientset kubernetes.Interface){
	fmt.Println("ListNameSpaces:")
	var nsItemList []corev1.Namespace
    nsItemList = getNamespaces(clientset)
    for _, n := range nsItemList {
        fmt.Println(n.Name)
    }
}

func createNamespace(clientset kubernetes.Interface, namespace string) bool{
	fmt.Println("CreateNamespace:")
	isCreated := false
	if clientset == nil {
		fmt.Println("Please connect to a k8s cluster")
	} else {
	    namespacesClient := clientset.CoreV1().Namespaces()
		existing, _ := namespacesClient.Get(context.Background(), namespace, metav1.GetOptions{})
		if existing == nil || namespace!=existing.ObjectMeta.Name {
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
				Status: corev1.NamespaceStatus{
					Phase: corev1.NamespaceActive,
				},
			}
			result, err := namespacesClient.Create(context.Background(), namespace, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			} else {
				isCreated = true
				fmt.Printf("Created Namespace %s on %s\n", result.ObjectMeta.Name, result.ObjectMeta.CreationTimestamp)
			}
	    } else {
	    	fmt.Printf("Namespace %s already exists\n",namespace)
	    }
    }
    return isCreated
}

func namespaceExists(clientset kubernetes.Interface, namespace string) bool{
	var nsItemList []corev1.Namespace
	var namespaceExists = false
    nsItemList = getNamespaces(clientset)
    for _, n := range nsItemList {
        if n.Name == namespace {
        	namespaceExists = true
        }
    }
    return namespaceExists
}

func podExists(clientset kubernetes.Interface, namespace string, podName string) bool {
	if clientset == nil {
		fmt.Println("Please connect to a k8s cluster")
		return false
	}
	pod, _ := clientset.CoreV1().Pods(namespace).Get(
	    context.Background(),
	    podName,
	    metav1.GetOptions{},
	)
	return pod.ObjectMeta.Name==podName;
}

func createPod(clientset kubernetes.Interface, namespace string, podName string){
	fmt.Println("CreatePod:")
	if ! namespaceExists(clientset, namespace) {
		fmt.Println("namespace does not exist")
	} else{
		if clientset == nil {
			fmt.Println("Please connect to a k8s cluster")
		} else {
			if ! podExists(clientset, namespace, podName) {
					pod := &corev1.Pod{
			    	ObjectMeta: metav1.ObjectMeta{Name: podName},
			    	Spec: corev1.PodSpec{
			        	RestartPolicy: corev1.RestartPolicyNever,
			        	Containers: []corev1.Container{
				            corev1.Container{
				                Name:    "main",
				                Image:   "python:3.8",
				                Command: []string{"python"},
				                Args:    []string{"-c", "print('hello world')"},
				            },
				        },
				    },
				}
				_, err := clientset.CoreV1().Pods(namespace).Create(
				    context.Background(),
				    pod,
				    metav1.CreateOptions{},
				)
				if err != nil {
					panic(err)
				} else {
					fmt.Printf("Created Podname:%s in Namespace:%s\n",podName, namespace)
				}
			} else {
				fmt.Printf("Did not create Podname:%s in Namespace:%s because it already exists\n",podName, namespace)
			}
		}
	}
}

func printPodInfo(clientset kubernetes.Interface, labelSelector string) {
	fmt.Println("PrintPodInfo:")
	if clientset == nil {
		fmt.Println("Please connect to a k8s cluster")
	} else {
		pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})

		for _, pod := range pods.Items {
	        fmt.Printf("Podname:%s, Namespace:%s\n",pod.Name, pod.Namespace)
	    }
	}
}

func deletePod(clientset kubernetes.Interface, namespace string, podName string){
	fmt.Println("DeletePod:")
	if clientset == nil {
		fmt.Println("Please connect to a k8s cluster")
	} else {
		if podExists(clientset, namespace, podName) {
			err := clientset.CoreV1().Pods(namespace).Delete(
			    context.Background(),
			    podName,
			    metav1.DeleteOptions{},
			)
			if err != nil {
				panic(err)
			} else {
				fmt.Printf("Deleted Podname:%s in Namespace:%s\n",podName, namespace)
			}
		} else {
			fmt.Printf("Could not delete Podname:%s in Namespace:%s because it does not exist\n",podName, namespace)
		}
	}
}

func onAdd(obj interface{}) {
	_ = obj.(*corev1.Pod)
}

func testInformer(clientset kubernetes.Interface, namespace string, podName string) {
	fmt.Println("TestInformer:")
	if clientset == nil {
		fmt.Println("Please connect to a k8s cluster")
	} else if ! podExists(clientset, namespace, podName){
		fmt.Println("Please create a pod ", podName )
	} else {
		// stop signal for the informer
		stopper := make(chan struct{})
		defer close(stopper)

		factory := informers.NewSharedInformerFactory(clientset, 0)
		podInformer := factory.Core().V1().Pods()
		informer := podInformer.Informer()

		defer runtime.HandleCrash()

		// start informer ->
		go factory.Start(stopper)

		// start to sync and call list
		if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
			runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
			return
		}

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: onAdd,  // register add eventhandler
			UpdateFunc: func(interface{}, interface{}) { fmt.Println("update not implemented") },
			DeleteFunc: func(interface{}) { fmt.Println("delete not implemented") },
		})

		podItem, _, _ := informer.GetIndexer().GetByKey(namespace + "/" + podName)
		fmt.Println(" podItem pod == ", podItem.(*corev1.Pod).ObjectMeta )
	}
}