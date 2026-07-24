package repo

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PrototypeRepo struct {
	coll   *mongo.Collection
	seeded []domain.CharacterPrototype
}

func NewPrototypeRepo(db *mongo.Database) *PrototypeRepo {
	return &PrototypeRepo{
		coll: db.Collection("character_prototypes"),
		seeded: []domain.CharacterPrototype{
			{ID: 1, Name: "Иван Петров", Age: 55, Profession: "дворецкий",
				ImagePath:   "assets/characters/ivan_petrov.png",
				Personality: "Консервативный, преданный семье, скрытный. Говорит медленно, с расстановкой. Предпочитает отмалчиваться, но если задеть за живое — срывается.",
				AudioToneID: "tone_male_deep"},
			{ID: 2, Name: "Елена Соколова", Age: 42, Profession: "домохозяйка",
				ImagePath:   "assets/characters/elena_sokolova.png",
				Personality: "Эмоциональная, вспыльчивая, но ранимая. Говорит быстро, часто перебивает. Хочет казаться безразличной, но на деле очень переживает.",
				AudioToneID: "tone_female_high"},
			{ID: 3, Name: "Майкл Браун", Age: 48, Profession: "деловой партнёр",
				ImagePath:   "assets/characters/michael_brown.png",
				Personality: "Харизматичный, уверенный в себе, умело манипулирует. Говорит спокойно, с лёгкой усмешкой. Всегда контролирует эмоции.",
				AudioToneID: "tone_male_mid"},
			{ID: 4, Name: "Анна Коваль", Age: 29, Profession: "горничная",
				ImagePath:   "assets/characters/anna_koval.png",
				Personality: "Застенчивая, тревожная, боится потерять работу. Говорит тихо, запинается. Старается быть незаметной, но глаза выдают страх.",
				AudioToneID: "tone_female_soft"},
			{ID: 5, Name: "Дмитрий Орлов", Age: 61, Profession: "адвокат",
				ImagePath:   "assets/characters/dmitry_orlov.png",
				Personality: "Циничный, расчётливый, за словом в карман не лезет. Говорит чётко, рублеными фразами. Привык контролировать ситуацию.",
				AudioToneID: "tone_male_raspy"},
		},
	}
}

func (r *PrototypeRepo) seed(ctx context.Context) error {
	count, err := r.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	docs := make([]any, len(r.seeded))
	for i, p := range r.seeded {
		docs[i] = p
	}
	_, err = r.coll.InsertMany(ctx, docs)
	return err
}

func (r *PrototypeRepo) GetAll(ctx context.Context) ([]domain.CharacterPrototype, error) {
	if err := r.seed(ctx); err != nil {
		return nil, fmt.Errorf("seed prototypes: %w", err)
	}
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find prototypes: %w", err)
	}
	defer cursor.Close(ctx)

	var items []domain.CharacterPrototype
	for cursor.Next(ctx) {
		var p domain.CharacterPrototype
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("decode prototype: %w", err)
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *PrototypeRepo) GetRandom(ctx context.Context, count int) ([]domain.CharacterPrototype, error) {
	all, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if count >= len(all) {
		return all, nil
	}
	shuffled := make([]domain.CharacterPrototype, len(all))
	copy(shuffled, all)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:count], nil
}

func (r *PrototypeRepo) ByID(ctx context.Context, id int) (*domain.CharacterPrototype, error) {
	if err := r.seed(ctx); err != nil {
		return nil, fmt.Errorf("seed prototypes: %w", err)
	}
	var p domain.CharacterPrototype
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("find prototype: %w", err)
	}
	return &p, nil
}
