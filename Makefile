run:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=count
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total

sonar-scan:
	go test ./... -coverprofile=coverage.out -covermode=count
	golangci-lint run --config .golangci.yml --out-format json > golangci-lint-report.json || true
	sonar-scanner

build-api:
	docker build -t video-processor-api .

build-worker:
	docker build -t video-processor-worker .

build-lambda:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap ./cmd/lambda
	zip lambda-handler.zip bootstrap

clean:
	rm -rf bin/ bootstrap lambda-handler.zip

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f
