# Feature 4: Drift Detection (`harnest drift`)

## Summary

Detects configuration drift between what Harnest *would* generate (based on current project state) and what is *actually* in the config files. Answers the question: "Is my AI assistant config still aligned with my project reality?"

---

## 1. What Exactly Drifts

### Drift Type 1: Stack Drift

The project's actual stack changed but config wasn't regenerated.

Examples:
- Added a new `webApp/` with Vue but no exec agent for `**/*.vue`
- Removed `iosApp/` directory but config still has `iosApp/**/*.swift` exec agent
- Upgraded from React to Next.js but config still says "react"

Detection: Re-run `detector.Detect()` and compare result with stacks listed in config's `## Stack` section.

### Drift Type 2: Agent Availability Drift

Agents referenced in config no longer exist on the machine (uninstalled plugin, deleted agent file).

Examples:
- Config references `voltagent-lang:vue-expert` but plugin was uninstalled
- Agent file `~/.claude/agents/my-architect.md` was deleted
- Plugin cache cleared

Detection: Re-run `agents.Discover()` and check every agent name in Consilium + Executing tables exists in discovered list.

### Drift Type 3: Scope Drift

File scope patterns in Executing table don't match actual project structure.

Examples:
- Config says `backend/**/*.kt` but directory was renamed to `server/`
- Config says `composeApp/**/*.kt` but the dir no longer exists
- New source directories appeared that aren't covered by any exec agent

Detection: For each exec scope glob, check if it matches at least one file in the project. For detected stacks, check if their expected scope is covered.

### Drift Type 4: Mapping Version Drift

Harnest's built-in mapping rules were updated (via `harnest update` or binary upgrade) but project config uses old suggestions.

Examples:
- Harnest v0.11 adds better agent matching for Go projects, but config was generated with v0.10
- New role added to Harnest (e.g., `mobile`) but config doesn't have it

Detection: Store `harnest_version` and `mapping_version` in a metadata comment inside generated config. Compare with current binary version.

### Drift Type 5: Model Tier Drift

Model tier defaults changed in Harnest update, but config uses old concrete model names.

Examples:
- Default for `architect` changed from `sonnet` to `opus` but config still says `sonnet`
- New harness tier mapping added but Models table uses legacy names

Detection: Compare Models table entries against `DefaultModelTiers()` + current `ResolveTier()` output.

### Drift Type 6: Profile Drift

Installed profiles are outdated compared to what current Harnest version ships.

Examples:
- `business-feature.md` profile was updated in Harnest v0.11 but local copy is from v0.9
- User has profiles that reference roles not in their project config

Detection: Compare checksums of `~/.claude/profiles/*.md` against embedded templates in binary.

---

## 2. Detection Algorithm

```
Input: project directory, harness name (auto-detected from existing config file)

Step 1: Read existing config
  - Find config file (CLAUDE.md / .cursorrules / .windsurfrules)
  - Parse: stacks, consilium, exec, models, metadata

Step 2: Compute expected state
  - Run detector.Detect(dir) → expected stacks
  - Run agents.Discover() → available agents
  - Run mapping.Resolve(stacks, discovered, harness) → expected config

Step 3: Diff each drift type
  For each type, produce a DriftItem:
    {Type, Severity, Current, Expected, AutoFixable, Description}

Step 4: Output
  - If no drift: "Config is up to date."
  - If drift found: show table with items grouped by severity
```

### Severity Levels

| Level | Meaning | Action |
|-------|---------|--------|
| `error` | Config references non-existent resources (missing agent, dead scope) | Must fix |
| `warning` | Config is valid but outdated (new stacks untracked, old version) | Should fix |
| `info` | Minor improvements available (better agent match, new tier defaults) | Optional |

---

## 3. Output Format

### Terminal (default)

```
$ harnest drift

Drift detected (3 issues):

 ERROR   Agent not found: "voltagent-lang:vue-expert" (role: frontend)
         Agent no longer installed. Available alternatives: vue-developer, frontend-expert

 WARNING Stack added: nuxt (typescript) [webApp/]
         No exec agent covers webApp/**/*.vue. Run `harnest drift --fix` to add.

 INFO    Mapping updated: architect suggestion changed
         Current: java-architect → New suggestion: spring-architect (better match for spring-boot)

Run `harnest drift --fix` to auto-fix 2 of 3 issues.
Run `harnest drift --fix --interactive` to review each fix.
```

### JSON (for CI)

```
$ harnest drift --json

{
  "version": "0.11.0",
  "harness": "claude-code",
  "config_file": "CLAUDE.md",
  "drifts": [
    {
      "type": "agent_availability",
      "severity": "error",
      "role": "frontend",
      "current": "voltagent-lang:vue-expert",
      "expected": null,
      "auto_fixable": false,
      "message": "Agent not found: voltagent-lang:vue-expert"
    },
    {
      "type": "stack",
      "severity": "warning",
      "stack": "nuxt",
      "current": null,
      "expected": {"agent": "vue-developer", "scope": "webApp/**/*.vue"},
      "auto_fixable": true,
      "message": "New stack detected but not in config"
    }
  ],
  "summary": {
    "errors": 1,
    "warnings": 1,
    "info": 1,
    "auto_fixable": 2
  }
}
```

---

## 4. Auto-Fix Capability

### Auto-fixable (no user decision needed)

| Drift Type | Fix Action |
|------------|------------|
| New stack detected | Add exec agent row with best-match agent from discovered list |
| Dead scope (dir removed) | Remove exec agent row |
| Profile outdated | Overwrite with latest embedded template |
| Mapping version outdated | Update metadata comment |
| Model name normalization | Replace deprecated model name with current tier equivalent |

### Requires interactive decision

| Drift Type | Why Manual |
|------------|-----------|
| Agent not found | User must pick replacement from available agents or install missing plugin |
| Better agent match available | User might prefer current agent intentionally |
| Role added by new Harnest version | User might not want all roles |
| Scope conflict (two agents match) | Ambiguous — user decides priority |

### Fix modes

```bash
harnest drift --fix                 # Auto-fix only auto-fixable items
harnest drift --fix --interactive   # Review each fix before applying
harnest drift --fix --all           # Fix everything (uses best suggestion for manual items)
```

---

## 5. CI Integration

### Command

```bash
harnest drift --ci [--fail-on warning|error]
```

### Behavior in CI mode

- Output is JSON to stdout (parseable by CI tools)
- Exit code 0 = no drift
- Exit code 1 = drift found at or above `--fail-on` level (default: `error`)
- No interactive prompts, no color codes
- Adds annotation format for GitHub Actions (`::error::`, `::warning::`)

### GitHub Actions Example

```yaml
name: Config Drift Check
on:
  push:
    paths:
      - 'package.json'
      - 'build.gradle.kts'
      - 'go.mod'
      - 'Cargo.toml'
      - 'pubspec.yaml'
      - 'requirements.txt'
      - 'pyproject.toml'

jobs:
  drift:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: AlexGladkov/setup-harnest@v1
      - run: harnest drift --ci --fail-on warning
```

### GitLab CI Example

```yaml
drift-check:
  stage: lint
  script:
    - harnest drift --ci --fail-on error
  rules:
    - changes:
        - package.json
        - "**/*.gradle.kts"
        - go.mod
```

### PR Comment (optional future extension)

```bash
harnest drift --ci --github-comment
```

Posts a PR comment with drift summary table using `$GITHUB_TOKEN`.

---

## 6. User Stories

### Story 1: "We added a microservice and forgot to update config"

Team adds `payments-service/` with Go + Chi framework. Two weeks later, AI assistant keeps suggesting wrong patterns because no exec agent covers `payments-service/**/*.go`. Developer runs `harnest drift` and sees:

```
WARNING  Stack added: chi (go) [payments-service/]
         No exec agent covers payments-service/**/*.go
```

Runs `harnest drift --fix` — gets the new row added automatically.

### Story 2: "Plugin author renamed their agents"

Developer uses `voltagent-research` plugin. Author releases v2 that renames agents from `market-researcher` to `voltagent-research:market-analyst`. Config now references dead names.

```
ERROR  Agent not found: "voltagent-research:market-researcher" (role: market-researcher)
       Did you mean: voltagent-research:market-analyst?
```

`--fix --interactive` prompts: "Replace with voltagent-research:market-analyst? [Y/n]"

### Story 3: "CI catches drift on dependency update"

PR updates `package.json`, adding `@remix-run/react`. CI runs `harnest drift --ci --fail-on warning`:

```json
{"type":"stack","severity":"warning","stack":"remix","message":"New framework detected: remix replaces react in webApp/"}
```

CI fails, PR author runs `harnest drift --fix` locally and commits the config update.

### Story 4: "Harnest upgrade brings better defaults"

Developer upgrades Harnest from v0.10 to v0.12. New version has improved agent matching and a new `api` role suggestion. Running `harnest drift`:

```
INFO  Mapping version drift: config generated with v0.10, current is v0.12
INFO  Better agent available for role "architect": spring-kotlin-architect (score: 3) vs java-architect (score: 1)
INFO  New role available: "api" — not in your config
```

All `info` level — optional to act on.

---

## 7. Edge Cases

### Intentional drift (user customized)

**Problem:** User manually edited config (e.g., set `architect → my-private-agent`). Drift detection should NOT flag this as "better agent available."

**Solution:** Metadata marker system.

```markdown
<!-- harnest:generated -->
| architect | java-architect |
<!-- harnest:manual -->
| frontend | my-custom-vue-expert |
```

Alternative (simpler, v1): `.harnest-lock` file that records which entries were user-overridden via `harnest agents set`:

```json
{
  "overrides": {
    "consilium.frontend": "my-custom-vue-expert",
    "exec.0": "my-spring-builder"
  },
  "generated_at": "2026-05-01T10:00:00Z",
  "harnest_version": "0.10.0"
}
```

Drift detection skips entries present in `overrides`. Reported as `info` at most: "Manually overridden, skipping."

### Dead agent that user wants to keep

User references an agent they plan to install later. Use `--ignore` flag:

```bash
harnest drift --ignore "my-future-agent"
```

Or `.harnestignore` file:

```
# Agents I'll install later
my-future-agent
experimental-security-agent
```

### Multi-harness projects

Project has both `CLAUDE.md` and `.cursorrules`. Drift check runs against the primary (auto-detected or specified):

```bash
harnest drift                        # checks primary config
harnest drift --harness cursor       # checks .cursorrules specifically
harnest drift --all-harnesses        # checks all found config files
```

### Monorepo with multiple stacks

Detection finds 8 stacks. Config only has 3 exec agents (user intentionally skipped others during wizard). Lock file records skipped items:

```json
{
  "skipped_stacks": ["docker", "github-actions", "terraform"],
  "skipped_roles": ["mobile", "devops"]
}
```

Drift detection won't flag skipped items.

---

## 8. CLI UX

### Commands

```
harnest drift [dir] [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--harness <name>` | auto-detect | Which config file to check |
| `--fix` | false | Auto-fix auto-fixable issues |
| `--fix --interactive` | false | Review each fix before applying |
| `--fix --all` | false | Fix everything using best suggestions |
| `--ci` | false | CI mode: JSON output, exit codes, no color |
| `--fail-on <level>` | `error` | CI exit code threshold (`error`, `warning`, `info`) |
| `--json` | false | JSON output (without CI exit code behavior) |
| `--ignore <agent>` | — | Skip specific agent from availability check |
| `--all-harnesses` | false | Check all found config files |
| `--verbose` | false | Show detailed comparison for each drift item |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No drift (or drift below `--fail-on` threshold) |
| 1 | Drift detected at or above threshold |
| 2 | Error reading config / running detection |

### Full interaction example

```
$ harnest drift

Checking CLAUDE.md against project state...

  Stacks:  4 detected, 4 in config    OK
  Agents:  12 referenced, 11 found    1 MISSING
  Scopes:  5 patterns, 5 valid        OK
  Models:  9 entries, 0 outdated       OK
  Version: v0.10.0 (config) vs v0.12.0 (binary)  OUTDATED

Drift detected (2 issues):

  ERROR    [agent_availability] Role: frontend
           Config: voltagent-lang:vue-expert
           Status: Agent not found on this machine
           Hint:   Available alternatives matching "vue": vue-developer, nuxt-expert
           Fix:    Interactive — run `harnest drift --fix --interactive`

  WARNING  [mapping_version] Config generated with v0.10.0
           Binary: v0.12.0 (newer mappings available)
           Fix:    Auto — run `harnest drift --fix`

Summary: 1 error, 1 warning, 0 info
Auto-fixable: 1 of 2

$ harnest drift --fix --interactive

[1/2] Agent "voltagent-lang:vue-expert" not found (role: frontend)
  Available matches:
    1) vue-developer (score: 2)
    2) nuxt-expert (score: 1)
    3) [skip — keep current]
    4) [type custom name]
  Select [1]: 1

  Updated: frontend → vue-developer

[2/2] Update mapping version metadata?
  Current: v0.10.0 → New: v0.12.0
  Accept? [Y/n]: Y

  Updated metadata.

Done. Fixed 2 issues in CLAUDE.md.
```

---

## 9. Data Model

### Internal Types

```go
type DriftType string

const (
    DriftStack           DriftType = "stack"
    DriftAgentAvail      DriftType = "agent_availability"
    DriftScope           DriftType = "scope"
    DriftMappingVersion  DriftType = "mapping_version"
    DriftModelTier       DriftType = "model_tier"
    DriftProfile         DriftType = "profile"
)

type Severity string

const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo    Severity = "info"
)

type DriftItem struct {
    Type        DriftType
    Severity    Severity
    Current     string      // what's in config now
    Expected    string      // what should be there (or "" if removal)
    Role        string      // for consilium drift
    Scope       string      // for exec drift
    Stack       string      // for stack drift
    AutoFixable bool
    Message     string
    Hint        string
}

type DriftResult struct {
    ConfigFile string
    Harness    string
    Version    string
    Items      []DriftItem
}
```

### Lock File: `.harnest-lock`

```json
{
  "schema_version": 1,
  "generated_at": "2026-05-05T14:30:00Z",
  "harnest_version": "0.12.0",
  "harness": "claude-code",
  "overrides": {
    "consilium.architect": "my-custom-architect",
    "exec.1.agent": "my-spring-builder"
  },
  "skipped_stacks": ["docker"],
  "skipped_roles": ["mobile"]
}
```

Created/updated by `harnest init` and `harnest agents set`. Read by `harnest drift` to distinguish intentional vs unintentional drift.

---

## 10. Implementation Plan

| Phase | Scope | Effort |
|-------|-------|--------|
| v1 (MVP) | Stack drift + Agent availability + Scope drift | 2 days |
| v1.1 | `--fix` auto-fix for safe items | 1 day |
| v1.2 | `--ci` mode + JSON output + exit codes | 1 day |
| v2 | Lock file + intentional drift tracking | 2 days |
| v2.1 | `--fix --interactive` for manual items | 1 day |
| v3 | Profile drift + Model tier drift + `--all-harnesses` | 2 days |

Total: ~9 days engineering effort.

---

## 11. Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Adoption | 30% of active users run `drift` monthly | CLI telemetry (opt-in) |
| CI integration | 10% of projects add drift check to CI | GitHub search for `harnest drift --ci` |
| Auto-fix rate | >60% of drift items are auto-fixable | Aggregate from `--json` output |
| False positive rate | <5% drift items are intentional | Lock file override frequency |
