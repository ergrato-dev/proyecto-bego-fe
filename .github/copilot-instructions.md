# 🎓 Instrucciones del Proyecto — NN Auth System

<!--
  ¿Qué? Archivo de instrucciones y convenciones para el desarrollo del proyecto.
  ¿Para qué? Garantizar coherencia de código, estilo y calidad en todo el equipo.
  ¿Impacto? Sin estas reglas, cada desarrollador usaría su propio criterio,
             generando un proyecto inconsistente y difícil de mantener.
-->

---

## 1. Identidad del Proyecto

| Campo        | Valor                                                                                    |
| ------------ | ---------------------------------------------------------------------------------------- |
| Nombre       | NN Auth System                                                                           |
| Tipo         | Proyecto educativo — SENA                                                                |
| Propósito    | Sistema de autenticación completo (registro, login, cambio y recuperación de contraseña) |
| Enfoque      | Aprendizaje guiado: cada línea de código y documentación debe enseñar                    |
| Fecha inicio | Marzo 2026                                                                               |

---

## 2. Stack Tecnológico

### 2.1 Backend (`be/`)

| Herramienta             | Versión | Rol                                        |
| ----------------------- | ------- | ------------------------------------------ |
| Go                      | 1.22+   | Lenguaje principal del backend             |
| Gin                     | 1.9+    | Framework HTTP de alto rendimiento         |
| GORM                    | 2.0+    | ORM para interactuar con la base de datos  |
| golang-migrate          | 4+      | Migraciones de base de datos versionadas   |
| golang-jwt/jwt          | v5      | Creación y verificación de tokens JWT      |
| golang.org/x/crypto     | latest  | Hashing seguro de contraseñas con bcrypt   |
| go-playground/validator | v10     | Validación de structs de request           |
| joho/godotenv           | latest  | Carga de variables de entorno desde `.env` |
| ulule/limiter           | v3      | Rate limiting por IP                       |
| jackc/pgx               | v5      | Driver PostgreSQL para Go                  |
| jordan-wright/email     | latest  | Envío de emails vía SMTP                   |
| stretchr/testify        | latest  | Assertions y mocks para tests              |
| golangci-lint           | latest  | Linter + analizador estático para Go       |

### 2.2 Frontend (`fe/`)

| Herramienta     | Versión | Rol                                          |
| --------------- | ------- | -------------------------------------------- |
| Node.js         | 20 LTS+ | Runtime de JavaScript                        |
| React           | 18+     | Biblioteca para interfaces de usuario        |
| Vite            | 6+      | Bundler y dev server ultrarrápido            |
| TypeScript      | 5.0+    | Superset tipado de JavaScript                |
| TailwindCSS     | 4+      | Framework CSS utility-first                  |
| React Router    | 7+      | Enrutamiento del lado del cliente            |
| Axios           | latest  | Cliente HTTP para comunicación con la API    |
| Vitest          | latest  | Framework de testing compatible con Vite     |
| Testing Library | latest  | Utilidades de testing para componentes React |
| ESLint          | latest  | Linter para TypeScript/React                 |
| Prettier        | latest  | Formateador de código                        |

### 2.3 Base de Datos

| Herramienta    | Versión | Rol                                             |
| -------------- | ------- | ----------------------------------------------- |
| PostgreSQL     | 17+     | Base de datos relacional principal              |
| Docker Compose | latest  | Orquestación de contenedores (BD en desarrollo) |

### 2.4 Autenticación

| Aspecto       | Detalle                                                       |
| ------------- | ------------------------------------------------------------- |
| Método        | JWT (JSON Web Tokens) — stateless                             |
| Access Token  | Duración: 15 minutos                                          |
| Refresh Token | Duración: 7 días                                              |
| Hashing       | bcrypt vía `golang.org/x/crypto/bcrypt`                       |
| Flujos        | Registro, Login, Cambio de contraseña, Recuperación por email |

---

## 3. Reglas de Lenguaje — OBLIGATORIAS

### 3.1 Nomenclatura técnica → INGLÉS

Todo lo que sea código debe estar en inglés:

- Variables, funciones, structs, métodos, interfaces
- Nombres de archivos y carpetas de código
- Endpoints y rutas de la API
- Nombres de tablas y columnas en la base de datos
- Nombres de componentes React y tipos TypeScript
- Mensajes de commits y ramas de git

```go
// ✅ CORRECTO
func GetUserByEmail(email string) (*models.User, error) { ... }

// ❌ INCORRECTO
func ObtenerUsuarioPorEmail(correo string) (*models.User, error) { ... }
```

### 3.2 Comentarios y documentación → ESPAÑOL

Todo lo que sea documentación o comentarios debe estar en español:

- Comentarios en el código (`//`, `/* */`)
- Comentarios de función (godoc style)
- Archivos de documentación (`.md`)
- README.md

### 3.3 Regla del comentario pedagógico — ¿QUÉ? ¿PARA QUÉ? ¿IMPACTO?

Cada comentario significativo debe responder tres preguntas:

```go
// ¿Qué? Función que hashea la contraseña del usuario usando bcrypt.
// ¿Para qué? Almacenar contraseñas de forma segura, nunca en texto plano.
// ¿Impacto? Si se omite el hashing, las contraseñas quedan expuestas ante
//            una filtración de la base de datos.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
```

### 3.4 Cabecera de archivo obligatoria

Cada archivo nuevo debe incluir un comentario de cabecera al inicio:

```go
// Archivo: security.go
// Descripción: Utilidades de seguridad — hashing de contraseñas y manejo de tokens JWT.
// ¿Para qué? Proveer funciones reutilizables de seguridad que se usan en todo el sistema de auth.
// ¿Impacto? Es la base de la seguridad del sistema. Un error aquí compromete toda la autenticación.
```

```typescript
/**
 * Archivo: AuthContext.tsx
 * Descripción: Contexto de React que gestiona el estado de autenticación global.
 * ¿Para qué? Proveer a toda la aplicación acceso al usuario autenticado, tokens y acciones de auth.
 * ¿Impacto? Sin este contexto, no habría forma de saber si el usuario está logueado
 *           ni de proteger rutas que requieren autenticación.
 */
```

---

## 4. Reglas de Entorno y Herramientas — OBLIGATORIAS

### 4.1 Go — SIEMPRE usar Go modules

```bash
# ✅ CORRECTO — Gestión de dependencias con Go modules
cd be
go mod download          # descargar dependencias
go mod tidy              # limpiar dependencias no usadas
go get github.com/gin-gonic/gin  # agregar una dependencia

# ❌ INCORRECTO — No usar GOPATH antiguo ni vendor manual sin go.mod
```

No se necesita `venv` como en Python — Go modules maneja el entorno de forma nativa con `go.mod` y `go.sum`.

### 4.2 Node.js — SIEMPRE usar `pnpm`

```bash
# ✅ CORRECTO
pnpm install
pnpm add axios
pnpm add -D vitest
pnpm dev
pnpm test
pnpm build

# ❌ INCORRECTO — NUNCA usar npm
npm install        # ← PROHIBIDO
npm run dev        # ← PROHIBIDO

# ❌ INCORRECTO — NUNCA usar yarn
yarn install       # ← PROHIBIDO
```

### 4.3 Variables de entorno

- NUNCA hardcodear credenciales, URLs de base de datos, secrets o configuración sensible.
- Usar archivos `.env` (no versionados en git).
- Proveer siempre un `.env.example` con las variables necesarias y valores de ejemplo.
- Validar las variables de entorno al iniciar la aplicación (struct de config en Go).

```bash
# be/.env.example
DATABASE_URL=postgresql://nn_user:nn_password@localhost:5432/nn_auth_db
SECRET_KEY=your-super-secret-key-change-in-production-min-32-chars
ALGORITHM=HS256
ACCESS_TOKEN_EXPIRE_MINUTES=15
REFRESH_TOKEN_EXPIRE_DAYS=7
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
FROM_EMAIL=noreply@nn-company.com
FRONTEND_URL=http://localhost:5173
ENVIRONMENT=development
```

---

## 5. Estructura del Proyecto

```
proyecto/                          # Raíz del monorepo
├── .github/
│   └── copilot-instructions.md    # ← ESTE ARCHIVO — reglas del proyecto
├── .gitignore
├── docker-compose.yml             # PostgreSQL 17 + Mailpit para desarrollo
├── README.md
│
├── _docs/                         # 📚 Documentación del proyecto
│   ├── referencia-tecnica/
│   │   ├── architecture.md
│   │   ├── api-endpoints.md
│   │   └── database-schema.md
│   ├── conceptos/
│   │   ├── owasp-top-10.md
│   │   ├── accesibilidad-aria-wcag.md
│   │   └── patrones-arquitectonicos.md
│   ├── requisitos/
│   │   ├── HUs/                   # Historias de Usuario
│   │   ├── RFs/                   # Requisitos Funcionales
│   │   ├── RNFs/                  # Requisitos No Funcionales
│   │   └── restricciones.md
│   ├── setup/
│   │   ├── con-docker.md
│   │   └── sin-docker.md
│   └── guion-video.md
│
├── _assets/                       # 🖼️ Recursos estáticos (diagramas)
│
├── be/                            # 🐹 Backend — Go + Gin + GORM
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go            # Punto de entrada — configura y arranca Gin
│   │   └── migrate/
│   │       └── main.go            # Herramienta CLI de migraciones
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go          # Configuración (struct + godotenv)
│   │   ├── database/
│   │   │   └── database.go        # GORM engine + pool de conexiones
│   │   ├── middleware/
│   │   │   ├── auth.go            # Middleware JWT: extrae y valida el token
│   │   │   ├── cors.go            # Middleware CORS
│   │   │   ├── ratelimit.go       # Rate limiting por IP
│   │   │   └── security.go        # Cabeceras de seguridad HTTP
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── password_reset_token.go
│   │   │   └── email_verification_token.go
│   │   ├── dto/
│   │   │   └── auth.go            # RegisterRequest, LoginRequest, TokenResponse, etc.
│   │   ├── handlers/
│   │   │   ├── auth.go            # Register, Login, Refresh, ChangePassword, etc.
│   │   │   └── user.go            # GetMe
│   │   ├── services/
│   │   │   └── auth_service.go    # Lógica de negocio
│   │   └── utils/
│   │       ├── security.go        # HashPassword, VerifyPassword, CreateToken, etc.
│   │       ├── email.go           # SendVerificationEmail, SendPasswordResetEmail
│   │       └── audit_log.go       # LogLoginSuccess, LogLoginFailed, etc.
│   ├── migrations/
│   │   ├── 000001_create_users.up.sql
│   │   ├── 000001_create_users.down.sql
│   │   ├── 000002_create_password_reset_tokens.up.sql
│   │   ├── 000002_create_password_reset_tokens.down.sql
│   │   ├── 000003_create_email_verification_tokens.up.sql
│   │   └── 000003_create_email_verification_tokens.down.sql
│   ├── go.mod
│   ├── go.sum
│   ├── .env.example
│   └── .env                       # Variables de entorno (NO versionado)
│
└── fe/                            # ⚛️ Frontend — React + Vite + TypeScript
    ├── .env
    ├── .env.example
    ├── index.html
    ├── package.json
    ├── pnpm-lock.yaml
    ├── vite.config.ts
    ├── tsconfig.json
    ├── eslint.config.js
    └── src/
        ├── main.tsx
        ├── App.tsx
        ├── index.css
        ├── api/
        │   ├── auth.ts
        │   └── axios.ts
        ├── components/
        │   ├── ui/
        │   └── layout/
        ├── pages/
        │   ├── LandingPage.tsx
        │   ├── LoginPage.tsx
        │   ├── RegisterPage.tsx
        │   ├── DashboardPage.tsx
        │   ├── ChangePasswordPage.tsx
        │   ├── ForgotPasswordPage.tsx
        │   └── ResetPasswordPage.tsx
        ├── hooks/
        │   └── useAuth.ts
        ├── context/
        │   └── AuthContext.tsx
        ├── types/
        │   └── auth.ts
        └── __tests__/
            └── auth.test.tsx
```

---

## 6. Convenciones de Código

### 6.1 Go (Backend)

| Aspecto              | Convención                                      |
| -------------------- | ----------------------------------------------- |
| Estilo               | `gofmt` — formato estándar de Go (obligatorio)  |
| Naming variables     | `camelCase`                                     |
| Naming tipos/structs | `PascalCase`                                    |
| Naming constantes    | `UpperCamelCase` o `UPPER_SNAKE_CASE`           |
| Naming archivos      | `snake_case.go`                                 |
| Errores              | Siempre manejar errores explícitamente (no `_`) |
| Comentarios          | Godoc style en español (funciones exportadas)   |
| Imports              | `goimports` — stdlib → external → internal      |
| Línea máxima         | 100 caracteres                                  |

```go
// ✅ Ejemplo de función bien documentada en Go
// HashPassword convierte la contraseña en texto plano a un hash bcrypt.
// ¿Para qué? Almacenar contraseñas de forma segura — nunca en texto plano.
// ¿Impacto? Sin este hashing, una filtración de BD expondría todas las contraseñas.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("hashing password: %w", err)
    }
    return string(bytes), nil
}
```

### 6.2 TypeScript/React (Frontend)

| Aspecto            | Convención                                                           |
| ------------------ | -------------------------------------------------------------------- |
| Estilo             | ESLint + Prettier                                                    |
| Naming variables   | `camelCase`                                                          |
| Naming componentes | `PascalCase`                                                         |
| Naming archivos    | `PascalCase` para componentes, `camelCase` para utilidades           |
| Naming tipos       | `PascalCase` con sufijo descriptivo (`UserResponse`, `LoginRequest`) |
| Componentes        | Funcionales con hooks — nunca clases                                 |
| CSS                | TailwindCSS utility classes — evitar CSS custom                      |
| Strict mode        | `"strict": true` en `tsconfig.json`                                  |

### 6.3 SQL / Base de Datos

| Aspecto      | Convención                 | Ejemplo                          |
| ------------ | -------------------------- | -------------------------------- |
| Tablas       | `snake_case`, plural       | `users`, `password_reset_tokens` |
| Columnas     | `snake_case`               | `hashed_password`, `created_at`  |
| Primary Keys | `id` UUID                  | `id UUID PK`                     |
| Foreign Keys | `<tabla_singular>_id`      | `user_id → users.id`             |
| Timestamps   | `created_at`, `updated_at` | Todas las tablas                 |
| Migraciones  | Siempre vía golang-migrate | Nunca alterar BD manualmente     |

---

## 7. Conventional Commits — OBLIGATORIO

### 7.1 Formato

```
type(scope): short description in english

What: Detailed description of what was done
For: Why this change is needed
Impact: What effect this has on the system
```

### 7.2 Tipos permitidos

| Tipo       | Cuándo usarlo                              |
| ---------- | ------------------------------------------ |
| `feat`     | Nueva funcionalidad                        |
| `fix`      | Corrección de bug                          |
| `docs`     | Solo documentación                         |
| `style`    | Formato, espacios (no afecta lógica)       |
| `refactor` | Reestructuración sin cambiar funcionalidad |
| `test`     | Agregar o corregir tests                   |
| `chore`    | Mantenimiento, configuración, dependencias |
| `ci`       | Cambios en CI/CD                           |
| `perf`     | Mejoras de rendimiento                     |

### 7.3 Scopes sugeridos

- `auth` — Autenticación y autorización
- `user` — Modelo/funcionalidad de usuario
- `db` — Base de datos y migraciones
- `api` — Handlers y rutas
- `ui` — Componentes y estilos del frontend
- `config` — Configuración y entorno
- `test` — Tests
- `deps` — Dependencias

### 7.4 Ejemplos

```bash
# ✅ Ejemplo de commit completo
git commit -m "feat(auth): add user registration handler

What: Creates POST /api/v1/auth/register endpoint with email validation,
password hashing with bcrypt, and duplicate email check
For: Allow new users to create accounts in the NN Auth System
Impact: Enables the user onboarding flow; stores hashed passwords in the users table"

# ✅ Ejemplo de fix
git commit -m "fix(auth): handle expired refresh token gracefully

What: Returns 401 with clear error message when refresh token is expired
For: Prevent confusing 500 errors when users try to refresh after 7 days
Impact: Improves UX by redirecting to login instead of showing error page"
```

---

## 8. Calidad — NO es Opcional, es OBLIGACIÓN

### 8.1 Principio fundamental

> Código que se genera, código que se prueba.

### 8.2 Testing — Backend (Go)

| Herramienta         | Propósito                                          |
| ------------------- | -------------------------------------------------- |
| `testing` (stdlib)  | Framework principal de testing de Go               |
| `net/http/httptest` | Servidor HTTP en memoria para tests de integración |
| `testify/assert`    | Assertions expresivas                              |
| `testify/mock`      | Mocks de interfaces                                |

```bash
# Ejecutar todos los tests del backend
cd be && go test ./... -v

# Ejecutar con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Ejecutar un paquete específico
go test ./internal/services/... -v -run TestRegisterUser
```

Cobertura mínima esperada: 80% en paquetes de lógica de negocio.

### 8.3 Testing — Frontend

| Herramienta            | Propósito                       |
| ---------------------- | ------------------------------- |
| Vitest                 | Test runner compatible con Vite |
| @testing-library/react | Testing de componentes React    |
| jsdom                  | Simular el DOM en Node.js       |

```bash
cd fe && pnpm test         # todos los tests
pnpm test:watch             # modo watch
pnpm test:coverage          # con cobertura
```

### 8.4 Linting y Formateo

```bash
# Backend — golangci-lint
cd be && golangci-lint run ./...    # verificar errores
cd be && gofmt -w .                 # formatear código

# Frontend — ESLint + Prettier
cd fe && pnpm lint                  # verificar errores
cd fe && pnpm format                # formatear código
```

### 8.5 Checklist antes de commit

- [ ] ¿El código de Go compila sin errores? (`go build ./...`)
- [ ] ¿Hay comentarios pedagógicos (¿Qué? ¿Para qué? ¿Impacto?)?
- [ ] ¿Los tests pasan? (`go test ./...` / `pnpm test`)
- [ ] ¿El linter no reporta errores? (`golangci-lint run` / `pnpm lint`)
- [ ] ¿El commit sigue Conventional Commits con What/For/Impact?
- [ ] ¿Las variables sensibles están en `.env` y no hardcodeadas?
- [ ] ¿El `.env.example` se actualizó si se agregaron nuevas variables?
- [ ] ¿Los errores de Go se manejan explícitamente (no ignorados con `_`)?

---

## 9. Seguridad — Mejores Prácticas

### 9.1 Contraseñas

- SIEMPRE hashear con bcrypt (`golang.org/x/crypto/bcrypt`) antes de almacenar
- NUNCA almacenar contraseñas en texto plano
- NUNCA loggear contraseñas ni incluirlas en responses
- Validar fortaleza mínima: ≥8 caracteres, al menos 1 mayúscula, 1 minúscula, 1 número

### 9.2 JWT (Tokens)

- Access Token: corta duración (15 min) — se envía en header `Authorization: Bearer <token>`
- Refresh Token: larga duración (7 días)
- Secret key: mínimo 32 caracteres, aleatoria, en variable de entorno
- Algoritmo: HS256 (`golang-jwt/jwt/v5`)
- NUNCA almacenar tokens en `localStorage` en producción

### 9.3 CORS

- Configurar orígenes permitidos explícitamente en el middleware de Gin
- En desarrollo: permitir `http://localhost:5173`
- En producción: NUNCA usar `AllowAllOrigins: true`

### 9.4 API

- Versionamiento: `/api/v1/...`
- Rate limiting en endpoints de auth (ulule/limiter)
- Validación de inputs con `go-playground/validator` (nunca confiar en datos del cliente)
- Mensajes de error genéricos en auth (no revelar si el email existe)

### 9.5 Base de datos

- Usar siempre GORM (nunca raw SQL sin parametrizar)
- Conexiones con pool configurado (`db.DB().SetMaxOpenConns()`)
- Credenciales en variables de entorno

---

## 10. Estructura de la API

### 10.1 Prefijo base

Todos los endpoints van bajo `/api/v1/`

### 10.2 Endpoints de autenticación (`/api/v1/auth/`)

| Método | Ruta               | Auth | Descripción                           |
| ------ | ------------------ | ---- | ------------------------------------- |
| POST   | `/register`        | No   | Registrar nuevo usuario               |
| POST   | `/login`           | No   | Iniciar sesión, obtener tokens        |
| POST   | `/refresh`         | No † | Renovar access token con refresh      |
| POST   | `/change-password` | Sí   | Cambiar contraseña (usuario logueado) |
| POST   | `/forgot-password` | No   | Solicitar email de recuperación       |
| POST   | `/reset-password`  | No † | Restablecer contraseña con token      |
| POST   | `/verify-email`    | No † | Verificar dirección de email          |

### 10.3 Endpoints de usuario (`/api/v1/users/`)

| Método | Ruta  | Auth | Descripción                       |
| ------ | ----- | ---- | --------------------------------- |
| GET    | `/me` | Sí   | Obtener perfil del usuario actual |

---

## 11. Esquema de Base de Datos

### 11.1 Tabla `users`

| Columna             | Tipo         | Restricciones                 |
| ------------------- | ------------ | ----------------------------- |
| `id`                | UUID         | PK, default gen_random_uuid() |
| `email`             | VARCHAR(255) | UNIQUE, NOT NULL, INDEXED     |
| `full_name`         | VARCHAR(255) | NOT NULL                      |
| `hashed_password`   | VARCHAR(255) | NOT NULL                      |
| `is_active`         | BOOLEAN      | DEFAULT TRUE                  |
| `is_email_verified` | BOOLEAN      | DEFAULT FALSE                 |
| `created_at`        | TIMESTAMPTZ  | DEFAULT NOW(), NOT NULL       |
| `updated_at`        | TIMESTAMPTZ  | DEFAULT NOW(), NOT NULL       |

### 11.2 Tabla `password_reset_tokens`

| Columna      | Tipo         | Restricciones             |
| ------------ | ------------ | ------------------------- |
| `id`         | UUID         | PK                        |
| `user_id`    | UUID         | FK → users.id, CASCADE    |
| `token`      | VARCHAR(255) | UNIQUE, NOT NULL, INDEXED |
| `expires_at` | TIMESTAMPTZ  | NOT NULL                  |
| `used`       | BOOLEAN      | DEFAULT FALSE             |
| `created_at` | TIMESTAMPTZ  | DEFAULT NOW(), NOT NULL   |

### 11.3 Tabla `email_verification_tokens`

Misma estructura que `password_reset_tokens`.

---

## 12. Flujos de Autenticación

### 12.1 Registro

```
Cliente → POST /api/v1/auth/register { email, full_name, password }
  → Validar datos (go-playground/validator)
  → Verificar email no duplicado (GORM query)
  → Hashear password (bcrypt)
  → Crear usuario en BD (GORM Create)
  → Crear EmailVerificationToken (24h expiración)
  → Enviar email con enlace de verificación
  → Retornar usuario creado (sin password, is_email_verified: false)
```

### 12.2 Login

```
Cliente → POST /api/v1/auth/login { email, password }
  → Buscar usuario por email (GORM First)
  → Verificar password contra hash (bcrypt.CompareHashAndPassword)
  → Verificar is_email_verified == true
  → Verificar is_active == true
  → Generar access_token (15 min) + refresh_token (7 días)
  → Registrar evento en audit log
  → Retornar { access_token, refresh_token, token_type: "bearer" }
```

### 12.3 Cambio de contraseña

```
Cliente → POST /api/v1/auth/change-password { current_password, new_password }
  → (Requiere Authorization: Bearer <access_token>)
  → Middleware extrae user_id del JWT
  → Verificar current_password contra hash
  → Hashear new_password con bcrypt
  → Actualizar en BD (GORM Save)
  → Registrar evento en audit log
  → Retornar confirmación
```

### 12.4 Recuperación de contraseña

```
Paso 1: POST /api/v1/auth/forgot-password { email }
  → Buscar usuario por email
  → Generar token UUID + expiración 1h
  → Guardar en password_reset_tokens
  → Enviar email con enlace
  → Retornar mensaje genérico (siempre igual)

Paso 2: POST /api/v1/auth/reset-password { token, new_password }
  → Buscar token en BD
  → Verificar expires_at > NOW()
  → Verificar used == false
  → Hashear new_password
  → Actualizar password del usuario
  → Marcar token used = true
  → Retornar confirmación
```

---

## 13. Configuración de Docker Compose

```yaml
# Solo para desarrollo local — PostgreSQL 17 + Mailpit
services:
  db:
    image: postgres:17-alpine
    container_name: nn_auth_db
    environment:
      POSTGRES_USER: nn_user
      POSTGRES_PASSWORD: nn_password
      POSTGRES_DB: nn_auth_db
    ports:
      - "5432:5432"
    volumes:
      - nn_auth_data:/var/lib/postgresql/data

  mailpit:
    image: axllent/mailpit
    container_name: nn_auth_mailpit
    ports:
      - "8025:8025" # Web UI
      - "1025:1025" # SMTP

volumes:
  nn_auth_data:
```

---

## 14. Mejores Prácticas — Resumen

### 14.1 Generales

- ✅ DRY (Don't Repeat Yourself) — reutilizar código
- ✅ KISS (Keep It Simple, Stupid) — preferir soluciones simples
- ✅ YAGNI (You Aren't Gonna Need It) — no agregar lo que no se necesita aún
- ✅ Separation of Concerns — cada paquete tiene una responsabilidad clara
- ✅ Fail fast — manejar errores al inicio de cada función

### 14.2 Backend (Go)

- ✅ Usar interfaces para desacoplar componentes (facilita testing con mocks)
- ✅ Retornar desde funciones con error como segundo valor (`value, err`)
- ✅ Usar `context.Context` en toda función que haga I/O
- ✅ Separar handlers (HTTP) de services (lógica de negocio)
- ✅ Usar GORM con DTOs separados de los modelos de BD
- ✅ Documentar toda función exportada con godoc

### 14.3 Frontend

- ✅ Componentes pequeños y reutilizables
- ✅ Estado global solo cuando es necesario (Context API para auth)
- ✅ Custom hooks para encapsular lógica reutilizable
- ✅ Rutas protegidas con componente `ProtectedRoute`
- ✅ Loading states para operaciones asíncronas

### 14.4 Diseño y UX/UI — OBLIGATORIO

| Aspecto           | Regla                                                           |
| ----------------- | --------------------------------------------------------------- |
| Temas             | Dark mode y Light mode con toggle — usar `prefers-color-scheme` |
| Tipografía        | Fuentes sans-serif exclusivamente (Inter, system-ui)            |
| Colores           | Sólidos y planos — SIN degradados (gradient) en ningún lugar    |
| Estilo visual     | Diseño moderno, limpio, minimalista con excelente UX/UI         |
| Botones de acción | Siempre alineados a la derecha (`justify-end`)                  |
| Spacing           | Escala consistente de Tailwind (p-4, gap-6, space-y-4)          |
| Responsividad     | Mobile-first — los formularios deben verse bien en móvil        |
| Accesibilidad     | Labels en inputs, aria-\* básicos, contraste WCAG AA            |

---

## 15. Reglas para Copilot / IA — Al Generar Código

1. **Dividir respuestas largas** — Si la implementación es extensa, dividirla en pasos incrementales.
2. **Código generado = código probado** — Siempre incluir o sugerir tests.
3. **Comentarios pedagógicos** — Cada bloque significativo con ¿Qué? ¿Para qué? ¿Impacto?
4. **Tipos obligatorios** — Nunca omitir tipos en Go ni en TypeScript.
5. **Formato correcto** — Respetar `gofmt` para Go y Prettier/ESLint para TypeScript.
6. **Usar las herramientas correctas** — Go modules para Go, `pnpm` para Node.js. Sin excepciones.
7. **Variables de entorno** — Toda configuración sensible en `.env`, nunca hardcodeada.
8. **Conventional Commits** — Sugerir mensajes de commit con formato correcto.
9. **Seguridad primero** — Nunca almacenar passwords en texto plano, nunca exponer secrets.
10. **Manejo de errores** — En Go, siempre manejar `error` explícitamente. Nunca usar `_` para ignorar errores.

---

## 16. Plan de Trabajo — Fases

> Cada fase es independiente y verificable. No avanzar a la siguiente sin completar y probar la actual.

### Fase 0 — Fundamentos y Configuración Base

- [ ] Crear `.github/copilot-instructions.md` (este archivo)
- [ ] Crear `.gitignore` raíz
- [ ] Crear `docker-compose.yml` con PostgreSQL 17 + Mailpit
- [ ] Crear `README.md`

### Fase 1 — Backend Setup (Go)

- [ ] Inicializar módulo Go: `go mod init nn-auth-system`
- [ ] Agregar dependencias en `go.mod`
- [ ] Crear `internal/config/config.go` — struct de configuración + godotenv
- [ ] Crear `internal/database/database.go` — conexión GORM + pool
- [ ] Crear `cmd/api/main.go` — app Gin con CORS y rutas base
- [ ] Crear `.env.example` y `.env`
- [ ] ✅ Verificar: `go run ./cmd/api/main.go` → servidor escucha en `:8000`

### Fase 2 — Modelos y Migraciones

- [ ] Crear `internal/models/user.go` — struct GORM User
- [ ] Crear `internal/models/password_reset_token.go`
- [ ] Crear `internal/models/email_verification_token.go`
- [ ] Crear scripts SQL en `migrations/`
- [ ] Crear `cmd/migrate/main.go` — herramienta de migración
- [ ] Ejecutar migraciones: `go run ./cmd/migrate/main.go up`
- [ ] ✅ Verificar: tablas creadas en PostgreSQL

### Fase 3 — Autenticación Backend

- [ ] Crear `internal/utils/security.go` — HashPassword, VerifyPassword, CreateToken, etc.
- [ ] Crear `internal/dto/auth.go` — structs de request/response con validación
- [ ] Crear `internal/services/auth_service.go` — lógica de negocio
- [ ] Crear `internal/utils/email.go` — envío de emails
- [ ] Crear `internal/middleware/auth.go` — middleware JWT
- [ ] Crear `internal/middleware/ratelimit.go` — rate limiting
- [ ] Crear `internal/middleware/security.go` — cabeceras de seguridad
- [ ] Crear `internal/handlers/auth.go` — handlers HTTP
- [ ] Crear `internal/handlers/user.go` — handler GET /me
- [ ] Registrar rutas en `cmd/api/main.go`
- [ ] ✅ Verificar: probar todos los endpoints con `curl` o Swagger

### Fase 4 — Tests Backend

- [ ] Crear tests en `internal/services/auth_service_test.go`
- [ ] Crear tests en `internal/handlers/auth_test.go`
- [ ] ✅ Verificar: `go test ./... -v` → todos los tests pasan

### Fase 5 — Frontend Setup

- [ ] Inicializar proyecto Vite con React + TypeScript en `fe/`
- [ ] Instalar dependencias con `pnpm`
- [ ] Configurar TailwindCSS
- [ ] Crear `.env.example`
- [ ] ✅ Verificar: `pnpm dev` → app visible en `http://localhost:5173`

### Fase 6 — Frontend Auth

- [ ] Crear tipos TypeScript (`types/auth.ts`)
- [ ] Crear cliente HTTP (`api/auth.ts`, `api/axios.ts`)
- [ ] Crear AuthContext + Provider
- [ ] Crear hook `useAuth`
- [ ] Crear componentes UI (InputField, Button, Alert)
- [ ] Crear ProtectedRoute
- [ ] Crear páginas: Landing, Login, Register, Dashboard, ChangePassword, ForgotPassword, ResetPassword
- [ ] Crear páginas legales y formulario de contacto
- [ ] ✅ Verificar: flujo completo funciona contra la API

### Fase 7 — Tests Frontend

- [ ] Configurar Vitest + Testing Library
- [ ] Crear tests para componentes y flujos de auth
- [ ] ✅ Verificar: `pnpm test` → todos los tests pasan

### Fase 8 — Documentación Final

- [ ] Crear/actualizar `_docs/referencia-tecnica/architecture.md`
- [ ] Crear/actualizar `_docs/referencia-tecnica/api-endpoints.md`
- [ ] Crear/actualizar `_docs/referencia-tecnica/database-schema.md`
- [ ] Actualizar `README.md` con instrucciones finales
- [ ] ✅ Verificar: documentación completa y coherente

---

## 17. Verificación Final del Sistema

```bash
# 1. Levantar base de datos
docker compose up -d

# 2. Levantar backend
cd be && go run ./cmd/api/main.go

# 3. Levantar frontend (en otra terminal)
cd fe && pnpm dev

# 4. Ejecutar tests backend
cd be && go test ./... -v -cover

# 5. Ejecutar tests frontend
cd fe && pnpm test

# 6. Flujo manual completo:
#    Registro → Verificar email → Login → Ver perfil → Cambiar contraseña →
#    Logout → Forgot password → Reset password → Login con nueva contraseña
```

> Recuerda: La calidad no es una opción, es una obligación. Cada línea de código es una oportunidad de aprender y enseñar.
