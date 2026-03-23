# Arquitectura del Sistema — NN Auth System

<!--
  ¿Qué? Documento que describe la arquitectura general del sistema.
  ¿Para qué? Que cualquier desarrollador entienda cómo están organizadas las capas,
             cómo fluyen las peticiones y por qué se tomaron las decisiones técnicas.
  ¿Impacto? Sin esta guía, cada desarrollador implementaría su propia estructura,
             generando un sistema inconsistente e imposible de mantener.
-->

---

## 1. Vista General

El sistema sigue una **arquitectura de tres capas** clásica (3-Tier Architecture):

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTE                                  │
│                   React + TypeScript + Vite                      │
│              (Navegador Web — Puerto 5173 en dev)                │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTP/HTTPS (JSON)
                             │ Authorization: Bearer <token>
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                       SERVIDOR (API)                             │
│                    Go 1.22+ + Gin v1.9+                         │
│              (Aplicación REST — Puerto 8000 en dev)              │
│                                                                  │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │
│   │Middleware│→ │ Handlers │→ │ Services │→ │  GORM / Utils│  │
│   │(CORS,JWT,│  │(auth.go, │  │(auth_    │  │(security.go, │  │
│   │ limiter) │  │ user.go) │  │service)  │  │ email.go)    │  │
│   └──────────┘  └──────────┘  └──────────┘  └──────────────┘  │
└────────────────────────────┬────────────────────────────────────┘
                             │ GORM (ORM)
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      BASE DE DATOS                               │
│                     PostgreSQL 17+                               │
│              (Puerto 5432 — Docker en desarrollo)                │
│                                                                  │
│  ┌──────────────┐  ┌───────────────────────┐  ┌─────────────┐  │
│  │    users      │  │ password_reset_tokens │  │email_verif. │  │
│  │  (usuarios)   │  │  (tokens de reset)    │  │  _tokens    │  │
│  └──────────────┘  └───────────────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Arquitectura del Backend (Go)

### 2.1 Estructura de capas

```
be/
├── cmd/
│   ├── api/
│   │   └── main.go          ← Punto de entrada: configura e inicia el servidor Gin
│   └── migrate/
│       └── main.go          ← Herramienta CLI: corre las migraciones SQL
│
├── internal/                ← Código privado (no importable desde fuera del módulo)
│   ├── config/
│   │   └── config.go        ← Carga variables de entorno y las mapea a un struct tipado
│   │
│   ├── database/
│   │   └── database.go      ← Inicializa la conexión GORM con pool de conexiones
│   │
│   ├── middleware/          ← Lógica transversal (ejecutada en CADA petición)
│   │   ├── auth.go          ← Valida el JWT y extrae el user_id al contexto de Gin
│   │   ├── cors.go          ← Configura los orígenes permitidos (CORS)
│   │   ├── ratelimit.go     ← Limita peticiones por IP (ulule/limiter)
│   │   └── security.go      ← Añade cabeceras HTTP de seguridad (X-Frame-Options, etc.)
│   │
│   ├── models/              ← Structs que representan las tablas de la BD (GORM)
│   │   ├── user.go
│   │   ├── password_reset_token.go
│   │   └── email_verification_token.go
│   │
│   ├── dto/                 ← Data Transfer Objects: structs para requests y responses
│   │   └── auth.go          ← RegisterRequest, LoginRequest, TokenResponse, etc.
│   │
│   ├── handlers/            ← Capa HTTP: recibe, valida y responde peticiones Gin
│   │   ├── auth.go          ← Register, Login, Refresh, ChangePassword, etc.
│   │   └── user.go          ← GetMe
│   │
│   ├── services/            ← Lógica de negocio (sin conocimiento del protocolo HTTP)
│   │   └── auth_service.go  ← RegisterUser, LoginUser, ResetPassword, etc.
│   │
│   └── utils/               ← Utilitarios independientes y reutilizables
│       ├── security.go      ← HashPassword, VerifyPassword, CreateToken, ParseToken
│       ├── email.go         ← SendVerificationEmail, SendPasswordResetEmail (SMTP)
│       └── audit_log.go     ← LogLoginSuccess, LogLoginFailed, LogPasswordChange
│
└── migrations/
    ├── 000001_create_users.up.sql
    ├── 000001_create_users.down.sql
    ├── 000002_create_password_reset_tokens.up.sql
    ├── 000002_create_password_reset_tokens.down.sql
    ├── 000003_create_email_verification_tokens.up.sql
    └── 000003_create_email_verification_tokens.down.sql
```

### 2.2 Flujo de una petición HTTP

Cada request al backend atraviesa las siguientes capas en orden:

```
Cliente HTTP
    │
    ▼
┌──────────────────────────────────────────────────────────────┐
│  Gin Router — Distribución de rutas                          │
│  cmd/api/main.go → r.POST("/api/v1/auth/login", handler)    │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Middleware Global (se ejecuta en TODAS las rutas)           │
│  1. CORS (cors.go)         → Verifica Origin permitido       │
│  2. Security (security.go) → Añade cabeceras de seguridad    │
│  3. RateLimit (ratelimit.go)→ Verifica límite por IP         │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Middleware de Ruta (solo en rutas protegidas)               │
│  4. Auth (auth.go) → Valida JWT, inyecta userID al contexto │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Handler (handlers/auth.go, handlers/user.go)                │
│  - Bind y validación del JSON body (go-playground/validator) │
│  - Llama al Service correspondiente                          │
│  - Responde con c.JSON(statusCode, responseBody)             │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Service (services/auth_service.go)                          │
│  - Lógica de negocio pura (sin conocer HTTP ni Gin)          │
│  - Usa utils/security.go para bcrypt y JWT                   │
│  - Usa utils/email.go para envío de emails                   │
│  - Usa modelos GORM para acceder a la BD                     │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  GORM + PostgreSQL (database/database.go)                    │
│  - Traduce structs Go a queries SQL                          │
│  - Maneja el pool de conexiones                              │
│  - Retorna errores tipados para fácil manejo                 │
└──────────────────────────────────────────────────────────────┘
```

### 2.3 Separación de responsabilidades

| Capa          | Responsabilidad                                     | Lo que NO debe hacer                      |
| ------------- | --------------------------------------------------- | ----------------------------------------- |
| `middleware/` | Validar, enriquecer o rechazar el request           | Lógica de negocio, acceso directo a BD    |
| `handlers/`   | Parsear HTTP request, llamar service, JSON response | Lógica de negocio, queries de BD directas |
| `services/`   | Orquestar la lógica de negocio                      | Saber que existe HTTP, Gin, o JSON        |
| `models/`     | Representar las tablas de la BD                     | Lógica de negocio                         |
| `dto/`        | Definir la forma de requests y responses            | Acceso a BD                               |
| `utils/`      | Operaciones atómicas reutilizables                  | Lógica de negocio específica del dominio  |

---

## 3. Arquitectura del Frontend (React)

### 3.1 Estructura de carpetas

```
fe/src/
├── main.tsx                 ← Punto de entrada: monta App en el DOM
├── App.tsx                  ← Define las rutas (React Router) y envuelve con AuthProvider
├── index.css                ← Estilos globales (TailwindCSS directives)
│
├── api/
│   ├── axios.ts             ← Instancia de Axios con baseURL, interceptors y refresh token
│   └── auth.ts              ← Funciones: loginUser(), registerUser(), getMeUser(), etc.
│
├── context/
│   └── AuthContext.tsx      ← Estado global de auth: user, tokens, isLoading
│
├── hooks/
│   └── useAuth.ts           ← Custom hook que consume AuthContext
│
├── components/
│   ├── ui/                  ← Componentes reutilizables: InputField, Button, Alert
│   └── layout/              ← Header, Footer, ProtectedRoute, Layout
│
├── pages/
│   ├── LandingPage.tsx
│   ├── LoginPage.tsx
│   ├── RegisterPage.tsx
│   ├── DashboardPage.tsx    ← Ruta protegida (requiere autenticación)
│   ├── ChangePasswordPage.tsx
│   ├── ForgotPasswordPage.tsx
│   └── ResetPasswordPage.tsx
│
└── types/
    └── auth.ts              ← Interfaces TypeScript: User, TokenResponse, LoginRequest, etc.
```

### 3.2 Flujo de estado de autenticación

```
┌─────────────────────────────────────────────────────────────────┐
│  AuthContext (context/AuthContext.tsx)                           │
│  Estado: { user, accessToken, refreshToken, isLoading }         │
│  Acciones: { login(), logout(), register() }                    │
└────────────────┬────────────────────────────────────────────────┘
                 │  Provee a través de <AuthProvider>
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  useAuth() hook (hooks/useAuth.ts)                              │
│  Simplifica el acceso al context desde cualquier componente     │
└────────┬──────────────────────────────────────────────────────┘
         │  Usado en
         ▼
┌──────────────────────────────────────────────────────────────────┐
│  Páginas y componentes                                           │
│  LoginPage → useAuth().login(credentials)                       │
│  ProtectedRoute → useAuth().user ? <Outlet> : <Navigate /login> │
│  DashboardPage → useAuth().user.fullName                        │
└──────────────────────────────────────────────────────────────────┘
         │  Llaman a
         ▼
┌──────────────────────────────────────────────────────────────────┐
│  Capa API (api/auth.ts + api/axios.ts)                          │
│  - Axios interceptor: adjunta access token a cada request       │
│  - Axios interceptor: si 401, intenta refresh automático        │
│  - Si refresh falla → logout() y redirect a /login              │
└──────────────────────────────────────────────────────────────────┘
```

---

## 4. Flujos de Autenticación (End-to-End)

### 4.1 Registro de usuario

```
Usuario               Frontend                 Backend (Go/Gin)           BD (PostgreSQL)
  │                      │                           │                          │
  │ Llena el form        │                           │                          │
  │─────────────────────►│                           │                          │
  │                      │ POST /api/v1/auth/register│                          │
  │                      │{email, full_name, password}                          │
  │                      │──────────────────────────►│                          │
  │                      │                           │ Valida (validator)        │
  │                      │                           │──────────────────────────►
  │                      │                           │ Busca email duplicado    │
  │                      │                           │◄──────────────────────────
  │                      │                           │ Hashea password (bcrypt) │
  │                      │                           │ Crea usuario GORM        │
  │                      │                           │──────────────────────────►
  │                      │                           │ Crea token verificación   │
  │                      │                           │──────────────────────────►
  │                      │                           │ Envía email (SMTP)       │
  │                      │ 201 Created {user}        │                          │
  │                      │◄──────────────────────────│                          │
  │ Mensaje de éxito     │                           │                          │
  │◄─────────────────────│                           │                          │
```

### 4.2 Login

```
Usuario               Frontend                 Backend (Go/Gin)           BD (PostgreSQL)
  │                      │                           │                          │
  │ Introduce credenciales                           │                          │
  │─────────────────────►│                           │                          │
  │                      │ POST /api/v1/auth/login   │                          │
  │                      │{email, password}          │                          │
  │                      │──────────────────────────►│                          │
  │                      │                           │ Busca user por email     │
  │                      │                           │──────────────────────────►
  │                      │                           │◄──────────────────────────
  │                      │                           │ bcrypt.CompareHashAndPassword
  │                      │                           │ Verifica is_active, is_email_verified
  │                      │                           │ Genera access_token (15 min)
  │                      │                           │ Genera refresh_token (7 días)
  │                      │                           │ Registra en audit log    │
  │                      │ 200 OK                    │                          │
  │                      │ {access_token,            │                          │
  │                      │  refresh_token}           │                          │
  │                      │◄──────────────────────────│                          │
  │                      │ Guarda tokens en memoria  │                          │
  │ Redirige a /dashboard│                           │                          │
  │◄─────────────────────│                           │                          │
```

### 4.3 Renovación automática de token (Refresh)

```
Frontend (Axios interceptor)          Backend (Go/Gin)
  │                                        │
  │ Request con access_token expirado      │
  │───────────────────────────────────────►│
  │                                        │ 401 Unauthorized
  │◄───────────────────────────────────────│
  │                                        │
  │ POST /api/v1/auth/refresh              │
  │ {refresh_token}                        │
  │───────────────────────────────────────►│
  │                                        │ Valida refresh_token (golang-jwt)
  │                                        │ Genera nuevo access_token
  │ 200 OK {new_access_token}              │
  │◄───────────────────────────────────────│
  │                                        │
  │ Reintenta el request original          │
  │ con el nuevo access_token              │
  │───────────────────────────────────────►│
```

---

## 5. Capas de Seguridad

```
Internet
    │
    ▼
┌─────────────────────────────────────┐
│ HTTPS / TLS (producción)            │ ← Cifrado en tránsito
│ (Nginx / Reverse proxy)             │
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ Rate Limiting (ulule/limiter)        │ ← Max N req/min por IP (A04 OWASP)
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ CORS (middleware/cors.go)           │ ← Solo orígenes permitidos (A05 OWASP)
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ Security Headers (middleware/       │ ← X-Frame-Options, X-XSS-Protection,
│ security.go)                        │   Content-Security-Policy (A05 OWASP)
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ JWT Verification (middleware/       │ ← Verifica firma y expiración (A07 OWASP)
│ auth.go) — solo rutas protegidas    │   golang-jwt/jwt v5
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ Input Validation (go-playground/    │ ← Sanitiza y valida todos los inputs
│ validator v10 en DTOs)              │   (A03, A07 OWASP)
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ GORM ORM (queries parametrizadas)  │ ← Previene SQL injection (A03 OWASP)
└──────────────────┬──────────────────┘
                   │
                   ▼
┌─────────────────────────────────────┐
│ bcrypt Hash (golang.org/x/crypto)  │ ← Contraseñas nunca en texto plano (A02)
└─────────────────────────────────────┘
```

---

## 6. Decisiones Técnicas

### ¿Por qué Go + Gin?

| Alternativa       | Por qué no se eligió                                       |
| ----------------- | ---------------------------------------------------------- |
| **Echo**          | Gin tiene mejor ecosistema para aprendizaje y más ejemplos |
| **Fiber**         | Más nuevo, menos maduro para proyectos educativos          |
| **Chi**           | Minimalista, requiere más configuración manual             |
| **net/http puro** | Muy bajo nivel para un proyecto de bootcamp                |
| **Gin ✅**        | Mayor adopción, buena documentación, fácil de aprender     |

### ¿Por qué GORM?

| Alternativa | Por qué no se eligió                                         |
| ----------- | ------------------------------------------------------------ |
| **sqlx**    | No es ORM completo, requiere SQL manual                      |
| **ent**     | Curva de aprendizaje alta, generación de código              |
| **sqlc**    | Requiere aprender sintaxis de anotaciones SQL                |
| **GORM ✅** | ORM completo, auto-migrate, scopes, más cercano a SQLAlchemy |

### ¿Por qué golang-migrate (SQL puro)?

- Las migraciones son archivos `.sql` numerados → más transparente para aprender
- No depende del ORM → funciona con cualquier driver PostgreSQL
- `alembic` (Python) genera Python → `golang-migrate` usa SQL estándar
- Versionado explícito y reversible (`.up.sql` y `.down.sql`)

### ¿Por qué JWT stateless?

- No requiere almacenar sesiones en la BD ni en Redis
- Escalable horizontalmente (sin estado compartido entre instancias)
- El access token corto (15 min) limita el riesgo si se compromete
- El refresh token largo (7 días) evita re-autenticación frecuente

### ¿Por qué React + Vite (no Next.js)?

- El foco es aprender autenticación, no SSR/SSG
- Vite tiene HMR ultrarrápido y configuración mínima
- React puro es más didáctico para entender hooks y context
- Next.js añade complejidad innecesaria para el objetivo del bootcamp

---

## 7. Configuración por Entorno

| Aspecto       | Desarrollo                      | Producción                            |
| ------------- | ------------------------------- | ------------------------------------- |
| Base de datos | Docker Compose (localhost:5432) | PostgreSQL gestionado (RDS, Supabase) |
| Backend URL   | http://localhost:8000           | https://api.tu-dominio.com            |
| Frontend URL  | http://localhost:5173           | https://tu-dominio.com                |
| Email         | Mailpit (localhost:1025/8025)   | Resend / SendGrid (SMTP real)         |
| JWT Secret    | Cualquier string ≥32 chars      | String criptográficamente aleatorio   |
| CORS          | http://localhost:5173           | https://tu-dominio.com (exacto)       |
| Logs          | Consola (desarrollo)            | Archivo + sistema de monitoreo        |
| Rate limit    | Permisivo (desarrollo)          | Estricto (10 req/min en auth)         |

---

## 8. Comunicación Backend ↔ Frontend

### Formato de respuesta estándar

- **Éxito:** `{ "data": { ... }, "message": "..." }`
- **Error:** `{ "detail": "Mensaje de error legible" }`
- **Validación:** `{ "detail": [{ "field": "email", "message": "email inválido" }] }`

### Autenticación en headers

```
# Request del frontend con JWT
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Respuesta de error de autenticación
HTTP/1.1 401 Unauthorized
Content-Type: application/json
{ "detail": "Token expirado o inválido" }
```

### CORS

El middleware de CORS permite peticiones solo desde el origen del frontend configurado en la variable de entorno `FRONTEND_URL`. En desarrollo: `http://localhost:5173`.

---

> Recuerda: Una buena arquitectura no es la más compleja, sino la más fácil de entender, cambiar y probar.
