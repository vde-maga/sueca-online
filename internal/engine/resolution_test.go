package engine

import (
    "sueca-online/internal/models"
    "testing"
)

// setupTrickRoom cria uma sala onde já foram jogadas 4 cartas, prontas para resolver
func setupTrickRoom() *models.GameRoom {
    room := &models.GameRoom{
        State:     models.StatePlaying,
        TrumpCard: models.Card{Suit: models.Paus}, // Trunfo é Paus
        CurrentTrick: models.Trick{
            LeadSuit: models.Ouros, // Naipe de saída é Ouros
            CardsPlayed: map[string]models.Card{
                "p1": {Suit: models.Ouros, Rank: models.Ace, Points: 11},   // 11 pts
                "p2": {Suit: models.Copas, Rank: models.Ace, Points: 11},   // Descarte (não é naipe nem trunfo)
                "p3": {Suit: models.Ouros, Rank: models.Seven, Points: 10}, // 10 pts
                "p4": {Suit: models.Paus, Rank: models.Two, Points: 0},     // Trunfo! (Vale 0 pts mas ganha)
            },
        },
        Players: [4]models.Player{
            {ID: "p1", Team: 1, Hand: []models.Card{}}, // Mão vazia para simular fim de jogo depois
            {ID: "p2", Team: 2, Hand: []models.Card{}},
            {ID: "p3", Team: 1, Hand: []models.Card{}},
            {ID: "p4", Team: 2, Hand: []models.Card{}},
        },
    }
    return room
}

// TestTrickWonByTrump verifica se um trunfo baixo ganha ao naipe de saída alto
func TestTrickWonByTrump(t *testing.T) {
    room := setupTrickRoom()

    winnerID, err := ResolveTrick(room)
    if err != nil {
        t.Fatalf("Esperava sucesso, obteve erro: %v", err)
    }

    // O vencedor deve ser o p4, que jogou o 2 de Paus (Trunfo)
    if winnerID != "p4" {
        t.Fatalf("Esperava vencedor p4, obteve %s", winnerID)
    }

    // Pontos na mesa: 11 (Ás Ouros) + 11 (Ás Copas) + 10 (7 Ouros) + 0 (2 Paus) = 32
    if room.Team2Score != 32 {
        t.Fatalf("Esperava Team2Score 32, obteve %d", room.Team2Score)
    }

    // A vez deve passar para o vencedor (p4 está no índice 3)
    if room.CurrentTurnIndex != 3 {
        t.Fatalf("Esperava CurrentTurnIndex 3, obteve %d", room.CurrentTurnIndex)
    }
}

// TestTrickWonByHighestLeadSuit verifica quem ganha se não houver trunfos
func TestTrickWonByHighestLeadSuit(t *testing.T) {
    room := setupTrickRoom()

    // Vamos alterar a jogada do p4 para um descarte (Copas) em vez de Trunfo
    room.CurrentTrick.CardsPlayed["p4"] = models.Card{Suit: models.Copas, Rank: models.King, Points: 4}

    winnerID, err := ResolveTrick(room)
    if err != nil {
        t.Fatalf("Esperava sucesso, obteve erro: %v", err)
    }

    // Sem trunfos, ganha a carta mais alta do naipe de saída (Ouros).
    // p1 jogou Ás (11 pts), p3 jogou 7 (10 pts). Vence o p1.
    if winnerID != "p1" {
        t.Fatalf("Esperava vencedor p1, obteve %s", winnerID)
    }

    // Pontos na mesa: 11 (Ás Ouros) + 11 (Ás Copas) + 10 (7 Ouros) + 4 (Rei Copas) = 36
    if room.Team1Score != 36 {
        t.Fatalf("Esperava Team1Score 36, obteve %d", room.Team1Score)
    }
}