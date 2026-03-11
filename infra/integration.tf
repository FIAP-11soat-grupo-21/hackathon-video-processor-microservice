resource "aws_lambda_event_source_mapping" "chunk_processor_trigger" {
  event_source_arn = data.terraform_remote_state.sqs_chunk_processor.outputs.sqs_queue_arn
  function_name    = module.chunk_processor_lambda.lambda_arn
  batch_size       = 1
  enabled          = var.lambda_event_source_enabled
}

resource "aws_lambda_event_source_mapping" "update_status_trigger" {
  event_source_arn = data.terraform_remote_state.sqs_update_video_chunk_status.outputs.sqs_queue_arn
  function_name    = module.update_video_chunk_status_lambda.lambda_arn
  batch_size       = 1
  enabled          = var.lambda_event_source_enabled
}

resource "aws_lambda_event_source_mapping" "zip_processor_trigger" {
  event_source_arn = data.terraform_remote_state.sqs_zip_processor.outputs.sqs_queue_arn
  function_name    = module.zip_processor_lambda.lambda_arn
  batch_size       = 1
  enabled          = var.lambda_event_source_enabled
}
