package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	httplayer "github.com/artcodefun/detective-game/backend/internal/interfaces/http"
)

type Module struct {
	Config   Config
	Adapters *Adapters
	Commands *Commands
	Queries  *Queries
	Server   *httplayer.Server
}

func NewModule(cfg Config) *Module {
	adapters := NewAdapters(cfg)
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
		Translator: adapters.Translator,
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	router := httplayer.NewRouter(handlers)
	server := httplayer.NewServer(router, addr)

	return &Module{
		Config:   cfg,
		Adapters: adapters,
		Commands: commands,
		Queries:  queries,
		Server:   server,
	}
}

func (m *Module) Run(ctx context.Context) error {
	log.Printf("Starting server on :%s", m.Config.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := m.Server.Start(); err != nil {
			errCh <- err
		}
	}()

	var runErr error
	select {
	case <-ctx.Done():
	case err := <-errCh:
		runErr = err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Shutting down...")

	if err := m.Server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if err := m.Adapters.Shutdown(shutdownCtx); err != nil {
		log.Printf("mongo shutdown error: %v", err)
	}

	log.Println("Stopped")
	return runErr
}
