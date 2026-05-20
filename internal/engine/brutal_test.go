package engine

import (
    "sueca-online/internal/models"
    "testing"
)

// TestZeroPointCardsHierarchy expõe o bug de cartas com 0 pontos (6 vs 5)
func TestZeroPointCardsHierarchy(t *testing.T) {
    room := &models.GameRoom{
        State:     models.StatePlaying,
        TrumpCard: models.Card{Suit: models.Paus}, // Trunfo irrelevante
        CurrentTrick: models.Trick{
            LeadSuit: models.Ouros,
            CardsPlayed: map[string]models.Card{
                "p1": {Suit: models.Ouros, Rank: models.Five, Points: 0}, // 5 de Ouros
                "p2": {Suit: models.Copas, Rank: models.Ace, Points: 11}, // Descarte
                "p3": {Suit: models.Ouros, Rank: models.Six, Points: 0},  // 6 de Ouros (Deve ganhar ao 5!)
                "p4": {Suit: models.Espadas, Rank: models.King, Points: 4}, // Descarte
            },
        },
        Players: [4]models.Player{
            {ID: "p1", Team: 1}, {ID: "p2", Team: 2},
            {ID: "p3", Team: 1}, {ID: "p4", Team: 2},
        },
    }

    winnerID, err := ResolveTrick(room)
    if err != nil {
        t.Fatalf("Erro inesperado: %v", err)
    }

    // O p3 jogou o 6, que é hierarquicamente superior ao 5 do p1
    if winnerID != "p3" {
        t.Fatalf("BUG DE HIERARQUIA: Esperava vencedor p3 (6 de Ouros), obteve %s", winnerID)
    }
}

// TestMultipleTrumpsPlayed garante que o trunfo mais alto ganha quando vários "pisam"
func TestMultipleTrumpsPlayed(t *testing.T) {
    room := &models.GameRoom{
        State:     models.StatePlaying,
        TrumpCard: models.Card{Suit: models.Paus},
        CurrentTrick: models.Trick{
            LeadSuit: models.Ouros,
            CardsPlayed: map[string]models.Card{
                "p1": {Suit: models.Ouros, Rank: models.Ace, Points: 11}, // Sai no naipe
                "p2": {Suit: models.Paus, Rank: models.Two, Points: 0},   // Pisa com 2 de Paus
                "p3": {Suit: models.Paus, Rank: models.Seven, Points: 10},// Pisa com 7 de Paus (Deve ganhar!)
                "p4": {Suit: models.Ouros, Rank: models.King, Points: 4}, // Acompanha naipe
            },
        },
        Players: [4]models.Player{
            {ID: "p1", Team: 1}, {ID: "p2", Team: 2},
            {ID: "p3", Team: 1}, {ID: "p4", Team: 2},
        },
    }

    winnerID, _ := ResolveTrick(room)

    if winnerID != "p3" {
        t.Fatalf("Esperava vencedor p3 (7 de Paus), obteve %s", winnerID)
    }
}

// TestResolveIncompleteTrick garante que o servidor crasha se tentar resolver vaza com 3 cartas
func TestResolveIncompleteTrick(t *testing.T) {
    room := &models.GameRoom{
        State:     models.StatePlaying,
        TrumpCard: models.Card{Suit: models.Paus},
        CurrentTrick: models.Trick{
            LeadSuit: models.Ouros,
            CardsPlayed: map[string]models.Card{
                "p1": {Suit: models.Ouros, Rank: models.Ace, Points: 11},
                "p2": {Suit: models.Ouros, Rank: models.King, Points: 4},
                "p3": {Suit: models.Ouros, Rank: models.Queen, Points: 3}, // Faltou o p4!
            },
        },
    }

    _, err := ResolveTrick(room)

    if err != ErrTrickNotComplete {
        t.Fatal("O motor deve rejeitar resolução de vazas incompletas (Defensive Programming)")
    }
}

// TestPlayCardWhenHandIsEmpty impede jogadas fantasma
func TestPlayCardWhenHandIsEmpty(t *testing.T) {
    room := setupTestRoom()
    room.Players[0].Hand = []models.Card{} // Jogador sem cartas

    cardToPlay := models.Card{Suit: models.Ouros, Rank: models.Ace, Points: 11}
    err := PlayCard(room, "p1", cardToPlay)

    if err != ErrCardNotInHand {
        t.Fatal("Um jogador com a mão vazia não conseguir jogar cartas fantasma")
    }
}

// TestTurnPassesToWinner garante que o vencedor da vaza anterior começa a próxima
func TestTurnPassesToWinner(t *testing.T) {
    room := &models.GameRoom{
        State:     models.StatePlaying,
        TrumpCard: models.Card{Suit: models.Paus},
        CurrentTrick: models.Trick{
            LeadSuit: models.Ouros,
            CardsPlayed: map[string]models.Card{
                "p1": {Suit: models.Ouros, Rank: models.Ace, Points: 11},
                "p2": {Suit: models.Ouros, Rank: models.King, Points: 4},
                "p3": {Suit: models.Ouros, Rank: models.Queen, Points: 3},
                "p4": {Suit: models.Ouros, Rank: models.Jack, Points: 2},
            },
        },
        Players: [4]models.Player{
            {ID: "p1", Team: 1, Hand: []models.Card{{Suit: models.Copas, Rank: models.Ace}}},
            {ID: "p2", Team: 2, Hand: []models.Card{{Suit: models.Copas, Rank: models.King}}},
            {ID: "p3", Team: 1, Hand: []models.Card{{Suit: models.Copas, Rank: models.Queen}}},
            {ID: "p4", Team: 2, Hand: []models.Card{{Suit: models.Copas, Rank: models.Jack}}},
        },
    }

    ResolveTrick(room) // p1 ganhou

    // Tentar jogar com o p2 (Não é a vez dele!)
    err := PlayCard(room, "p2", models.Card{Suit: models.Copas, Rank: models.King})
    if err != ErrNotYourTurn {
        t.Fatal("O vencedor da vaza deve ser o primeiro a jogar a próxima. p2 jogou fora de turno!")
    }

    // Tentar jogar com o p1 (Vencedor - Deve permitir)
    err = PlayCard(room, "p1", models.Card{Suit: models.Copas, Rank: models.Ace})
    if err != nil {
        t.Fatal("O vencedor p1 devia poder iniciar a próxima vaza")
    }
}