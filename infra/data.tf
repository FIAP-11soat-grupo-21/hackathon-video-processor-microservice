data "aws_region" "current" {}

data "terraform_remote_state" "network_vpc" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/Network/VPC/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "app_registry" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/AppRegistry/terraform.tfstate"
    region = "us-east-2"
  }
}