# Update worker nodes

In order to update worker nodes, you need to update their image in the autoscaling group and slowly taint/drain the node to use the new one.
This is what you have to do:
```
# Update ASG

# Scale down the cluster-autoscaler
kubectl -n kube-system scale deployment.apps/cluster-autoscaler-aws-cluster-autoscaler --replicas=0

# Double the number of nodes in the ASG (half of them will use the new image but have no pods, the other half use the old image and have all the pod)

# Taint all the old nodes so no pod get schedule on them (replace the version with the one you update from)
K8S_VERSION=1.11.5
nodes=$(kubectl get nodes -o jsonpath="{.items[?(@.status.nodeInfo.kubeletVersion==\"v$K8S_VERSION\")].metadata.name}")
for node in ${nodes[@]} ; do echo "Tainting $node" ; kubectl taint nodes $node key=value:NoSchedule ; done

# Drain the nodes so all the pod get migrated to new nodes
K8S_VERSION=1.11.5
nodes=$(kubectl get nodes -o jsonpath="{.items[?(@.status.nodeInfo.kubeletVersion==\"v$K8S_VERSION\")].metadata.name}")
for node in ${nodes[@]}; do echo "Draining $node" ; kubectl drain $node --ignore-daemonsets --delete-local-data ; done

# Delete the nodes with the previous image

# Scale back up the cluster-autoscaler. It will then remove the unused nodes slowly
kubectl -n kube-system scale deployment.apps/cluster-autoscaler-aws-cluster-autoscaler --replicas=1
```