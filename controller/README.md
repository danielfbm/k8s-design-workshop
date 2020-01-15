# ConfigMapReplica 

Sample controller for TDD study

The state of the code repository can be navigated using tags:

- step-1
- step-2 

etc.

## Step 1


Create base code 

- Install [kubebuilder](https://book.kubebuilder.io/quick-start.html#installation)
- Change below `domain`, `license`, `owner`, and `repo` flags


```shell
kubebuilder init --domain example.com --license MIT  --repo github.com/danielfbm/k8s-design-workshop/controller
```


Create a resource `ConfigMapReplica`


```shell
kubebuilder create api --group replica --version v1alpha1 --kind ConfigMapReplica --namespaced=false --resource --controller --example 
```


