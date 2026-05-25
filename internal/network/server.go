package network

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
    hub *Hub
}

func NewServer() *Server {
    return &Server{
        hub: NewHub(),
    }
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Erro no upgrade do WebSocket: %v", err)
        return
    }
    
    // Garantir que o servidor limpa o jogador quando a ligação fechar
    defer func() {
        s.hub.HandleDisconnect(conn)
        conn.Close()
    }()

    log.Printf("Novo cliente conectado!")

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            // Se houver erro (ex: desconectou), saímos do loop e o defer trata do resto
            log.Printf("Cliente desconectado ou erro de leitura.")
            break
        }

        var msg IncomingMessage
        err = json.Unmarshal(message, &msg)
        if err != nil {
            log.Printf("Payload inválido recebido: %v", err)
            sendErrorToConn(conn, "Formato de mensagem inválido. Deve ser JSON.")
            continue
        }

        switch msg.Action {
        case ActionJoinRoom:
            var payload JoinRoomPayload
            if err := json.Unmarshal(msg.Payload, &payload); err != nil {
                sendErrorToConn(conn, "Payload de JOIN_ROOM inválido")
                continue
            }
            s.hub.JoinRoom(conn, payload)

        case ActionPlayCard:
            var payload PlayCardPayload
            if err := json.Unmarshal(msg.Payload, &payload); err != nil {
                sendErrorToConn(conn, "Payload de PLAY_CARD inválido")
                continue
            }
            log.Printf("Jogador quer jogar %s de %s", payload.Rank, payload.Suit)

        default:
            sendErrorToConn(conn, "Ação desconhecida")
        }
    }
}

func (s *Server) Start(port string) {
    http.HandleFunc("/ws", s.HandleWebSocket)
    log.Printf("Servidor WebSocket a escuta na porta %s...", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatalf("Erro ao iniciar servidor: %v", err)
    }
}