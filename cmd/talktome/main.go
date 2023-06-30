package main

import (
	"fmt"
	"os"

	"talktome.com/internal/cmd/talktome"
)

func main() {
	openAIToken := mustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	resembleToken := mustReadEnvVar("TALKTOME_RESEMBLE_TOKEN")
	resembleProjectUUID := mustReadEnvVar("TALKTOME_RESEMBLE_PROJECT_UUID")
	serviceDomain := mustReadEnvVar("TALKTOME_SERVICE_DOMAIN")
	resembleCallBackURL := fmt.Sprintf("https://%s/callback/clip", serviceDomain)

	talktome := talktome.NewTalkToMe(openAIToken, resembleToken, resembleProjectUUID, resembleCallBackURL)

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

	rawAudioClip, err := talktome.DownloadSpeechClip("https://app.resemble.ai/rails/active_storage/blobs/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBCTUFNQ3cwPSIsImV4cCI6bnVsbCwicHVyIjoiYmxvYl9pZCJ9fQ==--b9d5ed35f6bf3f3deb0f5e44e9ad77b19ec8203b/CLI+Test-238771b5.wav")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("test.wav")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write(rawAudioClip)
	if err != nil {
		panic(err)
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
