# 4. Cracking

**Səhifələr:** 74-84

## Bu fəsil nədən bəhs edir?

Şifrələnmiş mətni açarını tapmadan ya da məzmumunu bilmədən açmağa çalışan `crack` funksiyasının yazılması. Test-driven yanaşma ilə həyata keçirilir.

## Əsas fikirlər

### 1. Crack Function
`Encipher` və `Decipher` funksiyalarına bənzər, lakin açarı məlum deyil. Yalnız şifrələnmiş mətn və mümkün açar dəyərləri verilir.

### 2. Score-based Cracking (Skor əsaslı kırma)
Mümkün açar dəyərlərinin hansının düzgün olduğunu təyin etmək üçün skorinq (scoring) sistemi işlətirlər.

### 3. Subtests
`testing` paketində subtests (`t.Run`) vasitəsilə testləri daha strukturlaşdırmaq.

## Əsas terminlər
- Cracking (Şifrəni kırma)
- Score (Skor)
- Subtests (Alt testlər)
- Test Table (Test cədvəli)

## Praktik nəticə

`crack` funksiyası üçün testlər yazılır və funksiya tərtib edilir.

## Mənbə

Pages: 74-84
