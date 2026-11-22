# Osa I — most z AstroMatu

## Stabilní ID vstupních karet (kontrakt s AstroMatem)

AstroMat smí tato ID mít natvrdo. Nepřejmenovávat; když už musí, záznam v `content/redirects.yaml`
(pack i web export ho nesou dál).

| Tělesa | Souhvězdí na ekliptice |
|---|---|
| `body.sun` | `constellation.aries` — Beran |
| `body.moon` | `constellation.taurus` — Býk |
| `body.mercury` | `constellation.gemini` — Blíženci |
| `body.venus` | `constellation.cancer` — Rak |
| `body.mars` | `constellation.leo` — Lev |
| `body.jupiter` | `constellation.virgo` — Panna |
| `body.saturn` | `constellation.libra` — Váhy |
| `body.uranus` | `constellation.scorpius` — Štír |
| `body.neptune` | `constellation.ophiuchus` — Hadonoš |
| `body.pluto` | `constellation.sagittarius` — Střelec |
| | `constellation.capricornus` — Kozoroh |
| | `constellation.aquarius` — Vodnář |
| | `constellation.pisces` — Ryby |

Slugy souhvězdí jsou latinské názvy IAU malými písmeny (`scorpius`, ne `scorpio`; `capricornus`,
ne `capricorn`) — znamení a souhvězdí jsou různé věci a ID je o souhvězdí.

Hotové v M0: `body.mars`, `constellation.aries`. Zbytek v M3.

**Mapování znamení → karta** (pro AstroMat): znamení Štír odkazuje na `constellation.scorpius`;
Hadonoš znamení nemá, ale karta existuje a odkazuje se na ni z karty Štíra a Střelce.

## Odkazy

- Univerzální odkaz: `https://<web>/n/<id>` (appka → web → store), schéma `kosmomat://node/<id>`.
- Neznámé ID otevře Hledat s předvyplněným dotazem.
- Zpětný odkaz KosmoMat → AstroMat je jen jeden řádek na kartách osy I.

## Efemeridy: zjištění z M0

AstroMat (`../Astro/packages/astro_core`) počítá polohy balíčkem `astronomia` 1.1.1 (MIT, čistý Dart:
VSOP87B, Meeus; Pluto z tabulky DE440), s golden testy proti skyfield/DE440.
Swiss Ephemeris je jen ve vývojovém nástroji `tools/golden/` a do appky se nelinkuje.
→ **Licenci Swiss Ephemeris kupovat nemusíme.**

**Rozhodnuto 2026-11-22 (Ol1n souhlasí, odchylka od defaultu plánu):** `sky_core` nestavět nad Astronomy Engine
přes `dart:ffi`, ale vyčlenit ho z `astro_core` (tropické polohy, východ/západ, fáze), který už má
ověřený řetězec (světelný čas, precese, aberace, nutace) a golden testy. Ušetří to nativní build
C knihovny pro iOS i Android a obě appky budou počítat stejným kódem. Astronomy Engine zůstává
záložní variantou, kdyby `astro_core` chyběla funkce (např. viditelnost nad obzorem pro konkrétní místo).

Postup v M3: v repu Astro vyčlenit `packages/sky_core` (polohy, východ/západ, souhvězdí podle hranic IAU)
z `astro_core`, `astro_core` na něm postavit a KosmoMat ho připojí jako git/path závislost.
