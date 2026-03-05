# ECS Service Variables
variable "application_name" {
  description = "Nome da aplicação ECS"
  type        = string
}

variable "image_name" {
  description = "Nome da imagem do container"
  type        = string
}

variable "image_port" {
  description = "Porta do container"
  type        = number
}

variable "desired_count" {
  description = "Número desejado de tarefas ECS"
  type        = number
  default     = 1
}

variable "container_environment_variables" {
  description = "Variáveis de ambiente do container"
  type        = map(string)
  default     = {}
}

variable "container_secrets" {
  description = "Segredos do container"
  type        = map(string)
  default     = {}
}

variable "health_check_path" {
  description = "Caminho de verificação de integridade do serviço"
  type        = string
  default     = "/health"
}

variable "task_role_policy_arns" {
  description = "Lista de ARNs de políticas para anexar à função da tarefa ECS"
  type        = list(string)
  default     = []
}

# API Gateway Variables
variable "apigw_integration_type" {
  description = "Tipo de integração do API Gateway"
  type        = string
  default     = "HTTP_PROXY"
}

variable "apigw_integration_method" {
  description = "Método de integração do API Gateway"
  type        = string
  default     = "ANY"
}

variable "apigw_payload_format_version" {
  description = "Versão do payload do API Gateway"
  type        = string
  default     = "1.0"
}

variable "apigw_connection_type" {
  description = "Tipo de conexão do API Gateway"
  type        = string
  default     = "VPC_LINK"
}

variable "app_path_pattern" {
  description = "Lista de padrões de caminho para o listener rule do ALB"
  type        = list(string)
}

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