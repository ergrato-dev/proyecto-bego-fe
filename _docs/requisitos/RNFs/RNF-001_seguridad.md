# RNF-001 — Seguridad

## Identificación

| Campo         | Valor                 |
| ------------- | --------------------- |
| **ID**        | RNF-001               |
| **Nombre**    | Seguridad del sistema |
| **Categoría** | Seguridad             |
| **Prioridad** | Crítica               |

---

## Descripción

El sistema debe implementar controles de seguridad robustos en todas las capas (backend, base de datos, comunicación y frontend) para proteger los datos de los usuarios y prevenir accesos no autorizados.

---

## Requisitos de Seguridad

### Contraseñas

| Requisito      | Detalle                                                     |
| -------------- | ----------------------------------------------------------- |
| Hashing        | Bcrypt con `bcrypt.DefaultCost` (≥ 10 rondas)               |
| Almacenamiento | Nunca en texto plano — solo hash en `users.hashed_password` |
| Validación     | Mínimo 8 chars, al menos 1 mayúscula, 1 minúscula, 1 dígito |
| Logs           | Jamás registrar contraseñas en logs                         |

### JWT (Tokens)

| Parámetro      | Valor                                         |
| -------------- | --------------------------------------------- |
| Algoritmo      | HS256 (`jwt.SigningMethodHS256`)              |
| Access Token   | TTL: 15 minutos                               |
| Refresh Token  | TTL: 7 días                                   |
| Secret key     | Mínimo 32 caracteres, generada aleatoriamente |
| Almacenamiento | En memoria (no `localStorage` en producción)  |

### Cabeceras HTTP de Seguridad

El middleware `security.go` debe incluir:

| Cabecera                    | Valor requerido                          |
| --------------------------- | ---------------------------------------- |
| `X-Content-Type-Options`    | `nosniff`                                |
| `X-Frame-Options`           | `DENY`                                   |
| `X-XSS-Protection`          | `1; mode=block`                          |
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains`    |
| `Content-Security-Policy`   | Restringir orígenes de scripts y estilos |
| `Referrer-Policy`           | `strict-origin-when-cross-origin`        |

### CORS

| Entorno    | Orígenes permitidos                   |
| ---------- | ------------------------------------- |
| Desarrollo | `http://localhost:5173` únicamente    |
| Producción | Dominio específico del frontend       |
| Prohibido  | `AllowAllOrigins: true` en producción |

### Rate Limiting

| Endpoint                     | Límite              |
| ---------------------------- | ------------------- |
| `POST /auth/login`           | 5 intentos / 15 min |
| `POST /auth/register`        | 10 intentos / hora  |
| `POST /auth/forgot-password` | 3 intentos / 15 min |
| `POST /auth/change-password` | 5 intentos / 15 min |

---

## Cobertura OWASP Top 10

| Categoría OWASP                 | Control implementado                     |
| ------------------------------- | ---------------------------------------- |
| A01 — Broken Access Control     | JWT middleware en rutas protegidas       |
| A02 — Cryptographic Failures    | Bcrypt para contraseñas, HS256 para JWT  |
| A03 — Injection                 | GORM con queries parametrizadas          |
| A04 — Insecure Design           | Mismo mensaje en login/forgot-password   |
| A05 — Security Misconfiguration | CORS restringido, cabeceras de seguridad |
| A07 — Auth & Session Failures   | Rate limiting, tokens de corta duración  |
| A09 — Security Logging Failures | `audit_log.go` registra eventos de auth  |

---

## Herramienta de Auditoría

```bash
# Detectar vulnerabilidades en dependencias Go
govulncheck ./...

# Análisis estático de seguridad
golangci-lint run --enable gosec ./...
```

---

## Criterio de Cumplimiento

- Sin contraseñas en texto plano en la BD (verificable)
- `golangci-lint run` sin errores críticos de seguridad
- Cabeceras HTTP verificadas con `curl -I http://localhost:8000`
- CORS no permite orígenes comodín en desarrollo
