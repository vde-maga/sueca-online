package engine

import (
    "errors"
    "sueca-online/internal/models"
)

var ErrNotEnoughPlayers = errors.New("need 4 players to deal")

// DealCards distribui as cartas e define o trunfo, transitando o estado do jogo
func DealCards(room *models.GameRoom) error {
    // 1. Verificar se estamos no Lobby
    if room.State != models.StateLobby {
        return ErrInvalidTransition
    }

    // 2. Validar se temos 4 jogadores conectados
    for _, p := range room.Players {
        if p.ID == "" || !p.IsConnected {
            return ErrNotEnoughPlayers
        }
    }

    // 3. Transitar para DEALING
    err := TransitionRoomState(room, models.StateDealing)
    if err != nil {
        return err
    }

    // 4. Gerar e baralhar o deck
    deck := NewDeck()
    deck = Shuffle(deck)

    // 5. Definir o Trunfo (A primeira carta do baralho define o naipe)
    // Em modo digital, a carta "física" vai para a mão de um jogador, mas o naipe fica registado.
    room.TrumpCard = deck[0]

    // 6. Distribuir 10 cartas a cada jogador
    for i := 0; i < 4; i++ {
        startIndex := i * 10
        room.Players[i].Hand = deck[startIndex : startIndex+10]
    }

    // 7. Transitar para PLAYING
    err = TransitionRoomState(room, models.StatePlaying)
    if err != nil {
        return err
    }

    // O primeiro jogador (índice 0) começa a jogar
    room.CurrentTurnIndex = 0
    room.CurrentTrick = models.Trick{
        CardsPlayed: make(map[string]models.Card),
    }

    return nil
}