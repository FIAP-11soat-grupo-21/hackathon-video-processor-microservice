.PHONY: test coverage coverage-html clean

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

clean:
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	find ./cmd/lambda -type f -name "bootstrap" -delete
	find ./cmd/lambda -type f -name "*.zip" -delete

