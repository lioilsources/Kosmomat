# Zdroje dat a licence

Přehled z plánu §8. **[OVĚŘIT]** = zkontrolovat proti aktuální dokumentaci, než na tom vznikne kód
(M2 harvester). Stav k 2026-11-22: nic z tabulky zatím nebylo ověřeno.

## Metadata a literatura

| Zdroj | Na co | Podmínky |
|---|---|---|
| arXiv (OAI-PMH, API) | preprinty: astro-ph.*, physics.plasm-ph, physics.space-ph, gr-qc, physics.hist-ph | metadata CC0; max 1 požadavek / 3 s, jedno spojení; plné texty nepřebírat |
| SciX / NASA ADS API | recenzovaná literatura, citace, bibkódy | token; denní limity **[OVĚŘIT]** |
| OpenAlex | afiliace a země, citační vazby, OA stav | API klíč povinný od 13. 2. 2026; data CC0; pro objem snapshot |
| Crossref REST | DOI metadata, retrakce | „polite pool“ **[OVĚŘIT]**; data Retraction Watch **[OVĚŘIT]** |
| Semantic Scholar | — | **nepoužívat** (licence omezuje komerční užití) |
| Wikidata | osoby, roky | CC0 |
| Wikipedie | jen k dohledání zdrojů | text nepřebírat |

## Obrázky — allowlist licencí (linter L5)

Hodnota `license` v `content/media/media.yaml` musí být jedna z těchto:

| `license` | Zdroj | Podmínky |
|---|---|---|
| `PD-NASA` | NASA Image and Video Library, NASA/STScI Hubble a Webb | kredit (např. „NASA, ESA, CSA, STScI“); pozor na díla třetích stran; žádná loga, žádný dojem podpory |
| `CC-BY-SA-3.0-IGO` | ESA — jen položky výslovně pod touto licencí | kredit + odkaz na licenci; share-alike |
| `CC-BY-4.0` | ESA/Hubble, ESA/Webb **[OVĚŘIT]**, ESO (VLT, ALMA) | plný kredit; úpravy označit; agenturní snímky AI neupravovat |
| `CC-BY-SA-4.0` | jednotlivě ověřená díla | kredit + share-alike |
| `CC0-1.0` | volná díla | — |
| `PD-old` | historická díla po uplynutí práv | skeny a překlady mívají vlastní práva **[OVĚŘIT kus od kusu]** |
| `own` | vlastní diagramy KosmoMatu | — |

Rubin Observatory: licence **[OVĚŘIT]**, do té doby se nepoužívá.

## Dizertace a regionální zdroje (M2+)

OATD, NDLTD, DART-Europe, HAL, Shodhganga, BDTD a SciELO, Repositorio UChile, ChinaXiv a RAA —
**[OVĚŘIT dostupnost a API každého]**. ProQuest a CNKI jsou uzavřené.

Zeměpisný původ se řeší metadatem (`country` u díla, z OpenAlex), ne oddělenými sběrači.
Podíl zdrojů podle regionu je ukazatel, ne kvóta.

## Efemeridy (osa I)

- AstroMat (`../Astro`) počítá polohy balíčkem **`astronomia` 1.1.1** (čistý Dart, VSOP87B + Meeus);
  Pluto z tabulky JPL DE440. **Swiss Ephemeris se do appky nelinkuje**, používá se jen v nástroji
  na referenční testovací vektory (`tools/golden/`). Licenční riziko AGPL tedy neexistuje.
- Licence `astronomia` 1.1.1: **MIT** (ověřeno v `LICENSE` balíčku, 2026-11-22).
- Viz `docs/AXIS-I.md` k volbě mezi `astronomia` a Astronomy Engine pro `sky_core`.
