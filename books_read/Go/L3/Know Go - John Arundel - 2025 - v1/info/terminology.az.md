# Know Go — Terminoloji Lüğət (Azərbaycanca)

## A

**Abstract type (abstrakt tip)** — funksiya body-sində istifadə olunan type
parameter (məs. `var top E`); instantiate olunanda konkret tip olur.

**any** — boş interfeys, hər tip implement edir; constraint kimi yalnız
ən zəif zəmanət verir (operatorlar YOX).

**any-nin limiti** — `x + y` any-də compile xətasıdır; operator üçün
daha dar constraint lazımdır.

## B

**Basic interface (əsas interfeys)** — yalnız metod elementləri olan
interface; həm adi tip, həm constraint kimi işləyə bilər.

**Bunch[E]** — kitabın əsas nümunə generic slice tipi: `type Bunch[E any] []E`.

## C

**cmp.Ordered** — standart paket constraint-i: bütün ordered built-in
tiplər + derived-ləri (<, >, min, max üçün).

**comparable** — predeclared constraint: == dəstəkləyən BÜTÜN tiplər
(sonsuz dəstə olduğundan Go-da ifadə olunmur — compiler daxili).

**Constraint (məhdudiyyət)** — type parameter üçün icazəli tiplər
dəsti; operator zəmanəti verir.

**Contention (rəqabət)** — mutex uğrunda gözləmə; goroutine-lər iş
əvəzinə lock gözləyir.

## D

**Data race (data yarışı)** — iki goroutine eyni data üzərində
konkurent yazı/oxu; crash və ya gizli pozuntu.

**Deadlock (ölüm kilidi)** — dövrəvi lock gözləməsi; iki goroutine
bir-birinin kilidini gözləyib əbədi bloklanır.

**Derived type (törəmə tip)** — `type MyInt int` kimi mövcud tipdən
yaradılan YENİ tip; named-only constraint-lərdən kənarda qalır → ~ lazımdır.

**Dynamic dispatch (dinamik yönləndirmə)** — runtime-da ad ilə funksiya
seçimi (FuncMap məşqi).

## E

**Empty type set (boş tip dəsti)** — heç bir tipin qane edə bilmədiyi
constraint (int ∩ string); instantiate xətası verir.

## F

**First-class function (birinci dərəcəli funksiya)** — dəyər kimi
ötürülən/qaytarılan funksiya.

**Fold/Inject** — Reduce-un başqa dillərdəki adları.

## G

**Generic tip** — type parameter alan tip; runtime-da yalnız
instantiate olunmuş formada mövcuddur.

**Goroutine leak (goroutine sızması)** — bloklanmış goroutine-in
yaddaşda əbədi qalması (channel-ə erkən çıxış halında).

## H

**Heterogeneous (qarışıq) kolleksiya** — []any: fərqli tiplər bir yerdə
(Bunch-dan FƏRLİ).

**Homogeneous (vahid) kolleksiya** — Bunch[int]: hamısı eyni tip.

## I

**Interface literal (interfeys literali)** — ad verilməmiş inline
interface; constraint kimi istifadə oluna bilər ([T ~int], [T interface{
Equal(T) bool }]).

**Instantiation (instansiasiya)** — generic tipin/funksiyanın konkret
tipə bağlanması; explicit ([int]) və ya inference ilə.

**Intersection (kəsişmə)** — çoxsətirli interface; tip HƏR elementi
qane etməlidir (io.Reader VƏ fmt.Stringer).

**Iterator (itorator)** — tələbə görə element verən funksiya:
`func(yield func(V) bool)`; Go 1.23.

## L

**Lazy evaluation (tənbəl hesablama)** — yalnız istənilən qədər
hesablama; iterator kompozisiyasının əsası.

**LIFO** — Last-In-First-Out; Stack-in prinsipi (Push/Pop).

## M

**mapFunc/keepFunc/reduceFunc** — kitabın adlandırdığı funksiya tipləri:
Map/Filter/Reduce üçün imzalar.

**Mutex (mutual exclusion)** — qarşılıqlı istisna kilidi; RWMutex:
RLock (paylaşılan oxu) + Lock (eksklüziv yazı).

## O

**Ordered type (sıralana bilən tip)** — <, >, <=, >= dəstəkləyən tip;
cmp.Ordered bunun constraint-idir.

## P

**Parameterised method (parametrləşdirilmiş metod) — QADAĞA** — metodun
öz type parameter-i ola BİLMƏZ; həll: səviyyəli generic funksiya.

**Predeclared (əvvəlcədən elan edilmiş)** — comparable kimi compiler
daxili identifikator (cmp paketi kimi Go-da yazıla bilmir).

## R

**Race detector** — `go test -race` / `go run -race`; data race-ləri
avtomatik aşkarlayan alət.

**RWMutex** — read-write mutex: eyni anda çoxlu oxu YOXSA tək yazı.

## S

**Set (çoxluq)** — unikal, sırasız elementlər; `map[E]struct{}`
üzərində tikilir.

**Stencilling (şablonlama)** — hər tip üçün ayrıca maşın kodu
yaratmaq (spray-paint metaforası); interface indirection YOX.

**Struct{} (zero-size type)** — sıfır yaddaş tutumlu tip; set value-si
kimi ideal.

**sync/atomic** — mutex-siz atomik əməliyyatlar (atomic.Uint64 .Add/.Load).

## T

**Type element (tip elementi)** — interface daxilində metod deyil, TİP
(məs. int, ~string, struct{...}).

**Type inference (tip çıxarımı)** — Go-nun E-ni çağırış parametrlərindən
müəyyən etməsi; `mapFunc[int]` mismatch xətası bunu göstərir.

**Type parameter (tip parametri)** — T: tip üçün placeholder; funksiya
parametrləri siyahısından əvvəl kvadrat mötərizədə.

**Type set (tip dəsti)** — constraint-i qane edən bütün tiplər.

## U

**Union (birləşmə)** — `|` ilə ayrılan tip elementləri: int | int8 | ...

## Y

**Yield (vermə)** — iteratorun növbəti dəyəri loop-a ötürməsi;
true = davam, false = loop bitdi (təmizlənmə şansı).

## ~

**~ (approximation/tilde)** — "underlying type T olan hər şey": ~int →
int + MyInt + ...; derived tipləri constraint-ə daxil edir.
