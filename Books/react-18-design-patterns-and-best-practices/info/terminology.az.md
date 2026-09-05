# Terminology — React 18 Design Patterns and Best Practices

Bu fayl kitabda rast gəlinən əsas texniki terminləri və onların Azərbaycanca qarşılıqlarını ehtiva edir.

## A

- Apollo Client — React üçün GraphQL client kitabxanası
- Apollo Server — GraphQL server implementasiyası
- Anti-pattern (Anti-pattern) — Tez-tez rast gəlinən, lakin yanlış və zərərli proqramlaşdırma təcrübəsi
- Atomic CSS modules — Çox kiçik, tək-tək istifadə olunan CSS sinifləri yanaşması

## B

- Babel — JavaScript/TypeScript transpiler
- Batching (Birləşdirmə) — Çoxlu state yeniləmələrinin bir render dövründə birləşdirilməsi
- Bundle (Buncl) — Webpack/Vite kimi alətlərin yaratdığı yekun JavaScript faylı

## C

- Children prop — Komponentin açılış-bağlanış teqləri arasındakı kontent
- Class component (Sinif komponenti) — ES6 sinfi şəklində yazılan React komponenti
- Component (Komponent) — Təkrar istifadə olunan React UI bloku
- Composition (Kompozisiya) — Kiçik komponentlərin birləşdirilərək böyük komponent yaratması
- Concurrent mode (Konkurent rejim) — React 18-in render prosesini kəsə bilən yeni rejimi
- Container component (Konteyner komponenti) — State və məntiqi saxlayan komponent
- Controlled component (Nəzarət edilən komponent) — Dəyəri React state ilə idarə olunan input
- CSS modules (CSS modulları) — Local scope'lu CSS sinifləri
- Currying (Karri) — Çox arqumentli funksiyanı tək arqumentli funksiyalara çevirmə

## D

- Declarative programming (Deklarativ proqramlaşdırma) — Nə etmək istədiyimizi bildirən yanaşma
- Dependency array (Asılılıq massivi) — useEffect/useMemo üçün asılılıqlar siyahısı
- DOM (Document Object Model) — Brauzerin HTML strukturunun obyekt modeli
- Droplet — DigitalOcean-un virtual server xidməti

## E

- Effect (Yan təsir) — useEffect ilə yerinə yetirilən yan təsir əməliyyatı
- Enum (Enumerator) — TypeScript-də sabit dəyərlər toplusu
- ESLint — JavaScript/TypeScript linter aləti

## F

- Falsy — Boolean kontekstində false sayılan dəyər
- Function as Child (Funksiya kimi uşaq) — Komponentə funksiya göndərmə patterni
- Functional programming (Funksional proqramlaşdırma) — Pure funksiyalara əsaslanan paradigm
- forwardRef — Ref-i uşaq komponentə yönləndirən HOC

## G

- GraphQL — API sorğu dili və runtime
- Grid (CSS Grid) — 2D layout CSS sistemi

## H

- HOC (Higher-Order Component) — Yüksək dərəcəli komponent, başqa komponenti qəbul edib yenisini qaytaran funksiya
- Hook (Hook) — React-in funksional komponentlərə state və lifecycle əlavə edən API
- hydrateRoot — React 18-in server render olunmuş kontenti canlandırmaq üçün yeni API
- Hydration (Canlandırma) — Server HTML-ini client React ağacına çevirmə

## I

- Imperative programming (İmperativ proqramlaşdırma) — Addım-addım əmrlər verən yanaşma
- Immutability (Dəyişməzlik) — Dəyərin yaradıldıqdan sonra dəyişdirilməməsi prinsipi
- Inline style (Daxili stil) — `style` prop'u vasitəsilə verilən CSS
- Interface (İnterfeys) — TypeScript-də obyekt formasının təyini
- Iteration (İterasiya) — Massiv və ya kolleksiyanın elementlər üzrə dövr edilməsi

## J

- JSX (JavaScript XML) — JavaScript daxilində HTML-ə oxşar sintaksis
- JWT (JSON Web Token) — İmzalı token standartı, auth üçün istifadə olunur

## K

- Key (Açar) — React-in siyahı elementlərini identifikasiya etməsi üçün xüsusi prop

## L

- Loader — React Router v6.4-də marşrut yüklənməmişdən əvvəl işləyən funksiya
- Linter — Kod keyfiyyətini yoxlayan alət

## M

- Memo (memo) — Komponenti prop dəyişmədikcə render etməməsi üçün yaddaş
- Memoization (Yaddaşlama) — Hesablama nəticəsinin saxlanılması
- Migration (Miqrasiya) — Köhnə API-dan yenisinə keçid
- MonoRepo (Mono-repozitori) — Birdən çox paketi bir repozitoriyada saxlama
- Mutation — GraphQL-də məlumat dəyişdirmə əməliyyatı

## N

- Namespace (Ad fəzası) — TypeScript-də identifikatorları qruplaşdırma mexanizmi
- Next.js — React üçün production framework (SSR, SSG, routing)
- Node.js — JavaScript runtime mühiti
- npm Workspaces — MonoRepo idarəsi üçün npm xüsusiyyəti

## P

- Purity (Saflıq) — Funksiyanın eyni giriş üçün həmişə eyni çıxış verməsi, yan təsirsiz olması
- PM2 — Node.js process manager
- Prettier — Kod formatlama aləti
- Presentational component (Təqdim komponenti) — Yalnız UI render edən komponent
- Props (Properties) — Komponentə xaricdən göndərilən parametrlər

## Q

- Query (Sorğu) — GraphQL-də məlumat oxuma əməliyyatı

## R

- Reconciliation (Uzlaşdırma) — React-in köhnə və yeni virtual DOM-u müqayisə etməsi
- Redux Toolkit — Redux-un müasir, sadələşdirilmiş versiyası
- Ref (Reference) — DOM elementinə və ya komponent instansiyasına birbaşa istinad
- Resolver — GraphQL sorğusunu həll edən funksiya
- Routing (Marşrutlama) — URL-ləri komponentlərə uyğunlaşdırma

## S

- Sequelize — Node.js üçün ORM kitabxanası
- Server-Side Rendering (SSR) — HTML-in server tərəfində generasiya olunması
- Slice (Redux Toolkit hissəsi) — Redux state-inin modul hissəsi
- Spread operator (Yayılma operatoru) — `...` ilə massiv/obyekti yayma
- Suspense — React-in async məlumat yükləməsini idarə edən komponenti
- SWR — Veri-fetching kitabxanası (stale-while-revalidate)

## T

- Template literal (Şablon literalı) — Backtick ilə yazılan string, interpolation dəstəyi ilə
- Transition (Keçid) — React 18-in uzunmüddətli state yeniləmələri üçün mexanizmi
- TypeScript — Microsoft tərəfindən yaradılmış JavaScript-in statik tipli superseti

## U

- Uncontrolled component (Nəzarətsiz komponent) — Form dəyəri DOM-da qalan input
- Universal application (Universal tətbiq) — Həm server, həm client tərəfdə işləyən tətbiq
- useCallback — Funksiya referensini memoize edən hook
- useDeferredValue — Dəyəri gecikmə ilə yeniləyən hook
- useEffect — Yan təsir həyata keçirən hook
- useId — Unikal ID generasiya edən hook
- useInsertionEffect — CSS injection üçün hook
- useMemo — Dəyəri memoize edən hook
- useReducer — State idarəsi üçün reducer-pattern hook
- useState — Lokal state yaratmaq üçün hook
- useTransition — Transition statusunu izləyən hook

## V

- Vite — Sürətli frontend build aləti (Create-React-App əvəzinə)
- Vitest — Vite ilə inteqrasiya edilmiş test runner
- Virtual DOM — Real DOM-un yaddaşdakı təmsili

## W

- Webpack — Modul bundler aləti
- Worktree — Git-in paralel işləmə mexanizmi

## Y

- Yield — Generator funksiyalarda dəyər qaytarma