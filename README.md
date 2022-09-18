# aerospike
Exercise for Aerospike interview

## Table of Contents

- [What's included](#whats-included)
- [How to use it](#how-to-use-it)

### What's included

The `main` package contains the exercise that does the following:
* connect to the k8s cluster
* print out the namespaces on the cluster
* create a new namespace
* create a pod in that namespace that runs a simple hello-world container
* print out pod names and the namespace they are in for any pods that have a label of ‘k8s-app=kube-dns’ or a similar label is ok as well
* delete the hello-world pod created from above
* extra credit - show how an client-go informer works


### How to use it
Run it with defaults:
```
% go run . 
```
Get usage:
```
% go run . -help                                                                        
Usage :

  -kubeconfig string
      (optional) absolute path to the kubeconfig file (default "/Users/dev/.kube/config")
  -labelSelector string
      (optional) labelSelector (default "k8s-app=kube-dns")
  -namespace string
      (optional) namespace (default "aerospike")
  -podName string
      (optional) podName (default "helloworldpod")
```
Example run with flags:
```
% go run . -namespace aerospike -podName helloaerospike -labelSelector=k8s-app=kube-dns

  0. exit
  1. connect to the k8s cluster
  2. print out the namespaces on the cluster
  3. create a new namespace
  4. create a pod in that namespace that runs a simple hello-world container
  5. print out pod names and the namespace they are in for any pods that have a label of ‘k8s-app=kube-dns’
  6. delete the hello-world pod created from above
  7. extra credit - show how an client-go informer works
  8. run all
Please enter a number from 0 to 8 : 8
Running all 1-7:
Connect
ListNameSpaces:
default
kube-node-lease
kube-public
kube-system
spike
CreateNamespace:
Created Namespace aerospike on 2022-09-18 16:42:41 +0200 CEST
CreatePod:
Created Podname:helloaerospike in Namespace:aerospike
PrintPodInfo:
Podname:coredns-95db45d46-6dkbb, Namespace:kube-system
Podname:coredns-95db45d46-bwspx, Namespace:kube-system
TestInformer:
 podItem pod ==  {helloaerospike  aerospike  d03eef22-7c9c-4e47-8322-5c43357d7235 250030 0 2022-09-18 16:42:41 +0200 CEST <nil> <nil> map[] map[] [] [] [{aerospike Update v1 2022-09-18 16:42:41 +0200 CEST FieldsV1 {"f:spec":{"f:containers":{"k:{\"name\":\"main\"}":{".":{},"f:args":{},"f:command":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:enableServiceLinks":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}} } {kubelet Update v1 2022-09-18 16:42:41 +0200 CEST FieldsV1 {"f:status":{"f:conditions":{"k:{\"type\":\"ContainersReady\"}":{".":{},"f:lastProbeTime":{},"f:lastTransitionTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Initialized\"}":{".":{},"f:lastProbeTime":{},"f:lastTransitionTime":{},"f:status":{},"f:type":{}},"k:{\"type\":\"Ready\"}":{".":{},"f:lastProbeTime":{},"f:lastTransitionTime":{},"f:message":{},"f:reason":{},"f:status":{},"f:type":{}}},"f:containerStatuses":{},"f:hostIP":{},"f:startTime":{}}} status}]}
DeletePod:
Deleted Podname:helloaerospike in Namespace:aerospike

  0. exit
  1. connect to the k8s cluster
  2. print out the namespaces on the cluster
  3. create a new namespace
  4. create a pod in that namespace that runs a simple hello-world container
  5. print out pod names and the namespace they are in for any pods that have a label of ‘k8s-app=kube-dns’
  6. delete the hello-world pod created from above
  7. extra credit - show how an client-go informer works
  8. run all
Please enter a number from 0 to 8 :
```

