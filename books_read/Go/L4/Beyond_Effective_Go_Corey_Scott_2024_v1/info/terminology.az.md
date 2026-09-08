# Beyond Effective Go — Part 2 — Terminoloji Lüğət (Azərbaycanca)

## A

**Abstract method pattern** — interface-də default (köməkçi) metod;
implementasiya yalnız məcburi metodu verir (Ch8).

**Adapter** — struct daxilində interface embed; xarici tipi istədiyin
formaya uyğunlaşdırır (Ch4).

**AST (go/ast)** — Go kodunun proqramatik təhlili üçün sintaksis ağacı (Ch9).

## B

**Behavior vs implementation test** — davranış (nəticə/effekt) yox,
daxili çağırış strukturu; sonuncu refaktoru partladır (Ch6).

**Be lazy** — productivlik fəlsəfəsi: avtomatlaşdır, alətləri öyrən,
təkrarı YOX (Ch7).

## C

**Clarity/Consistency/Predictability** — Code UX üç sütunu (Ch5).

**Closure** — əhatə etdiyi dəyişənlərə istinad saxlayan funksiya; loop
dəyişəni tələsinə diqqət (Ch8).

**Code UX** — kodun istifadəçi (developer) təcrübəsi; oxumaq = istifadə
etməkdir (Ch5).

**Composition over inheritance** — struct embed ilə davranış birləşdirmə;
miras hiyerarxiyası YOX (Ch4).

**Conventional Commits** — feat:/fix:/docs: prefiksli commit mesajları (Ch7).

**Currying** — çoxarqumentli funksiyanın tək-arqument zəncirinə çevrilməsi;
Go-da ağır və nadir (Ch8).

## D

**Decorator** — funksiyanı funksiya ilə bürüyen pattern:
func(http.Handler) http.Handler (Ch8).

**Delegation** — işi tanış komponentə ötürmə (embed + çağırış) (Ch4).

**DIP (Dependency Inversion)** — yüksək qatlar konkret DEYİL, interfeysə
asılıdır (Ch4).

**DRY vs KISS** — təkrarsızlıq vs sadəlik; DRY erkən və yanlış
tətbiq olunarsa zərərlidir (Ch4).

## E

**Empty struct (struct{})** — 0 bayt; signal kanalları: chan struct{} (Ch8).

**Early return** — nesting-i azaldan qoruma bəndi (Ch5).

## F

**Factory Method** — constructor funksiyası (NewX) / closure ilə obyekt
yaratma; class-hierarchy versiyası Go-da YOX (Ch4).

**Functional options** — Option variadic-ləri ilə konfiqurasiya;
AZ optionda sadə variantlar üstündür (Ch8).

**Futures** — asinxron nəticənin placeholder-ı; buffered channel + goroutine
(YOX: Go-nun idiomu birbaşa kanaldır) (Ch8).

## G

**go:generate + mockery** — interfeysdən mock/stub avtomatik generasiyası
(Ch6).

## H

**Higher-order function** — funksiya qəbul edən/qaytaran funksiya (Ch8).

## I

**ISP** — Interface Segregation: interfeysi istifadəçiyə görə böl (Ch4).

**Immutability** — FP konsepti; Go-da qismən (dəyişməz qəbul + kopya
qaytar) (Ch8).

## L

**Least astonishment** — predictability prinsipi: sürprizsiz API (Ch5).

## M

**Metaprogramming** — proqramın proqram haqqında işi: AST təhlili, kod
generasiyası, exec koordinasiyası (Ch9).

**Middleware** — handler-ı bürüyen zəncir funksiyaları (Ch8).

**Mock/Stub** — stub: hazır dəyər; mock: çağırış gözləntiləri yoxlanılır
(Ch6).

## N

**noCopy** — go vet copylock-ın tapması üçün Lock()/Unlock() marker metodu
olan struct (Ch8).

## O

**Observer** — hadisə abunəçiləri; Go-da channel mübadiləsi ilə (Ch4).

## P

**Principle of least astonishment** — gözləntiyə uyğun API (Ch5).

**Pure function** — yan təsirsiz; eyni girişlə eyni nəticə (Ch8).

## R

**Race detector (-race)** — data race aşkarlayıcısı; -count ilə təkrarla
(Ch6).

**Recorder (test)** — çağırışları qeyd edən yüngül test double; struktur
YOX davranış ardıcıllığı yoxlaması (Ch6).

## S

**Singleton** — tək nüsxə; sync.Once; amma explicit dependency adətən daha
yaxşıdır (Ch4).

**SRP** — Single Responsibility: hər şeyin tək səbəbi (Ch4).

**Staticcheck** — vet-dən geniş linter; SA/ST/S class-ları (Ch7).

**Stub** — bax Mock.

## T

**Table-driven test (TDT)** — ssenari cədvəli + t.Run subtest-lər (Ch6).

**t.Cleanup** — test-ömrü boyu təmizləmə; defer-in test ekvivalenti (Ch6).

**Test latch** — goroutine-in başlamasını gözləmək üçün wg.Add(1)/Done
idiomu (Ch6).

**Test pyramid** — 70/28/2 unit/integration/E2E balansı (Ch6).

**testdata/** — go tool-un İGNOR etdiyi fixture qovluğu (Ch6).

## U

**Unix Philosophy** — hər alət bir işi yaxşı görsün; birləşdir (Ch4).

## Z

**Zero value** — tipin hazır başlanğıcı; API-lər zero value ilə işləməli
(Ch5).
