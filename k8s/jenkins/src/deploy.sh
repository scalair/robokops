#!/bin/bash
kubectl apply -f manifests/ns.yaml
helm upgrade --install --force --namespace jenkins -f jenkins/values.yaml jenkins stable/jenkins --version 1.7.3 --wait