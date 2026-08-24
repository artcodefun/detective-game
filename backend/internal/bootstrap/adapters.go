package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/i18n"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/llm"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/readstore"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/repo"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Adapters struct {
	mongoClient *mongo.Client

	Users          ports.UserRepository
	Sessions       ports.SessionRepository
	Characters     ports.CharacterRepository
	Interrogations ports.InterrogationRepository
	Chat           ports.ChatMessageRepository
	Evidence       ports.EvidenceRepository
	Reports        ports.ActionReportRepository
	Chronology     ports.ChronologyRepository
	TxMgr          ports.TransactionManager

	ReadSessions ports.SessionReadRepository
	ReadChars    ports.CharacterReadRepository
	ReadEvidence ports.EvidenceReadRepository
	ReadReports  ports.ReportReadRepository
	ReadChron    ports.ChronologyReadRepository
	ReadChat     ports.ChatMessageReadRepository

	LLM        ports.LlmService
	Translator ports.Translator
}

func NewAdapters(cfg Config) *Adapters {
	client, db := newDatabase(cfg.MongoURI, cfg.MongoDatabase)

	llmService := llm.NewOpenRouterClient(cfg.OpenRouterKey, cfg.OpenRouterModel)
	translator := i18n.NewTranslator()
	log.Printf("using OpenRouter, model=%s", cfg.OpenRouterModel)

	return &Adapters{
		mongoClient: client,

		Users:          repo.NewUserRepo(db),
		Sessions:       repo.NewSessionRepo(db),
		Characters:     repo.NewCharacterRepo(db),
		Interrogations: repo.NewInterrogationRepo(db),
		Chat:           repo.NewChatRepo(db),
		Evidence:       repo.NewEvidenceRepo(db),
		Reports:        repo.NewReportRepo(db),
		Chronology:     repo.NewChronologyRepo(db),
		TxMgr:          repo.NewMongoTxManager(client),

		ReadSessions: readstore.NewSessionReadRepo(db),
		ReadChars:    readstore.NewCharacterReadRepo(db),
		ReadEvidence: readstore.NewEvidenceReadRepo(db),
		ReadReports:  readstore.NewReportReadRepo(db),
		ReadChron:    readstore.NewChronologyReadRepo(db),
		ReadChat:     readstore.NewChatReadRepo(db),

		LLM:        llmService,
		Translator: translator,
	}
}

func (a *Adapters) Setup(ctx context.Context) error {
	users, ok := a.Users.(*repo.UserRepo)
	if !ok {
		return fmt.Errorf("users adapter does not support setup")
	}
	if err := users.EnsureIndexes(ctx); err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}
	return nil
}

func (a *Adapters) Shutdown(ctx context.Context) error {
	return a.mongoClient.Disconnect(ctx)
}

func newDatabase(uri, dbName string) (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}

	return client, client.Database(dbName)
}
