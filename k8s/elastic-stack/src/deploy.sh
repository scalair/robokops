#!/bin/bash
kubectl apply -f manifests/ns.yaml

# Fluent-bit
#   Fluent-bit is installed first because it's a DS and it gets more chances to have enough spaces on nodes
#   if it's installed before ES. It will spin until Fluentd is up but it shouldn't be a big deal
helm upgrade --install --force --namespace elastic-stack -f fluent-bit/values.yaml fluent-bit stable/fluent-bit --version 2.5.0 --wait

if [ "${INSTALL_ELASTICSEARCH}" != "false" ]; then
	# ElasticSearch
	helm upgrade --install --force --namespace elastic-stack -f elasticsearch/values.yaml elasticsearch stable/elasticsearch --version 1.30.0 --wait
fi
if [ "${INSTALL_ELASTICSEARCH_CURATOR}" != "false" ]; then
	# ElasticSearch-curator
	helm upgrade --install --force --namespace elastic-stack -f elasticsearch-curator/values.yaml elasticsearch-curator stable/elasticsearch-curator --version 1.5.0 --wait
fi
if [ "${INSTALL_KIBANA}" != "false" ]; then
	# Kibana
	helm upgrade --install --force --namespace elastic-stack -f kibana/values.yaml kibana stable/kibana --version 3.2.3 --wait
fi
# Fluentd
helm upgrade --install --force --namespace elastic-stack -f fluentd/values.yaml fluentd stable/fluentd --version 2.4.2 --wait