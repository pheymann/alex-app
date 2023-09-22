package integrationtest_util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqualInterface(t *testing.T, expected interface{}, actual string) {
	if _, ok := expected.([]interface{}); ok {
		var actualUnmarshaled []interface{}
		err := json.Unmarshal([]byte(actual), &actualUnmarshaled)
		assert.NoError(t, err)

		expectedConverted := YamlToJson(expected)

		if !compare(t, expectedConverted, actualUnmarshaled) {
			t.Fail()
		}
	} else {
		var actualUnmarshaled interface{}
		err := json.Unmarshal([]byte(actual), &actualUnmarshaled)
		assert.NoError(t, err)

		expectedConverted := YamlToJson(expected)

		if !compare(t, expectedConverted, actualUnmarshaled) {
			t.Fail()
		}
	}
}

func compare(t *testing.T, expected, actual interface{}) bool {
	switch expectedTyped := expected.(type) {
	case map[string]interface{}:
		actualTyped, ok := actual.(map[string]interface{})
		if !ok {
			t.Errorf("expected %v to be a map[string]interface{}", actual)
			return false
		}

		for actualKey, actualValue := range actualTyped {
			if actualValue == nil {
				continue
			}

			expectedValue, ok := expectedTyped[actualKey]
			if !ok {
				continue
			}

			if expectedValue == nil {
				if str, ok := actualValue.(string); ok && str == "" {
					continue
				}
			}

			if !compare(t, expectedValue, actualValue) {
				return false
			}
		}
		return true

	case []interface{}:
		actualTyped, ok := actual.([]interface{})
		if !ok {
			t.Errorf("expected %v to be a []interface{}", actual)
			return false
		}

		if len(expectedTyped) != len(actualTyped) {
			t.Errorf("expected %v to have length %d, got %d", actual, len(expectedTyped), len(actualTyped))
			return false
		}

		for i, actualValue := range actualTyped {
			if !compare(t, expectedTyped[i], actualValue) {
				return false
			}
		}

		return true

	default:
		return assert.Equal(t, expected, actual)
	}
}
