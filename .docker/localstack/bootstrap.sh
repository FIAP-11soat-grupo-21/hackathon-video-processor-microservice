#!/bin/bash

awslocal sqs create-queue \
  --queue-name frame-extraction-queue \
  --region us-east-2

awslocal s3 mb s3://video-frames-bucket --region us-east-2

awslocal sqs list-queues --region us-east-2
awslocal s3 ls
