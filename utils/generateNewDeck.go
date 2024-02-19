package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateNewDeck() string {
	availableCards := []string{"cat", "defuse", "exploding_kitten", "shuffle"}
	deckSize := 5
	deck := make([]string, deckSize)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < deckSize; i++ {
		randomIndex := rand.Intn(len(availableCards))
		deck[i] = availableCards[randomIndex]
	}

	return strings.Join(deck, ",")
}
