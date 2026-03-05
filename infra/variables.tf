# Lambda Variables
variable "lambda_environment_variables" {
  description = "Variáveis de ambiente da função Lambda"
  type        = map(string)
  default     = {}
}

variable "lambda_bucket_name" {
  description = "Nome do bucket S3 para código da Lambda"
  type        = string
  default     = "fiap-tc-terraform-functions-846874"
}

variable "lambda_s3_key" {
  description = "Chave S3 do arquivo Lambda"
  type        = string
  default     = "video-frame-processor.zip"
}

variable "lambda_memory_size" {
  description = "Memória da função Lambda em MB"
  type        = number
  default     = 2048
}

variable "lambda_timeout" {
  description = "Timeout da função Lambda em segundos"
  type        = number
  default     = 300
}

# SQS Variables
variable "sqs_queue_name" {
  description = "Nome da fila SQS para frames"
  type        = string
  default     = "video-frame-queue"
}

# S3 Variables
variable "s3_bucket" {
  description = "Bucket S3 para vídeos e frames"
  type        = string
  default     = "fiap-tc-terraform-846874"
}