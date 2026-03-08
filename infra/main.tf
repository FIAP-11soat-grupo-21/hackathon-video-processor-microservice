# Lambda Chunk Processor
module "chunk_processor_lambda" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Lambda?ref=main"

  lambda_name = "chunk-processor"
  handler     = "bootstrap"
  runtime     = "provided.al2023"
  subnet_ids  = data.terraform_remote_state.network_vpc.outputs.private_subnets
  
  environment = merge(
    var.lambda_environment_variables,
    {
      S3_BUCKET                  = var.s3_bucket
      SNS_CHUNK_PROCESSED_TOPIC  = data.terraform_remote_state.sns_chunk_processed.outputs.topic_arn
    }
  )
  
  vpc_id      = data.terraform_remote_state.network_vpc.outputs.vpc_id
  memory_size = var.chunk_processor_memory_size
  timeout     = var.chunk_processor_timeout

  s3_bucket = var.lambda_bucket_name
  s3_key    = var.chunk_processor_s3_key

  role_permissions = {
    s3 = {
      actions = [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ]
      resources = [
        "arn:aws:s3:::${var.s3_bucket}",
        "arn:aws:s3:::${var.s3_bucket}/*"
      ]
    }
    sqs = {
      actions = [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:GetQueueUrl"
      ]
      resources = [
        data.terraform_remote_state.sqs_chunk_processor.outputs.sqs_queue_arn
      ]
    }
    sns = {
      actions = [
        "sns:Publish"
      ]
      resources = [
        data.terraform_remote_state.sns_chunk_processed.outputs.topic_arn
      ]
    }
  }

  tags = data.terraform_remote_state.app_registry.outputs.app_registry_application_tag
}

# Lambda Update Video Chunk Status
module "update_video_chunk_status_lambda" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Lambda?ref=main"

  lambda_name = "update-video-chunk-status"
  handler     = "bootstrap"
  runtime     = "provided.al2023"
  subnet_ids  = data.terraform_remote_state.network_vpc.outputs.private_subnets
  
  environment = merge(
    var.lambda_environment_variables,
    {
      DYNAMODB_TABLE_NAME              = var.dynamodb_table_name
      S3_BUCKET                        = var.s3_bucket
      SNS_ALL_CHUNKS_PROCESSED_TOPIC   = data.terraform_remote_state.sns_all_chunks_processed.outputs.topic_arn
    }
  )
  
  vpc_id      = data.terraform_remote_state.network_vpc.outputs.vpc_id
  memory_size = var.update_status_memory_size
  timeout     = var.update_status_timeout

  s3_bucket = var.lambda_bucket_name
  s3_key    = var.update_status_s3_key

  role_permissions = {
    dynamodb = {
      actions = [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ]
      resources = [
        "arn:aws:dynamodb:${data.aws_region.current.id}:*:table/${var.dynamodb_table_name}"
      ]
    }
    s3 = {
      actions = [
        "s3:GetObject",
        "s3:ListBucket"
      ]
      resources = [
        "arn:aws:s3:::${var.s3_bucket}",
        "arn:aws:s3:::${var.s3_bucket}/*"
      ]
    }
    sqs = {
      actions = [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:GetQueueUrl"
      ]
      resources = [
        data.terraform_remote_state.sqs_update_video_chunk_status.outputs.sqs_queue_arn
      ]
    }
    sns = {
      actions = [
        "sns:Publish"
      ]
      resources = [
        data.terraform_remote_state.sns_all_chunks_processed.outputs.topic_arn
      ]
    }
  }

  tags = data.terraform_remote_state.app_registry.outputs.app_registry_application_tag
}

# Lambda Zip Processor
module "zip_processor_lambda" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Lambda?ref=main"

  lambda_name = "zip-processor"
  handler     = "bootstrap"
  runtime     = "provided.al2023"
  subnet_ids  = data.terraform_remote_state.network_vpc.outputs.private_subnets
  
  environment = merge(
    var.lambda_environment_variables,
    {
      S3_BUCKET                           = var.s3_bucket
      SNS_VIDEO_PROCESSING_COMPLETE_TOPIC = data.terraform_remote_state.sns_video_processing_complete.outputs.topic_arn
    }
  )
  
  vpc_id      = data.terraform_remote_state.network_vpc.outputs.vpc_id
  memory_size = var.zip_processor_memory_size
  timeout     = var.zip_processor_timeout

  s3_bucket = var.lambda_bucket_name
  s3_key    = var.zip_processor_s3_key

  role_permissions = {
    s3 = {
      actions = [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ]
      resources = [
        "arn:aws:s3:::${var.s3_bucket}",
        "arn:aws:s3:::${var.s3_bucket}/*"
      ]
    }
    sqs = {
      actions = [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:GetQueueUrl"
      ]
      resources = [
        data.terraform_remote_state.sqs_zip_processor.outputs.sqs_queue_arn
      ]
    }
    sns = {
      actions = [
        "sns:Publish"
      ]
      resources = [
        data.terraform_remote_state.sns_video_processing_complete.outputs.topic_arn
      ]
    }
  }

  tags = data.terraform_remote_state.app_registry.outputs.app_registry_application_tag
}