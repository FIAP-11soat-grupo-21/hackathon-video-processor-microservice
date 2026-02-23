provider "aws" {
  region = "us-east-2"
}

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  backend "s3" {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/video-processor/terraform.tfstate"
    region = "us-east-2"
  }
}