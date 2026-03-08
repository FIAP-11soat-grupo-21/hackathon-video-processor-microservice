resource "aws_lambda_layer_version" "ffmpeg" {
  layer_name          = "ffmpeg"
  description         = "FFmpeg for video processing"
  s3_bucket           = var.lambda_bucket_name
  s3_key              = "ffmpeg-layer.zip"
  compatible_runtimes = ["provided.al2023", "provided.al2"]
}

resource "null_resource" "attach_ffmpeg_layer" {
  depends_on = [
    module.chunk_processor_lambda,
    aws_lambda_layer_version.ffmpeg
  ]

  triggers = {
    layer_version = aws_lambda_layer_version.ffmpeg.version
    lambda_name   = "chunk-processor"
  }

  provisioner "local-exec" {
    command = <<-EOT
      aws lambda update-function-configuration \
        --function-name ${self.triggers.lambda_name} \
        --layers ${aws_lambda_layer_version.ffmpeg.arn}
    EOT
  }
}
