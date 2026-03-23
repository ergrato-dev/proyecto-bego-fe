# RF-005 — Solicitar Recuperación de Contraseña

## Identificación

| Campo              | Valor                  |
| ------------------ | ---------------------- |
| **ID**             | RF-005                 |
| **Nombre**         | Solicitar recuperación |
| **Módulo**         | Autenticación          |
| **HU relacionada** | HU-005                 |
| **Prioridad**      | Alta                   |

---

## Descripción

El sistema debe permitir a cualquier visitante solicitar la recuperación de contraseña ingresando su email. El sistema genera un token de reset y envía un email con el enlace, sin revelar si el email está registrado.

---

## Entradas

| Campo   | Tipo   | Obligatorio | Validaciones         |
| ------- | ------ | ----------- | -------------------- |
| `email` | string | Sí          | Formato email válido |

---

## Proceso

1. El handler recibe y valida el email con `go-playground/validator`
2. El service busca el usuario por email en la tabla `users`
3. **Si el email NO existe:** NO se hace nada más. Se retorna HTTP 200 con mensaje genérico
4. **Si el email SÍ existe:**
   a. Genera un token UUID: `uuid.New().String()`
   b. Crea un registro en `password_reset_tokens` con `expires_at = NOW() + 1h`, `used = false`
   c. Envía email con el enlace: `{FRONTEND_URL}/reset-password?token={token}`
5. En ambos casos (email existe o no), retorna HTTP 200 con el mismo mensaje genérico

---

## Salidas

| Código HTTP | Situación               | Cuerpo                                                               |
| ----------- | ----------------------- | -------------------------------------------------------------------- |
| 200         | Siempre (por seguridad) | `{ message: "Si el email está registrado, recibirás un enlace..." }` |
| 429         | Rate limit excedido     | HTTP 429                                                             |

---

## Endpoint

`POST /api/v1/auth/forgot-password`

---

## Reglas de Negocio

| ID     | Regla                                                                                                         |
| ------ | ------------------------------------------------------------------------------------------------------------- |
| RN-015 | La respuesta SIEMPRE es la misma, independientemente de si el email existe                                    |
| RN-016 | El token de reset expira en 1 hora                                                                            |
| RN-017 | Rate limit: 3 peticiones / 15 minutos por IP (más estricto, protege spam)                                     |
| RN-018 | Se pueden crear múltiples tokens de reset para el mismo usuario (los anteriores siguen válidos hasta expirar) |
