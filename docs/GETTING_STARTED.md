# Getting started
This example creates a Kubernetes cluster on AWS using EKS with monitoring, log collector, administration dashboard and autoscaling.

## Prerequisites
* Downloald the latest version of Robokops: https://github.com/scalair/robokops/releases/latest
* Extract the .tar.gz and install it with `make`: `make install`
* You must have an AWS account to run that example
* Create a key-pair on AWS with the same value you will set on `ENV_NAME` below.
* To access the EKS cluster, you need [aws cli](https://aws.amazon.com/cli/) and [aws_iam_authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html)
* Finally, you need [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Configuration
Set the variables below:
```
# Your AWS credentials
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=
# Name of your environment, for instance eks.dev.scalair.eu-west-1
# This will be used to create the S3 bucket to store Terraform state.
# Therefore this name must be unique to you
export ENV_NAME=
# Name of your k8s cluster, for instance eks-dev-scalair-eu-west-1
# Cannot contains '.'
export CLUSTER_NAME=
```
And run [configure.sh](configure.sh) (this will configure the example folder for you):
```
cd example && ./configure.sh && cd -
```

## Run Robokops
Create the cluster and deploy the features (this will take around 20min):
```
robokops --config example/ --terraform apply --action deploy --target cluster-init --target cluster-autoscaler --target dashboard --target monitoring --target elastic-stack
```
If something goes wrong during creation or deployment, you can rerun the same command (which is idempotent).
For instance, creating the elastic-stack will take time and sometimes it can timeout (still the stack could be properly deployed). If it happens and not everything has been properly deployed, try to rerun that command. 

## Access it
Authenticate to the cluster:
```
source example/conf/common.conf
./example/conf/cluster-login.sh
```

To access Grafana run (username/password: `admin/admin`):
```
POD_NAME=$(kubectl get pods --namespace monitoring -l "app=grafana" -o jsonpath="{.items[0].metadata.name}")
echo "Visit http://127.0.0.1:3000 to use Grafana"
kubectl port-forward --namespace monitoring $POD_NAME 3000:3000
```

To access Kibana run:
```
POD_NAME=$(kubectl get pods --namespace elastic-stack -l "app=kibana,release=kibana" -o jsonpath="{.items[0].metadata.name}")
echo "Visit http://127.0.0.1:5601 to use Kibana"
kubectl port-forward --namespace elastic-stack $POD_NAME 5601:5601
```

To access the admin dashboard run:
```
POD_NAME=$(kubectl get pods -n kube-system -l "app=kubernetes-dashboard,release=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
echo http://127.0.0.1:9090/
TOKEN=$(kubectl -n kube-system get secret $(kubectl -n kube-system get secret | grep admin-token | awk '{print $1}') -o jsonpath="{.data.token}" | base64 --decode)
echo "Token: ${TOKEN}"
kubectl -n kube-system port-forward $POD_NAME 9090:9090
```

## Clean
To remove the cluster and all resources in AWS, just run and wait (this will take around another 20min):
```
robokops --config example/ --terraform destroy
```
