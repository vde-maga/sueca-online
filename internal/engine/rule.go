package engine

import (
    "errors"
    "sueca-online/internal/models"
)

var (
    ErrGameNotPlaying = errors.New("game is not in PLAYING state")
    ErrNotYourTurn    = errors.New("it is not your turn")
    ErrCardNotInHand  = errors.New("card is not in your hand")
)

// PlayCard tenta executar a ação de jogar uma carta no estado atual da sala
func PlayCard(room *models.GameRoom, playerID string, card models.Card) error {
    // 1. State Check
    if room.State != models.StatePlaying {
        return ErrGameNotPlaying
    }

    // 2. Turn Check
    currentPlayer := room.Players[room.CurrentTurnIndex]
    if currentPlayer.ID != playerID {
        return ErrNotYourTurn
    }

    // 3. Ownership Check
    cardIndex := -1
    for i, c := range currentPlayer.Hand {
        if c.Suit == card.Suit && c.Rank == card.Rank {
            cardIndex = i
            break
        }
    }

    if cardIndex == -1 {
        return ErrCardNotInHand
    }

    // Se chegámos aqui, a jogada é válida. Vamos mutar o estado.
    
    // Remover a carta da mão do jogador (Slice trick em Go)
    currentPlayer.Hand = append(currentPlayer.Hand[:cardIndex], currentPlayer.Hand[cardIndex+1:]...)
    room.Players[room.CurrentTurnIndex] = currentPlayer // Atualizar o array no room

    // Inicializar o mapa se for a primeira carta da vaza
    if room.CurrentTrick.CardsPlayed == nil {
        room.CurrentTrick.CardsPlayed = make(map[string]models.Card)
    }

    // Definir o LeadSuit se for a primeira carta jogada na vaza
    if len(room.CurrentTrick.CardsPlayed) == 0 {
        room.CurrentTrick.LeadSuit = card.Suit
    }

    // Adicionar a carta à vaza
    room.CurrentTrick.CardsPlayed[playerID] = card

    return nil
}