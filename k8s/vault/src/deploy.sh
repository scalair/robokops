#!/bin/bash
kubectl apply -f manifests/ns.yaml

# Consul
helm upgrade --install --force --namespace vault -f consul/values.yaml consul stable/consul --version 3.9.2 --wait

# Vault
helm repo add banzaicloud-stable http://kubernetes-charts.banzaicloud.com/branch/master
helm repo update
helm upgrade --install --force --namespace vault -f vault-operator/values.yaml vault-operator banzaicloud-stable/vault-operator --wait
kubectl apply -f manifests/rbac.yaml
kubectl apply -f manifests/vault.yaml