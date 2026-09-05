# 8. React Hooks

**Səhifələr:** 153-191

## Bu fəsil nədən bəhs edir?

React hooks 16.8 versiyasından functional komponentlərə state və lifecycle imkanları əlavə etdi. Bu fəsil hooks-un dərin analizini verir — useState, useEffect, memo, useMemo, useCallback və useReducer.

## Əsas fikirlər

### 1. Hooks-un üstünlükləri

Hooks class komponentlərə ehtiyacı aradan qaldırır — `this` problemi yoxdur, daha az kod, asan test. Heç bir breaking change yoxdur, mövcud class komponentlər dəstəklənir. Yeni hook-lar yazmaq asandır (məs. useWindowWidth).

### 2. Hooks qaydaları

(1) Hooks yalnız top-level-də çağırılmalıdır — şərt daxilində, dövrdə, nested funksiyada çağırılmamalıdır. (2) Hooks yalnız React funksiyalarından çağırılmalıdır (funksional komponentlər və ya digər custom hook-lar). Bu qaydalar React-in state saxlama mexanizmini qoruyur.

### 3. useState

`useState(initial)` `[value, setValue]` qaytarır. State yeniləmə funksiyasına həm dəyər, həm funksiya ötürə bilərik: `setCount(prev => prev + 1)` — bu async yeniləmələrdə köhnə state-ə görə düzgün işləyir.

### 4. useEffect

`useEffect(callback, deps)` hər render-dən sonra çağırılır (deps massivi boşdursa yalnız mount-da). Cleanup funksiyası return edilir. Async işlər üçün IIFE və ya `AbortController` istifadə olunur.

### 5. memo, useMemo, useCallback

Bunlar performans optimallaşdırması üçündür. `memo` komponenti — eyni props ilə yenidən render etmir. `useMemo` dəyəri hesablayır və nəticəni memoizasiya edir. `useCallback` funksiya referensini memoizasiya edir. Amma bunları hər yerdə istifadə etmək — öz-özünə yaddaş xərcləri yaradır, fayda verməz.

### 6. useReducer

`useReducer(reducer, initialState)` `[state, dispatch]` qaytarır. Reducer pure funksiyadır: `(state, action) => newState`. Complex state-lərdə useState-dən daha yaxşıdır. Redux-dan fərqli olaraq, yalnız lokal state üçün istifadə olunur.

## Əsas terminlər

- Hook
- useState
- useEffect
- useMemo
- useCallback
- memo
- useReducer
- Custom Hook
- Rules of Hooks
- Dependencies
- Cleanup
- Reducer

## Praktik nəticə

Bir class komponenti (məsələn, klassik Counter) Hooks-a miqrasiya et. Sonra `useWindowWidth` adlı custom hook yaz və onu `useEffect` ilə test et.

## Mənbə

Pages: 153-191