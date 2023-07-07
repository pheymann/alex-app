package talktome

import (
	"talktome.com/internal/art"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
)

type TalkToMe struct {
	textGen    textgeneration.TextGenerator
	speechGen  speechgeneration.SpeechGenerator
	artStorage art.StorageService
}

func NewTalkToMe(textGen textgeneration.TextGenerator, speechGen speechgeneration.SpeechGenerator, storage art.StorageService) TalkToMe {
	return TalkToMe{
		textGen:    textGen,
		speechGen:  speechGen,
		artStorage: storage,
	}
}
