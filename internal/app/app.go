package app

import server "github.com/davidgaspardev/usermes-backend/internal/infra/http/fiber"

func Run() {
	// Create a new server
	srv := server.NewServer(3000)

	// Start the server
	srv.Start()
}
