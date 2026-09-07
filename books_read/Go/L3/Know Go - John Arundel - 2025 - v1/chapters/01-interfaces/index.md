# Chapter 1 — Interfaces (İnterfeyslər)

## Bu fəsil nədən bəhs edir?

Generics-ə giriş hazırlığı: specific vs generic programming, interface tipləri
(io.Writer, Write metodu imzası), implicit implementation (MyWriter nümunəsi —
"heç nə edən" Write), interface parametrləri, BogusWriter compile xətası,
polymorphism ("many forms" — bir funksiya, çox konkret tip; hələ YARANMAMIŞ
tiplərlə belə işləyir!), öz interface yaratmaq (Stringer/Stringify, fmt.Stringer),
interfeysin məhdudiyyəti — yalnız metod ADI+İMZASI, DAVRANIŞI yox;
method set limitasiyası (AddNumbers/AddFloats dublikat kodu — int-in METODLARI
YOXDUR → built-in tiplər üçün interfeys YOLU YOX), any (boş interfeys — hamı
implement edir; amma `x + y` compile xətası: "operator + not defined on
interface"; həmçinin x/y fərqli tiplə ola bilər — məntiqi xəta), type
assertions/switch (hər tip üçün case = yenə dublikat!), math paketinin
float64-only problemini (Pow/Abs/Max), generics tarixi (2009 launch + 24
saatda ilk generics şərhi — Ian Lance Taylor; sadəlik + sürətli compile +
backwards compatibility maneələri; ~10 il; nəticə: "a lot with a little"
dizayn).

## Əsas fikirlər

### 1. Specific vs Generic Programming
- **Specific:** `func PrintString(s string)` — konkret tip; başqa tip ötürsən
  compile xətası
- **Generic:** istənilən YAXUD müəyyən TİP DƏSTİ qəbul edən funksiya —
  "PrintAnything" arzusu. Parametr siyahısına NƏ yazmalıyıq? — sualın özü
  problemdir (cavab: generics, ch2)

### 2. Interface Tipləri — Metod Dəsti Müqaviləsi
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

func PrintTo(w io.Writer, msg string) {
    fmt.Fprintln(w, msg)
}
```
- Parametrin dəqiq (dinamik) tipi RUNTIME-da bilinmir; amma EHTİYAC
  bildirilir: "io.Writer implement edir"
- **Implement etmək = müvafiq metod dəstinə malik olmaq**

### 3. Implicit Implementation — Bəyans Yoxdur
```go
type MyWriter struct{}

func (MyWriter) Write([]byte) (int, error) {
    return 0, nil       // heç nə edir — amma İMZA DÜZGÜNDÜR!
}
```
MyWriter-in `Write([]byte) (int, error)` metodu VAR → avtomatik io.Writer-dır.
Digər metodları da ola bilər — Go YALNIZ Write-ın mövcudluğunu yoxlayır.

Uğursuz nümunə:
```go
type BogusWriter struct{}
PrintTo(BogusWriter{}, "...")
// cannot use BogusWriter{} (type BogusWriter) as type io.Writer:
// BogusWriter does not implement io.Writer (missing Write method)
```
**Məna:** funksiya io.Writer qəbul etməklə "Write ÇAĞIRACAĞAM" deyir — Write-suz
tip İŞLƏMƏZ, Go bunu ƏVVƏLCƏDƏN bilir.

### 4. Polymorphism — "Many Forms"
Niyə konkret tip YOX, io.Writer? Çünki BİR funksiya — ÇOX tip: `*os.File`,
`*bytes.Buffer`, MyWriter... PrintToFile / PrintToBuffer / PrintToBuilder
AİLƏSİ yazmaq əvəzinə TƏK funksiya.

**Ən gözəl:** hələ MÖVCUD OLMAYAN tiplərlə də işləyir — sadəcə Write metodu
olan HƏR HANSI gələcək tip!

### 5. Öz İnterfeysini Yaratmaq
```go
type Stringer interface {
    String() string
}

func Stringify(s Stringer) string {
    return s.String()
}
```
(Standart kitabxanada mövcuddur: fmt.Stringer.) Interfeys parametri = eyni
kod, çox dinamik tip.

**Vacib məhdudiyyət:** interfeys yalnız metodun ADINI və İMZASINI tələb edə
bilər — DAVRANIŞINI YOX. MyWriter heç nə etmədiyi halda "implement edir" —
qəbul olunur.

### 6. Method Set Limitasiyası — Add Nümunəsi
```go
func AddNumbers(x, y int) int          { return x + y }
func AddFloats(x, y float64) float64 { return x + y }
```
MƏNTİQ EYNİDİR (x + y), yalnız tiplər fərqlənir — "type system is hurting
us more than it's helping". AddInt64s, AddInt32s, AddUints... — BORING, bu
müddəa proqramçı olmağın mahiyyəti deyil!

**İnterfeys xilas edə bilərmi?** — XEYR: **int-in (və bütün built-in
tiplərin) METODLARI YOXDUR!** int+float64+friends-i əhatə edəcək method set
MÖVCUD DEYİL. (Add metodu tələb edib, struct-lar yaratmaq olar — amma
built-in rəqəmləri KƏNARDA qoyar → "most inconvenient limitation".)

### 7. any — Boş İnterfeys və Onun Uğursuzluğu
```go
func AddAnything(x, y any) any {
    return x + y
}
// invalid operation: x + y (operator + not defined on interface)
```
- any = heç bir metod tələbi → HƏR tip implement edir
- Amma: **`+` operatoru interfeysdə TƏYİN OLUNMAMIŞDIR** — x struct olsa
  "toplama" nə deməkdir? Go TƏHLÜKƏSİZLİK seçir: QADAĞA.
- **Subtil problem:** x int, y string çağırışı MÜMKÜNDÜR (hər ikisi any!) —
  məntiqsiz, amma compile keçərdi (runtime problem).

### 8. Type Assertions / Type Switch — Yenə Dublikat
```go
switch v := x.(type) {
case int:
    return v + y
case float64:
    return v + y
case ...
}
```
Vədli görünür, amma HƏR TİP üçün case yazmaq = əvvəlki dublikatın MASKASI.
"An interface is no use here."

### 9. Real Dünya Zərəri — math Paketi
Standart `math`: Pow, Abs, Max... — YALNIZ float64. Başqa tip istəyirsən:
convert → çağır → geri convert. "That's just lame." General-purpose paket
yazmağı çətinləşdirən struktural problem.

### 10. Go, Meet Generics — Tarix
> "Go was released on November 10, 2009. Less than 24 hours later we saw
> the first comment about generics." — Ian Lance Taylor, "Why Generics?"

**Niyə gecikdi?**
1. **Sadəlik:** Go sürətlə öyrənilməli; hər yeni şey = yeni öyrəniləcək
   sintaksis
2. **Backwards compatibility:** heç bir breaking change; yeni sintaksis
   mövcud proqramlarla konflikt etməməli — "That's hard!"
3. Uzun illər təkliflər: compiler sürəti / runtime performans / complexity /
   uyğunluq maneələrinə çırpıldı

**~10 il sonra nəticə:** "a very nice design. Like Go itself, it does a lot
with a little" — minimal sintakslə yeni proqramlaşdırma dünyası. Növbəti
fəsil: nə əlavə olundu + ilk generic proqramlar.

## Əsas terminlələr
- Specific vs Generic Programming — konkret tip / tip dəsti qəbul edən
- Interface Type — metod dəsti müqaviləsi (io.Writer)
- Dynamic Type — runtime-da müəyyən olunan konkret tip
- Implement — müvafiq metod dəstinə sahib olmaq
- Implicit Implementation — bəyansız; metod imzası kifayətdir
- Method Set — tipin tələb olunan metodları
- Polymorphism ("many forms") — bir kod, çox tip; gələcək tiplər də daxil
- Interface Parameter — interfeys tipli funksiya parametri
- fmt.Stringer — String() string tələbi (standart)
- İmza vs Davranış — interfeys imzanı yoxlayır, məqsədi YOX
- Built-in Tiplərin Metodsuzluğu — int/fload64 method set-i YOXDUR
- any (Boş İnterfeys) — sıfır metod; HƏR tip implement edir
- "Operator + Not Defined on Interface" — any üzərində riyaziyyat QADAĞA
- Heterogen Çağırış Riski — x, y fərqli tiplə ola bilər
- Type Assertion / Type Switch — dinamik tip yoxlaması; dublikat maskası
- math float64-only — general-purpose çatışmazlıq nümunəsi
- Backwards Compatibility — breaking change qadağası
- "24 Saat" Şərhi — 2009 launch + ilk generics tələbi (Ian Lance Taylor)
- "A Lot with a Little" — Go/generics dizayn fəlsəfəsi

## Praktik nəticə
(1) Çox tip üçün bir funksiya: interfeys parametri — metod dəsti tələb et;
implicit implementation sayəsində gələcək tiplər də uyğun gələcək. (2)
İnterfeys YALNIZ ad+imza tələb edir — davranış zəmanəti YOX (MyWriter dərsləri).
(3) Built-in tiplərlə (int, float64) interfeys İŞLƏMİR — metodları yoxdur;
`+` kimi operatorlar any-də MÜMKÜNSÜZ. (4) Type switch = dublikatın gözəl
paltarından başqa bir şey deyil — sayı saxlaya bilməz. (5) Bu məhdudiyyətlərin
həlli GENERICS-dir (ch2+): ~10 illik dizayn problemi, "a lot with a little"
həlli. (6) General-purpose paket yazarkən float64-only tələsindən qaç —
generics API-ləri artıq mövcuddur.

## Mənbə
Pages: 14-24 (PDF 15-25)
