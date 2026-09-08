# Chapter 16 — Coins (səh. 268-277)

## Bu chapter nədən bəhs edir?

Kriptovalyuta və blockchain: "Bobcoin" düşüncə eksperimenti ilə mərkəzləşdirilmiş
bankın problemləri, distribüted ledger, double-spending hücumu, blok zənciri ilə
transaksiya sıralaması, konsensus (doomed stubs), 51% hücumu, proof-of-work və
proof-of-stake.

## Əsas fikirlər

### 1. Mərkəzi bank problemi (Bank of Bob)
**Ssenari:** Bob "Bobcoin" yaradır — Bank of Bob konversiya rate təyin edir
(1 Bobcoin = 1 USD), Alice $100 ödəyib 100 Bobcoin balansı alır.

**Mərkəzləşdirmənin zəiflikləri:**
1. **Yük:** Bütün dünya transaksiyaları Bob-un serverlərini yükləyir
2. **Single point of failure:** Bob-un elektriki çıxsa — qlobal ticarət dayanır
3. **Etibad problemi:** Bob-a niyə güvənək? O, oğurluq edə, təhdid edilə,
   açarları itirə bilər
4. **Sızma:** Kimsə Bob-un sistemlərinə giribsə — hamı əziyyət çəkir

### 2. Distributed ledger (paylanmış hesablaşma dəftəri)
**Həll:** Hər node ledger-in **tam nüsxəsini** saxlayır, transaksiyalar
yayımlanır.

**Eventual consistency:** Nəticə etibarlılığı "sonda" çatılır — hər node bir neçə
saniyə/dəqiqə ərzində bütün transaksiyaları alır, hamının balansı razılaşır.

### 3. Double-spending hücumu
**Problem:** Sam 10 Bobcoin var. Pizza alır (Pat-dan), sonra dərhal — şəbəkə
yayılmamışdan əvvəl — eyni 10 Bobcoin ilə Irene-dən dondurma alır. İrene
transaksiyaları gecikmə ilə və **bilinməyən sıra ilə** alır — hansı birinci?

**Həll yolu:** Transaksiya axınını bloklara böl → blokları **sıralamak** lazımdır
— dəyişdirilmiş/aralıqdan çıxarılmış blok aşkarlanabilsin. Bu, CBC zəncirinin
məntiqidir!

### 4. Blockchain = hash zənciri
**Necə işləyir:** Hər blok **əvvəlki blokun hashini** özündə saxlayır →
bərabər ardıcıllıqdan başqa düzgün sıralama mümkün deyil → Irene Sam-in
double-spend-ini aşkar edir. Uğur!

**Qalan problem:** İki node müstəqil olaraq fərqli "növbəti blok" hesablayıb
göndərirsə — zəncir **split** olar?

### 5. Konsensus və doomed stubs
**Necə işləyir:**
1. Node iki fərqli namizəd blok alır → **ikisini də saxlayır**, qərar gözləyir
2. Qısa müddətdə bir versiya üzərində **daha çox davam bloku** gəlir
3. "Qalib" zəncir davam edir; tək qalan budaq **doomed stub** kimi atılır
4. Şəbəkə tək avtoritet sıralamaya **tez yığılır** (converges)

**Mining fee (mükafat):** Bloku yaradan node az qala kiçik fee qazanır —
qalan fraksiyon Bobcoin-lər hər transaksiyadan. Fee-earning transaction öz
blokunda valideyt olunur.

### 6. 51% hücumu
**Problem:** "Hər node öz tapşırığını özü yoxlayır" — alkoqolik node saxta
blok yarada bilməzmi?

**Müdafiə:** Mallory saxta blok düzəldir, amma dürüst zəncir daha çox hashing
power daşıdığından **tez uzanır** → bütün node-lar qısa (saxta) zənciri atır.
Mallory ancaq şəbəkənin **51%-dən çox hash gücünü** ələ keçirsə qalib gələ
bilər — böyük şəbəkələrdə qeyri-mümkün qədər bahadır. **Nə qədər çox node →
o qədər təhlükəsiz** (safety in numbers).

### 7. Proof of Work (iş sübutu)
**Nədir:** Blok qəbul olunması üçün onun hashindən tələb: məs. hash-in
başlanğıcı müəyyən sayda sıfır olmalıdır. Node-lar **nonce** dəyişərək milyonlarla
hash hesablayır (mining) — sırf şans işidir.

**Təsirlər:**
- Saxta blok yaratmaq = eyni bahalı işi təkrar etmək → **daha baha başa gəlir
  ki, qazanc gətirir** (mining fee-ləri dürüst blokları mükafatlandırır)
- Bitcoin-də qəbul şansı təxminən **1 in 10^64**; hər ~10 dəqiqədə yeni blok
- Minlərlə node durmadan hash-ləyir — **astronomik enerji** sərf olunur
  ("planetə od vurulması")

### 8. Proof of Stake və qeydlər
**Alternativ:** Ən çox pulu olan node-lara güvən (böyük stake = sistemin
sökülməsinə ən az maraqlı). Amma demokratiklik baxımından tənqid: sistem ən
zənginlərin xeyrinə "qurulur" — bu artıq "kapitalizm" adlanan şeydir.

**Müəllifin qənaəti:** Planeti yandırmayan, ədalətli, səmərəli double-spend
qarşısı hələ yoxdur — "bəlkə bunu siz tapacaqsınız".

## Əsas terminlər

- Blockchain (blokzəncir)
- Distributed Ledger (paylanmış hesablaşma dəftəri)
- Eventual Consistency (sonlu tutarlılıq)
- Double-Spending (ikiqat xərcləmə)
- Doomed Stub (talehsiz budaq)
- Consensus (konsensus / razılıq)
- 51% Attack (51% hücumu)
- Proof of Work (iş sübutu)
- Nonce (təkrarlanmayan say)
- Proof of Stake (pay sübutu)

## Praktik nəticə

- CBC zənciri ilə blockchain eyni məntiqi paylaşır: hər blok əvvəlkinin
  hashinə bağlıdır → sıralama dəyişməz
- Double-spend-i sıralama ilə həll et; split-i "ən uzun zəncir qalib" ilə
- Mərkəzi etibad noktasını aradan qaldırmaq üçün konsensus + PoW/PoS
- 51% hashrate hücumu — yeganə praktik (böyük şəbəkədə qeyri-mümkün) hücum
- PoW təhlükəsizliyi enerjiyə əsaslanır; alternativlər hələ qeyri-kamildir

## Mənbə

Pages: 268-277 (Chapter 16, Explore Go: Cryptography)
