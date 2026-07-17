import '../models/character.dart';
import '../models/game_state.dart';
import '../models/report.dart';
import '../models/scenario.dart';
import 'llm_service.dart';

class MockLlmService implements LlmService {
  static final Map<String, CharacterData> characterById = {
    for (final c in characterPool) c.id: c,
  };

  static final List<CharacterData> characterPool = [
    CharacterData(
      id: 'char_01',
      name: 'Иван Петров',
      age: 55,
      profession: 'дворецкий',
      imagePath: 'assets/characters/ivan_petrov.png',
      personality:
          'Консервативный, преданный семье, скрытный. Говорит медленно, с расстановкой. '
          'Предпочитает отмалчиваться, но если задеть за живое — срывается.',
      audioToneId: 'tone_male_deep',
    ),
    CharacterData(
      id: 'char_02',
      name: 'Елена Соколова',
      age: 42,
      profession: 'домохозяйка',
      imagePath: 'assets/characters/elena_sokolova.png',
      personality:
          'Эмоциональная, вспыльчивая, но ранимая. Говорит быстро, часто перебивает. '
          'Хочет казаться безразличной, но на деле очень переживает.',
      audioToneId: 'tone_female_high',
    ),
    CharacterData(
      id: 'char_03',
      name: 'Майкл Браун',
      age: 48,
      profession: 'деловой партнёр',
      imagePath: 'assets/characters/michael_brown.png',
      personality:
          'Харизматичный, уверенный в себе, умело манипулирует. '
          'Говорит спокойно, с лёгкой усмешкой. Всегда контролирует эмоции.',
      audioToneId: 'tone_male_mid',
    ),
    CharacterData(
      id: 'char_04',
      name: 'Анна Коваль',
      age: 29,
      profession: 'горничная',
      imagePath: 'assets/characters/anna_koval.png',
      personality:
          'Застенчивая, тревожная, боится потерять работу. '
          'Говорит тихо, запинается. Старается быть незаметной, но глаза выдают страх.',
      audioToneId: 'tone_female_soft',
    ),
    CharacterData(
      id: 'char_05',
      name: 'Дмитрий Орлов',
      age: 61,
      profession: 'адвокат',
      imagePath: 'assets/characters/dmitry_orlov.png',
      personality:
          'Циничный, расчётливый, за словом в карман не лезет. '
          'Говорит чётко, рублеными фразами. Привык контролировать ситуацию.',
      audioToneId: 'tone_male_raspy',
    ),
  ];

  /* ─────────────── GROUND TRUTH ─────────────── */

  static final _groundTruth = Crime(
    type: CrimeType.murder,
    victim: 'Роберт Ланг',
    perpetratorId: 'char_03',
    motive: 'Ланг раскрыл хищение средств и собирался обратиться в полицию',
    method: 'отравление цианидом (подмешан в виски)',
    timeOfCrime: '22:15',
  );

  static final _timeline = Timeline(entries: [
    TimelineEntry(time: '19:00', event: 'Ужин в особняке. Все подозреваемые присутствуют.', characterId: null),
    TimelineEntry(time: '20:30', event: 'Ланг и Браун уходят в кабинет для разговора.', characterId: null),
    TimelineEntry(time: '21:00', event: 'Браун покидает кабинет. Ланг остается один.', characterId: 'char_03'),
    TimelineEntry(time: '21:15', event: 'Иван Петров заходит в кабинет с виски.', characterId: 'char_01'),
    TimelineEntry(time: '21:45', event: 'Орлов идёт в кабинет обсудить документы, но дверь заперта.', characterId: 'char_05'),
    TimelineEntry(time: '22:00', event: 'Елена Соколова слышит шум из кабинета.', characterId: 'char_02'),
    TimelineEntry(time: '22:15', event: 'Анна Коваль обнаруживает тело.', characterId: 'char_04'),
  ]);

  static final _evidence = [
    Evidence(
      id: 'ev_01',
      name: 'Бокал с остатками виски',
      description: 'Найден на столе в кабинете',
      iconAsset: 'assets/evidence/glass.png',
      detailedDescription:
          'Хрустальный бокал с остатками виски. Обнаружены следы цианида. '
          'Отпечатки пальцев: Роберт Ланг и Иван Петров.',
      type: EvidenceType.physical,
    ),
    Evidence(
      id: 'ev_02',
      name: 'Пузырёк с цианидом',
      description: 'Найден в мусорной корзине в кабинете',
      iconAsset: 'assets/evidence/vial.png',
      detailedDescription:
          'Стеклянный флакон из-под ингалятора с остатками цианида. '
          'Флакон без отпечатков — стёрты. Производитель — компания, '
          'поставляющая лекарства в аптечную сеть Брауна.',
      type: EvidenceType.physical,
    ),
    Evidence(
      id: 'ev_03',
      name: 'Финансовые документы',
      description: 'Разбросаны на столе',
      iconAsset: 'assets/evidence/documents.png',
      detailedDescription:
          'Бухгалтерские отчёты, указывающие на недостачу £500,000. '
          'Документы были скрыты, но кто-то их нашёл и разложил на столе. '
          'На полях пометки почерком Ланга: "мошенничество", "Браун", "полиция".',
      type: EvidenceType.document,
    ),
    Evidence(
      id: 'ev_04',
      name: 'Запертый ящик стола',
      description: 'Ключ торчит снаружи',
      iconAsset: 'assets/evidence/desk.png',
      detailedDescription:
          'Верхний ящик стола заперт, ключ оставлен в замке. '
          'Внутри — конверт с компроматом на Орлова. '
          'Видимо, Ланг кого-то шантажировал — или готовился к этому.',
      type: EvidenceType.physical,
    ),
    Evidence(
      id: 'ev_05',
      name: 'Записка с угрозой',
      description: 'Найдена под дверью кабинета',
      iconAsset: 'assets/evidence/note.png',
      detailedDescription:
          'Анонимная записка: «Не лезь не в своё дело, а то пожалеешь». '
          'Бумага дешёвая, текст напечатан на принтере. '
          'Анализ чернил ничего не дал.',
      type: EvidenceType.document,
    ),
  ];

  static final _characterStates = [
    CharacterState(
      base: characterById['char_01']!,
      knowledge: CharacterKnowledge(
        knownFacts: [
          'Я принёс виски в кабинет в 21:15, Ланг был жив и раздражён',
          'В 21:00 я видел, как Браун вышел из кабинета очень нервным',
        ],
        partialFacts: [
          'Кажется, Ланг с кем-то ссорился перед смертью — доносились голоса',
        ],
        falseBeliefs: [
          'Я думаю, что Елена могла отравить мужа — у них были проблемы в браке',
        ],
      ),
      secrets: [
        'Я должен Лангу крупную сумму (залез в кассу). Ланг знал, но не увольнял',
        'Я убрал свой бокал из кабинета, чтобы не было вопросов по отпечаткам',
      ],
      relationships: {
        'char_03': 'ненавидит — Браун хотел его уволить',
        'char_04': 'защищает — Анна его племянница',
      },
      memories: [
        Memory(id: 'm_01', content: 'Я видел, как Браун вышел из кабинета в 21:00, злой.', isTrue: true, timestamp: '21:00'),
        Memory(id: 'm_02', content: 'Ланг был жив, когда я заходил с виски в 21:15.', isTrue: true, timestamp: '21:15'),
      ],
      trust: 55,
      interrogationsRemaining: 3,
    ),
    CharacterState(
      base: characterById['char_02']!,
      knowledge: CharacterKnowledge(
        knownFacts: [
          'Я слышала крик из кабинета около 22:00',
          'Муж собирался разводиться со мной',
        ],
        partialFacts: [
          'Мне кажется, у мужа были проблемы с бизнесом',
        ],
        falseBeliefs: [
          'Я уверена, что Иван что-то скрывает — он слишком спокоен',
        ],
      ),
      secrets: [
        'У меня были отношения с Орловым. Ланг узнал и хотел развода',
        'В ночь убийства я видела, как Орлов выходил из кабинета в 21:45',
      ],
      relationships: {
        'char_05': 'любовники — Орлов был её любовником',
        'char_03': 'не доверяет — Браун странно влиял на мужа',
      },
      memories: [
        Memory(id: 'm_03', content: 'В 22:00 я слышала громкий спор, а потом глухой удар.', isTrue: true, timestamp: '22:00'),
        Memory(id: 'm_04', content: 'Муж собирался объявить о разводе на следующей неделе.', isTrue: true, timestamp: 'неделю назад'),
      ],
      trust: 45,
      interrogationsRemaining: 3,
    ),
    CharacterState(
      base: characterById['char_03']!,
      knowledge: CharacterKnowledge(
        knownFacts: [
          'Я был в кабинете с Лангом с 20:30 до 21:00',
          'Ланг обнаружил недостачу в отчётах',
        ],
        partialFacts: [
          'Иван должен Лангу деньги, я предлагал его уволить',
        ],
        falseBeliefs: [
          'Думаю, что записку с угрозой написала Анна — она боялась увольнения',
        ],
      ),
      secrets: [
        'Я отравил Ланга цианидом, подмешав яд в графин с виски в 20:45',
        'Я стёр отпечатки с пузырька и выбросил его в корзину',
      ],
      relationships: {
        'char_01': 'презирает — Иван воровал, Браун хотел его выгнать',
        'char_04': 'подозревает — думает, что Анна могла что-то видеть',
      },
      memories: [
        Memory(id: 'm_05', content: 'В 20:45, пока Ланг ненадолго вышел, я подсыпал яд в графин.', isTrue: false, timestamp: '20:45'),
        Memory(id: 'm_06', content: 'Ланг сам налил себе виски в 20:50 и выпил.', isTrue: true, timestamp: '20:50'),
      ],
      trust: 30,
      interrogationsRemaining: 3,
    ),
    CharacterState(
      base: characterById['char_04']!,
      knowledge: CharacterKnowledge(
        knownFacts: [
          'Я нашла тело в 22:15, когда зашла убрать кабинет',
          'В 21:30 я видела, как Орлов шёл к кабинету',
        ],
        partialFacts: [
          'Слышала, что у Ланга и Соколовой были проблемы',
        ],
        falseBeliefs: [
          'Думаю, что Иван мог отравить — он очень нервничал весь вечер',
        ],
      ),
      secrets: [
        'Я видела, как Иван выносил бокал из кабинета в 21:30',
        'Я боюсь, что меня уволят, если я расскажу слишком много',
      ],
      relationships: {
        'char_01': 'боится — Иван требует, чтобы молчала о бокале',
        'char_05': 'странный — Орлов слишком интересовался её показаниями',
      },
      memories: [
        Memory(id: 'm_07', content: 'В 21:30 я видела Ивана, выходящего из кабинета с пустыми руками.', isTrue: false, timestamp: '21:30'),
        Memory(id: 'm_08', content: 'Тело лежало лицом вниз, на столе стоял пустой бокал.', isTrue: true, timestamp: '22:15'),
      ],
      trust: 40,
      interrogationsRemaining: 3,
    ),
    CharacterState(
      base: characterById['char_05']!,
      knowledge: CharacterKnowledge(
        knownFacts: [
          'Я подходил к кабинету в 21:45, дверь была заперта',
          'У меня были деловые отношения с Лангом',
        ],
        partialFacts: [
          'Ланг кого-то шантажировал, судя по тому, что я нашёл в документах',
        ],
        falseBeliefs: [
          'Уверен, что убийца — Иван. У него был мотив (долг) и возможность',
        ],
      ),
      secrets: [
        'У меня роман с Еленой Соколовой, женой Ланга',
        'Ланг нашёл в кабинете мой конверт с компроматом на сделку с Брауном',
      ],
      relationships: {
        'char_02': 'любовники — роман с Еленой',
        'char_03': 'деловой партнёр — провернули сделку с недвижимостью',
      },
      memories: [
        Memory(id: 'm_09', content: 'Я подошёл к кабинету в 21:45, дверь была заперта изнутри.', isTrue: true, timestamp: '21:45'),
        Memory(id: 'm_10', content: 'Ланг нашёл конверт с документами и вызвал меня на разговор.', isTrue: true, timestamp: '20:00'),
      ],
      trust: 50,
      interrogationsRemaining: 3,
    ),
  ];

  final _interrogationResponses = {
    'char_01': {
      'neutral': LlmInterrogationResponse(
        answer: 'Я уже всё рассказал полиции. В 21:15 я занёс мистеру '
            'Лангу виски, как обычно. Он был жив и здоров. '
            'Больше я ничего не видел.',
        attitudeDelta: -2,
        statements: ['Я заносил виски в кабинет в 21:15'],
      ),
      'aggressive': LlmInterrogationResponse(
        answer: 'Послушайте, я 30 лет работаю в этом доме! '
            'Я не обязан отвечать на ваши подозрения. '
            'Если хотите знать, спросите у Брауна — он последним вышел от него.',
        attitudeDelta: -15,
        statements: ['Браун последним вышел из кабинета'],
      ),
      'kind': LlmInterrogationResponse(
        answer: 'Хорошо, я скажу, но это между нами. '
            'Мистер Браун был очень зол, когда вышел из кабинета. '
            'Я слышал, как он говорил по телефону о каких-то "пропавших деньгах".',
        attitudeDelta: 10,
        statements: ['Браун говорил о пропавших деньгах'],
      ),
    },
  };

  @override
  Future<GameSession> generateScenario(List<CharacterData> selectedCharacters) async {
    await Future.delayed(const Duration(seconds: 1));

    return GameSession(
      id: 'session_${_sessionId++}',
      crime: _groundTruth,
      timeline: _timeline,
      evidence: _evidence,
      characters: _characterStates,
      actionPoints: 5,
      phase: GamePhase.investigating,
    );
  }

  int _sessionId = 1;

  @override
  Future<LlmInterrogationResponse> respondInInterrogation({
    required CharacterState characterState,
    required String playerMessage,
  }) async {
    await Future.delayed(const Duration(milliseconds: 500));

    final msg = playerMessage.toLowerCase();
    final characterId = characterState.base.id;

    final isAggressive = msg.contains('убил') ||
        msg.contains('лжёшь') ||
        msg.contains('вы врёте') ||
        msg.contains('докажите');
    final isKind = msg.contains('пожалуйста') ||
        msg.contains('расскажите') ||
        msg.contains('помогите');

    final responses = _interrogationResponses[characterId];
    if (responses != null) {
      if (isAggressive) {
        return responses['aggressive']!;
      }
      if (isKind) {
        return responses['kind']!;
      }
    }

    return _genericResponse(characterState);
  }

  LlmInterrogationResponse _genericResponse(CharacterState character) {
    final name = character.base.name;
    final trust = character.trust;

    if (trust < 25) {
      return LlmInterrogationResponse(
        answer:
            '$name: Я не хочу больше говорить. Вызовите моего адвоката.',
        attitudeDelta: -5,
        statements: ['Отказался отвечать на вопросы'],
      );
    }
    if (trust < 50) {
      return LlmInterrogationResponse(
        answer:
            '$name: Я уже говорил полиции всё, что знаю. '
            'У меня нет времени на эти игры.',
        attitudeDelta: -3,
        statements: ['Повторил, что уже говорил полиции'],
      );
    }

    return LlmInterrogationResponse(
      answer:
          '$name: Я ничего не скрываю. Задавайте ваши вопросы, '
          'я отвечу, если смогу.',
      attitudeDelta: 0,
      statements: ['Готов отвечать на вопросы'],
    );
  }

  @override
  Future<LlmFeedbackResponse> evaluateReport({
    required FinalReport playerReport,
    required Crime groundTruth,
  }) async {
    await Future.delayed(const Duration(milliseconds: 500));

    final perpName = characterById[groundTruth.perpetratorId]?.name.toLowerCase() ?? '';
    final perpParts = perpName.split(' ');
    final whoCorrect = perpParts.any((part) => playerReport.who.toLowerCase().contains(part));

    final whyWords = groundTruth.motive.toLowerCase().split(RegExp(r'[\s,]+'));
    final whyCorrect = whyWords.any((w) => w.length > 3 && playerReport.why.toLowerCase().contains(w));

    final howWords = groundTruth.method.toLowerCase().split(RegExp(r'[\s,]+'));
    final howCorrect = howWords.any((w) => w.length > 3 && playerReport.how.toLowerCase().contains(w));

    final whenCorrect = playerReport.when.contains(groundTruth.timeOfCrime) ||
        playerReport.when.contains(groundTruth.timeOfCrime.replaceAll(':', ''));

    final evidenceCorrect =
        _evidence.any((e) => playerReport.evidence.toLowerCase().contains(e.name.toLowerCase().substring(0, (e.name.length / 2).ceil())));

    final perpFullName = characterById[groundTruth.perpetratorId]?.name ?? 'преступник';

    return LlmFeedbackResponse(
      narrativeFeedback: _generateFeedbackText(
        whoCorrect, whyCorrect, howCorrect, whenCorrect, evidenceCorrect,
      ),
      breakdownDetails: {
        'who': whoCorrect
            ? 'Верно, преступник — $perpFullName.'
            : 'Нет, убийца — $perpFullName. ${groundTruth.motive}',
        'why': whyCorrect
            ? 'Правильно — ${groundTruth.motive}.'
            : 'На самом деле $perpFullName убил, чтобы скрыть ${groundTruth.motive.toLowerCase()}.',
        'how': howCorrect
            ? 'Да, ${groundTruth.method} — верный способ.'
            : 'Нет, орудие — ${groundTruth.method}.',
        'when': whenCorrect
            ? 'Время указано верно — ${groundTruth.timeOfCrime}.'
            : 'Преступление произошло около ${groundTruth.timeOfCrime}.',
        'evidence': evidenceCorrect
            ? 'Хорошо, ключевые улики учтены.'
            : 'Упущены ${_evidence.map((e) => e.name.toLowerCase()).join(', ')}.',
      },
    );
  }

  String _generateFeedbackText(bool who, bool why, bool how, bool when, bool evidence) {
    final correct = [who, why, how, when, evidence].where((c) => c).length;
    if (correct >= 4) {
      return 'Отличная работа! Вы практически полностью восстановили '
          'картину преступления. Ваш анализ заслуживает похвалы.';
    }
    if (correct >= 2) {
      return 'Неплохо, но есть над чем поработать. '
          'Обратите внимание на время преступления и способ.';
    }
    return 'К сожалению, ваша версия далека от истины. '
        'Попробуйте ещё раз изучить улики и показания.';
  }
}
