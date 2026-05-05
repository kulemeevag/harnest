# Feature 5: Team Sharing (`harnest.yaml` as Source of Truth)

## Summary

A single committed `harnest.yaml` file becomes the declarative source of truth for AI assistant configuration. All harness-specific files (CLAUDE.md, .cursorrules, etc.) become generated artifacts. Teams share one config; individual devs layer personal preferences on top.

---

## 1. File Format (Full YAML Schema)

```yaml
# harnest.yaml — source of truth for AI assistant configuration
# Committed to repo. Generated files are gitignored.

version: 1                    # Schema version for future migrations

project:
  name: "my-app"              # Project display name
  description: "E-commerce platform"  # Optional context for AI assistants

# Explicit stack declaration (optional — overrides auto-detection)
# If omitted, `harnest generate` runs detection automatically
stacks:
  - name: spring-boot
    lang: kotlin
    category: backend
    path: server/
  - name: vue
    lang: typescript
    category: frontend
    path: web/
  - name: docker
    lang: dockerfile
    category: infra
    path: .

# Agent assignments
agents:
  consilium:
    architect: voltagent-lang:java-architect
    frontend: voltagent-lang:vue-expert
    ui: ui-designer
    security: security-kotlin
    devops: devops-orchestrator
    api: api-designer
    diagnostics: kotlin-diagnostics
    test: test-spring
    # Omitted roles are excluded from generated config

  executing:
    - agent: builder-spring-feature
      scope: server/**/*.kt
    - agent: voltagent-lang:vue-expert
      scope: web/**/*.vue
    - agent: voltagent-lang:vue-expert
      scope: web/**/*.ts

  models:
    architect: high
    security: high
    frontend: medium
    ui: medium
    api: medium
    devops: medium
    diagnostics: medium
    test: medium
    default: medium           # Fallback for roles not listed

# Target harnesses to generate (array — supports multi-harness)
harnesses:
  - claude-code
  - cursor

# Design system reference (optional)
design_system: linear

# Profiles configuration (optional — overrides which profiles are active)
profiles:
  enabled:
    - business-feature
    - bug-hunting
    - research
    - refactoring
    - e2e-testing
  # custom profiles can be defined inline or referenced
  custom:
    - name: release-check
      file: .harnest/profiles/release-check.md

# Team-level settings
settings:
  # Auto-detect stacks on generate (false = use only explicit stacks above)
  auto_detect: true
  # Merge strategy for stacks: "replace" (yaml wins) or "merge" (yaml + detected)
  stack_strategy: merge
  # Lock file generation
  lock_file: true
```

### Minimal valid `harnest.yaml`

```yaml
version: 1
agents:
  consilium:
    architect: java-architect
  executing:
    - agent: spring-builder
      scope: src/**/*.kt
harnesses:
  - claude-code
```

---

## 2. Relationship to Existing Config

### Current state (without harnest.yaml)

```
harnest init → wizard → CLAUDE.md (source of truth, committed)
                       → .cursorrules (if selected)
```

Config files ARE the source of truth. Editing them directly is the workflow.

### New state (with harnest.yaml)

```
harnest.yaml (source of truth, committed)
    ↓ harnest generate
CLAUDE.md (generated, gitignored)
.cursorrules (generated, gitignored)
.windsurfrules (generated, gitignored)
```

Config files are DERIVED. Editing them is overwritten on next `generate`.

### Coexistence Rules

| Scenario | Behavior |
|----------|----------|
| `harnest.yaml` exists | `harnest init` refuses: "Use `harnest generate` instead" |
| `harnest.yaml` does NOT exist | `harnest init` works as before (current flow) |
| `harnest generate` without `harnest.yaml` | Error: "No harnest.yaml found. Run `harnest init` for wizard or `harnest export` to create from existing config" |
| Both `harnest.yaml` and manual `CLAUDE.md` exist | Warning on generate: "CLAUDE.md exists and will be overwritten. Backup created at CLAUDE.md.bak" |

---

## 3. Workflow: How a Team Adopts This

### Fresh project (greenfield)

```bash
# 1. Install framework (unchanged)
harnest install

# 2. Generate harnest.yaml with wizard
cd my-project
harnest init --yaml

# Wizard runs as before, but output goes to harnest.yaml instead of CLAUDE.md
# At the end:
#   Created: harnest.yaml
#   Run `harnest generate` to produce config files.

# 3. Generate config files
harnest generate

# Output:
#   Generated: CLAUDE.md (from harnest.yaml)
#   Generated: .cursorrules (from harnest.yaml)
#   Added to .gitignore: CLAUDE.md, .cursorrules

# 4. Commit
git add harnest.yaml .gitignore
git commit -m "Add harnest config"
```

### Existing project (migration)

```bash
# Export current config to harnest.yaml
harnest export

# Output:
#   Exported CLAUDE.md → harnest.yaml
#   Review harnest.yaml and run `harnest generate` to verify.
#
#   Next steps:
#     1. Review harnest.yaml
#     2. Run `harnest generate` to test
#     3. Add CLAUDE.md to .gitignore
#     4. Commit harnest.yaml

# Verify round-trip produces same output
harnest generate --dry-run --diff
# Shows diff between current CLAUDE.md and what would be generated
```

---

## 4. Conflict Resolution

### harnest.yaml vs local overrides

**Rule: harnest.yaml always wins on `harnest generate`.**

Local overrides live in a separate layer (see section 7). If someone edits CLAUDE.md directly:

```bash
$ harnest generate

WARNING: CLAUDE.md has local modifications not in harnest.yaml.
  Modified: consilium.frontend = "my-local-agent" (harnest.yaml says "vue-expert")

Options:
  1. Overwrite (harnest.yaml wins) — default
  2. Abort (keep local changes)
  3. Import changes to harnest.yaml

Select [1]:
```

### Team member A and B edit harnest.yaml concurrently (git merge conflict)

Standard YAML merge conflict resolution. The file is structured to minimize conflicts:
- Each role is on its own line
- Exec agents are array items (append-friendly)
- No auto-generated content that changes on every run

### Harness-specific overrides

Some harnesses need extra content (e.g., Cursor needs `.cursorrules` with specific formatting). Handle via `harnesses` section extension:

```yaml
harnesses:
  - name: claude-code
    extra_sections:
      - file: .harnest/claude-extra.md    # Appended to generated CLAUDE.md
  - name: cursor
    extra_sections:
      - file: .harnest/cursor-extra.md
```

---

## 5. Onboarding Flow

### New developer joins team

```bash
git clone git@github.com:team/project.git
cd project

# Option A: Has harnest installed
harnest generate
# Output:
#   Generated: CLAUDE.md
#   Generated: .cursorrules
#   Ready. Your AI assistant is configured.

# Option B: No harnest installed
# IDE shows CLAUDE.md is gitignored but referenced in docs
# README says: "Run `harnest generate` to set up AI assistant config"
# Or: CI/pre-commit hook runs it automatically
```

### Pre-commit hook (optional)

```bash
# .husky/post-checkout or .git/hooks/post-checkout
#!/bin/sh
if command -v harnest &> /dev/null && [ -f harnest.yaml ]; then
  harnest generate --quiet
fi
```

### IDE integration (future)

VS Code task in `.vscode/tasks.json`:

```json
{
  "label": "Generate AI config",
  "type": "shell",
  "command": "harnest generate",
  "runOptions": {"runOn": "folderOpen"}
}
```

---

## 6. Gitignore Strategy

### Committed (source of truth)

| File | Purpose |
|------|---------|
| `harnest.yaml` | Team config |
| `.harnest/profiles/*.md` | Custom team profiles |
| `.harnest/extras/*.md` | Harness-specific extra content |
| `.harnest-local.yaml` | NOT committed (personal, in .gitignore) |

### Generated (gitignored)

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Generated for Claude Code |
| `.cursorrules` | Generated for Cursor |
| `.windsurfrules` | Generated for Windsurf |
| `AGENTS.md` | Generated for Codex/OpenCode |
| `QWEN.md` | Generated for Qwen Code |
| `opencode.json` | Generated for OpenCode |
| `.harnest-lock` | Lock file with generation metadata |

### Auto-managed .gitignore

`harnest generate` automatically adds generated files to `.gitignore` if not already present:

```gitignore
# Harnest generated (do not edit — run `harnest generate`)
CLAUDE.md
.cursorrules
.windsurfrules
AGENTS.md
QWEN.md
opencode.json
.harnest-lock
.harnest-local.yaml
```

---

## 7. Personal Overrides

### File: `.harnest-local.yaml`

Gitignored. Layered on top of team `harnest.yaml` during generation.

```yaml
# .harnest-local.yaml — personal preferences (gitignored)

# Override specific agents (I prefer my fork of vue-expert)
agents:
  consilium:
    frontend: my-forked-vue-expert

# Override model tiers (I have Opus access, team uses Sonnet)
  models:
    architect: high
    default: high

# Additional harness (I also use Windsurf personally)
harnesses:
  - windsurf

# My preferred design system override
design_system: claude
```

### Merge behavior

```
Final config = deep_merge(harnest.yaml, .harnest-local.yaml)
```

Rules:
- Scalar values: local wins
- Arrays (harnesses, exec agents): union (local adds, doesn't remove)
- Maps (consilium, models): local overrides specific keys, team keys preserved
- `agents.executing`: local can ADD entries, not remove team entries

### CLI to manage local overrides

```bash
# Set a personal override
harnest local set agents.consilium.frontend my-agent
# Creates/updates .harnest-local.yaml

# View effective (merged) config
harnest config show
# Shows final merged result

# View what's local vs team
harnest config diff
# Shows only your local overrides
```

---

## 8. Monetization Model

### Free Tier (open source, always free)

- `harnest init` (wizard, single config)
- `harnest detect`
- `harnest agents list/set`
- `harnest convert`
- `harnest drift` (basic — stack + scope drift only)
- `harnest.yaml` with single harness target
- Up to 1 team member (solo dev)

### Pro Tier ($9/month per repo or $29/month unlimited)

- Multi-harness generation from single `harnest.yaml`
- `harnest drift` full (all 6 drift types + CI mode)
- `.harnest-local.yaml` personal overrides
- `harnest generate --watch` (auto-regenerate on yaml change)
- Priority agent mapping updates
- Custom profiles in `harnest.yaml`

### Team Tier ($49/month per org, up to 20 repos)

- Everything in Pro
- `harnest.yaml` team workflow (git-based sharing)
- Org-wide agent registry (`harnest registry push/pull`)
- Drift CI integration with PR comments
- Usage analytics dashboard
- SSO for CLI auth

### Enforcement

```yaml
# In harnest.yaml — declares required tier
version: 1
license: pro          # "free", "pro", "team"
```

On `harnest generate`:
- If `license: pro` and user not authenticated: generate with watermark comment + warning
- Features degrade gracefully (multi-harness falls back to first only, etc.)
- No hard lock-out: generated files always work, just missing advanced features

Authentication:
```bash
harnest auth login          # Opens browser, gets token
harnest auth status         # Shows current plan
harnest auth logout
```

Token stored in `~/.harnest/credentials.json` (gitignored, never committed).

---

## 9. CLI UX

### New commands

```
harnest generate [dir] [flags]       Generate config files from harnest.yaml
harnest export [dir] [flags]         Export existing config to harnest.yaml
harnest config show [dir]            Show effective (merged) config
harnest config diff [dir]            Show local overrides vs team
harnest local set <key> <value>      Set personal override
harnest local unset <key>            Remove personal override
harnest auth login|status|logout     Manage authentication
```

### `harnest init --yaml`

```
$ harnest init --yaml

Detected stack:
  - spring-boot (kotlin) [server/]
  - vue (typescript) [web/]
  - docker (dockerfile) [.]

── Agent Wizard ──
Found 42 agents on this machine

[Consilium: architect]
  Suggestion: voltagent-lang:java-architect
  Enter=accept, s=skip, ?=search: <enter>

[Consilium: frontend]
  Suggestion: voltagent-lang:vue-expert
  Enter=accept, s=skip, ?=search: <enter>

... (same wizard as before) ...

Target harnesses:
  1) claude-code
  2) cursor
  3) windsurf
  4) codex
  5) opencode
  6) qwen-code

Select (comma-separated) [1]: 1,2

Created: harnest.yaml
Next: run `harnest generate` to produce config files.
```

### `harnest generate`

```
$ harnest generate

Reading harnest.yaml...
  Stacks: 3 (2 explicit + 1 auto-detected)
  Harnesses: claude-code, cursor

Generating:
  CLAUDE.md ........... done (42 lines)
  .cursorrules ........ done (38 lines)

Updated .gitignore (added CLAUDE.md, .cursorrules)

Done. 2 config files generated.
```

### `harnest generate --dry-run`

```
$ harnest generate --dry-run

Would generate:
  CLAUDE.md (42 lines, 1.2 KB)
  .cursorrules (38 lines, 1.0 KB)

No files written. Remove --dry-run to generate.
```

### `harnest generate --diff`

```
$ harnest generate --diff

--- CLAUDE.md (current)
+++ CLAUDE.md (from harnest.yaml)
@@ -5,7 +5,7 @@
 ### Consilium
 | Role | Agent |
 |------|-------|
-| frontend | old-vue-agent |
+| frontend | voltagent-lang:vue-expert |
 | architect | java-architect |
```

### `harnest export`

```
$ harnest export

Reading CLAUDE.md...
  Found: 9 consilium roles, 4 exec agents, 9 model entries
  Detected harness: claude-code

Exporting to harnest.yaml...

Created: harnest.yaml

Review the file and run `harnest generate --diff` to verify round-trip.
Add to .gitignore: CLAUDE.md (now generated)
```

### `harnest config show`

```
$ harnest config show

Effective config (team + local):

project:
  name: my-app

stacks:
  - spring-boot (kotlin) [server/]    # from harnest.yaml
  - vue (typescript) [web/]           # from harnest.yaml
  - docker (dockerfile) [.]           # auto-detected

agents:
  consilium:
    architect: java-architect                 # team
    frontend: my-forked-vue-expert            # LOCAL override
    ui: ui-designer                           # team
    ...

  models:
    architect: high                           # LOCAL override (team: medium)
    default: high                             # LOCAL override (team: medium)

harnesses:
  - claude-code                              # team
  - cursor                                   # team
  - windsurf                                 # LOCAL addition
```

### `harnest local set`

```
$ harnest local set agents.consilium.frontend my-agent

Set in .harnest-local.yaml:
  agents.consilium.frontend = my-agent

Run `harnest generate` to apply.
```

---

## 10. Migration Path

### `harnest export` algorithm

```
1. Find config file (CLAUDE.md / .cursorrules / .windsurfrules)
2. Parse with existing config.ReadProject()
3. Determine harness type from file name
4. Run detector.Detect() for stacks (or parse ## Stack section)
5. Construct YAML structure:
   - version: 1
   - project.name: from first heading or dir name
   - stacks: from ## Stack section or detection
   - agents.consilium: from ### Consilium table
   - agents.executing: from ### Executing table
   - agents.models: from ### Models table (convert model names back to tiers)
   - harnesses: [detected harness type]
6. Write harnest.yaml
7. Print next steps
```

### Handling edge cases in export

| Case | Behavior |
|------|----------|
| Config has custom content outside tables | Exported to `.harnest/extras/claude-extra.md` and referenced in `harnesses[0].extra_sections` |
| Config has `<!-- harnest-managed -->` markers | Only export managed content; warn about unmanaged |
| Multiple config files exist | Ask user which to use as primary, export others as additional harnesses |
| Config references agents not on machine | Export as-is with comment `# WARN: agent not found locally` |
| Models table uses concrete names (opus, sonnet) | Convert back to tiers (high, medium, low) |

### Rollback

If team wants to go back to direct-edit workflow:

```bash
# Remove harnest.yaml, un-gitignore config files
harnest eject

# Output:
#   Removed harnest.yaml
#   Removed generated file markers from .gitignore
#   CLAUDE.md and .cursorrules are now source files (edit directly)
```

---

## 11. Implementation Plan

| Phase | Scope | Effort |
|-------|-------|--------|
| v1 (MVP) | `harnest.yaml` parsing + `harnest generate` for single harness | 3 days |
| v1.1 | `harnest export` from existing config | 2 days |
| v1.2 | Multi-harness generation | 1 day |
| v2 | `.harnest-local.yaml` personal overrides | 2 days |
| v2.1 | `harnest config show/diff` | 1 day |
| v2.2 | Auto .gitignore management | 0.5 day |
| v3 | `harnest init --yaml` (wizard outputs yaml) | 1 day |
| v3.1 | Auth + license enforcement | 3 days |
| v4 | CI hooks, `--watch` mode, PR comments | 3 days |

Total: ~16.5 days engineering effort.

---

## 12. Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Team adoption | 20% of multi-dev projects use harnest.yaml within 6 months | Config file type in telemetry |
| Migration rate | 40% of existing users run `harnest export` after feature launch | Command usage |
| Onboarding time | New dev productive in <2 minutes (clone + generate) | Time from clone to first AI interaction |
| Conflict rate | <5% of `harnest generate` runs hit local modification warnings | Warning frequency |
| Pro conversion | 8% of free users upgrade for multi-harness | Subscription funnel |

---

## 13. Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Users resist gitignoring CLAUDE.md | High — breaks existing workflow | Coexistence mode: harnest.yaml optional, old flow still works |
| YAML syntax errors block team | Medium — broken generate | `harnest validate` command + JSON Schema for IDE autocomplete |
| Lock-in concern | Medium — "what if Harnest dies?" | `harnest eject` always works, generated files are plain text |
| Merge conflicts in harnest.yaml | Low — YAML is merge-friendly | Structured format, one value per line, documented merge strategy |
| Personal overrides diverge too much | Low — defeats team alignment purpose | `harnest config diff --team` shows divergence; team leads can audit |
