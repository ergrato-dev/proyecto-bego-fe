# RF-001 — Registro de Usuario

## Identificación

| Campo              | Valor               |
| ------------------ | ------------------- |
| **ID**             | RF-001              |
| **Nombre**         | Registro de usuario |
| **Módulo**         | Autenticación       |
| **HU relacionada** | HU-001              |
| **Prioridad**      | Alta                |
| **Estado**         | Pendiente           |

---

## Descripción

El sistema debe permitir a un visitante crear una nueva cuenta de usuario proporcionando su nombre completo, dirección de email y contraseña. El sistema debe validar los datos, hashear la contraseña con bcrypt, crear el registro en la base de datos y enviar un email de verificación.

---

## Entradas

| Campo       | Tipo   | Obligatorio | Validaciones                                                |
| ----------- | ------ | ----------- | ----------------------------------------------------------- |
| `email`     | string | Sí          | Formato email RFC 5322, máximo 255 caracteres               |
| `full_name` | string | Sí          | Mínimo 2 caracteres, máximo 255 caracteres                  |
| `password`  | string | Sí          | Mínimo 8 chars, al menos 1 mayúscula, 1 minúscula, 1 dígito |

---

## Proceso

1. El handler Go recibe el JSON body y lo bindea al struct `RegisterRequest` (DTO)
2. El validador `go-playground/validator v10` verifica todas las restricciones del struct
3. Si la validación falla → retornar HTTP 422 con los mensajes de error
4. El service busca en la tabla `users` si ya existe un registro con el mismo email (GORM: `db.Where("email = ?", req.Email).First()`)
5. Si el email ya existe → retornar HTTP 409 con "Ya existe una cuenta con este email"
6. La contraseña se hashea: `bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)`
7. Se crea el registro en la tabla `users` con `is_email_verified = false`, `is_active = true`
8. Se genera un token UUID: `uuid.New().String()`
9. Se crea el registro en `email_verification_tokens` con `expires_at = NOW() + 24h`, `used = false`
10. Se envía el email de verificación usando `jordan-wright/email` via SMTP con el enlace: `{FRONTEND_URL}/verify-email?token={token}`
11. Se retorna HTTP 201 con el `UserResponse` DTO (sin `HashedPassword`)

---

## Salidas

| Código HTTP | Situación                   | Cuerpo de respuesta                        |
| ----------- | --------------------------- | ------------------------------------------ |
| 201         | Usuario creado exitosamente | `{ message, data: UserResponse }`          |
| 409         | Email ya registrado         | `{ detail: "Ya existe una cuenta..." }`    |
| 422         | Validación fallida          | `{ detail: "campo: mensaje de error" }`    |
| 429         | Rate limit excedido         | `{ detail: "Demasiadas peticiones" }`      |
| 500         | Error interno del servidor  | `{ detail: "Error interno del servidor" }` |

---

## Endpoint

`POST /api/v1/auth/register`

---

## Reglas de Negocio

| ID     | Regla                                                                     |
| ------ | ------------------------------------------------------------------------- |
| RN-001 | El email debe ser único en la tabla `users`                               |
| RN-002 | La contraseña NUNCA se almacena en texto plano                            |
| RN-003 | El usuario recién creado no puede iniciar sesión hasta verificar su email |
| RN-004 | La contraseña no puede ser igual al email del usuario                     |
| RN-005 | El token de verificación expira en 24 horas                               |
