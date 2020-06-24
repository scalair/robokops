#!/bin/bash

# Vault
kubectl delete -f manifests/vault.yaml
helm delete vault-operator --purge
kubectl delete -f manifests/rbac.yaml

# Consul
helm delete consul --purge

kubectl -n vault delete pvc --all
kubectl delete -f manifests/ns.yaml