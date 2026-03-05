.PHONY: test coverage coverage-html build-lambda clean

LAMBDA_BUCKET=fiap-tc-terraform-functions-846874
LAMBDA_KEY=video-frame-processor.zip
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

test:
	go test -v ./...

coverage:
	go test ./... -coverprofile=$(COVERAGE_FILE)
	go tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

build-lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap ./cmd/lambda
	zip -j lambda-handler.zip bootstrap

clean:
	rm -f bootstrap lambda-handler.zip $(COVERAGE_FILE) $(COVERAGE_HTML)
