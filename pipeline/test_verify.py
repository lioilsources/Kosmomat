"""Tests for the verify step; no network, the model is a fake.

    python3 -m unittest discover -s pipeline
"""

import json
import os
import sys
import tempfile
import unittest

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)

import verify  # noqa: E402
from numcheck import numbers_agree  # noqa: E402

VECTORS = os.path.join(HERE, "..", "internal", "evalset", "testdata", "numbers.json")


class NumbersShared(unittest.TestCase):
    def test_same_vectors_as_go(self):
        with open(VECTORS) as f:
            cases = json.load(f)
        self.assertGreater(len(cases), 3)
        for c in cases:
            ok, missing = numbers_agree(c["sentence"], c["passages"], c["allowed"])
            self.assertEqual((ok, missing), (c["ok"], c["missing"]), c["sentence"])


def item(i, sentence, text, year=2016):
    return {"id": i, "sentence": sentence, "passages": [{"work": "w", "year": year, "text": text}]}


class Verify(unittest.TestCase):
    def test_parse_reply(self):
        self.assertEqual(verify.parse_reply('Sure: {"verdict": "Yes", "reason": "ok"}'), ("yes", "ok"))
        self.assertEqual(verify.parse_reply('```json\n{"verdict":"partial","reason":"r"}\n```')[0], "partial")
        self.assertEqual(verify.parse_reply("I think it is fine")[0], "no")
        self.assertEqual(verify.parse_reply('{"verdict": "maybe"}')[0], "no")

    def test_numbers_downgrade_yes(self):
        chat = lambda msgs: '{"verdict": "yes", "reason": "fine"}'
        v = verify.verify_item(item("a", "We see it in 306 galaxies.", "We see it in 153 galaxies."), chat, "fake")
        self.assertEqual((v["verdict"], v["numbers_ok"]), ("partial", False))
        v = verify.verify_item(item("b", "We see it in 153 galaxies in 2016.", "We see it in 153 galaxies."), chat, "fake")
        self.assertEqual((v["verdict"], v["numbers_ok"]), ("yes", True))
        self.assertEqual(v["model"], "fake/" + verify.PROMPT_VERSION)

    def test_prompt_carries_sentence_and_passages(self):
        seen = []

        def chat(msgs):
            seen.append(msgs)
            return '{"verdict": "no", "reason": "x"}'

        verify.verify_item(item("c", "Speeds stay flat.", "Rotation curves are flat.", 1980), chat, "fake")
        user = seen[0][1]["content"]
        self.assertIn("Speeds stay flat.", user)
        self.assertIn("w (1980)", user)
        self.assertIn("Rotation curves are flat.", user)
        self.assertEqual(seen[0][0]["role"], "system")

    def test_run_resumes(self):
        with tempfile.TemporaryDirectory() as d:
            items, out = os.path.join(d, "items.jsonl"), os.path.join(d, "verdicts.jsonl")
            with open(items, "w") as f:
                for i in ("x", "y", "z"):
                    f.write(json.dumps(item(i, "Flat.", "Flat.")) + "\n")
            with open(out, "w") as f:
                f.write(json.dumps({"id": "x", "verdict": "yes", "numbers_ok": True}) + "\n")
            calls = []
            chat = lambda msgs: calls.append(1) or '{"verdict": "yes", "reason": ""}'
            n = verify.run(items, out, chat, "fake", log=lambda *_: None)
            self.assertEqual((n, len(calls)), (2, 2))
            with open(out) as f:
                self.assertEqual([json.loads(l)["id"] for l in f], ["x", "y", "z"])


if __name__ == "__main__":
    unittest.main()
