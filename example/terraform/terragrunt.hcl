terraform {
  extra_arguments "custom_vars" {
    commands = get_terraform_commands_that_need_vars()
  }
  after_hook "copy_common_main_providers" {
    commands = ["init-from-module"]
    execute  = ["cp", "${get_terragrunt_dir()}/../main_providers.tf", "."]
  }
}

remote_state {
  backend = "s3"
  config = {
    bucket         = "<ENV_NAME>"
    key            = "eks/terraform/${path_relative_to_include()}/terraform.tfstate"
    region         = "<AWS_DEFAULT_REGION>"
    encrypt        = true
    dynamodb_table = "<ENV_NAME>"
    
    s3_bucket_tags = {
      env = "<ENV_NAME>"
    }

    dynamodb_table_tags = {
      env = "<ENV_NAME>"
    }
  }
}