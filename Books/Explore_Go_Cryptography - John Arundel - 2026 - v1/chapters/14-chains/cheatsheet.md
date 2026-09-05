# Cheat Sheet — Chains

## Hash Chain

### Zəncir hash quruluşu
```
H(Block₁) = H₀
H(Block₂) = H(H₀ + Data₂)
H(Block₃) = H(H₁ + Data₃)
```

**İzah:**
Hər yeni blok əvvəlki hash dəyərini və öz məlumatını götürərək yeni hash yaradır.

**Mənbə:** Chapter 14, page 226
