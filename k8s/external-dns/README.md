# External DNS
[ExternalDNS](https://github.com/helm/charts/tree/master/stable/external-dns) is a Kubernetes addon that configures public DNS servers with information about exposed Kubernetes services to make them discoverable.
When using external-dns, you must specify the domain want to update:
```
# external-dns/external-dns/values.yaml
domainFilters: ["example.com"]
```
## AWS
When using external-dns with AWS route53, you must specify the role ARN and the region:
```
# external-dns/external-dns/values.yaml
aws:
  roleArn: route53.my.eks.cluster.eu-west-1
  region: "eu-west-1"
```