# Chapter 7 — Event-Carried State Transfer (səh. 192-215)

## Bu chapter nədən bəhs edir?

Sinxron gRPC çağırışlarının event-carried state transfer ilə əvəz olunması:
lokal keş repoları, sync fallback-lər, AsyncAPI sənədləşməsi və çox mənbəli
read model (Search modulu).

## Əsas fikirlər

### 1. Refactor strategiyası
Store Management məlumatı (Store, Product) digər modullara axır. Hər istehlakçı
modul:
1. Öz **lokal cache repository**-sini qurur (`stores.products_cache` cədvəli)
2. Integration hadisələrinə abunə olur, hadisədə gələn state-i keşə yazır
3. Sinxron gRPC çağırışı əvəzinə lokal keşdən oxuyur

### 2. Upsert + idempotent Add
```go
// Add() metodunda unikal constraint pozuntusu xətaları IGNORE olunur
// səbəb: at-least-once çatdırılma → eyni hadisə 2 dəfə gələ bilər
```
Cache yazımı dedupe-li və idempotentdir.

### 3. Sinxron fallback (qaçış yolu)
Keşdə məlumat yoxdursa (yeni abunəçi, replay edilməmiş keçmiş): əvvəlcədən
mövcud gRPC endpoint-lər **fallback** kimi saxlanılır. Məcburi deyil, amma
soyuq başlanğıc/miss halını xilas edir.

### 4. Modul-modul dəyişikliklər
- **Store Management → Baskets/Depot:** lokal keş + fallback
- **Customers:** inteqrasiya hadisələrinin publisher-i olur
- **Order Processing → Notifications:** callback-yönlü gRPC əvəzinə
  `OrderStatusChanged` hadisəsi; notification hadisəni dinləyir
- **Payments → Order Processing:** invoice status hadisə ilə "itələnir" —
  Order Processing payment-in bitməsini gözləyib order-i tamamlayır

### 5. AsyncAPI sənədləşməsi
Asinxron API-lər də sənədlənməlidir (bad docs ≈ no docs):
- **AsyncAPI** spesifikasiyası — OpenAPI-nin async versiyası
- Generator: HTML sənəyə + boilerplate kod (typescript də daxil)
- Sync API-lər də eyni sənəddə birləşdirilə bilər

### 6. Yeni Search modulu
Sıfırdan modul: `/search` qovluğu, module.go pattern, composition root-da
product/order/customer cache repoları. **Çox mənbəli read model:**
- `orders` oxu cədvəli + `product_ids`, `store_ids` sütunları (axtarış üçün)
- Customer/store/product datası gələndə yazılır; order gələndə birləşdirilir
- Sadə SELECT ilə axtarış — join-lər yox, çünki hamısı bir read cədvəlində

## Əsas terminlər

- Event-Carried State Transfer
- Local Cache Repository
- Sync Fallback
- Idempotent Upsert
- AsyncAPI
- Multi-Source Read Model
- Cold Start (soyuq başlanğıc) problemi

## Praktik nəticə

- Hər istehlakçı öz nüsxəsini saxlayır → producer yoxdursa da işləyir
- Cache yazımları idempotent olmalı (at-least-once + dedupe)
- Fallback-i saxla — keçmiş məlumat üçün gRPC xilasedici olur
- Asinxron API-ləri AsyncAPI ilə sənədləşdir; axtarış modulları üçün
  birləşmiş read model qur

## Mənbə

Pages: 192-215 (Chapter 7, Event-Driven Architecture in Golang)
