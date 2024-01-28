package services

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSService struct {
	EC2 *ec2.EC2
}

func NewAWSService(sess *session.Session) *AWSService {
	return &AWSService{
		EC2: ec2.New(sess),
	}
}

func (s *AWSService) CreateInstance(imageID string, instanceType string, subnetID string) (*ec2.Instance, error) {
	// Load the AWS SDK configuration
	sess := session.Must(session.NewSession())
	s = NewAWSService(sess)
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
	result, err := svc.RunInstances(context.TODO(), runInstancesInput)
	if err != nil {
		fmt.Println("Error launching instance:", err)
		return
	}
}

func terminateInstanceBySessionID(db *sql.DB, sessionID string) error {
	awsInstanceID, err := getAWSInstanceID(db, sessionID)
	if err != nil {
		return fmt.Errorf("error getting AWS instance ID: %v", err)
	}

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		// Specify AWS Config if needed
	})
	if err != nil {
		return fmt.Errorf("error creating AWS session: %v", err)
	}

	// Create new EC2 client
	svc := ec2.New(sess)

	// Terminate the instance
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(awsInstanceID)},
	}

	_, err = svc.TerminateInstancesWithContext(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error terminating instance: %v", err)
	}

	return nil
}

func (s *AWSService) TerminateInstance(instanceID string) error {
	sess := session.Must(session.NewSession())
	s = NewAWSService(sess)
	svc := s.EC2
	terminateInstanceBySessionID(db, sessionID)
}
