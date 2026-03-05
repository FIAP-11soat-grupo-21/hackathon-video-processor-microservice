output "frame_queue_url" {
  value       = aws_sqs_queue.frame_queue.url
  description = "URL da fila SQS para frames"
}

output "frame_queue_arn" {
  value       = aws_sqs_queue.frame_queue.arn
  description = "ARN da fila SQS para frames"
}

output "frame_dlq_url" {
  value       = aws_sqs_queue.frame_dlq.url
  description = "URL da Dead Letter Queue"
}

output "frame_dlq_arn" {
  value       = aws_sqs_queue.frame_dlq.arn
  description = "ARN da Dead Letter Queue"
}

output "lambda_arn" {
  value       = module.frame_processor_lambda.lambda_arn
  description = "ARN da função Lambda"
}
