package main

import (
	"fmt"
	"os"

	textgeneration "talktome.com/internal/textgenerator"
)

func main() {
	openAIToken := mustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)

	content, err := textGen.GenerateArtDescription("Caspar David Friedrich", "Der Wanderer über dem Wolkenmeer")
	if err != nil {
		panic(err)
	}

	fmt.Println(content)
}

func mustReadEnvVar(name string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	} else {
		panic(fmt.Sprintf("FATAL: env var %s does not exists", name))
	}
}
