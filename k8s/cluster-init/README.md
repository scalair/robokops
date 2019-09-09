# Cluster Init
Initialise the cluster:
* Create rbac `admin` role for helm in `kube-system` namespace
* Initialise Helm

## Limit resources usage
To limit resources usage in the kube-system namespace, you have to create the following file in your conf folder:
```
cluster-init/manifests/kube-system-limit-range.yaml
```
Here is an example:
```
apiVersion: v1
kind: LimitRange
metadata:
  name: kube-system-limit-range
  namespace: kube-system
spec:
  limits:
  - default:
      memory: 128Mi
      cpu: 200m
    defaultRequest:
      memory: 128Mi
      cpu: 200m
    type: Container
```

## Storage class
To create the storage class for your persistent volumes, you have to create the following file in your conf folder:
```
cluster-init/manifests/storage-class.yaml
```
Here is an example:
```
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ssd
provisioner: kubernetes.io/aws-ebs
allowVolumeExpansion: true
parameters:
  type: gp2
  fsType: ext4
```