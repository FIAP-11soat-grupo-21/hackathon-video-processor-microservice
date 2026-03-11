lambda_environment_variables = {
  GO_ENV              = "production"
  LOG_LEVEL           = "INFO"
  FRAME_QUALITY       = "85"
  FRAME_FORMAT        = "jpg"
  MAX_CONCURRENT_JOBS = "10"
  PROCESS_TIMEOUT_SEC = "300"
}

bucket_name = "fiap-tc-terraform-functions-846874"

sqs_queue_name = "video-frame-queue"

ecr_repository_url = "846874.dkr.ecr.us-east-2.amazonaws.com/video-processor-api"
container_image_tag = "latest"

service_discovery_namespace = "video-processor.local"