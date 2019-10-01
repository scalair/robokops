#!/bin/bash
helm upgrade --install --force --namespace kube-system -f kubewatch/values.yaml kubewatch stable/kubewatch --version 0.8.9 --wait