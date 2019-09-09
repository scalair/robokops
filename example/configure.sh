#!/bin/bash

# This script configure the "conf" and "terraform" folder with
# the variables you set before
#
# If you want to rollback the example folder to change something,
# use Git to reset it.

###################################
### Check variables are defined ###
###################################
if [[ -z "${AWS_ACCESS_KEY_ID}" ]]; then
	echo "AWS_ACCESS_KEY_ID must be defined"
	exit 1
fi
if [[ -z "${AWS_SECRET_ACCESS_KEY}" ]]; then
	echo "AWS_SECRET_ACCESS_KEY must be defined"
	exit 1
fi
if [[ -z "${AWS_DEFAULT_REGION}" ]]; then
	echo "AWS_DEFAULT_REGION must be defined"
	exit 1
fi
if [[ -z "${ENV_NAME}" ]]; then
	echo "ENV_NAME must be defined"
	exit 1
fi
if [[ -z "${CLUSTER_NAME}" ]]; then
	echo "CLUSTER_NAME must be defined"
	exit 1
fi

##########################
### Search and replace ###
##########################
find terraform -name ".terragrunt-cache" -exec rm -f {} \;
FILES=$(find . -type f)
for file in ${FILES}; do
	# We use `-i.bk` for this command to work on Linux and MacOS
	sed -i.bk -e "s|AKIAJLWY7TR42BGIM76A|${AWS_ACCESS_KEY_ID}|g" ${file}
	sed -i.bk -e "s|nVS1y3ziOX3q1LO64tARWqg1/DLaiA1R/P3XSgx0|${AWS_SECRET_ACCESS_KEY}|g" ${file}
	sed -i.bk -e "s|eu-west-1|${AWS_DEFAULT_REGION}|g" ${file}
	sed -i.bk -e "s|eks.dev.scalair.eu-west-1|${ENV_NAME}|g" ${file}
	sed -i.bk -e "s|eks-dev-scalair-eu-west-1|${CLUSTER_NAME}|g" ${file}
done
# Clean the .bk files created
find . -name "*.bk" -exec rm -f {} \;
