data "aws_region" "current" {}

data "terraform_remote_state" "api_gateway" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/GatewayAPI/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "network_vpc" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/Network/VPC/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "ecs" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/ECS/Cluster/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "alb" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/Network/ALB/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "ghcr_secret" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/Secrets/GHCR/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "auth" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/auth/terraform.tfstate"
    region = "us-east-2"
  }
}