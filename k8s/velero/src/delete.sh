#!/bin/bash
helm delete velero --purge
kubectl delete crds -l app.kubernetes.io/instance=velero
kubectl delete -f manifests/ns.yaml