# Spec: Глобальный профиль "Бизнес-исследование"

## Суть

Новый профиль для глобального CLAUDE.md — бизнес-исследование. Отличие от существующего "Исследование": фокус на бизнес-аналитику, сокращение костов/увеличение прибыли, использование специализированных бизнес-ролей вместо технических.

Код не пишется. Результат — структурированный отчёт с action items.

## Workflow

```
Research → Analysis → Plan → Report → Done
```

### Allowed transitions
```
Research  → Analysis
Analysis  → Plan
Analysis  → Research   (недостаточно данных)
Plan      → Report
Plan      → Analysis   (пересмотр выводов)
Report    → Done
```

## Стадии

### Research — консилиум (5 ролей параллельно)

| Роль                | Ответственность                                    |
|---------------------|----------------------------------------------------|
| `market-researcher` | Анализ рынка, TAM/SAM/SOM, тренды, объём           |
| `competitive-analyst`| Конкуренты, позиционирование, бенчмарки            |
| `business-analyst`  | Бизнес-процессы, юнит-экономика, ROI               |
| `product-manager`   | Продуктовая стратегия, приоритизация, PMF           |
| `architect`         | Техническая реализуемость идей, оценка сложности   |

Все роли запускаются параллельно через Task tool. Каждая роль резолвится в конкретный агент через проектный CLAUDE.md или defaults.

### Analysis — синтез данных

Роль: `business-analyst`
Модель: opus

Вход: результаты всех 5 ролей из Research.
Выход: структурированные выводы, ключевые инсайты, риски, возможности.

### Plan — action items

Роль: `business-analyst`
Модель: opus

Вход: результаты Analysis.
Выход: конкретные action items с приоритетами и ожидаемым ROI.

Формат:
- Что делать
- Приоритет (P0/P1/P2)
- Ожидаемый ROI / impact
- Зависимости

### Report — итоговый отчёт

Роль: `technical-writer`
Модель: haiku

Сохранение: `./swarm-report/research-<slug>-<YYYY-MM-DD>.md`

Содержимое отчёта:
- Название исследования и дата
- Исходный запрос
- Резюме Research (по каждой роли)
- Выводы Analysis
- Action items из Plan
- Риски и ограничения
- Статус: Done

## Agents per stage

| Stage    | Role(s)                                                              | Model  |
|----------|----------------------------------------------------------------------|--------|
| Research | market-researcher, competitive-analyst, business-analyst, product-manager, architect | opus |
| Analysis | business-analyst                                                     | opus   |
| Plan     | business-analyst                                                     | opus   |
| Report   | technical-writer                                                     | haiku  |
| Done     | —                                                                    | —      |

## Автодетект — ключевые слова

Рынок, конкуренты, аналитика, бизнес-модель, монетизация, юнит-экономика, TAM, SAM, SOM, unit economics, market research, бизнес-исследование, анализ рынка, анализ конкурентов, маржинальность, revenue, pricing

## Изменения в глобальном CLAUDE.md

### 1. Новые роли в таблицу ролей

| Роль                 | Назначение                                     |
|----------------------|------------------------------------------------|
| `market-researcher`  | Анализ рынка, объёмы, тренды                   |
| `competitive-analyst`| Конкуренты, бенчмарки, позиционирование        |
| `business-analyst`   | Бизнес-процессы, ROI, юнит-экономика           |
| `product-manager`    | Продуктовая стратегия, приоритеты, PMF          |
| `technical-writer`   | Структурирование и оформление отчётов          |

### 2. Defaults для новых ролей

| Роль                 | Default агент                              |
|----------------------|--------------------------------------------|
| `market-researcher`  | voltagent-research:market-researcher       |
| `competitive-analyst`| voltagent-research:competitive-analyst     |
| `business-analyst`   | voltagent-biz:business-analyst             |
| `product-manager`    | voltagent-biz:product-manager              |
| `technical-writer`   | voltagent-biz:technical-writer             |

### 3. Новая строка в таблице профилей

| Бизнес-исследование | `~/.claude/profiles/business-research.md` | Анализ рынка, конкурентов, бизнес-модели, монетизации |

### 4. Новая строка в автодетекте

- Рынок, конкуренты, аналитика, бизнес-модель, монетизация, юнит-экономика, TAM/SAM/SOM, анализ рынка, анализ конкурентов → **Бизнес-исследование**
