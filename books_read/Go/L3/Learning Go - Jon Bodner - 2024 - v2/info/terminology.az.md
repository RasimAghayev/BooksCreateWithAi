# Learning Go, Second Edition — Terminologiya (Azərbaycanca)

Format: **English Term (Azərbaycanca qarşılıq)** — qısa izah

## Mühit və Alətlər (Ch1)
- **Workspace (iş sahəsi)** — üçüncü tərəf Go alətlərinin saxlandığı `$GOPATH` qovluğu
- **Linter (kod təmizləyici)** — stil qaydalarını yoxlayan alət (golint, golangci-lint)
- **Vet (yozyoxlama)** — məntiqən şübhəli, amma leqal kodu tutan go vet aləti
- **Language Server (dil serveri)** — editərlər üçün ağıllı kod xidmətləri standartı (gopls)
- **Compatibility Promise (uyğunluq vədi)** — Go 1.x boyu backward-compatibility zəmanəti

## Tiqlər (Ch2-3)
- **Zero Value (sıfır dəyəri)** — təyin edilməmiş dəyişənin default dəyəri (0/false/"")
- **Literal (literal)** — kodda yazılı dəyər; untyped — tipi kontekstdən gəlir
- **Rune (simvol kodu)** — int32 alias; tək Unicode code point
- **Explicit Type Conversion (açıq tip çevirməsi)** — `T(v)`; avtomatik promotion YOX
- **Typed/Untyped Constant (tipli/tipsiz konstant)** — const-un tip məcburiyyəti var/yox
- **Backing Store (arxa anbar)** — slice-in altında duran array
- **Capacity (tutum)** — slice üçün ayrılmış yaddaş ölçüsü (cap)
- **Full Slice Expression (tam slice ifadəsi)** — `x[a:b:c]` — subslice tutum məhdudiyyəti
- **Comma-ok İdiom (vergül-ok idiomu)** — `v, ok := ...` dəyər+mövcudluq oxuması
- **Hash Map (xəş xəritəsi)** — açarı hash edərək bucket-də saxlayan map implementasiyası
- **Collision (toqquşma)** — fərqli açarların eyni bucket-ə düşməsi

## Kontrol və Quraşdırma (Ch4-5)
- **Shadowing (kölgələmə)** — daxili blokda eyni adlı identifier-in xaricini gizlətməsi
- **Universe Block (kainat bloku)** — predeclared identifier-lərin (nil, int, true) bloku
- **Blank Switch (boş switch)** — müqayisə dəyərsiz, bool case-li switch
- **Variadic (çoxsaylı)** — `...T` ilə istənilən sayda arqument
- **Named Return Value (adlı qaytarma)** — funksiya imzasında adlandırılmış nəticə
- **Blank/Naked Return (boş qaytarma)** — dəyərsiz `return` — QADAĞAN
- **Closure (bağlanma)** — xarici dəyişənləri capture edən daxili funksiya
- **Call By Value (dəyərlə çağırma)** — parametr həmişə kopya kimi ötürülür

## Pointer-lər (Ch6)
- **Dereference (dəyərə keçid)** — `*p` ilə pointer hədəfinin oxunması
- **Address Operator (ünvan operatoru)** — `&x` — dəyişənin yaddaş ünvanı
- **Escape Analysis (çıxış analizi)** — datanın stack/heap qərarı
- **Mechanical Sympathy (mexaniki uyğunluq)** — hardware-ə uyğun kod (ardıcıl yaddaş)
- **Buffer (tampon)** — təkrar istifadə üçün bir dəfə ayrılmış slice

## Metodlar və İnterfeyslər (Ch7)
- **Receiver (qəbul edici)** — metodun bağlılıq parametri `(p Person)`
- **Method Set (metod dəsti)** — tipin/pointer-in metodlarının toplusu
- **Method Value/Expression** — metodu dəyər kimi / tip üzərindən funksiya kimi saxlama
- **Embedded Field (daxil edilmiş sahə)** — adsız sahə; promote mexanizmi
- **Promotion (təbliğ)** — embedded sahənin metodlarının xarici struct-a keçməsi
- **Implicit Interface (implicit interfeys)** — bəyansız method-set uyğunluğu
- **Duck Typing (ördək yazısı)** — "metodları varsa, odur" dinamik yanaşma
- **Type Assertion (tip iddiası)** — `i.(T)` runtime yoxlaması
- **Type Switch (tip açarı)** — interface dəyərinin tipinə görə switch
- **Optional Interface (seçimə bağlı interfeys)** — assertion ilə aşkarlanan əlavə imkan
- **Dependency Injection (asılılıq inyeksiyası)** — asılılıqların xarici təminatı
- **Decorator Pattern (bəzədici nümunəsi)** — eyni interfeysi qaytaran sarğı

## Xətalar (Ch8)
- **Sentinel Error (gözcü xətası)** — konkret dəyərlə davamsızlıq bildirən error
- **Error Wrapping (xəta bükülməsi)** — `%w` ilə orijinalı saxlayan kontekst əlavəsi
- **Error Chain (xəta zənciri)** — wrap olunmuş errorlar ardıcıllığı
- **Golden Path (qızıl yol)** — xətasız əsas icra yolu
- **panic/recover** — dayanma mexanizmi və yalnız defer-də tutulması

## Modullar (Ch9)
- **Module Path (modul yolu)** — modulun qlobal unikal identifikatoru
- **Package Clause (paket bəndi)** — faylın `package NAME` sətri
- **Blank Import (boş import)** — `_ "pkg"` — yalnız init-i işə salır
- **Type Alias (tip aliası)** — `type B = A` — eyni tipin yeni adı
- **Semantic Versioning (semantik versiyalama)** — vMAJOR.MINOR.PATCH
- **Minimum Version Selection** — ən yeni tələb olunan asılılıq versiyası seçilir
- **Import Compatibility Rule** — minor/patch 100% backward-compatible olmalı
- **Vendoring (kənarlaşdırma)** — asılılıqların vendor/ qovluğunda saxlanması
- **Sum Database (cəm bazası)** — imzalı modul hash-lərinin mərkəzi qeydi

## Paralellik (Ch10)
- **CSP (Communicating Sequential Processes)** — Go-nun paralellik modeli (Hoare 1978)
- **Goroutine** — runtime idarəli yüngül icra vahidi
- **Channel (kanal)** — goroutine-lərarası ünsiyyət quruluşu
- **Buffered/Unbuffered Channel** — buferli/bufersiz kanal
- **Starvation (aclıq)** — select-də həmişə eyni case-in qalib gəlməsi problemi
- **Deadlock (ölüm kilidi)** — qarşılıqlı gözləmə; hamı yatırsa runtime öldürür
- **Done Channel (bitiş kanalı)** — close ilə çıxış siqnalı verən struct{} kanalı
- **Backpressure (geri təzyiq)** — daxil olan işin token-bufer ilə məhdudlaşdırılması
- **Goroutine Leak (goroutine itkisi)** — əbədi çıxa bilməyən goroutine
- **Critical Section (kritik bölmə)** — mutex ilə qorunan kod/dəyər
- **Reentrant Lock (yenidən giriş kilidi)** — Go-da OLMAYAN xüsusiyyət
- **RWMutex** — oxu paralel, yazma tək sahibli kilid

## Standart Kitabxana (Ch11-12)
- **io.Reader/io.Writer** — buffer parametrli oxuma/yazma interfeysləri
- **io.EOF** — axın sonu sentinel error-u (data ilə birgə gələ bilər)
- **Monotonic Clock (monoton saat)** — sıçramayan, boot-dan sayılan vaxt
- **Reference Time (istinad vaxtı)** — "Jan 2 15:04:05 2006 MST" format şablonu
- **Marshaling/Unmarshaling** — Go tipindən JSON-a / əksinə çevrilmə
- **Struct Tag (strukt etiketi)** — `` `json:"name"` `` metadata sətri
- **Middleware (ara proqram)** — `func(http.Handler) http.Handler` sarğı funksiyası
- **ServeMux (sörvyu mükss)** — standart HTTP router
- **Context (kontekst)** — request metadata; ilk parametr konvensiyası
- **CancelFunc (ləğv funksiyası)** — context-i ləğv edən; defer ilə mütləq çağırılır
- **Done Channel** — context ləğvində bağlanan kanal
- **Unexported Key Pattern** — context dəyəri üçün unexported tip/konstanta açarı

## Testlər (Ch13)
- **Table Test (cədvəl testi)** — anonim struct slice + t.Run subtestləri
- **Code Coverage (kod örtüyü)** — testlərin toxunduğu kod nisbəti
- **Benchmark (sürət sınağı)** — b.N looplu performans ölçmə funksiyası
- **Stub (stAB)** — hazır dəyər qaytaran saxta asılılıq
- **Mock (mok)** — çağırış ardıcıllığını təsdiqləyən saxta
- **Build Tag (düzmə etiketi)** — `// +build` fayl səviyyəli kompilyasiya şərti
- **Data Race (data yarışı)** — kilidsiz paralel yaddaş çıxışı
- **httptest** — random portlu test HTTP server paketi

## Dragonlar (Ch14-15)
- **Reflection (refleksiya)** — runtime tip/dəyər manipulyasiyası
- **Type/Kind/Value** — refleksiyanın 3 sütunu (ad / quruluş / dəyər)
- **unsafe.Pointer** — hər tip pointer üçün körpü tipi
- **uintptr** — pointer riyaziyyatı üçün tam ədəd (GC-dən qorunmur)
- **runtime.KeepAlive** — GC-nin vaxtından əvvəl toplamasının bloku
- **cgo** — Go↔C FFI (xarici funksiya interfeysi)
- **Wire Format (şəbəkə formatı)** — big-endian bayt sırası
- **Type Parameter (tip parametri)** — `[T any]` generik dəyişəni
- **Type List (tip siyahısı)** — interfeysdaxili operator icazəsi siyahısı
- **Sum Type (cəm tipi)** — məhdud tip dəsti (gələcək istiqamət)
