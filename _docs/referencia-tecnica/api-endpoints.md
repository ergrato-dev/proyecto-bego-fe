# Referencia de Endpoints — NN Auth System API

<!--
  ¿Qué? Documentación completa de todos los endpoints REST de la API.
  ¿Para qué? Que el equipo de frontend y backend tengan una fuente única de verdad
             sobre el contrato de la API: rutas, métodos, parámetros y respuestas.
  ¿Impacto? Sin esto, cada desarrollador asumiría formatos diferentes,
             generando bugs de integración difíciles de rastrear.
-->

---

## Información General

| Campo         | Valor                                   |
| ------------- | --------------------------------------- |
| URL Base      | `http://localhost:8000/api/v1`          |
| Formato       | JSON (`Content-Type: application/json`) |
| Autenticación | JWT Bearer Token                        |
| Versión API   | v1                                      |

---

## Convenciones de Respuesta

### Respuesta exitosa

```json
{
  "message": "Operación exitosa",
  "data": { ... }
}
```

### Respuesta de error

```json
{
  "detail": "Descripción del error legible para el usuario"
}
```

### Respuesta de error de validación

```json
{
  "detail": "email: must be a valid email address; password: min length is 8"
}
```

### Códigos HTTP usados

| Código | Significado           | Cuándo se usa                               |
| ------ | --------------------- | ------------------------------------------- |
| 200    | OK                    | Petición exitosa (GET, login, refresh)      |
| 201    | Created               | Recurso creado (registro de usuario)        |
| 400    | Bad Request           | Datos inválidos o lógica de negocio fallida |
| 401    | Unauthorized          | Token ausente, expirado o inválido          |
| 404    | Not Found             | Recurso no encontrado                       |
| 409    | Conflict              | Email ya registrado                         |
| 422    | Unprocessable Entity  | Error de validación de campos               |
| 429    | Too Many Requests     | Rate limit excedido                         |
| 500    | Internal Server Error | Error inesperado del servidor               |

---

## Endpoints de Autenticación (`/api/v1/auth/`)

### POST `/register` — Registrar nuevo usuario

Crea una nueva cuenta de usuario y envía un email de verificación.

**Autenticación requerida:** No
**Rate limit:** 5 peticiones / 15 minutos por IP

#### Request Body

```json
{
  "email": "maria@ejemplo.com",
  "full_name": "María García",
  "password": "MiContraseña123"
}
```

#### Validaciones del body

| Campo       | Tipo   | Obligatorio | Reglas                                                      |
| ----------- | ------ | ----------- | ----------------------------------------------------------- |
| `email`     | string | Sí          | Formato email válido, máx 255 chars                         |
| `full_name` | string | Sí          | Mínimo 2 chars, máximo 255 chars                            |
| `password`  | string | Sí          | Mínimo 8 chars, al menos 1 mayúscula, 1 minúscula, 1 número |

#### Respuestas

**201 Created** — Usuario creado exitosamente

```json
{
  "message": "Usuario registrado. Revisa tu email para verificar tu cuenta.",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "maria@ejemplo.com",
    "full_name": "María García",
    "is_active": true,
    "is_email_verified": false,
    "created_at": "2026-03-15T10:30:00Z"
  }
}
```

**409 Conflict** — Email ya registrado

```json
{
  "detail": "Ya existe una cuenta con este email"
}
```

**422 Unprocessable Entity** — Campos inválidos

```json
{
  "detail": "password: debe tener al menos 8 caracteres"
}
```

---

### POST `/login` — Iniciar sesión

Autentica al usuario y devuelve los tokens de acceso.

**Autenticación requerida:** No
**Rate limit:** 10 peticiones / 15 minutos por IP

#### Request Body

```json
{
  "email": "maria@ejemplo.com",
  "password": "MiContraseña123"
}
```

#### Respuestas

**200 OK** — Login exitoso

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "bearer"
}
```

**401 Unauthorized** — Credenciales incorrectas

```json
{
  "detail": "Credenciales inválidas"
}
```

**401 Unauthorized** — Email no verificado

```json
{
  "detail": "Debes verificar tu email antes de iniciar sesión"
}
```

**401 Unauthorized** — Cuenta inactiva

```json
{
  "detail": "Tu cuenta ha sido desactivada"
}
```

> **Nota de seguridad:** Los mensajes de error de login son genéricos intencionalmente para no revelar si un email está registrado (previene enumeración de usuarios — OWASP A07).

---

### POST `/refresh` — Renovar access token

Genera un nuevo access token usando el refresh token.

**Autenticación requerida:** No (usa refresh_token)
**Rate limit:** 20 peticiones / 15 minutos por IP

#### Request Body

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Respuestas

**200 OK** — Token renovado

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "bearer"
}
```

**401 Unauthorized** — Refresh token inválido o expirado

```json
{
  "detail": "Refresh token inválido o expirado"
}
```

---

### POST `/change-password` — Cambiar contraseña

Cambia la contraseña del usuario autenticado.

**Autenticación requerida:** Sí (`Authorization: Bearer <access_token>`)
**Rate limit:** 5 peticiones / 15 minutos por IP

#### Headers

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Request Body

```json
{
  "current_password": "MiContraseñaActual123",
  "new_password": "NuevaContraseñaSegura456"
}
```

#### Validaciones

| Campo              | Reglas                                                      |
| ------------------ | ----------------------------------------------------------- |
| `current_password` | Obligatorio                                                 |
| `new_password`     | Mínimo 8 chars, al menos 1 mayúscula, 1 minúscula, 1 número |

#### Respuestas

**200 OK** — Contraseña cambiada

```json
{
  "message": "Contraseña actualizada exitosamente"
}
```

**400 Bad Request** — Contraseña actual incorrecta

```json
{
  "detail": "La contraseña actual es incorrecta"
}
```

**401 Unauthorized** — Token ausente o inválido

```json
{
  "detail": "Token de autenticación requerido"
}
```

---

### POST `/forgot-password` — Solicitar recuperación de contraseña

Envía un email con enlace de recuperación.

**Autenticación requerida:** No
**Rate limit:** 3 peticiones / 15 minutos por IP

#### Request Body

```json
{
  "email": "maria@ejemplo.com"
}
```

#### Respuestas

**200 OK** — Siempre (independientemente de si el email existe)

```json
{
  "message": "Si el email está registrado, recibirás un enlace de recuperación"
}
```

> **Nota de seguridad:** La respuesta es siempre la misma para no revelar si un email existe en el sistema (previene enumeración — OWASP A07).

---

### POST `/reset-password` — Restablecer contraseña con token

Restablece la contraseña usando el token recibido por email.

**Autenticación requerida:** No (usa token del email)
**Rate limit:** 5 peticiones / 15 minutos por IP

#### Request Body

```json
{
  "token": "abc123def456ghi789...",
  "new_password": "NuevaContraseña789"
}
```

#### Respuestas

**200 OK** — Contraseña restablecida

```json
{
  "message": "Contraseña restablecida exitosamente. Ya puedes iniciar sesión."
}
```

**400 Bad Request** — Token inválido

```json
{
  "detail": "Token inválido o expirado"
}
```

**400 Bad Request** — Token ya usado

```json
{
  "detail": "Este enlace de recuperación ya fue utilizado"
}
```

---

### POST `/verify-email` — Verificar dirección de email

Confirma la dirección de email con el token recibido.

**Autenticación requerida:** No (usa token del email)
**Rate limit:** 10 peticiones / 15 minutos por IP

#### Request Body

```json
{
  "token": "xyz789abc123..."
}
```

#### Respuestas

**200 OK** — Email verificado

```json
{
  "message": "Email verificado exitosamente. Ya puedes iniciar sesión."
}
```

**400 Bad Request** — Token inválido o expirado

```json
{
  "detail": "Token de verificación inválido o expirado"
}
```

---

## Endpoints de Usuario (`/api/v1/users/`)

### GET `/me` — Obtener perfil del usuario autenticado

Devuelve los datos del usuario actualmente autenticado.

**Autenticación requerida:** Sí (`Authorization: Bearer <access_token>`)

#### Headers

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Respuestas

**200 OK** — Perfil del usuario

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "maria@ejemplo.com",
    "full_name": "María García",
    "is_active": true,
    "is_email_verified": true,
    "created_at": "2026-03-15T10:30:00Z",
    "updated_at": "2026-03-15T15:45:00Z"
  }
}
```

**401 Unauthorized** — Token ausente o inválido

```json
{
  "detail": "Token de autenticación requerido"
}
```

---

## Tabla Resumen de Endpoints

| Método | Ruta                           | Auth | Rate Limit      | Descripción                    |
| ------ | ------------------------------ | ---- | --------------- | ------------------------------ |
| POST   | `/api/v1/auth/register`        | No   | 5/15min por IP  | Registrar nuevo usuario        |
| POST   | `/api/v1/auth/login`           | No   | 10/15min por IP | Iniciar sesión                 |
| POST   | `/api/v1/auth/refresh`         | No†  | 20/15min por IP | Renovar access token           |
| POST   | `/api/v1/auth/change-password` | Sí   | 5/15min por IP  | Cambiar contraseña (auth)      |
| POST   | `/api/v1/auth/forgot-password` | No   | 3/15min por IP  | Solicitar recuperación email   |
| POST   | `/api/v1/auth/reset-password`  | No†  | 5/15min por IP  | Restablecer contraseña (token) |
| POST   | `/api/v1/auth/verify-email`    | No†  | 10/15min por IP | Verificar email (token)        |
| GET    | `/api/v1/users/me`             | Sí   | Sin límite      | Obtener perfil propio          |

†: No requiere JWT, pero usa un token de uso único enviado por email.

---

## Errores del Servidor

### Formato de errores inesperados

```json
{
  "detail": "Error interno del servidor. Por favor intenta más tarde."
}
```

> **Regla:** Los errores internos (500) NUNCA revelan detalles técnicos al cliente (stack traces, queries SQL, etc.). Solo se loggean en el servidor.

---

## Pruebas con curl

```bash
# Registrar usuario
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@ejemplo.com","full_name":"Test User","password":"Password123"}'

# Login
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@ejemplo.com","password":"Password123"}'

# Obtener perfil (reemplaza <TOKEN> con el access_token del login)
curl -X GET http://localhost:8000/api/v1/users/me \
  -H "Authorization: Bearer <TOKEN>"

# Cambiar contraseña
curl -X POST http://localhost:8000/api/v1/auth/change-password \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"current_password":"Password123","new_password":"NewPassword456"}'
```
