#!/bin/bash
kubectl delete -f manifests/rbac-admin.yaml
kubectl delete deployment tiller-deploy --namespace kube-system
kubectl delete service tiller-deploy --namespace kube-system

if [ -f /conf/cluster-init/manifests/kube-system-limit-range.yaml ]; then
	kubectl delete -f /conf/cluster-init/manifests/kube-system-limit-range.yaml
fi

if [ -f /conf/cluster-init/manifests/storage-class.yaml ]; then
	kubectl delete -f /conf/cluster-init/manifests/storage-class.yaml
fi