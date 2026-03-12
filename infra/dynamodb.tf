resource "aws_dynamodb_table" "video_chunks" {
  name           = var.dynamodb_table_name
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "video_id"
  range_key      = "chunk_part"

  attribute {
    name = "video_id"
    type = "S"
  }

  attribute {
    name = "chunk_part"
    type = "N"
  }

  tags = merge(
    data.terraform_remote_state.app_registry.outputs.app_registry_application_tag,
    {
      Name        = var.dynamodb_table_name
      Environment = "production"
      ManagedBy   = "terraform"
    }
  )
}
