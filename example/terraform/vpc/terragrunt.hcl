include {
  path = "${find_in_parent_folders()}"
}

terraform {
  source = "github.com/terraform-aws-modules/terraform-aws-vpc?ref=v2.9.0"
}

inputs = {
  name = "<ENV_NAME>"
  cidr = "10.0.0.0/16"

  azs             = ["<AWS_DEFAULT_REGION>a", "<AWS_DEFAULT_REGION>b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_nat_gateway = true

  tags = {
    Env = "<ENV_NAME>"
    KubernetesCluster = "<CLUSTER_NAME>",
    "kubernetes.io/cluster/<CLUSTER_NAME>" = "shared",
    "kubernetes.io/role/elb" = ""
  }
}