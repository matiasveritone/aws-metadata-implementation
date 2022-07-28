package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
)
import "log"
import "github.com/aws/aws-sdk-go-v2/config"

func main() {
	// defining context
	ctx := context.TODO()
	checkAWSmachine := true

	// checking if metadata exists
	if checkAWSmachine {
		log.Printf("Checking for metadata of AWS EC2.")

		// creating metadata client
		err := startMetadataClient(ctx)
		if err != nil {
			log.Printf("ERROR: Unable to create AWS Metadata Client: %s\n", err)
			return
		}

	}

	// starting AWS Comprehend client
	clientComprehend, err := startComprehendClient(ctx, checkAWSmachine)
	if err != nil {
		log.Printf("ERROR: Unable to create AWS Comprehend Client: %s\n", err)
		return
	}

	// executing API
	// create parameters for DetectDominantLanguage API
	text := "English as dominant language"

	detectDominantLanguageInput := comprehend.DetectDominantLanguageInput{
		Text: &text,
	}

	_, err = clientComprehend.DetectDominantLanguage(ctx, &detectDominantLanguageInput)
	if err != nil {
		log.Printf("ERROR: Something failed while running aws API. %v", err)
		return
	}

	log.Printf("Everything worked fine. ")
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

func startComprehendClient(ctx context.Context, metadata bool) (comprehend.Client, error) {
	var client *comprehend.Client

	// checking if custom credentials exist
	var credsConfig config.LoadOptionsFunc

	// loading creds if not in metadata
	if !metadata {
		log.Printf("Creating AWS Comprehend Client with credentials file.")

		credsConfig = config.WithSharedCredentialsFiles(
			[]string{"./credentials/credentials"},
		)

		// Load the Shared AWS Configuration and Credentials
		cfg, err := config.LoadDefaultConfig(ctx,
			credsConfig,
			config.WithRegion(
				"us-east-1",
			))
		if err != nil {
			log.Printf("Error: failed to load configuration, %v\n", err)
			return comprehend.Client{}, err
		} else {
			log.Println("AWS Client configuration loaded successfully")
		}

		// Using the Config value, to create Comprehend client
		client = comprehend.NewFromConfig(cfg)
		if err != nil {
			log.Printf("Error: failed to create client, %v\n", err)
			return *client, err
		}

	} else {
		log.Printf("Creating AWS Comprehend Client with metadata.")

		// basic config for EC2 metadata
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Printf("ERROR: loading config failed: %v", err)
			return comprehend.Client{}, err
		}

		// Using the Config value, to create Comprehend client
		client = comprehend.NewFromConfig(cfg)
		if err != nil {
			log.Printf("Error: failed to create client, %v\n", err)
			return *client, err
		}

	}

	return *client, nil

}
