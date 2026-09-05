# 11. Enumeration

**Səhifələr:** 178-198

## Bu fəsil nədən bəhs edir?

Açar məkanını (keyspace) sıralama (enumeration) və parallel kırma. 256-bit açarlar üçün brute force hücumunun praktiki olaraq mümkün olmadığı göstərilir.

## Əsas fikirlər

### 1. Keyspace Enumeration
Bütün mümkün açar dəyərlərini sistematik olaraq yoxlamaq. Kiçik keyspace-lərdə effektiv, böyük keyspace-lərdə impraktikdir.

### 2. Parallel Enumeration
Çoxsaylı CPU-lar / thread-lər vasitəsilə kırma prosesini sürətləndirmək. Go-da goroutine-lərdən istifadə.

### 3. Practical Limits (Praktik Limitlər)
256-bit açar üçün 2²⁵⁶ mümkün dəyər var. Hətta dünyanın ən güclü superkompüteri ilə belə kırmaq mümkün deyil.

## Əsas terminlər
- Enumeration (Sıralama)
- Keyspace (Açar məkanı)
- Parallel Processing (Paralell emal)
- Goroutine
- Brute Force (Zorla hücum)

## Praktik nəticə

Parallel enumeration testi yazılır. Performans ölçülmələri aparılır.

## Mənbə

Pages: 178-198
