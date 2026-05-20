package engine

import (
    "crypto/rand"
    "math/big"
    "sueca-online/internal/models"
)

// NewDeck cria um baralho padrão de 40 cartas de Sueca
func NewDeck() []models.Card {
    suits := []models.Suit{models.Ouros, models.Copas, models.Espadas, models.Paus}
    ranks := []models.Rank{
        models.Two, models.Three, models.Four, models.Five, models.Six, models.Seven,
        models.Jack, models.Queen, models.King, models.Ace,
    }

    // Pontuação clássica da Sueca de Portugal:
    // Ás=11, 7=10, Rei=4, Dama=3, Valete=2, Restantes=0
    // Total por naipe = 30. Total no baralho = 120.
    pointsMap := map[models.Rank]int{
        models.Ace:   11,
        models.Seven: 10,
        models.King:  4,
        models.Queen: 3,
        models.Jack:  2,
        models.Two:   0,
        models.Three: 0,
        models.Four:  0,
        models.Five:  0,
        models.Six:   0,
    }

    deck := make([]models.Card, 0, 40)
    for _, suit := range suits {
        for _, rank := range ranks {
            deck = append(deck, models.Card{
                Suit:   suit,
                Rank:   rank,
                Points: pointsMap[rank],
            })
        }
    }
    return deck
}

// Shuffle embaralha o slice usando Fisher-Yates com crypto/rand (CSPRNG)
func Shuffle(deck []models.Card) []models.Card {
    n := len(deck)
    for i := n - 1; i > 0; i-- {
        // Gera um número aleatório seguro entre 0 e i (inclusive)
        jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
        if err != nil {
            // Em caso de erro no crypto/rand (extremamente raro), fazemos panic
            // pois a segurança do jogo ficou comprometida (Defensive Programming)
            panic("failed to generate secure random number: " + err.Error())
        }
        j := int(jBig.Int64())

        // Swap (Algoritmo Fisher-Yates)
        deck[i], deck[j] = deck[j], deck[i]
    }
    return deck
}