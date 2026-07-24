package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/llm"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/readstore"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/repo"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Adapters struct {
	Users          ports.UserRepository
	Sessions       ports.SessionRepository
	Characters     ports.CharacterRepository
	Interrogations ports.InterrogationRepository
	Chat           ports.ChatMessageRepository
	Evidence       ports.EvidenceRepository
	Reports        ports.ActionReportRepository
	Chronology     ports.ChronologyRepository

	ReadSessions ports.SessionReadRepository
	ReadChars    ports.CharacterReadRepository
	ReadEvidence ports.EvidenceReadRepository
	ReadReports  ports.ReportReadRepository
	ReadChron    ports.ChronologyReadRepository
	ReadChat     ports.ChatMessageReadRepository

	LLM        ports.LlmService
	Prototypes ports.CharacterPrototypeRepository
}

func NewAdapters(cfg Config) *Adapters {
	db := newDatabase(cfg.MongoURI, cfg.MongoDatabase)

	return &Adapters{
		Users:          repo.NewUserRepo(db),
		Sessions:       repo.NewSessionRepo(db),
		Characters:     repo.NewCharacterRepo(db),
		Interrogations: repo.NewInterrogationRepo(db),
		Chat:           repo.NewChatRepo(db),
		Evidence:       repo.NewEvidenceRepo(db),
		Reports:        repo.NewReportRepo(db),
		Chronology:     repo.NewChronologyRepo(db),

		ReadSessions: readstore.NewSessionReadRepo(db),
		ReadChars:    readstore.NewCharacterReadRepo(db),
		ReadEvidence: readstore.NewEvidenceReadRepo(db),
		ReadReports:  readstore.NewReportReadRepo(db),
		ReadChron:    readstore.NewChronologyReadRepo(db),
		ReadChat:     readstore.NewChatReadRepo(db),

		LLM:        llm.NewMockLlmService(),
		Prototypes: repo.NewPrototypeRepo(db),
	}
}

func newDatabase(uri, dbName string) *mongo.Database {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}

	return client.Database(dbName)
}
