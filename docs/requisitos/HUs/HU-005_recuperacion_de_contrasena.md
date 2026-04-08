# HU-005 — Recuperación de Contraseña

## Identificación

| Campo         | Valor                                     |
| ------------- | ----------------------------------------- |
| **ID**        | HU-005                                    |
| **Módulo**    | Autenticación                             |
| **Prioridad** | Alta                                      |
| **Estado**    | Pendiente                                 |
| **Sprint**    | Sprint 2                                  |
| **RFs**       | RF-005, RF-006                            |
| **RNFs**      | RNF-001 (Seguridad), RNF-003 (Usabilidad) |

---

## Historia

**Como** usuario que olvidó su contraseña,
**quiero** poder solicitar un enlace de recuperación por email
**y** establecer una nueva contraseña usando ese enlace,
**para** poder volver a acceder a mi cuenta sin asistencia manual.

---

## Criterios de Aceptación

### CA-005.1 — Solicitud de recuperación (Forgot Password)

**Dado** que estoy en la página `/forgot-password`
**Y** ingreso el email de mi cuenta
**Cuando** envío el formulario
**Entonces** el sistema siempre muestra: "Si el email está registrado, recibirás un enlace de recuperación"
**Y** si el email existe, envía un email con un enlace que expira en 1 hora

### CA-005.2 — Seguridad: misma respuesta para email existente y no existente

**Dado** que ingreso un email no registrado en el sistema
**Cuando** envío el formulario de recuperación
**Entonces** el sistema muestra exactamente el mismo mensaje que si el email SÍ existiera
**Y** no revela si el email está o no registrado

### CA-005.3 — Restablecimiento exitoso (Reset Password)

**Dado** que tengo un enlace de recuperación válido y no expirado
**Y** accedo a la página `/reset-password?token=<token>`
**Y** ingreso y confirmo una nueva contraseña
**Cuando** envío el formulario
**Entonces** el sistema actualiza mi contraseña con el nuevo hash bcrypt
**Y** marca el token como `used = true`
**Y** muestra "Contraseña restablecida. Ya puedes iniciar sesión."
**Y** me redirige al login

### CA-005.4 — Token expirado

**Dado** que el token de recuperación ha expirado (más de 1 hora)
**Cuando** intento usar ese token
**Entonces** el sistema muestra "Token inválido o expirado"
**Y** ofrece la opción de solicitar un nuevo enlace

### CA-005.5 — Token ya usado

**Dado** que ya usé el token de recuperación para cambiar mi contraseña
**Cuando** intento usar el mismo token de nuevo
**Entonces** el sistema muestra "Este enlace de recuperación ya fue utilizado"

### CA-005.6 — Rate limiting en forgot-password

**Dado** que alguien envía más de 3 solicitudes de forgot-password en 15 minutos
**Entonces** el sistema devuelve HTTP 429

---

## Notas Técnicas

- **Paso 1:** `POST /api/v1/auth/forgot-password` — genera token UUID, guarda en `password_reset_tokens`, envía email
- **Paso 2:** `POST /api/v1/auth/reset-password` — verifica token, actualiza contraseña, marca `used=true`
- El token de reset es un UUID generado con `uuid.New()` (no un JWT)
- El email incluye el enlace: `{FRONTEND_URL}/reset-password?token={token}`
- El frontend extrae el token del query parameter y lo incluye en el body del POST
