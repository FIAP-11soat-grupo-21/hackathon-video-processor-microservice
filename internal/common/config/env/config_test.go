package env

import (
	"os"
	"sync"
	"testing"
)

func resetConfig() {
	config = nil
	once = sync.Once{}
}

func TestGetConfig_WithDefaults(t *testing.T) {
	resetConfig()
	os.Clearenv()
	
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "test-queue")

	cfg := GetConfig()

	if cfg == nil {
		t.Fatal("Expected config to be initialized")
	}

	if cfg.Environment != "development" {
		t.Errorf("Expected Environment to be 'development', got '%s'", cfg.Environment)
	}

	if cfg.API.Host != "0.0.0.0" {
		t.Errorf("Expected API.Host to be '0.0.0.0', got '%s'", cfg.API.Host)
	}

	if cfg.API.Port != "8080" {
		t.Errorf("Expected API.Port to be '8080', got '%s'", cfg.API.Port)
	}

	expectedURL := "0.0.0.0:8080"
	if cfg.API.URL != expectedURL {
		t.Errorf("Expected API.URL to be '%s', got '%s'", expectedURL, cfg.API.URL)
	}

	if cfg.AWS.Region != "us-east-1" {
		t.Errorf("Expected AWS.Region to be 'us-east-1', got '%s'", cfg.AWS.Region)
	}

	if cfg.S3.Bucket != "fiap-tc-terraform-846874" {
		t.Errorf("Expected S3.Bucket to be 'fiap-tc-terraform-846874', got '%s'", cfg.S3.Bucket)
	}
}

func TestGetConfig_WithCustomEnvironmentVariables(t *testing.T) {
	resetConfig()
	os.Clearenv()

	os.Setenv("GO_ENV", "production")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("API_PORT", "3000")
	os.Setenv("AWS_REGION", "sa-east-1")
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "test-queue")
	os.Setenv("S3_BUCKET", "test-bucket")

	cfg := GetConfig()

	if cfg.Environment != "production" {
		t.Errorf("Expected Environment to be 'production', got '%s'", cfg.Environment)
	}

	if cfg.API.Host != "localhost" {
		t.Errorf("Expected API.Host to be 'localhost', got '%s'", cfg.API.Host)
	}

	if cfg.API.Port != "3000" {
		t.Errorf("Expected API.Port to be '3000', got '%s'", cfg.API.Port)
	}

	expectedURL := "localhost:3000"
	if cfg.API.URL != expectedURL {
		t.Errorf("Expected API.URL to be '%s', got '%s'", expectedURL, cfg.API.URL)
	}

	if cfg.AWS.Region != "sa-east-1" {
		t.Errorf("Expected AWS.Region to be 'sa-east-1', got '%s'", cfg.AWS.Region)
	}

	if cfg.AWS.Endpoint != "http://localhost:4566" {
		t.Errorf("Expected AWS.Endpoint to be 'http://localhost:4566', got '%s'", cfg.AWS.Endpoint)
	}

	if cfg.AWS.SQS.Queues.FrameExtraction != "test-queue" {
		t.Errorf("Expected AWS.SQS.Queues.FrameExtraction to be 'test-queue', got '%s'", cfg.AWS.SQS.Queues.FrameExtraction)
	}

	if cfg.S3.Bucket != "test-bucket" {
		t.Errorf("Expected S3.Bucket to be 'test-bucket', got '%s'", cfg.S3.Bucket)
	}
}

func TestGetConfig_Singleton(t *testing.T) {
	resetConfig()
	os.Clearenv()
	
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "test-queue")

	cfg1 := GetConfig()
	cfg2 := GetConfig()

	if cfg1 != cfg2 {
		t.Error("Expected GetConfig to return the same instance (singleton pattern)")
	}
}

func TestIsProduction_WhenProduction(t *testing.T) {
	resetConfig()
	os.Clearenv()
	os.Setenv("GO_ENV", "production")
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "test-queue")

	cfg := GetConfig()

	if !cfg.IsProduction() {
		t.Error("Expected IsProduction to return true for production environment")
	}
}

func TestIsProduction_WhenDevelopment(t *testing.T) {
	resetConfig()
	os.Clearenv()
	os.Setenv("GO_ENV", "development")
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "test-queue")

	cfg := GetConfig()

	if cfg.IsProduction() {
		t.Error("Expected IsProduction to return false for development environment")
	}
}

func TestIsProduction_WhenStaging(t *testing.T) {
	resetConfig()
	os.Clearenv()
	os.Setenv("GO_ENV", "staging")
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "test-queue")

	cfg := GetConfig()

	if cfg.IsProduction() {
		t.Error("Expected IsProduction to return false for staging environment")
	}
}

func TestGetEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")

	if result != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", result)
	}
}

func TestGetEnv_WithDefault(t *testing.T) {
	os.Unsetenv("NON_EXISTENT_VAR")

	result := getEnv("NON_EXISTENT_VAR", "default-value")

	if result != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", result)
	}
}

func TestConfig_AllFieldsPopulated(t *testing.T) {
	resetConfig()
	os.Clearenv()
	
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_HOST", "testhost")
	os.Setenv("API_PORT", "9999")
	os.Setenv("AWS_REGION", "eu-west-1")
	os.Setenv("AWS_ENDPOINT", "https://aws.amazon.com")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "frame-queue-prod")
	os.Setenv("S3_BUCKET", "production-bucket")

	cfg := GetConfig()

	if cfg.Environment != "test" {
		t.Errorf("Expected Environment 'test', got '%s'", cfg.Environment)
	}
	
	if cfg.API.Host != "testhost" {
		t.Errorf("Expected API.Host 'testhost', got '%s'", cfg.API.Host)
	}
	
	if cfg.API.Port != "9999" {
		t.Errorf("Expected API.Port '9999', got '%s'", cfg.API.Port)
	}
	
	if cfg.API.URL != "testhost:9999" {
		t.Errorf("Expected API.URL 'testhost:9999', got '%s'", cfg.API.URL)
	}
	
	if cfg.AWS.Region != "eu-west-1" {
		t.Errorf("Expected AWS.Region 'eu-west-1', got '%s'", cfg.AWS.Region)
	}
	
	if cfg.AWS.Endpoint != "https://aws.amazon.com" {
		t.Errorf("Expected AWS.Endpoint 'https://aws.amazon.com', got '%s'", cfg.AWS.Endpoint)
	}
	
	if cfg.AWS.SQS.Queues.FrameExtraction != "frame-queue-prod" {
		t.Errorf("Expected FrameExtraction queue 'frame-queue-prod', got '%s'", cfg.AWS.SQS.Queues.FrameExtraction)
	}
	
	if cfg.S3.Bucket != "production-bucket" {
		t.Errorf("Expected S3.Bucket 'production-bucket', got '%s'", cfg.S3.Bucket)
	}
}
