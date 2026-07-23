package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	httplayer "github.com/artcodefun/detective-game/backend/internal/interfaces/http"
)

type Module struct {
	Port     string
	Adapters *Adapters
	Commands *Commands
	Queries  *Queries
	Server   *httplayer.Server
}

func NewModule() *Module {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	adapters := NewAdapters()
	commands := NewCommands(adapters)
	queries := NewQueries(adapters)

	handlers := &httplayer.Handlers{
		Scenario:      commands.Scenario,
		Interrogation: commands.Interrogation,
		Evaluation:    commands.Evaluation,
		Actions:       commands.Actions,
		Notebook:      commands.Notebook,

		Session:    queries.Session,
		Character:  queries.Character,
		Evidence:   queries.Evidence,
		Chronology: queries.Chronology,
		Chat:       queries.Chat,
	}

	addr := fmt.Sprintf(":%s", port)
	router := httplayer.NewRouter(handlers)
	server := httplayer.NewServer(router, addr)

	return &Module{
		Port:     port,
		Adapters: adapters,
		Commands: commands,
		Queries:  queries,
		Server:   server,
	}
}

func (m *Module) Run(ctx context.Context) error {
	log.Printf("Starting server on :%s", m.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := m.Server.Start(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server...")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.Server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server stopped cleanly")
	return nil
}
