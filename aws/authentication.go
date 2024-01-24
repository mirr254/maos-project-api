package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateAwsSession (){
	
    _, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "ao"),
	})

	if err != nil {
		fmt.Printf("Ann error occured while creating sesssion %v", err.Error())
		return
	}
}


