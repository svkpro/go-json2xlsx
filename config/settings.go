package config

import(
	"encoding/json"
	"fmt"
	"os"
)

type Settings struct {
	Environment   string `json:"ENVIRONMENT"`
	LoggingLevel  string `json:"LOGGING_LEVEL"`
	HttpPort      string `json:"HTTP_PORT"`

	AwsURL        string `json:"AWS_URL"`
	AwsAccessKey  string `json:"AWS_ACCESS_KEY"`
	AwsSecretKey  string `json:"AWS_SECRET_KEY"`
	AwsBucketName string `json:"AWS_BUCKET_NAME"`
	AwsRegion     string `json:"AWS_REGION"`
}

func New() Settings {
	file, err := os.Open("config/conf.json")
	defer file.Close()

	if err != nil {
		fmt.Println("Error:", err)
	}

	d := json.NewDecoder(file)
	s := Settings{}
	err = d.Decode(&s)

	if err != nil {
		fmt.Println("Error:", err)
	}

	return s
}