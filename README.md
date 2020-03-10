This is a proof-of-concept repo for a custom kubernetes controller that manages a new CRD called NamespaceTemplate.



Install the CRD
----------------

# install the crd
`make install` 

# verify that the CRD (NamespaceTemplate) got installed
```
$ kubectl get crd/namespacetemplates.mega.aragunathan.com
NAME                                      CREATED AT
namespacetemplates.mega.aragunathan.com   2020-03-10T20:16:49Z
```

# start the controller in a terminal
`make run`

Create objects for the NamespaceTemplate CRD
---------------------------------------------
# Use another terminal for kubectl client commands

# note that there are no objects for NamespaceTemplate yet
```
$ kubectl get namespacetemplates.mega.aragunathan.com
No resources found in default namespace.
```

# create an object for the NamespaceTemplate
```
$ kubectl apply -f config/samples/mega_v1_namespacetemplate.yaml 
namespacetemplate.mega.aragunathan.com/namespacetemplate-sample created
```

Create namespace using a NamespaceTemplate object
-------------------------------------------------
# create a namespace that uses the NamespaceTemplate. This is a label for the namespace.
```
$ kubectl apply -f config/samples/namespace.yaml
namespace/namespace-sample created
```

# notice how the controller provisions additionalresources, postStartHooks, etc that are 
# specified in the NamespaceTemplate
```
$ kubectl get pods,secrets,limitrange -n namespace-sample
NAME           READY   STATUS    RESTARTS   AGE
pod/test-pod   1/1     Running   0          20s

NAME                                      TYPE                                  DATA   AGE
secret/default-token-hjgbj                kubernetes.io/service-account-token   3      30s
secret/nginx-serviceaccount-token-29jlt   kubernetes.io/service-account-token   3      21s
secret/test-secret                        Opaque                                2      20s

NAME                         CREATED AT
limitrange/test-limitrange   2020-03-10T21:44:08Z
```

# create another namespace under the same NamespaceTemplate
```
$ kubectl apply -f config/samples/namespace2.yaml
namespace/namespace-sample2 created
```   

# again, how notice the resources mentioned in the NamespaceTemplate get provisioned in the namespace
```
$ kubectl get pods,secrets,limitrange -n namespace-sample2
NAME           READY   STATUS    RESTARTS   AGE
pod/test-pod   1/1     Running   0          34s

NAME                                      TYPE                                  DATA   AGE
secret/default-token-rgn66                kubernetes.io/service-account-token   3      5m55s
secret/nginx-serviceaccount-token-hj9wf   kubernetes.io/service-account-token   3      34s
secret/test-secret                        Opaque                                2      34s

NAME                         CREATED AT
limitrange/test-limitrange   2020-03-10T22:12:11Z
```

Update the NamespaceTemplate object and observe the Renciliation loop update the namespaces
---------------------------------------------------------------------------------------------
# Update the NST obj and see if the pod spec is updated in both namespaces
# change pod container port from 8080 to 80
```
$ kubectl apply -f config/samples/mega_v1_namespacetemplate.yaml 
namespacetemplate.mega.aragunathan.com/namespacetemplate-sample updated

$ kubectl get namespacetemplates.mega.aragunathan.com/namespacetemplate-sample -o=jsonpath={.spec.addresources.pod.spec}
map[containers:[map[image:nginx imagePullPolicy:Always name:web ports:[map[containerPort:80]]]]]

$ kubectl get pod/test-pod -o=jsonpath={.spec.containers[0].ports} -n namespace-sample
[map[containerPort:80 protocol:TCP]]

$ kubectl get pod/test-pod -o=jsonpath={.spec.containers[0].ports} -n namespace-sample2
[map[containerPort:80 protocol:TCP]]
```