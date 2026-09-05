# 5. Keys

**Səhifələr:** 85-103

## Bu fəsil nədən bəhs edir?

Kriptoqrafik açarlar (keys) və keyspace (açar məkanı) anlayışı. Açar uzunluğunun təhlükəsizlik üzərində təsiri və skorinq mexanizmləri.

## Əsas fikirlər

### 1. Keyspace (Açar Məkanı)
Mümkün açar dəyərlərinin ümumi sayı. Keyspace neçə böyük olarsa, brute force hücumla kırmaq bir o qədər çətindir.

### 2. Açar uzunluğu
- 1-byte (8 bit) açar: 256 mümkün dəyər — asan kırılır
- 32-byte (256 bit) açar: 2²⁵⁶ mümkün dəyər — praktiki olaraq qırılmaz

### 3. Scoring Candidates (Namizədləri skorlama)
Düzgün açarı tapmaq üçün mümkün açarların skorunu hesablamaq (məsələn, ingilis dilinin tezlik cədvəli ilə müqayisə).

## Əsas terminlər
- Key (Açar)
- Keyspace (Açar məkanı)
- Brute Force (Zorla hücum)
- Scoring (Skorlama)
- Candidate (Namizəd)

## Praktik nəticə

Açar uzunluğu artırıldıqda crack funksiyasının performansı ölçülür. Real dünya açarları üçün 256-bit və daha yüksək tövsiyə olunur.

## Mənbə

Pages: 85-103
