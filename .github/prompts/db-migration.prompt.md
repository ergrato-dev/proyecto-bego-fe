---
description: "Genera una migración golang-migrate (SQL puro) y actualiza el modelo GORM correspondiente. Usar cuando se necesite crear o modificar tablas, columnas, índices o constraints en PostgreSQL."
name: "Migración golang-migrate + GORM"
argument-hint: "Describe el cambio en la BD: qué tabla/columna se crea, modifica o elimina, y por qué"
agent: "agent"
---

# Migración golang-migrate — NN Auth System

Genera la migración SQL y actualiza el modelo GORM para el cambio de esquema indicado.

## Convenciones obligatorias

- **NUNCA** modificar la base de datos manualmente — siempre vía golang-migrate
- **NUNCA** editar una migración ya ejecutada — crear una nueva
- Nombres de tablas: `snake_case`, plural (`users`, `password_reset_tokens`)
- Nombres de columnas: `snake_case` (`created_at`, `hashed_password`)
- Toda tabla debe tener `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- Tablas actualizables deben tener `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- Primary keys: UUID con `DEFAULT gen_random_uuid()` (extensión `pgcrypto`)
- Foreign keys: nombrar `fk_<tabla>_<columna>` y siempre declarar `ON DELETE CASCADE` o `RESTRICT` explícitamente

## Archivos de referencia

- Modelo User: [be/internal/models/user.go](../../../be/internal/models/user.go)
- Modelos existentes: [be/internal/models/](../../../be/internal/models/)
- Migraciones existentes: [be/migrations/](../../../be/migrations/)
- Configuración GORM: [be/internal/database/database.go](../../../be/internal/database/database.go)

## Lo que debes generar

### 1. Archivo `.up.sql` (`be/migrations/`)

El nombre sigue el patrón: `<número>_<descripcion_snake_case>.up.sql`
El número es el siguiente en la secuencia (ej. `000004_...` si ya existen 3 migraciones).

```sql
-- ¿Qué? Migración que agrega la columna locale a la tabla users.
-- ¿Para qué? Persistir la preferencia de idioma del usuario entre sesiones.
-- ¿Impacto? Sin esta columna, el idioma se resetea en cada login.

ALTER TABLE users
    ADD COLUMN locale VARCHAR(10) NOT NULL DEFAULT 'es';
```

### 2. Archivo `.down.sql` (`be/migrations/`)

Nombre: `<mismo_número>_<misma_descripcion>.down.sql`

**CRÍTICO**: debe deshacer EXACTAMENTE lo que hace el `.up.sql`.

```sql
-- ¿Qué? Revierte la adición de la columna locale.
-- ¿Para qué? Permitir rollback limpio de la migración.
-- ¿Impacto? Se pierde el dato de locale de todos los usuarios existentes.

ALTER TABLE users
    DROP COLUMN IF EXISTS locale;
```

### 3. Actualizar el modelo GORM (`be/internal/models/`)

Si la columna nueva debe aparecer en queries, actualizar el struct correspondiente:

```go
// ¿Qué? Struct GORM que mapea la tabla 'users' en PostgreSQL.
// ¿Para qué? Permite a GORM hacer queries y operaciones CRUD sin SQL manual.
// ¿Impacto? Si el struct no refleja la tabla, GORM puede generar queries incorrectos.
type User struct {
    // ... campos existentes ...

    // ¿Qué? Preferencia de idioma del usuario (código ISO 639-1).
    // ¿Para qué? Persistir el idioma seleccionado entre sesiones.
    // ¿Impacto? Sin este campo, el Frontend no puede recuperar la preferencia guardada.
    Locale string `gorm:"not null;default:es" json:"locale"`
}
```

### 4. Actualizar DTOs si corresponde (`be/internal/dto/`)

Si la columna nueva debe aparecer en requests o responses, actualizar los structs
en [be/internal/dto/auth.go](../../../be/internal/dto/auth.go) u otro archivo DTO.

## Cómo aplicar la migración

```bash
cd be

# Aplicar todas las migraciones pendientes
go run ./cmd/migrate/main.go up

# Verificar estado actual
go run ./cmd/migrate/main.go version

# Si algo falla, revertir la última migración
go run ./cmd/migrate/main.go down 1
```

## Tipos de datos comunes en este proyecto

| GORM (Go)                    | PostgreSQL         | Uso                          |
| ---------------------------- | ------------------ | ---------------------------- |
| `string` + `gorm:"size:255"` | `VARCHAR(255)`     | Texto corto (email, nombre)  |
| `string`                     | `TEXT`             | Texto largo sin límite       |
| `bool`                       | `BOOLEAN`          | Flags activo/inactivo        |
| `time.Time`                  | `TIMESTAMPTZ`      | Fechas con timezone          |
| `uuid.UUID`                  | `UUID`             | Primary keys                 |
| `*time.Time`                 | `TIMESTAMPTZ NULL` | Fechas opcionales (nullable) |

## Descripción del cambio de esquema a implementar

$input
