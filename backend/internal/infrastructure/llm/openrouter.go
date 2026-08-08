package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

=== ТАЙМЛАЙН ===
Это ПОЛНАЯ ОБЪЕКТИВНАЯ ХРОНИКА мира, известная только тебе. Это не отчёт
криминалистов и не загадка: в ней нет неизвестных обстоятельств, догадок,
формулировок «было найдено» без объяснения и скрытых причин.

Дай 14-20 событий в строгом хронологическом порядке. Каждое event — ясное,
развёрнутое описание факта: КТО совершил действие, ЧТО именно сделал, ГДЕ,
с каким предметом и к чему это привело. Указывай конкретные локации в тексте.
Не пиши «предмет появился», «улика найдена», «смартфон упал» без действующего
лица и причины.

Таймлайн обязан содержать:
- перемещения и существенные действия всех пяти персонажей;
- подготовку преступления, сам момент преступления, действия преступника
  после него и обнаружение тела конкретным человеком;
- полную причинную цепочку для КАЖДОЙ релевантной улики: кто и когда создал,
  принёс, использовал, переместил, спрятал или оставил предмет; где предмет
  оказался к началу расследования; кто его затем увидел или нашёл;
- все инсценировки, ложные следы и попытки отвлечь подозрения: кто их создал,
  зачем и почему выбранная деталь должна вести к конкретной версии;
- непрерывную, непротиворечивую последовательность. Не меняй владельца,
  местоположение, назначение или автора предмета в других частях JSON.

Таймлайн — источник истины для evidence и memories. Ни одна релевантная улика,
никакое действие из memories или secrets не может появиться, если его нет в
timeline. Исключение только для явно случайной, не связанной с делом улики.

=== УЛИКИ (evidence) ===
Это предметы, найденные НА МЕСТЕ ПРЕСТУПЛЕНИЯ. Только то, что криминалист видит при первом осмотре:
- physical: орудие убийства, одежда со следами, разбитые предметы, следы обуви, волокна ткани, бокалы, пузырьки
- digital: телефон, ноутбук, флешка, диск — САМИ УСТРОЙСТВА. НЕ ИХ СОДЕРЖИМОЕ. Контент устройств — предмет действий игрока (call_history, transaction_check и т.д.)
- document: записки, письма, финансовые отчёты, дневники, контракты

Сначала закончи timeline, затем создавай evidence ТОЛЬКО из предметов и следов,
которые уже упомянуты в timeline. Описание улики обязано совпадать с её
происхождением, владельцем, местом и состоянием в timeline. Не меняй автора
записки, владельца одежды или расположение предмета. Если предмет подброшен
для отвода подозрений, timeline должен объяснять, на кого и почему он указывает.

Примеры того, что НЕЛЬЗЯ класть в evidence:
- ❌ "запись с камеры наблюдения" (это результат camera_review)
- ❌ "расшифровка звонков" (это результат call_history)
- ❌ "отпечатки принадлежат X" (это результат fingerprints)
- ❌ "ДНК жертвы на орудии" (это результат dna_analysis)

=== DETAILED DESCRIPTION ===
Развёрнутое описание улики — как в отчёте криминалиста. 3-5 предложений: где именно нашли, в каком состоянии предмет, что видно невооружённым глазом (цвет, форма, повреждения, следы). НЕ пиши «требуется анализ», «необходимо проверить» — это детектив решит сам. НЕ раскрывай результаты анализов (ДНК, отпечатки, содержимое телефона).

✅ Хорошие примеры:
- "Хрустальный бокал из набора «Леннокс», стоит на прикроватной тумбе справа от кровати. На донышке — мутный белёсый осадок. На внешней стенке — частичный след губной помады кораллового оттенка. Рядом с бокалом — лужица жидкости, уже подсохшая."
- "Смартфон «Самсунг», лежит на полу возле входной двери, экраном вниз. Угол экрана разбит — похоже, падал. Корпус в чехле из крокодиловой кожи с инициалами «А.Р.». Аппарат не реагирует на кнопку включения."
- "Записка на листе жёлтой бумаги формата А5, приколота кухонным ножом к дверце холодильника. Текст написан синей шариковой ручкой, размашистым почерком, с наклоном вправо. Бумага слегка помята, в левом верхнем углу — бурое пятно размером с монету."

❌ Плохие примеры:
- "Бокал с остатками жидкости. Требуется анализ содержимого и снятие отпечатков."
- "Отпечатки пальцев принадлежат Ивану."
- "В телефоне найдена переписка с жертвой."

=== ПРИВАТНЫЕ ДАННЫЕ ПЕРСОНАЖЕЙ ===
На этом этапе НЕ генерируй opinions, secrets и memories. Оставь эти поля
пустыми массивами. Они будут сгенерированы отдельным запросом после проверки
crime, timeline, characters и evidence.

=== ОБЩИЕ ПРАВИЛА ===
- characters.id — порядковый номер (1,2,3,4,5)
- perpetrator_id — номер персонажа-убийцы (совпадает с characters.id)
- timeline[].character_id — номер персонажа или null
- Генерируй 5 УНИКАЛЬНЫХ персонажей с разными именами, возрастами и профессиями
- personality — подробное описание характера и манеры речи (2-3 предложения)
- gender — "male" или "female"
- Все значимые улики ДОЛЖНЫ появляться в timeline с полной причинной цепочкой
- relationships — ключ это номер персонажа, значение — описание отношений
- trust — начальное отношение к детективу (0-100). Определи его по характеру,
  личному опыту и предыстории персонажа. Оно не должно зависеть от того,
  является ли персонаж преступником или невиновным. Не делай преступника
  автоматически более закрытым, нервным или враждебным.
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

	private, err := c.generatePrivateCharacterData(ctx, llmResp)
	if err != nil {
		return nil, fmt.Errorf("generate private character data: %w", err)
	}

	perpetratorID := uuid.Nil
	idByLLMID := make(map[int]uuid.UUID, len(llmResp.Characters))
	domainChars := make([]domain.Character, len(llmResp.Characters))
	for i, lc := range llmResp.Characters {
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

	caseName, caseBrief := c.generateCaseBrief(ctx, domainChars, llmResp.Crime, llmResp.Evidence, llmResp.Timeline)

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

func (c *OpenRouterClient) generatePrivateCharacterData(ctx context.Context, scenario llmScenarioResponse) (map[int]llmPrivateCharacter, error) {
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

	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "Сгенерируй приватные данные персонажей для этого сценария:\n" + string(contextJSON)},
	}
	content, err := c.chat(ctx, messages, true)
	if err != nil {
		return nil, err
	}

	var response llmPrivateResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		retryContent, retryErr := c.chat(ctx, messages, true)
		if retryErr != nil {
			return nil, fmt.Errorf("parse private character data (retry failed): %w\n%s", err, content)
		}
		if retryParseErr := json.Unmarshal([]byte(retryContent), &response); retryParseErr != nil {
			return nil, fmt.Errorf("parse private character data: %w\n%s", err, content)
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

func (c *OpenRouterClient) generateCaseBrief(ctx context.Context, chars []domain.Character, crime llmCrime, evidence []llmEvidence, timeline []llmTimelineEntry) (string, string) {
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

	content, err := c.chat(ctx, []chatMessage{
		{Role: "system", Content: "Ты — помощник детектива. Генерируешь официальные документы."},
		{Role: "user", Content: briefPrompt},
	}, true)
	if err != nil {
		return fmt.Sprintf("Дело №%d", caseNumber), "_документ не сгенерирован_"
	}

	var resp struct {
		CaseName  string `json:"case_name"`
		CaseBrief string `json:"case_brief"`
	}
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
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

func (c *OpenRouterClient) RespondInInterrogation(ctx context.Context, character domain.Character, playerMessage string) (*ports.LlmInterrogationResponse, error) {
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
