include {
  path = "${find_in_parent_folders()}"
}

terraform {
  source = "github.com/scalair/terraform-aws-eks?ref=v1.0.4"
}

dependencies {
  paths = ["../vpc"]
}

inputs = {
  vpc_bucket = "<ENV_NAME>"
  vpc_state_key = "eks/terraform/vpc/terraform.tfstate"
  vpc_state_region = "<AWS_DEFAULT_REGION>"

  subnet_bucket = "<ENV_NAME>"
  subnet_state_key = "eks/terraform/vpc/terraform.tfstate"
  subnet_state_region = "<AWS_DEFAULT_REGION>"

  eks_worker_groups  = [
    {
      name                  = "default"
      instance_type         = "t3.medium"
      asg_min_size          = "1"
      asg_max_size          = "10"
      root_volume_size      = "20"
      root_volume_type      = "gp2"
      key_name              = "<ENV_NAME>"
      ami_id                = "ami-00ac2e6b3cb38a9b9"
      autoscaling_enabled   = true
      protect_from_scale_in = true
    }
  ]
  eks_cluster_version = "1.13"
  eks_cluster_name= "<CLUSTER_NAME>"

  # For this example, we don't create an admin user for the cluster
  iam_user_create_user = false

  tags = {
    Env = "<ENV_NAME>",
    "k8s.io/cluster/<CLUSTER_NAME>" = "owned"
  }
}