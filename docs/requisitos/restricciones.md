# Restricciones del Proyecto — NN Auth System

<!--
  ¿Qué? Lista de restricciones obligatorias que delimitan las decisiones técnicas
        y de proceso del proyecto.
  ¿Para qué? Que todo el equipo entienda qué NO es negociable y por qué.
  ¿Impacto? Sin restricciones definidas, cada desarrollador podría usar
             herramientas diferentes, generando un proyecto imposible de mantener.
-->

---

## Tipos de Restricciones

| Código | Tipo                        | Descripción                           |
| ------ | --------------------------- | ------------------------------------- |
| RT     | Restricción Tecnológica     | Herramientas y versiones obligatorias |
| RH     | Restricción de Herramientas | Comandos y flujos de trabajo          |
| RD     | Restricción de Diseño       | Estilo visual y UX                    |
| RL     | Restricción de Lenguaje     | Idioma del código y documentación     |
| RO     | Restricción Organizacional  | Proceso y calidad                     |
| RS     | Restricción de Seguridad    | Prácticas de seguridad obligatorias   |

---

## Restricciones Tecnológicas (RT)

### RT-001 — Backend: Go + Gin + GORM + golang-migrate

**Restricción:** El backend DEBE implementarse con el siguiente stack exacto:

| Componente  | Tecnología                 | Versión |
| ----------- | -------------------------- | ------- |
| Lenguaje    | Go                         | 1.22+   |
| Framework   | Gin                        | 1.9+    |
| ORM         | GORM                       | 2.0+    |
| Migraciones | golang-migrate             | 4+      |
| JWT         | golang-jwt/jwt             | v5      |
| Bcrypt      | golang.org/x/crypto/bcrypt | latest  |
| Validación  | go-playground/validator    | v10     |
| Config      | joho/godotenv              | latest  |
| Rate limit  | github.com/ulule/limiter   | v3      |
| Email       | jordan-wright/email        | latest  |
| Driver BD   | jackc/pgx                  | v5      |

**Motivación:** El stack fue elegido por su madurez, adopción en la industria y adecuación a proyectos educativos. Cambiar cualquier componente rompería la coherencia del proyecto y los materiales de enseñanza.

**No permitido:**

- Usar Echo, Fiber, Chi u otro framework HTTP en lugar de Gin
- Usar sqlx o ent en lugar de GORM
- Usar `database/sql` directamente para queries de negocio
- Usar python-jose o librerías de otros lenguajes

---

### RT-002 — Frontend: React + Vite + TypeScript + TailwindCSS + React Router

**Restricción:** El frontend DEBE implementarse con:

| Componente   | Tecnología             | Versión |
| ------------ | ---------------------- | ------- |
| UI library   | React                  | 18+     |
| Build tool   | Vite                   | 6+      |
| Lenguaje     | TypeScript             | 5.0+    |
| CSS          | TailwindCSS            | 4+      |
| Enrutamiento | React Router           | 7+      |
| HTTP client  | Axios                  | latest  |
| Testing      | Vitest                 | latest  |
| Testing DOM  | @testing-library/react | latest  |

**No permitido:**

- Usar Next.js, Remix u otro meta-framework
- Usar Vue, Angular u otra librería UI
- Usar CSS Modules, styled-components o SASS en lugar de TailwindCSS
- Usar fetch nativo en lugar de Axios para llamadas a la API

---

### RT-003 — Base de Datos: PostgreSQL 17+

**Restricción:** La base de datos DEBE ser PostgreSQL versión 17 o superior.

**No permitido:**

- MySQL, MariaDB, SQLite u otros motores
- Bases de datos NoSQL (MongoDB, Redis como BD principal)
- Versiones anteriores a PostgreSQL 17

**Motivación:** PostgreSQL 17 incluye mejoras de performance y seguridad. La función `gen_random_uuid()` (del módulo `pgcrypto`) se usa para UUIDs.

---

### RT-004 — Autenticación: JWT Stateless

**Restricción:** La autenticación DEBE ser JWT stateless con:

- Access token: duración máxima de 15 minutos
- Refresh token: duración máxima de 7 días
- Algoritmo: HS256 (`golang-jwt/jwt v5`)
- Header de autenticación: `Authorization: Bearer <token>`

**No permitido:**

- Sesiones en servidor (cookies de sesión, session store)
- Almacenar tokens en `localStorage` (solo en memoria o `sessionStorage` temporal)
- Tokens de larga duración como único mecanismo (sin refresh token)

---

### RT-005 — Hashing de Contraseñas: bcrypt

**Restricción:** Todas las contraseñas DEBEN hashearse con bcrypt usando `golang.org/x/crypto/bcrypt` con `DefaultCost` (10).

**No permitido:** MD5, SHA1, SHA256, SHA512 u otros algoritmos de hash simples para contraseñas.

**Motivación:** bcrypt es un algoritmo de derivación de clave (KDF) diseñado para ser lento, haciendo inviables los ataques de fuerza bruta. Los hashes simples como SHA256 pueden calcularse millions de veces por segundo.

---

## Restricciones de Herramientas (RH)

### RH-001 — Go Modules Obligatorio

**Restricción:** El proyecto de Go DEBE usar Go Modules (`go.mod` y `go.sum`).

```bash
# ✅ CORRECTO
go mod download     # descargar dependencias
go mod tidy         # limpiar dependencias no usadas
go get github.com/gin-gonic/gin  # agregar dependencia

# ❌ PROHIBIDO
# No usar vendor/ manual sin go.mod
# No usar GOPATH sin módulos
```

**Motivación:** Go modules es el estándar oficial desde Go 1.11 y obligatorio desde Go 1.16. Garantiza builds reproducibles.

---

### RH-002 — pnpm Obligatorio para Frontend

**Restricción:** El frontend DEBE usar `pnpm` como gestor de paquetes.

```bash
# ✅ CORRECTO
pnpm install
pnpm add axios
pnpm add -D vitest
pnpm dev
pnpm test
pnpm build

# ❌ PROHIBIDO
npm install     # No usar npm
npm run dev     # No usar npm
yarn install    # No usar yarn
yarn dev        # No usar yarn
```

**Motivación:** pnpm es significativamente más eficiente en espacio de disco (hard links compartidos) y más rápido que npm/yarn. Es el standard en este proyecto.

---

### RH-003 — golangci-lint Obligatorio

**Restricción:** Todo código Go DEBE pasar `golangci-lint` sin errores antes de hacer commit.

```bash
# Verificar errores de lint
cd be && golangci-lint run ./...

# Formatear código
cd be && gofmt -w .
```

**Configuración requerida:** El proyecto DEBE incluir un `.golangci.yml` con la configuración de linters habilitados.

---

### RH-004 — ESLint + Prettier Obligatorio para Frontend

**Restricción:** Todo código TypeScript/React DEBE pasar ESLint y Prettier antes de hacer commit.

```bash
# Verificar errores de lint
cd fe && pnpm lint

# Formatear código
cd fe && pnpm format
```

---

### RH-005 — Variables de Entorno en .env

**Restricción:** Toda configuración sensible DEBE estar en archivos `.env` (no versionados). NUNCA hardcodear en el código.

**Obligatorio:**

- `be/.env` — Variables del backend (no en git)
- `fe/.env` — Variables del frontend (no en git)
- `be/.env.example` — Plantilla con valores de ejemplo (sí en git)
- `fe/.env.example` — Plantilla con valores de ejemplo (sí en git)

---

## Restricciones de Diseño (RD)

### RD-001 — Sin Gradientes

**Restricción:** PROHIBIDO usar degradados (`gradient`) en cualquier elemento del UI.

```tsx
// ❌ PROHIBIDO
<div className="bg-gradient-to-r from-blue-500 to-purple-500">
<div style={{ background: 'linear-gradient(...)' }}>

// ✅ CORRECTO — colores sólidos y planos
<div className="bg-blue-600">
```

---

### RD-002 — Dark Mode y Light Mode Obligatorios

**Restricción:** La aplicación DEBE implementar ambos modos con toggle. Debe respetar `prefers-color-scheme` del sistema operativo como estado inicial.

---

### RD-003 — Fuentes Sans-Serif Exclusivamente

**Restricción:** Solo fuentes sans-serif (Inter, system-ui, sans). Prohibidas las fuentes serif y monospace en textos del UI (solo en bloques de código).

---

### RD-004 — Botones de Acción Alineados a la Derecha

**Restricción:** En formularios, los botones de acción principal (Submit, Guardar, Siguiente) DEBEN alinearse a la derecha (`justify-end`).

---

### RD-005 — Mobile First

**Restricción:** El diseño DEBE implementarse mobile first. Los formularios deben verse correctos en pantallas de 320px de ancho.

---

## Restricciones de Lenguaje (RL)

### RL-001 — Código en Inglés

**Restricción:** Todo código DEBE estar en inglés: nombres de variables, funciones, structs, métodos, endpoints, tablas de BD, columnas, nombres de archivos.

```go
// ✅ CORRECTO
func GetUserByEmail(email string) (*models.User, error) {}

// ❌ INCORRECTO
func ObtenerUsuarioPorEmail(correo string) (*models.User, error) {}
```

---

### RL-002 — Comentarios y Documentación en Español

**Restricción:** Todos los comentarios en el código, archivos de documentación (`.md`) y el README DEBEN estar en español.

---

### RL-003 — Comentario Pedagógico Obligatorio

**Restricción:** Todo bloque de código significativo DEBE tener un comentario con el formato ¿Qué? / ¿Para qué? / ¿Impacto?

```go
// ¿Qué? Función que hashea la contraseña con bcrypt.
// ¿Para qué? Almacenar contraseñas de forma segura, nunca en texto plano.
// ¿Impacto? Si se omite, una filtración de BD expone todas las contraseñas.
func HashPassword(password string) (string, error) { ... }
```

---

## Restricciones Organizacionales (RO)

### RO-001 — Conventional Commits

**Restricción:** Todos los commits DEBEN seguir el formato Conventional Commits con descripción What/For/Impact en el cuerpo.

```
feat(auth): add user registration handler

What: Creates POST /api/v1/auth/register endpoint
For: Allow new users to create accounts
Impact: Enables the user onboarding flow
```

---

### RO-002 — Tests Obligatorios

**Restricción:** Todo código nuevo de lógica de negocio (`services/`) DEBE tener tests con cobertura mínima del 80%.

---

### RO-003 — No se puede avanzar de fase sin completar la anterior

**Restricción:** El plan de trabajo está dividido en 8 fases. No se puede iniciar la Fase N+1 sin que los criterios de verificación de la Fase N hayan sido cumplidos.

---

## Restricciones de Seguridad (RS)

### RS-001 — Contraseñas Nunca en Texto Plano

**Restricción:** Está PROHIBIDO almacenar, loggear o retornar contraseñas en texto plano en cualquier parte del sistema.

---

### RS-002 — Secrets Nunca Hardcodeados

**Restricción:** Está PROHIBIDO hardcodear en el código: JWT secrets, contraseñas de BD, API keys, URLs de servicios externos.

---

### RS-003 — Manejo Explícito de Errores en Go

**Restricción:** En Go, NUNCA se puede ignorar un error con `_`. Todo error DEBE ser manejado explícitamente.

```go
// ✅ CORRECTO
bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return "", fmt.Errorf("hashing password: %w", err)
}

// ❌ PROHIBIDO
bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

---

### RS-004 — Rate Limiting en Endpoints de Auth

**Restricción:** Los endpoints de autenticación (`/register`, `/login`, `/forgot-password`) DEBEN tener rate limiting configurado vía `ulule/limiter`.

---

## Resumen de Restricciones

| Código | Categoría    | Restricción Clave                                 |
| ------ | ------------ | ------------------------------------------------- |
| RT-001 | Tecnológica  | Go + Gin + GORM + golang-migrate (obligatorio)    |
| RT-002 | Tecnológica  | React + Vite + TS + TailwindCSS (obligatorio)     |
| RT-003 | Tecnológica  | PostgreSQL 17+ (obligatorio)                      |
| RT-004 | Tecnológica  | JWT stateless, access 15min, refresh 7días        |
| RT-005 | Tecnológica  | bcrypt para contraseñas (obligatorio)             |
| RH-001 | Herramientas | Go modules (no GOPATH sin módulos)                |
| RH-002 | Herramientas | pnpm (nunca npm ni yarn)                          |
| RH-003 | Herramientas | golangci-lint antes de commit                     |
| RH-004 | Herramientas | ESLint + Prettier antes de commit                 |
| RH-005 | Herramientas | Variables sensibles en .env (no hardcoded)        |
| RD-001 | Diseño       | Sin gradientes en la UI                           |
| RD-002 | Diseño       | Dark mode + Light mode obligatorios               |
| RD-005 | Diseño       | Mobile first                                      |
| RL-001 | Lenguaje     | Código en inglés                                  |
| RL-002 | Lenguaje     | Comentarios en español                            |
| RL-003 | Lenguaje     | Comentario ¿Qué?/¿Para qué?/¿Impacto? obligatorio |
| RO-001 | Organización | Conventional Commits con What/For/Impact          |
| RO-002 | Organización | Tests con 80% de cobertura en services/           |
| RS-001 | Seguridad    | Contraseñas nunca en texto plano                  |
| RS-003 | Seguridad    | Nunca ignorar errores con \_ en Go                |
| RS-004 | Seguridad    | Rate limiting en endpoints de auth                |
