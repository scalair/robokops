#!/bin/bash
./build.sh

# Extract the namespace manifest first, so we don't delete the namespace before everything else
mv manifests/00namespace-namespace.yaml .

kubectl delete -f manifests/

kubectl -n monitoring delete --all pvc

# Finally, the namespace can be deleted
kubectl delete -f 00namespace-namespace.yaml
