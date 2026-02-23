module "frame_processor_lambda" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Lambda?ref=main"

  api_id      = data.terraform_remote_state.api_gateway.outputs.api_id
  lambda_name = "video-frame-processor"
  handler     = "bootstrap"
  runtime     = "provided.al2023"
  subnet_ids  = data.terraform_remote_state.network_vpc.outputs.private_subnets
  
  environment = merge(
    var.lambda_environment_variables,
    {
      FRAME_QUEUE_URL  = aws_sqs_queue.frame_queue.url
      S3_INPUT_BUCKET  = var.s3_input_bucket
      S3_OUTPUT_BUCKET = var.s3_output_bucket
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
        "arn:aws:s3:::${var.s3_input_bucket}",
        "arn:aws:s3:::${var.s3_input_bucket}/*",
        "arn:aws:s3:::${var.s3_output_bucket}",
        "arn:aws:s3:::${var.s3_output_bucket}/*"
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
}

module "video_processor_api" {
  source     = "git::ssh://git@github.com/FIAP-11soat-grupo-21/infra-core.git//modules/ECS-Service?ref=main"
  depends_on = [aws_lb_listener.listener]

  cluster_id            = data.terraform_remote_state.ecs.outputs.cluster_id
  ecs_security_group_id = data.terraform_remote_state.ecs.outputs.ecs_security_group_id

  cloudwatch_log_group     = data.terraform_remote_state.ecs.outputs.cloudwatch_log_group
  ecs_container_image      = var.image_name
  ecs_container_name       = var.application_name
  ecs_container_port       = var.image_port
  ecs_service_name         = var.application_name
  ecs_desired_count        = var.desired_count
  registry_credentials_arn = data.terraform_remote_state.ghcr_secret.outputs.secret_arn

  ecs_container_environment_variables = merge(
    var.container_environment_variables,
    {
      FRAME_QUEUE_URL  = aws_sqs_queue.frame_queue.url
      S3_INPUT_BUCKET  = var.s3_input_bucket
      S3_OUTPUT_BUCKET = var.s3_output_bucket
    }
  )

  ecs_container_secrets = var.container_secrets

  private_subnet_ids      = data.terraform_remote_state.network_vpc.outputs.private_subnets
  task_execution_role_arn = data.terraform_remote_state.ecs.outputs.task_execution_role_arn
  task_role_policy_arns   = var.task_role_policy_arns
  alb_target_group_arn    = aws_alb_target_group.target_group.arn
  alb_security_group_id   = data.terraform_remote_state.alb.outputs.alb_security_group_id
}

module "VideoProcessorAPIRoutes" {
  source     = "git::ssh://git@github.com/FIAP-11soat-grupo-21/infra-core.git//modules/API-Gateway-Routes?ref=main"
  depends_on = [module.video_processor_api]

  api_id       = data.terraform_remote_state.api_gateway.outputs.api_id
  alb_proxy_id = aws_apigatewayv2_integration.alb_proxy.id

  endpoints = {
    process_video = {
      route_key  = "POST /videos/process"
      restricted = false
    }
  }
}