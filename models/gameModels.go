package game_models

type Game struct {
	EMAIL              string `json:"email"`
	DECK               string `json:"deck"`
	POINTER            int    `json:"pointer"`
	CURRENT_HIGH_SCORE int    `json:"currentHighScore"`
}