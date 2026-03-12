package env

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	API         APIConfig
	AWS         AWSConfig
	S3          S3Config
}

type APIConfig struct {
	Host string
	Port string
	URL  string
}

type AWSConfig struct {
	Region   string
	Endpoint string
	SQS      SQSConfig
}

type SQSConfig struct {
	Queues QueuesConfig
}

type QueuesConfig struct {
	FrameExtraction string
}

type S3Config struct {
	Bucket string
}

var (
	config *Config
	once   sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		apiHost := getEnv("API_HOST", "0.0.0.0")
		apiPort := getEnv("API_PORT", "8080")

		config = &Config{
			Environment: getEnv("GO_ENV", "development"),
			API: APIConfig{
				Host: apiHost,
				Port: apiPort,
				URL:  apiHost + ":" + apiPort,
			},
			AWS: AWSConfig{
				Region:   getEnvOptional("AWS_REGION", "us-east-1"),
				Endpoint: getEnvOptional("AWS_ENDPOINT", ""),
				SQS: SQSConfig{
					Queues: QueuesConfig{
						FrameExtraction: getEnvOptional("AWS_SQS_FRAME_EXTRACTION_QUEUE", ""),
					},
				},
			},
			S3: S3Config{
				Bucket: getEnvOptional("S3_BUCKET", "fiap-hackathon-chunk-bucket-54351"),
			},
		}
	})

	return config
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			log.Fatalf("Environment variable %s is required", key)
		}
		return defaultValue
	}
	return value
}

func getEnvOptional(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
