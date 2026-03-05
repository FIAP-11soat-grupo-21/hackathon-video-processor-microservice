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

output "lambda_integration_id" {
  value       = module.frame_processor_lambda.lambda_integration_id
  description = "ID da integração Lambda com API Gateway"
}


output "target_group_arn" {
  value       = aws_alb_target_group.target_group.arn
  description = "ARN do Target Group"
}

output "alb_listener_arn" {
  value       = aws_lb_listener.listener.arn
  description = "ARN do Listener do ALB"
}
