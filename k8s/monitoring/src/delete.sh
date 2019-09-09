#!/bin/bash
./build.sh

# Extract the namespace manifest first, so we don't delete the manifests before everything else
mv manifests/00namespace-namespace.yaml .

kubectl delete -f manifests/

# Finally, the namespace can be deleted
kubectl delete -f 00namespace-namespace.yaml