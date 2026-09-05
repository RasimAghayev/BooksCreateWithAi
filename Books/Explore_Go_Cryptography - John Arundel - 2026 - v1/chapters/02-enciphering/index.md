# 2. Enciphering

**Səhifələr:** 45-60

## Bu fəsil nədən bəhs edir?

Shift cipher-ni real məlumatlarla işləyən bir proqram çevirmək. Standart giriş/çıxış (stdin/stdout) istifadə edərək filtr (filter) kimi işləyən enciphering utility quraşdırılır.

## Əsas fikirlər

### 1. Filter Interface (Filtr İnterfeysi)
Proqramın stdin-dən oxuyub stdout-a yazması. Bu, Unix pipeline-larında istifadə imkanı verir.

### 2. Test-Driven Development (TDD)
Test yazmaqdan əvvəl funksiyanın davranışını müəyyənləşdirmək:
- Hərfi birbaşa əvəz etmək
- Key (açar) dəyərinə görə hərfi dəyişdirmək

### 3. Test strukturı
- Test case-lərin cədvəlləşdirilməsi
- Subtests (alt testlər) istifadəsi
- Edge case-lərin yoxlanılması

## Əsas terminlər
- Enciphering (Şifrələmə)
- Filter (Filtr)
- stdin (Standart giriş)
- stdout (Standart çıxış)
- Test Case (Test zamanı)
- Subtests (Alt testlər)

## Praktik nəticə

`encipher` funksiyası və onun üçün testlər yazılır. Proqram komanda sətiri vasitəsilə real mətnləri şifrələyə bilir.

## Mənbə

Pages: 45-60
