package art

import "fmt"

type ArtPiece struct {
	ArtistName string `json:"artist_name" dynamodbav:"artist_name"`
	Name       string `json:"art_piece_name" dynamodbav:"art_piece_name"`
}

func (piece ArtPiece) String() string {
	return fmt.Sprintf(`%s's "%s"`, piece.ArtistName, piece.Name)
}
