# EKS
## Configure local access to a cluster
If you need to access the cluster locally, you have to configure your environment like that:
* Install [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/install-linux-al2017.html)
* Install [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html)
* Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) (if needed)

Then, configure your credentials in `~/.aws/credentials`, for instance:
```
[admin.eks.dev.example.eu-west-1]
aws_default_region=eu-west-1
aws_access_key_id=<YOUR_KEY>
aws_secret_access_key=<YOUR_SECRET>
```
Configure your environment variable to use the right profile, for instance:
```
export AWS_PROFILE=admin.eks.dev.example.eu-west-1
```
Update your kube config, for instance (replace the last part with the name of your cluster):
```
aws eks --region eu-west-1 update-kubeconfig --name eks-dev-example-eu-west-1
```
Then you're ready to go, you can test your configuration by running:
```
kubectl get namespace
```
Which should return something like that:
```
NAME          STATUS   AGE
default       Active   36m
kube-public   Active   36m
kube-system   Active   36m
```
