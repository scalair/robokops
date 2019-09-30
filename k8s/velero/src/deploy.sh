#!/bin/bash
kubectl apply -f manifests/ns.yaml

helm upgrade --install --force --namespace velero -f velero/values.yaml velero stable/velero --version 2.1.6 --wait

# velero install \
#     --provider aws \
#     --bucket eks.stage.scalair.eu-west-1 \
#     --prefix velero \
#     --secret-file ./credentials-velero \
#     --backup-location-config region=eu-west-1 \
#     --snapshot-location-config region=eu-west-1