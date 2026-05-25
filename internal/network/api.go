package network

import "encoding/json"

// ActionType define os tipos de ações que o cliente pode enviar
type ActionType string

const (
    ActionJoinRoom ActionType = "JOIN_ROOM"
    ActionPlayCard ActionType = "PLAY_CARD"
)

// IncomingMessage é o wrapper para todas as mensagens vindas do Cliente
type IncomingMessage struct {
    Action  ActionType          `json:"action"`
    Payload json.RawMessage     `json:"payload"` // O payload bruto, seremos nós a fazer o parse depois
}

// JoinRoomPayload representa os dados necessários para entrar numa sala
type JoinRoomPayload struct {
    RoomCode string `json:"roomCode"`
    Nickname string `json:"nickname"`
    Password string `json:"password"`
}

// PlayCardPayload representa os dados necessários para jogar uma carta
type PlayCardPayload struct {
    Suit string `json:"suit"`
    Rank string `json:"rank"`
}

// --- Respostas do Servidor ---

// OutgoingMessageType define os tipos de eventos que o servidor envia
type OutgoingMessageType string

const (
    TypeError       OutgoingMessageType = "ERROR"
    TypeRoomUpdate  OutgoingMessageType = "ROOM_UPDATE"
    TypeGameStart   OutgoingMessageType = "GAME_START"
    TypeCardPlayed  OutgoingMessageType = "CARD_PLAYED"
    TypeTrickEnd    OutgoingMessageType = "TRICK_END"
    TypeGameOver    OutgoingMessageType = "GAME_OVER"
)

// OutgoingMessage é o wrapper para todas as mensagens enviadas para o Cliente
type OutgoingMessage struct {
    Type    OutgoingMessageType `json:"type"`
    Payload interface{}         `json:"payload"` // Interface aceita qualquer struct
}

// ErrorPayload é a estrutura para enviar erros para o cliente
type ErrorPayload struct {
    Message string `json:"message"`
}