# Lambda Variables
variable "lambda_environment_variables" {
  description = "Variáveis de ambiente da função Lambda"
  type        = map(string)
  default     = {}
}

variable "lambda_bucket_name" {
  description = "Nome do bucket S3 para código da Lambda"
  type        = string
  default     = "fiap-hackathon-lambda-content-44573"
}

# S3 Variables
variable "s3_bucket" {
  description = "Bucket S3 para vídeos e frames"
  type        = string
  default     = "fiap-tc-terraform-846874"
}

# DynamoDB Variables
variable "dynamodb_table_name" {
  description = "Nome da tabela DynamoDB para status dos vídeos"
  type        = string
  default     = "video-chunks"
}

# Chunk Processor Lambda
variable "chunk_processor_s3_key" {
  description = "Chave S3 do arquivo Lambda chunk-processor"
  type        = string
  default     = "chunk-processor.zip"
}

variable "chunk_processor_memory_size" {
  description = "Memória da função Lambda chunk-processor em MB"
  type        = number
  default     = 3008
}

variable "chunk_processor_timeout" {
  description = "Timeout da função Lambda chunk-processor em segundos"
  type        = number
  default     = 900
}

# Update Video Chunk Status Lambda
variable "update_status_s3_key" {
  description = "Chave S3 do arquivo Lambda update-video-chunk-status"
  type        = string
  default     = "update-video-chunk-status.zip"
}

variable "update_status_memory_size" {
  description = "Memória da função Lambda update-video-chunk-status em MB"
  type        = number
  default     = 512
}

variable "update_status_timeout" {
  description = "Timeout da função Lambda update-video-chunk-status em segundos"
  type        = number
  default     = 60
}

# Zip Processor Lambda
variable "zip_processor_s3_key" {
  description = "Chave S3 do arquivo Lambda zip-processor"
  type        = string
  default     = "zip-processor.zip"
}

variable "zip_processor_memory_size" {
  description = "Memória da função Lambda zip-processor em MB"
  type        = number
  default     = 2048
}

variable "zip_processor_timeout" {
  description = "Timeout da função Lambda zip-processor em segundos"
  type        = number
  default     = 900
}