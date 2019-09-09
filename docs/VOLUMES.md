# Volumes

## Resize a persistent volumes
This section describe how to resize a persistent volume. More info [here](https://kubernetes.io/blog/2018/07/12/resizing-persistent-volumes-using-kubernetes/)

### Prerequise
The storage class used to create the volume must have been created with `allowVolumeExpansion: true`

### Steps
If you created the PVC through a STS, you should first update the STS to reflate the reality.
At the moment, Helm will not let you directly update the size of the volume, so you have to recreate the STS to update it:
```
kubectl delete sts --cascade=false [your_sts]
# recreate the sts, either directly or using helm upgrade
```
Then, resize the PVC:
```
kubectl edit pvc [your_pvc] # edit the spec.resources.requests.storage
```
Wait for the cloud provider to resize the volume.
Then, you must recreate the pod that is using that volume, you can just delete it amd it will be recreated by the STS:
```
kubectl delete pod [your_pod]
```
