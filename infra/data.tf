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

data "terraform_remote_state" "sns_chunk_processed" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/SNS/chunk-processed/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "sns_all_chunks_processed" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/SNS/all-chunks-processed/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "sns_video_processing_complete" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/SNS/video-processing-complete/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "sqs_chunk_processor" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/SQS/chunk-processor/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "sqs_update_video_chunk_status" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/SQS/update-video-chunk-status/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "sqs_zip_processor" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/SQS/zip-processor/terraform.tfstate"
    region = "us-east-2"
  }
}