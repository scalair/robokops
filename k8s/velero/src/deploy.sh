#!/bin/bash
kubectl apply -f manifests/ns.yaml

helm upgrade --install --force --namespace velero -f velero/values.yaml velero stable/velero --version 2.1.6 --wait