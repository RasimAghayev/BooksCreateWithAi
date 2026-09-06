# Terminologiya — Go in Action, Second Edition

Bu fayl kitabda işlənən əsas texniki terminləri Azərbaycan dilində izah edir.

## A

- **Array** (Massiv) — Sabit ölçülü, eyni tip elementlər kolleksiyası
- **Alias** (Ləqəb) — `byte` (`uint8`), `rune` (`int32`) kimi başqa tipin adı

## B

- **Bool** (Məntiqi tip) — `true`/`false` dəyərlərini saxlayan tip
- **Byte** (Bayt) — `uint8` alias, 8-bit data (məlumat) üçün
- **Buffered I/O** (Buferli G/Ç) — `bufio` paketi, data-nı buferlər vasitəsilə emal edir

## C

- **Channel** (Kanal) — Goroutine-lər arasında data mübadiləsi üçün tipli strukturlar
- **Complex Number** (Kompleks Ədəd) — Real və imagery hissələrə malik tiplər (`complex64`, `complex128`)
- **Compile** (Kompilyasiya) — Source kodu executable binary-yə çevirmək
- **Concurrency** (Paralellik) — Eyni anda bir neçə tapşırığı idarəetmək qabiliyyəti
- **Constructor** (Konstruktor) — `NewT()` kimi tipi initialize edən funksiya
- **Cobra** (İlan) — Bitwise operator `&^` (bit clear)

## D

- **Dereference** (Göstəricini İzləmə) — `*pointer` ilə göstəricinin göstərdiyi dəyəri almaq
- **Dynamic Typing** (Dinamik Tip) — Tip runtime-də müəyyən olunan sistem (JavaScript, Python)
- **Dependency** (Asılılıq) — Xarici kitabxana və ya paket

## E

- **Error Handling** (Xəta İdarəetməsi) — `error` tipi ilə xətaları idarəetmə, `if err != nil` pattern-i
- **Export** (İxrac) — Açar söz ilə başlayan identifier (böyük hərf), başqa paketlərdən əlçatan
- **EOF** (End of File / Fayl Sonu) — Faylın sonu, `io.EOF`

## F

- **Float** (Həqiqi Ədəd) — Kəsr hissəsi olan ədəd tipi (`float32`, `float64`)
- **For Loop** (Döngü) — Go-da yeganə döngü konstruksiyası, 3 variety (növ) var
- **Fuzzing** (Qarışdırma) — Automate (avtomatik) testləşdirmə texnikası, Go 1.18+ daxildir

## G

- **Garbage Collection** (Zibil Toplayıcı) — Avtomatik yaddaş idarəetmə sistemi
- **Generics** (Generik Proqramlaşdırma) — Type parameter ilə universal funksiyalar
- **Git** — Distributed (paylanmış) version control (versiya nəzarəti) sistemi
- **Go Module** (Go Modulu) — Asılılıqları idarəetmək sistemi, `go.mod` faylı
- **Goroutine** (Gözaltı Proses) — `go` açar sözü ilə işə salınan concurrency vahidi
- **GOPATH** — Köhnə Go workspace (iş sahəsi) dəyişeni, indi `GO Modules` əvəzinə işlənir
- **GOPATH** — Köhnə Go workspace (iş sahəsi) dəyişeni, indi `GO Modules` əvəzinə işlənir

## I

- **Immutable** (Dəyişməz) — Dəyəri dəyişməyən tip (məsələn, string)
- **Import** (Daxil Etmə) — Başqa package-i proyektə daxil etmək
- **Implicit Interface** (Daxili İnterfeys) — Tip interfeysi implement edərsə, avtomatik olaraq ötənə bilər
- **Inference** (Çıxarma) — Tipin dəyərdən avtomatik müəyyən edilməsi
- **Init Statement** (Başlanğıc İfadəsi) — `for` döngüsünün ilk hissəsi
- **Interface** (İnterfeys) — Davranış təyin edən tip, metodlar qrupu
- **io.Reader** — Oxuma interfeysi, `Read(p []byte)` metodu
- **io.Writer** — Yazma interfeysi, `Write(p []byte)` metodu
- **Integer** (Tam Ədəd) — Signed/unsigned, müxtəlif ölçülü (`int`, `int64`, `uint` və s.)
- **Interpreted String** (Tərcümə Edilmiş Mətn) — `"Hello"` kimi, escape sequence-ləri işləyir

## K

- **Keyword** (Açar Söz) — Go-nun rezervi sözləri (`func`, `var`, `return`, `go` və s.)
- **Kebab** — Bitwise operator `&^` (bit clear)

## L

- **Literal** (Sabit Dəyər) — Mənbə kodda birbaşa yazılmış dəyər (`5`, `"hello"`, `true`)
- **Loop** (Döngü) — Təkrarlanan əməliyyat, `for` konstruksiyası

## M

- **Map** (Xüritə) — Key-value kolleksiyası, hash table əsasında
- **MEAP** (Manning Early Access Program) — Kitabın yarımçıq versiyası, çapdan əvvəl əlçatan
- **Method** (Metod) — Tipə bağlı funksiya, `func (p Person) Name()`
- **Module** (Modul) — Asılılıqları idarəetmək vahidi, `go.mod` faylı
- **Mutability** (Dəyişkənlik) — Dəyərin dəyişə bilməsi (mutable / immutable)

## N

- **Nil** (Sıfır Göstərici) — Pointer, map, slice, interface üçün sıfır dəyəri
- **Numeric Literal** (Ədəd Sabit Dəyəri) — `1000`, `0x3E8`, `0b111`, `1_000`

## O

- **OS Thread** (Əməliyyat Sistemi Thread-i) — OS tərəfindən idarə olunan icra vahidi

## P

- **Package** (Paket) — Funksiya və tiplərin təşkilat vahidi
- **Panic** (Çöküş) — Go runtime tərəfindən baş verən ağır xəta
- **Pointer** (Göstərici) — Yaddaş adresini saxlayan tip, `*Type`
- **Primitive Type** (Əsas Tip) — Daxili tiplər (`int`, `string`, `bool` və s.)
- **Property** (Xüsusiyyət) — Struct sahəsi
- **Polymorphism** (Çoxformlılıq) — Eyni interfeys üzərində fərqli tiplərin işləməsi

## R

- **Range Loop** (Aralıq Döngüsü) — `for i, v := range collection`
- **Raw String** (İşarətli Mətn) — `` `Hello` `` kimi, heç bir interpretasiya yox
- **Receiver** (Qəbul Edən) — Metodun hansı tipə aid olduğu, `(p Person)`
- **Reference Type** (Referans Tipi) — Slice, map, pointer — kopya deyil, reference ötürülür
- **Reflection** (İnfissas) — Runtime-də tipləri yoxlamaq imkanı (`reflect` paketi)
- **Rune** (Simvol) — `int32` alias, Unicode code point

## S

- **Scope** (Görüş Sahəsi) — Dəyişənin görünür olduğu kod sahəsi
- **Shadowing** (Kölgələmə) — Eyni adlı dəyişənin daxili scope-da yenidən yaradılması
- **Short Declaration** (Qısa Elan) — `:=` operatoru ilə elan və təyinat
- **Slice** (Kəsik) — Dinamik massiv, `[]int` kimi
- **Slice Expression** (Kəsmi İfadə) — `array[low:high]` kimi alt kolleksiya yaratmaq
- **SplitFunc** (Ayırma Funksiyası) — `bufio.Scanner` üçün tokenization növü
- **Standard Library** (Standart Kitabxana) — Go ilə birlikdə gələn paketlər (`fmt`, `os`, `strings`)
- **Static Typing** (Statik Tip) — Tip compile vaxtında müəyyən olunur
- **String** (Mətn) — Read-only byte array, UTF-8
- **Struct** (Strukt) — Composite tip, sahələr qrupu
- **Sub-kod** — Kod blokunun hissə izahı

## T

- **Tag** (Etiket) — Struct sahəsinin metadata-sı, `` json:"name" ``
- **Terminal** (Terminal) — Command-line interface (komanda sətri interfeysi)
- **Toolchain** (Alət Zənciri) — Go kompilyatoru və yardımçı alətlər
- **Type** (Tip) — Dəyərin növü (`int`, `string`, `bool`, `Person` və s.)
- **Type Conversion** (Tip Çevrilməsi) — `float64(i)` kimi açıq çevrilmə
- **Type Inference** (Tip Çıxarma) — `text := "hello"` kimi dəyərdən tipin avtomatik müəyyən edilməsi

## U

- **Unsafe** (Təhlükəsiz Olmayan) — `unsafe` paketi, low-level yaddaş əməliyyatları
- **Unicode** (Yunikod) — Dünya yazıları standartı, UTF-8 encoding
- **Unsigned Integer** (İşarətsiz Tam Ədəd) — `uint`, `uint8`, `uint16`, `uint32`, `uint64`

## V

- **Value Type** (Dəyər Tipi) — Array, struct — kopyalanır
- **Variable** (Dəyişən) — `var` və ya `:=` ilə yaradılan dəyər saxlama vahidi
- **Vendor** (Təchizatçı) — Third-party (üçüncü tərəf) asılılıqlar
- **Version Control** (Versiya Nəzarəti) — Git, code history (kod tarixçəsi)

## W

- **wc** — Word count (söz sayma) CLI utility (yardımçı proqram), Unix/Linux standart
- **Wildcard** (Ümumi maska) — `*` simvolu, bütününə uyğun gələn

## Z

- **Zero Value** (Sıfır Dəyəri) — Tipin default dəyəri (`0`, `false`, `""`, `nil`)
