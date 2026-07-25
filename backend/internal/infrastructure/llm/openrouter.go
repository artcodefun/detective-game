package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

type OpenRouterClient struct {
	cli    *http.Client
	apiKey string
	model  string
}

func NewOpenRouterClient(apiKey, model string) *OpenRouterClient {
	return &OpenRouterClient{
		cli:    &http.Client{Timeout: 2 * 60 * time.Second},
		apiKey: apiKey,
		model:  model,
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model          string         `json:"model"`
	Messages       []chatMessage  `json:"messages"`
	ResponseFormat map[string]any `json:"response_format,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *OpenRouterClient) chat(ctx context.Context, messages []chatMessage, jsonMode bool) (string, error) {
	req := chatRequest{
		Model:    c.model,
		Messages: messages,
	}
	if jsonMode {
		req.ResponseFormat = map[string]any{"type": "json_object"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.cli.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("openrouter error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

type llmMemory struct {
	Content   string `json:"content"`
	IsTrue    bool   `json:"is_true"`
	Timestamp string `json:"timestamp"`
}

type llmCharacter struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	Age           int               `json:"age"`
	Profession    string            `json:"profession"`
	Personality   string            `json:"personality"`
	Gender        string            `json:"gender"`
	KnownFacts    []string          `json:"known_facts"`
	PartialFacts  []string          `json:"partial_facts"`
	FalseBeliefs  []string          `json:"false_beliefs"`
	Secrets       []string          `json:"secrets"`
	Relationships map[string]string `json:"relationships"`
	Memories      []llmMemory       `json:"memories"`
	Trust         int               `json:"trust"`
}

type llmEvidence struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	DetailedDescription string `json:"detailed_description"`
	Type                string `json:"type"`
}

type llmCrime struct {
	CrimeType     string `json:"crime_type"`
	Victim        string `json:"victim"`
	PerpetratorID int    `json:"perpetrator_id"`
	Motive        string `json:"motive"`
	Method        string `json:"method"`
	TimeOfCrime   string `json:"time_of_crime"`
}

type llmTimelineEntry struct {
	Time        string `json:"time"`
	Event       string `json:"event"`
	CharacterID *int   `json:"character_id"`
}

type llmScenarioResponse struct {
	Crime      llmCrime           `json:"crime"`
	Timeline   []llmTimelineEntry `json:"timeline"`
	Characters []llmCharacter     `json:"characters"`
	Evidence   []llmEvidence      `json:"evidence"`
}

func (c *OpenRouterClient) GenerateScenario(ctx context.Context) (*ports.ScenarioOutput, error) {
	systemPrompt := `Ты — генератор детективных сценариев для игры. Создай запутанное дело об убийстве.

Верни ТОЛЬКО JSON, без текста до или после, без Markdown-форматирования:

{
  "crime": {
    "crime_type": "murder",
    "victim": "имя жертвы",
    "perpetrator_id": 1,
    "motive": "мотив преступления",
    "method": "способ убийства",
    "time_of_crime": "22:15"
  },
  "timeline": [
    {"time": "19:00", "event": "событие", "character_id": 1},
    {"time": "20:30", "event": "событие", "character_id": null}
  ],
  "characters": [
    {
      "id": 1,
      "name": "Имя Фамилия",
      "age": 55,
      "profession": "профессия",
      "personality": "описание характера и манеры речи (2-3 предложения)",
      "gender": "male или female",
      "known_facts": ["факт, который персонаж знает точно"],
      "partial_facts": ["факт, известный частично"],
      "false_beliefs": ["ложное убеждение"],
      "secrets": ["секрет"],
      "relationships": {"2": "описание отношений с персонажем 2"},
      "memories": [
        {"content": "воспоминание", "is_true": true, "timestamp": "21:00"}
      ],
      "trust": 55
    }
  ],
  "evidence": [
    {
      "name": "название улики",
      "description": "краткое описание (что видно невооружённым глазом)",
      "detailed_description": "что предстоит выяснить", 
      "type": "physical"
    }
  ]
}

=== ТАЙМЛАЙН ===
Это ОБЪЕКТИВНАЯ хроника событий, известная только тебе (игрок её не видит). Она не должна быть "отчётом" — это истина, что происходило на самом деле. Должен содержать:
- 8-12 событий в хронологическом порядке
- Каждое появление/исчезновение значимой улики ДОЛЖНО быть в таймлайне (кто и когда оставил предмет, передвинул, спрятал)
- Перемещения всех персонажей в ключевые моменты
- Момент преступления с указанием убийцы
- Момент обнаружения тела
- events — конкретные действия, а не "персонаж нервничал"

=== УЛИКИ (evidence) ===
Это предметы, найденные НА МЕСТЕ ПРЕСТУПЛЕНИЯ. Только то, что криминалист видит при первом осмотре:
- physical: орудие убийства, одежда со следами, разбитые предметы, следы обуви, волокна ткани, бокалы, пузырьки
- digital: телефон, ноутбук, флешка, диск — САМИ УСТРОЙСТВА. НЕ ИХ СОДЕРЖИМОЕ. Контент устройств — предмет действий игрока (call_history, transaction_check и т.д.)
- document: записки, письма, финансовые отчёты, дневники, контракты

Примеры того, что НЕЛЬЗЯ класть в evidence:
- ❌ "запись с камеры наблюдения" (это результат camera_review)
- ❌ "расшифровка звонков" (это результат call_history)
- ❌ "отпечатки принадлежат X" (это результат fingerprints)
- ❌ "ДНК жертвы на орудии" (это результат dna_analysis)

=== DETAILED DESCRIPTION ===
Описывай ТОЛЬКО то, что видно при осмотре. Не результаты анализов. Примеры:
- ✅ "Бокал с остатками жидкости. Требуется анализ содержимого и снятие отпечатков."
- ✅ "Телефон, заблокирован паролем. Требуется доступ к звонкам и сообщениям."
- ✅ "Пятна на ковре, похожие на кровь. Требуется анализ ДНК."
- ❌ "Отпечатки пальцев принадлежат Ивану."
- ❌ "Анализ показал наличие цианида."
- ❌ "В телефоне найдена переписка с жертвой."

=== ОБЩИЕ ПРАВИЛА ===
- characters.id — порядковый номер (1,2,3,4,5)
- perpetrator_id — номер персонажа-убийцы (совпадает с characters.id)
- timeline[].character_id — номер персонажа или null
- Генерируй 5 УНИКАЛЬНЫХ персонажей с разными именами, возрастами и профессиями
- personality — подробное описание характера и манеры речи (2-3 предложения)
- gender — "male" или "female"
- Все значимые улики ДОЛЖНЫ появляться в таймлайне
- relationships — ключ это номер персонажа, значение — описание отношений
- secrets — у каждого персонажа должен быть хотя бы один секрет
- memories — 1-2 воспоминания на персонажа
- trust — начальное доверие к детективу (0-100). У убийцы низкое (10-30)
- Всего 5 персонажей и ровно 5 улик
- ВАЖНО: проверь, что все JSON-массивы и объекты корректно закрыты. Не оставляй ключи без значений.`

	content, err := c.chat(ctx, []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "Сгенерируй детективный сценарий."},
	}, true)
	if err != nil {
		return nil, err
	}

	var llmResp llmScenarioResponse
	if err := json.Unmarshal([]byte(content), &llmResp); err != nil {
		content2, retryErr := c.chat(ctx, []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: "Сгенерируй детективный сценарий."},
		}, true)
		if retryErr != nil {
			return nil, fmt.Errorf("parse scenario (retry failed): %w\n%s", err, content)
		}
		if err2 := json.Unmarshal([]byte(content2), &llmResp); err2 != nil {
			return nil, fmt.Errorf("parse scenario: %w\n%s", err, content)
		}
	}

	perpetratorID := uuid.Nil
	idByLLMID := make(map[int]uuid.UUID, len(llmResp.Characters))
	domainChars := make([]domain.Character, len(llmResp.Characters))
	for i, lc := range llmResp.Characters {
		charID := uuid.New()
		idByLLMID[lc.ID] = charID

		var memories []domain.Memory
		for _, m := range lc.Memories {
			memories = append(memories, domain.Memory{
				ID:        uuid.New(),
				Content:   m.Content,
				IsTrue:    m.IsTrue,
				Timestamp: m.Timestamp,
			})
		}

		domainChars[i] = domain.Character{
			ID:          charID,
			Name:        lc.Name,
			Age:         lc.Age,
			Profession:  lc.Profession,
			Personality: lc.Personality,
			Gender:      domain.Gender(lc.Gender),
			Knowledge: domain.CharacterKnowledge{
				KnownFacts:   lc.KnownFacts,
				PartialFacts: lc.PartialFacts,
				FalseBeliefs: lc.FalseBeliefs,
			},
			Secrets:                 lc.Secrets,
			Relationships:           lc.Relationships,
			Memories:                memories,
			Trust:                   lc.Trust,
			InterrogationsRemaining: domain.MaxInterrogations,
		}
	}
	if id, ok := idByLLMID[llmResp.Crime.PerpetratorID]; ok {
		perpetratorID = id
	}

	domainEvidence := make([]domain.Evidence, len(llmResp.Evidence))
	for i, le := range llmResp.Evidence {
		domainEvidence[i] = domain.Evidence{
			ID:                  uuid.New(),
			Name:                le.Name,
			Description:         le.Description,
			DetailedDescription: le.DetailedDescription,
			Type:                domain.EvidenceType(le.Type),
		}
	}

	domainTimeline := make([]domain.TimelineEntry, len(llmResp.Timeline))
	for i, te := range llmResp.Timeline {
		var charID *uuid.UUID
		if te.CharacterID != nil {
			if id, ok := idByLLMID[*te.CharacterID]; ok {
				cid := id
				charID = &cid
			}
		}
		domainTimeline[i] = domain.TimelineEntry{
			Time:        te.Time,
			Event:       te.Event,
			CharacterID: charID,
		}
	}

	caseName, caseBrief := c.generateCaseBrief(ctx, domainChars, llmResp.Crime, llmResp.Evidence)

	return &ports.ScenarioOutput{
		CaseName:  caseName,
		CaseBrief: caseBrief,
		Crime: domain.Crime{
			Type:          domain.CrimeType(llmResp.Crime.CrimeType),
			Victim:        llmResp.Crime.Victim,
			PerpetratorID: perpetratorID,
			Motive:        llmResp.Crime.Motive,
			Method:        llmResp.Crime.Method,
			TimeOfCrime:   llmResp.Crime.TimeOfCrime,
		},
		Timeline:   domain.Timeline{Entries: domainTimeline},
		Characters: domainChars,
		Evidence:   domainEvidence,
	}, nil
}

func (c *OpenRouterClient) generateCaseBrief(ctx context.Context, chars []domain.Character, crime llmCrime, evidence []llmEvidence) (string, string) {
	var charList strings.Builder
	for _, ch := range chars {
		charList.WriteString(fmt.Sprintf("- %s, %d лет, %s\n", ch.Name, ch.Age, ch.Profession))
	}

	var evList strings.Builder
	for _, ev := range evidence {
		evList.WriteString(fmt.Sprintf("- %s: %s\n", ev.Name, ev.Description))
	}

	briefPrompt := fmt.Sprintf(`Сгенерируй название дела и официальный документ для детектива.

Контекст:
Жертва: %s
Время преступления: %s
Способ: %s
Мотив: %s

Подозреваемые:
%s
Улики:
%s

Верни JSON:
{
  "case_name": "название дела (короткое, атмосферное)",
  "case_brief": "Markdown-документ по формату ниже"
}

Формат case_brief:

# Дело №[номер]

**Место преступления:** [придумай адрес]
**Предполагаемое время:** %s
**Тип преступления:** Убийство

## Обстоятельства

[2-3 абзаца: где найдено тело, в каком состоянии, что заметили криминалисты]

## Список подозреваемых

| Имя | Возраст | Роль | Связь с жертвой |
|-----|---------|------|-----------------|
[по одной строке на каждого персонажа из контекста]

## Улики с места преступления

[маркированный список улик из контекста]

---
*Дело передано:* [вымышленное имя старшего следователя]
*Дата открытия:* [сегодняшняя дата]`,
		crime.Victim, crime.TimeOfCrime, crime.Method, crime.Motive,
		charList.String(), evList.String(), crime.TimeOfCrime,
	)

	content, err := c.chat(ctx, []chatMessage{
		{Role: "system", Content: "Ты — помощник детектива. Генерируешь официальные документы."},
		{Role: "user", Content: briefPrompt},
	}, true)
	if err != nil {
		return "Дело №" + crime.Victim, "_документ не сгенерирован_"
	}

	var resp struct {
		CaseName  string `json:"case_name"`
		CaseBrief string `json:"case_brief"`
	}
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return "Дело №" + crime.Victim, content
	}

	return resp.CaseName, resp.CaseBrief
}

func (c *OpenRouterClient) RespondInInterrogation(ctx context.Context, character domain.Character, playerMessage string) (*ports.LlmInterrogationResponse, error) {
	systemPrompt := fmt.Sprintf(`Ты — персонаж в детективной игре. Отвечай в соответствии с характером и знаниями.

Имя: %s
Возраст: %d
Профессия: %s
Характер: %s
Текущий уровень доверия к следователю: %d из 100.

Известные факты: %s
Частично известные факты: %s
Ложные убеждения: %s
Секреты: %s

ВАЖНО: доверие определяет ТОН и ОТКРОВЕННОСТЬ ответа. Это главный параметр.

Правила:
- Доверие 0-15: персонаж вас НЕНАВИДИТ. Отвечай грубо, отказывайся говорить, требуй адвоката, угрожай.
- Доверие 16-30: персонаж вам не доверяет. Отвечай враждебно, уклончиво, односложно.
- Доверие 31-50: персонаж насторожен. Отвечай нейтрально, с неохотой, увиливай.
- Доверие 51-75: персонаж готов сотрудничать. Отвечай открыто, но без деталей.
- Доверие 76-100: персонаж доверяет. Отвечай подробно, делись наблюдениями (но НЕ секретами).
- На прямые обвинения ("ты убил", "ты лжёшь") — реагируй острее, attitude_delta дополнительно отрицательный
- На вежливые просьбы ("пожалуйста", "расскажите") — attitude_delta дополнительно положительный
- Секреты — самая ценная информация. Раскрывай их только когда персонаж ЗАГНАН В УГОЛ: доверие >60 И игрок предъявил улику или поймал на лжи.
- Никогда не говори "я не знаю" — отвечай в роли`,
		character.Name,
		character.Age,
		character.Profession,
		character.Personality,
		character.Trust,
		strings.Join(character.Knowledge.KnownFacts, "; "),
		strings.Join(character.Knowledge.PartialFacts, "; "),
		strings.Join(character.Knowledge.FalseBeliefs, "; "),
		strings.Join(character.Secrets, "; "),
	)

	content, err := c.chat(ctx, []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf(`Ответь на реплику следователя. Верни ТОЛЬКО JSON:

{
  "answer": "твой ответ (1-3 предложения)",
  "attitude_delta": -5..5 (изменение доверия),
  "statements": ["ключевое утверждение из ответа"]
}

Реплика следователя: %s`, playerMessage)},
	}, true)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Answer        string   `json:"answer"`
		AttitudeDelta int      `json:"attitude_delta"`
		Statements    []string `json:"statements"`
	}
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w\n%s", err, content)
	}

	return &ports.LlmInterrogationResponse{
		Answer:        resp.Answer,
		AttitudeDelta: resp.AttitudeDelta,
		Statements:    resp.Statements,
	}, nil
}

func (c *OpenRouterClient) EvaluateReport(ctx context.Context, playerReport domain.FinalReport, groundTruth domain.Crime) (*ports.LlmFeedbackResponse, error) {
	systemPrompt := fmt.Sprintf(`Ты оцениваешь финальный отчёт детектива. Вот что известно об убийстве на самом деле:

Жертва: %s
Мотив преступника: %s
Способ убийства: %s
Время преступления: %s

Сравни отчёт игрока с этими фактами и верни ТОЛЬКО JSON:

{
  "narrative_feedback": "развёрнутый текстовый отзыв (3-5 предложений на русском). Оцени общее качество расследования.",
  "breakdown_details": {
    "who": "совпадает ли имя преступника в отчёте с фактами? дай краткую оценку",
    "why": "насколько мотив в отчёте близок к реальному? дай краткую оценку",
    "how": "правильно ли указан способ? дай краткую оценку",
    "when": "правильно ли указано время? дай краткую оценку",
    "evidence": "насколько убедительно описаны улики? дай краткую оценку"
  },
  "missed_facts": ["какие важные детали игрок упустил"]
}`,
		groundTruth.Victim,
		groundTruth.Motive,
		groundTruth.Method,
		groundTruth.TimeOfCrime,
	)

	content, err := c.chat(ctx, []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf(`Отчёт игрока:
Кто: %s
Почему: %s
Как: %s
Когда: %s
Улики: %s`, playerReport.Who, playerReport.Why, playerReport.How, playerReport.When, playerReport.Evidence)},
	}, true)
	if err != nil {
		return nil, err
	}

	var resp struct {
		NarrativeFeedback string            `json:"narrative_feedback"`
		BreakdownDetails  map[string]string `json:"breakdown_details"`
		MissedFacts       []string          `json:"missed_facts"`
	}
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return nil, fmt.Errorf("parse feedback: %w\n%s", err, content)
	}

	return &ports.LlmFeedbackResponse{
		NarrativeFeedback: resp.NarrativeFeedback,
		BreakdownDetails:  resp.BreakdownDetails,
		MissedFacts:       resp.MissedFacts,
	}, nil
}

func (c *OpenRouterClient) RunAction(ctx context.Context, actionName string, evidenceID *uuid.UUID, characterID *uuid.UUID, alibiText *string) (string, error) {
	var contextParts []string
	if evidenceID != nil {
		contextParts = append(contextParts, fmt.Sprintf("улика ID: %s", evidenceID))
	}
	if characterID != nil {
		contextParts = append(contextParts, fmt.Sprintf("персонаж ID: %s", characterID))
	}
	if alibiText != nil {
		contextParts = append(contextParts, fmt.Sprintf("текст алиби: %s", *alibiText))
	}

	actionLabels := map[string]string{
		"dna_analysis":      "анализ ДНК",
		"fingerprints":      "отпечатки пальцев",
		"alibi_check":       "проверка алиби",
		"camera_review":     "записи с камер наблюдения",
		"call_history":      "история звонков",
		"transaction_check": "банковские транзакции",
	}

	label := actionLabels[actionName]
	if label == "" {
		label = actionName
	}

	content, err := c.chat(ctx, []chatMessage{
		{Role: "system", Content: "Ты — криминалистическая лаборатория. Отвечай коротко, по делу, на русском языке. Не используй JSON."},
		{Role: "user", Content: fmt.Sprintf("Выполни запрос: %s. Контекст: %s. Верни результат в 2-3 предложениях.", label, strings.Join(contextParts, ", "))},
	}, false)
	if err != nil {
		return "", err
	}

	return content, nil
}
