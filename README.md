# MdLint

![MdLint icon](docs/img/logo.png)

MdLint is a fast, extensible Markdown linter written in Go.

## Installation

```bash
go install github.com/asymmetric-effort/mdlint@latest
```

Binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/asymmetric-effort/mdlint/releases).

## Usage

```bash
mdlint [flags] <path|glob>...
```

Common flags:

| Flag | Description |
| --- | --- |
| `-c, --config <file>` | Use a specific config file |
| `-o, --output <format>` | Output format: `json` or `text` |
| `--fail-level <sev>` | Minimum severity that causes a non-zero exit |

## Configuration

MdLint reads options from `.mdlintrc.yaml` in your project root. Example:

```yaml
version: 1
ignored:
  - MD1500
severity:
  MD1100: error
paths:
  "docs/**":
    ignored: [MD1400]
spell:
  lang: en_US
  add_words: [GoLand]
heading:
  style: atx
  allow_mixed: false
failure_threshold: warning
```

## Pre-commit

Use MdLint as a [pre-commit](https://pre-commit.com/) hook:

```yaml
repos:
  - repo: https://github.com/asymmetric-effort/mdlint
    rev: v1.0.0
    hooks:
      - id: mdlint
        args: ["--fail-level", "warning"]
```

## Editor Integration

* **VS Code:** add a task running `mdlint` and enable the problem matcher.
* **JetBrains IDEs:** configure `mdlint` as an External Tool and enable
  "Show in output".

## Built-in Rules

Rule documentation is generated from source annotations. See the [rule reference](docs/rules.md).
