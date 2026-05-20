package engine

import (
    "sueca-online/internal/models"
    "testing"
)

// TestDealCardsSuccess verifica a distribuição correta de cartas e definição do trunfo
func TestDealCardsSuccess(t *testing.T) {
    room := &models.GameRoom{
        State: models.StateLobby,
        Players: [4]models.Player{
            {ID: "p1", IsConnected: true},
            {ID: "p2", IsConnected: true},
            {ID: "p3", IsConnected: true},
            {ID: "p4", IsConnected: true},
        },
    }

    err := DealCards(room)
    if err != nil {
        t.Fatalf("Esperava sucesso, obteve erro: %v", err)
    }

    // Verificar se o estado mudou para PLAYING
    if room.State != models.StatePlaying {
        t.Fatalf("Esperava estado PLAYING, obteve %s", room.State)
    }

    // Verificar se o Trunfo foi definido (a primeira carta do baralho antes de distribuir)
    if room.TrumpCard.Suit == "" {
        t.Fatal("O Trunfo não foi definido na sala")
    }

    // Verificar se cada jogador tem exatamente 10 cartas
    for _, p := range room.Players {
        if len(p.Hand) != 10 {
            t.Fatalf("Jogador %s deveria ter 10 cartas, tem %d", p.ID, len(p.Hand))
        }
    }

    // Verificar se as cartas do jogador incluem o trunfo (a carta está na mão de alguém)
    // A soma das mãos deve totalizar 40 cartas (o baralho todo)
    totalCards := 0
    for _, p := range room.Players {
        totalCards += len(p.Hand)
    }
    if totalCards != 40 {
        t.Fatalf("Esperava 40 cartas distribuidas, obteve %d", totalCards)
    }
}

// TestDealCardsNotEnoughPlayers garante que não se começa sem 4 jogadores
func TestDealCardsNotEnoughPlayers(t *testing.T) {
    room := &models.GameRoom{
        State: models.StateLobby,
        Players: [4]models.Player{
            {ID: "p1", IsConnected: true},
            {ID: "p2", IsConnected: true},
            // p3 e p4 faltam (Zero value)
        },
    }

    err := DealCards(room)
    if err == nil {
        t.Fatal("Não deveria ser possível distribuir cartas sem 4 jogadores")
    }
}