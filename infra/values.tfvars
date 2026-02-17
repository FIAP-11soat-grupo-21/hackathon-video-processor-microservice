application_name = "video-processor-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8080
app_path_pattern = ["/process*"]

container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8080"
  API_HOST : "0.0.0.0"
}

container_secrets = {}
health_check_path = "/health"
task_role_policy_arns = []

apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"
