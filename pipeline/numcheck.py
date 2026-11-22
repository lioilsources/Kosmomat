"""Deterministic number check ("číselná shoda", plan §9).

Mirror of internal/evalset/plant.go (Numbers, NumbersAgree); both are tested
against internal/evalset/testdata/numbers.json so they cannot drift apart.
A model is never asked whether numbers match.
"""

import re

_NUMBER = re.compile(r"\b\d+(?:[.,]\d+)?\b", re.ASCII)  # ASCII \b like Go: "4.2σ" → 4.2


def norm_number(m: str) -> str:
    i = m.find(",")
    if i >= 0:
        if len(m) - i - 1 == 3:  # thousands separator
            return m.replace(",", "")
        return m.replace(",", ".", 1)
    return m


def numbers(s: str) -> list[str]:
    return [norm_number(m) for m in _NUMBER.findall(s)]


def numbers_agree(sentence: str, passages: list[str], allowed: list[str] = ()) -> tuple[bool, list[str]]:
    """Every number of the sentence must occur in a passage or in `allowed`
    (publication years of the cited works, which come from citation metadata)."""
    have = set(allowed)
    for p in passages:
        have.update(numbers(p))
    missing = [n for n in numbers(sentence) if n not in have]
    return not missing, missing
