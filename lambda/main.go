package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/bhagashetti/job-alert-system/scraper"
)

func HandleRequest(ctx context.Context) (string, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("config error: %v", err)
	}

	// Create SNS client
	snsClient := sns.NewFromConfig(cfg)

	// Get SNS topic ARN from environment variable
	topicARN := os.Getenv("SNS_TOPIC_ARN")
	if topicARN == "" {
		return "", fmt.Errorf("SNS_TOPIC_ARN environment variable is not set")
	}

	// Scrape jobs using broad keyword for testing
	jobs, err := scraper.ScrapeJobs("software developer")
	if err != nil {
		return "", fmt.Errorf("scraper error: %v", err)
	}

	// Loop through jobs and send alerts
	var sentCount int
	for _, job := range jobs {
		// Construct message
		msg := fmt.Sprintf("📢 Job Alert!\nTitle: %s\nCompany: %s\nLink: %s",
			job.Title, job.Company, job.URL)

		// Publish to SNS
		_, err := snsClient.Publish(ctx, &sns.PublishInput{
			Message:  aws.String(msg),
			TopicArn: aws.String(topicARN),
		})
		if err != nil {
			fmt.Printf("Failed to send alert for job: %s | Error: %v\n", job.Title, err)
			continue
		}
		sentCount++
	}

	return fmt.Sprintf("Job scan complete. %d alerts sent.", sentCount), nil
}

func main() {
	lambda.Start(HandleRequest)
}
