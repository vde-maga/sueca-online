package engine

import (
    "errors"
    "sueca-online/internal/models"
)

// ErrInvalidTransition é retornado quando uma transição de estado não é permitida
var ErrInvalidTransition = errors.New("invalid state transition")

// TransitionRoomState tenta mudar o estado da GameRoom baseado nas regras do FSM
func TransitionRoomState(room *models.GameRoom, desiredState models.GameState) error {
    current := room.State
    valid := false

    switch current {
    case models.StateLobby:
        if desiredState == models.StateDealing {
            valid = true
        }
    case models.StateDealing:
        if desiredState == models.StatePlaying {
            valid = true
        }
    case models.StatePlaying:
        if desiredState == models.StateTrickResolution {
            valid = true
        }
    case models.StateTrickResolution:
        // Pode voltar a Playing (próxima vaza) ou ir para GameOver (fim do jogo)
        if desiredState == models.StatePlaying || desiredState == models.StateGameOver {
            valid = true
        }
    case models.StateGameOver:
        if desiredState == models.StateLobby {
            valid = true
        }
    }

    if !valid {
        return ErrInvalidTransition
    }

    room.State = desiredState
    return nil
}