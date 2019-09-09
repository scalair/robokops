#!/bin/bash
./build.sh

kubectl apply -f manifests/

# At creation, it can take a few seconds for the above command to fully create the following resources, so verify the resources are ready before proceeding.
until kubectl get customresourcedefinitions servicemonitors.monitoring.coreos.com ; do date; sleep 1; echo ""; done
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done

kubectl apply -f manifests/ # This command sometimes may need to be done twice (to workaround a race condition).