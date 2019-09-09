#!/bin/bash
cp -r /conf/terraform/* .
source terraform.conf

if [ "${KUBERNETES_CLUSTER}" = "EKS" ]; then
	# Specific treatment for delete with AWS
	# We have to manually clean resources created by Kubernetes
	echo "Running scripts/aws-clean-resources.py (this will take several minutes)"
	python3 scripts/aws-clean-resources.py --env_name ${ENV_NAME}
fi

# Answer yes to destroy
terragrunt destroy-all --terragrunt-non-interactive