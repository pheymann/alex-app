package integrationtest_cdc

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"talktome.com/internal/conversation"
	integrationtest_util "talktome.com/internal/intergrationtest/util"
	"talktome.com/internal/shared"
	"talktome.com/internal/user"
)

type ClientDefinedContract[R any] struct {
	Request       CDCRequest           `yaml:"Request"`
	ResponseCases []CDCResponseCase[R] `yaml:"ResponseCases"`
}

type CDCRequest struct {
	Uri          string        `yaml:"uri"`
	PathPameters []CDCKeyValue `yaml:"pathParameters"`
	Method       string        `yaml:"method"`
	Headers      []CDCKeyValue `yaml:"headers"`
	Body         string        `yaml:"body"`
}

func (request CDCRequest) GetMapHeaders() map[string]string {
	headers := make(map[string]string)
	for _, header := range request.Headers {
		headers[header.Name] = header.Value
	}

	return headers
}

type CDCResponseCase[R any] struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	StatusCode int         `yaml:"statusCode"`
	Body       *R          `yaml:"body"`
	ErrorBody  string      `yaml:"errorBody"`
	Database   CDCDatabase `yaml:"database"`
}

type CDCDatabase struct {
	Users         []user.User                 `yaml:"users"`
	Conversations []conversation.Conversation `yaml:"conversations"`
}

type CDCKeyValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func ReadCDC[R any](contractPath string) (*ClientDefinedContract[R], error) {
	contractFile, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, err
	}

	var contract ClientDefinedContract[R]
	if err := yaml.Unmarshal([]byte(contractFile), &contract); err != nil {
		return nil, err
	}

	return &contract, nil
}

func MustReadCDC[R any](contractPath string) ClientDefinedContract[R] {
	contract, err := ReadCDC[R](contractPath)
	if err != nil {
		panic(err)
	}

	return *contract
}

func MustLoadCDC[R any](cdcPath string) ClientDefinedContract[R] {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return MustReadCDC[R](rootPath + "/../../../cdc" + cdcPath)
}

func RunContracts[R any](
	t *testing.T,
	cdcPath string,
	eval func(
		t *testing.T,
		responseCase CDCResponseCase[R],
		event events.APIGatewayProxyRequest,
		users map[string]*user.User,
		conversations map[string]*conversation.Conversation,
	) (events.APIGatewayProxyResponse, error),
) {
	contracts := MustLoadCDC[R](cdcPath)

	testUserRequestCtx := shared.NewAwsTestRequestContext("0")

	pathParameters := map[string]string{}
	for _, param := range contracts.Request.PathPameters {
		pathParameters[param.Name] = param.Value
	}

	event := events.APIGatewayProxyRequest{
		Resource:       contracts.Request.Uri,
		Path:           contracts.Request.Uri,
		PathParameters: pathParameters,
		HTTPMethod:     contracts.Request.Method,
		Headers:        contracts.Request.GetMapHeaders(),
		Body:           contracts.Request.Body,
		RequestContext: testUserRequestCtx,
	}

	for _, responseCase := range contracts.ResponseCases {
		t.Logf("Running contract test case: %s", responseCase.Name)
		users := make(map[string]*user.User)
		for _, user := range responseCase.Database.Users {
			userRef := user
			users[userRef.ID] = &userRef
		}

		conversations := make(map[string]*conversation.Conversation)
		for _, conv := range responseCase.Database.Conversations {
			convRef := conv
			conversations[conv.ID] = &convRef
		}

		response, err := eval(t, responseCase, event, users, conversations)
		assert.NoError(t, err)

		assert.Equal(t, responseCase.StatusCode, response.StatusCode)
		if responseCase.Type == "success" {
			integrationtest_util.AssertEqualJson(t, *responseCase.Body, response.Body, "body does not match")
		} else {
			assert.Equal(t, responseCase.ErrorBody, response.Body, "error body does not match")
		}
	}
}
