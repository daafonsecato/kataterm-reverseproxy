package services

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSService struct {
	EC2 *ec2.EC2
}

func NewAWSService() *AWSService {

	sess := session.Must(session.NewSession())
	return &AWSService{
		EC2: ec2.New(sess),
	}
}

func (s *AWSService) CreateInstance(imageID string, instanceType string, subnetID string) (*ec2.Instance, error) {
	// Load the AWS SDK configuration

	s = NewAWSService()
	svc := s.EC2

	userData := `#!/bin/bash
		cd /home/ubuntu/kateterm
		echo "DB_HOST=10.0.0.175" >> backend/.env
		docker compose up -d
		# More commands
	`
	userDataBase64 := base64.StdEncoding.EncodeToString([]byte(userData))

	// Define the instance launch configuration
	runInstancesInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(imageID),
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		UserData:     aws.String(userDataBase64),
		SubnetId:     aws.String(subnetID),
	}

	// Launch the instance
	result, err := svc.RunInstances(runInstancesInput)
	if err != nil {
		fmt.Println("Error launching instance:", err)
		return nil, err
	}

	// Return the instance from the reservation
	return result.Instances[0], nil
}

func (s *AWSService) TerminateInstance(awsInstanceID string) error {
	// Initialize AWS session
	s = NewAWSService()
	svc := s.EC2
	// Terminate the instance
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(awsInstanceID)},
	}

	_, err := svc.TerminateInstancesWithContext(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error terminating instance: %v", err)
	}

	return nil
}
