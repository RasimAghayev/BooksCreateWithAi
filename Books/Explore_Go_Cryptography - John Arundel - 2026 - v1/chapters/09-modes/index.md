# 9. Modes

**Səhifələr:** 146-161

## Bu fəsil nədən bəhs edir?

Blok şifrələrinin işləmə rejimləri (modes of operation). CBC, GCM (Galois/Counter Mode) və authenticated encryption (autentifikasiya olunmuş şifrələmə).

## Əsas fikirlər

### 1. CBC Mode
Cipher Block Chaining — hər blok əvvəlki blokla XOR edilir. IV ilə başlanır.

### 2. GCM Mode
Galois/Counter Mode — parallel işləyən, autentifikasiya ilə birlikdə şifrələmə.

### 3. Authenticated Encryption
Şifrələmə və mətn bütövlüyünü yoxlamaq (message integrity) eyni anda təmin edir.

## Əsas terminlər
- Mode of Operation (İşləmə rejimi)
- CBC (Cipher Block Chaining)
- GCM (Galois/Counter Mode)
- Authenticated Encryption (Autentifikasiya olunmuş şifrələmə)
- Nonce (Yanlız bir dəfə istifadə edilən rəqəm)

## Praktik nəticə

GCM rejimində şifrələmə və deşifrələmə testləri yazılır.

## Mənbə

Pages: 146-161
