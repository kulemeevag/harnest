# Feature ROI Analysis: Harnest Top-5 Features
**Дата:** 2026-05-05
**Версия продукта:** v0.10.0
**Методология:** Weighted scoring (Build Cost, Acquisition, Retention, Monetization, Moat)

---

## Executive Summary

Из 7 известных gaps выбраны 5 фич с наибольшим ROI. Приоритизация учитывает реалии solo/small team на Go CLI, open-source модель (CC-BY-NC-4.0), и взрывной рост рынка AI coding tools (2024-2026).

**Ключевой инсайт:** Harnest находится в уникальной позиции "meta-layer" над AI-инструментами. Максимальный ROI дают фичи, усиливающие эту позицию (team sharing, custom harness), а не углубляющие индивидуальный use case (cost optimization, drift detection).

---

## Приоритизированный список

### #1. Custom Harness Support (Plugin System)

**Business Case:**
Рынок AI coding tools фрагментирован и растёт — новые инструменты появляются ежемесячно (Aider, Continue, Cline, Zed AI, JetBrains AI). Жёсткая привязка к 6 harness-ам создаёт потолок роста. Plugin system делает Harnest future-proof и привлекает community contributions.

**Effort:** Medium (2-3 недели)
- Абстракция `Generator` interface уже есть
- Нужен: спецификация формата плагина (YAML/TOML template), discovery из `~/.harnest/harnesses/`, CLI команда `harnest harness add`
- Не нужен plugin runtime — достаточно template-based generation

**Revenue Potential:** Высокий (indirect)
- Расширяет TAM на пользователей любого AI-инструмента
- Community harness-ы = бесплатный маркетинг (каждый PR = новый сегмент)
- Блокирует fork-и (зачем форкать если можно плагин)

**Risk if NOT built:**
- Каждый новый AI-tool = ручная работа maintainer-а
- Конкурент с plugin-системой быстро обгонит по coverage
- Community не может контрибьютить → проект воспринимается как single-maintainer hobby

**Score: 9.2/10**

---

### #2. Team Sharing / Org Configs

**Business Case:**
Команды из 5+ разработчиков — главный monetization сегмент. Боль: каждый настраивает AI-tools по-своему, нет consistency. "Org config" = shared baseline (какие агенты, какие тиры, какие profiles). Это единственная фича, которая естественно конвертируется в paid tier.

**Effort:** Medium (3-4 недели)
- `harnest org init` → создаёт `.harnest/org.yaml` в repo
- `harnest org sync` → мерджит org config с local overrides
- Inheritance: org → team → project → user
- Хранение: git repo (самый простой вариант), позже — cloud registry

**Revenue Potential:** Прямой
- Freemium модель: solo = free, team (5+) = $15/user/month, org (50+) = enterprise
- Sticky: после внедрения в оргструктуру — switching cost высок
- Upsell: analytics dashboard (кто какие модели использует, cost tracking)

**Risk if NOT built:**
- Cursor Teams, Windsurf Teams уже имеют свои config sharing
- Без team story — Harnest остаётся инструментом для одиночек
- Монетизация без team tier крайне сложна для CLI tool

**Score: 9.0/10**

---

### #3. Project Templates (Starter Configs)

**Business Case:**
Снижение time-to-value с минут до секунд. Новый пользователь: `harnest init --template fullstack-react-node` → готовый config с best practices. Templates = SEO магнит ("harnest template for Next.js"), community growth engine, onboarding accelerator.

**Effort:** Small (1-2 недели)
- Template = JSON/YAML с предзаполненными agents + models + profiles
- `harnest templates list` / `harnest init --template <name>`
- Registry: GitHub repo с templates (community-contributed)
- Hardcoded starter set: 10-15 популярных стеков

**Revenue Potential:** Средний (indirect)
- Ускоряет adoption → больше пользователей → больше team upgrades
- Curated/premium templates для enterprise стеков (возможный paywall)
- Партнёрства: "Official Vercel template", "Official Supabase template"

**Risk if NOT built:**
- Высокий friction при onboarding остаётся
- Конкуренты с "zero-config" experience выигрывают у casual users
- Auto-detect хорош, но templates дают "opinionated best practice" which users trust more

**Score: 8.5/10**

---

### #4. Custom Model Mapping (Beyond Tiers)

**Business Case:**
Power users (целевая аудитория Harnest) хотят: "architect = claude-opus-4, frontend = gpt-4o, test = claude-haiku". Текущая система high/medium/low слишком грубая. Это низко-висящий фрукт, который усиливает core value proposition.

**Effort:** Small (3-5 дней)
- В `config.go` уже есть `Models map[string]string`
- Нужно: убрать validation `high/medium/low` only, разрешить произвольные строки
- `harnest agents set-model architect claude-opus-4` (уже почти работает)
- Добавить `harnest models list` — показать tier→model mapping для текущего harness
- Добавить fallback: если значение не в tier map → использовать as-is

**Revenue Potential:** Низкий (direct), Высокий (retention)
- Не монетизируется напрямую
- Но удерживает power users, которые потом приводят команды (→ team tier)
- Differentiator: ни один конкурент не даёт per-role model control

**Risk if NOT built:**
- Низкий immediate risk — но frustration у power users нарастает
- Feature requests будут множиться
- Простота реализации делает отсутствие необъяснимым

**Score: 8.2/10**

---

### #5. Drift Detection (Config Staleness)

**Business Case:**
Проект эволюционирует: добавляются новые языки, фреймворки, зависимости. Config генерируется один раз и устаревает. Drift detection: `harnest check` → "Your project now uses Prisma but no DB agent is configured. Your React version upgraded from CRA to Next.js but scope still points to old paths."

**Effort:** Medium (2-3 недели)
- Re-run detector, сравнить с текущим config
- Diff: new stacks detected, removed stacks, changed paths
- `harnest check` → human-readable report
- `harnest check --fix` → interactive wizard для обновления
- Optional: git hook integration (`pre-commit`)

**Revenue Potential:** Средний
- Retention mechanism: пользователь возвращается регулярно
- Team tier upsell: org-wide drift monitoring dashboard
- Unique value: никто не делает "AI config health check"

**Risk if NOT built:**
- Configs деградируют тихо → пользователь винит Harnest за плохие результаты
- Churn по причине "настроил и забыл, потом перестало работать"
- Без drift detection — Harnest one-time tool, а не continuous companion

**Score: 7.8/10**

---

## Features NOT in Top-5 (и почему)

### Cross-Machine Sync
**Почему отложено:** Решается через git (org config) + dotfiles. Отдельный sync-сервис = инфра overhead, который не оправдан на текущем этапе. Вернуться после team tier.

### Cost Optimization
**Почему отложено:** Требует интеграцию с billing API каждого провайдера. Complexity высокая, данные быстро устаревают (цены меняются), а user demand не подтверждён — power users обычно знают свои costs.

---

## Рекомендуемый Roadmap

```
v0.11.0  Custom Model Mapping        (1 неделя)   ← quick win, unblocks power users
v0.12.0  Project Templates           (2 недели)   ← onboarding accelerator
v0.13.0  Custom Harness Support      (3 недели)   ← TAM expansion
v0.14.0  Drift Detection             (2 недели)   ← retention
v1.0.0   Team Sharing / Org Configs  (4 недели)   ← monetization unlock
```

**Обоснование порядка:**
1. Custom Model Mapping — минимальный effort, максимальный goodwill у существующих пользователей
2. Templates — снижает friction для новых пользователей перед расширением harness-ов
3. Custom Harness — расширяет TAM, привлекает community до запуска monetization
4. Drift Detection — retention перед тем как начать брать деньги
5. Team Sharing — венчает v1.0 как monetizable product с clear value prop для teams

---

## Monetization Model (рекомендация)

| Tier | Price | Features |
|------|-------|----------|
| Free | $0 | All current features + custom models + templates + drift check |
| Pro | $9/mo | Custom harness authoring, priority template access, advanced drift rules |
| Team | $15/user/mo | Org configs, shared profiles, team analytics, SSO |
| Enterprise | Custom | Self-hosted registry, audit logs, compliance templates |

**Ключевой принцип:** CLI остаётся бесплатным и open-source. Монетизация через collaboration layer и managed registry.

---

## Риски и ограничения анализа

1. **Market data uncertainty** — рынок AI coding tools движется быстро, приоритеты могут сместиться за квартал
2. **Solo maintainer risk** — roadmap рассчитан на 10-12 недель, реалистичен для одного разработчика но без буфера
3. **Competitor response** — Cursor/Windsurf могут встроить config portability нативно, обесценив часть value prop
4. **License tension** — CC-BY-NC-4.0 может отпугивать enterprise (рекомендация: рассмотреть dual licensing для core CLI)

---

**Статус:** Done
