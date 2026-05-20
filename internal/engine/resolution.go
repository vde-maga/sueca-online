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

        // Usar a nova função de comparação robusta em vez de comparar pontos!
        if isCardStronger(card, winningCard, trumpSuit, leadSuit) {
            winningPlayerID = playerID
            winningCard = card
        }
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

// rankHierarchy define a força absoluta de cada carta na Sueca
var rankHierarchy = map[models.Rank]int{
    models.Ace:   10,
    models.Seven: 9,
    models.King:  8,
    models.Queen: 7,
    models.Jack:  6,
    models.Six:   5,
    models.Five:  4,
    models.Four:  3,
    models.Three: 2,
    models.Two:   1,
}

// isCardStronger verifica se a carta A é mais forte que a carta B, dado o contexto do trunfo e naipe de saída
func isCardStronger(cardA, cardB models.Card, trumpSuit, leadSuit models.Suit) bool {
    aIsTrump := cardA.Suit == trumpSuit
    bIsTrump := cardB.Suit == trumpSuit

    // Se A é trunfo e B não é, A ganha
    if aIsTrump && !bIsTrump {
        return true
    }
    // Se B é trunfo e A não é, B ganha
    if !aIsTrump && bIsTrump {
        return false
    }
    // Se ambos são trunfos, ganha a hierarquia maior
    if aIsTrump && bIsTrump {
        return rankHierarchy[cardA.Rank] > rankHierarchy[cardB.Rank]
    }

    // Se nenhum é trunfo, verificamos o naipe de saída
    aIsLead := cardA.Suit == leadSuit
    bIsLead := cardB.Suit == leadSuit

    // Se A é do naipe e B não, A ganha
    if aIsLead && !bIsLead {
        return true
    }
    if !aIsLead && bIsLead {
        return false
    }
    // Se ambos são do naipe de saída, ganha a hierarquia maior
    if aIsLead && bIsLead {
        return rankHierarchy[cardA.Rank] > rankHierarchy[cardB.Rank]
    }

    // Se chegámos aqui, A é um descarte e B também. O primeiro a jogar (B) ganha por padrão.
    // Logo, A não é mais forte.
    return false
}