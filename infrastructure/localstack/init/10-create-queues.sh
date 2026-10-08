#!/bin/sh

set -eu

endpoint="http://localhost:4566"
region="${AWS_DEFAULT_REGION:-sa-east-1}"
dlq_name="${WAGER_TRANSACTIONS_DLQ:-wager-transactions-dlq.fifo}"
queue_name="${WAGER_TRANSACTIONS_QUEUE:-wager-transactions.fifo}"

fifo_attributes='{"FifoQueue":"true","ContentBasedDeduplication":"true"}'

awslocal --endpoint-url "$endpoint" --region "$region" sqs create-queue \
  --queue-name "$dlq_name" \
  --attributes "$fifo_attributes" >/dev/null

dlq_url="$(awslocal --endpoint-url "$endpoint" --region "$region" sqs get-queue-url \
  --queue-name "$dlq_name" --query QueueUrl --output text)"
dlq_arn="$(awslocal --endpoint-url "$endpoint" --region "$region" sqs get-queue-attributes \
  --queue-url "$dlq_url" --attribute-names QueueArn --query Attributes.QueueArn --output text)"

redrive_policy="$(printf '{"deadLetterTargetArn":"%s","maxReceiveCount":"5"}' "$dlq_arn")"
escaped_redrive_policy="$(printf '%s' "$redrive_policy" | sed 's/"/\\"/g')"
queue_attributes="$(printf '{"FifoQueue":"true","ContentBasedDeduplication":"true","RedrivePolicy":"%s"}' "$escaped_redrive_policy")"

awslocal --endpoint-url "$endpoint" --region "$region" sqs create-queue \
  --queue-name "$queue_name" \
  --attributes "$queue_attributes" >/dev/null

echo "LocalStack SQS queues ready: $queue_name, $dlq_name"