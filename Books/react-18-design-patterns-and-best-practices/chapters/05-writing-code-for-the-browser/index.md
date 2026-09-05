# 5. Writing Code for the Browser

**Səhifələr:** 94-114

## Bu fəsil nədən bəhs edir?

Browser üçün React kodu yazarkən qarşılaşacağımız vacib mövzular: formlar, event-lər, refs, forwardRef, animasiyalar və SVG.

## Əsas fikirlər

### 1. Controlled vs Uncontrolled komponentlər

Controlled komponentin dəyəri React state-də saxlanılır (hər dəyişiklikdə state yenilənir). Uncontrolled komponentdə isə dəyər DOM-da qalır, biz yalnız lazım olanda `ref` vasitəsilə oxuyuruq. Controlled daha çox React-yə uyğundur, amma bəzi hallarda (məs. file input) uncontrolled məcburidir.

### 2. Event handling

React event-ləri camelCase ilə yazılır (`onClick`, `onChange`, `onSubmit`). SyntheticEvent React-in öz event sistemidir, brauzerlər arasında uyğunluq təmin edir. Event handler-lər `useCallback` ilə memoizasiya olunmalıdır ki, uşaq komponentlərə göndəriləndə referens dəyişməsin.

### 3. useRef Hook

`useRef` DOM elementlərinə və ya mutable dəyərlərə istinad saxlamaq üçündür. Render arasında qalır, amma state kimi render-ə səbəb olmur. `inputRef.current?.focus()` kimi DOM API-lərini çağırmaq olur.

### 4. forwardRef

`forwardRef` ref-i uşaq komponentdən dərin DOM elementə yönləndirməyə imkan verir. Müasir React-də bunun əvəzinə `useImperativeHandle` istifadə olunur. forwardRef hələ də işləyir, lakin yeni kodlarda `ref` prop kimi ötürmək daha üstündür.

### 5. SVG və Animasiyalar

React SVG elementlərini birbaşa JSX ilə yazmağa imkan verir. Animasiyalar üçün CSS transitions və ya `requestAnimationFrame` istifadə olunur. React 18-in concurrent mode-u animasiyaları daha smooth edir.

## Əsas terminlər

- Controlled Component
- Uncontrolled Component
- Event Handler
- SyntheticEvent
- useRef
- forwardRef
- useImperativeHandle
- SVG
- Animation

## Praktik nəticə

Login formunu həm controlled, həm uncontrolled variantlarda yaz. Sonra `useRef` ilə autoFocus effektini implementasiya et.

## Mənbə

Pages: 94-114