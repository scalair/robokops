#!/bin/bash
helm delete jenkins --purge
kubectl -n jenkins delete --all pvc
kubectl delete -f manifests/ns.yaml