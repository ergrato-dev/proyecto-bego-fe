# RF-006 — Restablecer Contraseña con Token

## Identificación

| Campo              | Valor                  |
| ------------------ | ---------------------- |
| **ID**             | RF-006                 |
| **Nombre**         | Restablecer contraseña |
| **Módulo**         | Autenticación          |
| **HU relacionada** | HU-005                 |
| **Prioridad**      | Alta                   |

---

## Descripción

El sistema debe permitir al usuario establecer una nueva contraseña usando el token recibido por email, verificando que el token sea válido, no expirado y no haya sido usado previamente.

---

## Entradas

| Campo          | Tipo   | Obligatorio | Validaciones                          |
| -------------- | ------ | ----------- | ------------------------------------- |
| `token`        | string | Sí          | No vacío                              |
| `new_password` | string | Sí          | Mínimo 8 chars, 1 mayúscula, 1 dígito |

---

## Proceso

1. El handler bindea y valida el JSON con `go-playground/validator`
2. El service busca el token en `password_reset_tokens`: `db.Where("token = ?", req.Token).First(&resetToken)`
3. Si no se encuentra → HTTP 400 "Token inválido o expirado"
4. Verificar `resetToken.ExpiresAt > time.Now()` → si expiró → HTTP 400
5. Verificar `resetToken.Used == false` → si ya fue usado → HTTP 400 "Este enlace ya fue utilizado"
6. Hashea la nueva contraseña con bcrypt
7. Actualiza `hashed_password` y `updated_at` del usuario correspondiente (vía `resetToken.UserID`)
8. Marca el token como usado: `db.Model(&resetToken).Update("used", true)` con GORM
9. Retorna HTTP 200 con mensaje de confirmación

---

## Salidas

| Código HTTP | Situación                | Cuerpo                                                              |
| ----------- | ------------------------ | ------------------------------------------------------------------- |
| 200         | Contraseña restablecida  | `{ message: "Contraseña restablecida. Ya puedes iniciar sesión." }` |
| 400         | Token inválido/expirado  | `{ detail: "Token inválido o expirado" }`                           |
| 400         | Token ya usado           | `{ detail: "Este enlace de recuperación ya fue utilizado" }`        |
| 422         | Validación de contraseña | `{ detail: "mensaje de validación" }`                               |

---

## Endpoint

`POST /api/v1/auth/reset-password`

---

## Reglas de Negocio

| ID     | Regla                                                     |
| ------ | --------------------------------------------------------- |
| RN-019 | Un token de reset solo puede usarse una vez               |
| RN-020 | Después de restablecer, el token se marca `used=true`     |
| RN-021 | La nueva contraseña se hashea con bcrypt antes de guardar |
