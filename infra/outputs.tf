output "chunk_processor_lambda_arn" {
  value       = module.chunk_processor_lambda.lambda_arn
  description = "ARN da função Lambda chunk-processor"
}

output "update_status_lambda_arn" {
  value       = module.update_video_chunk_status_lambda.lambda_arn
  description = "ARN da função Lambda update-video-chunk-status"
}

output "zip_processor_lambda_arn" {
  value       = module.zip_processor_lambda.lambda_arn
  description = "ARN da função Lambda zip-processor"
}
