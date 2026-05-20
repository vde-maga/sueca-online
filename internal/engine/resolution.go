package engine

import (
    "errors"
    "sueca-online/internal/models"
)

var ErrTrickNotComplete = errors.New("trick does not have 4 cards yet")

// ResolveTrick determina o vencedor da vaza atual, atualiza a pontuação e prepara a próxima
func ResolveTrick(room *models.GameRoom) (string, error) {
    if len(room.CurrentTrick.CardsPlayed) != 4 {
        return "", ErrTrickNotComplete
    }

    leadSuit := room.CurrentTrick.LeadSuit
    trumpSuit := room.TrumpCard.Suit

    var winningPlayerID string
    var winningCard models.Card
    trickPoints := 0
    firstCard := true

    // Iterar sobre as cartas jogadas para encontrar o vencedor
    for playerID, card := range room.CurrentTrick.CardsPlayed {
        trickPoints += card.Points

        if firstCard {
            winningPlayerID = playerID
            winningCard = card
            firstCard = false
            continue
        }

        isNewCardTrump := card.Suit == trumpSuit
        isCurrentWinnerTrump := winningCard.Suit == trumpSuit

        if isNewCardTrump {
            if !isCurrentWinnerTrump {
                // A nova carta é trunfo e a atual não. A nova ganha.
                winningPlayerID = playerID
                winningCard = card
            } else {
                // Ambas são trunfos. A de maior pontuação ganha.
                if card.Points > winningCard.Points {
                    winningPlayerID = playerID
                    winningCard = card
                }
            }
        } else if card.Suit == leadSuit {
            // A nova carta não é trunfo, mas é do naipe de saída.
            if !isCurrentWinnerTrump && card.Points > winningCard.Points {
                // A atual também não é trunfo e está no naipe de saída. A maior ganha.
                winningPlayerID = playerID
                winningCard = card
            }
        }
        // Se não for trunfo nem do naipe de saída, é descarte e nunca ganha.
    }

    // Encontrar a equipa do vencedor para adicionar os pontos
    winnerIndex := findPlayerIndexByID(room, winningPlayerID)
    winningTeam := room.Players[winnerIndex].Team

    if winningTeam == 1 {
        room.Team1Score += trickPoints
    } else {
        room.Team2Score += trickPoints
    }

    // O vencedor torna-se o próximo a jogar
    room.CurrentTurnIndex = winnerIndex

    // Limpar a mesa (Preparar para a próxima vaza)
    room.CurrentTrick = models.Trick{
        CardsPlayed: make(map[string]models.Card),
    }

    // Transição de Estado: Verificar se o jogo acabou (mãos vazias) ou continuar
    if len(room.Players[0].Hand) == 0 {
        TransitionRoomState(room, models.StateGameOver)
    } else {
        TransitionRoomState(room, models.StatePlaying)
    }

    return winningPlayerID, nil
}

// findPlayerIndexByID é um helper para encontrar o índice no array de jogadores
func findPlayerIndexByID(room *models.GameRoom, playerID string) int {
    for i, p := range room.Players {
        if p.ID == playerID {
            return i
        }
    }
    return -1 // Não deve acontecer com IDs válidos
}