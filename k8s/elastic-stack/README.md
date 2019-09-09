# Elastic Stack
Setup Elastic stack with:
* ElasticSearch
* ElasticSearch Curator
* Kibana
* Fluentd
* Fluent-bit

Fluent-bit are installed as daemon set on every worker nodes, they collect logs and send them to Fluentd. Fluentd act as the aggregator and send the logs to ElasticSearch which can then be visualise using Kibana.

## Persistent volumes
ElasticSearch and Fluentd should be configured to use persistent volumes.
Data from ElasticSearch should be backed up in case of a disaster. ElasticSearch Curator can be used to do that.

## Logs retention
Elasticsearch Curator can be configured to automatically purged ElasticSearch indices older than a certain numbers of days.
To configure elasticsearch curator for retention do:
```
# elastic-stack/elasticsearch-curator/values.yaml
cronjob:
  schedule: '* */6 * * *'

configMaps:
  action_file_yml: |-
    ---
    actions:
      1:
        action: delete_indices
        description: "Clean up ES by deleting old indices"
        options:
          timeout_override:
          continue_if_exception: False
          disable_action: False
          ignore_empty_list: True
        filters:
        - filtertype: age
          source: name
          direction: older
          timestring: '%Y.%m.%d'
          unit: days
          unit_count: 7
          field:
          stats_result:
          epoch:
          exclude: False
```

## External Elasticsearch
You can use the elastic-stack with an Elasticsearch outside of your cluster. In order to do that, disable the installation of this Elasticsearch and configure you Elasticsearch endpoint
```
# elastic-stack/elastic-stack.conf
export INSTALL_ELASTICSEARCH=false
export ELASTICSEARCH_ENDPOINT=my-awesome-elasticsearch-on-aws-abcdef12345.eu-west-1.es.amazonaws.com
```
Then configure Fluentd and Kibana to use that endpoint:
```
# elastic-stack/fluentd/values.yaml
output:
  host: ${ELASTICSEARCH_ENDPOINT}
  port: 443
  scheme: https
  sslVersion: TLSv1
```
```
# elastic-stack/kibana/values.yaml
files:
  kibana.yml:
    elasticsearch.hosts: http://${ELASTICSEARCH_ENDPOINT}:80
```
### AWS ES
When using AWS ES, you have to configure Fluentd like that:
```
# elastic-stack/fluentd/values.yaml
configMaps:
  output.conf: |
    <match **>
      @id elasticsearch
      @type elasticsearch
      host "${ELASTICSEARCH_ENDPOINT}"
      port "443"
      scheme "https"
      ssl_version "TLSv1"
      logstash_format true
      logstash_prefix fluentd
      # Prevent reloading connections to AWS ES
      reload_on_failure false
      reload_connections false
      ...
    </match>
```
The important parameters are `reload_on_failure` and `reload_connections` which must be set to `false`.