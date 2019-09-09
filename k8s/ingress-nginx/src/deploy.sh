#!/bin/bash
kubectl apply -f manifests/ingress-configmap.yaml
kubectl apply -f manifests/ingress-rbac.yaml
kubectl apply -f manifests/ingress-deployment.yaml
kubectl apply -f manifests/ingress-service.yaml