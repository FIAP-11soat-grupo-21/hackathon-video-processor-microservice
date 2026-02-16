run:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

test:
	go test ./...

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
