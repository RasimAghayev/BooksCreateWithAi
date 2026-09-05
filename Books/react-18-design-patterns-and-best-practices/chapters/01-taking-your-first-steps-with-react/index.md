# 1. Taking Your First Steps with React

**Səhifələr:** 30-43

## Bu fəsil nədən bəhs edir?

Bu fəsil React ilə ilk addımları atmaq üçün nəzərdə tutulub. Declarative və imperative proqramlaşdırma arasındakı fərq, React elements-in iş prinsipi və JavaScript fatigue mövzusu izah edilir.

## Əsas fikirlər

### 1. Declarative vs Imperative programming

Imperative yanaşmada addım-addım əmrlər verilir (məsələn: 'elementi tap, dəyərini dəyişdir'). Declarative yanaşmada isə nə istədiyimizi bildiririk, React bunu necə edəcəyini özü həll edir. React declarative-dir, biz yalnız istədiyimiz vəziyyəti təsvir edirik.

### 2. React elements-in iş prinsipi

React hər render zamanı React elements (sadə JavaScript obyektləri) yaradır və React DOM bu elementləri real DOM-a çevirir. Bu virtual DOM yanaşması performans üçün çox vacibdir.

### 3. Komponent əsaslı yanaşma

React tətbiqləri 'components' adlanan kiçik, təkrar istifadə olunan parçalardan qurulur. Bu, klassik 'separation of concerns' (məsuliyyətin ayrılması) anlayışını 'separation of technologies'-dən fərqləndirir. Hər komponent öz məntiqi və görünüşünü ehtiva edir.

### 4. JavaScript fatigue

JavaScript dünyasında alət və framework çoxluğu inkişaf etdiriciləri yora bilər. Kitab bu 'yorğunluqla' mübarizə üçün kiçik, idarəolunan addımlarla irəliləməyi tövsiyə edir.

## Əsas terminlər

- Declarative
- Imperative
- React Element
- Virtual DOM
- Component
- JSX
- JavaScript Fatigue
- Vite

## Praktik nəticə

Bu fəsildən sonra `npm create vite@latest my-app -- --template react-ts` ilə ilk Vite + React + TypeScript layihəni qur və brauzerdə açılan ilk səhifəni dəyişdir.

## Mənbə

Pages: 30-43