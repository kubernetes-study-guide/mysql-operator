# mysql-operator
mysql kubernetes operator built using the operator-sdk

```
brew install operator-sdk
```

```
$ operator-sdk version
operator-sdk version: "v0.15.1", commit: "e35ec7b722ba095e6438f63fafb9e7326870b486", go version: "go1.13.6 darwin/amd64"
```

```
operator-sdk new mysql-operator --repo=github.com/Sher-Chowdhury/mysql-operator
cd mysql-operator
```


```
operator-sdk add api --api-version=cache.codingbee.net/v1alpha1 --kind=MySQL  # kind needs to start with uppercase
```

This would have created a crd and a sample cr file, that you can try deploying at this stage. They look like this:

```
$ cat deploy/crds/cache.codingbee.net_mysqls_crd.yaml 
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mysqls.cache.codingbee.net
spec:
  group: cache.codingbee.net
  names:
    kind: MySQL
    listKind: MySQLList
    plural: mysqls
    singular: mysql
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: MySQL is the Schema for the mysqls API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: MySQLSpec defines the desired state of MySQL
          type: object
        status:
          description: MySQLStatus defines the observed state of MySQL
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
```

and the sample custom resource (cr):

```
$ cat deploy/crds/cache.codingbee.net_v1alpha1_mysql_cr.yaml
apiVersion: cache.codingbee.net/v1alpha1
kind: MySQL
metadata:
  name: example-mysql
spec:
  # Add fields here
  size: 3
```

Added custom setting to pkg/apis/cache/v1alpha1/mysql_types.go. 

For these changes to get propated to the rest of the repo, run:

```
operator-sdk generate k8s
```

Next run the following to get the changes reflected in the crd file:

```
operator-sdk generate crds
```

Now let's create our operator's controller:

```
operator-sdk add controller --api-version=cache.codingbee.net/v1alpha1 --kind=MySQL
```
This ends up creating the file `pkg/controller/add_mysql.go` and the folder `pkg/controller/mysql/` along with all it's content.


Next update controller to make use of the mysql env vars - https://github.com/Sher-Chowdhury/mysql-operator/commit/2ca59bc874221c67080cb8094952d3282ed1ba58


Now deploy the crd (you can also deploy the example cr too if you want too):

```
kubectl apply -f deploy/crds/cache.codingbee.net_mysqls_crd.yaml
```




Now deploy the operator. Theres 2 ways to do that. deploy it as a pod, or run it locally. 

### Deploy operator as a container

First build an image that has your controller baked in:

```
account=sher_chowdhury
image_name=mysql-operator
tag_version=v0.0.1
docker login quay.io
operator-sdk build quay.io/${account}/${image_name}:${tag_version}
docker push quay.io/${account}/${image_name}:${tag_version}
sed -i "" "s|REPLACE_IMAGE|quay.io/${account}/${image_name}:${tag_version}|g" deploy/operator.yaml
```



```
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/operator.yaml
```







# References

https://github.com/operator-framework/operator-sdk/blob/v0.15.1/doc/user-guide.md