package art

import (
	"fmt"

	"talktome.com/internal/textgeneration"
)

type ArtConversation struct {
	ID                        string                      `json:"id" dynamodbav:"id"`
	Piece                     ArtPiece                    `json:"art_piece" dynamodbav:"art_piece"`
	ConversationStart         textgeneration.Conversation `json:"conversation_start" dynamodbav:"conversation_start"`
	ConversationStartClipUUID string                      `json:"conversation_start_clip_uuid" dynamodbav:"conversation_start_clip_uuid"`
}

func CreateArtConversationID(piece ArtPiece) string {
	return fmt.Sprintf("%s::%s", piece.ArtistName, piece.Name)
}
