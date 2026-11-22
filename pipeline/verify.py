"""verify — does each sentence follow from the passages it cites? (plan §9, step 5)

Input:  JSONL from `kosmoeval export` — {"id", "sentence", "passages": [{"work", "year", "text"}]}
Output: JSONL for `kosmoeval score` — {"id", "verdict": yes|partial|no, "numbers_ok", "reason", "model"}

The model judges support only; numbers are checked in code (numcheck.py), and a
"yes" with a failed number check is downgraded to "partial".

Model endpoint: any OpenAI-compatible chat API (LiteLLM on SPARK via the broker):
  LLM_BASE_URL  e.g. http://<gateway>/v1   (no default: hosts live in config, not code)
  LLM_MODEL     model name as the gateway lists it
  LLM_API_KEY   optional bearer token

  python3 pipeline/verify.py out/eval/items.jsonl out/eval/verdicts.jsonl
"""

import argparse
import json
import os
import re
import sys
import time
import urllib.error
import urllib.request
from typing import Callable

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from numcheck import numbers_agree  # noqa: E402

PROMPT_VERSION = "verify-v1"

SYSTEM = """You check citations in a popular-science encyclopedia.
You get one SENTENCE and the PASSAGES it cites. Judge only from the passages; do not use your own knowledge,
even if you know the sentence is true.

verdict:
- "yes": everything the sentence asserts is stated in, or directly follows from, the passages.
- "partial": the core claim is supported, but some detail (a number, qualifier, name, date, scope) is not.
- "no": the passages contradict the sentence, or do not support its core claim.

A publication year shown in a passage header counts as support for when that work appeared.
Reply with JSON only: {"verdict": "yes|partial|no", "unsupported": ["..."], "reason": "one sentence"}"""

VERDICTS = ("yes", "partial", "no")


def user_message(item: dict) -> str:
    parts = [f"SENTENCE:\n{item['sentence']}\n", "PASSAGES:"]
    for i, p in enumerate(item["passages"], 1):
        parts.append(f"[{i}] {p['work']} ({p.get('year', '?')}):\n{p['text']}\n")
    return "\n".join(parts)


def parse_reply(text: str) -> tuple[str, str]:
    """Extract verdict and reason from a model reply; anything unparsable is "no"."""
    m = re.search(r"\{.*\}", text, re.S)
    if m:
        try:
            obj = json.loads(m.group(0))
            v = str(obj.get("verdict", "")).strip().lower()
            if v in VERDICTS:
                return v, str(obj.get("reason", ""))
        except json.JSONDecodeError:
            pass
    return "no", f"unparsable reply: {text[:200]!r}"


Chat = Callable[[list[dict]], str]


def http_chat(base_url: str, model: str, api_key: str = "", timeout: float = 120, retries: int = 3) -> Chat:
    url = base_url.rstrip("/") + "/chat/completions"

    def chat(messages: list[dict]) -> str:
        body = json.dumps({"model": model, "messages": messages, "temperature": 0}).encode()
        headers = {"Content-Type": "application/json"}
        if api_key:
            headers["Authorization"] = f"Bearer {api_key}"
        last = None
        for attempt in range(retries + 1):
            req = urllib.request.Request(url, data=body, headers=headers, method="POST")
            try:
                with urllib.request.urlopen(req, timeout=timeout) as resp:
                    return json.load(resp)["choices"][0]["message"]["content"]
            except (urllib.error.URLError, TimeoutError, KeyError) as e:
                last = e
                time.sleep(min(5 * 2**attempt, 60))
        raise RuntimeError(f"chat failed after {retries + 1} attempts: {last}")

    return chat


def verify_item(item: dict, chat: Chat, model: str) -> dict:
    passages = [p["text"] for p in item["passages"]]
    years = [str(p["year"]) for p in item["passages"] if p.get("year")]
    numbers_ok, missing = numbers_agree(item["sentence"], passages, years)
    verdict, reason = parse_reply(chat([
        {"role": "system", "content": SYSTEM},
        {"role": "user", "content": user_message(item)},
    ]))
    if verdict == "yes" and not numbers_ok:
        verdict, reason = "partial", f"numbers not in passages: {missing}; model: {reason}"
    return {"id": item["id"], "verdict": verdict, "numbers_ok": numbers_ok, "reason": reason,
            "model": f"{model}/{PROMPT_VERSION}"}


def run(items_path: str, out_path: str, chat: Chat, model: str, log=print) -> int:
    done = set()
    if os.path.exists(out_path):  # resume: skip ids already judged
        with open(out_path) as f:
            done = {json.loads(line)["id"] for line in f if line.strip()}
    n = 0
    with open(items_path) as f, open(out_path, "a") as out:
        for line in f:
            if not line.strip():
                continue
            item = json.loads(line)
            if item["id"] in done:
                continue
            v = verify_item(item, chat, model)
            out.write(json.dumps(v, ensure_ascii=False) + "\n")
            out.flush()
            n += 1
            log(f"{v['verdict']:8} {'' if v['numbers_ok'] else '#'} {item['id']}")
    return n


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("items", nargs="?", default="out/eval/items.jsonl")
    ap.add_argument("out", nargs="?", default="out/eval/verdicts.jsonl")
    a = ap.parse_args()
    base, model = os.environ.get("LLM_BASE_URL"), os.environ.get("LLM_MODEL")
    if not base or not model:
        sys.exit("set LLM_BASE_URL and LLM_MODEL (LiteLLM on SPARK, leased via sparkctl)")
    n = run(a.items, a.out, http_chat(base, model, os.environ.get("LLM_API_KEY", "")), model)
    print(f"{n} items judged → {a.out}", file=sys.stderr)


if __name__ == "__main__":
    main()
