# 2. Introducing TypeScript

**Səhifələr:** 44-57

## Bu fəsil nədən bəhs edir?

TypeScript Microsoft tərəfindən yaradılmış JavaScript-in statik tipli supersetidir. Bu fəsil TypeScript-in əsaslarını və React layihələrində necə istifadə olunacağını öyrədir.

## Əsas fikirlər

### 1. TypeScript-in üstünlükləri

Statik tipləmə compile-time-da xətaları tutur, daha yaxşı IDE dəstəyi verir, kodun niyyətini aydınlaşdırır, refactoring-i asanlaşdırır. Anders Hejlsberg (C# yaradıcısı) tərəfindən dizayn edilib.

### 2. Types (Tiplər)

TypeScript dəyişənlərə tip verməyə imkan verir: `string`, `number`, `boolean`, `array`, `tuple`, `enum`, `any`, `unknown`, `void`, `never`. Məsələn: `const name: string = 'Carlos'`. Hər tipin istifadəsi fərqlidir — `any` hər şeyi qəbul edir, `unknown` isə əvvəlcce yoxlama tələb edir.

### 3. Interfaces (İnterfeyslər)

İnterfeyslər obyektin formasını təyin edir: `interface User { id: number; name: string; email?: string }`. `?` optional field bildirir. İnterfeyslər `type` açar sözü ilə də yaradıla bilər, lakin daha yaxşı adlandırma dəstəyi verir.

### 4. Extending və Merging

İnterfeysləri `extends` ilə genişləndirmək olur: `interface Admin extends User { role: 'admin' }`. Eyni adda iki interface avtomatik merge olunur — bu xüsusiyyət MonoRepo-larda çox faydalıdır.

### 5. Enums

Enum sabit dəyərlər toplusudur: `enum Status { Active, Inactive, Pending }`. Hər enum üzvünün arxasında rəqəm dayanır, lakin string dəyərlər də təyin etmək olar.

### 6. tsconfig.json

`tsconfig.json` TypeScript konfiqurasiya faylıdır. MonoRepo-larda `common` və `specific` hissələrə bölmək tövsiyə olunur ki, bütün paketlər eyni əsas qaydaları paylaşsın.

## Əsas terminlər

- TypeScript
- Type
- Interface
- Enum
- Namespace
- Template Literal
- tsconfig.json
- Generics

## Praktik nəticə

Mövcud JavaScript faylını `.ts` yeniden adlandır və tip xətalarını düzəlt. Sonra `interface User { id, name, email }` yarat və komponentə prop olaraq ötür.

## Mənbə

Pages: 44-57