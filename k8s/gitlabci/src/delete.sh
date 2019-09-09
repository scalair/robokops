#!/bin/bash
helm delete gitlabci --purge

kubectl delete -f manifests/ns.yaml