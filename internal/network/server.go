package network

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

// Upgrader configura o processo de upgrade de HTTP para WebSocket
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // Por agora, permitimos todas as origens (CORS) para desenvolvimento.
    // Em produção (Fase 6), isto deve ser restrito ao teu domínio!
    CheckOrigin: func(r *http.Request) bool { return true },
}

// Server representa o servidor WebSocket
type Server struct {
    // Hub será adicionado no próximo passo
    clients map[*websocket.Conn]bool // Temporário: lista de clientes conectados
}

// NewServer cria uma nova instância do servidor
func NewServer() *Server {
    return &Server{
        clients: make(map[*websocket.Conn]bool),
    }
}

// HandleWebSocket lida com as novas ligações WebSocket
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Faz o upgrade da ligação HTTP para WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Erro no upgrade do WebSocket: %v", err)
        return
    }
    // Garantir que a ligação fecha no fim da função
    defer conn.Close()

    // Registar o cliente (temporário)
    s.clients[conn] = true
    log.Printf("Novo cliente conectado! Total: %d", len(s.clients))

    // Loop de leitura: O servidor fica à escuta de mensagens deste cliente
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            // Se houver erro (ex: desconectou), removemos o cliente e saímos do loop
            delete(s.clients, conn)
            log.Printf("Cliente desconectado. Total: %d", len(s.clients))
            break
        }

        // Por agora, fazemos apenas Echo (devolvemos a mesma mensagem)
        log.Printf("Recebido: %s", message)
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            log.Printf("Erro ao enviar mensagem: %v", err)
            break
        }
    }
}

// Start inicia o servidor HTTP na porta indicada
func (s *Server) Start(port string) {
    http.HandleFunc("/ws", s.HandleWebSocket)

    log.Printf("Servidor WebSocket a escuta na porta %s...", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatalf("Erro ao iniciar servidor: %v", err)
    }
}