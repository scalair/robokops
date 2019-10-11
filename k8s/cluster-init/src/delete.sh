#!/bin/bash
kubectl delete -f manifests
kubectl delete deployment tiller-deploy --namespace kube-system
kubectl delete service tiller-deploy --namespace kube-system