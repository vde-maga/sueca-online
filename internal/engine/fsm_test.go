package engine

import (
    "sueca-online/internal/models"
    "testing"
)

// TestInvalidStateTransition verifica se o FSM bloqueia saltos proibidos
func TestInvalidStateTransition(t *testing.T) {
    // Cenário: Uma sala acabou de ser criada, está no estado LOBBY.
    room := models.GameRoom{
        State: models.StateLobby,
    }

    // Tentativa inválida: Tentar ir do LOBBY diretamente para PLAYING
    // (O correto seria LOBBY -> DEALING -> PLAYING)
    err := TransitionRoomState(&room, models.StatePlaying)

    // O teste passa se ocorrer um erro (porque a transição é inválida)
    if err == nil {
        t.Fatal("Esperava-se um erro ao tentar transitar de LOBBY para PLAYING, mas nenhum erro ocorreu")
    }

    // Verifica se o estado não foi alterado indevidamente
    if room.State != models.StateLobby {
        t.Fatalf("O estado da sala foi alterado para %s, mas deveria ter permanecido em LOBBY", room.State)
    }
}

// TestValidStateTransition verifica se o FSM permite transições corretas
func TestValidStateTransition(t *testing.T) {
    room := models.GameRoom{
        State: models.StateLobby,
    }

    // Tentativa válida: LOBBY -> DEALING
    err := TransitionRoomState(&room, models.StateDealing)

    if err != nil {
        t.Fatalf("Esperava-se que a transição de LOBBY para DEALING fosse bem-sucedida, mas obteve-se: %v", err)
    }

    if room.State != models.StateDealing {
        t.Fatalf("O estado da sala é %s, mas deveria ser DEALING", room.State)
    }
}