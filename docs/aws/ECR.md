# ECR
AWS container registry

## Push an image to ECR
In order to push image to ECR, you need to configure your environment with credentials to access ECR.
Configure your `~/.aws/credentials`:
```
[docker]
aws_default_region=eu-west-1
aws_access_key_id=<YOUR_KEY>
aws_secret_access_key=<YOUR_SECRET>
```
Configure your environment variable:
```
export AWS_PROFILE=docker
```
Login to ECR:
```
$(aws ecr get-login --no-include-email --region eu-west-1)
```
Then you can build, tag and push your image, for instance:
```
docker build -t my-app .
docker tag my-app:latest 123456789012.dkr.ecr.eu-west-1.amazonaws.com/my-app:0.0.1
docker push 123456789012.dkr.ecr.eu-west-1.amazonaws.com/my-app:0.0.1
```

## Access images from containers
The easiest way to access ECR images from containers it to attach a role to the worker nodes that give permission to read ECR. Then you don't have to manage the registry access.