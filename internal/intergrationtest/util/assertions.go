package integrationtest_util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualJson[B any](t *testing.T, expected B, actual string, description string) {
	var actualUnmarshaled B
	err := json.Unmarshal([]byte(actual), &actualUnmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, expected, actualUnmarshaled, description)
}
