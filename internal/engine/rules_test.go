package engine

import (
    "sueca-online/internal/models"
    "testing"
)

// setupTestRoom cria uma sala genérica no estado PLAYING para testarmos as jogadas
func setupTestRoom() *models.GameRoom {
    room := &models.GameRoom{
        State:            models.StatePlaying,
        CurrentTurnIndex: 0, // O Player 1 (índice 0) começa
        CurrentTrick: models.Trick{
            CardsPlayed: make(map[string]models.Card),
        },
        Players: [4]models.Player{
            {ID: "p1", Hand: []models.Card{{Suit: models.Ouros, Rank: models.Ace, Points: 11}}},
            {ID: "p2", Hand: []models.Card{{Suit: models.Copas, Rank: models.Ace, Points: 11}}},
            {ID: "p3", Hand: []models.Card{{Suit: models.Espadas, Rank: models.Ace, Points: 11}}},
            {ID: "p4", Hand: []models.Card{{Suit: models.Paus, Rank: models.Ace, Points: 11}}},
        },
    }
    return room
}

// TestPlayOutOfTurn verifica se o jogador não pode jogar quando não é a sua vez
func TestPlayOutOfTurn(t *testing.T) {
    room := setupTestRoom()
    
    // O turno é do "p1" (índice 0), mas o "p2" tenta jogar
    cardToPlay := models.Card{Suit: models.Copas, Rank: models.Ace, Points: 11}
    err := PlayCard(room, "p2", cardToPlay)

    if err != ErrNotYourTurn {
        t.Fatalf("Esperava ErrNotYourTurn, obtive: %v", err)
    }
}

// TestPlayCardNotInHand verifica se o jogador não pode jogar uma carta que não tem na mão
func TestPlayCardNotInHand(t *testing.T) {
    room := setupTestRoom()
    
    // O "p1" tem Ouros Ás, mas tenta jogar Paus Ás (carta que não tem)
    cardToPlay := models.Card{Suit: models.Paus, Rank: models.Ace, Points: 11}
    err := PlayCard(room, "p1", cardToPlay)

    if err != ErrCardNotInHand {
        t.Fatalf("Esperava ErrCardNotInHand, obtive: %v", err)
    }
}

// TestValidPlay verifica se uma jogada válida é aceite e altera o estado corretamente
func TestValidPlay(t *testing.T) {
    room := setupTestRoom()
    
    cardToPlay := models.Card{Suit: models.Ouros, Rank: models.Ace, Points: 11}
    err := PlayCard(room, "p1", cardToPlay)

    if err != nil {
        t.Fatalf("Jogada deveria ser válida, mas obteve-se erro: %v", err)
    }

    // Verificar se a carta foi removida da mão do jogador
    if len(room.Players[0].Hand) != 0 {
        t.Fatal("A carta não foi removida da mão do jogador após jogada válida")
    }

    // Verificar se a carta foi adicionada à vaza atual
    if len(room.CurrentTrick.CardsPlayed) != 1 {
        t.Fatal("A carta não foi adicionada à vaza atual")
    }

    // Verificar se o naipe de saída (LeadSuit) foi definido
    if room.CurrentTrick.LeadSuit != models.Ouros {
        t.Fatalf("LeadSuit esperado Ouros, obtido %s", room.CurrentTrick.LeadSuit)
    }
}