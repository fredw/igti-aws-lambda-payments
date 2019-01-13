# Variables
OUTPUT=out/main
LAMBDA_NAME=payments
SQS_QUEUE_NAME=payments.fifo
# Colors
GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m

all: clean build update

# Build the project
build:
	@echo "${GREEN}* Building...${NC}"
	@GOOS=linux go build -o ${OUTPUT} main.go
	@zip -r -j ${OUTPUT}.zip ./${OUTPUT} > /dev/null

# Invoke the Lambda function on AWS
# Usage: make invoke debug=1
invoke:
	@echo "${GREEN}* Invoking...${NC}"
	@if [ "$(debug)" == "1" ]; then \
		echo "${GREEN}* Logs:${NC}"; \
		aws lambda invoke --function-name ${LAMBDA_NAME} --log-type Tail --payload '' out/response.txt | jq -r .LogResult | base64 --decode; \
	else \
		aws lambda invoke --function-name ${LAMBDA_NAME} --payload '' out/response.txt | jq; \
	fi
	@echo "${GREEN}* Lambda response:${NC}"
	@cat out/response.txt

# Update the Lambda function on AWS
update:
	@echo "${GREEN}* Updating...${NC}"
	@aws lambda update-function-code --function-name ${LAMBDA_NAME} --zip-file fileb://${OUTPUT}.zip > /dev/null

# Run tests
test:
	go vet ./...
	go test -v -cover ./...

# Clean
clean:
	@rm -rf out/*

# SQS: receive messages with no visibility timeout
sqs_receive_messages:
	@aws sqs receive-message \
		--queue-url ${shell aws sqs get-queue-url --queue-name ${SQS_QUEUE_NAME} | jq -r .QueueUrl} \
		--visibility-timeout 0 \
		--max-number-of-messages 10 \
		| jq

# SQS: create a test message
sqs_create_test_message:
	@aws sqs send-message \
		--queue-url ${shell aws sqs get-queue-url --queue-name ${SQS_QUEUE_NAME} | jq -r .QueueUrl} \
		--message-body "{\"provider\": \"Example\"}" \
		--message-group-id ${shell openssl rand -base64 6} \
		--message-deduplication-id ${shell openssl rand -base64 6} \
		| jq

# SQS: purge the queue (delete all messages)
sqs_purge_queue:
	@aws sqs purge-queue --queue-url ${shell aws sqs get-queue-url --queue-name ${SQS_QUEUE_NAME} | jq -r .QueueUrl}