# publicistika-scraper
Tvorba prehledu hostu temat a moderatoru vybranych poradu publicistiky.

[![Go](https://github.com/czech-radio/publicistika-scraper/actions/workflows/go.yml/badge.svg)](https://github.com/czech-radio/publicistika-scraper/actions/workflows/go.yml)

# Zadání

Získej tabulku pro zadané období (obvykle měsíc) s informacemi o zadaných pořadech vysílaných na stanicích Plus a Radiožurnál:

1) 20 minut Radiožurnálu, [odkaz](https://radiozurnal.rozhlas.cz/dvacet-minut-radiozurnalu-5997743)
2) Interview Plus, [odkaz](https://plus.rozhlas.cz/interview-plus-6504167)
3) Pro a Proti (Plus) [odkaz](https://plus.rozhlas.cz/pro-a-proti-6482952)
4) Hlavní zprávy: Rozhovory komentáře (Radiožurnál, Plus), [odkaz](https://radiozurnal.rozhlas.cz/hlavni-zpravy-rozhovory-a-komentare-5997846)
5) Speciál (Radiožurnál, Plus?) [odkaz](https://radiozurnal.rozhlas.cz/special-radiozurnalu-7770703)
   (prozatím není nutné implementovat)

Jako informace bereme:
1. datum vysílání
2. čas vysílání
3. název pořadu
4. název epizody (titulek)
5. moderátor
6. hosté
7. upoutávka (teaser)
8. popis (pokud existuje)
9. odkaz na zdroj, ze kterého byly inforpace čerpány

Data chceme ukládat vše do jednoho dokumentu, avšak aktualizovat průběžne po dnech.

Ukazuje se, že každý pořad bude mít vlastní zdroj, popř. více zdrojů, ze kterých se budou o nich čerpat potřebné informace např.
- Hlavní zprávy - Rozhovory a komentáře: data získáme ze stránek pořadu
- 20 minut, Interview Plus a Pro a Proti: data získáme s REST API programu
