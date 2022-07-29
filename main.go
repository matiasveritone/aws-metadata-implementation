package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"log"
)

func main() {
	log.Printf("Starting process")

	// defining context and parameters
	ctx := context.TODO()
	isAWSAccount := true
	region := "us-east-1"
	regionConfig := config.WithRegion(
		region,
	)

	// checking if we are on AWS EC2 before using AWS Account credentials
	if isAWSAccount {
		log.Printf("Checking for AWS Account inside of AWS EC2 environment.")

		// creating metadata client
		err := startMetadataClient(ctx)
		if err != nil {
			log.Printf("ERROR: Unable to create AWS Metadata Client: %s\n", err)
			return
		}

	}

	// starting AWS Comprehend client
	clientComprehend, err := startComprehendClient(ctx, isAWSAccount, regionConfig)
	if err != nil {
		log.Printf("ERROR: Unable to create AWS Comprehend Client: %s\n", err)
		return
	}

	// executing API
	// create parameters for DetectDominantLanguage API
	textSample := "This is just a test with English as dominant language"

	detectDominantLanguageInput := comprehend.DetectDominantLanguageInput{
		Text: &textSample,
	}

	_, err = clientComprehend.DetectDominantLanguage(ctx, &detectDominantLanguageInput)
	if err != nil {
		log.Printf("ERROR: Something failed while running aws API. %v", err)
		return
	} else {
		log.Printf("Everything worked fine, connection was successful. ")
	}
}

func startMetadataClient(ctx context.Context) error {
	log.Printf("Creating metadata client")

	// basic config for EC2 metadata
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("ERROR: loading config failed: %v", err)
		return err
	}

	// creating EC2 metadata client
	clientMetadata := imds.NewFromConfig(cfg)

	// checking for local EC2 instance name
	log.Printf("Checking for metadata info (only works inside AWS EC2 Environment).")
	instanceId, err := clientMetadata.GetMetadata(ctx, &imds.GetMetadataInput{
		Path: "instance-id",
	})
	if err != nil {
		log.Printf("ERROR: Unable to retrieve the EC2 instance name: %s\n", err)
		return err
	}

	log.Printf("Connected into AWS EC2 instance: %v\n", instanceId)

	return nil

}

// startComprehendClient defines whether we create the client using a credentials file
// or checking if we are in an AWS EC2 environment that uses a AWS Account
func startComprehendClient(ctx context.Context, metadata bool, region config.LoadOptionsFunc) (comprehend.Client, error) {
	var err error
	var cfg aws.Config
	var client *comprehend.Client
	var credsConfig config.LoadOptionsFunc

	// creating client
	if !metadata {
		log.Printf("Creating AWS Comprehend Client with credentials file.")

		// when not in EC2 environment, a file with correspondant credentials is used
		credsConfig = config.WithSharedCredentialsFiles(
			[]string{"./credentials/credentials"},
		)

		// load the Shared AWS Configuration and credentials
		cfg, err = config.LoadDefaultConfig(ctx,
			credsConfig,
			region,
		)
		if err != nil {
			log.Printf("ERROR: failed to load configuration with custom credentials, %v\n", err)
			return comprehend.Client{}, err
		} else {
			log.Println("AWS Client configuration loaded successfully")
		}

	} else {
		log.Printf("Creating AWS Comprehend Client with AWS Account inside EC2 environment.")

		// basic config for EC2 metadata
		cfg, err = config.LoadDefaultConfig(
			ctx,
			region,
		)
		if err != nil {
			log.Printf("ERROR: failed to load configuration within AWS EC2 environment: %v", err)
			return comprehend.Client{}, err
		} else {
			log.Println("AWS Comprehend Client configuration loaded successfully")
		}

	}

	// Using the respective config value to create Comprehend client
	client = comprehend.NewFromConfig(cfg)
	if err != nil {
		log.Printf("ERROR: failed to create comprehend client, %v\n", err)
		return comprehend.Client{}, err
	} else {
		log.Printf("AWS Comprehend client created successfully")
		return *client, nil
	}
}
