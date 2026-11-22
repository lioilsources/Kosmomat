# KosmoMat

Mobilní encyklopedie vesmíru, kde každé tvrzení ukazuje, odkud se vzalo, o co se opírá, kdo s ním
nesouhlasí a jak dopadly jeho předpovědi.

- Redakční ústava: [`docs/EDITORIAL.md`](docs/EDITORIAL.md)
- Formát obsahu: [`docs/SCHEMA.md`](docs/SCHEMA.md)
- Zdroje a licence: [`docs/SOURCES.md`](docs/SOURCES.md)

## Stav

- **M0 — kostra a zlatý vzorek:** linter, builder packu, 15 karet (cs + en), všechny `verified_by: draft`.
- **M1 — čtečka:** Flutter appka (`app/`) — Dnes, Karta, list zdroje, Rodokmen, kniha předpovědí, Hledat
  (offline FTS), Moje, Nastavení, Kredity, Redakční ústava; přibalený pack.

## Použití

```
go run ./cmd/kosmolint -warnings
go run ./cmd/kosmobuild -version 0.1.0
```

Výstup v `out/`: `pack.sqlite`, `pack.sqlite.zst`, `manifest.json`, `web/` (JSON pro Astro web).
