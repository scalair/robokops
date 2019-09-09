#!/bin/bash
kubectl apply -f manifests/rbac-admin.yaml
helm init --service-account admin --upgrade --wait

if [ -f /conf/cluster-init/manifests/kube-system-limit-range.yaml ]; then
	kubectl apply -f /conf/cluster-init/manifests/kube-system-limit-range.yaml
fi

if [ -f /conf/cluster-init/manifests/storage-class.yaml ]; then
	kubectl apply -f /conf/cluster-init/manifests/storage-class.yaml
fi