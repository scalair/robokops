# Monitoring
Setup monitoring using [prometheus-operator](https://github.com/coreos/kube-prometheus)
This will setup:
* Prometheus
* Grafana
* AlertManager
* Node Exporter
* Kube State Metrics

## Jsonnet
The kube-prometheus project use jsonnet as a templating solution. So in order to configure the monitoring, you have to create a `monitoring.jsonnet` file

## Kube Prometheus library version
The following version of kube-prometheus library is currently installed in `src`.
```
jb install github.com/coreos/kube-prometheus/jsonnet/kube-prometheus@release-0.1
```
When a new release is available, we will have to run the update in `src` and create a new image with it.

## Customize dashboards and alerts
You can add your own dashboards and alerts easily, instruction can be found [here](https://github.com/coreos/kube-prometheus/blob/master/docs/developing-prometheus-rules-and-grafana-dashboards.md)

## Complete example
This is a complete example of the `monitoring.jsonnet` file:
```
local k = import 'ksonnet/ksonnet.beta.3/k.libsonnet';
local secret = k.core.v1.secret;
local ingress = k.extensions.v1beta1.ingress;
local ingressTls = ingress.mixin.spec.tlsType;
local ingressRule = ingress.mixin.spec.rulesType;
local httpIngressPath = ingressRule.mixin.http.pathsType;
local pvc = k.core.v1.persistentVolumeClaim;

// Alert rules filtered out
//   KubeControllerManagerDown & KubeSchedulerDown are to be ignored
//   because these pods are not accessible in EKS 
local filter = {
  prometheusAlerts+:: {
    groups: std.map(
      function(group)
        if group.name == 'kubernetes-absent' then
          group {
            rules: std.filter(function(rule)
              rule.alert != "KubeControllerManagerDown" && rule.alert != "KubeSchedulerDown",
              group.rules
            ),
          }
        else
          group,
      super.groups
    ),
  },
};

local kp =
  (import 'kube-prometheus/kube-prometheus.libsonnet') + filter +
  {
    _config+:: {
      namespace: 'monitoring',

      // Namespace to monitore
      prometheus+:: {
        namespaces+: ['gitlabci', 'elastic-stack'],
      },

      // Customize alertmanager
      alertmanager+: {
        config: importstr 'alerts/alertmanager-config.yaml',
      },
    },

    // Custom dashboard
    grafanaDashboards+:: {
      'my-app.json': (import 'dashboards/my-app.json'),
    },

    prometheus+:: {
      prometheus+: {
        spec+: { // https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#prometheusspec
          retention: '30d',

          storage: {
            volumeClaimTemplate:
              pvc.new() +
              pvc.mixin.spec.withAccessModes('ReadWriteOnce') +
              pvc.mixin.spec.resources.withRequests({ storage: '30Gi' }) +
              pvc.mixin.spec.withStorageClassName('ssd'),
          },
        },
      },
    },

    // Create ingress objects per application
    ingress+:: {
      'alertmanager-main':
        ingress.new() +
        ingress.mixin.metadata.withName('alertmanager-main') +
        ingress.mixin.metadata.withNamespace($._config.namespace) +
        ingress.mixin.metadata.withAnnotations({
          'nginx.ingress.kubernetes.io/force-ssl-redirect': 'true',
          'nginx.ingress.kubernetes.io/auth-type': 'basic',
          'nginx.ingress.kubernetes.io/auth-secret': 'basic-auth',
          'nginx.ingress.kubernetes.io/auth-realm': 'Authentication Required',
        }) +
        ingress.mixin.spec.withRules(
          ingressRule.new() +
          ingressRule.withHost('alertmanager.example.com') +
          ingressRule.mixin.http.withPaths(
            httpIngressPath.new() +
            httpIngressPath.mixin.backend.withServiceName('alertmanager-main') +
            httpIngressPath.mixin.backend.withServicePort('web')
          ),
        ),
      grafana:
        ingress.new() +
        ingress.mixin.metadata.withName('grafana') +
        ingress.mixin.metadata.withNamespace($._config.namespace) +
        ingress.mixin.metadata.withAnnotations({
          'nginx.ingress.kubernetes.io/force-ssl-redirect': 'true',
        }) +
        ingress.mixin.spec.withRules(
          ingressRule.new() +
          ingressRule.withHost('grafana.example.com') +
          ingressRule.mixin.http.withPaths(
            httpIngressPath.new() +
            httpIngressPath.mixin.backend.withServiceName('grafana') +
            httpIngressPath.mixin.backend.withServicePort('http')
          ),
        ),
      'prometheus-k8s':
        ingress.new() +
        ingress.mixin.metadata.withName('prometheus-k8s') +
        ingress.mixin.metadata.withNamespace($._config.namespace) +
        ingress.mixin.metadata.withAnnotations({
          'nginx.ingress.kubernetes.io/force-ssl-redirect': 'true',
          'nginx.ingress.kubernetes.io/auth-type': 'basic',
          'nginx.ingress.kubernetes.io/auth-secret': 'basic-auth',
          'nginx.ingress.kubernetes.io/auth-realm': 'Authentication Required',
        }) +
        ingress.mixin.spec.withRules(
          ingressRule.new() +
          ingressRule.withHost('prometheus.example.com') +
          ingressRule.mixin.http.withPaths(
            httpIngressPath.new() +
            httpIngressPath.mixin.backend.withServiceName('prometheus-k8s') +
            httpIngressPath.mixin.backend.withServicePort('web')
          ),
        ),
    },
  } + {
    // Create basic auth secret
    ingress+:: {
      'basic-auth-secret':
        secret.new('basic-auth', { auth: std.base64(importstr 'auth') }) +
        secret.mixin.metadata.withNamespace($._config.namespace),
    },
  };

{ ['00namespace-' + name]: kp.kubePrometheus[name] for name in std.objectFields(kp.kubePrometheus) } +
{ ['0prometheus-operator-' + name]: kp.prometheusOperator[name] for name in std.objectFields(kp.prometheusOperator) } +
{ ['node-exporter-' + name]: kp.nodeExporter[name] for name in std.objectFields(kp.nodeExporter) } +
{ ['kube-state-metrics-' + name]: kp.kubeStateMetrics[name] for name in std.objectFields(kp.kubeStateMetrics) } +
{ ['alertmanager-' + name]: kp.alertmanager[name] for name in std.objectFields(kp.alertmanager) } +
{ ['prometheus-' + name]: kp.prometheus[name] for name in std.objectFields(kp.prometheus) } +
{ ['prometheus-adapter-' + name]: kp.prometheusAdapter[name] for name in std.objectFields(kp.prometheusAdapter) } +
{ ['grafana-' + name]: kp.grafana[name] for name in std.objectFields(kp.grafana) } +
{ ['ingress-' + name]: kp.ingress[name] for name in std.objectFields(kp.ingress) }
```