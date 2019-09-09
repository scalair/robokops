#!/bin/bash
cp -r /conf/terraform/* .
source terraform.conf

# Plan-all will not work if you have dependencies to modules that haven't been apply yet
terragrunt plan-all
