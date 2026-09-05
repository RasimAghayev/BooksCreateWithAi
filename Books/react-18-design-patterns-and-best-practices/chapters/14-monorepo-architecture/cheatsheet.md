# Cheat Sheet — MonoRepo Architecture

## MonoRepo (NPM Workspaces)

### Root package.json

**Nə edir:** Workspace konfiqurasiyası — alt paketlər bir yerdə qurulur.

**Kod:**

```json
{
  "name": "my-monorepo",
  "private": true,
  "workspaces": ["packages/*", "apps/*"]
}
```

**Mənbə:** Chapter 14, page 354-412

### Workspace install

**Nə edir:** Müəyyən workspace-ə dependency əlavə etmək.

**Kod:**

```shell
npm install lodash -w @my-org/utils
```

**Mənbə:** Chapter 14, page 354-412

### tsconfig.base.json

**Nə edir:** Bütün paketlər üçün ümumi TS qaydaları.

**Kod:**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "strict": true
  }
}
```

**Mənbə:** Chapter 14, page 354-412

### tsconfig.json (paket)

**Nə edir:** Hər paket extends edir, spesifik ayarları əlavə edir.

**Kod:**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": { "outDir": "./dist" },
  "include": ["src"]
}
```

**Mənbə:** Chapter 14, page 354-412
