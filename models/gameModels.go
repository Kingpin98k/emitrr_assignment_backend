package game_models

type Game struct {
	EMAIL              string `json:"email"`
	DECK               string `json:"deck"`
	POINTER            int    `json:"pointer"`
	CURRENT_HIGH_SCORE int    `json:"currentHighScore"`
	DEFFUSECARDS       int    `json:"deffuseCards"`
	BOMBCARDS          int    `json:"bombCards"`
}

type GameResponse struct {
	EMAIL              string `json:"email"`
	POINTER            int    `json:"pointer"`
	CURRENT_HIGH_SCORE int    `json:"currentHighScore"`
	TOP_CARD           string `json:"topCard"`
	DEFFUSECARDS       int    `json:"deffuseCards"`
	BOMBCARDS          int    `json:"bombCards"`
}