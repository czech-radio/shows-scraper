---
title:    Selected shows collector  
author:   Czech Radio  
created:  2023-04-15  
---

# Selected shows collector

**This application is used to crowle/scrape data for internal usage.**

## Motivace a zadání

Pro potřeby vnitřních analýz se pravidelně reportuje o výskytech respondentů ve vybraných pořadech na stanici Plus a Radiožurnál. Protože manuální získávání a ověřování takových informací je časově náročné, potřebujeme alespoň část tohoto procesu automatizovat.

Výsledkem práce programu  by měla být tabulka obsahující informacemi o epizodách vybraných pořadů vysílaných na stanicích Plus a Radiožurnál. Tabulka musí být snadno importovatelná do programu Microsoft Excel. Pro každou vysílanou epizodu, tzn. řádek tabylky, očekáváme alespoň následující informace:

### Přehled sloupců tabulky

|Název|Datový typ|Komentář|
|-----|----------|--------|
|název stanice | string | Lze uvažovat o číselném kódu (11 \| 13)
|datum vysílání | string (DD:MM:YYYY)|
|čas začátku vysílání | string (HH:MM)|
|čas konce vysílání | string (HH:MM) |
|název pořadu | string|
|název epizody| string|
|typ epizody| string (premiéra \| repríza) |
|jméno a příjmení moderátora| string |
|jméno a příjmení respondentů| string | oddělit středníkem (`;`) např. `Jiří Sova;Jan Sokol`
|upoutávka či popis (teaser)| string |
|webový odkaz na použitý zdroj | string (URL) | (pokud byl použit/existuje)

### Přehled vysílaných pořadů

|Pořad|Stanice|Respondenti|Moderátor|Vysíláno| Délka|Poznámka
|-----|-------|-----------|---------|--------|------|--------
|[Dvacet minut Radiožurnálu](https://radiozurnal.rozhlas.cz/dvacet-minut-radiozurnalu-5997743)| Plus, Radiožurnal | 1 | 1 | po-pá 17:06 (Plus) / po-pá 17:06; repríza 0:10 (Radiožurnál) | 20 minut | Moderuje *Vladimír Kroc* nebo *Tomáš Pancíř*.
|[Interview Plus](https://plus.rozhlas.cz/interview-plus-6504167)| Plus | 1 | 1 | po-pá 11:34; repríza po-pá 16:34, út-pá 05:05, so 06:05 | 25 minut | Moderuje *Veronika Sedláčková* nebo *Jan Bumba*.
|[Pro a proti](https://plus.rozhlas.cz/pro-a-proti-6482952) | Plus | 2 |1 | po-pá 09:33; repríza po-pá 14:33, út-so 22:35, 03:33 | 24 minut |  Moderuje *Karolína Koubová*.
|[Hlavní zprávy - rozhovory a komentáře](https://radiozurnal.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846) | Plus, Radiožurnál| 3 | 1 | po-pá 12:10 a 18:10 | 20 minut | Moderuje *Věra Štechrová*, *Vladimír Kroc*, &hellip;
| [Speciál](https://radiozurnal.rozhlas.cz/special-radiozurnalu-7770703) | Plus?, Radiožurnál | různé | 1 | občasně | 30 a více minut | Moderuje *Jan Pokorný*, &hellip;

## Zpracování

Část informací je možné získat z interních informačních systémů, kde však nejsou systematicky zaznamenáni respondenti v příslušných pořadech. Jedním z dobrých zdrojů  informací se však ukazují být webové stránky příslušných pořadů. Pro potřeby oddělení tedy byl vytvořen jednodychý program, který na požádání stáhne a připraví data pro další zpracování.

- [x] Protoyp v Jupyter notebooku.
- [ ] Vytvoř rozhraní pro příkazovou řádku.
- [ ] Ukládej (cashuj) již jednou stažené/zpracované pořady.
- [ ] Rozhodni jak a jestli lze automatizovaně zpracovat speciály.

Výstupní tabulka Excel by měl a obsahovat sloupce 

- datum
- čas
- pořad
- moderátor
- příjmení
- jméno
- strana
- popis_funkce
- funkce
- popis_téma např.

```csv
datum	čas	pořad	moderátor	příjmení	jméno	strana	popis_funkce	funkce	popis_téma
2/1/2023		20 m RŽ	Kroc, Vladimír	Berg	Michal	Ostatní	spolupředseda Strany zelených	senátor	kauza Bečva
```

## Použití

Program pracuje následujícím způsobem:

1. Stáhni progam pro dan stanice a období.
   Data se ukládájí do souboru `./schedule.csv`
2. Získej pouze epizody/premiéry daných pořadů.
   Data se zpracovávají pomocí přiloženého [notebooku](./process.ipynb).
3. Pro každou epizodu se snažíme získat:
   a. moderátora (vždy jeden)
   b. očekávaný počet respondentů (1 až 3)  
4. Stáhni z webu pro Hlavní zprávy - rozhovory a komentáře pomocí <https://github.com/czech-radio/rozhovory.scraper> a získaná data přidej do výsledné tabulky.


## Instalace

- Vytvoř virtuální prostředí.

  ```powershell
  py -m venv --upgrade-deps .venv
  ```

- Aktivuje virtuální prostředí.

  ```powershell
  .\.venv\Scritpts\activate
  ```

- Instaluj závislosti do virtuálního prostředí.

  ```powershell
  py -m pip install -r .\requirements.txt
  ```

Pro potřeby prototypování a ukázek používáme nástroj [Jupyter](https://jupyter.org/). Notebooky je možné prohlížet v prostředí GitHub a upravovat přímo v editoru Visual Studio Code.

### Odkazy

- <https://rapidoc.croapp.cz/>
- <https://help.geneea.com/api_media2>
- <https://github.com/czech-radio/schedule>
- <https://github.com/czech-radio/rozhovory.scraper>
