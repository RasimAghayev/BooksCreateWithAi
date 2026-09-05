# 13. Randomness

**Səhifələr:** 213-225

## Bu fəsil nədən bəhs edir?

Təsadüfi rəqəm generatorları (PRNG), kriptoqrafik təsadüfiyyət və Go-da `crypto/rand` paketinin istifadəsi.

## Əsas fikirlər

### 1. Pseudorandom Number Generator (PRNG)
Kompyuterlər həqiqi təsadüfiyyət yarada bilmirlər, yalnız təxmin edilən (pseudorandom) dəyərlər. Bu, kriptoqrafiya üçün kifayət deyil.

### 2. crypto/rand
Go-nun kriptoqrafik təsadüfiyyət üçün paketi. OS-dən gələn ehtiyatı istifadə edir.

### 3. Randomness Tests
Təsadüfi dəyərlərin keyfiyyətini yoxlamaq üçün testlər (məsələn, chi-squared test).

## Əsas terminlər
- PRNG (Pseudorandom Number Generator)
- crypto/rand
- Randomness (Təsadüfiyyət)
- Seed ( toxum)

## Praktik nəticə

`crypto/rand` istifadə edərək təhlükəsiz təsadüfi açar və nonce-lər yaratmaq.

## Mənbə

Pages: 213-225
