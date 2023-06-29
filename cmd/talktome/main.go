package main

import (
	"fmt"
	"os"

	"talktome.com/internal/cmd/talktome"
)

func main() {
	openAIToken := mustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	talktome := talktome.NewTalkToMe(openAIToken)

	artistName := "Caspar David Friedrich"
	artName := "Der Wanderer über dem Wolkenmeer"

	presentation, err := talktome.GenerateArtPresentation(artistName, artName)
	if err != nil {
		panic(err)
	}

	fmt.Println("= " + presentation.Description)
	for _, task := range presentation.Tasks {
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
