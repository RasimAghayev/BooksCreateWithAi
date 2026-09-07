# Chapter 15 — The Tao of Go (Go-nun Tao-su)

## Bu fəsil nədən bəhs edir?

Kitabın fəlsəfi finalı. Tao = şeylərin daxili təbiəti/meyli (su aşağı axar —
bənd çək, pumpla, yenə gedəcəyi yerə gedər); "poor swimmer thrashes, Taoist
surfs" — dənəz qayıqla mübarizə YOX, dalğa ilə sürüş. Go-nun Tao-su: Dàodé
Jīng-in üç xəzinəsi — KINDNESS (insanlar üçün kod: istifadəçilərə — descriptive
adlar, asan import, yaxşı sənəd, cömüc lisenziyalar, dərin abstraksiyalar/kiçik
API; proqramı İŞLƏDƏNLƏRƏ — asan install/update, minimum konfiq, səhvləri
tutub DOSTCASA izah; kodu OXUYANLARA — clear/simple/explicit, mənalı adlar,
konvensiyalar, standart interfeyslər, "obvious thing the obvious way"; ÖZÜMÜZƏ
— əla testlər + mikro-refaktorlar: "spaghetti-yə çevrilən kod DÜZƏLƏ BİLMƏZ"
(Ousterhout)); SIMPLICITY (frugal/modest dil — minimal sintaksis; "Simplify,
simplify!" Thoreau; bir işi yaxşı görmək, sensible defaults, YARARSIZ
flexibility/extensibility QADAĞA — "sadə proqram mürəkkəbdən asan genişlənir");
HUMILITY (alçaq yerlər axtaran su kimi; Go-nun özü də humble — fəndlərdən
imtina, bizni xətalardan qoruyan dil; "ən təhlükəli xəta — öz xətaya meylimizi
görməmək" (Liddell Hart); clever olma, óbvio et, standart kitabxana birinci,
pre-engineering YOX, öz kodunu review et — "özümüz əziyyət çəkmiriksə, niyə
onlar?", sənədləşmə + işlək nümunələr); NOT STRIVING / wúwéi (tənbəllik DEYİL —
əksinə! "chronically busy" insanlar heç nə əldə etmir; ən yaxşı iş — gəzintidə/
verandada gəlir; it-qapı tələsi — PUSH YOX PULL; problem-solving < problem-
ELİMİNATİNG: "reframe et ki, problem YOX OLSUN"; "ən yaxşı optimizasiya — işi
HEÇ ETMƏMƏK"; programming ≠ typing — bəzən lazım olan yeganə düymə DELETE-dır;
büffalo hara istəyir orada getdər — istiqaməti dəyişmək üçün onu DİNİQ).
Sonda: "for the love of Go!" + About (müəllif, feedback, beta readers,
növbəti kitablar: Tools/Tests/Generics).

## Əsas fikirlər

### 1. Tao — Təbiətlə Mübarizəsizlik
**Tao** = daxili meyl. Su aşağı axar — bənd, kanal, pump heç nə dəyişməz:
yenə öz hədəfinə çatar. Taoçu yanaşma: təbii konturlarla İŞLƏ (wood grain),
əleyhinə DÖYÜŞMƏ. Yaxşı üzgüçü qulaq çalıb səs-küylə boğulur; Taoçu — SURF
edir (dalğanın gücündən istifadə). Go-nun Tao-su: dilin və problemin təbii
formasını görüb bulldozer YOX.

### 2. Kindness — İnsanlar Üçün Kod (Xəzinə 1)
**"Kod kompüterlər üçün YOX, insanlar üçün"** — insanlar səhvlə meylli,
sabırsız, təcrübəsiz, diqqətsizdirlər. Dörd qrup üçün mərhəmət:

**İstifadəçilərə (paket/API):**
- Descriptive adlar, asan import, yaxşı dokumentasiya, cömüd open-source
  lisenziya
- **Dərin abstraksiyalar:** kiçik sadə API → güclü davranış (deep module!)

**Proqramı işlədənlərə:**
- Asan install/update; minimum konfiqurasiya və asılılıq
- Adi istifadə səhvlərini tutmaq + DOST, dəqiq "nə xəta var, nəcə düzəlt" mesajı

**Kodu oxuyanlara:**
- Clear, sadə, EXPLICIT; kontekstdə mənalı adlar
- Abstraksiyalar düzgün kombinasiyalarda birləşsin
- Konvensiyalara bağlılıq, standart interfeyslər, "obvious thing the
  obvious way" — kognitiv yolları/əngəlləri LƏĞV

**Özümüzə (gələcək self):**
- Əla testlər (bu kitabın öyrəndiyi!) — gələcəkdə anlamaq/düzəltmək/
  təkmilləşdirmək asan; testsiz "things can go downhill fast"
- **Mikro-təkmilləşdirmələr:** hər ziyarətdə bir az refaktor — "spaghetti-yə
  çevrilmiş kod təmirsizdir" (Ousterhout); sıfırdan rewrite nadir halqdır —
  uzun müddət kiçik investisiyalar SAĞLAMLIĞIN yeganə praktik yoludur

### 3. Simplicity — Frugal Dildən Frugal Koda (Xəzinə 2)
"Üç xəzinəm: kindness, simplicity, humility" (Dàodé Jīng). **Go frugal dildir:**
minimal sintaksis, kiçik səth; hər şeyi etməyə, hər kəni razı salmağa
ÇALIŞMIR.

**"Simplify, simplify!" (Thoreau — Walden):**
- Proqramlar kiçik, fokuslu, clutter-sız — BİR işi yaxşı görsün
- Dərin abstraksiyalar = güclü maşın + sadə interfeys
- İstifadəçiyə "paperwork" YOX — API çağırmaq üçün icazə QAZANMAQ yox
- Sensible defaults — ən adi hallar üçün sadə yol
- **Flexibility tələsi:** HƏR halı, HƏR feature-u dəstəkləmə YOX;
  **Extensibility tələsi:** hələ lazım olmayan üçün sadə dizaynı POZMA
  YOX — "sadə proqram mürəkkəb olandan asan genişlənir"

### 4. Humility — Suyun Aşağı Yolu (Xəzinə 3)
Taoçu su kimi AŞAĞI yeri axtarır — mübarizəsiz, yarışsız, təsir etmədən.
**Go özü humble:** high-tech feature-ləri, teoretik üstünlükləri YOXDUR;
digər dillərin satış nöqtələrini BİLƏRƏKDƏN buraxıb — təsirli dil/populyarlıq
yarışı yerinə "kiçik sadə alət, faydalı iş, praktik yol".

**Go bizni qoruyur:** yaddaş allocatıon, təmizlik, unused import/var
xəbərdarlıqları — "özünü hər şeyi bilməyən, xətaya meylli bilən insanlar üçün
dil" — yəni HUMBLE insanlar üçün.

**"Ən təhlükəli xəta: öz xətaya meylimizi TANIMAMAQ" (Liddell Hart)**

**Humble Go proqramçısı:**
- **Clever olma** — kiminsə təsir etmək üçün YOX; óbvio et; şəxsiyyəti koda
  SOXMƏ
- Standart kitabxana həll edirsə — O; 3-cü tərəf yalnız işləmirsə; de-facto
  standart varsa — "başqalarına bəsdirsə, mənə də bəsdirir"
- **Pre-engineering YOX:** gələcəyi proqnozlaşdıra bilmirik — ehtiyatsız
  şeylər üçün vaxt/enerji ISRAF; digər proqramçıların nə istədiyini də bilmirik
  → hardwired asılılıqlar YOX
- Hər koddа BUG var fərziyyəsi → unexpected davranışı TUTAN testlər
- Status quo üçün həddindən artıq optimizasiya YOX (işlər boşa gedəcək)
- **Öz kodunu review ETMƏDƏN başqasına vermə** — "mən əziyyət çəkmirəmsə,
  onlar niyə çəksin?"; sətir-sətir, YENİ developer gözü ilə: haradan başlamaq?
  vacib tiplər əvvəldə? adlar hələ də dəqiqdirmi? struktura nezət saxlayır?
- **Genius deyilik** — self-explanatory kod yazmaq cəhdi YOX → İZAHLARA əmək:
  nə YOX, NECƏ istifadə; scratch-dən realistik nümunələr; gözlənilən nəticə
  + növbəti addım; nümunələri REGULAR yoxla (işlədiyinə əmin ol)

### 5. Not Striving — Wúwéi (推 YOX, Çəkmə)
**Tələq:** tənbəllik/pasivlik DEYİL — tam ƏKSİNƏ! "Çox çalışmaq" ≠ "çox yaxşı
iş". Hər kəs tanıyır: chronic busy, həmişə tələsik, fury-activity — amma HEÇ
NƏ əldə etməyən (və bunu bilən, buna görə də faciəvi vəziyyətdə olan) insanlar.

**Ən yaxşı iş "heç nə edərkən" gəlir:** çay kənarında gəzinti, verandada
hörümçək toru izləmə. "Bir dəqiqə zorlamayı dayandırmaq qədər ağıllı olsaq —
düz ideya dərhal gələr."

**İt-qapı dərsi:** inadlı PUSH edən, sonra qapının PULL olduğunu anlayan —
hər kəs yaşayıb. **Gündəlik işdə hansı "pull" işarələrini itiririk?**

**Problem-eliminating > problem-solving:**
- Reframe: problemi elə ifadə et ki, YOX OLSUN
- Requirement-ləri elə yenidən yaz ki, həll TRIVIAL/óbvio olsun
- İrrelevant detalın arxasınca qaçıb GÖRMƏDİYİMİZ sadə-elegant dizayn varmı?
- "Bu problemlə ÜMUMİYYƏTLƏ məşğul olmaya bilərəmmi?"
- **"Ən yaxşı optimizasiya — işi HEÇ ETMƏMƏK"**

**Programming ≠ Typing:**
- Boşluğa baxan adam "işləmir" görünür; klavyəni döyən "işləyir" — İKİSİ DƏ
  YANLIŞ fərziyyə!
- Real proqramlaşdırma YAZDAN ƏVVƏL (və bəzən YAZIN YERİNƏ) baş verir
- **"Əla proqramlaşdırma parçasından sonra lazım olan yeganə düymə — DELETE"**

**Buffalo qanunu (Weinberg):** "Buffalonu istədiyin yerə ötürə bilərsən —
yalnız ora GETMƏK İSTƏSƏ." Zorlamayı dayandır: buffalo HARA İSTƏYİR? Bəlkə o
yer daha yaxşıdır!

## Əsas terminlələr
- Tao — şeylərin daxili təbiəti/meyli
- "Water Flows Downhill" — mübarizəsızlıq metaforası
- "Thrashing vs Surfing" — mübarizə / dalğadan istifadə
- Kindness / Compassion — insan-fokuslu kod (Xəzinə 1)
- Deep Abstraction — kiçik API, güclü maşın (kitabın Tools hissəsində dərin)
- Friendly Errors — "nə xəta var, necə düzəlt" mesajı
- Obvious Thing the Obvious Way — konvensiya prinsipi
- Micro-Improvements — hər ziyarətdə kiçik refaktor
- "Spaghetti İmkansız Təmir" — Ousterhout qanunu
- Simplicity / Frugality — azla çox (Xəzinə 2)
- "Simplify, Simplify!" — Thoreau sitatı
- Sensible Defaults — adi hal üçün qısa yol
- Extensibility Tələsi — lazımsız genişlənmə üçün sadəliyi pozmaq
- Humility — alçaq yol, su kimi (Xəzinə 3)
- "Humble Go" — feature-lərdən bilərəkdən imtina
- Liddell Hart Qanunu — öz xəta meylini görməmək = ən təhlükəli
- Pre-Engineering — gələcək üçün erkən qurmaq — ISRAF
- "Özünü Review Et" — başqasından ƏVVƏL öz kodunu oxu
- Genius Fikri — self-explanatory YOX, izah + nümunələr
- Wúwéi / Not Striving — zorlamamaq (tənbəllik DEYİL)
- "Push-Pull Qapı" — yanlış istiqamətdə zorlama dərsi
- Problem-Eliminating — problemi ləğv etmək > həll etmək
- "Ən Yaxşı Optimallaşdırma — Etmemək" — 0 iş = 0 xəta
- Programming ≠ Typing — fikir klaviaturadan əvvəl
- "Yeganə Düymə — DELETE" — yazaq az, silək çox
- Buffalo Qanunu — Weinberg: istiqaməti istifadəçi istəyi müəyyən edir
- Dàodé Jīng — üç xəzinə mənbəyi

## Praktik nəticə
(1) Kodu insanlar üçün yaz: istifadəçi (asan API + sənəd), işlədən (asan
install + dostcasə xəta), oxuyan (óbvio + konvensiya), ÖZÜN (test + kiçik
refaktorlar). (2) Bir işi yaxşı; defaults; hər feature-u dəstəkləmə — "sadə
proqram mürəkkəbdən genişlənir". (3) Clever YOX — óbvio; standart kitabxana
əvvəl; pre-engineering və erkən optimizasiya ISRAF. (4) Öz kodunu review et
— sətir-sətir, yeni göz ilə; izah + realistik nümunələr yaz, nümunələri
regular yoxla. (5) Zorlanma (wúwéi): problem görəndə əvvəl SORUŞ — problemi
ləğv etmək olarmı? Reframe? Heç etməmək? "Ən yaxşı optimallaşdırma = işi
etməmək". (6) Fikir gəlmirsə — klavyədən uzaqlaş (gəzinti, veranda); real
proqramlaşdırma typingdən əvvəl gəlir; əla kod çox vaxt SİLİNMİŞ koddur.
(7) İnadla push etdiyin qapılar — pull olsun? Buffalo hara istəyir? — işi
onun istiqamətinə çevir. (8) Kitabın məqsədi: master proqramçı yox — Tao-su
təbii Go yazan, dilin dənəzinə uyğun işləyən mühəndis. "For the love of Go!"

## Mənbə
Pages: 185-197 (PDF 186-198)
