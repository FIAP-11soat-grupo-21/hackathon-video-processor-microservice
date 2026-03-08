# Video Processor Service

Microserviço de processamento de vídeo usando AWS Lambda com arquitetura hexagonal.

## Lambdas

### chunk-processor
Processa chunks de vídeo usando FFmpeg, extrai frames e salva no S3.
Publica mensagem quando chunk é processado.

### update-video-chunk-status
Atualiza status do chunk no DynamoDB após processamento.
Quando todos os chunks estão processados, publica mensagem para iniciar compactação.

### zip-processor
Compacta todos os frames processados em um arquivo ZIP.
Salva o ZIP no S3 e publica mensagem de conclusão.

## Fluxo de Mensageria

### SNS Topics
- `chunk-uploaded` → Trigger para chunk-processor
- `chunk-processed` → Trigger para update-video-chunk-status
- `all-chunks-processed` → Trigger para zip-processor
- `video-processing-complete` → Notificação final
- `video-processed-error` → Erros de processamento

### SQS Queues
- `chunk-processor` → Consome de chunk-uploaded
- `update-video-chunk-status` → Consome de chunk-processed
- `zip-processor` → Consome de all-chunks-processed
- `update-video-status` → Consome de video-processing-complete
- `notificate-user` → Notificações ao usuário

### Fluxo
```
SNS: chunk-uploaded
  ↓
SQS: chunk-processor → Lambda: chunk-processor
  ↓
SNS: chunk-processed
  ↓
SQS: update-video-chunk-status → Lambda: update-video-chunk-status
  ↓
SNS: all-chunks-processed
  ↓
SQS: zip-processor → Lambda: zip-processor
  ↓
SNS: video-processing-complete
```

## Estrutura do Projeto

```
├── cmd/lambda/                    # Entrypoints das Lambdas (apenas bootstrap)
│   ├── chunk-processor/
│   ├── update-video-chunk-status/
│   └── zip-processor/
├── internal/
│   ├── adapters/
│   │   ├── driven/                # Implementações (S3, SNS/SQS, DynamoDB, FFmpeg)
│   │   └── driver/lambda/         # Handlers das Lambdas (lógica de negócio)
│   ├── core/
│   │   ├── domain/ports/          # Interfaces
│   │   ├── dto/                   # Data Transfer Objects
│   │   ├── use_cases/             # Casos de uso
│   │   └── factory/               # Injeção de dependências
│   └── common/config/             # Configurações
└── infra/                         # Terraform (IaC)
```
