package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
)
import "log"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

func main() {
	// defining context
	ctx := context.TODO()

	// basic config for EC2 metadata
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	// creating EC2 metadata client
	log.Printf("Creating metadata client")
	clientMetadata := imds.NewFromConfig(cfg)

	// checking for local EC2 ip
	instanceId, err := clientMetadata.GetMetadata(ctx, &imds.GetMetadataInput{
		Path: "instance-id",
	})
	if err != nil {
		log.Printf("ERROR: Unable to retrieve the EC2 instance name: %s\n", err)
	}

	log.Printf("Connected into AWS EC2 instance: %v\n", instanceId)

	// starting AWS Comprehend client
	clientComprehend, err := StartClient(ctx)
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

func StartClient(ctx context.Context) (comprehend.Client, error) {

	// checking if custom credentials exist
	var credsConfig config.LoadOptionsFunc

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
	client := comprehend.NewFromConfig(cfg)
	if err != nil {
		log.Printf("Error: failed to create client, %v\n", err)
		return *client, err
	}

	return *client, nil

}
