# Robokops
## TL;DR
* Downloald the latest version of Robokops: https://github.com/scalair/robokops/releases/latest
* Extract the .tar.gz and install it with `make`: `make install`
```
robokops --config example/ --terraform apply --action deploy --target cluster-init --target cluster-autoscaler --target dashboard --target monitoring --target elastic-stack
```
Checkout the [Getting started guide](docs/GETTING_STARTED.md) for details on how to set up `example`

## Introduction
Robokops is an opensource product developped by [Scalair](https://www.scalair.fr/) that helps you to manage Kubernetes clusters and deploy common features.

With Robokops, you can easily automate:
* Creation & destruction of Kubernetes cluster
* Deployment of [Prometheus Operator](https://github.com/coreos/kube-prometheus)
* Deployment of an Elastic stack
* Deployment of an ingress
* Management of DNS
* etc...

Your clusters can either be hosted and managed by [us](https://www.scalair.fr/) or you can use that product to build and manage your owns.

### Terraform
Robokops is managing infrastructure using [Terraform](https://www.terraform.io/) wrapped by [Terragrunt](https://github.com/gruntwork-io/terragrunt). Robokops does not have its own way of creating clusters, instead its relaying on existing Terraform modules to do so. The only thing Robokops is doing here is calling the code.
You can still use this product's features if you have existing clusters or if you manage them with something else than Terraform.
More info in [terraform](/terraform).

### K8s
Robokops has a set a prepacked features ready to deploy to solve common problems. Robokops abstract the complexity of deploying features by providing a common interface to manage all of them.
More info in [k8s](/k8s).

### Docker
All features deployed by Robokops are deployed using Docker containers. Containers have been chosen because they provide identical environments to run commands, such as `terraform`, `kubectl` or `helm`. All containers are based on the same image defined in [docker](/docker).

## Getting started
[Getting started guide](docs/GETTING_STARTED.md).

## Robokops parameters
Here is the list of parameters for `robokops`:

| Name        | Description                                                                                                                  | Required |
|-------------|------------------------------------------------------------------------------------------------------------------------------|----------|
| config      | Customer configuration folder. Checkout [example](/example) for more details.                                                | yes      |
| terraform   | Plan, apply or destroy infrastructure. Choose between `plan`, `apply` or `destroy`                                           | no       |
| action      | Action to execute. Choose between: `deploy`, `delete` or `dry-run`                                                           | no       |
| target      | Targets of the action. If not provided will execute against all matching configuration folders                               | no       |
| env         | Define environment variables to pass to containers. You can use `--env all` to map all env vars available in your OS context | no       |
| ssh         | Path of the .ssh directory (use only by Terraform to clone private modules)                                                  | no       |
| dev         | Add this flag to use local docker image instead of the remote registry                                                       | no       |
| version     | Return the installed version of Robokops                                                                                     | no       |

## Changelog
See [CHANGELOG.md](CHANGELOG.md)

## Documentation
The documentation can be find in [docs](/docs)