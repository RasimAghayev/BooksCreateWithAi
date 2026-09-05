# 14. MonoRepo Architecture

**Səhifələr:** 354-412

## Bu fəsil nədən bəhs edir?

MonoRepo (monorepository) — birdən çox paket və tətbiqi bir git reposunda saxlamaq strategiyasıdır. Bu fəsil NPM Workspaces, TypeScript, Webpack ilə production-grade MonoRepo qurur.

## Əsas fikirlər

### 1. MonoRepo-nun üstünlükləri

Birdən çox paket bir yerdə: kod paylaşmaq asanlaşır, refactoring tək repoda, dependency-lər sinxron, paketlər arası koordinasyon yaxşılaşır. Nx, Turborepo, Lerna alternativləri var.

### 2. NPM Workspaces ilə MonoRepo

`package.json` daxilində `workspaces: ['packages/*', 'apps/*']` təyin et. Bu, npm install zamanı bütün alt paketləri bir yerdə qurur. `-w <workspace>` flag ilə konkret workspace-də əmr işlət.

### 3. TypeScript MonoRepo konfiqurasiyası

`tsconfig.base.json` ilə əsas qaydaları (`compilerOptions` — target, module, strict) yaz, hər paketdə `tsconfig.json` extends edir. Project references ilə build performansı yaxşılaşır.

### 4. devtools paketi — Webpack konfiqurasiya

Ümumi Webpack konfiqurasiyasını `@my-org/devtools` paketinə köçür. Hər tətbiq bu paketdən `common`, `dev`, `prod` konfiqurasiyaları import edir. Colorful logger də bu paketdə olur.

### 5. utils paketi

Ümumi köməkçi funksiyalar (`formatDate`, `validateEmail` və s.) burada saxlanılır. Hər paket bu paketdən idxal edə bilər.

### 6. API paketi

GraphQL schema, resolvers, Sequelize modelləri, xidmət konfiqurasiyası burada qurulur. CRM xidməti kimi bir backend burada yaşayır. Shared model (User) və shared GraphQL types/resolvers yaradılır.

### 7. Frontend paketi

Next.js istifadə edilən frontend tətbiqi. Sites sistemi (məsələn, san-pancho saytı), Page Switcher komponenti, Login sistemi, sites konfiqurasiyası burada yerləşir.

### 8. Putting everything together

Bütün paketlər `npm run build` ilə build olunur, lokal-da `npm run dev` ilə development işlədilir. Demo zamanı həm API, həm frontend paralel işləyir.

## Əsas terminlər

- MonoRepo
- Monorepository
- NPM Workspaces
- package.json
- TypeScript Project References
- Webpack
- tsconfig
- Shared Package
- Frontend
- API
- Sites
- Next.js
- GraphQL

## Praktik nəticə

Sıfırdan bir MonoRepo qur: `apps/web`, `packages/utils`, `packages/api`, `packages/devtools`. Workspace-ləri qur, TypeScript konfiqurasiya et, Webpack ilə build et.

## Mənbə

Pages: 354-412