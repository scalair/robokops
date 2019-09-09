# Example
This folder contain an example of Robokops configuration. This is used by the getting started guide.

## More examples
Create a Kubernetes cluster with some features from scratch:
```
robokops --config example/ --terraform apply --action deploy --target cluster-init --target cluster-autoscaler --target dashboard --target monitoring --target elastic-stack
```

Update some features:
```
robokops --config example/ --action deploy --target dashboard --target monitoring
```

Delete one feature. Cannot be undone:
```
robokops --config example/ --action delete --target monitoring
```

Delete everything. Cannot be undone:
```
robokops --config example/ --terraform destroy
```

Generate all manifests without making modifications (dry run)
```
robokops --config example/ --action dry-run
```
