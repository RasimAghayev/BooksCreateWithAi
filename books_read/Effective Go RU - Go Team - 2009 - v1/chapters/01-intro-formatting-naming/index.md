# Chapter 1 — Giriş, formatlaşdırma və adlandırma (səh. 1-6)

## Bu chapter nədən bəhs edir?

Effective Go-nun giriş hissəsi: Go-nun fəlsəfəsi (C++/Java-dan birbaşa tərcümə YOX), gofmt formatlaşdırma yanaşması, kod commentləri, paket/funksiya/interfeys adlandırma konvensiyaları vəGo-nun avtomatik nöqtəli vergül (semicolon) qaydası.

---

## Əsas fikirlər

### 1. Go fəlsəfəsi — birbaşa tərcümə YOX

Go — yeni dildir: mövcud dillərdən ideya götürür, amma **effektiv Go proqramları xarakterinə görə fərqlidir**. C++/Java proqramının birbaşa tərcüməsi qənaətbəxş nəticə vermir — Java proqramı Java üçün yazılıb, Go üçün YOX. Problemi **Go prizmasından** düşünmək tamamilə fərqli, amma uğurlu proqrama gətirib çıxarır.

Effektivliyin açarı: (1) dilin özəllikləri və idiomları, (2) adlandırma/formatlaşdırma/struktur konvensiyaları — başqa Go proqramçıları üçün anlaşılan kod.

**2022 qeydi:** sənəd 2009 buraxılışı üçün yazılıb və ciddi yenilənməyib — dil stabil qaldığından özü aktual, amma ekosistem (build, test, modullar, polimorfizm) haqqında YOX (issue 28782).

### 2. Formatlaşdırma — maşina qərar verir

Format ən mübahisəli, amma ən az əhəmiyyətli mövzudur. Go-nun yanaşması: **gofmt maşını işə sal** — standart stildə çıxış verir (indent, vertikal düzləndirmə, commentlər). Sual yarananda: gofmt işlət; nəticə səhv görünürsə proqramı yenidən qur (və ya gofmt bug yaz), əksinə keçin.

**Nümunə — sahə commentləri avtomatik sütunlanır:**

```go
type T struct {
    name    string  // имя объекта
    value   int     // его значение
}
```

Standart paketlərin HAMISI gofmt ilə formatlaınıb.

**Qalan detallar:**
- **İndent:** TAB (gofmt default); space yalnız zərurətdə
- **Sətir uzunluğu:** limit YOX — "perforator kartından qorxma"; uzun sətir → köçür + əlavə TAB
- **Mötərizə:** if/for/switch-də mötərizə YOX; operator prioriteti qısadır: `x<<8 + y<<16` — interval göstərdiyi kimi oxunur

### 3. Commentlər

- Sətir commenti `//` — norma; blok `/* */` — əsasən paket sənədləri / böyük kod söndürmə
- Yuxarı səviyyə elanlardan ƏVVƏL (ara sətirsiz) comment = həmin elanın sənədi (**doc comment**) — paketin əsas dokumentasiyası

### 4. Adlandırma — semantik effekt!

Adın **ilk hərfi BÖYÜKDÜRSƏ** paketdən kənara görünür (export). Konvensiyalar:

**Paket adları:** qısa, tək söz, lowercase; alt xətt/mixedCaps YOX. `src/encoding/base64` → paket adı `base64` (`encoding_base64` YOX). Toqquşma nadirdir — import alias həlli: `strings2 "..."`.

**Export təkrarsızlığı:** `bufio.Reader` YOX `BufReader` (istifadəçi `bufio.Reader` görür — aydındır). `ring.Ring` üçün konstruktor sadəcə `New` → `ring.New`. internal/server nümunəsi: `server.New()` — `NewServer()` lazımsız.

**Getter/Setter:** `Get` prefiksi YOX — sahə `owner` → getter `Owner()`:

```go
owner := obj.Owner()
if owner != user {
    obj.SetOwner(user)
}
```

**İnterfeys adları:** 1-metodlu interfeyslər `metod-adı + -er`: `Reader`, `Writer`, `Formatter`, `CloseNotifier`. **Kanonik imzalar:** `Read/Write/Close/Flush/String` — eyni semantika daşımayan metoda bu adları VERMƏ; daşıyırsa eyni ad+imza (String YOX ToString).

**MixedCaps:** çoxsözlü adlar `MixedCaps`/`mixedCaps` — alt xətt YOX.

### 5. Nöqtəli vergül (semicolon) — lekserin qaydası

Go-da nöqtəli vergül mənbədə YOX — **lekser avtomatik daxil edir**: sətir sonundaki token id/literal/`break, continue, fallthrough, return, ++, --, ), }` → sonra `;` əlavə olunur.

**Nəticə:** if/for/switch/select-in açılış `{`-i NÖVBƏTİ SƏTİRƏ QOYMAQ OLMAZ (əvvəlinə `;` daxil olunar):

```go
if i < f() {     // DOĞRU
    g()
}
if i < f()       // YANLIŞ!
{                // YANLIŞ!
```

İdiomatik kodda `;` yalnız `for` init/cond/post və bir sətirdə çox operator ayrıcı kimi.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| gofmt | Standart stil verən formatter — format müzakirələrini bitirir |
| Doc comment | Elandan əvvəlki comment — paket sənədi |
| Export (böyük hərf) | İlk hərf böyükdürsə paketdən görünür |
| Kanonik imza | Read/Write/String-in standart mənası — eyni adda uyğunluq şərt |
| MixedCaps | Çoxsözlü ad konvensiyası (alt xətt YOX) |
| Semicolon insertion | Leksik qayda: sətir sonu tokeni → avtomatik `;` |
| Import alias | Ad toqquşması həlli (`strings2 "..."`) |

---

## Praktik nəticə

1. **Format müzakirəsi ETMƏ** — gofmt işlət; nəticə xoşuna gəlmirsə strukturu dəyiş.
2. **Adlar semantikdir:** böyük hərf = export; getter `Owner` (Get YOX); 1-metodlu interfeys `-er` suffiksi.
3. **Paket adı ilə export təkrarı:** `ring.New`, `bufio.Reader` — ad çərçivəsi istifadəçi tərəfindən tamamlanır.
4. **`{` eyni sətirdə** — lekserin `;` qaydası səbəbindən.
5. **C++/Java-dan köçürmə YOX** — problemi Go idiomları ilə yenidən qur.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi) — go.dev rəsmi dokumentasiya, 2009
- PDF səhifələri: 1-6
