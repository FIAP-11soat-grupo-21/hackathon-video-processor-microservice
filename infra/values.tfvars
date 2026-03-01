application_name = "video-processor-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8080
app_path_pattern = ["/v1/videos*", "/v1/videos/*"]

desired_count = 2

container_environment_variables = {
  GO_ENV   = "production"
  API_PORT = "8080"
  API_HOST = "0.0.0.0"
  LOG_LEVEL = "INFO"
}

container_secrets       = {}
health_check_path       = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonS3FullAccess",
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
]
apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"

lambda_environment_variables = {
  LOG_LEVEL = "INFO"
}

lambda_memory_size = 2048
lambda_timeout     = 300

sqs_queue_name = "video-frame-queue"

s3_bucket = "fiap-tc-terraform-846874"
