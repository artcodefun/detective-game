package llm

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

var groundTruth = domain.Crime{
	Type:        domain.CrimeTypeMurder,
	Victim:      "Роберт Ланг",
	Motive:      "Ланг раскрыл хищение средств и собирался обратиться в полицию",
	Method:      "отравление цианидом (подмешан в виски)",
	TimeOfCrime: "22:15",
}

var scenarioTimeline = domain.Timeline{Entries: []domain.TimelineEntry{
	{Time: "19:00", Event: "Ужин в особняке. Все подозреваемые присутствуют."},
	{Time: "20:30", Event: "Ланг и Браун уходят в кабинет для разговора.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000103"); return &id }()},
	{Time: "21:00", Event: "Браун покидает кабинет. Ланг остается один.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000103"); return &id }()},
	{Time: "21:15", Event: "Иван Петров заходит в кабинет с виски.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000101"); return &id }()},
	{Time: "21:45", Event: "Орлов идёт в кабинет обсудить документы, но дверь заперта.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000105"); return &id }()},
	{Time: "22:00", Event: "Елена Соколова слышит шум из кабинета.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000102"); return &id }()},
	{Time: "22:15", Event: "Анна Коваль обнаруживает тело.", CharacterID: func() *uuid.UUID { id := uuid.MustParse("00000000-0000-0000-0000-000000000104"); return &id }()},
}}

var scenarioEvidence = []domain.Evidence{
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Бокал с остатками виски", Description: "Найден на столе в кабинете",
		DetailedDescription: "Хрустальный бокал с остатками жидкости. Требуется анализ содержимого и снятие отпечатков.",
		Type:                domain.EvidenceTypePhysical},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Пузырёк", Description: "Найден в мусорной корзине",
		DetailedDescription: "Стеклянный флакон без этикетки. Требуется анализ содержимого.",
		Type:                domain.EvidenceTypePhysical},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Финансовые документы", Description: "Разбросаны на столе",
		DetailedDescription: "Бухгалтерские отчёты с пометками. Требуется изучение содержимого.",
		Type:                domain.EvidenceTypeDocument},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "Запертый ящик стола", Description: "Ключ торчит снаружи",
		DetailedDescription: "Верхний ящик стола заперт, ключ оставлен в замке. Требуется вскрытие и осмотр содержимого.",
		Type:                domain.EvidenceTypePhysical},
	{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005"), Name: "Записка с угрозой", Description: "Найдена под дверью кабинета",
		DetailedDescription: "Анонимная записка. Бумага дешёвая, текст напечатан на принтере. Требуется анализ бумаги и чернил.",
		Type:                domain.EvidenceTypeDocument},
}

func makeCharacters() []domain.Character {
	chars := []domain.Character{
		{
			Name: "Иван Петров", Age: 55, Profession: "дворецкий", Gender: domain.GenderMale,
			Personality: "Консервативный, преданный семье, скрытный. Говорит медленно, с расстановкой.",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я принёс виски в кабинет в 21:15, Ланг был жив и раздражён", "В 21:00 я видел, как Браун вышел из кабинета очень нервным"},
				PartialFacts: []string{"Кажется, Ланг с кем-то ссорился перед смертью — доносились голоса"},
				FalseBeliefs: []string{"Я думаю, что Елена могла отравить мужа — у них были проблемы в браке"},
			},
			Secrets:       []string{"Я должен Лангу крупную сумму (залез в кассу). Ланг знал, но не увольнял", "Я убрал свой бокал из кабинета, чтобы не было вопросов по отпечаткам"},
			Relationships: map[string]string{"3": "ненавидит — Браун хотел его уволить", "4": "защищает — Анна его племянница"},
			Memories: []domain.Memory{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Content: "Я видел, как Браун вышел из кабинета в 21:00, злой.", IsTrue: true, Timestamp: "21:00"},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Content: "Ланг был жив, когда я заходил с виски в 21:15.", IsTrue: true, Timestamp: "21:15"},
			},
			Trust: 55, InterrogationsRemaining: 3,
		},
		{
			Name: "Елена Соколова", Age: 42, Profession: "домохозяйка", Gender: domain.GenderFemale,
			Personality: "Эмоциональная, вспыльчивая, но ранимая. Говорит быстро, часто перебивает.",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я слышала крик из кабинета около 22:00", "Муж собирался разводиться со мной"},
				PartialFacts: []string{"Мне кажется, у мужа были проблемы с бизнесом"},
				FalseBeliefs: []string{"Я уверена, что Иван что-то скрывает — он слишком спокоен"},
			},
			Secrets:       []string{"У меня были отношения с Орловым. Ланг узнал и хотел развода", "В ночь убийства я видела, как Орлов выходил из кабинета в 21:45"},
			Relationships: map[string]string{"5": "любовники — Орлов был её любовником", "3": "не доверяет — Браун странно влиял на мужа"},
			Memories: []domain.Memory{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Content: "В 22:00 я слышала громкий спор, а потом глухой удар.", IsTrue: true, Timestamp: "22:00"},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Content: "Муж собирался объявить о разводе на следующей неделе.", IsTrue: true, Timestamp: "неделю назад"},
			},
			Trust: 45, InterrogationsRemaining: 3,
		},
		{
			Name: "Майкл Браун", Age: 48, Profession: "деловой партнёр", Gender: domain.GenderMale,
			Personality: "Харизматичный, уверенный в себе, умело манипулирует. Говорит спокойно, с лёгкой усмешкой.",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я был в кабинете с Лангом с 20:30 до 21:00", "Ланг обнаружил недостачу в отчётах"},
				PartialFacts: []string{"Иван должен Лангу деньги, я предлагал его уволить"},
				FalseBeliefs: []string{"Думаю, что записку с угрозой написала Анна — она боялась увольнения"},
			},
			Secrets:       []string{"Я отравил Ланга цианидом, подмешав яд в графин с виски в 20:45", "Я стёр отпечатки с пузырька и выбросил его в корзину"},
			Relationships: map[string]string{"1": "презирает — Иван воровал, Браун хотел его выгнать", "4": "подозревает — думает, что Анна могла что-то видеть"},
			Memories: []domain.Memory{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000005"), Content: "В 20:45, пока Ланг ненадолго вышел, я подсыпал яд в графин.", IsTrue: false, Timestamp: "20:45"},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000006"), Content: "Ланг сам налил себе виски в 20:50 и выпил.", IsTrue: true, Timestamp: "20:50"},
			},
			Trust: 30, InterrogationsRemaining: 3,
		},
		{
			Name: "Анна Коваль", Age: 29, Profession: "горничная", Gender: domain.GenderFemale,
			Personality: "Застенчивая, тревожная, боится потерять работу. Говорит тихо, запинается.",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я нашла тело в 22:15, когда зашла убрать кабинет", "В 21:30 я видела, как Орлов шёл к кабинету"},
				PartialFacts: []string{"Слышала, что у Ланга и Соколовой были проблемы"},
				FalseBeliefs: []string{"Думаю, что Иван мог отравить — он очень нервничал весь вечер"},
			},
			Secrets:       []string{"Я видела, как Иван выносил бокал из кабинета в 21:30", "Я боюсь, что меня уволят, если я расскажу слишком много"},
			Relationships: map[string]string{"1": "боится — Иван требует, чтобы молчала о бокале", "5": "странный — Орлов слишком интересовался её показаниями"},
			Memories: []domain.Memory{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000007"), Content: "В 21:30 я видела Ивана, выходящего из кабинета с пустыми руками.", IsTrue: false, Timestamp: "21:30"},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000008"), Content: "Тело лежало лицом вниз, на столе стоял пустой бокал.", IsTrue: true, Timestamp: "22:15"},
			},
			Trust: 40, InterrogationsRemaining: 3,
		},
		{
			Name: "Дмитрий Орлов", Age: 61, Profession: "адвокат", Gender: domain.GenderMale,
			Personality: "Циничный, расчётливый, за словом в карман не лезет. Говорит чётко, рублеными фразами.",
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   []string{"Я подходил к кабинету в 21:45, дверь была заперта", "У меня были деловые отношения с Лангом"},
				PartialFacts: []string{"Ланг кого-то шантажировал, судя по тому, что я нашёл в документах"},
				FalseBeliefs: []string{"Уверен, что убийца — Иван. У него был мотив (долг) и возможность"},
			},
			Secrets:       []string{"У меня роман с Еленой Соколовой, женой Ланга", "Ланг нашёл в кабинете мой конверт с компроматом на сделку с Брауном"},
			Relationships: map[string]string{"2": "любовники — роман с Еленой", "3": "деловой партнёр — провернули сделку с недвижимостью"},
			Memories: []domain.Memory{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000009"), Content: "Я подошёл к кабинету в 21:45, дверь была заперта изнутри.", IsTrue: true, Timestamp: "21:45"},
				{ID: uuid.MustParse("00000000-0000-0000-0000-00000000000a"), Content: "Ланг нашёл конверт с документами и вызвал меня на разговор.", IsTrue: true, Timestamp: "20:00"},
			},
			Trust: 50, InterrogationsRemaining: 3,
		},
	}
	for i := range chars {
		chars[i].ID = uuid.New()
	}
	return chars
}

type MockLlmService struct{}

func NewMockLlmService() *MockLlmService {
	return &MockLlmService{}
}

func (m *MockLlmService) GenerateScenario(_ context.Context) (*ports.ScenarioOutput, error) {
	chars := makeCharacters()
	crime := groundTruth
	crime.PerpetratorID = chars[2].ID
	return &ports.ScenarioOutput{
		Crime:      crime,
		Timeline:   scenarioTimeline,
		CaseName:   "Дело об убийстве Роберта Ланга",
		CaseBrief:  "# Дело №2024-07-003\n\n**Место преступления:** Особняк Лангов, ул. Садовая, 15\n**Предполагаемое время:** 22:15\n**Тип преступления:** Убийство\n\n## Обстоятельства\n\nВ особняке Лангов обнаружено тело главы семейства Роберта Ланга. Тело найдено горничной в кабинете в 22:15. Криминалисты зафиксировали следы борьбы и несколько улик, включая бокал с остатками жидкости и финансовые документы. Все находившиеся в доме задержаны для допроса.\n\n## Список подозреваемых\n\n| Имя | Возраст | Роль | Связь с жертвой |\n|-----|---------|------|------------------|\n| Иван Петров | 55 | дворецкий | Прислуга, должен жертве деньги |\n| Елена Соколова | 42 | домохозяйка | Жена, на грани развода |\n| Майкл Браун | 48 | деловой партнёр | Бизнес-партнёр, финансовые разногласия |\n| Анна Коваль | 29 | горничная | Прислуга, племянница Ивана |\n| Дмитрий Орлов | 61 | адвокат | Юрист семьи, имеет дела с жертвой |\n\n## Улики с места преступления\n\n- Бокал с остатками жидкости — требуется анализ\n- Пузырёк без этикетки — требуется анализ содержимого\n- Финансовые документы с пометками — требуется изучение\n- Запертый ящик стола — требуется вскрытие\n- Записка с угрозой — требуется анализ бумаги\n\n---\n*Дело передано:* ст. следователь Громов А.В.\n*Дата открытия:* 15 июля 2024 г.",
		Characters: chars,
		Evidence:   scenarioEvidence,
	}, nil
}

func (m *MockLlmService) RespondInInterrogation(_ context.Context, character domain.Character, playerMessage string) (*ports.LlmInterrogationResponse, error) {
	msg := strings.ToLower(playerMessage)
	isAggressive := strings.Contains(msg, "убил") || strings.Contains(msg, "лжёшь") || strings.Contains(msg, "врёте")
	isKind := strings.Contains(msg, "пожалуйста") || strings.Contains(msg, "расскажите") || strings.Contains(msg, "помогите")

	if isAggressive {
		return m.responseByTrust(&character, "aggressive"), nil
	}
	if isKind {
		return m.responseByTrust(&character, "kind"), nil
	}
	return m.responseByTrust(&character, "neutral"), nil
}

func (m *MockLlmService) responseByTrust(char *domain.Character, tone string) *ports.LlmInterrogationResponse {
	switch {
	case char.Trust < 25:
		return &ports.LlmInterrogationResponse{
			Answer:        fmt.Sprintf("%s: Я не хочу больше говорить. Вызовите моего адвоката.", char.Name),
			AttitudeDelta: -5,
			Statements:    []string{"Отказался отвечать на вопросы"},
		}
	case char.Trust < 50:
		if tone == "aggressive" {
			return &ports.LlmInterrogationResponse{
				Answer:        fmt.Sprintf("%s: Да как вы смеете меня обвинять?! У вас нет доказательств!", char.Name),
				AttitudeDelta: -10,
				Statements:    []string{"Возмущён обвинениями"},
			}
		}
		return &ports.LlmInterrogationResponse{
			Answer:        fmt.Sprintf("%s: Я уже всё сказал. Честно, я не знаю, кто это сделал.", char.Name),
			AttitudeDelta: -2,
			Statements:    []string{"Повторяет ранее сказанное"},
		}
	default:
		if tone == "kind" {
			return &ports.LlmInterrogationResponse{
				Answer:        fmt.Sprintf("%s: Хорошо, я расскажу вам кое-что... Только между нами. %s", char.Name, randomFact(char)),
				AttitudeDelta: 5,
				Statements:    []string{"Поделился информацией"},
			}
		}
		return &ports.LlmInterrogationResponse{
			Answer:        fmt.Sprintf("%s: Я отвечу на ваши вопросы, насколько смогу.", char.Name),
			AttitudeDelta: 0,
			Statements:    []string{"Готов сотрудничать"},
		}
	}
}

func randomFact(char *domain.Character) string {
	if len(char.Knowledge.KnownFacts) > 0 {
		return char.Knowledge.KnownFacts[rand.Intn(len(char.Knowledge.KnownFacts))]
	}
	if len(char.Knowledge.PartialFacts) > 0 {
		return char.Knowledge.PartialFacts[rand.Intn(len(char.Knowledge.PartialFacts))]
	}
	return "Не могу сейчас сказать."
}

var _ ports.LlmService = (*MockLlmService)(nil)

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

	evidenceCorrect := len(strings.TrimSpace(playerReport.Evidence)) > 10

	return &ports.LlmFeedbackResponse{
		NarrativeFeedback: generateFeedbackText(whoCorrect, whyCorrect, howCorrect, whenCorrect, evidenceCorrect),
		BreakdownDetails: map[string]string{
			"who":      feedbackDetail("who", whoCorrect, fmt.Sprintf("Верно, преступник — %s.", perpName), fmt.Sprintf("Нет, убийца — %s.", perpName)),
			"why":      feedbackDetail("why", whyCorrect, fmt.Sprintf("Правильно — %s.", groundTruth.Motive), fmt.Sprintf("На самом деле убийство совершено, чтобы скрыть %s.", strings.ToLower(groundTruth.Motive))),
			"how":      feedbackDetail("how", howCorrect, fmt.Sprintf("Да, %s — верный способ.", groundTruth.Method), fmt.Sprintf("Нет, орудие — %s.", groundTruth.Method)),
			"when":     feedbackDetail("when", whenCorrect, fmt.Sprintf("Время указано верно — %s.", groundTruth.TimeOfCrime), fmt.Sprintf("Преступление произошло около %s.", groundTruth.TimeOfCrime)),
			"evidence": feedbackDetail("evidence", evidenceCorrect, "Хорошо, ключевые улики учтены.", "Упущены важные улики с места преступления."),
		},
		MissedFacts: []string{},
	}, nil
}

func (m *MockLlmService) RunAction(_ context.Context, actionName string, _ *uuid.UUID, _ *uuid.UUID, _ *string) (string, error) {
	return fmt.Sprintf("Запрос «%s» выполнен. Результаты добавлены в дело.", actionName), nil
}

func feedbackDetail(_ string, correct bool, correctText, wrongText string) string {
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
		return "Отличная работа! Вы практически полностью восстановили картину преступления."
	case correct >= 2:
		return "Неплохо, но есть над чем поработать. Обратите внимание на время и способ."
	default:
		return "К сожалению, ваша версия далека от истины. Попробуйте ещё раз изучить улики."
	}
}
