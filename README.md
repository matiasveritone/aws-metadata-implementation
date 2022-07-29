# aws-metadata-implementation
## Purpose

Implementation of EC2 metadata library (AWS virtual machine instances) in Golang, 
in order to check the following items:
- If we are in an AWS virtual environment
- If the user that's logged has permissions to run an API.

### Requirements
- AWS EC2 Instance ([docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EC2_GetStarted.html)).
- AWS IAM Role with necesary permissions ([docs](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-service.html#roles-creatingrole-service-console)).
  - Amazon EC2 access
  - Amazon {API} access
- AWS IAM Role attached to the intance where the engine is running ([docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#attach-iam-role))

### Parameters

`isAWSAccount` ---> Checks if AWS Account is allowed to run API inside the EC2 machine. Value `true` for account checking, `false` for custom credentials.

`region` ---> Defines the region of the endpoint, `us-east-1` is the one we always use.


### Documentation 
- [AWS SDK for Go v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2#section-documentation)
- [EC2 metadata Go Package](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/ec2metadata#pkg-overview)
- [EC2 metadata AWS docs](https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/ec2-imds/)
- [EC2 metadata AWS categories](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html#dynamic-data-categories)
