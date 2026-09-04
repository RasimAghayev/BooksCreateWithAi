# Chapter 2 — Разработка микросервиса авторизации и аутентификации (Auth)

Pages: 87-110

## Bu chapter nədən bəhs edir?
Bu chapter Authorization (avtorizasiya) və Authentication (autentifikasiya)
mikroservisini hazırlamağı izah edir. JWT, şəxsiyyət vəsiqələri, bir dəfə
parol və əsas təhlükəsizlik tədbilləri haqqında bəhs edir.

## Əsas fikirlər
### 1. Autentifikasiya növləri
- Parol əsaslı autentifikasiya
- Sertifikat əsaslı autentifikasiya
- Bir dəfə parol (OTP) autentifikasiya
- Açar diqqəti autentifikasiya
- Token əsaslı autentifikasiya (JWT)

### 2. Təhlükəsizlik tədbilləri
- Buffer overflow (buffer daşması)
- Race condition (yarış vəziyyəti)
- Input validation (giriş validasiyası)
- Authentication attacks (autentifikasiya hücumları)
- Authorization attacks (avtorizasiya hücumları)
- Client-side attacks (client tərəfi hücumları)

### 3. Auth modulu
- JWT tokenları
- Refresh token mexanizmi
- Rol əsaslı giriş nəzarəti (RBAC)

## Əsas terminlər
- Authentication (autentifikasiya / kimlik təsdiqi)
- Authorization (avtorizasiya / icazə yoxlanışı)
- JWT (JSON Web Token)
- OTP (One-Time Password / bir dəfəlik parol)
- RBAC (Role-Based Access Control / rol əsaslı giriş nəzarəti)
- Password Hashing (parol hashləmə)
- Salt (duz / təsadüfi dəyər)
- Certificate (sertifikat)
- Vulnerability (zəiflik)

## Praktik nəticə
Auth mikroservisi hazırlanır və User mikroservisi ilə inteqrasiya olunur.
Token əsaslı autentifikasiya mexanizmi işləyir.

## Mənbə
Pages: 87-110
