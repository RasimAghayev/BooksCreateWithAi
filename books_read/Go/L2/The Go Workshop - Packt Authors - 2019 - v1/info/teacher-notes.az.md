# The Go Workshop — Müəllim Qeydləri (Azərbaycanca)

## 📖 Kitab deyir:

Öyrənmə WORKSHOP metodu ilədir — hər anlayış Exercise (addım-addım) və Activity
(müstəqil) ilə möhkəmlənir. Kitabın tədris fəlsəfəsi: "Learning by doing is a
fundamental principle" — passiv oxu əvəzinə KOD YAZARAQ öyrənmə. 19 fəsil
tədris planı: fundamentallar → abstraksiya (funksiya/error/interface) →
modul/paket → real dünya (vaxt/JSON/fayl/DB/HTTP) → concurrency → alətlər →
security → xüsusiyyətlər. Hər fəsilın əvvəlində "By the end of this chapter, you
will be able to..." — aydın tədris məqsədləri.

Kitabın məqamlı tövsiyələri:
- "A function should perform one task" — single responsibility ~25 sətirdə
- "Do not use any of its features prematurely" — paketi/interfeysi ehtiyac
  olanda yarat
- "Security cannot be an afterthought"
- "Share by communicating, do not communicate by sharing"
- "Premature optimization is the root of all evil"

## 👨‍🏫 Müəllim qeydi:

Kitab əla başlanğıc dərslikdir, amma praktikada bir neçə cəhəti tamamlamaq lazımdır:

1. **Kitab Go 1.12/1.13 dövründə yazılıb (2019) — bəzi tövsiyələr köhnəlib.**
   - `ioutil.WriteFile/ReadFile` — Go 1.16-dan `os.WriteFile/os.ReadFile` əvəz
     edir; ioutil DEPRECATED-dir
   - `rand.Seed()` — Go 1.20-dən lazımsız; kitabın şablonu hər yerdə Seed çağırır
   - Go 1.22 ServeMux metod+wildcard pattern-ləri (`GET /users/{id}`) kitabın
     routing bölməsini sadələşdirir
   - `// +build` — yeni forma `//go:build` (Go 1.17+); gofmt avtomatik çevirir

2. **Xəta interaktivliyi zəif qalır.** Kitab panic-i öyrədir, amma müasir
   error-wrapping (`fmt.Errorf("%w")`, errors.Is/As) tam öyrədilmir. Tədrisdə
   Ch6-ya 30 dəqiqə %w/Is/As əlavəsi tövsiyə olunur.

3. **Concurrency fəsli (Ch16) müasir Go-da az-çox dəyişib:** Go 1.22-də loop
   dəyişəni tələsi avtomatik həll olunur — kitabın parametr-ötürmə tövsiyəsi hələ
   də good practice-dir, amma RACE ARTIQ yaranmır. `errgroup` paketi kitabda YOX —
   prakitkada WaitGroup-dan daha çox istifadə olunur.

4. **Security fəsli aktualdır, amma səthi:** OWASP Top 10 xatırlanır, amma CSRF,
   authorization, secrets management (env vs fayl) yoxdur. bcrypt tövsiyəsi
   DÜZGÜNDÜR (cost 10-12); MD5-ün "checksum üçün YOX" statusu da dəqiqdir.

5. **Activity həlləri appendix-də (səh. 683-783) — 100 səhifəlik HƏLL SƏNƏDİ.**
   Bu, tədris üçün əla: tələbə özü cəhd edir, sonra müqayisə edir. Raw/ fəsillərdə
   fərddir: həllər ayrıca cavab bölməsindədir, chapter index.md-lərimizdə YOX.

6. **PostgreSQL quraşdırma tələbi (Ch13) Windows-da əngəlidir** — SQLite ilə
   əvəzetmə mümkündür; database/sql API eyni qalır, yalnız driver fərqlənir.
   Amma kitabın PostgreSQL seçimi real production-a daha yaxındır.

7. **cgo bölməsi (Ch19) Windows-da MinGW tələb edir** — tədrisdə optional
   saxla; əsas mesaj (unsafe yalnız son çarə) onsuz da çatdırılır.

## Ən vacib 5 fikir

1. Workshop metodu: anlayış → Exercise → Activity — ən yaxşı öyrənmə tsikli.
2. Go-nun sadəliyi: 25 keyword, implicit interfeys, error-as-value — sistem
   BÜTÜN proqramçılar üçün oxunaqlıdır.
3. Slice/channel daxili mexanika — bug-ların 90%-nin kökü; kitab bunu ƏN YAXŞI
   izah edir (noLink ssenariləri).
4. Concurrency birinci sinif vətəndaşdır — amma race/atomic/mutex/channel
   VOKABULARI olmadan istifadə ETMƏ.
5. Alət+security mədəniyyəti: vet/-race/html-template/Prepare/bcrypt —
   professional Go-nun minimum standartı.

## Kitabın ən dəyərli hissəsi

Chapter 4 (Complex Types) — slice hidden array mexanizması 6 ssenari ilə (linked,
noLink, capLink, capNoLink, copyNoLink, appendNoLink) — Go tədrisində ən yaxşı
bölmələrdən biri. Chapter 16 (Concurrent Work) — race-ən atomic-ə mutex-ə
channel-ə context-ə tam praktik yol. Bu ikili + Ch7 interfaces = Go biliyinin
ümumi strukturu.

⚠️ Uyğunsuzluq qeydləri: (1) Ch3 raw-string nümunəsində RTL simvollar (Fars)
gorünüşü pozur — nümunə mətnində çap səhvi var. (2) Ch11 `json:"enrolled, omitempty "`
tagində fazlalıq boşluq var — go vet bu tag-i xəta kimi tutar. (3) Kitabın
təkrar icra ssenarilərində "table already exists" xətaları NÖVBƏTİ icrada baş
verir — həllərdə DROP IF EXISTS istifadə olunmayıb, bu İRƏLİLİYƏN klasiikadır.
