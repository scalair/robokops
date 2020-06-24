# Vault
Deploy a cluster of Hashicorp Vault with Consul as a backend.
This solution use [bank-vaults](https://github.com/banzaicloud/bank-vaults) to deploy a Vault operator and use [helm](https://github.com/helm/charts/tree/master/stable/consul) to deploy Consul.

## Consul snapshot
In order to backup Consul automatically, you can use [consul-snapshot](https://github.com/pshima/consul-snapshot) with a CronJob like this:
```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: consul-snapshot
  namespace: vault
spec:
  schedule: '0 * * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: consul-snapshot
            image: anchorfree/consul-snapshot:v1.1
            args:
            - backup
            env:
              - name: CONSUL_HTTP_ADDR
                value: consul:8500
              - name: S3BUCKET
                value: my.awesome.bucket
              - name: S3REGION
                value: eu-west-1
              - name: CONSUL_SNAPSHOT_S3_SSE
                value: AES256
          restartPolicy: OnFailure
```

## Restore Consul snapshot from S3
To restore a snapshot previously backed up to S3, you can run a `consul-snapshot` image:
```
kubectl -n vault run -it --rm --restart=Never consul-snapshot --image=anchorfree/consul-snapshot:v1.1 \
--env "CONSUL_HTTP_ADDR=consul:8500" \
--env "S3BUCKET=my.awesome.bucket" \
--env "S3REGION=eu-west-1" \
--env "CONSUL_SNAPSHOT_S3_SSE=AES256" \
--command -- sh
```
And then run the `restore` command that will fetch the snapshot from S3 and restore the backend to it:
```
consul-snapshot restore backups/2019/11/15/ip-10-42-0-57.eu-west-1.compute.internal.consul.snapshot.1573830934.tar.gz
```

## Prometheus exporter
Vault automatically exposed Prometheus metrics. Using the Prometheus operator, you can create a ServiceMonitor like this to scrape the metrics:
```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app: vault
  name: vault
  namespace: vault
spec:
  endpoints:
  - interval: 30s
    port: statsd
  jobLabel: vault
  namespaceSelector:
    matchNames:
    - vault
  selector:
    matchLabels:
      app: vault
```

## AWS
### KMS
With AWS KMS, you can use the following commands to encrypt or decrypt keys stored in S3:
```
# Encrypt (change the key-id with your KMS ID)
aws kms encrypt --key-id "abcdefgh-1234-5678-90ij-klmnopqrstuv" --encryption-context "Tool=bank-vaults"  --plaintext fileb://vault-unseal-0.txt --output text --query CiphertextBlob | base64 -D > vault-unseal-0
# Decrypt
aws kms decrypt --encryption-context "Tool=bank-vaults" --ciphertext-blob fileb://vault-unseal-0 --output text --query Plaintext | base64 --decode > vault-unseal-0.txt
```
