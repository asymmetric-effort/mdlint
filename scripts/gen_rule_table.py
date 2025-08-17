"""Generate Markdown rule reference tables from annotated Go rule files."""

from __future__ import annotations

import pathlib
import re
from dataclasses import dataclass
from typing import Iterable, List


# Copyright (c) 2024 MdLint contributors.
# SPDX-License-Identifier: MIT

ROOT = pathlib.Path(__file__).resolve().parents[1]
RULES_DIR = ROOT / "internal" / "rules"
DOC_PATH = ROOT / "docs" / "rules.md"


@dataclass
class Rule:
    """Metadata about a single lint rule."""

    rule_id: str
    name: str
    summary: str
    severity: str
    options: str


def parse_rule_file(path: pathlib.Path) -> Rule:
    """Parse a Go rule file and extract rule metadata."""
    text = path.read_text(encoding="utf-8")
    data = {}
    for key in ("RuleID", "Name", "Summary", "Severity", "Options"):
        match = re.search(rf"^// {key}: (.+)$", text, re.MULTILINE)
        if not match:
            raise ValueError(f"missing {key} in {path}")
        data[key.lower()] = match.group(1).strip()
    return Rule(
        rule_id=data["ruleid"],
        name=data["name"],
        summary=data["summary"],
        severity=data["severity"],
        options=data["options"],
    )


def load_rules(directory: pathlib.Path) -> List[Rule]:
    """Load all rule files from the provided directory."""
    rules: List[Rule] = []
    for path in sorted(directory.glob("md*.go")):
        rules.append(parse_rule_file(path))
    return rules


def render_table(rules: Iterable[Rule]) -> str:
    """Render a Markdown table of rules."""
    header = "| ID | Name | Summary | Default Severity | Options |\n| --- | --- | --- | --- | --- |"
    rows = [
        f"| {r.rule_id} | {r.name} | {r.summary} | {r.severity} | {r.options} |"
        for r in rules
    ]
    return "\n".join([header, *rows])


def main() -> None:
    """Entry point for the script."""
    rules = load_rules(RULES_DIR)
    table = render_table(rules)
    content = (
        "<!-- markdownlint-disable MD013 -->\n"
        "# Built-in Rules\n\n"
        + table
        + "\n<!-- markdownlint-enable MD013 -->\n"
    )
    DOC_PATH.write_text(content, encoding="utf-8")


if __name__ == "__main__":  # pragma: no cover
    main()
