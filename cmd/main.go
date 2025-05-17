package main

import (
	"log"

	"github.com/thespecialone1/aicare/config"
	"github.com/thespecialone1/aicare/internal/api"
)

func main() {
	cfg := config.Load()
	db := config.ConnectDB(cfg.DBUrl)
	defer db.Close()

	// 1. Build Low level pieces
	gemClient, err := api.NewGeminiClient() // uses ENV inside
	if err != nil {
		log.Fatalf("creating Gemini client: %v", err)
	}
	qRepo := api.NewQuestionRepo(db)
	mRepo := api.NewMessageRepo(db)
	qaSvc := api.NewQAService(gemClient, qRepo, mRepo)

	// 2. Build and run server
	server, err := api.NewServer(":8080", qaSvc)
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
	log.Fatal(server.Run())
}
