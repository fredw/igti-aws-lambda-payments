# Variables
OUTPUT=out/main
LAMBDA_NAME=payments
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
