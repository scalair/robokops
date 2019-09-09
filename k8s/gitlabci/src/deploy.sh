#!/bin/bash
helm repo add gitlab https://charts.gitlab.io
helm repo update

kubectl apply -f manifests/ns.yaml

helm upgrade --install --force --namespace gitlabci -f gitlab-runner/values.yaml gitlabci gitlab/gitlab-runner --version 0.7.0 --wait