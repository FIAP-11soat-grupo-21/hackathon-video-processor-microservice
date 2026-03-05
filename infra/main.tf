module "frame_processor_lambda" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Lambda?ref=main"

  lambda_name = "video-frame-processor"
  handler     = "bootstrap"
  runtime     = "provided.al2023"
  subnet_ids  = data.terraform_remote_state.network_vpc.outputs.private_subnets
  
  environment = merge(
    var.lambda_environment_variables,
    {
      AWS_SQS_FRAME_EXTRACTION_QUEUE  = aws_sqs_queue.frame_queue.url
      S3_BUCKET                       = var.s3_bucket
    }
  )
  
  vpc_id      = data.terraform_remote_state.network_vpc.outputs.vpc_id
  memory_size = var.lambda_memory_size
  timeout     = var.lambda_timeout

  s3_bucket = var.lambda_bucket_name
  s3_key    = var.lambda_s3_key

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
        "sqs:GetQueueUrl",
        "sqs:SendMessage"
      ]
      resources = [
        aws_sqs_queue.frame_queue.arn,
        aws_sqs_queue.frame_dlq.arn
      ]
    }
  }

  tags = data.terraform_remote_state.app_registry.outputs.app_registry_application_tag
}