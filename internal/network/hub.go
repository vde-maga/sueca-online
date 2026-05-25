package network

import (
    "encoding/json"
    "log"
    "sueca-online/internal/engine"
    "sueca-online/internal/models"
    "sync"

    "github.com/google/uuid"
    "github.com/gorilla/websocket"
)

// Client representa um jogador conectado ao servidor
type Client struct {
    Conn     *websocket.Conn
    PlayerID string
    RoomCode string
}

// Hub gere todas as salas de jogo e clientes conectados
type Hub struct {
    mu      sync.Mutex               // Protege o acesso concorrente aos maps
    Rooms   map[string]*models.GameRoom
    Clients map[*websocket.Conn]*Client
}

func NewHub() *Hub {
    return &Hub{
        Rooms:   make(map[string]*models.GameRoom),
        Clients: make(map[*websocket.Conn]*Client),
    }
}

// HandleDisconnect remove um jogador que fechou a ligação
func (h *Hub) HandleDisconnect(conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()

    client, exists := h.Clients[conn]
    if !exists {
        return // O jogador não estava numa sala
    }

    // Remover do mapa de clientes
    delete(h.Clients, conn)

    // Atualizar a sala (marcar como desconectado)
    room, exists := h.Rooms[client.RoomCode]
    if !exists {
        return
    }

    for i, p := range room.Players {
        if p.ID == client.PlayerID {
            room.Players[i].IsConnected = false
            log.Printf("Jogador %s desconectou da sala %s", p.Nickname, client.RoomCode)
            break
        }
    }

    // Avisar a sala que alguém saiu
    h.broadcastRoomUpdate(room)
}

// JoinRoom lida com a lógica de um cliente a tentar entrar numa sala
func (h *Hub) JoinRoom(conn *websocket.Conn, payload JoinRoomPayload) {
    h.mu.Lock()
    defer h.mu.Unlock()

    // 1. Validar Input
    if payload.RoomCode == "" || payload.Nickname == "" {
        sendErrorToConn(conn, "Código de sala e nickname são obrigatórios")
        return
    }

    // 2. Verificar se o jogador já está numa sala
    if client, exists := h.Clients[conn]; exists && client.RoomCode != "" {
        sendErrorToConn(conn, "Já estás numa sala!")
        return
    }

    // 3. Obter ou criar a sala
    room, exists := h.Rooms[payload.RoomCode]
    if !exists {
        room = &models.GameRoom{
            ID:    payload.RoomCode,
            State: models.StateLobby,
        }
        h.Rooms[payload.RoomCode] = room
        log.Printf("Sala %s criada", payload.RoomCode)
    }

    // 4. Verificar se a sala está cheia
    connectedPlayers := 0
    for _, p := range room.Players {
        if p.IsConnected {
            connectedPlayers++
        }
    }

    if connectedPlayers >= 4 {
        sendErrorToConn(conn, "A sala está cheia (Máximo 4 jogadores)")
        return
    }

    // 5. Adicionar jogador à sala
    playerID := uuid.New().String()
    player := models.Player{
        ID:          playerID,
        Nickname:    payload.Nickname,
        IsConnected: true,
        Team:        connectedPlayers%2 + 1,
    }

    for i := 0; i < 4; i++ {
        if room.Players[i].ID == "" {
            room.Players[i] = player
            break
        }
    }

    // 6. Registar o cliente no Hub
    h.Clients[conn] = &Client{
        Conn:     conn,
        PlayerID: playerID,
        RoomCode: payload.RoomCode,
    }

    log.Printf("Jogador %s entrou na sala %s. Total na sala: %d", payload.Nickname, payload.RoomCode, connectedPlayers+1)

    // 7. Notificar toda a sala
    h.broadcastRoomUpdate(room)

    // 8. Verificar se o jogo pode começar
    if connectedPlayers+1 == 4 {
        log.Printf("Sala %s está cheia! A distribuir cartas...", payload.RoomCode)
        err := engine.DealCards(room)
        if err != nil {
            log.Printf("Erro ao distribuir cartas na sala %s: %v", payload.RoomCode, err)
            return
        }
        h.broadcastGameStart(room)
    }
}

// broadcastRoomUpdate envia o estado do lobby para todos os jogadores da sala
func (h *Hub) broadcastRoomUpdate(room *models.GameRoom) {
    type LobbyPlayer struct {
        Nickname string `json:"nickname"`
        Team     int    `json:"team"`
    }
    type LobbyPayload struct {
        RoomCode string        `json:"roomCode"`
        Players  []LobbyPlayer `json:"players"`
    }

    players := []LobbyPlayer{}
    for _, p := range room.Players {
        if p.IsConnected {
            players = append(players, LobbyPlayer{Nickname: p.Nickname, Team: p.Team})
        }
    }

    msg := OutgoingMessage{
        Type: TypeRoomUpdate,
        Payload: LobbyPayload{
            RoomCode: room.ID,
            Players:  players,
        },
    }

    h.broadcastToRoom(room.ID, msg)
}

// broadcastGameStart envia o estado inicial do jogo
// broadcastGameStart envia o estado inicial do jogo MASCARADO (Data Masking)
func (h *Hub) broadcastGameStart(room *models.GameRoom) {
    // NOTA: Esta função é chamada de dentro do JoinRoom, que já tem o Lock ativo.
    // Não podemos usar h.mu.Lock() aqui senão causamos Deadlock!

    // Vamos iterar sobre os clientes conectados a esta sala
    for conn, client := range h.Clients {
        if client.RoomCode == room.ID {
            // Procurar a mão deste jogador específico no estado do Motor
            var playerHand []models.Card
            for _, p := range room.Players {
                if p.ID == client.PlayerID {
                    playerHand = p.Hand
                    break
                }
            }

            // Construir o payload seguro (Só a mão deste jogador, trunfo e vez)
            payload := GameStartPayload{
                Hand:        playerHand,
                TrumpCard:   room.TrumpCard,
                CurrentTurn: room.Players[room.CurrentTurnIndex].ID,
                PlayerID:    client.PlayerID,
            }

            msg := OutgoingMessage{
                Type:    TypeGameStart,
                Payload: payload,
            }

            data, err := json.Marshal(msg)
            if err != nil {
                log.Printf("Erro ao criar mensagem de GameStart: %v", err)
                continue
            }

            // Enviar APENAS para este cliente
            conn.WriteMessage(websocket.TextMessage, data)
        }
    }
}

// broadcastToRoom envia uma mensagem para todos os WebSockets de uma sala
func (h *Hub) broadcastToRoom(roomCode string, msg OutgoingMessage) {
    data, err := json.Marshal(msg)
    if err != nil {
        log.Printf("Erro ao fazer marshal da mensagem: %v", err)
        return
    }

    for conn, client := range h.Clients {
        if client.RoomCode == roomCode {
            conn.WriteMessage(websocket.TextMessage, data)
        }
    }
}

// sendErrorToConn é um helper para enviar erros
func sendErrorToConn(conn *websocket.Conn, message string) {
    errMsg := OutgoingMessage{
        Type:    TypeError,
        Payload: ErrorPayload{Message: message},
    }
    data, _ := json.Marshal(errMsg)
    conn.WriteMessage(websocket.TextMessage, data)
}