resource "aws_sqs_queue" "frame_queue" {
  name                       = var.sqs_queue_name
  delay_seconds              = 0
  max_message_size           = 262144
  message_retention_seconds  = 1209600
  receive_wait_time_seconds  = 0
  visibility_timeout_seconds = 900

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.frame_dlq.arn
    maxReceiveCount     = 3
  })

  tags = {
    Project     = "video-processor-lambda"
    Environment = "production"
    ManagedBy   = "terraform"
  }
}

resource "aws_sqs_queue" "frame_dlq" {
  name                      = "${var.sqs_queue_name}-dlq"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 1209600

  tags = {
    Project     = "video-processor-lambda"
    Environment = "production"
    ManagedBy   = "terraform"
  }
}

resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = aws_sqs_queue.frame_queue.arn
  function_name    = module.frame_processor_lambda.lambda_arn
  batch_size       = 1
  enabled          = true
}
