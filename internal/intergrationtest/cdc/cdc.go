package integrationtest_cdc

import (
	"encoding/json"
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

type ClientDefinedContract struct {
	Name               string                  `yaml:"name"`
	View               string                  `yaml:"view"`
	AuthorizationToken string                  `yaml:"authorizationToken"`
	Database           CDCDatabase             `yaml:"database"`
	CallChain          []CDCRequestAndResponse `yaml:"callChain"`
}

type CDCDatabase struct {
	Users         []user.User                 `yaml:"users"`
	Conversations []conversation.Conversation `yaml:"conversations"`
}

type CDCRequestAndResponse struct {
	Request  CDCRequest  `yaml:"request"`
	Response CDCResponse `yaml:"response"`
}

type CDCRequest struct {
	Uri          string        `yaml:"uri"`
	PathPameters []CDCKeyValue `yaml:"pathParameters"`
	Method       string        `yaml:"method"`
	Headers      []CDCKeyValue `yaml:"headers"`
	Body         *interface{}  `yaml:"body"`
}

type CDCKeyValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (request CDCRequest) GetMapHeaders() map[string]string {
	headers := make(map[string]string)
	for _, header := range request.Headers {
		headers[header.Name] = header.Value
	}

	return headers
}

type CDCResponse struct {
	Type       string       `yaml:"type"`
	StatusCode int          `yaml:"status"`
	Body       *interface{} `yaml:"body"`
	ErrorBody  string       `yaml:"errorBody"`
}

func ReadCDC[R any](contractPath string) (*ClientDefinedContract, error) {
	contractFile, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, err
	}

	var contract ClientDefinedContract
	if err := yaml.Unmarshal([]byte(contractFile), &contract); err != nil {
		return nil, err
	}

	return &contract, nil
}

func MustReadCDC[R any](contractPath string) ClientDefinedContract {
	contract, err := ReadCDC[R](contractPath)
	if err != nil {
		panic(err)
	}

	return *contract
}

func MustLoadCDC[R any](cdcPath string) ClientDefinedContract {
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
		event events.APIGatewayProxyRequest,
		users map[string]*user.User,
		conversations map[string]*conversation.Conversation,
	) (events.APIGatewayProxyResponse, error),
) {
	contract := MustLoadCDC[R](cdcPath)

	t.Logf("Running contract test case: %s", contract.Name)

	testUserRequestCtx := shared.NewAwsTestRequestContext("0")

	users := make(map[string]*user.User)
	for _, user := range contract.Database.Users {
		userRef := user
		users[userRef.ID] = &userRef
	}

	conversations := make(map[string]*conversation.Conversation)
	for _, conv := range contract.Database.Conversations {
		convRef := conv
		conversations[conv.ID] = &convRef
	}

	for _, requestAndResponse := range contract.CallChain {
		pathParameters := map[string]string{}
		for _, param := range requestAndResponse.Request.PathPameters {
			pathParameters[param.Name] = param.Value
		}

		var body = ""
		if requestAndResponse.Request.Body != nil {
			bodyBytes, err := json.Marshal(integrationtest_util.YamlToJson(*requestAndResponse.Request.Body))
			if err != nil {
				panic(err)
			}

			body = string(bodyBytes)
		}

		event := events.APIGatewayProxyRequest{
			Resource:       requestAndResponse.Request.Uri,
			Path:           requestAndResponse.Request.Uri,
			PathParameters: pathParameters,
			HTTPMethod:     requestAndResponse.Request.Method,
			Headers:        requestAndResponse.Request.GetMapHeaders(),
			Body:           body,
			RequestContext: testUserRequestCtx,
		}

		response, err := eval(t, event, users, conversations)
		assert.NoError(t, err)

		assert.Equal(t, requestAndResponse.Response.StatusCode, response.StatusCode)
		if requestAndResponse.Response.Type == "success" && response.Body != "" {
			integrationtest_util.AssertEqualInterface(t, *requestAndResponse.Response.Body, response.Body)
		} else {
			assert.Equal(t, requestAndResponse.Response.ErrorBody, response.Body, "error body does not match")
		}
	}
}
