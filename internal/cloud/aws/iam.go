package aws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

type PolicyDocument struct {
	Version   string            `json:"Version"`
	Statement []PolicyStatement `json:Statement`
}

type PolicyStatement struct {
	Effect    string            `json:"Effect"`
	Principal map[string]string `json:"Principal"`
	Action    []string          `json:"Action"`
}

func (a CloudProvider) retrieveRole() string {
	svc := iam.New(a.session)

	output, err := svc.GetRole(&iam.GetRoleInput{
		RoleName: aws.String("eksClusterRole"),
	})

	if err != nil {
		fmt.Println("[🐶] Error describing Role: ", err)
		return ""
	}

	fmt.Printf("[🐶] Found role %s with Id: %s\n", *output.Role.RoleName, *output.Role.RoleId)
	return *output.Role.Arn

}

func (a CloudProvider) createRole() string {
	svc := iam.New(a.session)

	trustPolicy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:    "Allow",
			Principal: map[string]string{"Service": "eks.amazonaws.com"},
			Action:    []string{"sts:AssumeRole"},
		}},
	}

	policyBytes, err := json.Marshal(trustPolicy)
	if err != nil {
		log.Fatal("[🐶] Couldn't create trust policy for 'eks.amazonaws.com': ", err)
	}

	roleOutput, err := svc.CreateRole(&iam.CreateRoleInput{
		RoleName:                 aws.String("eksClusterRole"),
		AssumeRolePolicyDocument: aws.String(string(policyBytes)),
	})

	if err != nil {
		log.Fatal("[🐶] Error creating Role: ", err)
	}

	return *roleOutput.Role.Arn
}

func (a CloudProvider) attachRolePolicy() {
	svc := iam.New(a.session)

	fmt.Println("[🐶] Attaching Role Policy to eksClusterRole...")

	_, err := svc.AttachRolePolicy(&iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
		RoleName:  aws.String("eksClusterRole"),
	})

	if err != nil {
		log.Fatal("[🐶] Error attaching Role Policy: ", err)
	}

	fmt.Println("[🐶] Role Policy Successfully Attached")
}
