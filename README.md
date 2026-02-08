# 🎬 Video Processor Service

## Arquitetura

```
API REST (Orquestrador)
  ↓
Calcula timestamps baseado no intervalo
  ↓
Enfileira mensagens no SQS (1 por frame)
  ↓
Lambda Workers (paralelo)
  ↓
Extrai frames e salva no S3
```
