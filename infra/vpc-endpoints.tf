data "aws_vpc" "main" {
  id = data.terraform_remote_state.network_vpc.outputs.vpc_id
}

resource "aws_vpc_endpoint" "sns" {
  vpc_id              = data.terraform_remote_state.network_vpc.outputs.vpc_id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.sns"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = data.terraform_remote_state.network_vpc.outputs.private_subnets
  security_group_ids  = [aws_security_group.vpc_endpoints.id]
  private_dns_enabled = true

  tags = merge(
    data.terraform_remote_state.app_registry.outputs.app_registry_application_tag,
    {
      Name        = "video-processor-sns-endpoint"
      Environment = "production"
      ManagedBy   = "terraform"
    }
  )
}

resource "aws_security_group" "vpc_endpoints" {
  name        = "video-processor-vpc-endpoints-sg"
  description = "Security group for VPC endpoints"
  vpc_id      = data.terraform_remote_state.network_vpc.outputs.vpc_id

  ingress {
    description = "HTTPS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [data.aws_vpc.main.cidr_block]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    data.terraform_remote_state.app_registry.outputs.app_registry_application_tag,
    {
      Name        = "video-processor-vpc-endpoints-sg"
      Environment = "production"
      ManagedBy   = "terraform"
    }
  )
}
