# Telegram Channel Template Guide

Bu template istənilən yeni kanal üçün başlanğıc nümunəsidir.
Bütün mövcud sahələr və qəbul edilən dəyərlər aşağıda izah edilib.

## config.json strukturü

```json
{
  "channel": {
    "name": "RA_EXAMPLE",
    "enabled": true,
    "description": "Optional description of what this channel is for"
  },
  "rules": {
    "match": "*",
    "languages": ["Go", "JavaScript", "Python"],
    "technologies": ["Docker", "Kubernetes", "gRPC", "Redis"],
    "domains": ["Backend", "Frontend", "DevOps", "Microservices", "Database"],
    "levels": [1, 2, 3, 4, 5]
  },
  "telegram": {
    "send_file": true,
    "send_metadata": true,
    "send_caption": true,
    "send_teacher_notes": false
  },
  "behavior": {
    "priority": 100,
    "max_retries": 3,
    "retry_delay_seconds": 60
  }
}
```

## Sahələr izahı

### `channel`
| Sahə | Tipe | Təsviri |
|------|------|---------|
| `name` | string | Kanalın unikal adı. `RA_` prefiksi ilə başlamalıdır. |
| `enabled` | bool | `false` olarsa, bu kanal upload prosesindən tamamilə çıxarılır. |
| `description` | string | Təsvir. Xarici API üçün deyil, insan oxuması üçündür. |

### `rules`
Burada kitabın bu kanala düşüb düşməməsini təyin edən qaydalar var.

| Sahə | Tipe | Təsviri |
|------|------|---------|
| `match` | string | `"*"` yazılırsa, **bütün kitablar** bu kanala düşür. Digər qaydalar yoxlanılmır. |
| `languages` | string array | Kitabın `primary_language` dəyəri bunlardan biri ilə eynidirsə uyğun gəlir. Məsələn: `["Go"]` |
| `technologies` | string array | Kitabın `technologies` siyahısında bunlardan ən az biri varsa uyğun gəlir. Məsələn: `["Docker", "Kubernetes"]` |
| `domains` | string array | Kitabın `domains` siyahısında bunlardan ən az biri varsa uyğun gəlir. Məsələn: `["Backend", "Microservices"]` |
| `levels` | number array | Kitabın `level` dəyəri bunlardan biri ilə eynidirsə uyğun gəlir. Məsələn: `[3, 4, 5]` |

**Qayda qarışdırma qeydi:** `match: "*"` varsa, digər bütün qaydalar nəzərə alınmır. `match` olmadan isə qaydalar **OR** məntiqi ilə işləyir — yalnız birinin uyğun gəlməsi kifayətdir.

### `telegram`
| Sahə | Tipe | Təsviri |
|------|------|---------|
| `send_file` | bool | `true` olarsa, kitab faylı (`.pdf` və s.) Telegram-a göndərilir. |
| `send_metadata` | bool | `true` olarsa, kitab haqqında məlumat mesajda göstərilir. |
| `send_caption` | bool | `true` olarsa, faylın altında başlıq (caption) yazılır. |
| `send_teacher_notes` | bool | `true` olarsa, müəllim qeydləri də mesaja əlavə olunur. |

### `behavior`
| Sahə | Tipe | Təsviri |
|------|------|---------|
| `priority` | number | Birdən çox kanala uyğun gələndə yüksək prioritetli kanal əvvəl işlənir. Məsələn: `100` |
| `max_retries` | number | Upload uğursuz olsa neçə dəfə təkrar cəhd edilsin. Məsələn: `3` |
| `retry_delay_seconds` | number | Təkrar cəhdlər arasındakı gözləmə müddəti (saniyə). Məsələn: `60` |

## Nümunə konfiqurasiyalar

### 1. Bütün kitablar kanalı (Catch-All)
```json
{
  "channel": {"name": "RA_ALL", "enabled": true},
  "rules": {"match": "*"}
}
```

### 2. Yalnız Go kitabları
```json
{
  "channel": {"name": "RA_GO", "enabled": true},
  "rules": {"languages": ["Go"]}
}
```

### 3. Go + Intermediate/Advanced
```json
{
  "channel": {"name": "RA_GO_ADVANCED", "enabled": true},
  "rules": {
    "languages": ["Go"],
    "levels": [3, 4, 5]
  }
}
```

### 4. Backend + Docker
```json
{
  "channel": {"name": "RA_BACKEND_TOOLING", "enabled": true},
  "rules": {
    "technologies": ["Docker", "Kubernetes"],
    "domains": ["Backend", "DevOps"]
  }
}
```

### 5. Sadəcə Beginner kitablar
```json
{
  "channel": {"name": "RA_BEGINNER", "enabled": true},
  "rules": {"levels": [1, 2]}
}
```

## Yeni kanal yaratmaq addım-addım

1. Yeni qovluq yarad: `telegram/channels/RA_{AD}/`
2. Həmin qovluğa `config.json` və `.env` yerləşdir
3. `.env`-ə `TELEGRAM_BOT_TOKEN` və `TELEGRAM_CHAT_ID` yaz
4. `config.json`-u yuxarıdakı nümunələrə əsasən doldur
5. Proqramı yenidən compile etməyə ehtiyac yox — uploader avtomatik skan edir

## Təhlükəsizlik qeydləri

- Bot token və chat ID `.env` faylında saxlanmalıdır, `config.json`-da YOX.
- `.env` faylını Git-ə əlavə etmə.
- `telegram/channels/*/.env` `.gitignore`-da olmalıdır.
- Bir bot tokeni bir kanal üçün nəzərdə tutulub, amma eyni tokeni birdən çox kanalda istifadə etmək mümkündür.
