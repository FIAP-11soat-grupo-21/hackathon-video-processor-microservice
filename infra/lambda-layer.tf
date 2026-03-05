resource "aws_lambda_layer_version" "ffmpeg" {
  layer_name          = "ffmpeg"
  description         = "FFmpeg for video processing"
  s3_bucket           = var.lambda_bucket_name
  s3_key              = "ffmpeg-layer.zip"
  compatible_runtimes = ["provided.al2023", "provided.al2"]
}

data "aws_s3_object" "lambda_code" {
  bucket = var.lambda_bucket_name
  key    = var.lambda_s3_key
}

resource "null_resource" "update_lambda_code" {
  depends_on = [module.frame_processor_lambda]

  triggers = {
    code_hash   = data.aws_s3_object.lambda_code.etag
    lambda_name = "video-frame-processor"
  }

  provisioner "local-exec" {
    command = <<-EOT
      aws lambda update-function-code \
        --function-name ${self.triggers.lambda_name} \
        --s3-bucket ${var.lambda_bucket_name} \
        --s3-key ${var.lambda_s3_key}
    EOT
  }
}

resource "null_resource" "attach_layer_to_lambda" {
  depends_on = [
    null_resource.update_lambda_code,
    aws_lambda_layer_version.ffmpeg
  ]

  triggers = {
    layer_version = aws_lambda_layer_version.ffmpeg.version
    lambda_name   = "video-frame-processor"
  }

  provisioner "local-exec" {
    command = <<-EOT
      aws lambda update-function-configuration \
        --function-name ${self.triggers.lambda_name} \
        --layers ${aws_lambda_layer_version.ffmpeg.arn}
    EOT
  }
}
