package models

// Suit representa o naipe da carta
type Suit string

const (
    Ouros   Suit = "Ouros"
    Copas   Suit = "Copas"
    Espadas Suit = "Espadas"
    Paus    Suit = "Paus"
)

// Rank representa o valor/face da carta
type Rank string

const (
    Two   Rank = "2"
    Three Rank = "3"
    Four  Rank = "4"
    Five  Rank = "5"
    Six   Rank = "6"
    Seven Rank = "7" // Necessário para fechar as 40 cartas
    Jack  Rank = "Valete"
    Queen Rank = "Rainha"
    King  Rank = "Rei"
    Ace   Rank = "Ás"
)

// Card representa uma carta de jogar
type Card struct {
    Suit   Suit `json:"suit"`
    Rank   Rank `json:"rank"`
    Points int  `json:"points"`
}

// GameState define os estados possíveis do FSM
type GameState string

const (
    StateLobby           GameState = "LOBBY"
    StateDealing         GameState = "DEALING"
    StatePlaying         GameState = "PLAYING"
    StateTrickResolution GameState = "TRICK_RESOLUTION"
    StateGameOver        GameState = "GAME_OVER"
)

// GameRoom representa a sala do jogo (O Estado Principal)
type GameRoom struct {
    ID               string    `json:"id"`
    Password         string    `json:"-"`
    Players          [4]Player `json:"players"`
    Deck             []Card    `json:"-"`
    TrumpCard        Card      `json:"trump_card"`
    CurrentTrick     Trick     `json:"current_trick"`
    Team1Score       int       `json:"team1_score"`
    Team2Score       int       `json:"team2_score"`
    CurrentTurnIndex int       `json:"current_turn_index"`
    State            GameState `json:"state"`
}

// Player representa um jogador
type Player struct {
    ID          string `json:"id"`
    Nickname    string `json:"nickname"`
    Hand        []Card `json:"hand"`
    Team        int    `json:"team"`
    IsConnected bool   `json:"is_connected"`
}

// Trick representa a vaza atual
type Trick struct {
    CardsPlayed map[string]Card `json:"cards_played"`
    LeadSuit    Suit            `json:"lead_suit"`
}