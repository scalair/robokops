#!/bin/bash
cp -r /conf/terraform/* .
source terraform.conf

# Answer yes to apply
terragrunt apply-all --terragrunt-non-interactive
