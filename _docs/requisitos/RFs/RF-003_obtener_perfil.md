# RF-003 — Obtener Perfil del Usuario

## Identificación

| Campo              | Valor          |
| ------------------ | -------------- |
| **ID**             | RF-003         |
| **Nombre**         | Obtener perfil |
| **Módulo**         | Usuario        |
| **HU relacionada** | HU-003         |
| **Prioridad**      | Alta           |

---

## Descripción

El sistema debe permitir a un usuario autenticado obtener sus propios datos de perfil mediante un endpoint protegido.

---

## Entradas

| Fuente                        | Campo           | Descripción                          |
| ----------------------------- | --------------- | ------------------------------------ |
| Header HTTP                   | `Authorization` | `Bearer <access_token>` (JWT válido) |
| JWT (extraído por middleware) | `sub`           | UUID del usuario autenticado         |

---

## Proceso

1. El middleware `AuthMiddleware` extrae y valida el JWT del header `Authorization`
2. Si el token es inválido o ausente → HTTP 401
3. El middleware inyecta el `user_id` al contexto de Gin: `c.Set("user_id", userID)`
4. El handler extrae el `user_id` del contexto: `c.Get("user_id")`
5. El service busca el usuario: `db.First(&user, "id = ?", userID)`
6. Si no existe → HTTP 404 (no debería ocurrir si el JWT es válido)
7. Se retorna el `UserResponse` DTO (sin `HashedPassword`, sin `used_at`, etc.)

---

## Salidas

| Código HTTP | Situación              | Cuerpo                                                                                                  |
| ----------- | ---------------------- | ------------------------------------------------------------------------------------------------------- |
| 200         | Perfil obtenido        | `{ data: UserResponse }` con id, email, full_name, is_active, is_email_verified, created_at, updated_at |
| 401         | Token ausente/inválido | `{ detail: "Token de autenticación requerido" }`                                                        |
| 500         | Error interno          | `{ detail: "Error interno del servidor" }`                                                              |

---

## Endpoint

`GET /api/v1/users/me`

---

## Reglas de Negocio

| ID     | Regla                                                                     |
| ------ | ------------------------------------------------------------------------- |
| RN-010 | El campo `hashed_password` NUNCA se incluye en la respuesta               |
| RN-011 | Un usuario solo puede ver su propio perfil (el JWT identifica al usuario) |
