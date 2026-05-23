package main

import (
    "sueca-online/internal/network"
)

func main() {
    server := network.NewServer()
    
    // Arrancar o servidor na porta 8080
    server.Start("8080")
}