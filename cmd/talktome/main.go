package main

import (
	"fmt"
	"os"

	textgeneration "talktome.com/internal/textgenerator"
)

func main() {
	openAIToken := mustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)

	artistName := "Caspar David Friedrich"
	artName := "Der Wanderer über dem Wolkenmeer"

	fmt.Printf("Generate description for %s's \"%s\"\n", artistName, artName)

	// description, err := textGen.GenerateArtDescription(artistName, artName)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(description)

	fmt.Printf("Generate tasks for %s's \"%s\"\n", artistName, artName)

	tasks, err := textGen.GenerateTasks(artistName, artName)
	if err != nil {
		panic(err)
	}
	for _, task := range tasks {
		fmt.Println("-----> " + task)
	}
}

func mustReadEnvVar(name string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	} else {
		panic(fmt.Sprintf("FATAL: env var %s does not exists", name))
	}
}
