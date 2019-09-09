#!/bin/bash
helm upgrade --install --force --namespace kube-system -f cluster-autoscaler/values.yaml cluster-autoscaler stable/cluster-autoscaler --version 3.0.0 --wait