package integrationtest_util

import (
	"fmt"
)

func YamlToJson(m interface{}) interface{} {
	if array, ok := m.([]interface{}); ok {
		res := []interface{}{}
		for _, v := range array {
			res = append(res, YamlToJson(v))
		}
		return res
	}

	return yamlToJsonMap(m.(map[interface{}]interface{}))
}

func yamlToJsonMap(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = YamlToJson(v2)

		case []interface{}:
			res[fmt.Sprint(k)] = YamlToJson(v2)

		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}
