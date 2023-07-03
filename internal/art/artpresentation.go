package art

import "fmt"

type ArtPresentation struct {
	ID                  string                `json:"id" dynamodbav:"id"`
	Piece               ArtPiece              `json:"art_piece" dynamodbav:"art_piece"`
	Description         string                `json:"description" dynamodbav:"description"`
	DescriptionClipUUID string                `json:"description_clip_uuid" dynamodbav:"description_clip_uuid"`
	Tasks               []ArtPresentationTask `json:"tasks" dynamodbav:"tasks"`
}

func CreatePresentationID(piece ArtPiece) string {
	return fmt.Sprintf("%s::%s", piece.ArtistName, piece.Name)
}

type ArtPresentationTask struct {
	Task         string `json:"task_text" dynamodbav:"task_text"`
	TaskClipUUID string `json:"task_clip_uuid" dynamodbav:"task_clip_uuid"`
}
