package engine

import (
    "sueca-online/internal/models"
    "testing"
)

// TestNewDeckSizeAndSum valida a integridade do baralho
func TestNewDeckSizeAndSum(t *testing.T) {
    deck := NewDeck()

    // 1. Verificar se tem exatamente 40 cartas
    if len(deck) != 40 {
        t.Fatalf("Esperado 40 cartas no baralho, obtido %d", len(deck))
    }

    // 2. Verificar se a soma total de pontos é 120 (conforme documento)
    totalPoints := 0
    for _, card := range deck {
        totalPoints += card.Points
    }

    if totalPoints != 120 {
        t.Fatalf("Esperado 120 pontos totais no baralho, obtido %d", totalPoints)
    }
}

// TestShuffleRandomizesDeck valida que o shuffle altera a ordem e não perde cartas
func TestShuffleRandomizesDeck(t *testing.T) {
    originalDeck := NewDeck()
    
    // Fazer uma cópia para comparar
    copyDeck := make([]models.Card, len(originalDeck))
    copy(copyDeck, originalDeck)

    shuffledDeck := Shuffle(originalDeck)

    // 1. Verificar se o baralho baralhado tem o mesmo tamanho
    if len(shuffledDeck) != len(copyDeck) {
        t.Fatal("O baralho baralhado tem um tamanho diferente do original")
    }

    // 2. Verificar se a ordem mudou (estatisticamente impossível ser igual após Fisher-Yates)
    // A menos que tenhas uma sorte absurda com o crypto/rand!
    areEqual := true
    for i := range shuffledDeck {
        if shuffledDeck[i] != copyDeck[i] {
            areEqual = false
            break
        }
    }

    if areEqual {
        t.Fatal("O baralho permaneceu na mesma ordem após o shuffle. O algoritmo não está a funcionar?")
    }
}