#!/bin/bash

# Fluent-bit
helm delete fluent-bit --purge

# Fluentd
helm delete fluentd --purge

if [ "${INSTALL_KIBANA}" != "false" ]; then
	# Kibana
	helm delete kibana --purge
fi
if [ "${INSTALL_ELASTICSEARCH_CURATOR}" != "false" ]; then
	# ElasticSearch-curator
	helm delete elasticsearch-curator --purge
fi
if [ "${INSTALL_ELASTICSEARCH}" != "false" ]; then
	# ElasticSearch
	helm delete elasticsearch --purge
fi

kubectl -n elastic-stack delete --all pvc
kubectl delete -f manifests/ns.yaml
