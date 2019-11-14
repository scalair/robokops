
# K8s
Each folder here is a feature that can be deployed in an existing cluster.
There are all packaged using `Docker` container that supports 3 actions:
* deploy -> call `deploy.sh`
* delete -> call `delete.sh`
* dry-run -> generate and return merged and templated manifests

The entrypoint of each container is [entrypoint.sh](../docker/entrypoint.sh) of `robokops-base`.

The configuration provided to Robokops is mounted into the container (in `/conf`) and environment variables are passed.
Robokops should be used to run those containers, but you can run them by yourself if necessary by running something like:
```
docker run -v "$(pwd)/../example/conf/:/conf" robokops-feature [deploy|delete|dry-run]
```

## Cluster login
The access to your cluster must be configured in the `conf/cluster-login.sh` file. This script must configure the `~/.kube/config` file, it could either get it from a secret management tool or you can run authentification commands in order to create it.

## Common configuration
You can set in `conf/common.conf` file variables that can be used by any features or by the `cluster-login.sh` script. This file can be used to fetch the cloud provider credentials for instance.

## Override configuration
Yaml files can easily be overwritten, you just need to create a file with the same name, under the same file hierarchie in your configuration folder. If present, the Yaml file you provided will be merged with the one existing (using [yq](https://github.com/mikefarah/yq)), so you don't have to duplicate the entiere file, you can just override what's specific for you. It doesn't matter if yaml files are k8s manifests or Helm configuration, it will work the same.

### Example
If you want to override `domainFilters` in [ingress-nginx/src/external-dns/values.yaml](ingress-nginx/src/external-dns/values.yaml)

just create the file:
```
example/conf/ingress-nginx/external-dns/values.yaml
```
with the content:
```
domainFilters: ["your-company.com"]
```
and run the deployment:
```
robokops --config example/ --action deploy --target ingress-nginx
```

## Secrets
Managing secrets is another tricky part with Kubernetes, since you don't want to store them in plain text directly in manifests, you need a templating mechanism.
With Robokops, you can use variables in your manifests which will be templated using envsubst.

### Example
To change Grafana credentials, create:
```
example/conf/monitoring/post-manifests/grafana-credentials.yaml
```
with the following content:
```
apiVersion: v1
kind: Secret
metadata:
  name: grafana-credentials
  namespace: monitoring
data:
  user: ${GRAFANA_USERNAME}
  password: ${GRAFANA_PASSWORD}
```
and set `GRAFANA_USERNAME` and `GRAFANA_PASSWORD` in:
```
example/conf/monitoring/monitoring.conf
```
using [Vault](https://www.vaultproject.io/):
```
export GRAFANA_USERNAME=$(vault read -field username secret/.../grafana)
export GRAFANA_PASSWORD=$(vault read -field password secret/.../grafana)
```
You will also have to create:
```
example/conf/monitoring/manifests/grafana-deployment.yaml
```
to override the Grafana deployment object to read credentials from that secret (you can find an example below).

## Kubectl
To manage k8s objects with `kubectl`, we recommand to use [declarative management](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/), which means to use:
```
kubectl apply -f manifests/
```
to create or update manifests and to use:
```
kubectl delete -f manifests/
```
to delete them.

## Helm
When using `helm`, in order to make creation and update idempotent, the recommanded way is to use:
```
helm upgrade --install --force [...]
```
which will either install or update the chart and will replace a failed deployment. To delete chart, we recommand to use:
```
helm delete [...] --purge
```
Also, you should always specify a chart version and use `--wait` option to avoid concurrent issues.

## Limitation

### Yq
#### Special characters
Some special chars in manifests are not always well managed by yq, for instance if you set something like:
```
cronjob:
  schedule: "* */6 * * *"
```
yq will poorly crash. The workarround is to change double quote with single:
```
cronjob:
  schedule: '* */6 * * *'
```
We recommand to use `dry-run` if you want to make sure your manifests get properly merged first. 

### Yaml list of objects override
Let's say you want to add environment variables to a deployment object, you cloud do something like that:
```
apiVersion: apps/v1beta2
kind: Deployment
spec:
  template:
    spec:
      containers:
      - image: grafana/grafana:6.0.1
        env:
        - name: GF_SECURITY_ADMIN_USER
          valueFrom:
            secretKeyRef:
              name: grafana-credentials
              key: user
        - name: GF_SECURITY_ADMIN_PASSWORD
          valueFrom:
            secretKeyRef:
              name: grafana-credentials
              key: password
```
and hope that the rest of the content (such as `ports`, `volumeMounts`, etc) will remains... well unfortunalely it's not the case. If you want to override a list of objects you gonna have to rewritte the all thing. In this case something like that:
```
apiVersion: apps/v1beta2
kind: Deployment
spec:
  template:
    spec:
      containers:
      - image: grafana/grafana:6.0.1
        name: grafana
        env:
        - name: GF_SECURITY_ADMIN_USER
          valueFrom:
            secretKeyRef:
              name: grafana-credentials
              key: user
        - name: GF_SECURITY_ADMIN_PASSWORD
          valueFrom:
            secretKeyRef:
              name: grafana-credentials
              key: password
        ports:
        ...
        readinessProbe:
        ...
        resources:
        ...
        volumeMounts:
        ...
```