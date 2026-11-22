# KOSMOMAT_PLAN.md

> Handoff pro Opus (Claude Code). Pracovní název **KosmoMat** — vědecký protějšek horoskopové appky AstroMat/EzoMat.
> Verze plánu: v1.1, 2026-11-22. Zadání: Ol1n. Jazyk obsahu při startu: cs + en.
> Změny v1.1: AstroMat/EzoMat je horoskopová appka → most (osa I) je součást prvního vydání; redakční ústava a externí kontrola odsouhlaseny (§13).

---

## 0. Jak s plánem pracovat

- Přečti celý plán před prvním commitem. Milníky M0 → M6 jdou v tomto pořadí; každý má akceptační kritéria, M2 má go/no-go bránu.
- **Nevymýšlej věcná tvrzení z hlavy.** Všechno věcné v `content/` musí mít zdroj podle §3. Roky a čísla uvedené v tomto plánu jsou vodítko pro seed, ne zdroj — při psaní karet je ověř proti primárnímu zdroji.
- Položky označené **[OVĚŘIT]** zkontroluj proti aktuální dokumentaci, než na nich postavíš kód (verze balíčků, limity API, licence).
- Stack: Go (harvester, builder, linter, API), Python jen pro LLM kroky, Flutter (appka), Svelte (redakční UI), Astro (veřejný web). Žádný bash.
- Využití Sparku (modely, priority, fronta) konzultuj s Director session.
- Otevřené otázky jsou v §13. Kde není odpověď, platí uvedený default.

---

## 1. Co stavíme

**Jedna věta:** mobilní encyklopedie vesmíru, kde každé tvrzení ukazuje, odkud se vzalo, o co se opírá, kdo s ním nesouhlasí a jak dopadly jeho předpovědi.

Není to další „vesmírná wiki“. Odlišují ji tři věci:

1. **Rodokmen myšlenky** — od každé karty se dá dojít k počátku (kdo, kdy, v jakém díle) a dopředu k tomu, co z ní vzešlo.
2. **Spory stejným metrem** — standardní model i jeho konkurenti mají stejnou šablonu: tvrzení, předpovědi, výsledky testů, otevřené problémy.
3. **Dějiny omylů a obratů** — jak lidstvo měnilo obraz vesmíru, co rozhodlo a proč to trvalo.

Čím appka **není**: horoskopy, zpravodajský agregátor, platforma jedné teorie.

Zásady produktu (Ol1nův standard): bez reklam, bez trackingu, zdarma s jasnou IAP politikou, čtení funguje offline.

---

## 2. Výchozí teze zadavatele a jak s ní appka naloží

Zadavatel vychází z osobní pozice: kosmologie velkého třesku s gravitací jako hlavní silou a s „neviditelnými“ složkami mu nedává smysl; bližší je mu Elektrický vesmír (D. E. Scott). Appka tuhle pozici **nepřebírá ani nevyvrací redakčně** — pouští ji stejným strojem jako všechno ostatní (§3). Aby seed nebyl zkreslený, drž se těchto výchozích faktů a ověř je ve zdrojích:

**2.1 Tři různé věci, které se pletou**

| Co | Kdo | Pozice v literatuře |
|---|---|---|
| Plazmová astrofyzika | Birkeland, Alfvén, Parker, dnešní heliofyzika | běžná součást oboru; magnetická pole, MHD, rekonexe, Birkelandovy proudy v polárních zářích |
| Plazmová kosmologie | Alfvén–Klein, Peratt, Lerner | publikovaná nestandardní kosmologie; většina oboru ji odmítla, protože nesedí na pozorování (CMB, zastoupení lehkých prvků) |
| Elektrický vesmír (EU) | Juergens → Thornhill, Talbott, Scott; kořeny u Velikovského | mimo recenzovanou literaturu; knihy a web |

Scottova kniha je v originále *The Electric Sky* (2006); *The Electric Universe* (2007) napsali Thornhill a Talbott. Veď obě jako `work`.

**2.2 Opravy pojmů (patří do slovníku hned v M0)**

- **Antihmota** není hypotéza — vyrábí se a měří v laboratořích (pozitrony v PET, antivodík v CERN). „Neviditelné“ složky jsou **temná hmota** a **temná energie**; ty jsou odvozené z pozorování, ne přímo detekované.
- Standardní kosmologie **není jen gravitace**. Raný vesmír je plazma; reliktní záření je světlo z okamžiku, kdy plazma zprůhlednělo; baryonové akustické oscilace jsou zvukové vlny v tom plazmatu.
- Karta „Proč gravitace vyhrává, i když je o ~36–39 řádů slabší“: hmota je ve velkém měřítku elektricky neutrální (náboje se ruší a stíní), gravitace se jen sčítá. Tohle je jádro sporu s EU — udělat poctivě, s výpočtem.

**2.3 Kde má intuice „elektřina a magnetismus“ skutečnou oporu** (samostatná osa obsahu E, §5)

Ohřev sluneční koróny, urychlování slunečního větru, magnetická rekonexe, výtrysky z okolí černých děr, magnetary a rychlé rádiové záblesky, původ galaktických a mezigalaktických magnetických polí, kosmické počasí. Tady jsou živé otevřené otázky a EM síly tu dominují.

**2.4 Skutečné trhliny standardního modelu** (osa C, §5) — jsou silnější materiál než EU, protože jsou v recenzované literatuře a mají čísla. Stav k 11/2026:

- **Hubbleovo napětí:** lokální měření ~73 km/s/Mpc (H0DN 2026: 73,50 ± 0,81) vs. hodnota z reliktního záření ~67; část týmů (CCHP, ~70,4) tvrdí, že napětí mizí. Přes 5σ, trvá.
- **Vyvíjející se temná energie:** DESI DR2 (3/2025) dává podle kombinace dat 2,8–4,2σ; práh objevu 5σ nepřekročen, výklad se zpochybňuje.
- **Hmotnost neutrin z kosmologie** tlačí na spodní mez z oscilací.
- **JWST a raný vesmír:** „malé červené tečky“ — původně „galaxie, které rozbíjejí kosmologii“, dnes převážně vykládané jako černé díry v hustém plynu; nedořešeno.
- **MOND vs. temná hmota:** MOND dobře sedí na rotační křivky galaxií, selhává u kup a CMB; testy na širokých dvojhvězdách z Gaia vycházejí v rigoróznějších analýzách pro Newtona.
- Dále **[OVĚŘIT stav]:** S8 napětí, kosmický dipól, problém lithia, nulové výsledky přímé detekce temné hmoty, timescape kosmologie, problém kosmologické konstanty.
- Čerstvá data: Euclid DR1-Foundation (12. 11. 2026; plné DR1 s kosmologií 6/2027), Rubin/LSST (ostrý průzkum od 29. 6. 2026).

**2.5 Očekávání:** při poctivém vedení bude karta EU ukazovat neúspěšné testy (hlavně sluneční neutrina — viz §5 D). To je vlastnost systému, ne chyba. Zadavatel to odsouhlasil 2026-11-22 (§13).

---

## 3. Redakční ústava (`docs/EDITORIAL.md`, vzniká v M0, je veřejná i v appce)

Tohle je jádro produktu. Linter (§7) většinu pravidel vynucuje strojově.

**P1 — Pozorování ≠ výklad.** „Rotační křivky galaxií jsou ploché“ je `observation`. „Způsobuje to temná hmota“ je `claim` navázaný na `theory`. Jedno pozorování může mít víc výkladů; karta pozorování je všechny vyjmenuje.

**P2 — Dvě nezávislé osy u každého tvrzení a teorie:**

- `status` (síla důkazů): `established` | `strong` | `contested` | `speculative` | `refuted` | `historical` | `philosophical`
- `literature` (kde se o tom vede debata): `mainstream` | `minority` | `outside`

„Menšinový“ neznamená „špatný“ a „mainstream“ neznamená „jistý“. Appka tyhle osy nikdy neslévá do jedné nálepky.

**P3 — Pravidla pro `status`** (mechanická, ne dojmová):

- `established`: ≥ 2 nezávislé metody nebo týmy, zdroje úrovně T1, žádný živý spor v T1 za posledních 10 let.
- `contested`: existují T1 zdroje na obou stranách.
- `speculative`: zatím neexistuje rozhodující test.
- `refuted`: selhala jasně formulovaná předpověď (záznam v knize předpovědí, P5) nebo existuje přímo odporující pozorování s T1 zdrojem.
- `philosophical`: empiricky nerozhodnutelné; karta přiřazuje postoje konkrétním myslitelům.

**P4 — Úrovně zdrojů** (`tier`):

| Tier | Co | Poznámka |
|---|---|---|
| T1 | recenzovaný článek, data-release paper kolaborace | |
| T2 | preprint | v appce štítek „nerecenzováno“ |
| T3 | dizertace, habilitace | |
| T4 | tisková zpráva / stránka agentury či observatoře | nikdy jediný zdroj pro `established` |
| T5 | kniha, monografie | i Scott, Thornhill/Talbott, Kuhn |
| T6 | popularizace, blog, žurnalistika | jen jako „ohlas“, nikdy jako opora tvrzení |

Příznak `retracted` z Crossref **[OVĚŘIT dostupnost dat Retraction Watch v Crossref API]**.

**P5 — Kniha předpovědí.** Každá teorie má seznam předpovědí: co tvrdila, kdy, **zda před měřením** (`made_before_observation`), kdo testoval, jak to dopadlo (`confirmed` | `refuted` | `partial` | `open` | `untestable`). Dodatečné vysvětlení už známého faktu se počítá jako retrodikce a je tak označeno. Tohle je hlavní nástroj férového srovnání.

**P6 — Symetrie.** Každá `theory` má stejné povinné sekce: hlavní tvrzení, ≥ 3 předpovědi, ≥ 2 otevřené problémy, ≥ 1 konkurent, „co by ji vyvrátilo“ (pokud možno slovy zastánců). Platí pro ΛCDM stejně jako pro MOND, plazmovou kosmologii i EU.

**P7 — Žádné tvrzení bez zdroje.** Každý odstavec mimo `philosophical` má ≥ 1 citaci s lokátorem (strana, sekce, rovnice, obrázek).

**P8 — Steelman.** Pozici popisujeme tak, aby ji její zastánce uznal za věrnou; kritiku připisujeme konkrétním autorům, ne „vědcům“.

**P9 — Jazyk.** Zakázané obraty v těle karet: „vědci dokázali“, „je jisté“, „mainstream tají“, „pavěda“ bez atribuce. Místo „vědci se mýlili“ vždy kdo, kdy, v čem a co ho opravilo.

**P10 — Datum ověření.** Každá karta má `last_verified`. U `contested` a `speculative` starší než 12 měsíců appka ukáže „čeká na revizi“.

**P11 — Autorská práva.** Citát ≤ 20 slov s lokátorem, nejvýš jeden na zdroj a kartu; jinak vlastní formulace. Abstrakty se nepřebírají. Obrázky jen s licencí z allowlistu (§8) a s kreditem.

**P12 — Historie karty je veřejná.** Zdroj pravdy je git; každá změna stavu má záznam v `changelog` s důvodem a zdrojem.

---

## 4. Obsahový model

### 4.1 Typy uzlů

`question`, `observation`, `claim`, `theory`, `prediction`, `case` (dějinná případová studie), `person`, `work` (bibliografický záznam), `mission` (mise, přístroj, observatoř), `term` (slovník), `experiment` (ověř si sám), `era`, `worldview` (mýtus / historický model světa), `body` (těleso sluneční soustavy), `constellation`.

ID: `<typ>.<slug>`, např. `obs.flat-rotation-curves`, `theory.lcdm`, `pred.cmb-relic-radiation`.

### 4.2 Hrany

`first_proposed_in` (→ work), `proposed_by` (→ person), `builds_on`, `inspired_by`, `supported_by` (→ observation), `contradicted_by` (→ observation), `interprets` (claim → observation), `competes_with`, `superseded_by`, `predicts` (theory → prediction), `tested_by` (prediction → observation/mission), `answers` (theory/worldview → question), `opposed_by` (→ person, s citací), `measured_by` (observation → mission).

Každá hrana `supported_by` / `contradicted_by` / `opposed_by` nese vlastní `source` (work + lokátor).

### 4.3 Soubory (zdroj pravdy = git, složka `content/`)

```
content/
  nodes/<typ>/<slug>/
    node.yaml
    cs.lay.md   en.lay.md        # povinné
    cs.kid.md   en.kid.md        # volitelné
    cs.student.md en.student.md  # volitelné
  works/<rok>/<klíč>.yaml
  media/<sha256>.<ext>  +  media.yaml
  paths/<slug>.yaml               # čtenářské trasy
```

`node.yaml` (příklad, zkráceno):

```yaml
id: theory.electric-universe
type: theory
status: refuted            # vyplyne z knihy předpovědí, viz L3
literature: outside
title: { cs: "Elektrický vesmír", en: "Electric Universe" }
summary:
  cs: "Soubor myšlenek, podle nichž vesmírné děje řídí hlavně elektrické proudy v plazmatu, ne gravitace."
  en: "..."
first_proposed: { year: 1972, by: person.ralph-juergens, in: work.juergens-1972 }
edges:
  - { rel: inspired_by,  to: person.immanuel-velikovsky, source: { work: ..., loc: "..." } }
  - { rel: inspired_by,  to: theory.plasma-cosmology }
  - { rel: competes_with, to: theory.lcdm }
  - { rel: predicts,     to: pred.eu-sun-no-fusion-neutrinos }
  - { rel: contradicted_by, to: obs.solar-neutrino-spectrum, source: { work: work.borexino-2020, loc: "abstract" } }
required_sections: [claims, predictions, open_problems, competitors, would_refute]
last_verified: 2026-11-22
verified_by: human
```

`prediction` (příklad):

```yaml
id: pred.cmb-relic-radiation
type: prediction
theory: theory.hot-big-bang
statement: { cs: "Vesmír vyplňuje tepelné záření o teplotě několika kelvinů.", en: "..." }
made: { year: 1948, by: [person.ralph-alpher, person.robert-herman], in: work.alpher-herman-1948 }
made_before_observation: true
test: { year: 1965, obs: obs.cmb, by: [person.arno-penzias, person.robert-wilson] }
outcome: confirmed
```

Konvence v Markdownu: citace `[@work-klic s. 12]`, odkaz na uzel `[[obs.flat-rotation-curves]]`, odkaz s textem `[[term.sigma|5σ]]`. Builder převádí Markdown na **bloky JSON** (odstavec = pole spanů: text / cite / link / term); appka nevykresluje Markdown, takže nezávisí na markdown balíčku a citace jsou klikatelné.

### 4.4 Pack (SQLite, generuje builder)

`nodes`, `node_i18n(node_id, lang, level, title, summary, blocks_json)`, `edges`, `works`, `node_sources(node_id, work_id, locator)`, `predictions`, `media(id, node_id, sha256, credit, license, source_url, w, h)`, `changelog(node_id, date, kind, note, work_id)`, `paths`, FTS5 tabulka nad title + summary + plain text per jazyk.

---

## 5. Obsahové osy a seed

Cíl pro první release: **~170 uzlů** (bez `work` a `term`), z toho 23 vstupních karet osy I. Kvalita před objemem.

### A. Velké otázky (`question`) — vstupní brána, 10 ks

Má smysl je zahrnout: lidé nepřicházejí kvůli článkům, ale kvůli otázkám. Podmínka je **oddělit vrstvy** a každou otázku označit `testability: empirical | partial | philosophical`.

Odkud se vzal vesmír? · Měl začátek a co bylo „před“? · Z čeho je? · Jak je velký a kde jsme? · Jak skončí? · Jsme sami? · Co je čas? · Proč jsou zákony takové, jaké jsou (jemné vyladění)? · Proč je něco, a ne nic? · Jak poznám, komu věřit, když nejsem odborník?

Každá otázka má tři vrstvy odpovědí na jedné časové ose: **mýtus a filozofie** (`worldview`) → **historické modely** → **současná věda a její spory**. Vrstvy se nemíchají; mýtus se nehodnotí jako věda a věda se nepodává jako mýtus.

Poslední otázka je praktická výbava čtenáře: kontrolní seznam (předpověď před měřením? čísla? nezávislé potvrzení? kde je to publikováno? co by autora přesvědčilo o opaku?).

### B. Dějiny omylů a obratů (`case`) — 20 ks

Šablona každé studie:

1. Co se tehdy vědělo a proč to dávalo smysl
2. Kdo přišel s čím a odkud to vzal
3. Námitky — které byly rozumné a které ne
4. Co rozhodlo (měření, přístroj, předpověď)
5. Jak dlouho to trvalo a kdo se přesvědčit nedal
6. Poučení — a kde už neplatí

**Důležité vyvážení:** vedle „kacířů, kteří měli pravdu“ musí být stejnou šablonou „kacíři, kteří ji neměli“ a „omyly mainstreamu opravené během měsíců“. Jinak appka učí chybu přeživších („smáli se Galileovi, tedy mám pravdu“).

Potvrzení odmítaných:
- Heliocentrismus: Aristarchos → Koperník 1543 → Galileo 1610 → Kepler 1609/1619 → Newton 1687; aberace (Bradley 1729), paralaxa (Bessel 1838). Pozor na férovost: Tychonovy námitky (chybějící paralaxa, velikosti hvězd) byly ve své době rozumné a Koperníkovy kružnice nebyly přesnější než Ptolemaios.
- Kameny z nebe: Chladni 1794, pád u L'Aigle 1803.
- Stáří Země: Kelvin (desítky milionů let) vs. geologové; rozhodlo radiometrické datování.
- Kontinentální drift: Wegener 1912 → přijato až v 60. letech (rozpínání mořského dna, paleomagnetismus).
- Složení Slunce: Cecilia Payne 1925, odrazena Russellem, který totéž potvrdil o pár let později.
- Velká debata 1920 (Shapley–Curtis) → Hubbleovy cefeidy v M31.
- Rozpínání vesmíru: Lemaître 1927, Hubble 1929. Hubbleova první hodnota byla ~7× mimo — správná myšlenka, špatné číslo.
- Velký třesk vs. stacionární vesmír: název „Big Bang“ razil odpůrce Hoyle; rozhodlo reliktní záření 1965 předpovězené 1948.
- Chandrasekhar vs. Eddington 1935.
- Černé díry: od Einsteinových pochyb po LIGO 2015 a EHT 2019.
- Temná hmota: Zwicky 1933, ~40 let bez odezvy, Rubin a Ford v 70. letech.
- **Birkelandovy proudy:** Birkeland 1908, odpor (Chapman), potvrzení družicemi 1966–1974.
- **Alfvén:** MHD vlny 1942 přijaty se skepsí, Nobelova cena 1970 — a tentýž člověk později s kosmologií neuspěl. Jeden z nejcennějších případů: stejný autor na obou stranách.
- Sluneční neutrina: 30 let „chybějících“ neutrin; model Slunce byl správně, neúplná byla částicová fyzika (oscilace, SNO 2001–2002).

Neptun vs. Vulkán (klíčová dvojice pro spor o temnou hmotu): Le Verrier z odchylek Uranu předpověděl Neptun (nalezen 1846) — „neviditelná hmota“ vyšla. Tentýž Le Verrier z odchylek Merkuru předpověděl planetu Vulkán — neexistuje, řešením byla nová teorie gravitace (1915). Obě cesty v dějinách jednou vyhrály.

Odmítnutí, která se nepotvrdila: éter (Michelson–Morley 1887), kanály na Marsu, stacionární vesmír, unavené světlo, Velikovskij 1950 (včetně nefér reakce vědců — nátlak na nakladatele; selhání na obou stranách).

Omyly mainstreamu opravené rychle: BICEP2 2014 (prach, ne gravitační vlny z inflace), OPERA 2011 (uvolněný kabel, ne nadsvětelná neutrina).

Mýtus o mýtu: **„Středověk věřil v placatou Zemi“ je sám omyl** — kulatost znali Řekové od 5. stol. př. n. l., Eratosthenés ~240 př. n. l. změřil obvod; pověst o placatém středověku rozšířilo 19. století. Placatá Země proto není dobrý příklad „těžkého přesvědčování vědců“; lepší jsou případy výše. Samostatná karta + karta o moderním hnutí placaté Země.

Neevropské linie (aby dějiny nebyly jen Evropa): babylonská astronomie; Árjabhata (rotace Země, 499); čínské školy kaj-tchien / chun-tchien / süan-jie (nekonečný prázdný prostor); al-Bírúní; Túsího pár a jeho stopa u Koperníka; Ulugbek; mayský kalendář; polynéská navigace.

### C. Živé spory (`claim` + `observation`) — 12 ks

Viz §2.4. Každý spor má obrazovku „Spor“ (§6): společná pozorování uprostřed, výklady po stranách, kniha předpovědí pod tím.

### D. Teorie vedle sebe (`theory`) — 10 ks

ΛCDM · inflace · MOND a relativistická rozšíření · plazmová kosmologie (Alfvén–Klein, Peratt, Lerner) · Elektrický vesmír · stacionární vesmír · unavené světlo · timescape **[OVĚŘIT]** · cyklické modely · geocentrismus (jako historická teorie se stejnou šablonou).

Povinné záznamy v knize předpovědí pro EU (ověřit ve zdrojích zastánců i kritiků):

- Slunce bez jaderné fúze v jádru → žádná nebo jiná neutrina. Měření: spektrum slunečních neutrin odpovídá pp řetězci i CNO cyklu (Borexino 2014–2020). `refuted`.
- Neexistence neutronových hvězd a černých děr → časování pulsarů, gravitační vlny, snímky EHT.
- Komety jako elektricky nabitá skalnatá tělesa (Deep Impact 2005) — zastánci hlásí úspěch; uvést jejich tvrzení i protiargumenty, `outcome` podle zdrojů.
- Rudý posuv není mírou vzdálenosti (Arp) → dilatace času u supernov, BAO.

Povinné otevřené problémy pro ΛCDM: Hubbleovo napětí, žádná detekovaná částice temné hmoty, problém kosmologické konstanty, lithium, původ inflace.

### E. Elektřina a magnetismus ve vesmíru — 12 ks

Polární záře a Birkelandovy proudy · magnetosféry planet · sluneční cyklus a dynamo · ohřev koróny · sluneční vítr a Parker Solar Probe · magnetická rekonexe · kosmické počasí · magnetary · rychlé rádiové záblesky · výtrysky aktivních galaxií · galaktická magnetická pole a jejich původ · plazma v laboratoři vs. ve vesmíru (škálování, Debyeovo stínění).

### F. Mise a přístroje (`mission`) — 15 ks, **[OVĚŘIT aktuální stav každé]**

NASA: JWST, Hubble, Parker Solar Probe, Chandra · ESA: Planck, Gaia, Euclid, Solar Orbiter · Chile: ALMA, VLT, Rubin/LSST · Čína: FAST, Chang'e, Tianwen, Einstein Probe · Indie: Aditya-L1, Chandrayaan · Japonsko: XRISM, Hayabusa2 · mezinárodní: LIGO–Virgo–KAGRA, DESI, EHT, IceCube, SKA.

SpaceX veď jako dopravce a inženýrskou firmu, ne jako zdroj vědeckých výsledků; samostatná karta „Starlink a astronomie“ (stopy družic ve snímcích) je legitimní živý spor.

### G. Ověř si sám (`experiment`) — 10 ks

Eratosthenův stín (dvě města, jeden den) · fáze Venuše dalekohledem · Jupiterovy měsíce triedrem · paralaxa palcem a pak hvězdná · Foucaultovo kyvadlo · spektrum Slunce z CD · proč loď mizí za obzorem od spodu · zatmění Měsíce a kulatý stín · magnet a kompas vs. polární záře · odhad rychlosti světla z Jupiterových měsíců (Rømer).

### H. Slovník (`term`) — ~25 ks

Sigma a „5σ“ · preprint vs. recenze · rudý posuv · standardní svíčka · standardní pravítko · temná hmota · temná energie · antihmota · plazma · MHD · paradigma (Kuhn) · falzifikace (Popper) · výzkumný program (Lakatos) · retrodikce · systematická chyba · replikace · Occamova břitva…

### I. Most k AstroMat/EzoMat — součást prvního vydání, 23 uzlů

AstroMat/EzoMat je horoskopová appka; KosmoMat je její vědecký protějšek. AstroMat je proto hlavní akviziční kanál a most musí fungovat od vydání (M5).

**Vstupní karty (cílové stránky odkazů z AstroMatu)**

- 10 těles (`body.*`): Slunce, Měsíc, Merkur, Venuše, Mars, Jupiter, Saturn, Uran, Neptun, Pluto.
- 13 souhvězdí na ekliptice (`constellation.*`): 12 zvířetníkových + Hadonoš.
- Obsah každé karty: co to fyzicky je · jak je to daleko a jak dlouho odtud letí světlo · kde je to dnes na obloze · co zjistily mise · jedna otevřená otázka. U souhvězdí navíc sekce „Znamení a souhvězdí“.

**Pravidla tónu (doplněk redakční ústavy)**

- Vstupní karta nezačíná hodnocením astrologie. Čtenář přišel ze zvědavosti; dostane oblohu.
- Astrologie jako téma má vlastní karty `worldview.astrology` a `case.astrology-tests` se stejnou šablonou jako cokoli jiného: rodokmen (Babylon → helénismus → arabský svět → renesance → dnešek), předpovědi, testy (Carlson 1985, Forerův efekt 1949, spor o „Mars efekt“). Bez posměchu, se zdroji.
- Na vstupních kartách je na ně jeden odkaz dole („Co říkají testy astrologie“). Nevnucovat, neschovávat.
- Povinná férová nuance: západní (tropický) zvěrokruh je definovaný jarní rovnodenností, ne souhvězdími — „posunutá znamení“ nejsou omyl, o kterém by astrologové nevěděli; védská astrologie používá zvěrokruh siderický. Precese posouvá obojí vůči sobě zhruba o jedno znamení za ~2000 let.
- Do rodokmenu patří, že astronomie a astrologie byly do 17. století jedno řemeslo (Kepler sestavoval horoskopy).

**Technický kontrakt**

- Univerzální odkaz `https://<web>/n/<id>` (appka → web → store) a schéma `kosmomat://node/<id>`. iOS Universal Links a Android App Links přes soubory na webu.
- ID uzlů osy I jsou stabilní (`body.mars`, `constellation.aries`); AstroMat je smí mít natvrdo. Přejmenování = přesměrování v manifestu packu.
- Zpětný odkaz KosmoMat → AstroMat jen jako jeden řádek na kartách osy I, nikde jinde.
- Polohy těles: sdílený Dart balíček pro obě appky (pracovně `sky_core`) nad **Astronomy Engine** — MIT, verze v C jsou dva soubory (`astronomy.h`, `astronomy.c`), volat přes `dart:ffi`; přesnost ±1′ stačí na znamení i na „kde je to dnes“.
- **Swiss Ephemeris** má duální licenci: AGPL (celý projekt pak musí být pod AGPL) nebo placená Professional License (CHF 700, platnost 6 let, neomezený počet projektů). První krok v M0: zjistit, čím polohy počítá AstroMat. Pokud Swiss Ephemeris, buď koupit licenci, nebo AstroMat převést na `sky_core`.

---

## 6. Appka (Flutter)

Spodní navigace: **Dnes · Otázky · Spory · Dějiny · Hledat**; vpravo nahoře **Moje**.

**Dnes** — 3–5 karet: co se změnilo (z `changelog`), případ dne, otázka dne. Žádný nekonečný feed.

**Karta (detail uzlu)**
- Titulek, shrnutí, čipy `status` + `literature` (klepnutí vysvětlí, co nálepka znamená a podle jakého pravidla ji karta má).
- Přepínač úrovně: Dítě / Laik / Student (jen pokud varianta existuje).
- Tělo z bloků; citace jsou horní indexy → **list zdroje**: název, autoři, rok, tier, země pracoviště, DOI/arXiv odkaz, lokátor, případně krátký citát, příznak „nerecenzováno“ / „staženo“.
- **Rodokmen** (níže), „Kdo oponuje“, „Co z toho vzešlo“.
- Média s kreditem a licencí.
- „Historie karty“ (changelog) a `last_verified`.

**Rodokmen** — svislá časová osa: předci nahoře až k počátku (`first_proposed`), potomci dole; boční větve = konkurenti. Vlastní `CustomPaint`, žádná force-graph knihovna. Klepnutí na uzel = přechod na jeho kartu. Akceptační požadavek: z každého `claim`/`theory` se dá proklikat k uzlu s `first_proposed`.

**Spor** — vodorovně přejížděné sloupce (2–4 výklady), uprostřed připíchnutá společná pozorování; pod tím tabulka knihy předpovědí (řádek = předpověď, sloupec = teorie, buňka = výsledek + rok + „před měřením ano/ne“). Na konci „Co by koho přesvědčilo o opaku“.

**Dějiny** — časová osa epoch; případové studie podle šablony §5 B; filtr „potvrzení / nepotvrzení / rychlé opravy“.

**Otázky** — seznam; detail = tři vrstvy odpovědí na ose (mýtus → historické modely → dnešek).

**Hledat** — offline FTS, filtry typ / status / éra.

**Ověř si sám** — dostupné z karet a z Hledat; kroky, co budeš potřebovat, co uvidíš, co to dokazuje a co ne.

**Moje** — záložky, poznámky, rozečtené trasy, výsledky kvízů. Jen lokálně, bez účtu.

**Nastavení** — jazyk, výchozí úroveň, aktualizace packu, Redakční ústava, Kredity a licence (povinné kvůli CC BY atribucím).

Technicky:
- Riverpod, go_router, SQLite pack přes `drift` nebo `sqlite3` s FTS5 **[OVĚŘIT verze]**.
- Pack v1 přibalený v binárce; při startu dotaz na manifest a tiché stažení novější verze; přepnutí při příštím spuštění. Média líně s diskovou cache.
- Žádná analytika, žádné SDK třetích stran mimo nákupy.
- Univerzální odkazy a schéma `kosmomat://` (§5 I); neznámé ID otevře Hledat s předvyplněným dotazem.
- Sdílení: karta jako obrázek (titulek, shrnutí, status, 1 zdroj, QR/odkaz na web). Použitelné i pro PromoHub.
- Přístupnost: dynamická velikost písma, kontrast čipů ne jen barvou.
- Typografie: `cute_kid_fonts` jen pro úroveň Dítě, jinak čitelný serif pro tělo **[rozhodne Ol1n při kuraci]**.

---

## 7. Server a nástroje (Go)

```
kosmomat/
  app/        Flutter
  content/    zdroj pravdy (git)
  cmd/
    kosmolint/     linter obsahu
    kosmobuild/    content → pack.sqlite(.zst) + manifest + web export
    harvester/     stahování metadat
    reviewd/       API pro redakční UI
  internal/...
  pipeline/   Python (LLM kroky)
  review/     Svelte
  web/        Astro
  docs/       EDITORIAL.md, SCHEMA.md, SOURCES.md
```

**Linter (`kosmolint`) — pravidla**

- L1 každý odstavec mimo `philosophical` má ≥ 1 citaci
- L2 každý citovaný `work` existuje; T1/T2 má DOI nebo arXiv ID
- L3 `status` odpovídá pravidlům P3 (počty a tiery zdrojů, existence hran, záznam `refuted` v knize předpovědí)
- L4 `theory` má všechny povinné sekce podle P6
- L5 každé médium má kredit, licenci z allowlistu a `source_url`
- L6 existuje `cs.lay` i `en.lay`; shrnutí ≤ 280 znaků
- L7 `last_verified` u `contested`/`speculative` ≤ 12 měsíců (jinak varování)
- L8 žádné sirotky; z každého `claim`/`theory` vede cesta k `first_proposed`
- L9 citát ≤ 20 slov, s lokátorem, max 1 na zdroj a kartu
- L10 zakázané obraty (P9)
- L11 každá hrana `supported_by`/`contradicted_by`/`opposed_by` má `source`

**Builder (`kosmobuild`)** — validace → bloky JSON → SQLite + FTS → komprese → `manifest.json {version, sha256, size, min_app, url}` → export pro Astro.

**Distribuce** — pack, média i web jako **statické soubory na Cloudflare (Pages/R2)**. Čtení v appce nesmí záviset na dostupnosti JODA ani Sparku; domácí infrastruktura obsah jen vyrábí.

**Postgres na JODA** (sqlc, golang-migrate) jen pro výrobní linku: `works_raw`, `harvest_runs`, `candidates`, `extractions`, `drafts`, `reviews`.

**Redakční UI (Svelte, za Cloudflare Access)** — fronta kandidátů, vedle sebe návrh karty / zdrojové pasáže / výsledek verifikace; schválit → commit do `content/`; vrátit s poznámkou. Musí být pohodlné na telefonu.

---

## 8. Zdroje dat a licence (`docs/SOURCES.md`)

| Zdroj | Na co | Podmínky |
|---|---|---|
| arXiv (OAI-PMH, API) | preprinty; kategorie astro-ph.*, physics.plasm-ph, physics.space-ph, gr-qc, physics.hist-ph | metadata CC0; max 1 požadavek / 3 s, jedno spojení; plné texty mají vlastní licence — nepřebírat |
| SciX / NASA ADS API | recenzovaná literatura, citace, bibkódy | token; denní limity **[OVĚŘIT]**; ADS přechází na SciX, odkazy ADS dál fungují |
| OpenAlex | afiliace a země, citační vazby, OA stav | od 13. 2. 2026 povinný API klíč; zdarma ~1 USD využití denně; data CC0; pro objem použít snapshot |
| Crossref REST | DOI metadata, retrakce | „polite pool“ s mailto **[OVĚŘIT]** |
| Semantic Scholar | — | licence API omezuje komerční užití bez dohody → **nepoužívat** |
| NASA Image and Video Library (`images-api.nasa.gov`) | obrázky | bez klíče; obsah NASA zpravidla volný, ale pozor na díla třetích stran; žádná loga, žádný dojem podpory |
| NASA Hubble / Webb (STScI) | obrázky | volné s kreditem NASA, ESA, CSA, STScI |
| ESA | obrázky | jen položky výslovně pod CC BY-SA 3.0 IGO (kredit + odkaz na licenci; share-alike); ostatní nepoužívat |
| ESA/Hubble, ESA/Webb | obrázky | CC BY 4.0 **[OVĚŘIT]** |
| ESO (VLT, ALMA) | obrázky | CC BY 4.0, plný kredit; úpravy (i AI) výslovně označit → agenturní snímky AI neupravovat |
| Rubin Observatory | obrázky | **[OVĚŘIT licenci]**; data sama jsou 2 roky neveřejná |
| Wikidata | data o osobách, roky | CC0 |
| Wikipedie | jen k dohledání zdrojů | text nepřebírat |
| Knihy (Scott, Thornhill/Talbott, Kuhn…) | T5 | parafráze + citát podle P11 |
| Historická díla (Koperník, Galileo, Kepler, Newton) | originály volné | překlady a skeny mívají vlastní práva **[OVĚŘIT kus od kusu]** |

**Dizertace a regionální zdroje [OVĚŘIT dostupnost a API každého]:** OATD, NDLTD, DART-Europe, HAL (Francie), Shodhganga (Indie), BDTD a SciELO (Brazílie, Latinská Amerika), Repositorio UChile, ChinaXiv a časopis RAA (Čína). ProQuest a CNKI jsou uzavřené — nepočítat s nimi.

Praktická poznámka k zadání „americké/čínské/indické/brazilské/chilské/evropské univerzity“: naprostá většina astro výzkumu ze všech těchto zemí končí na arXivu a v SciX. Zeměpisný původ proto řeš **metadatem** (země pracoviště z OpenAlex, zobrazená v listu zdroje), ne šesti oddělenými sběrači. Regionální repozitáře přidej až pro dizertace. Sleduj metriku podílu zdrojů podle regionu; je to ukazatel, ne kvóta.

---

## 9. Výrobní linka obsahu

```
harvester (Go) → triage (LLM) → extract (LLM) → draft (LLM) → verify (LLM, jiný prompt/model) → review (člověk) → build (Go)
```

1. **harvester** — denně nová metadata (arXiv, SciX, RSS agentur) do `works_raw`; deduplikace přes DOI / arXiv ID / bibkód.
2. **triage** — pro každý záznam: týká se některého existujícího uzlu? mění jeho stav? je to nový uzel? Výstup `candidate {work, node_ids[], kind: update|new|ignore, reason}`.
3. **extract** — z abstraktu (a z plného textu jen u licencí CC BY) strukturovaná tvrzení: `{text, kind: observation|interpretation, quantities[], method, dataset, evidence_span}`. `evidence_span` je doslovná pasáž uložená jen interně pro verifikaci.
4. **draft** — návrh `node.yaml` + `cs.lay.md` + `en.lay.md` nebo diff existující karty; každá věta má citaci.
5. **verify** — pro každou větu: vyplývá z citované pasáže? (ano / částečně / ne). Čísla v textu se musí rovnat číslům v pasáži. Kontrola P9 a L-pravidel. Výstup jde do redakčního UI.
6. **review** — Ol1n schvaluje. Nic se nepublikuje bez člověka.
7. **build** — lint → pack → publikace.

Modely: lokální přes LiteLLM/vLLM na Sparku podle doporučení Director session; `verify` jiným modelem nebo aspoň jiným promptem než `draft`.

**Evaluační sada (vzniká v M2):** 50 ručně ověřených tvrzení ze seed uzlů + 20 záměrně chybných (špatné číslo, obrácený závěr, citace mimo). Metriky: přesnost citací (podíl vět skutečně podložených), číselná shoda, záchyt nastražených chyb.

**Externí kontrola (odsouhlaseno):** před vydáním projde placený člověk z oboru (doktorand MFF UK / AsÚ AV ČR) 30 náhodných karet + všechny karty `theory` + `case.astrology-tests`. Oslovit ho už během M2, ať je v M3 k dispozici. Zadání pro něj: věcné chyby (špatně / nepřesně / v pořádku), nefér podání některé strany, chybějící klíčový zdroj. Bez toho nevydávat.

---

## 10. Milníky

**M0 — Kostra a zlatý vzorek** (bez LLM)
- Repo, `EDITORIAL.md`, `SCHEMA.md`, `kosmolint`, `kosmobuild`.
- 12 ručně psaných uzlů pokrývajících všechny typy, mezi nimi: spor rotační křivky → temná hmota vs. MOND; případ Neptun/Vulkán; teorie ΛCDM a EU se stejnou šablonou a knihou předpovědí; slovník antihmota / temná hmota / temná energie; 1 experiment; 1 otázka se třemi vrstvami.
- Zjistit, čím AstroMat počítá polohy těles, a dohodnout stabilní ID uzlů osy I (§5 I). Do zlatého vzorku přidat `body.mars` a `constellation.aries` (celkem 14 uzlů).
- Akceptace: lint bez chyb; pack se sestaví; Ol1n zlatý vzorek přečetl a odsouhlasil tón i šablony.

**M1 — Čtečka**
- Dnes, Karta, list zdroje, Rodokmen, Hledat, přibalený pack, Kredity.
- Akceptace: všech 14 uzlů čitelných v letadlovém režimu; každá citace otevře list zdroje; z každého claim/theory dojdu k počátku; studený start < 2 s na středním Androidu.

**M2 — Spike výrobní linky + brána**
- harvester (arXiv + SciX + Crossref), triage → verify, evaluační sada.
- Oslovit externího recenzenta (§9).
- **Go/no-go:** přesnost citací ≥ 95 %, číselná shoda 100 %, záchyt nastražených chyb ≥ 90 %.
- Při nesplnění ani po druhé iteraci: linka jen navrhuje zdroje a pasáže, text píše člověk s asistencí. Rozsah prvního vydání se sníží na osy B + C + D + I (~100 uzlů).

**M3 — Obsahový sprint na ~170 uzlů**
- Redakční UI; média s manifestem licencí; cs + en; 23 vstupních karet osy I a balíček `sky_core`.
- Akceptace: lint bez chyb; externí kontrola 30 karet bez věcné chyby stupně „špatně“; `last_verified` u všech.
- Stop kritérium: propustnost review < 20 uzlů týdně po dvou týdnech → škrtat rozsah, ne kvalitu.

**M4 — Spor, kniha předpovědí, Dějiny, Otázky, Ověř si sám, úrovně**
- Akceptace: 12 sporů s vyplněnou tabulkou předpovědí; 20 případových studií, z toho ≥ 6 „nepotvrzených“ a ≥ 2 „rychlé opravy“.

**M5 — Web a vydání**
- Astro web z téhož exportu (SEO, univerzální odkazy), sdílecí karty, store release cs/en přes existující pipeline (screenshoty, lokalizace).
- Most: univerzální odkazy ověřené z AstroMatu na obou platformách (appka nainstalovaná / nenainstalovaná).
- Kategorie Vzdělávání / Reference. V metadatech „astronomie, kosmologie, vesmír“ a pro publikum z AstroMatu „znamení a souhvězdí, planety“. Žádný příslib horoskopů.

**M6 — Čerstvost a rozšíření**
- Týdenní běh linky → „Co se změnilo“; upozornění „karta čeká na revizi“.
- IAP (§11), volitelně „Zeptej se“ (RAG jen nad packem, odpověď výhradně s citacemi uzlů).

Dogfooding po celou dobu: Ol1n čte denně 5 karet v appce; chyby hlásí přes redakční UI.

---

## 11. Monetizace

- Zdarma: celá encyklopedie úrovně Laik, spory, dějiny, hledání, offline.
- Jednorázové odemknutí „Hloubka“: úroveň Student, čtenářské trasy, kvízy, HD média offline, „Zeptej se“ s denním limitem.
- Zahrnout do Ol1n Pass (entitlement `all_apps` podle `PASS_PLAN.md`).
- Žádné předplatné v1 — „Zeptej se“ běží na Sparku, takže limit a férové upozornění na dostupnost.
- Křížová propagace: AstroMat → KosmoMat přes karty těles a souhvězdí; obě appky pod Ol1n Pass.

Upřímně: encyklopedie se monetizuje špatně (Wikipedie je zdarma). Hodnota projektu je spíš v obsahu pro krátká videa („Kdo se mýlil“, „Spor týdne“ přes PromoHub), v důvěryhodnosti značky a jako vstup do Passu.

---

## 12. Rizika

| Riziko | Dopad | Opatření |
|---|---|---|
| Zaujatost redakce (kterýmkoli směrem) | ztráta důvěry obou publik | ústava §3, symetrické šablony, linter, externí kontrola, veřejná historie karet |
| Halucinace LLM | věcné chyby se zdrojem „na oko“ | verify, brána M2, člověk schvaluje vše |
| Zastarávání | „nejnovější poznatky“ rok staré | `last_verified`, týdenní linka, stav „čeká na revizi“ |
| Licence obrázků a textů | stažení ze storu | allowlist, manifest, P11 |
| Rozsah | nikdy nevyjde | 170 uzlů, stop kritéria |
| Publikum z horoskopové appky čeká potvrzení a dostane fyziku | odinstalace, špatné recenze | vstupní karty jsou o obloze, ne o hodnocení; pravidla tónu §5 I; store popis neslibuje horoskopy |
| Dvě appky jedné značky říkají o astrologii různé věci | ztráta důvěry u obou publik | KosmoMat neútočí a nic neschovává; AstroMat rámovat jako zábavu a sebereflexi, ne jako předpověď (rozhodnutí §13) |
| Licence efemerid | stažení AstroMatu i KosmoMatu ze storu | `sky_core` nad Astronomy Engine (MIT), nebo koupená licence Swiss Ephemeris |
| Falešná vyváženost | okrajová teorie vypadá jako rovnocenná | dvě osy P2 + kniha předpovědí; rovný je metr, ne výsledek |
| Závislost na domácí infrastruktuře | appka nejde číst | statická distribuce přes Cloudflare |

---

## 13. Rozhodnutí a otevřené otázky

**Rozhodnuto 2026-11-22**

- AstroMat/EzoMat je horoskopová appka; KosmoMat je její vědecký protějšek → osa I je v prvním vydání.
- Redakční ústava platí bez výjimek, včetně toho, že karta Elektrického vesmíru ponese výsledky testů, jaké jsou. Totéž platí pro kartu o testech astrologie.
- Externí odborná kontrola před vydáním: ano, placená.

**Otevřené (platí default, dokud Ol1n neřekne jinak)**

1. **Čím AstroMat počítá polohy těles?** Rozhoduje o licenci a o balíčku `sky_core`. Default: Astronomy Engine pro obě appky.
2. **Rámování AstroMatu** vůči KosmoMatu. Default: AstroMat jako zábava a sebereflexe; KosmoMat na vstupních kartách astrologii nehodnotí, karta testů je dohledatelná jedním odkazem.
3. **Název.** Default: KosmoMat (řada AstroMat / EzoMat / KosmoMat).
4. **Publikum.** Default: zvídavý dospělý laik; úroveň Dítě jen u vybraných karet.
5. **Jazyky při vydání.** Default: cs + en, další přes lokalizační pipeline po M5.
6. **„Zeptej se“ v první verzi?** Default: ne, až M6.
7. **Vizuální styl** — fotografie agentur + vlastní diagramy; generované ilustrace jen u mýtů a dějin a vždy označené. Default: ano.

---

## 14. Seed: prvních 14 uzlů pro M0

`term.antimatter` · `term.dark-matter` · `term.dark-energy` · `obs.flat-rotation-curves` · `claim.dark-matter-halos` · `theory.lcdm` · `theory.mond` · `theory.electric-universe` · `pred.eu-sun-no-fusion-neutrinos` (+ `obs.solar-neutrino-spectrum`) · `case.neptune-and-vulcan` · `q.what-is-the-universe-made-of` · `exp.eratosthenes-shadow` · `body.mars` · `constellation.aries`

---

## 15. Výchozí odkazy z rešerše (ověřit a doplnit primárními zdroji)

- Test EU na neutrinech a spektru: https://briankoberlein.com/blog/testing-electric-universe/
- EU a recenzní řízení (reportáž z konference EU2015): https://www.scientificamerican.com/article/the-difference-between-science-and-pseudoscience/
- Borexino, CNO neutrina: https://arxiv.org/abs/2006.15115
- DESI DR2, temná energie: https://newscenter.lbl.gov/2025/03/19/new-desi-results-strengthen-hints-that-dark-energy-may-evolve/
- Hubbleovo napětí, H0DN 2026: https://earthsky.org/space/local-universes-expansion-rate-clear-but-questions/ ; přehled: https://arxiv.org/abs/2606.20434
- Široké dvojhvězdy vs. MOND: https://arxiv.org/abs/2602.24035
- Malé červené tečky: https://cerncourier.com/a/the-mystery-of-the-little-red-dots/
- Euclid DR1: https://arxiv.org/abs/2611.16510
- Rubin, raná data: https://rubinobservatory.org/for-scientists/resources/early-science
- OpenAlex, klíče a ceny: https://developers.openalex.org/api-reference/authentication
- arXiv API, podmínky: https://info.arxiv.org/help/api/tou.html
- NASA Image API: https://images.nasa.gov/docs/images.nasa.gov_api_docs.pdf
- ESO licence: https://www.eso.org/public/copyright/ ; ESA: https://www.esa.int/ESA_Multimedia/Copyright_Notice_Images
- Mýtus o placaté Zemi: https://en.wikipedia.org/wiki/Myth_of_the_flat_Earth
- Swiss Ephemeris, licence: https://www.astro.com/swisseph/swephinfo_e.htm ; smlouva Professional License: https://www.astro.com/swisseph/secont_e.pdf
- Astronomy Engine: https://github.com/cosinekitty/astronomy
