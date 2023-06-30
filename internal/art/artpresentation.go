package art

type ArtPresentation struct {
	ID                  string                `json:"id" dynamodbav:"id"`
	Piece               ArtPiece              `json:"art_piece" dynamodbav:"art_piece"`
	Description         string                `json:"description" dynamodbav:"description"`
	DescriptionClipUUID string                `json:"description_clip_uuid" dynamodbav:"description_clip_uuid"`
	DescriptionClipURL  string                `json:"description_clip_url" dynamodbav:"description_clip_url"`
	Tasks               []ArtPresentationTask `json:"tasks" dynamodbav:"tasks"`
}

type ArtPresentationTask struct {
	Task         string `json:"task_text" dynamodbav:"task_text"`
	TaskClipUUID string `json:"task_clip_uuid" dynamodbav:"task_clip_uuid"`
	TaskClipURL  string `json:"tasl_clip_url" dynamodbav:"tasl_clip_url"`
}
