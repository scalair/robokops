#!/bin/bash

# Dashboard
helm delete kubernetes-dashboard --purge

# Metrics server
if [ "${METRICS_SERVER}" = "true" ]; then
	kubectl delete -f heapster/heapster.yaml
	kubectl delete -f heapster/influxdb.yaml
	kubectl delete -f heapster/heapster-rbac.yaml
fi
