# T120B165-Web-Application-Design

## Nekilnomojo turto valdymo sistema

### Uždavinio aprašymas

### Funkciniai reikalavimai

1. **Svečio funkcijos:**
    - Peržiūrėti temas
    - Peržiūrėti edukacijų skelbimus
    - Aktyvuoti naujienlaiškio prenumeratą
    - Registruotis kaip organizatoriui
2. **Organizatoriaus funkcijos:**
    - Sukurti temą
    - Paslėpti/Rodyti temą
    - Pašalinti temą
    - Sukurti skelbimą
    - Paslėpti/Rodyti skelbimą
    - Pašalinti skelbimą
3. **Administratoriaus funkcijos:**
    - Patvirtinti temą/skelbimą
    - Pašalinti paskyrą
    - Patvirtinti naują paskyrą

### Pasirinktų technologijų aprašymas:
1. Backend technologijos:
   - Programavimo kalba: Golang
   - Duomenų bazė: MariaDB
   - Autentifikacija: JWT

2. Frontend technologijos:
   - React.js


### Objektai:

1. **Tema**
    - Pavadinimas
    - Aprašymas
2. **Veikla**
    - Pavadinimas
    - Aprašymas
    - Kaina
    - Vieta (gatvė, miestas, šalis arba koordinatės).
    - Bazinė kaina
    - Kontaktai
3. **Paketas**
    - Pavadinimas
    - Aprašymas
    - Kaina
    - Kontaktai


### Rolės:

- **Svečias**: gali ieškoti skelbimų, aktyvuoti naujienlaiškio prenumeratą.
- **Organizatorius**: gali skelbti skelbimus ir temas.
- **Administratorius**: turi visas valdymo funkcijas, įskaitant vartotojų administravimą, skelbimų ir temų priežiūrą.

