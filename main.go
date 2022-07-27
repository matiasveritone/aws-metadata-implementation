package main

import (
	"context"
)
import "log"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := imds.NewFromConfig(cfg)

	log.Printf("Lalala")
	log.Printf("Client %v", *client)

	localip, err := client.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: "local-ipv4",
	})
	if err != nil {
		log.Printf("Unable to retrieve the private IP address from the EC2 instance: %s\n", err)
		return
	}

	log.Printf("local-ip: %v\n", localip)

}
