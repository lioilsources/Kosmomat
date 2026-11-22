# Redakční ústava KosmoMatu

Tahle pravidla platí pro každou kartu bez výjimky — pro standardní kosmologii stejně jako pro její
konkurenty, pro vědu stejně jako pro astrologii. Jsou veřejná a v appce dostupná z Nastavení.
Většinu z nich hlídá strojově linter `kosmolint`; číslo pravidla linteru je v závorce.

## P1 — Pozorování ≠ výklad

„Rotační křivky galaxií jsou ploché“ je **pozorování** (`observation`). „Způsobuje to temná hmota“
je **tvrzení** (`claim`) navázané na teorii. Jedno pozorování může mít víc výkladů a karta
pozorování je všechny vyjmenuje.

## P2 — Dvě nezávislé osy

Každé tvrzení a teorie má dvě nálepky, které se nikdy neslévají do jedné:

- **síla důkazů** (`status`): ověřené · silné · sporné · spekulativní · vyvrácené · historické · filozofické
- **kde se o tom vede debata** (`literature`): hlavní proud · menšinová literatura · mimo recenzovanou literaturu

„Menšinový“ neznamená „špatný“ a „hlavní proud“ neznamená „jistý“.

## P3 — Pravidla pro sílu důkazů (L3)

Nálepka se neuděluje dojmem, ale podle počítatelných pravidel:

- **ověřené**: aspoň dvě nezávislé metody nebo týmy, recenzované zdroje (T1), žádný živý spor
  v recenzované literatuře za posledních 10 let;
- **sporné**: recenzované zdroje existují na obou stranách;
- **spekulativní**: zatím neexistuje rozhodující test;
- **vyvrácené**: selhala jasně formulovaná předpověď (zapsaná v knize předpovědí, P5),
  nebo existuje přímo odporující pozorování s recenzovaným zdrojem;
- **filozofické**: empiricky nerozhodnutelné; karta přiřazuje postoje konkrétním myslitelům.

## P4 — Úrovně zdrojů

| Úroveň | Co | Jak s ní zacházíme |
|---|---|---|
| T1 | recenzovaný článek, datová publikace kolaborace | |
| T2 | preprint | v appce štítek „nerecenzováno“ |
| T3 | dizertace, habilitace | |
| T4 | tisková zpráva, stránka agentury nebo observatoře | nikdy jediný zdroj pro „ověřené“ |
| T5 | kniha, monografie | včetně knih zastánců menšinových teorií i filozofů vědy |
| T6 | popularizace, blog, žurnalistika | jen jako ohlas, nikdy jako opora tvrzení |

Stažené (retracted) práce se necitují (L2).

## P5 — Kniha předpovědí

Každá teorie má seznam předpovědí: co tvrdila, kdy, **zda před měřením**, kdo ji testoval a jak to
dopadlo (potvrzeno · vyvráceno · částečně · otevřené · netestovatelné). Dodatečné vysvětlení už
známého faktu je **retrodikce** a je tak označeno. Kniha předpovědí je hlavní nástroj férového
srovnání: rovný je metr, ne výsledek.

## P6 — Symetrie (L4)

Každá teorie má stejné povinné sekce: hlavní tvrzení, aspoň 3 předpovědi, aspoň 2 otevřené
problémy, aspoň jednoho konkurenta a „co by ji vyvrátilo“ — pokud možno slovy jejích zastánců.

## P7 — Žádné tvrzení bez zdroje (L1, L2)

Každý odstavec mimo filozofické karty má aspoň jednu citaci s lokátorem (strana, sekce, rovnice,
obrázek). Výjimkou jsou jen návody (kroky pokusu) a údaje, které appka sama počítá (poloha na obloze).

## P8 — Steelman

Pozici popisujeme tak, aby ji její zastánce uznal za věrnou. Kritiku připisujeme konkrétním
autorům, ne „vědcům“.

## P9 — Jazyk (L10)

V textu karet nepoužíváme: „vědci dokázali“, „je jisté“, „mainstream tají“, „pavěda“ bez
atribuce. Místo „vědci se mýlili“ píšeme vždy kdo, kdy, v čem a co ho opravilo.

## P10 — Datum ověření (L7)

Každá karta má datum posledního ověření. U sporných a spekulativních karet starších než
12 měsíců appka ukáže „čeká na revizi“.

## P11 — Autorská práva (L5, L9)

Citát má nejvýš 20 slov a lokátor, nejvýš jeden na zdroj a kartu; jinak vlastní formulace.
Abstrakty nepřebíráme. Obrázky jen s licencí ze seznamu v [`SOURCES.md`](SOURCES.md) a s kreditem;
agenturní snímky neupravujeme umělou inteligencí.

## P12 — Historie karty je veřejná

Zdrojem pravdy je git. Každá změna stavu karty má záznam v historii karty s důvodem a zdrojem.

## Doplněk pro vstupní karty těles a souhvězdí (osa I)

- Vstupní karta nezačíná hodnocením astrologie. Čtenář přišel ze zvědavosti; dostane oblohu.
- Astrologie má vlastní karty se stejnou šablonou jako cokoli jiného: rodokmen, předpovědi, testy.
  Bez posměchu, se zdroji. Na vstupní kartě je na ně jeden odkaz dole — nevnucovat, neschovávat.
- Tropický zvěrokruh je definovaný jarní rovnodenností, ne souhvězdími; „posunutá znamení“ nejsou
  omyl, o kterém by astrologové nevěděli.
- Astronomie a astrologie byly do 17. století jedno řemeslo.

## Jak se karta dostane do appky

Návrh (člověk nebo výrobní linka) → linter → kontrola člověkem → sestavení packu. Nic se
nepublikuje bez lidského schválení. Před prvním vydáním projde náhodný vzorek karet a všechny
teorie placený recenzent z oboru.
