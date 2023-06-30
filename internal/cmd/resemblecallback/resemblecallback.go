package resemblecallback

import (
	"fmt"

	"talktome.com/internal/art"
	"talktome.com/internal/speechgeneration"
)

type ResembleCallBack struct {
	speechGen  *speechgeneration.SpeechGenerator
	artStorage *art.StorageCtx
}

func NewResembleCallBack(speechGen *speechgeneration.SpeechGenerator, storage *art.StorageCtx) ResembleCallBack {
	return ResembleCallBack{
		speechGen:  speechGen,
		artStorage: storage,
	}
}

func (callback ResembleCallBack) StoreSpeechClip(clipUUID string, clipURL string) error {
	fmt.Printf("[DEBUG] Store speech clip %s\n", clipUUID)

	clipAudio, err := callback.downloadSpeechClip(clipURL)
	if err != nil {
		return err
	}

	return callback.artStorage.StoreClip(clipUUID, clipAudio)
}

func (callback ResembleCallBack) downloadSpeechClip(clipURL string) ([]byte, error) {
	fmt.Printf("[DEBUG] Download speech clip %s\n", clipURL)
	return callback.speechGen.DownloadSpeechClip(clipURL)
}
