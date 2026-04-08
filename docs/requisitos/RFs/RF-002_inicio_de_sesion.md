# RF-002 — Inicio de Sesión

## Identificación

| Campo              | Valor            |
| ------------------ | ---------------- |
| **ID**             | RF-002           |
| **Nombre**         | Inicio de sesión |
| **Módulo**         | Autenticación    |
| **HU relacionada** | HU-002           |
| **Prioridad**      | Alta             |

---

## Descripción

El sistema debe permitir a un usuario registrado y con email verificado autenticarse con su email y contraseña. Si las credenciales son válidas, el sistema genera y retorna un access token (15 min) y un refresh token (7 días).

---

## Entradas

| Campo      | Tipo   | Obligatorio | Validaciones         |
| ---------- | ------ | ----------- | -------------------- |
| `email`    | string | Sí          | Formato email válido |
| `password` | string | Sí          | No vacío             |

---

## Proceso

1. El handler bindea el JSON a `LoginRequest` DTO
2. El validador verifica formato básico
3. El service busca el usuario por email: `db.Where("email = ?", req.Email).First(&user)`
4. Si no existe → retornar HTTP 401 con mensaje genérico (no revelar si el email existe)
5. Verificar la contraseña: `bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))`
6. Si la contraseña no coincide → retornar HTTP 401 con mensaje genérico
7. Verificar `user.IsEmailVerified == true` → si no, HTTP 401 con mensaje específico
8. Verificar `user.IsActive == true` → si no, HTTP 401 con mensaje específico
9. Generar access token: `jwt.NewWithClaims(jwt.SigningMethodHS256, claims)` con `exp = 15min`
10. Generar refresh token: mismo proceso con `exp = 7días`
11. Registrar evento `LOGIN_SUCCESS` en el audit log (IP, user-agent, timestamp)
12. Retornar HTTP 200 con `TokenResponse` DTO

---

## Salidas

| Código HTTP | Situación              | Cuerpo                                                  |
| ----------- | ---------------------- | ------------------------------------------------------- |
| 200         | Login exitoso          | `{ access_token, refresh_token, token_type: "bearer" }` |
| 401         | Credenciales inválidas | `{ detail: "Credenciales inválidas" }` (genérico)       |
| 401         | Email no verificado    | `{ detail: "Debes verificar tu email..." }`             |
| 429         | Rate limit excedido    | HTTP 429                                                |

---

## Endpoint

`POST /api/v1/auth/login`

---

## Reglas de Negocio

| ID     | Regla                                                                                                  |
| ------ | ------------------------------------------------------------------------------------------------------ |
| RN-006 | El mensaje de error para email no existente y contraseña incorrecta es el mismo (previene enumeración) |
| RN-007 | Solo usuarios con `is_email_verified = true` pueden iniciar sesión                                     |
| RN-008 | Solo usuarios con `is_active = true` pueden iniciar sesión                                             |
| RN-009 | Rate limit: 10 intentos / 15 minutos por IP                                                            |
