# `mdlint` – Technical Specification (Markdown Linter in Go)

**Status:** Draft v1.0
**Language:** Go **1.24.6+**
**Targets:** Linux/macOS/Windows (amd64, arm64)
**Project style:** **Align with** `asymmetric-effort/docker-lint` repository conventions (layout, build, release, CI, config)
**Binary name:** `mdlint`
**Module path (proposed):** `github.com/asymmetric-effort/mdlint`

---

## 1. Purpose & Scope

`mdlint` is an **extensible Markdown linter** focused on **deterministic, fast, single-binary** execution suitable for local use and CI. It provides:

* **Rule engine** with Markdown‑aware AST analysis (paragraphs, headings, links, code spans/blocks, lists, tables) powered by Goldmark.
* **Extensible ruleset** with IDs prefixed **`MD`** (e.g., `MD1000`), modeled after prose/consistency checks commonly used in Vale.
* **Configurable severities** (suggestion, warning, error), ignores, and per‑path rule config.
* **Zero‑network, offline** operation by default (optional extras opt‑in).
* **JSON/tty** outputs and stable **exit codes** for CI.

Non‑goals (v1): auto‑fixes; collaborative review server; rendering HTML/CSS fidelity.

---

## 2. Repository Layout (match `docker-lint` patterns)

```
mdlint/
├─ cmd/
│  └─ mdlint/                 # main package for CLI
├─ internal/
│  ├─ engine/                 # scheduler, walker, rule orchestration
│  ├─ rules/                  # built-in rule impls (MDxxxx)
│  ├─ parser/                 # goldmark AST parsing and node facades
│  ├─ config/                 # load/merge/validate .mdlint.yaml
│  ├─ findings/               # result model, formatting (json, text)
│  └─ sys/                    # os/path/glob, charset, io helpers
├─ docs/                      # docs, rule reference, examples
├─ scripts/                   # helper scripts (fmt, release, etc.)
├─ tests/                     # end-to-end and golden tests
├─ testdata/                  # fixtures (markdown samples, dictionaries)
├─ .goreleaser.yml            # release pipeline (like docker-lint)
├─ Makefile                   # make build/test/lint, mirrors docker-lint targets
├─ SECURITY.md                # report channel and policy
├─ README.md                  # usage and quick start
├─ go.mod / go.sum
└─ LICENSE (MIT)
```

**Makefile targets (parity with docker-lint):** `make clean|lint|test|build`, optional `release` via Goreleaser.

---

## 3. Configuration

* **File:** `.mdlint.yaml` (project root).
* **Merge order:** CLI flags > project `.mdlint.yaml` > user `$XDG_CONFIG_HOME/mdlint/config.yaml` (optional) > built‑ins.
* **Schema (v1):**

```yaml
version: 1
ignored:            # list of MD rule IDs to skip globally
  - MD1100
severity:           # default severity per rule (overrides built-ins)
  MD1000: warning
  MD1100: error
paths:              # per-path overrides (globs)
  "docs/**":
    ignored: [MD1500]
    severity:
      MD1300: suggestion
spell:              # MD1000 options
  lang: en_US
  add_words: [CRSCE, GoLand, yaml]
  reject_words: []
  filters: ["[pP]y.*\\b"]
heading:
  style: atx        # MD1100: atx|setext|consistent
  allow_mixed: false
output:
  format: json      # json|text
  color: auto       # auto|always|never
failure_threshold: warning # suggestion|warning|error
```

---

## 4. CLI

```
mdlint [flags] <path|glob>...

Flags
  -c, --config <file>         Use a specific config file
  -o, --output <format>       json|text (default json)
  --color <mode>              auto|always|never
  --fail-level <sev>          suggestion|warning|error (default from config)
  --rules                     Print enabled rules and exit
  --list                      List all rules supported by this build
  --version                   Print version
  -q, --quiet                 Suppress non-error logs
```

**Exit codes**
`0`: no findings at or above failure threshold
`1`: findings at/above threshold
`2`: usage/config error
`3`: internal error

**Output (JSON example)**

```json
[
  {"rule":"MD1100","severity":"error","message":"Inconsistent heading style (expected atx)","file":"README.md","line":3,"column":1},
  {"rule":"MD1000","severity":"warning","message":"Possible misspelling: 'teh'","file":"README.md","line":42,"column":8}
]
```

---

## 5. Parsing & Scoping

* **Parser:** `github.com/yuin/goldmark` with extensions: tables, strikethrough, definition lists, task lists, footnotes.
* **Scopes:** Block/inline aware rule targeting; exclude code spans/blocks and fenced regions from prose rules by default.
* **Front‑matter:** YAML front‑matter parsed (if present) and **excluded** from prose rules unless a rule opts in.
* **Files:** UTF‑8 expected; fallback via BOM/charset sniffing; binary-safe skip.

---

## 6. Rule Engine

### 6.1 Rule Model

```go
type Rule interface {
    ID() string                 // e.g. "MD1100"
    Name() string               // human label
    DefaultSeverity() Severity  // suggestion|warning|error
    Apply(doc *Document, cfg RuleConfig) ([]Finding, error)
}

type Finding struct {
    Rule, Message, File string
    Severity Severity
    Line, Column int
    // optional: EndLine, EndColumn, Suggestions []string
}

type Severity int // iota: Suggestion, Warning, Error
```

* **Registration:** built‑ins register via `init()` in `internal/rules`, added to a central registry.
* **Execution:** engine walks Goldmark AST once, builds **indexes** (headings, links, words), then runs enabled rules, each consuming indexes (minimize re-walks).
* **Performance:** concurrent per‑file execution; per‑rule execution inside a file is sequential to preserve deterministic ordering (stable sort by position at the end).

### 6.2 Configuration Injection

Each rule receives a typed `RuleConfig` view (merged from defaults + file‑level overrides + path overrides). Unknown keys are rejected in **strict** mode.

---

## 7. Built‑in Rules (MD‑prefix)

The initial ruleset mirrors the **capabilities** of common prose/consistency checks. Rule numbers reserve ranges for categories and include two explicitly requested rules.

| ID         | Name                          | Summary                                                                                   | Inputs / Notes                                                    |
| ---------- | ----------------------------- | ----------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| **MD1000** | **Spelling**                  | Flag tokens not in dictionary; supports **allow**/**reject** lists and regex **filters**. | Uses wordlist + optional Hunspell adapters; ignores code regions. |
| **MD1100** | **Consistent Heading Style**  | Enforce `atx` (`# H1`), `setext` (`===`/`---`), or **consistent** style per doc.          | Checks Goldmark heading nodes; ignores code blocks.               |
| MD1200     | Heading Capitalization        | Sentence/title case policy for headings; configurable acronyms whitelist.                 | Tokenizer + case heuristics.                                      |
| MD1300     | Readability Metric            | Report too‑long sentences / Flesch‑Kincaid thresholds; skip code/links.                   | Sentence splitter + syllable estimator.                           |
| MD1400     | Repetition                    | Detect repeated words (`the the`) and near‑adjacent duplicates.                           | Token stream.                                                     |
| MD1500     | Consistency (Preferred Terms) | Enforce term variants (e.g., "email" vs "e‑mail"); project vocabulary.                    | Trie lookup; project term map.                                    |
| MD1600     | Occurrence Limits             | Control density (e.g., max exclamations per 1000 chars).                                  | Character counters.                                               |
| MD1700     | Sequence                      | Enforce sequences (e.g., punctuation then space) via regex/state machine.                 | Inline text scan.                                                 |
| MD1800     | Conditional                   | Contextual checks (e.g., forbid passive voice in headings only).                          | Scope filters.                                                    |
| MD1900     | Metric                        | Report file‑level metrics (reading time, word count) as suggestions.                      | Aggregation only.                                                 |
| MD1950     | Script Hook (Opt‑in)          | External script rule (disabled by default) for custom org checks.                         | Sandbox; off by default.                                          |

> Additional MD rules can extend this table; IDs remain stable once released.

---

## 8. Spell Checking (MD1000)

* **Backends:**

  * **Wordlist engine (default, pure Go):** hashed dictionary, case‑insensitive; supports add/reject lists and per‑project allowlists.
  * **Hunspell adapter (optional build tag `hunspell`):** enables `.dic/.aff` dictionaries (CGO); disabled by default to keep a pure‑Go build.
* **Filters:** configurable regex list to ignore patterns (e.g., code‑like identifiers).
* **Tokenization:** Unicode word boundaries; skip code spans/blocks, inline code, autolinks, and fenced regions tagged with languages.

---

## 9. Consistent Heading Style (MD1100)

* **Config:** `heading.style: atx|setext|consistent`; `allow_mixed: bool`
* **Behavior:**

  * `atx`: all headings use `#`, `##`, ...
  * `setext`: H1/H2 use `===` / `---`; deeper levels must be `atx` (Markdown spec).
  * `consistent`: infer from first heading, enforce for rest.
* **Edge cases:** skip headings in code blocks or blockquotes if configured to ignore.

---

## 10. Findings & Reporting

* **Formats:**

  * **JSON**: machine‑readable array (one finding per object).
  * **Text**: `FILE:LINE:COL MDxxxx[SEV] message` with optional color.
* **Sorting:** file → line → column → rule → message (stable).
* **Thresholding:** findings below `failure_threshold` do **not** affect exit code.

---

## 11. Performance

* **I/O:** memory‑mapped reads where available; buffered readers elsewhere.
* **Parsing:** single Goldmark parse per file; shared indices (headings, sentences, words).
* **Concurrency:** worker pool over files; CPU‑bound; bounded parallelism to avoid thrash.
* **Caches:** dictionary and vocabulary LRU; compiled regex cache.

---

## 12. Security Considerations

* **No network** access during linting.
* **External scripts (MD1950)** disabled by default; when enabled: explicit allowlist path, **no arguments**, **no environment**, timeout, stdout-only, exit status handling.
* **Sandbox input**: parse only text files; limit max file size (default 5 MiB, configurable).
* **Config loading**: strict YAML decoder; unknown fields error in strict mode; path globs normalized.
* **Reproducibility**: deterministic ordering; no time‑dependent logic in findings.

---

## 13. Testing Strategy

* **Unit tests**: tokenization, filters, dictionary lookups, heading style transitions, per‑path overrides.
* **Golden tests**: markdown fixtures → expected JSON findings; run with `-race`.
* **Property tests**: random whitespace/markdown permutations preserve rule invariants.
* **Performance tests**: 10k‑line documents, large tables.
* **Cross‑platform CI**: Linux/macOS/Windows on Go 1.24.x via GitHub Actions.

---

## 14. Rule Authoring (Extensibility)

* **API:** drop a new file in `internal/rules/mdXXXX_<name>.go` exposing `Rule` and `init()` registration.
* **Config binding:** define a typed struct and register a decoder hook in `internal/config`.
* **Docs:** add markdown under `docs/rules/MDXXXX.md` with examples and configuration.

Template:

```go
func init() { rules.Register(NewMD1234()) }

type md1234 struct{}
func (r md1234) ID() string { return "MD1234" }
func (r md1234) Name() string { return "My Rule" }
func (r md1234) DefaultSeverity() Severity { return Warning }
func (r md1234) Apply(doc *Document, cfg RuleConfig) ([]Finding, error) { /* ... */ }
```

---

## 15. Integration & Tooling

* **Pre‑commit hook:** sample `.pre-commit-config.yaml` shipped.
* **Editor integrations:** VS Code task + problem matcher; JetBrains External Tools example.
* **Containers:** publish `ghcr.io/asymmetric-effort/mdlint:latest` with `ENTRYPOINT ["/mdlint"]` mirroring docker‑lint image usage.

---

## 16. Example

```bash
mdlint --fail-level warning "docs/**/*.md"
```

Output (text):

```
README.md:3:1 MD1100[error] Inconsistent heading style (expected atx)
README.md:42:8 MD1000[warning] Possible misspelling: 'teh'
```

---

## 17. Roadmap (v1 → v1.2)

* **v1.0**: MD1000, MD1100, MD1400, MD1500 shipped; JSON/text output; Windows/macOS/Linux builds; pure‑Go defaults.
* **v1.1**: Add MD1200 (heading capitalization), MD1300 (readability), MD1600 (occurrence), per‑language dictionaries.
* **v1.2**: Optional Hunspell build tag; cached dictionary loaders; basic autofix for trivial spacing issues.

---

## 18. Compliance with `docker-lint` Conventions

* **Output contract**: JSON array of findings; fields `rule`, `message`, `line`, `file` similar to docker-lint.
* **Config**: `.mdlint.yaml` mirrors docker‑lint’s `.docker-lint.yaml` structure (ignored, failure‑threshold, severities) with Markdown‑specific sections.
* **Build/Release**: Make + Goreleaser; tagged releases; reproducible builds; `--version` variants supported.
* **Docs layout**: top‑level README, SECURITY.md, rule docs mirror (DL→MD) structure.

---

## 19. Glossary

* **Finding**: a single rule violation with position and severity.
* **Failure threshold**: minimum severity that triggers non‑zero exit code.
* **Scope**: AST region(s) a rule operates on (e.g., headings, paragraphs, inline text).
* **Dictionary**: wordlist or Hunspell resources used by MD1000.

