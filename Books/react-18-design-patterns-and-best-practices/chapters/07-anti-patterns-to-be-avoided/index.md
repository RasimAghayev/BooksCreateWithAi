# 7. Anti-Patterns to Be Avoided

**Səhifələr:** 142-152

## Bu fəsil nədən bəhs edir?

React-də tez-tez rast gəlinən və qarşısını almaq lazım olan 3 anti-pattern: state-i prop ilə initializasiya, index-i key kimi istifadə, qeyri-standart prop-ları DOM-a spread etmək.

## Əsas fikirlər

### 1. State-i properties ilə initializasiya etmək

`useState(props.initial)` çağıranda state yalnız ilk render-də prop-dan oxunur. Sonradan prop dəyişsə belə state köhnə dəyərdə qalır. Bu, uşaq komponentin state-inin valideynin prop-u ilə sinxronlaşmamasına səbəb olur. Həll: ya state-i qaldır (lift up), ya da açıq şəkildə adlandır ki, niyyət aydın olsun.

### 2. Index-i key kimi istifadə etmək

Siyahılarda `key={index}` istifadəsi React-in reconciliation mexanizmini pozur. Məsələn, siyahının əvvəlinə yeni element əlavə etsən, bütün elementlərin key-ləri dəyişir, React bütün siyahını yenidən render edir və state itir. Həmişə unikal, sabit ID istifadə edin.

### 3. Qeyri-standart prop-ları DOM-a spread etmək

`{...props}` yazanda DOM-a React tanımayan atributlar (`prop1`, `callback1` və s.) ötürülür. Bu, brauzerdə DOM warning yaradır və performansı aşağı salır. Həll: destrukturla `const { known, ...rest } = props` və yalnız `...rest` ötür.

## Əsas terminlər

- Anti-pattern
- State initialization
- Key Prop
- Index as Key
- Spread operator
- Reconciliation
- DOM warning

## Praktik nəticə

Mövcud layihəndə `key={index}` istifadə edən bütün yerləri tap və `key={item.id}` ilə əvəz et. React DevTools ilə yoxla ki, siyahıda element əlavə edəndə bütün elementlər yenidən render olunmur.

## Mənbə

Pages: 142-152