#!/bin/bash
kubectl delete -f manifests/ingress-service.yaml
kubectl delete -f manifests/ingress-deployment.yaml
kubectl delete -f manifests/ingress-rbac.yaml
kubectl delete -f manifests/ingress-configmap.yaml
