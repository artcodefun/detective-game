package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
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

func (c *OpenRouterClient) chat(ctx context.Context, stage string, messages []chatMessage, jsonMode bool) (content string, err error) {
	startedAt := time.Now()
	requestChars := 0
	for _, message := range messages {
		requestChars += len(message.Content)
	}
	slog.InfoContext(ctx, "llm request started", "stage", stage)
	slog.DebugContext(ctx, "llm request payload", "stage", stage, "json_mode", jsonMode, "message_count", len(messages), "request_chars", requestChars, "messages", messages)
	defer func() {
		attrs := []any{"stage", stage, "duration", time.Since(startedAt)}
		if err != nil {
			attrs = append(attrs, "error", err)
			slog.ErrorContext(ctx, "llm request failed", attrs...)
			return
		}
		slog.InfoContext(ctx, "llm request completed", attrs...)
		slog.DebugContext(ctx, "llm response payload", "stage", stage, "response_chars", len(content), "content", content)
	}()

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
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("openrouter returned HTTP %d: %s", resp.StatusCode, truncateForError(string(respBody)))
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

func truncateForError(value string) string {
	const limit = 500
	if len(value) <= limit {
		return value
	}
	return value[:limit] + "…"
}

type llmMemory struct {
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type llmCharacter struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	Age           int               `json:"age"`
	Profession    string            `json:"profession"`
	Personality   string            `json:"personality"`
	Gender        string            `json:"gender"`
	Opinions      []string          `json:"opinions"`
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

type llmPrivateCharacter struct {
	ID       int         `json:"id"`
	Opinions []string    `json:"opinions"`
	Secrets  []string    `json:"secrets"`
	Memories []llmMemory `json:"memories"`
}

type llmPrivateResponse struct {
	Characters []llmPrivateCharacter `json:"characters"`
}

func (c *OpenRouterClient) generateScenarioDraft(ctx context.Context, locale domain.Locale) (string, error) {
	systemPrompt := `Ты — автор и логический редактор детективного сценария. Создай внутренний объективный план нового дела об убийстве. Он известен только автору и прямо раскрывает всю истину: жертву, убийцу, мотив, подготовку, точный момент и способ убийства, сокрытие следов и обнаружение тела. Не создавай интригу и не оставляй неизвестных обстоятельств.

Верни Markdown строго с тремя разделами. Не используй JSON и не добавляй другие разделы.

## Участники
Сначала кратко опиши жертву и ровно пять уникальных живых подозреваемых; жертва не входит в их число. Для каждого подозреваемого укажи имя, возраст, профессию, характер, публичную связь с жертвой, значимые отношения с другими участниками и его исходную цель на день убийства. Исходная цель должна быть правдоподобна и не раскрывать заранее тайный преступный замысел, сокрытие следов или иной решающий поворот дела: такие намерения могут сформироваться и быть объяснены только в последующих событиях таймлайна. Цель должна объяснять, почему он находится в этом месте и какие действия готов предпринять. Каждая указанная цель обязана вызвать хотя бы одно последующее действие в таймлайне; не добавляй декоративных целей, которые не влияют на события.

## Объективный таймлайн дня
Дай 14-20 событий в строгом хронологическом порядке с конкретным временем. Это полная объективная истина, а не отчёт следствия. Каждое существенное событие должно ясно отвечать на вопросы: кто, что, где сделал, зачем сделал и к чему это привело. Непосредственная мотивация должна вытекать из цели персонажа, его отношений или предыдущего события.

Не описывай действие без причины: персонаж не может брать, переносить, прятать, оставлять или использовать предмет «просто так». Для обычного действия достаточно короткой практической причины; для подготовки убийства, сокрытия следов и инсценировок обязательно объясни цель и ожидаемый эффект. Причина действия должна быть известна к этому моменту таймлайна: нельзя мотивировать поступок событием, которое произойдёт позже. События должны охватить действия всех подозреваемых, подготовку преступления, убийство, действия после него и обнаружение тела.

Один и тот же предмет всегда называй одинаково; не заменяй и не переименовывай его по ходу истории. Для каждого из пяти финальных предметов построй непрерывную цепочку: первое появление → каждое получение, перенос, использование или изменение → последнее состояние перед прибытием полиции. Не пропускай звенья этой цепочки. Предмет не может находиться в двух местах, быть одновременно у двух людей, быть и внутри другого предмета, и вне его, или получить новое повреждение без отдельного события передачи, перемещения или изменения состояния. Если предмет меняет владельца, местоположение, содержимое или состояние, явно опиши соответствующее действие. Орудие убийства и любые его части должны иметь один и тот же путь и финальное состояние: не подменяй его другим предметом. Для потенциальных ДНК, отпечатков и иных экспертиз укажи объективную физическую причину контакта или следа, но не объявляй лабораторные результаты.

## Предметы к прибытию полиции
Перечисли ровно пять физических предметов или документов, находящихся на месте к прибытию полиции. Для каждого укажи неизменное короткое название, точное финальное место и видимое состояние. Этот список только повторяет последнее состояние предмета из таймлайна: место, состояние, содержимое и название должны в точности совпадать. Не добавляй сюда новые перемещения, повреждения, упаковку, части предмета или иные факты, которых нет в последнем относящемся к нему событии. Большая часть предметов должна быть связана с раскрытием, но один или два могут быть случайными или второстепенными и не относиться к убийству. Не повторяй здесь полную цепочку предмета — она уже должна быть в таймлайне.

Не включай результаты действий игрока: видеозаписи, расшифровки звонков, банковские данные, результаты ДНК или отпечатков. Для цифровой улики допустимо только само устройство, но не его содержимое.

Перед ответом молча проверь, что у всех существенных действий есть мотивация, цели персонажей не противоречат их поступкам, а путь каждого из пяти предметов непрерывен до прибытия полиции. Для каждого пункта финального списка найди последнее событие о нём в таймлайне и сверь с ним название, место, состояние и содержимое; при любом расхождении исправь таймлайн или список до ответа.`

	content, err := c.chat(ctx, "scenario.draft", []chatMessage{
		{Role: "system", Content: systemPrompt + languageInstruction(locale)},
		{Role: "user", Content: "Создай внутренний план нового детективного дела об убийстве."},
	}, false)
	if err != nil {
		return "", fmt.Errorf("generate scenario draft: %w", err)
	}
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("generate scenario draft: empty response")
	}
	return content, nil
}

func (c *OpenRouterClient) GenerateScenario(ctx context.Context, locale domain.Locale) (output *ports.ScenarioOutput, err error) {
	startedAt := time.Now()
	slog.InfoContext(ctx, "scenario generation started")
	defer func() {
		if err != nil {
			slog.ErrorContext(ctx, "scenario generation failed", "duration", time.Since(startedAt), "error", err)
			return
		}
		slog.InfoContext(ctx, "scenario generation completed", "duration", time.Since(startedAt))
	}()

	draft, err := c.generateScenarioDraft(ctx, locale)
	if err != nil {
		return nil, err
	}

	systemPrompt := `Ты преобразуешь готовый внутренний план детективного дела в строгую структуру. План является единственным источником истины: не добавляй новых персонажей, предметов, событий, контактов или результатов анализов и не меняй причинную цепочку.

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
      "opinions": [],
      "secrets": [],
      "relationships": {"2": "описание отношений с персонажем 2"},
      "memories": [],
      "trust": 55
    }
  ],
  "evidence": [
    {
      "name": "название улики (например: Бокал с виски)",
      "description": "краткое описание: что за предмет, где найден (1 предложение)",
      "detailed_description": "развёрнутое описание: состояние, детали, следы (3-5 предложений, см. правила ниже)", 
      "type": "physical"
    }
  ]
}

Правила преобразования:
- Только извлекай и нормализуй сведения из синопсиса. Не исправляй сюжет и не добавляй факты.
- Сохрани имена людей и названия предметов дословно.
- Назначь пяти подозреваемым id от 1 до 5; perpetrator_id и character_id должны ссылаться на эти id. Жертва не входит в characters.
- Разбей прозу на 14-20 атомарных событий в хронологическом порядке. Сохраняй в event описанные место, действие, результат и непосредственную мотивацию существенного действия.
- Верни ровно пять улик из синопсиса. description — одно предложение о предмете и месте обнаружения; detailed_description — 3-5 предложений только о видимом состоянии и деталях.
- Не помещай в evidence результаты ДНК, отпечатков, истории звонков, транзакций или просмотра камер. Для digital evidence описывай устройство, а не его содержимое.
- relationships используют id других персонажей как ключи.
- opinions, secrets и memories оставь пустыми массивами.
- trust выбери по характеру персонажа, независимо от его виновности.
- gender — только "male" или "female"; crime_type — "murder"; тип улики — "physical", "digital" или "document".
- Проверь только корректность и полноту JSON перед ответом.`

	systemPrompt += languageInstruction(locale)

	content, err := c.chat(ctx, "scenario.structure", []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "Преобразуй этот план в JSON, сохранив все факты и связи:\n\n" + draft},
	}, true)
	if err != nil {
		return nil, fmt.Errorf("structure scenario: %w", err)
	}

	var llmResp llmScenarioResponse
	if err := json.Unmarshal([]byte(content), &llmResp); err != nil {
		slog.WarnContext(ctx, "scenario structure parse failed; retrying", "error", err, "response_preview", truncateForError(content))
		content2, retryErr := c.chat(ctx, "scenario.structure_retry", []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: "Преобразуй этот план в корректный JSON. Проверь синтаксис перед ответом:\n\n" + draft},
		}, true)
		if retryErr != nil {
			return nil, fmt.Errorf("parse scenario (retry failed): %w; response: %s", err, truncateForError(content))
		}
		if err2 := json.Unmarshal([]byte(content2), &llmResp); err2 != nil {
			return nil, fmt.Errorf("parse scenario after retry: %w; response: %s", err2, truncateForError(content2))
		}
	}

	if err := validateScenario(llmResp); err != nil {
		return nil, fmt.Errorf("validate structured scenario: %w", err)
	}
	review, err := c.reviewScenario(ctx, locale, llmResp)
	if err != nil {
		return nil, err
	}
	if len(review.Issues) > 0 {
		slog.WarnContext(ctx, "scenario review found issues; repairing", "issue_count", len(review.Issues), "issues", review.Issues)
		llmResp, err = c.repairScenario(ctx, locale, llmResp, review.Issues)
		if err != nil {
			return nil, err
		}
		if err := validateScenario(llmResp); err != nil {
			return nil, fmt.Errorf("validate repaired scenario: %w", err)
		}
	}
	private, err := c.generatePrivateCharacterData(ctx, locale, llmResp)
	if err != nil {
		return nil, fmt.Errorf("generate private character data: %w", err)
	}

	perpetratorName := ""
	idByLLMID := make(map[int]uuid.UUID, len(llmResp.Characters))
	domainChars := make([]domain.Character, len(llmResp.Characters))
	for i, lc := range llmResp.Characters {
		if lc.ID == llmResp.Crime.PerpetratorID {
			perpetratorName = lc.Name
		}
		if details, ok := private[lc.ID]; ok {
			lc.Opinions = details.Opinions
			lc.Secrets = details.Secrets
			lc.Memories = details.Memories
		}
		charID := uuid.New()
		idByLLMID[lc.ID] = charID

		var memories []domain.Memory
		for _, m := range lc.Memories {
			memories = append(memories, domain.Memory{
				Content:   m.Content,
				Timestamp: m.Timestamp,
			})
		}

		domainChars[i] = domain.Character{
			ID:                      charID,
			Name:                    lc.Name,
			Age:                     lc.Age,
			Profession:              lc.Profession,
			Personality:             lc.Personality,
			Gender:                  domain.Gender(lc.Gender),
			Opinions:                lc.Opinions,
			Secrets:                 lc.Secrets,
			Relationships:           lc.Relationships,
			Memories:                memories,
			Trust:                   lc.Trust,
			InterrogationsRemaining: domain.MaxInterrogations,
		}
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

	caseName, caseBrief := c.generateCaseBrief(ctx, locale, domainChars, llmResp.Crime, llmResp.Evidence, llmResp.Timeline)

	output = &ports.ScenarioOutput{
		CaseName:  caseName,
		CaseBrief: caseBrief,
		Crime: domain.Crime{
			Type:            domain.CrimeType(llmResp.Crime.CrimeType),
			Victim:          llmResp.Crime.Victim,
			PerpetratorName: perpetratorName,
			Motive:          llmResp.Crime.Motive,
			Method:          llmResp.Crime.Method,
			TimeOfCrime:     llmResp.Crime.TimeOfCrime,
		},
		Timeline:   domain.Timeline{Entries: domainTimeline},
		Characters: domainChars,
		Evidence:   domainEvidence,
	}
	return output, nil
}

type llmScenarioReview struct {
	Issues []string `json:"issues"`
}

func (c *OpenRouterClient) reviewScenario(ctx context.Context, locale domain.Locale, scenario llmScenarioResponse) (llmScenarioReview, error) {
	scenarioJSON, err := json.Marshal(scenario)
	if err != nil {
		return llmScenarioReview{}, fmt.Errorf("marshal scenario for review: %w", err)
	}

	systemPrompt := `Ты проверяешь внутреннюю логическую согласованность уже готового JSON детективного сценария. Не изменяй сценарий и не возвращай его заново. Верни ТОЛЬКО JSON формата {"issues":["конкретная проблема"]}.

Включай проблему только если она прямо подтверждается самим JSON: противоречие в мотивации существенного действия, разрыв или противоречие в пути существующей улики, несовпадение финального состояния существующей улики с последним событием о ней, либо лабораторный результат в evidence. Не проверяй, не потерялись ли события, персонажи или предметы по отношению к какому-либо внешнему плану; не требуй добавлять события или улики. Не отмечай стиль, полноту художественного описания или детали, которые не следуют однозначно из JSON. Если подтверждённых проблем нет, верни {"issues":[]}.

Каждое замечание должно быть исправимо изменением уже существующих полей, без добавления или удаления персонажей, событий и улик. Не предлагай новую сюжетную линию. Максимум пять коротких, адресных замечаний с указанием поля или события, которое нужно исправить.`

	messages := []chatMessage{
		{Role: "system", Content: systemPrompt + languageInstruction(locale)},
		{Role: "user", Content: "СЦЕНАРИЙ ДЛЯ ПРОВЕРКИ:\n" + string(scenarioJSON)},
	}
	content, err := c.chat(ctx, "scenario.review", messages, true)
	if err != nil {
		return llmScenarioReview{}, fmt.Errorf("review scenario: %w", err)
	}

	var review llmScenarioReview
	if err := json.Unmarshal([]byte(content), &review); err != nil {
		slog.WarnContext(ctx, "scenario review parse failed; retrying", "error", err, "response_preview", truncateForError(content))
		retryContent, retryErr := c.chat(ctx, "scenario.review_retry", messages, true)
		if retryErr != nil {
			return llmScenarioReview{}, fmt.Errorf("parse scenario review (retry failed): %w", err)
		}
		if retryParseErr := json.Unmarshal([]byte(retryContent), &review); retryParseErr != nil {
			return llmScenarioReview{}, fmt.Errorf("parse scenario review after retry: %w; response: %s", retryParseErr, truncateForError(retryContent))
		}
	}
	return review, nil
}

func (c *OpenRouterClient) repairScenario(ctx context.Context, locale domain.Locale, scenario llmScenarioResponse, issues []string) (llmScenarioResponse, error) {
	scenarioJSON, err := json.Marshal(scenario)
	if err != nil {
		return llmScenarioResponse{}, fmt.Errorf("marshal scenario for repair: %w", err)
	}
	issuesJSON, err := json.Marshal(issues)
	if err != nil {
		return llmScenarioResponse{}, fmt.Errorf("marshal scenario review issues: %w", err)
	}

	systemPrompt := `Ты точечно исправляешь JSON детективного сценария по списку подтверждённых замечаний. Верни ТОЛЬКО полный исправленный JSON точно той же схемы, что и исходный JSON.

Исходный JSON — единственный источник истины. Исправь только перечисленные проблемы, сохрани все остальные значения, персонажей, id, порядок и количество событий, а также порядок, количество и названия улик без изменений. Не добавляй и не удаляй персонажей, улики, события или новые сюжетные факты. Не изменяй opinions, secrets и memories: они должны остаться пустыми массивами. Перед ответом проверь корректность JSON.`
	messages := []chatMessage{
		{Role: "system", Content: systemPrompt + languageInstruction(locale)},
		{Role: "user", Content: "ИСХОДНЫЙ JSON:\n" + string(scenarioJSON) + "\n\nЗАМЕЧАНИЯ ДЛЯ ИСПРАВЛЕНИЯ:\n" + string(issuesJSON)},
	}
	content, err := c.chat(ctx, "scenario.repair", messages, true)
	if err != nil {
		return llmScenarioResponse{}, fmt.Errorf("repair scenario: %w", err)
	}

	var repaired llmScenarioResponse
	if err := json.Unmarshal([]byte(content), &repaired); err != nil {
		slog.WarnContext(ctx, "repaired scenario parse failed; retrying", "error", err, "response_preview", truncateForError(content))
		retryContent, retryErr := c.chat(ctx, "scenario.repair_retry", messages, true)
		if retryErr != nil {
			return llmScenarioResponse{}, fmt.Errorf("parse repaired scenario (retry failed): %w", err)
		}
		if retryParseErr := json.Unmarshal([]byte(retryContent), &repaired); retryParseErr != nil {
			return llmScenarioResponse{}, fmt.Errorf("parse repaired scenario after retry: %w; response: %s", retryParseErr, truncateForError(retryContent))
		}
	}
	return repaired, nil
}

func validateScenario(scenario llmScenarioResponse) error {
	if len(scenario.Characters) != 5 {
		return fmt.Errorf("expected 5 characters, got %d", len(scenario.Characters))
	}
	if len(scenario.Evidence) != 5 {
		return fmt.Errorf("expected 5 evidence items, got %d", len(scenario.Evidence))
	}
	if len(scenario.Timeline) < 14 || len(scenario.Timeline) > 20 {
		return fmt.Errorf("expected 14-20 timeline entries, got %d", len(scenario.Timeline))
	}

	characterIDs := make(map[int]struct{}, len(scenario.Characters))
	characterNames := make(map[string]struct{}, len(scenario.Characters))
	for _, character := range scenario.Characters {
		if character.ID < 1 || character.ID > 5 {
			return fmt.Errorf("character id %d is outside 1-5", character.ID)
		}
		if _, exists := characterIDs[character.ID]; exists {
			return fmt.Errorf("duplicate character id %d", character.ID)
		}
		characterIDs[character.ID] = struct{}{}

		name := strings.ToLower(strings.TrimSpace(character.Name))
		if name == "" {
			return fmt.Errorf("character %d has empty name", character.ID)
		}
		if _, exists := characterNames[name]; exists {
			return fmt.Errorf("duplicate character name %q", character.Name)
		}
		characterNames[name] = struct{}{}
	}
	if _, exists := characterIDs[scenario.Crime.PerpetratorID]; !exists {
		return fmt.Errorf("perpetrator id %d does not reference a character", scenario.Crime.PerpetratorID)
	}
	if _, exists := characterNames[strings.ToLower(strings.TrimSpace(scenario.Crime.Victim))]; exists {
		return fmt.Errorf("victim %q is also present in characters", scenario.Crime.Victim)
	}

	previousMinutes := -1
	for i, entry := range scenario.Timeline {
		minutes, err := parseClock(entry.Time)
		if err != nil {
			return fmt.Errorf("timeline entry %d: %w", i+1, err)
		}
		if minutes < previousMinutes {
			return fmt.Errorf("timeline entry %d at %s is out of order", i+1, entry.Time)
		}
		previousMinutes = minutes
		if strings.TrimSpace(entry.Event) == "" {
			return fmt.Errorf("timeline entry %d has empty event", i+1)
		}
		if entry.CharacterID != nil {
			if _, exists := characterIDs[*entry.CharacterID]; !exists {
				return fmt.Errorf("timeline entry %d references unknown character id %d", i+1, *entry.CharacterID)
			}
		}
	}

	evidenceNames := make(map[string]struct{}, len(scenario.Evidence))
	for i, evidence := range scenario.Evidence {
		name := strings.ToLower(strings.TrimSpace(evidence.Name))
		if name == "" {
			return fmt.Errorf("evidence %d has empty name", i+1)
		}
		if _, exists := evidenceNames[name]; exists {
			return fmt.Errorf("duplicate evidence name %q", evidence.Name)
		}
		evidenceNames[name] = struct{}{}
		if strings.TrimSpace(evidence.Description) == "" || strings.TrimSpace(evidence.DetailedDescription) == "" {
			return fmt.Errorf("evidence %q has an empty description", evidence.Name)
		}
	}
	return nil
}

func parseClock(value string) (int, error) {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time %q", value)
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, fmt.Errorf("invalid time %q", value)
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid time %q", value)
	}
	return hour*60 + minute, nil
}

func (c *OpenRouterClient) generatePrivateCharacterData(ctx context.Context, locale domain.Locale, scenario llmScenarioResponse) (map[int]llmPrivateCharacter, error) {
	contextJSON, err := json.Marshal(scenario)
	if err != nil {
		return nil, fmt.Errorf("marshal private character context: %w", err)
	}

	systemPrompt := `Ты — генератор внутреннего состояния персонажей детективной игры.
Тебе передан уже готовый объективный сценарий: crime, timeline, characters и evidence.
Верни ТОЛЬКО JSON следующего формата:

{
  "characters": [
    {
      "id": 1,
      "opinions": ["мнение или предположение персонажа"],
      "secrets": ["то, что персонаж намеренно скрывает"],
      "memories": [
        {"content": "субъективное воспоминание", "timestamp": "21:00"}
      ]
    }
  ]
}

Правила:
- Верни ровно одного результата для каждого из пяти characters, сохранив их id.
- Сгенерируй 8-12 воспоминаний для каждого персонажа. Они должны охватывать
  события до преступления, момент преступления и события после него.
- Воспоминание описывает только личную перспективу: что персонаж видел, слышал,
  делал или узнал. Не выдавай объективный комментарий о его достоверности.
- Все воспоминания, мнения и секреты должны быть основаны на timeline. Не добавляй
  действий, предметов, мест или мотивов, которых нет в timeline.
- Учитывай, что персонаж может неправильно интерпретировать увиденное или иметь
  ошибочное предположение. Не помечай это как «ложное» — просто формулируй его
  как мнение или личное восприятие.
- Сгенерируй 2-5 opinions и 2-4 secrets для каждого персонажа.
- Для персонажа, чей id равен crime.perpetrator_id, обязательно включи в memories
  его собственные воспоминания о подготовке, совершении преступления и сокрытии
  следов. В secrets обязательно укажи факт совершённого убийства и важные действия
  по сокрытию. Он знает, что виновен, но не обязан признаваться детективу.
- Для остальных персонажей не добавляй знаний, которых они не могли получить из
  личного опыта, отношений или событий timeline.
- Не добавляй поля id или source в memories.
- Проверь, что JSON корректно закрыт.`

	systemPrompt += languageInstruction(locale)

	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "Сгенерируй приватные данные персонажей для этого сценария:\n" + string(contextJSON)},
	}
	content, err := c.chat(ctx, "scenario.private_characters", messages, true)
	if err != nil {
		return nil, err
	}

	var response llmPrivateResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		slog.WarnContext(ctx, "private character data parse failed; retrying", "error", err, "response_preview", truncateForError(content))
		retryContent, retryErr := c.chat(ctx, "scenario.private_characters_retry", messages, true)
		if retryErr != nil {
			return nil, fmt.Errorf("parse private character data (retry failed): %w; response: %s", err, truncateForError(content))
		}
		if retryParseErr := json.Unmarshal([]byte(retryContent), &response); retryParseErr != nil {
			return nil, fmt.Errorf("parse private character data after retry: %w; response: %s", retryParseErr, truncateForError(retryContent))
		}
	}

	byID := make(map[int]llmPrivateCharacter, len(response.Characters))
	for _, character := range response.Characters {
		byID[character.ID] = character
	}
	if len(byID) != len(scenario.Characters) {
		return nil, fmt.Errorf("private character data contains %d characters, expected %d", len(byID), len(scenario.Characters))
	}
	return byID, nil
}

func (c *OpenRouterClient) generateCaseBrief(ctx context.Context, locale domain.Locale, chars []domain.Character, crime llmCrime, evidence []llmEvidence, timeline []llmTimelineEntry) (string, string) {
	var charList strings.Builder
	for _, ch := range chars {
		charList.WriteString(fmt.Sprintf("- %s, %d лет, %s\n", ch.Name, ch.Age, ch.Profession))
	}

	var evList strings.Builder
	for _, ev := range evidence {
		evList.WriteString(fmt.Sprintf("- %s: %s\n", ev.Name, ev.Description))
	}
	var timelineList strings.Builder
	for _, entry := range timeline {
		timelineList.WriteString(fmt.Sprintf("- %s — %s\n", entry.Time, entry.Event))
	}

	caseNumber := rand.Intn(9000) + 1000

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
Объективная хронология:
%s

КРИТИЧЕСКИ ВАЖНО: хронология — источник истины для проверки согласованности,
но она содержит внутренние события, неизвестные следствию на старте. Не переноси
в документ скрытые действия персонажей, их намерения, мотивы, алиби, попытки
скрыть следы или создать ложную версию. Не раскрывай виновного и причинную
цепочку преступления.

Документ должен описывать только информацию, доступную следствию при первом
осмотре: место и время, тело, общие обстоятельства происшествия и физические
объекты, действительно обнаруженные на месте. Не добавляй факты, которых нет
в контексте, и не противоречь хронологии.

Верни JSON:
{
  "case_name": "название дела (короткое, атмосферное)",
  "case_brief": "Markdown-документ по формату ниже"
}

Формат case_brief:

# Дело №%d

**Место преступления:** [используй только локацию из хронологии; не придумывай адрес]

**Предполагаемое время:** %s

**Тип преступления:** Убийство

## Обстоятельства

[2-3 абзаца: только обстоятельства, которые следствие может установить по
хронологии и списку улик; не добавляй новые факты]

## Список подозреваемых

[маркированный список, по одной строке на каждого. Формат: - **Имя**, возраст,
профессия — нейтральная связь с жертвой или местом дела]

Связь должна быть общей и не раскрывать роль персонажа в происшествии: должность,
родство, знакомство, деловые или личные отношения, присутствие на месте работы
или иная публично известная связь. Не описывай конкретные действия персонажа,
обращение с предметами, перемещение объектов, создание следов, обнаружение тела
или другие события из внутренней хронологии.

## Улики с места преступления

[маркированный список улик из контекста]

---
*Дело передано:* [вымышленное имя старшего следователя]

*Дата открытия:* [сегодняшняя дата]`,
		crime.Victim, crime.TimeOfCrime, crime.Method, crime.Motive,
		charList.String(), evList.String(), timelineList.String(), caseNumber, crime.TimeOfCrime,
	)

	content, err := c.chat(ctx, "scenario.case_brief", []chatMessage{
		{Role: "system", Content: "Ты — помощник детектива. Генерируешь официальные документы." + languageInstruction(locale)},
		{Role: "user", Content: briefPrompt},
	}, true)
	if err != nil {
		if locale == domain.LocaleEN {
			return fmt.Sprintf("Case #%d", caseNumber), "_document was not generated_"
		}
		return fmt.Sprintf("Дело №%d", caseNumber), "_документ не сгенерирован_"
	}

	var resp struct {
		CaseName  string `json:"case_name"`
		CaseBrief string `json:"case_brief"`
	}
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		slog.WarnContext(ctx, "case brief parse failed", "error", err, "response_preview", truncateForError(content))
		return fmt.Sprintf("Дело №%d", caseNumber), content
	}

	return resp.CaseName, resp.CaseBrief
}

func formatMemories(memories []domain.Memory) string {
	if len(memories) == 0 {
		return "нет записей"
	}

	var b strings.Builder
	for _, memory := range memories {
		fmt.Fprintf(&b, "- [%s] %s\n", memory.Timestamp, memory.Content)
	}
	return b.String()
}

func (c *OpenRouterClient) RespondInInterrogation(ctx context.Context, locale domain.Locale, character domain.Character, playerMessage string) (*ports.LlmInterrogationResponse, error) {
	systemPrompt := fmt.Sprintf(`Ты — персонаж в детективной игре. Отвечай в соответствии с характером и знаниями.

Имя: %s
Возраст: %d
Профессия: %s
Характер: %s
Текущий уровень доверия к следователю: %d из 100.

Мнения и предположения: %s
Секреты: %s
Воспоминания: %s

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
		strings.Join(character.Opinions, "; "),
		strings.Join(character.Secrets, "; "),
		formatMemories(character.Memories),
	)

	systemPrompt += languageInstruction(locale)

	content, err := c.chat(ctx, "interrogation.response", []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf(`Ответь на реплику следователя. Верни ТОЛЬКО JSON:

{
  "answer": "твой ответ (1-3 предложения)",
  "attitude_delta": -5..5 (изменение доверия),
  "statements": ["краткое значимое утверждение из ответа"]
}

Поле statements предназначено для блокнота следователя. Добавляй в него только новые, конкретные и потенциально значимые для расследования утверждения подозреваемого: о действиях, наблюдениях, времени, месте, участниках или связях между ними. Утверждение не обязательно является правдой.
- Если в ответе нет такого утверждения (например, это приветствие, эмоция, угроза, уклончивый ответ, общее отрицание или оценочное суждение), верни пустой массив [].
- Не включай в statements содержание вопроса следователя, эмоции, догадки без фактической основы и повтор уже сказанного в этом же ответе.
- Если ответ содержит несколько независимых утверждений, верни их отдельными короткими записями.
- Формулируй записи нейтрально и самодостаточно, без имени персонажа и без вводных слов вроде «подозреваемый утверждает».

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
		return nil, fmt.Errorf("parse interrogation response: %w; response: %s", err, truncateForError(content))
	}

	return &ports.LlmInterrogationResponse{
		Answer:        resp.Answer,
		AttitudeDelta: resp.AttitudeDelta,
		Statements:    resp.Statements,
	}, nil
}

func (c *OpenRouterClient) EvaluateReport(ctx context.Context, locale domain.Locale, playerReport domain.FinalReport, groundTruth domain.Crime) (*ports.LlmFeedbackResponse, error) {
	systemPrompt := fmt.Sprintf(`Ты оцениваешь финальный отчёт детектива. Вот что известно об убийстве на самом деле:

Жертва: %s
Преступник: %s
Мотив преступника: %s
Способ убийства: %s
Время преступления: %s

Сравни отчёт игрока с этими фактами и верни ТОЛЬКО JSON:

{
  "narrative_feedback": "развёрнутый текстовый отзыв (3-5 предложений на русском). Оцени общее качество расследования.",
  "breakdown_details": {
    "who": {"correct": true, "comment": "краткая оценка выбора преступника"},
    "why": {"correct": true, "comment": "краткая оценка мотива"},
    "how": {"correct": true, "comment": "краткая оценка способа"},
    "when": {"correct": true, "comment": "краткая оценка времени"},
    "evidence": {"correct": true, "comment": "краткая оценка улик"}
  },
  "missed_facts": ["какие важные детали игрок упустил"]
}

Для каждого поля correct укажи, соответствует ли часть отчёта реальным фактам. Для who допустимы имя без фамилии, обычный порядок имени и фамилии и небольшая очевидная опечатка. Каждый comment обязан соответствовать своему correct.`,
		groundTruth.Victim,
		groundTruth.PerpetratorName,
		groundTruth.Motive,
		groundTruth.Method,
		groundTruth.TimeOfCrime,
	)

	systemPrompt += languageInstruction(locale)

	content, err := c.chat(ctx, "report.evaluation", []chatMessage{
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

	type feedbackDetail struct {
		Correct bool   `json:"correct"`
		Comment string `json:"comment"`
	}
	type feedbackBreakdownDetails struct {
		Who      feedbackDetail `json:"who"`
		Why      feedbackDetail `json:"why"`
		How      feedbackDetail `json:"how"`
		When     feedbackDetail `json:"when"`
		Evidence feedbackDetail `json:"evidence"`
	}
	var resp struct {
		NarrativeFeedback string                   `json:"narrative_feedback"`
		BreakdownDetails  feedbackBreakdownDetails `json:"breakdown_details"`
		MissedFacts       []string                 `json:"missed_facts"`
	}
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return nil, fmt.Errorf("parse feedback: %w; response: %s", err, truncateForError(content))
	}

	return &ports.LlmFeedbackResponse{
		NarrativeFeedback: resp.NarrativeFeedback,
		Breakdown: domain.ScoreBreakdown{
			WhoCorrect:      resp.BreakdownDetails.Who.Correct,
			WhyCorrect:      resp.BreakdownDetails.Why.Correct,
			HowCorrect:      resp.BreakdownDetails.How.Correct,
			WhenCorrect:     resp.BreakdownDetails.When.Correct,
			EvidenceCorrect: resp.BreakdownDetails.Evidence.Correct,
		},
		BreakdownDetails: domain.ScoreBreakdownDetails{
			Who:      resp.BreakdownDetails.Who.Comment,
			Why:      resp.BreakdownDetails.Why.Comment,
			How:      resp.BreakdownDetails.How.Comment,
			When:     resp.BreakdownDetails.When.Comment,
			Evidence: resp.BreakdownDetails.Evidence.Comment,
		},
		MissedFacts: resp.MissedFacts,
	}, nil
}

func formatTimeline(timeline domain.Timeline) string {
	if len(timeline.Entries) == 0 {
		return "нет записей"
	}

	var b strings.Builder
	for _, entry := range timeline.Entries {
		fmt.Fprintf(&b, "- %s: %s\n", entry.Time, entry.Event)
	}
	return b.String()
}

func timelineForSelection(timeline domain.Timeline) string {
	var b strings.Builder
	b.WriteString("{\n")
	for i, entry := range timeline.Entries {
		encoded, _ := json.Marshal(fmt.Sprintf("%s — %s", entry.Time, entry.Event))
		fmt.Fprintf(&b, "  \"%d\": %s", i+1, encoded)
		if i < len(timeline.Entries)-1 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
	}
	b.WriteByte('}')
	return b.String()
}

func (c *OpenRouterClient) selectActionTimelineEntries(ctx context.Context, locale domain.Locale, actionName, requestContext string, timeline domain.Timeline) ([]domain.TimelineEntry, error) {
	content, err := c.chat(ctx, "action.timeline_selection", []chatMessage{
		{Role: "system", Content: "Выбирай только объективные записи таймлайна, относящиеся к запросу. Для предмета включи всю релевантную цепочку: создание или появление, контакты людей, использование, перемещения и финальное местоположение. Для алиби включи события проверяемого человека и события, способные независимо подтвердить или опровергнуть его слова. Не добавляй факты и не меняй тексты. Верни ТОЛЬКО JSON: {\"timeline_entry_ids\":[\"1\"]}. Укажи от 0 до 8 существующих ключей. Если оснований нет, верни пустой массив." + languageInstruction(locale)},
		{Role: "user", Content: fmt.Sprintf("Действие: %s. Контекст запроса: %s.\nТаймлайн: %s", actionName, requestContext, timelineForSelection(timeline))},
	}, true)
	if err != nil {
		return nil, err
	}
	var response struct {
		TimelineEntryIDs []string `json:"timeline_entry_ids"`
	}
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("parse timeline selection: %w", err)
	}
	entries := make([]domain.TimelineEntry, 0, len(response.TimelineEntryIDs))
	seen := make(map[string]struct{})
	for _, id := range response.TimelineEntryIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		index, err := strconv.Atoi(id)
		if err != nil || index < 1 || index > len(timeline.Entries) {
			continue
		}
		seen[id] = struct{}{}
		entries = append(entries, timeline.Entries[index-1])
		if len(entries) == 8 {
			break
		}
	}
	return entries, nil
}

func (c *OpenRouterClient) RunAction(ctx context.Context, locale domain.Locale, actionName string, crime domain.Crime, timeline domain.Timeline, evidence *domain.Evidence, character *domain.Character, alibiText *string) (string, error) {
	var contextParts []string
	if evidence != nil {
		contextParts = append(contextParts, fmt.Sprintf("улика: %s — %s", evidence.Name, evidence.DetailedDescription))
	}
	if character != nil {
		contextParts = append(contextParts, fmt.Sprintf("персонаж: %s, %s", character.Name, character.Profession))
	}
	if alibiText != nil {
		contextParts = append(contextParts, fmt.Sprintf("текст алиби: %s", *alibiText))
	}
	requestContext := strings.Join(contextParts, ", ")
	relevantEntries, err := c.selectActionTimelineEntries(ctx, locale, actionName, requestContext, timeline)
	if err != nil {
		return "", err
	}
	relevantTimeline := formatTimeline(domain.Timeline{Entries: relevantEntries})

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

	content, err := c.chat(ctx, "action.report", []chatMessage{
		{Role: "system", Content: `Ты составляешь итоговый отчёт уже выполненной экспертизы или проверки для детектива. Входные объективные события — скрытая внутренняя основа результата, недоступная детективу.

Правила ответа:
- Пиши только итог исследования в 2-3 предложениях, как законченный официальный результат.
- Никогда не упоминай таймлайн, переданные события, входные данные, скрытый контекст или процесс выбора информации. Не используй формулировки «по переданным записям», «согласно таймлайну», «из представленных данных» и похожие. Запись камеры или история звонков могут упоминаться только как непосредственный объект соответствующей проверки.
- Не пересказывай, кто и когда физически касался объекта, если запрос был лабораторным анализом. Сообщай непосредственно обнаруженные профили ДНК, отпечатки, вещества или иные результаты.
- Не пиши о том, что анализ «может» или «мог бы» что-либо обнаружить: действие уже выполнено, поэтому дай фактический результат.
- Не упоминай внутренние идентификаторы, UUID или технические поля.
- Не придумывай следы, совпадения, звонки, транзакции или участников, для которых нет объективного основания во входных событиях.
- Называй конкретного человека только при наличии явного объективного основания для соответствующего результата.
- Не превращай отсутствие события во входных данных в доказательство отсутствия следов.
- Если объективной основы недостаточно, сформулируй это как результат выполненной экспертизы: например, что пригодный для идентификации материал не выделен или установить обстоятельство не удалось. Не объясняй это отсутствием информации во входном контексте.` + languageInstruction(locale)},
		{Role: "user", Content: fmt.Sprintf("Выполненное действие: %s.\nОбъект исследования: %s\nСкрытая объективная основа результата:\n%s\nСоставь только текст итогового отчёта.", label, requestContext, relevantTimeline)},
	}, false)
	if err != nil {
		return "", err
	}

	return content, nil
}

func languageInstruction(locale domain.Locale) string {
	if locale == domain.LocaleEN {
		return "\n\nReturn every human-readable value in English."
	}
	return "\n\nВсе текстовые значения возвращай на русском языке."
}
