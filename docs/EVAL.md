# Evaluační sada a brána M2

Plán §9 a §10: krok `verify` musí na sadě 50 ověřených a 20 nastražených tvrzení dosáhnout
**přesnosti citací ≥ 95 %**, **číselné shody 100 %** a **záchytu nastražených chyb ≥ 90 %**.

## Soubory

- `eval/claims.yaml` — pravá tvrzení: úseky anglických karet, jejichž citace míří na abstrakt
  (`[@klíč abstract]`). Vygenerováno `kosmoeval candidates`, všechna `verified_by: draft`.
- `eval/planted.yaml` — 20 nastražených chyb (`kosmoeval plant`): posunuté číslo, obrácené tvrzení,
  cizí citace. Očekávaný verdikt `no`.
- `out/eval/abstracts.json` — stažené abstrakty (Crossref → arXiv → SciX). **Necommitují se** —
  abstrakty vydavatelů mají vlastní licence. Položky je pinují otiskem `passage_sha`; změní-li se
  text u zdroje, `export` i `score` selžou a položku je třeba znovu posoudit.

## Postup

```
go run ./cmd/kosmoeval candidates -out /dev/null   # jen stáhne abstrakty do out/eval/ (claims.yaml nepřepíše)
go run ./cmd/kosmoeval numbers                     # pravé věty s čísly, která v abstraktu nejsou
go run ./cmd/kosmoeval export                      # → out/eval/items.jsonl
LLM_BASE_URL=… LLM_MODEL=… python3 pipeline/verify.py   # → out/eval/verdicts.jsonl (navazuje, kde skončil)
go run ./cmd/kosmoeval score                       # metriky; exit 1 = no-go
```

`verify.py` volá jakékoli OpenAI-kompatibilní API (LiteLLM na SPARKu — pronájem přes `sparkctl`).
Model posuzuje jen oporu v pasáži; čísla kontroluje kód (`numcheck.py`) a `yes` s nesedícím
číslem sníží na `partial`.

## Co musí udělat člověk (Ol1n)

1. Projít `eval/claims.yaml`: plyne věta opravdu z abstraktu? Nastavit `expected` (`yes`/`partial`/`no`)
   a `verified_by: human`. Nevhodné položky smazat.
2. Projít `eval/planted.yaml`: je každá chyba opravdu chyba? (Obrácení „is → is not“ může náhodou
   dát pravdivou větu.)
3. Doplnit sadu na 50 pravých: chybí abstrakty 15 děl (starší články APS/Nature) — s `ADS_API_TOKEN`
   je doplní SciX.

Dokud nejsou položky `verified_by: human`, je skóre jen orientační (`score` to vypíše).

## Výchozí stav (2026-11-22)

38 pravých + 20 nastražených. Kontrolní běh s „modelem“, který vždy odpoví `yes`: přesnost
86,8 %, číselná shoda 86,5 %, záchyt 55 % (jen díky číselné kontrole) → no-go, jak má být.
Skutečný model zatím neběžel — broker SPARKu byl nedostupný.

## První běh s modelem (2026-11-22)

Qwen3.6-35B-A3B-FP8 na L40S (vLLM, `qwen36/verify-v1`), 58 položek, ~15 min. Výsledek: **no-go**.

| metrika | výsledek | práh |
|---|---|---|
| přesnost citací | 34,2 % (13/38 `yes`) | ≥ 95 % |
| číselná shoda | 86,5 % (5 vět) | 100 % |
| záchyt nastražených | 95,0 % (19/20) | ≥ 90 % |

Co z toho plyne:
- **Záchyt funguje.** Unikla jen `planted-20` („is not disputed“). Model přehlédl obrácení a
  ohodnotil jen dvojici výsledků.
- **Nízká přesnost je hlavně problém sady, ne modelu.** Všechny pravé položky mají
  `expected: yes` jen jako výchozí návrh. Řada vět ale říká víc než abstrakt:
  - závěr karty („vyvrací předpověď EU“, `theory.lcdm#18`);
  - obecné tvrzení podepřené jedním experimentem („žádný experiment“, `theory.lcdm#10`);
  - citaci na nevhodný abstrakt (`term.antimatter#7` → LZ, `theory.lcdm#1` → Planck).

  Model tu má pravdu, že to z pasáže neplyne. 18 z 25 chyb je `partial`: jádro sedí, jen detail ne.
- **Číselná shoda** je deterministická (`numcheck`) a selhává na týchž 5 větách jako dřív.
  Je potřeba opravit karty, nikoli model.

Další krok: Ol1n nastaví `expected` u pravých položek (u mnohých `partial`, karty s chybnou
citací opravit) a rozhodne, zda `partial` u pravé věty počítat do přesnosti. Pak se `score` spustí
znovu bez nového běhu modelu.
