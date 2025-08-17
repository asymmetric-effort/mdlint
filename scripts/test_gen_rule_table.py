"""Tests for rule table generation script."""

# Copyright (c) 2024 MdLint contributors.
# SPDX-License-Identifier: MIT

from pathlib import Path

from gen_rule_table import load_rules, render_table


def test_load_rules(tmp_path: Path) -> None:
    rules = load_rules(Path(__file__).resolve().parents[1] / "internal" / "rules")
    ids = [r.rule_id for r in rules]
    assert "MD1000" in ids and "MD1500" in ids


def test_render_table() -> None:
    rules = load_rules(Path(__file__).resolve().parents[1] / "internal" / "rules")
    table = render_table(rules)
    assert "MD1100" in table
    assert table.count("|") > 10
