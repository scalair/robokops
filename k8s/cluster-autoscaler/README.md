# Cluster Autoscaler
Cluster Autoscaler is a tool that automatically adjusts the size of the Kubernetes cluster when one of the following conditions is true:
* there are pods that failed to run in the cluster due to insufficient
  resources,
* there are nodes in the cluster that have been underutilized for an extended period of time and their pods can be placed on other existing nodes.

More info [here](https://github.com/helm/charts/tree/master/stable/cluster-autoscaler)