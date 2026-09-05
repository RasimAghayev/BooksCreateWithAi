# 4. Exploring Popular Composition Patterns

**Səhifələr:** 81-93

## Bu fəsil nədən bəhs edir?

Komponentləri necə bir-biri ilə əlaqələndirmək olar? Bu fəsil React-in ən populyar composition patternlərini — Container/Presentational, HOC və FunctionAsChild-i izah edir.

## Əsas fikirlər

### 1. Props ilə ünsiyyət

Komponentlər arasında məlumat ötürməyın əsas yolu props-dur. Props yalnız valideyn-dən uşağa axır (one-way data flow). State lift etmək üçün state-i yuxarıya qaldırmaq lazımdır.

### 2. Container / Presentational pattern

Bu pattern komponenti iki hissəyə bölür: Container (məntiq, state, data fetching) və Presentational (yalnız UI). Adətən sonuna `Container` əlavə olunur. Üstünlüyü: presentational komponentlər təkrar istifadə edilə bilər, container isə data mənbəyini təcrid edir.

### 3. HOC (Higher-Order Component)

HOC funksiyadır ki, komponenti arqument kimi qəbul edib yeni komponent qaytarır. Məqsəd: cross-cutting concern-ləri (auth, logging, theme) komponentlər arasında paylaşmaq. Məsələn: `withAuth(Component)`. Hooks dövründə HOC daha az istifadə olunur, custom Hooks onu əvəz edir.

### 4. FunctionAsChild pattern

Bu pattern komponentə render funksiyasını `children` prop kimi göndərməyə imkan verir: `<DataProvider>{data => <Display data={data} />}</DataProvider>`. Bu, HOC-dan daha elastikdir və 'render prop' da adlanır.

## Əsas terminlər

- Props
- Composition
- Container Pattern
- Presentational Pattern
- HOC
- Higher-Order Component
- FunctionAsChild
- Render Prop
- Children

## Praktik nəticə

Bir siyahı komponentini Container + Presentational şəklində refaktor et. Sonra `withLoading` adlı HOC yaz ki, istənilən komponentə loading state əlavə edə bilsin.

## Mənbə

Pages: 81-93