# Esquema de Base de Datos — NN Auth System

<!--
  ¿Qué? Documentación del esquema de PostgreSQL usando golang-migrate.
  ¿Para qué? Que todo el equipo entienda la estructura de la BD,
             las relaciones entre tablas y cómo evolucionan las migraciones.
  ¿Impacto? Un esquema mal diseñado es difícil o imposible de cambiar en producción
             sin perder datos. Estas decisiones afectan todo el ciclo de vida del sistema.
-->

---

## 1. Herramienta de Migraciones: golang-migrate

### ¿Qué es golang-migrate?

`golang-migrate` es la herramienta estándar para versionado de esquema de BD en Go. A diferencia de Alembic (Python) que genera código Python, `golang-migrate` usa **archivos SQL puros** numerados secuencialmente.

### Convención de nombres de archivos

```
migrations/
├── 000001_create_users.up.sql           ← Aplica la migración
├── 000001_create_users.down.sql         ← Revierte la migración
├── 000002_create_password_reset_tokens.up.sql
├── 000002_create_password_reset_tokens.down.sql
├── 000003_create_email_verification_tokens.up.sql
└── 000003_create_email_verification_tokens.down.sql
```

Formato: `{número_secuencial}_{descripción}.{up|down}.sql`

### Comandos de migración

```bash
# Desde la raíz del proyecto be/

# Aplicar todas las migraciones pendientes
go run ./cmd/migrate/main.go up

# Revertir la última migración aplicada (un paso)
go run ./cmd/migrate/main.go down

# Ver el estado actual
go run ./cmd/migrate/main.go version
```

> **Nota:** `down` revierte exactamente **un paso** (la migración más reciente).
> Para revertir más pasos, ejecutarlo múltiples veces.

---

## 2. Diagrama Entidad-Relación

```
┌──────────────────────────────────────┐
│               users                  │
├──────────────────────────────────────┤
│ id               UUID  PK            │
│ email            VARCHAR(255) UNIQUE │
│ full_name        VARCHAR(255)        │
│ hashed_password  VARCHAR(255)        │
│ is_active        BOOLEAN             │
│ is_email_verified BOOLEAN            │
│ created_at       TIMESTAMPTZ         │
│ updated_at       TIMESTAMPTZ         │
└──────────┬──────────────────┬────────┘
           │                  │
           │ 1                │ 1
           │ N                │ N
           ▼                  ▼
┌─────────────────────┐  ┌──────────────────────────┐
│ password_reset_     │  │ email_verification_       │
│ tokens              │  │ tokens                    │
├─────────────────────┤  ├──────────────────────────┤
│ id       UUID PK    │  │ id       UUID PK          │
│ user_id  UUID FK──►users  user_id  UUID FK──►users │
│ token    VARCHAR    │  │ token    VARCHAR           │
│ expires_at TSTZ     │  │ expires_at TSTZ           │
│ used     BOOLEAN    │  │ used     BOOLEAN           │
│ created_at TSTZ     │  │ created_at TSTZ           │
└─────────────────────┘  └──────────────────────────┘
```

---

## 3. Tabla `users`

Almacena a todos los usuarios registrados del sistema.

### Definición SQL (`000001_create_users.up.sql`)

```sql
CREATE TABLE IF NOT EXISTS users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email               VARCHAR(255) NOT NULL UNIQUE,
    full_name           VARCHAR(255) NOT NULL,
    hashed_password     VARCHAR(255) NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    is_email_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índice para búsquedas rápidas por email (login)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

### Reversión SQL (`000001_create_users.down.sql`)

```sql
DROP TABLE IF EXISTS users;
```

### Descripción de columnas

| Columna             | Tipo         | Default             | Descripción                                           |
| ------------------- | ------------ | ------------------- | ----------------------------------------------------- |
| `id`                | UUID         | `gen_random_uuid()` | Identificador único. UUID es más seguro que INTEGER   |
| `email`             | VARCHAR(255) | —                   | Email del usuario. UNIQUE garantiza unicidad          |
| `full_name`         | VARCHAR(255) | —                   | Nombre completo del usuario                           |
| `hashed_password`   | VARCHAR(255) | —                   | Hash bcrypt. NUNCA almacenar texto plano              |
| `is_active`         | BOOLEAN      | `TRUE`              | Permite deshabilitar cuentas sin borrarlas            |
| `is_email_verified` | BOOLEAN      | `FALSE`             | Solo usuarios verificados pueden hacer login          |
| `created_at`        | TIMESTAMPTZ  | `NOW()`             | Fecha de registro (`TZ` = with timezone, recomendado) |
| `updated_at`        | TIMESTAMPTZ  | `NOW()`             | Última actualización. Actualizar manualmente          |

### ¿Por qué UUID en lugar de INTEGER auto-incremental?

| Aspecto      | UUID                                    | INTEGER Serial                         |
| ------------ | --------------------------------------- | -------------------------------------- |
| Seguridad    | No predecible (no se puede adivinar)    | Predecible (1, 2, 3...)                |
| Distribución | Funciona en multi-servidor sin colisión | Requiere coordinación entre servidores |
| Exposición   | No revela volumen de usuarios           | Revela cuántos usuarios hay (id=1000)  |
| Performance  | Ligeramente más lento                   | Ligeramente más rápido                 |

---

## 4. Tabla `password_reset_tokens`

Almacena tokens temporales para el flujo de recuperación de contraseña.

### Definición SQL (`000002_create_password_reset_tokens.up.sql`)

```sql
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índice para búsqueda rápida por token (reset-password endpoint)
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token
    ON password_reset_tokens(token);

-- Índice para buscar tokens por usuario
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id
    ON password_reset_tokens(user_id);
```

### Reversión SQL (`000002_create_password_reset_tokens.down.sql`)

```sql
DROP TABLE IF EXISTS password_reset_tokens;
```

### Descripción de columnas

| Columna      | Tipo         | Descripción                                                          |
| ------------ | ------------ | -------------------------------------------------------------------- |
| `id`         | UUID         | Identificador del registro                                           |
| `user_id`    | UUID FK      | Usuario al que pertenece. CASCADE: si se borra el user, borra tokens |
| `token`      | VARCHAR(255) | Token UUID generado aleatoriamente. UNIQUE para búsqueda segura      |
| `expires_at` | TIMESTAMPTZ  | Expiración del token (1 hora tras creación)                          |
| `used`       | BOOLEAN      | Si ya se usó. Previene reutilización del mismo token                 |
| `created_at` | TIMESTAMPTZ  | Cuándo se creó el token de reset                                     |

### Flujo de uso

```
1. Usuario solicita reset → Se crea un registro (used=FALSE, expires_at=NOW()+1h)
2. Usuario hace clic en el enlace o introd. el token
3. Backend verifica: token existe, expires_at > NOW(), used=FALSE
4. Se hashea y actualiza la contraseña en users
5. Se actualiza: used=TRUE (el token no puede volver a usarse)
```

---

## 5. Tabla `email_verification_tokens`

Almacena tokens para verificar la dirección de email de nuevos usuarios.

### Definición SQL (`000003_create_email_verification_tokens.up.sql`)

```sql
CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_token
    ON email_verification_tokens(token);

CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id
    ON email_verification_tokens(user_id);
```

### Reversión SQL (`000003_create_email_verification_tokens.down.sql`)

```sql
DROP TABLE IF EXISTS email_verification_tokens;
```

### Descripción de columnas

Misma estructura que `password_reset_tokens`, pero:

| Campo        | Diferencia                                          |
| ------------ | --------------------------------------------------- |
| `expires_at` | Expiración de 24 horas (en lugar de 1 hora)         |
| `used`       | Al usarse, actualiza `users.is_email_verified=TRUE` |

### Flujo de uso

```
1. Registro → Se crea un registro (used=FALSE, expires_at=NOW()+24h)
2. Usuario hace clic en el enlace del email de bienvenida
3. Backend verifica: token existe, expires_at > NOW(), used=FALSE
4. Se actualiza: users.is_email_verified=TRUE
5. Se actualiza: used=TRUE en el token
6. El usuario ya puede hacer login
```

---

## 6. Modelo GORM en Go

### `internal/models/user.go` (referencia)

El struct de Go que mapea a la tabla `users`:

```go
// Modelo GORM que representa la tabla users.
// ¿Para qué? Permite a GORM mapear columnas SQL a campos del struct Go.
// ¿Impacto? Este struct es la fuente de verdad para GORM; debe reflejar
//            exactamente las columnas de la migración SQL.
type User struct {
    ID                UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Email             string    `gorm:"uniqueIndex;not null;size:255"`
    FullName          string    `gorm:"not null;size:255"`
    HashedPassword    string    `gorm:"not null;size:255"`
    IsActive          bool      `gorm:"not null;default:true"`
    IsEmailVerified   bool      `gorm:"not null;default:false"`
    CreatedAt         time.Time `gorm:"not null;autoCreateTime"`
    UpdatedAt         time.Time `gorm:"not null;autoUpdateTime"`
}
```

> **Importante:** Los modelos GORM NO se usan para crear las tablas en este proyecto. Las tablas se crean con las **migraciones SQL de golang-migrate**. GORM solo los usa para hacer queries.

---

## 7. Reglas de Integridad Referencial

| Relación                            | Comportamiento      | Razón                                      |
| ----------------------------------- | ------------------- | ------------------------------------------ |
| `password_reset_tokens.user_id`     | `ON DELETE CASCADE` | Al borrar un usuario, se borran sus tokens |
| `email_verification_tokens.user_id` | `ON DELETE CASCADE` | Al borrar un usuario, se borran sus tokens |

---

## 8. Índices y Performance

| Índice                                | Columna       | Justificación                                  |
| ------------------------------------- | ------------- | ---------------------------------------------- |
| `idx_users_email`                     | `users.email` | Login requiere buscar por email frecuentemente |
| `idx_password_reset_tokens_token`     | `token`       | Reset-password busca por token en cada intento |
| `idx_email_verification_tokens_token` | `token`       | Verify-email busca por token en cada intento   |

---

## 9. Verificar el Esquema

```bash
# Verificar que las tablas fueron creadas (requiere psql instalado o Docker)
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db -c "\dt"

# Ver la estructura de la tabla users
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db -c "\d users"

# Ver el estado actual de las migraciones
go run ./cmd/migrate/main.go version
```

Salida esperada de `\dt`:

```
             List of relations
 Schema |             Name              | Type  |  Owner
--------+-------------------------------+-------+---------
 public | schema_migrations             | table | nn_user
 public | users                         | table | nn_user
 public | password_reset_tokens         | table | nn_user
 public | email_verification_tokens     | table | nn_user
```

> `schema_migrations` es la tabla interna de golang-migrate que registra qué migraciones ya se aplicaron.
