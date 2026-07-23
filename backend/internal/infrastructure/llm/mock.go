package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

var characterPool = []domain.CharacterPrototype{
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
}

var groundTruth = domain.Crime{
	Type:        domain.CrimeTypeMurder,
	Victim:      "Роберт Ланг",
	Motive:      "Ланг раскрыл хищение средств и собирался обратиться в полицию",
	Method:      "отравление цианидом (подмешан в виски)",
	TimeOfCrime: "22:15",
}

var scenarioTimeline = domain.Timeline{Entries: []domain.TimelineEntry{
	{Time: "19:00", Event: "Ужин в особняке. Все подозреваемые присутствуют."},
	{Time: "20:30", Event: "Ланг и Браун уходят в кабинет для разговора."},
	{Time: "21:00", Event: "Браун покидает кабинет. Ланг остается один.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000103"); return &id }()},
	{Time: "21:15", Event: "Иван Петров заходит в кабинет с виски.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000101"); return &id }()},
	{Time: "21:45", Event: "Орлов идёт в кабинет обсудить документы, но дверь заперта.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000105"); return &id }()},
	{Time: "22:00", Event: "Елена Соколова слышит шум из кабинета.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000102"); return &id }()},
	{Time: "22:15", Event: "Анна Коваль обнаруживает тело.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000104"); return &id }()},
}}

var scenarioEvidence = []domain.Evidence{
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Бокал с остатками виски", Description: "Найден на столе в кабинете",
		IconAsset:           "assets/evidence/glass.png",
		DetailedDescription: "Хрустальный бокал с остатками виски. Обнаружены следы цианида. Отпечатки пальцев: Роберт Ланг и Иван Петров.",
		Type:                domain.EvidenceTypePhysical},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Пузырёк с цианидом", Description: "Найден в мусорной корзине в кабинете",
		IconAsset:           "assets/evidence/vial.png",
		DetailedDescription: "Стеклянный флакон из-под ингалятора с остатками цианида. Флакон без отпечатков — стёрты. Производитель — компания, поставляющая лекарства в аптечную сеть Брауна.",
		Type:                domain.EvidenceTypePhysical},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Финансовые документы", Description: "Разбросаны на столе",
		IconAsset:           "assets/evidence/documents.png",
		DetailedDescription: "Бухгалтерские отчёты, указывающие на недостачу £500,000. Документы были скрыты, но кто-то их нашёл и разложил на столе. На полях пометки почерком Ланга: \"мошенничество\", \"Браун\", \"полиция\".",
		Type:                domain.EvidenceTypeDocument},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "Запертый ящик стола", Description: "Ключ торчит снаружи",
		IconAsset:           "assets/evidence/desk.png",
		DetailedDescription: "Верхний ящик стола заперт, ключ оставлен в замке. Внутри — конверт с компроматом на Орлова. Видимо, Ланг кого-то шантажировал — или готовился к этому.",
		Type:                domain.EvidenceTypePhysical},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005"), Name: "Записка с угрозой", Description: "Найдена под дверью кабинета",
		IconAsset:           "assets/evidence/note.png",
		DetailedDescription: "Анонимная записка: «Не лезь не в своё дело, а то пожалеешь». Бумага дешёвая, текст напечатан на принтере. Анализ чернил ничего не дал.",
		Type:                domain.EvidenceTypeDocument},
}

func makeCharacters() []domain.Character {
	chars := []domain.Character{
		{
			PrototypeID: 1, Name: "Иван Петров", Age: 55, Profession: "дворецкий",
			ImagePath:   "assets/characters/ivan_petrov.png",
			Personality: "Консервативный, преданный семье, скрытный. Говорит медленно, с расстановкой. Предпочитает отмалчиваться, но если задеть за живое — срывается.",
			AudioToneID: "tone_male_deep",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я принёс виски в кабинет в 21:15, Ланг был жив и раздражён", "В 21:00 я видел, как Браун вышел из кабинета очень нервным"},
				PartialFacts: []string{"Кажется, Ланг с кем-то ссорился перед смертью — доносились голоса"},
				FalseBeliefs: []string{"Я думаю, что Елена могла отравить мужа — у них были проблемы в браке"},
			},
			Secrets:                 []string{"Я должен Лангу крупную сумму (залез в кассу). Ланг знал, но не увольнял", "Я убрал свой бокал из кабинета, чтобы не было вопросов по отпечаткам"},
			Relationships:           map[int]string{3: "ненавидит — Браун хотел его уволить", 4: "защищает — Анна его племянница"},
			Memories:                []domain.Memory{{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Content: "Я видел, как Браун вышел из кабинета в 21:00, злой.", IsTrue: true, Timestamp: "21:00"}, {ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Content: "Ланг был жив, когда я заходил с виски в 21:15.", IsTrue: true, Timestamp: "21:15"}},
			Trust:                   55,
			InterrogationsRemaining: 3,
		},
		{
			PrototypeID: 2, Name: "Елена Соколова", Age: 42, Profession: "домохозяйка",
			ImagePath:   "assets/characters/elena_sokolova.png",
			Personality: "Эмоциональная, вспыльчивая, но ранимая. Говорит быстро, часто перебивает. Хочет казаться безразличной, но на деле очень переживает.",
			AudioToneID: "tone_female_high",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я слышала крик из кабинета около 22:00", "Муж собирался разводиться со мной"},
				PartialFacts: []string{"Мне кажется, у мужа были проблемы с бизнесом"},
				FalseBeliefs: []string{"Я уверена, что Иван что-то скрывает — он слишком спокоен"},
			},
			Secrets:                 []string{"У меня были отношения с Орловым. Ланг узнал и хотел развода", "В ночь убийства я видела, как Орлов выходил из кабинета в 21:45"},
			Relationships:           map[int]string{5: "любовники — Орлов был её любовником", 3: "не доверяет — Браун странно влиял на мужа"},
			Memories:                []domain.Memory{{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Content: "В 22:00 я слышала громкий спор, а потом глухой удар.", IsTrue: true, Timestamp: "22:00"}, {ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Content: "Муж собирался объявить о разводе на следующей неделе.", IsTrue: true, Timestamp: "неделю назад"}},
			Trust:                   45,
			InterrogationsRemaining: 3,
		},
		{
			PrototypeID: 3, Name: "Майкл Браун", Age: 48, Profession: "деловой партнёр",
			ImagePath:   "assets/characters/michael_brown.png",
			Personality: "Харизматичный, уверенный в себе, умело манипулирует. Говорит спокойно, с лёгкой усмешкой. Всегда контролирует эмоции.",
			AudioToneID: "tone_male_mid",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я был в кабинете с Лангом с 20:30 до 21:00", "Ланг обнаружил недостачу в отчётах"},
				PartialFacts: []string{"Иван должен Лангу деньги, я предлагал его уволить"},
				FalseBeliefs: []string{"Думаю, что записку с угрозой написала Анна — она боялась увольнения"},
			},
			Secrets:                 []string{"Я отравил Ланга цианидом, подмешав яд в графин с виски в 20:45", "Я стёр отпечатки с пузырька и выбросил его в корзину"},
			Relationships:           map[int]string{1: "презирает — Иван воровал, Браун хотел его выгнать", 4: "подозревает — думает, что Анна могла что-то видеть"},
			Memories:                []domain.Memory{{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005"), Content: "В 20:45, пока Ланг ненадолго вышел, я подсыпал яд в графин.", IsTrue: false, Timestamp: "20:45"}, {ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"), Content: "Ланг сам налил себе виски в 20:50 и выпил.", IsTrue: true, Timestamp: "20:50"}},
			Trust:                   30,
			InterrogationsRemaining: 3,
		},
		{
			PrototypeID: 4, Name: "Анна Коваль", Age: 29, Profession: "горничная",
			ImagePath:   "assets/characters/anna_koval.png",
			Personality: "Застенчивая, тревожная, боится потерять работу. Говорит тихо, запинается. Старается быть незаметной, но глаза выдают страх.",
			AudioToneID: "tone_female_soft",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я нашла тело в 22:15, когда зашла убрать кабинет", "В 21:30 я видела, как Орлов шёл к кабинету"},
				PartialFacts: []string{"Слышала, что у Ланга и Соколовой были проблемы"},
				FalseBeliefs: []string{"Думаю, что Иван мог отравить — он очень нервничал весь вечер"},
			},
			Secrets:                 []string{"Я видела, как Иван выносил бокал из кабинета в 21:30", "Я боюсь, что меня уволят, если я расскажу слишком много"},
			Relationships:           map[int]string{1: "боится — Иван требует, чтобы молчала о бокале", 5: "странный — Орлов слишком интересовался её показаниями"},
			Memories:                []domain.Memory{{ID: uuid.MustParse("00000000-0000-0000-0000-000000000007"), Content: "В 21:30 я видела Ивана, выходящего из кабинета с пустыми руками.", IsTrue: false, Timestamp: "21:30"}, {ID: uuid.MustParse("00000000-0000-0000-0000-000000000008"), Content: "Тело лежало лицом вниз, на столе стоял пустой бокал.", IsTrue: true, Timestamp: "22:15"}},
			Trust:                   40,
			InterrogationsRemaining: 3,
		},
		{
			PrototypeID: 5, Name: "Дмитрий Орлов", Age: 61, Profession: "адвокат",
			ImagePath:   "assets/characters/dmitry_orlov.png",
			Personality: "Циничный, расчётливый, за словом в карман не лезет. Говорит чётко, рублеными фразами. Привык контролировать ситуацию.",
			AudioToneID: "tone_male_raspy",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я подходил к кабинету в 21:45, дверь была заперта", "У меня были деловые отношения с Лангом"},
				PartialFacts: []string{"Ланг кого-то шантажировал, судя по тому, что я нашёл в документах"},
				FalseBeliefs: []string{"Уверен, что убийца — Иван. У него был мотив (долг) и возможность"},
			},
			Secrets:                 []string{"У меня роман с Еленой Соколовой, женой Ланга", "Ланг нашёл в кабинете мой конверт с компроматом на сделку с Брауном"},
			Relationships:           map[int]string{2: "любовники — роман с Еленой", 3: "деловой партнёр — провернули сделку с недвижимостью"},
			Memories:                []domain.Memory{{ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"), Content: "Я подошёл к кабинету в 21:45, дверь была заперта изнутри.", IsTrue: true, Timestamp: "21:45"}, {ID: uuid.MustParse("00000000-0000-0000-0000-00000000000a"), Content: "Ланг нашёл конверт с документами и вызвал меня на разговор.", IsTrue: true, Timestamp: "20:00"}},
			Trust:                   50,
			InterrogationsRemaining: 3,
		},
	}
	for i := range chars {
		chars[i].ID = uuid.New()
	}
	return chars
}

var interrogationResponses = map[int]map[string]ports.LlmInterrogationResponse{
	1: {
		"neutral":    {Answer: "Я уже всё рассказал полиции. В 21:15 я занёс мистеру Лангу виски, как обычно. Он был жив и здоров. Больше я ничего не видел.", AttitudeDelta: -2, Statements: []string{"Я заносил виски в кабинет в 21:15"}},
		"aggressive": {Answer: "Послушайте, я 30 лет работаю в этом доме! Я не обязан отвечать на ваши подозрения. Если хотите знать, спросите у Брауна — он последним вышел от него.", AttitudeDelta: -15, Statements: []string{"Браун последним вышел из кабинета"}},
		"kind":       {Answer: "Хорошо, я скажу, но это между нами. Мистер Браун был очень зол, когда вышел из кабинета. Я слышал, как он говорил по телефону о каких-то \"пропавших деньгах\".", AttitudeDelta: 10, Statements: []string{"Браун говорил о пропавших деньгах"}},
	},
}

type MockLlmService struct{}

func NewMockLlmService() *MockLlmService {
	return &MockLlmService{}
}

func (m *MockLlmService) GenerateScenario(_ context.Context, characters []domain.CharacterPrototype) (*ports.ScenarioOutput, error) {
	chars := makeCharacters()
	crime := groundTruth
	crime.PerpetratorID = chars[2].ID
	return &ports.ScenarioOutput{
		Crime:      crime,
		Timeline:   scenarioTimeline,
		CaseName:   "Дело об убийстве Роберта Ланга",
		CaseBrief:  "В особняке Лангов обнаружено тело главы семейства. Все подозреваемые — члены семьи и прислуга.",
		Characters: chars,
		Evidence:   scenarioEvidence,
	}, nil
}

func (m *MockLlmService) RespondInInterrogation(_ context.Context, character domain.Character, playerMessage string) (*ports.LlmInterrogationResponse, error) {
	msg := strings.ToLower(playerMessage)
	charID := character.PrototypeID

	isAggressive := strings.Contains(msg, "убил") || strings.Contains(msg, "лжёшь") || strings.Contains(msg, "вы врёте") || strings.Contains(msg, "докажите")
	isKind := strings.Contains(msg, "пожалуйста") || strings.Contains(msg, "расскажите") || strings.Contains(msg, "помогите")

	if responses, ok := interrogationResponses[charID]; ok {
		if isAggressive {
			r := responses["aggressive"]
			return &r, nil
		}
		if isKind {
			r := responses["kind"]
			return &r, nil
		}
		if r, ok := responses["neutral"]; ok {
			return &r, nil
		}
	}

	return m.genericResponse(&character), nil
}

func (m *MockLlmService) genericResponse(char *domain.Character) *ports.LlmInterrogationResponse {
	switch {
	case char.Trust < 25:
		return &ports.LlmInterrogationResponse{
			Answer:        fmt.Sprintf("%s: Я не хочу больше говорить. Вызовите моего адвоката.", char.Name),
			AttitudeDelta: -5,
			Statements:    []string{"Отказался отвечать на вопросы"},
		}
	case char.Trust < 50:
		return &ports.LlmInterrogationResponse{
			Answer:        fmt.Sprintf("%s: Я уже говорил полиции всё, что знаю. У меня нет времени на эти игры.", char.Name),
			AttitudeDelta: -3,
			Statements:    []string{"Повторил, что уже говорил полиции"},
		}
	default:
		return &ports.LlmInterrogationResponse{
			Answer:        fmt.Sprintf("%s: Я ничего не скрываю. Задавайте ваши вопросы, я отвечу, если смогу.", char.Name),
			AttitudeDelta: 0,
			Statements:    []string{"Готов отвечать на вопросы"},
		}
	}
}

func (m *MockLlmService) EvaluateReport(_ context.Context, playerReport domain.FinalReport, groundTruth domain.Crime) (*ports.LlmFeedbackResponse, error) {
	perpName := "Майкл Браун"
	perpParts := strings.Split(strings.ToLower(perpName), " ")

	whoCorrect := false
	for _, part := range perpParts {
		if strings.Contains(strings.ToLower(playerReport.Who), part) {
			whoCorrect = true
			break
		}
	}

	whyWords := strings.FieldsFunc(strings.ToLower(groundTruth.Motive), func(r rune) bool { return r == ' ' || r == ',' })
	whyCorrect := false
	for _, w := range whyWords {
		if len(w) > 3 && strings.Contains(strings.ToLower(playerReport.Why), w) {
			whyCorrect = true
			break
		}
	}

	howCorrect := strings.Contains(strings.ToLower(playerReport.How), strings.ToLower(groundTruth.Method[:min(len(groundTruth.Method), 10)]))
	whenCorrect := strings.Contains(playerReport.When, groundTruth.TimeOfCrime)

	evidenceCorrect := false
	for _, ev := range scenarioEvidence {
		halfLen := len(ev.Name) / 2
		if halfLen < 1 {
			halfLen = 1
		}
		if strings.Contains(strings.ToLower(playerReport.Evidence), strings.ToLower(ev.Name[:halfLen])) {
			evidenceCorrect = true
			break
		}
	}

	return &ports.LlmFeedbackResponse{
		NarrativeFeedback: generateFeedbackText(whoCorrect, whyCorrect, howCorrect, whenCorrect, evidenceCorrect),
		BreakdownDetails: map[string]string{
			"who":      feedbackDetail("who", whoCorrect, fmt.Sprintf("Верно, преступник — %s.", perpName), fmt.Sprintf("Нет, убийца — %s. %s", perpName, groundTruth.Motive)),
			"why":      feedbackDetail("why", whyCorrect, fmt.Sprintf("Правильно — %s.", groundTruth.Motive), fmt.Sprintf("На самом деле %s убил, чтобы скрыть %s.", perpName, strings.ToLower(groundTruth.Motive))),
			"how":      feedbackDetail("how", howCorrect, fmt.Sprintf("Да, %s — верный способ.", groundTruth.Method), fmt.Sprintf("Нет, орудие — %s.", groundTruth.Method)),
			"when":     feedbackDetail("when", whenCorrect, fmt.Sprintf("Время указано верно — %s.", groundTruth.TimeOfCrime), fmt.Sprintf("Преступление произошло около %s.", groundTruth.TimeOfCrime)),
			"evidence": feedbackDetail("evidence", evidenceCorrect, "Хорошо, ключевые улики учтены.", fmt.Sprintf("Упущены %s.", joinEvidenceNames(scenarioEvidence))),
		},
		MissedFacts: []string{},
	}, nil
}

func feedbackDetail(key string, correct bool, correctText, wrongText string) string {
	if correct {
		return correctText
	}
	return wrongText
}

func generateFeedbackText(who, why, how, when, evidence bool) string {
	correct := 0
	for _, b := range []bool{who, why, how, when, evidence} {
		if b {
			correct++
		}
	}
	switch {
	case correct >= 4:
		return "Отличная работа! Вы практически полностью восстановили картину преступления. Ваш анализ заслуживает похвалы."
	case correct >= 2:
		return "Неплохо, но есть над чем поработать. Обратите внимание на время преступления и способ."
	default:
		return "К сожалению, ваша версия далека от истины. Попробуйте ещё раз изучить улики и показания."
	}
}

func joinEvidenceNames(ev []domain.Evidence) string {
	names := make([]string, len(ev))
	for i, e := range ev {
		names[i] = strings.ToLower(e.Name)
	}
	return strings.Join(names, ", ")
}

func (m *MockLlmService) RunAction(_ context.Context, actionName string, evidenceID *uuid.UUID, characterID *uuid.UUID, alibiText *string) (string, error) {
	if evidenceID != nil {
		return fmt.Sprintf("Результат %s для улики:\n\nСледов ДНК, соответствующих подозреваемым, не обнаружено. Материал передан в лабораторию для дальнейшего анализа.", strings.ToLower(actionName)), nil
	}
	if characterID != nil {
		return fmt.Sprintf("Результат запроса %s:\n\nДанные получены. Анализ показал, что в указанный период никаких подозрительных операций не зафиксировано.", strings.ToLower(actionName)), nil
	}
	if alibiText != nil {
		return fmt.Sprintf("Результат проверки алиби:\n\nЗапрос: %s\n\nАгенты проверили указанное место. Показания свидетелей частично подтверждают алиби, но есть расхождения по времени. Рекомендуется продолжить расследование.", *alibiText), nil
	}
	return fmt.Sprintf("Результат %s:\n\nЗапрос выполнен. Полученные данные добавлены в дело.", strings.ToLower(actionName)), nil
}

var _ ports.LlmService = (*MockLlmService)(nil)

func strPtr(s string) *string {
	return &s
}
