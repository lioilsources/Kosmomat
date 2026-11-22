# KosmoMat

<!-- Max 80 řádků. Co se nevejde, jde do docs/ a sem jeden odkaz. -->

## Účel
Offline mobilní encyklopedie vesmíru: každé tvrzení se zdrojem, rodokmenem a knihou předpovědí.
Vědecký protějšek AstroMatu (`../Astro`); efemeridy přes `sky_core` vyčleněný z `astro_core` (docs/AXIS-I.md). Zadání a milníky: `Prompts/00-KOSMOMAT_PLAN.md`.

## Build a test (gate)
```
go vet ./... && go test ./...
python3 -m unittest discover -s pipeline
go run ./cmd/kosmolint                 # jen chyby; -warnings ukáže i varování
go run ./cmd/kosmobuild -version X -app app/assets/pack   # lint → pack + kopie do appky
cd app && flutter analyze && flutter test
```
Po změně `content/` vždy přestavět pack s `-app`; CI kontroluje, že přibalený pack odpovídá obsahu.

## Mapa
- `content/` — zdroj pravdy: `nodes/<typ>/<slug>/`, `works/<rok>/<klíč>.yaml`. Formát: `docs/SCHEMA.md`
- `docs/EDITORIAL.md` — redakční ústava (P1–P12), veřejná
- `docs/SOURCES.md` — zdroje dat a allowlist licencí; `docs/AXIS-I.md` — kontrakt ID s AstroMatem
- `internal/content` — loader + parser Markdownu na bloky JSON
- `internal/lint` — pravidla L1–L11 (+ S strukturální, V ověření)
- `internal/pack` — SQLite pack (`schema.sql`, FTS5), manifest, web export
- `app/` — Flutter čtečka (Riverpod, go_router, sqlite3 s FTS5); pack v `app/assets/pack/`, texty UI v ARB (cs, en)
- `cmd/harvester` + `internal/harvest` — arXiv/Crossref/SciX → `out/harvest.sqlite` (works_raw); `verify` kontroluje `content/works`
- `cmd/kosmoeval` + `internal/evalset` — evaluační sada M2 (`eval/`), nastražené chyby, skóre brány; postup `docs/EVAL.md`
- `pipeline/` — LLM kroky v Pythonu (jen stdlib); `verify.py` + `numcheck.py` (zrcadlo Go, sdílené vektory `internal/evalset/testdata/numbers.json`)
- Plánované: `cmd/reviewd`, triage/extract/draft v `pipeline/`, `review/` Svelte, `web/` Astro

## Invarianty
- Žádné věcné tvrzení bez zdroje s lokátorem. Nevymýšlet; nejistotu zapsat do `todo_verify`.
- `verified_by: human` smí nastavit jen Ol1n. Stroj píše `draft`.
- Works `verified: true` jen po kontrole proti Crossref/SciX/arXiv.
- ID osy I (`body.*`, `constellation.*`) jsou kontrakt s AstroMatem — nepřejmenovávat bez `redirects.yaml`.
- T6 nikdy jako jediná opora; agenturní snímky neupravovat AI.
- Žádný bash ve skriptech: nástroje v Go, LLM kroky v Pythonu.
- SQLite přes `modernc.org/sqlite` (CGO_ENABLED=0); testy bez sítě a bez `time.Now` (lint bere `-today`).

## Odchylky
- Odchylka od plánu §7 L2: T1 bez DOI/arXiv projde s bibkódem; `kind: classic` bez identifikátoru jen varuje.
