# RF-013 — Renovación de Token de Acceso

## Identificación

| Campo              | Valor                              |
| ------------------ | ---------------------------------- |
| **ID**             | RF-013                             |
| **Nombre**         | Renovación de token (refresh)      |
| **Módulo**         | Autenticación                      |
| **HU relacionada** | HU-002 (extiende la sesión activa) |
| **Prioridad**      | Alta                               |

---

## Descripción

El sistema debe permitir renovar el access token usando un refresh token válido, sin requerir que el usuario ingrese sus credenciales nuevamente. El interceptor de Axios en el frontend maneja este proceso de forma transparente.

---

## Entradas

| Campo           | Tipo   | Obligatorio | Descripción                           |
| --------------- | ------ | ----------- | ------------------------------------- |
| `refresh_token` | string | Sí          | JWT de larga duración (7 días de TTL) |

---

## Proceso (Backend)

1. El handler recibe el refresh token en el body del request
2. `golang-jwt/jwt` parsea y verifica la firma del token con la `SECRET_KEY`
3. Verificar que el tipo de claim sea `refresh` (no `access`)
4. Verificar que el token no haya expirado (`exp` claim)
5. Extraer el `user_id` de los claims
6. Buscar el usuario en la BD para confirmar que sigue activo (`is_active=true`)
7. Generar nuevo `access_token` (15 min) y nuevo `refresh_token` (7 días)
8. Retornar ambos tokens en la respuesta

---

## Proceso (Frontend — Interceptor Axios)

1. El interceptor de respuesta en `api/axios.ts` detecta HTTP 401
2. Recupera el `refresh_token` del almacenamiento
3. Hace POST a `/api/v1/auth/refresh` con el refresh token
4. Si tiene éxito: guarda los nuevos tokens y reintenta el request original
5. Si falla (refresh también expirado): llama a `logout()` y redirige a `/login`

---

## Salidas

| Código HTTP | Situación              | Cuerpo                                                  |
| ----------- | ---------------------- | ------------------------------------------------------- |
| 200         | Tokens renovados       | `{ access_token, refresh_token, token_type: "bearer" }` |
| 401         | Refresh token inválido | `{ detail: "Token de renovación inválido" }`            |
| 401         | Refresh token expirado | `{ detail: "Sesión expirada, inicia sesión de nuevo" }` |

---

## Endpoint

`POST /api/v1/auth/refresh`

---

## Reglas de Negocio

| ID     | Regla                                                                      |
| ------ | -------------------------------------------------------------------------- |
| RN-041 | Solo un refresh token (no un access token) puede renovar la sesión         |
| RN-042 | El refresh token también se renueva en cada llamada exitosa                |
| RN-043 | Si el usuario está inactivo (`is_active=false`), retornar 401              |
| RN-044 | El interceptor de Axios maneja el refresh de forma transparente al usuario |
