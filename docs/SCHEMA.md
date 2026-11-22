# Schéma obsahu

Zdroj pravdy je složka `content/`. Builder (`cmd/kosmobuild`) z ní dělá pack (SQLite),
linter (`cmd/kosmolint`) hlídá pravidla z [`EDITORIAL.md`](EDITORIAL.md).

```
content/
  nodes/<typ>/<slug>/
    node.yaml
    cs.lay.md  en.lay.md          # povinné
    cs.kid.md  en.kid.md          # volitelné
    cs.student.md en.student.md   # volitelné
  works/<rok>/<klíč>.yaml
  media/<sha256>.<ext> + media.yaml
  paths/<slug>.yaml
  redirects.yaml                  # staré ID → nové ID
```

## ID a typy

ID = `<prefix>.<slug>`, slug `[a-z0-9-]+`. Složka `nodes/<typ>/<slug>` musí sedět s ID.

| typ | prefix | `status` povinný |
|---|---|---|
| question | `q` | ne (místo něj `testability`) |
| observation | `obs` | ano |
| claim | `claim` | ano |
| theory | `theory` | ano |
| prediction | `pred` | ano (u předpovědi = stav tvrzení) |
| case | `case` | ne |
| person | `person` | ne |
| mission | `mission` | ne |
| term | `term` | ne |
| experiment | `exp` | ne |
| era | `era` | ne |
| worldview | `worldview` | ne |
| body | `body` | ne |
| constellation | `constellation` | ne |

`work` není uzel, ale bibliografický záznam (`content/works/`). V YAML se na něj odkazuje
`work.<klíč>`, v Markdownu `[@<klíč> lokátor]`.

ID osy I (`body.*`, `constellation.*`) jsou **stabilní kontrakt s AstroMatem**: nepřejmenovávat;
když už musí, záznam do `redirects.yaml`.

## node.yaml

```yaml
id: theory.lcdm
type: theory
status: strong              # established|strong|contested|speculative|refuted|historical|philosophical
literature: mainstream      # mainstream|minority|outside
title:   { cs: "…", en: "…" }
summary: { cs: "…", en: "…" }   # ≤ 280 znaků
first_proposed:             # povinné u claim/theory, pokud k počátku nevede hrana (L8)
  year: 1983
  by: [person.mordehai-milgrom]   # ID osoby; uzel ještě nemusí existovat (varování)
  in: work.milgrom-1983
edges:
  - { rel: competes_with, to: theory.mond }
  - { rel: supported_by, to: obs.flat-rotation-curves, source: { work: work.rubin-1980, loc: "obr. 4" } }
predictions: []             # jen theory; inline záznamy knihy předpovědí (viz níže)
changelog:
  - { date: 2026-11-22, kind: created, note: "Zlatý vzorek M0" }
last_verified: 2026-11-22
verified_by: draft          # draft|human — draft = napsal stroj, čeká na Ol1na
todo_verify:                # co musí člověk ověřit proti primárnímu zdroji (linter hlásí jako varování)
  - "Lokátor u work.scott-2006"
```

Typově specifická pole:

| typ | pole |
|---|---|
| question | `testability: empirical\|partial\|philosophical` |
| case | `case_kind: vindicated\|not_vindicated\|quick_fix\|pair\|myth` |
| prediction | `theory`, `statement`, `made`, `made_before_observation`, `retrodiction`, `test`, `outcome`, `note` |
| body, constellation | `astromat: true` (vstupní karta osy I) |
| experiment | `difficulty: easy\|medium\|hard`, `needs: { cs: […], en: […] }` |

### Hrany

`rel` ∈ `first_proposed_in, proposed_by, builds_on, inspired_by, supported_by, contradicted_by,
interprets, competes_with, superseded_by, predicts, tested_by, answers, opposed_by, measured_by,
related`.

- `supported_by`, `contradicted_by`, `opposed_by` musí mít `source: { work, loc }` (L11).
- `to` je ID uzlu. Cíl, který ještě neexistuje, je varování (osnova pro další milníky), ne chyba.
- Zpětné hrany se nepíšou; builder je dopočítá (`interprets` z claimu se ukáže na kartě pozorování).

### Předpověď (kniha předpovědí, P5)

Buď samostatný uzel `pred.*` (pole nahoře v souboru, `theory:` = teorie), nebo inline položka
v `predictions:` teorie se stejnými poli a vlastním `id: pred.*`.

```yaml
id: pred.cmb-relic-radiation
statement: { cs: "…", en: "…" }
made: { year: 1948, by: [person.ralph-alpher], in: work.alpher-herman-1948, loc: "s. 774" }
made_before_observation: true
retrodiction: false          # true = dodatečné vysvětlení známého faktu
test: { year: 1965, obs: obs.cmb, by: [person.arno-penzias], in: work.penzias-wilson-1965, loc: "abstrakt" }
outcome: confirmed           # confirmed|refuted|partial|open|untestable
note: { cs: "…", en: "…" }   # u sporných výsledků: kdo tvrdí co
```

Kontroly: `confirmed`/`refuted`/`partial` potřebuje `test.in`; `made_before_observation: true`
vyžaduje `made.year ≤ test.year`; `retrodiction: true` vylučuje `made_before_observation: true`.

## works/<rok>/<klíč>.yaml

Klíč: `<příjmení-prvního-autora>-<rok>[-přípona]`, malými písmeny, ASCII.

```yaml
key: milgrom-1983
tier: T1                     # T1 recenze · T2 preprint · T3 dizertace · T4 agentura · T5 kniha · T6 popularizace
kind: article                # article|data-release|preprint|thesis|press|book|chapter|popular|classic
title: "A modification of the Newtonian dynamics as a possible alternative to the hidden mass hypothesis"
authors: ["Milgrom, M."]
year: 1983
venue: "The Astrophysical Journal 270, 365–370"
ids: { doi: "10.1086/161130", arxiv: "", bibcode: "1983ApJ...270..365M", isbn: "", url: "" }
country: [IL]                # země pracoviště (ISO 3166), z OpenAlex; volitelné
retracted: false
verified: false              # true = metadata ověřena proti Crossref/SciX/arXiv
```

T1/T2 musí mít DOI nebo arXiv ID (L2). **Odchylka od plánu:** u děl před érou DOI stačí `bibcode`
(SciX) — jinak by Zwicky 1933 nešel citovat jako T1. Díla `kind: classic` (např. Le Verrier 1859)
bez jakéhokoli identifikátoru linter jen varuje; patří k nim odkaz na sken (`url`).

## Markdown karet

Builder převádí Markdown na bloky JSON; appka Markdown nevykresluje.

| zápis | význam |
|---|---|
| `## Nadpis {#id}` | sekce; `id` je strojový klíč (povinné sekce, viz níže) |
| odstavec | blok `para` |
| `- položka` | blok `list`, každá položka je samostatný odstavec pro L1 |
| `> citát [@klíč lok]` | blok `quote` (P11: ≤ 20 slov, lokátor, max 1 na zdroj a kartu) |
| `[@klíč lokátor]` | citace; lokátor povinný (strana, §, rovnice, obr., tab., „abstrakt“) |
| `[@klíč lok; @klíč2 lok]` | víc citací naráz |
| `[[obs.flat-rotation-curves]]` | odkaz na uzel (text = titulek cíle) |
| `[[term.sigma\|5σ]]` | odkaz s vlastním textem; na `term.*` je to span `term` |
| `**tučně**`, `*kurzíva*` | zvýraznění |

Odstavec bez citace je chyba L1, kromě sekcí s příznakem `.nocite`, povoleným jen tady:

| typ | sekce, které smí mít `.nocite` |
|---|---|
| experiment | `needs`, `steps`, `you_will_see` |
| body, constellation | `sky_today`, `astrology_link` |
| question | `checklist` |
| všechny | sekce karty se `status: philosophical` |

Zápis: `## Kroky {#steps .nocite}`.

### Povinné sekce podle typu

| typ | sekce (id) |
|---|---|
| theory | `claims`, `predictions`, `open_problems` (≥ 2 položky), `competitors`, `would_refute` |
| case | `context`, `who`, `objections`, `decided`, `duration`, `lesson` |
| question | `myth`, `historical`, `science` |
| experiment | `needs`, `steps`, `you_will_see`, `proves`, `does_not_prove` |
| body | `what`, `distance`, `sky_today`, `missions`, `open_question`, `astrology_link` |
| constellation | jako body + `sign_vs_constellation` |
| observation | `observation`, `interpretations` |

U `theory` navíc ≥ 3 předpovědi (inline + uzly `pred.*` s `theory:` = tato teorie) a ≥ 1 hrana `competes_with`.

## Bloky JSON (v packu, `node_i18n.blocks_json`)

```json
[
  {"t":"h","id":"claims","text":"Hlavní tvrzení"},
  {"t":"p","s":[{"t":"text","v":"Rotační křivky "},{"t":"cite","w":"rubin-1980","loc":"obr. 4"}]},
  {"t":"ul","items":[[{"t":"text","v":"…"}]]},
  {"t":"q","s":[{"t":"text","v":"…"},{"t":"cite","w":"…","loc":"…"}]}
]
```

Spany: `text {v, b?, i?}`, `cite {w, loc}`, `link {to, v}`, `term {to, v}`.

## Pack

SQLite (`pack.sqlite`, komprimovaný `pack.sqlite.zst`) s tabulkami `meta`, `nodes`, `node_i18n`,
`edges`, `works`, `node_sources`, `predictions`, `media`, `changelog`, `paths`, `redirects`
a FTS5 `search(node_id, lang, title, summary, body)`. Přesné DDL: `internal/pack/schema.sql`.

`manifest.json`: `{version, sha256, size, min_app, url, nodes, built_from}`.
Web export: `out/web/n/<id>.json` + `out/web/index.json` pro Astro.
