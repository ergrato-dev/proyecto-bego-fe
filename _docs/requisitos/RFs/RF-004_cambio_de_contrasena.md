# RF-004 — Cambio de Contraseña

## Identificación

| Campo              | Valor                |
| ------------------ | -------------------- |
| **ID**             | RF-004               |
| **Nombre**         | Cambio de contraseña |
| **Módulo**         | Autenticación        |
| **HU relacionada** | HU-004               |
| **Prioridad**      | Alta                 |

---

## Descripción

El sistema debe permitir a un usuario autenticado cambiar su contraseña actual por una nueva, verificando primero que conoce la contraseña actual.

---

## Entradas

| Campo                  | Tipo   | Obligatorio | Validaciones                          |
| ---------------------- | ------ | ----------- | ------------------------------------- |
| `current_password`     | string | Sí          | No vacío                              |
| `new_password`         | string | Sí          | Mínimo 8 chars, 1 mayúscula, 1 dígito |
| Header `Authorization` | —      | Sí          | `Bearer <access_token>` JWT válido    |

---

## Proceso

1. El middleware `AuthMiddleware` valida el JWT y extrae el `user_id`
2. El handler bindea el JSON a `ChangePasswordRequest` DTO
3. `go-playground/validator` verifica las restricciones de `new_password`
4. El service busca el usuario en BD por `user_id`
5. Verifica la contraseña actual: `bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.CurrentPassword))`
6. Si no coincide → HTTP 400 "La contraseña actual es incorrecta"
7. Hashea la nueva contraseña: `bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)`
8. Actualiza `hashed_password` y `updated_at` en la tabla `users` con GORM
9. Registra evento `PASSWORD_CHANGED` en el audit log
10. Retorna HTTP 200 con mensaje de confirmación

---

## Salidas

| Código HTTP | Situación                        | Cuerpo                                               |
| ----------- | -------------------------------- | ---------------------------------------------------- |
| 200         | Contraseña cambiada exitosamente | `{ message: "Contraseña actualizada exitosamente" }` |
| 400         | Contraseña actual incorrecta     | `{ detail: "La contraseña actual es incorrecta" }`   |
| 401         | Token ausente o inválido         | `{ detail: "Token de autenticación requerido" }`     |
| 422         | Validación fallida               | `{ detail: "mensaje de validación" }`                |

---

## Endpoint

`POST /api/v1/auth/change-password`

---

## Reglas de Negocio

| ID     | Regla                                                                  |
| ------ | ---------------------------------------------------------------------- |
| RN-012 | La nueva contraseña debe cumplir los mismos requisitos que el registro |
| RN-013 | La nueva contraseña puede ser igual a la actual (no hay restricción)   |
| RN-014 | Se requiere verificar la contraseña actual antes de cambiarla          |
