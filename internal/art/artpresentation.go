package art

type ArtPresentation struct {
	Piece              ArtPiece              `json:"art_piece" dynamodbav:"art_piece"`
	Description        string                `json:"description" dynamodbav:"description"`
	DescriptionClipURL string                `json:"description_clip_url" dynamodbav:"description_clip_url"`
	Tasks              []ArtPresentationTask `json:"tasks" dynamodbav:"tasks"`
}

type ArtPresentationTask struct {
	Task        string `json:"task_text" dynamodbav:"task_text"`
	TaskClipURL string `json:"tasl_clip_url" dynamodbav:"tasl_clip_url"`
}
