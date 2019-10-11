#!/bin/bash
kubectl apply -f manifests
helm init --service-account admin --upgrade --wait