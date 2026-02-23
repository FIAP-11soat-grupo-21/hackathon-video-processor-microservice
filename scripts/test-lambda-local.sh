#!/bin/bash

# Script para testar a Lambda localmente
# Uso: ./scripts/test-lambda-local.sh

set -e

echo "🧪 Testing Lambda locally..."

# Configurar variáveis de ambiente de teste
export AWS_REGION="us-east-2"
export S3_INPUT_BUCKET="test-input-bucket"
export S3_OUTPUT_BUCKET="test-output-bucket"
export FRAME_QUEUE_URL="https://sqs.us-east-2.amazonaws.com/123456/test-queue"
export LOG_LEVEL="DEBUG"

# Criar evento de teste
cat > /tmp/lambda-test-event.json <<EOF
{
  "Records": [
    {
      "messageId": "test-message-1",
      "body": "{\"jobId\":\"test-job-123\",\"bucket\":\"test-bucket\",\"key\":\"test-video.mp4\",\"timestamp\":1.0,\"index\":0}"
    },
    {
      "messageId": "test-message-2",
      "body": "{\"jobId\":\"test-job-123\",\"bucket\":\"test-bucket\",\"key\":\"test-video.mp4\",\"timestamp\":5.0,\"index\":1}"
    }
  ]
}
EOF

echo "📝 Test event created at /tmp/lambda-test-event.json"
echo ""
echo "Event content:"
cat /tmp/lambda-test-event.json | jq .
echo ""

# Executar Lambda com evento de teste
echo "🚀 Executing Lambda handler..."
echo ""

go run cmd/lambda/main.go < /tmp/lambda-test-event.json

echo ""
echo "✅ Lambda test completed!"
