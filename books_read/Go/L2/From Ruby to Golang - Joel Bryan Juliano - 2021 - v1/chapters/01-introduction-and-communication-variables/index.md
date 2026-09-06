# Chapter 1 — Introduction: kommunikasiya və dəyişənlər (book səh. 1-2)

## Bu chapter nədən bəhs edir?

Kitabın fəlsəfəsi (Rubyist perspektivindən Go), müəllifin keçid hekayəsi və proqramlaşdırmada "kommunikasiya" anlayışı: lokal dəyişənlərin funksiya-daxili scope-u vs Ruby-da qlobal çıxış.

---

## Əsas fikirlər

### 1. Kitabın yanaşması — analogiya öyrənmə

Müəllif (Joel Bryan Juliano, 12+ il təcrübə, Amsterdam) 2018-də Go istifadə edən şirkətə keçəndə öz şəxsi öyrənmə qeydlərini onlayn məqalələr seriyasına çevirib. **Metod:** hər Go anlayışını Ruby-dəki analoqu ilə izah etmək — "biznesdə bildiklərinlə yeni dili qurmaq".

**Go-nun seçilmə səbəbləri (müəllifin siyahısı):**
- Asan oxunan, aydın sintaksis; tək binari fayla kompayl — sürətli və kompakt, cross-platform
- Statik tip + garbage collection — effektiv
- "Modern C": paket dəstəyi, yaddaş təhlükəsizliyi, konkurrentlik daxildir
- Cloud-native ekosistem: Docker, Kubernetes, Ethereum, Terraform Go ilə; Google, Netflix, Dropbox, Uber production-da

**OOP mövqeyi:** Robert Sessions (1992): "OOP, strukturlaşdırılmış proqramlaşdırmanın sadəcə sağlam-düşüncəli uzantısıdır". Go OOP dili DEYİL (seçimdir), amma Ruby-dəki OOP texnikaları Go-ya tətbiq oluna bilər.

### 2. Kommunikasiya → dəyişən metaforası

Bütün canlılar — nəhəng heyvanlardan bakteriya/qrirose-lərə — mesaj mübadiləsi aparır. Proqramda da kodun hissələri kommunikasiya edir: **dəyişən = kodun konkret hissəsinə bağlı dəyərlərin ötürülməsi/kordinasiyası.**

**Lokal dəyişən (Ruby):**

```ruby
def say_hello(message)
  hello = "Hello"       # LOKAL — funksiya daxilindən başqa yerdən çatmır
  puts hello + message
end
say_hello("World")       # Hello World
```

Xaricdən `hello = "Hi"` əlavə etsək belə — funksiya DAXİLİNDƏKİ hello təsirlenmir (əks halda başqa dəyişən kölgəsidir) → nəticə "Hello World" qalır.

**Qlobal dəyişən** — kodun hər yerindən çatılaşan.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Analogiya öyrənmə | Ruby anlayışını Go-yə xəritələmə — kitabın metodu |
| Lokal dəyişən | Yalnız funksiya daxilindən çatılaşan dəyər |
| Qlobal dəyişən | Kod boyu çatılaşan dəyər |
| "Modern C" | Paketlar + memory safety + GC + konkurrentlik daxili |
| Cloud-native | Docker/K8s/Terraform ekosistemi — Go-nun evi |

---

## Praktik nəticə

1. **Ruby biliklərini xəritələ:** Instans dəyişəni→struct, module→package, each→for range, splat→variadic, mixin→embed — kitabın izlədiyi xəritə.
2. **Keçid dövründə OOP-təfəkkürü itirmə:** Go OOP deyil, amma strukturlar+metodlar+interfeyslər eyni məsələləri həll edir.
3. **Kommunikasiya prizması:** dəyişənlər, sonra instans dəyişənləri (struct), map/array, paketlər, interfeyslər — hər fəsil bir kommunikasiya qatıdır.

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021 (ISBN 978-1080944002)
- Chapter: Introduction, book səh. 1-2
- PDF səhifələri: 7-9
