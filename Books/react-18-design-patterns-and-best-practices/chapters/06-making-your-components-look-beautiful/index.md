# 6. Making Your Components Look Beautiful

**Səhifələr:** 115-141

## Bu fəsil nədən bəhs edir?

React komponentlərinə stil vermək üçün müxtəlif yanaşmalar: inline styles, CSS modules, atomic CSS modules və styled-components.

## Əsas fikirlər

### 1. CSS-in JavaScript-də problemi

Böyük layihələrdə CSS selektorları çox qarışır. Meta-nın təcrübəsi göstərir ki, CSS sinif adları təkrarlanır, specificity müharibələri yaranır. Həll yolu: komponentə aid CSS-i komponentlə birlikdə saxlamaq (CSS-in-JS).

### 2. Inline styles

React `style` prop'unu qəbul edir: `<div style={{ color: 'red', fontSize: 16 }}>`. Üstünlüyü: lokallaşdırılmış, dinamik asanlıqla dəyişdirilə bilər. Çatışmazlığı: media queries, pseudo-classes (`:hover`) dəstəklənmir, performans baxımından bəzi edge-case-lərdə yavaş ola bilər.

### 3. CSS Modules

CSS Modules komponentə aid CSS faylını lokal scope ilə yaradır: `Button.module.css` daxilindəki `.primary` sinfi avtomatik olaraq unikallaşdırılır (məsələn, `.primary__abc123`). Bu, klassik BEM kimi adlandırma ehtiyacını aradan qaldırır. Webpack 5 bunu default dəstəkləyir.

### 4. Atomic CSS modules

Hər stil üçün atomik (çox kiçik, tək-tək istifadə olunan) siniflər yaradılır: `.bg-blue`, `.p-2`. Bu yanaşma Tailwind CSS-in əsasını təşkil edir. Üstünlüyü: daha az CSS faylı, yaxşı cache.

### 5. styled-components

`styled-components` runtime CSS-in-JS kitabxanasıdır. Template literal ilə stil yazılır: `const Button = styled.button`padding: 0.5rem;```. Dinamik prop-larla stil dəyişdirilə bilər: `<Button primary />`. SSR üçün Next.js ilə inteqrasiya mövcuddur.

## Əsas terminlər

- Inline Style
- CSS Modules
- Atomic CSS
- styled-components
- CSS-in-JS
- Webpack 5
- BEM
- Tailwind
- Dynamic Styling

## Praktik nəticə

Eyni Button komponentini 3 variantda yaz: inline style, CSS Module və styled-components. Hansının daha oxunaqlı və maintainable olduğunu qiymətləndir.

## Mənbə

Pages: 115-141