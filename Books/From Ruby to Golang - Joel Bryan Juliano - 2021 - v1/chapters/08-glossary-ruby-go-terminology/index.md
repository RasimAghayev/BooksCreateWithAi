# Chapter 8 — Glossary: Ruby↔Go Terminologiya Lüğəti (book səh. 112-158)

## Bu chapter nədən bəhs edir?

Kitabın geniş lüğəti (təxminən 46 PDF səhifə): simvollar (`%+v`, `%s`, `&&`, `$GOPATH`), operatorlar, API/CLI/abstraksiya/attribute kimi proqramlaşdırma anlayışları — Wiktionary/Wikimedia mənbəli təriflərlə. Ruby-dən gələn oxucu üçün istinad materialıdır.

---

## Əsas məzmun strukturu

Lüğət 3 kateqoriyaya bölünür:

**1. Simvol və konstruktlar (nümunələr):**

| Element | Tərif (kitabdan) |
|---------|------------------|
| `\n` | Newline işarə edən control character ardıcıllığı |
| `!` | Məntiqi NOT operatoru — boolean əksinə çevirir |
| `""` | Boş string |
| `#<Set: {...}>` | Ruby Set obyektinin çap forması |
| `$` | Shell-də xüsusi simvol — dəyişən prefiksi |
| `$GOPATH` | Köhnə modul sisteminin kitabxana axtarış yolu |
| `%+v` | fmt konstruktu — struct dəyər + SAHƏ ADLARI ilə çap |
| `%s` | fmt — string formatı |
| `%` | Go-da xüsusi simvol — qayıdış dəyəri gözləyən tip prefiksi |
| `&&` | Məntiqi AND |
| `&MapName{}` | İlkinləşdirilməmiş boş go map — dəyişən kimi assign oluna bilər |
| `//` | Comment göstəricisi (C/C++/Java da daxil) |
| `=` | Bərabərlik/əmr simvolu |

**2. Proqramlaşdırma anlayışları (alpha-sıralı, Wiktionary əsaslı):**

Ability (uyğunluq) · Accessibility · API · Array · Assembler · Assignment · Attribute · Binary file · Boolean · Branch · Bug · Byte · Cache · Channel · Character · Class · Client · Cloud computing · Code · Collection · Command-line interface · Compiler · Computer · Concurrent computing · Control flow · CPU · Database · Data structure · Data type · Declaration · Dependency · Dereference operator · Device driver · Documentation · Dynamic programming language · Encoding · Enumerator · Error · Executable · Expression · Field · File system · Firmware · Floating point · Function · Garbage collection · Go · Gopher · Hardware · Hash table · IDE · Immutable object · Index · Inheritance · Instruction · Interface · Interpreter · Iteration · JSON · Library · Linker · Linter · Machine code · Map · Memory · Method · Microservice · Module · Namespace · Network · Object · Object-oriented programming · Operating system · Operator · Package management · Parameter · Pointer · Process · Programming language · Queue · RAM · Recursion · Reference · Register · Runtime · Scope · Scripting language · Server · Slice · Software · Sorting · Source code · Stack · Statement · Static typing · String · Struct · Syntax · Thread · Type system · Unix · Variable · Virtual machine · Web server və s.

Hər tərif 1-3 cümləlik sadə izahdır (Wiktionary/Wikipedia-dan adaptasiya — Credits bölməsində qeyd olunur).

**3. Kitabə xas anlayışlar:** produce/basket nümunələrindən gələn terminlər və Ruby-Go xəritələnmə istiqamətləri.

---

## Lüğətin istifadə qaydası

1. **Keçid dövründə istinad:** Ruby termini (mixin, block, enumerator) Go kontekstində nə kimi səslənir — əsas mətnlərin sözlüyü.
2. **Sadəlik prinsipi:** Təriflər başlanğıc səviyyəyə uyğun — heç bir ön-şərt tələb etmir.
3. **İki dilli praktika:** `Hash → map`, `module → package`, `each → range`, `splat → variadic` kimi əsas xəritələmələr termin səviyyəsində də görünür.

**Qiymətləndirmə:** 46 səhifəlik lüğət kitabın təxminən 30%-i — müəllifin "öyrənərkən qarşılaşdığım HƏR söz" yanaşmasının məhsulu. Sistematik Ruby→Go fərqlilik cədvəli YOXDUR; klassik ensiklopedik sadalama formatındadır.

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 8: Glossary, book səh. 112-158
- PDF səhifələri: 118-165
- Credits: Wiktionary/Wikimedia Foundation (lüğət təriflərinin mənbəyi)
