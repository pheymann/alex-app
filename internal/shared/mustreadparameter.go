package shared

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func MustReadParameter(name string, ssmClient *ssm.SSM) string {
	param, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	return *param.Parameter.Value
}
