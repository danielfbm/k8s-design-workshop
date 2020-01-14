# clientgo

All `client-go` and automatically generated client code does provide a interface called `Interface`  and a fake implemntation in a subpackage called `fake`.[check this](https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md) 

Here are some examples:

Kubernetes client-go [interface](https://godoc.org/k8s.io/client-go/kubernetes#Interface), [implementation](https://godoc.org/k8s.io/client-go/kubernetes#Clientset) , and [fake](https://godoc.org/k8s.io/client-go/kubernetes/fake#Clientset)


In its core, `fake` will use a response mechanism that can be heavily customized. Check the test cases for more details


