# rozhovory-scraper

**Program stáhne data o pořadu [Hlavní zprávy - rozhovory a komentáře](https://radiozurnal.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846) z webu Českého rozhlasu**.

[![main](https://github.com/czech-radio/rozhovory-scraper/actions/workflows/main.yml/badge.svg)](https://github.com/czech-radio/rozhovory-scraper/actions/workflows/main.yml) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/518b8ee5b79240e78d3b955beb19d393)](https://app.codacy.com/gh/czech-radio/rozhovory-scraper/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

- [x] Název pořadu
- [x] Název epizody
- [x] Popis epizody
- [ ] Počet témat epizody (z názvu).
- [x] Čas vysílání epizody (polední = 12:10, odpolední = 18:10).

- [ ] Získej moderátora (částěčně splněno).
- [ ] Získej hosty.

## Build

```bash
git clone git@github.com:czech-radio/rozhovory-scraper.git
cd rozhovory-scraper
```

### Unix

```bash
...
```

### Windows

```powershell
.\build.ps1
```

## Usage

```bash
./rozhovory-scraper -p [počet_stran]
```

- `p` Počet stran ke stažení, výchozí 1.
