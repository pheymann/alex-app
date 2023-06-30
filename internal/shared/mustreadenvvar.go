package shared

import (
	"fmt"
	"os"
)

func MustReadEnvVar(name string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	} else {
		panic(fmt.Sprintf("FATAL: env var %s does not exists", name))
	}
}
