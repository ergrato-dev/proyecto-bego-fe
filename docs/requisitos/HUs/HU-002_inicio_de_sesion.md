# HU-002 — Inicio de Sesión

## Identificación

| Campo         | Valor                                      |
| ------------- | ------------------------------------------ |
| **ID**        | HU-002                                     |
| **Módulo**    | Autenticación                              |
| **Prioridad** | Alta                                       |
| **Estado**    | Pendiente                                  |
| **Sprint**    | Sprint 1                                   |
| **RFs**       | RF-002                                     |
| **RNFs**      | RNF-001 (Seguridad), RNF-002 (Rendimiento) |

---

## Historia

**Como** usuario registrado con email verificado,
**quiero** iniciar sesión con mi email y contraseña,
**para** acceder al dashboard y funcionalidades protegidas.

---

## Criterios de Aceptación

### CA-002.1 — Login exitoso

**Dado** que soy un usuario con email verificado y cuenta activa
**Cuando** ingreso mis credenciales correctas en el formulario de login
**Y** hago clic en "Iniciar sesión"
**Entonces** el sistema genera un access token (15 min) y un refresh token (7 días)
**Y** me redirige al Dashboard
**Y** mi nombre aparece en la barra de navegación

### CA-002.2 — Credenciales incorrectas

**Dado** que ingreso una contraseña incorrecta
**Cuando** intento iniciar sesión
**Entonces** el sistema muestra el mensaje genérico "Credenciales inválidas"
**Y** NO revela si el email existe o si solo la contraseña es incorrecta
**Y** el contador de intentos aumenta para el rate limiting

### CA-002.3 — Email no verificado

**Dado** que soy un usuario registrado pero no he verificado mi email
**Cuando** intento iniciar sesión con credenciales correctas
**Entonces** el sistema muestra "Debes verificar tu email antes de iniciar sesión"
**Y** opcionalmente ofrece la opción de reenviar el email de verificación

### CA-002.4 — Cuenta desactivada

**Dado** que mi cuenta ha sido desactivada por el administrador (is_active = false)
**Cuando** intento iniciar sesión
**Entonces** el sistema muestra "Tu cuenta ha sido desactivada"

### CA-002.5 — Rate limiting

**Dado** que un atacante hace más de 10 intentos de login en 15 minutos desde la misma IP
**Cuando** intenta el intento número 11
**Entonces** el sistema devuelve HTTP 429 "Demasiadas peticiones. Intenta en unos minutos."

### CA-002.6 — Redirección desde ruta protegida

**Dado** que intento acceder a `/dashboard` sin estar autenticado
**Cuando** el sistema me redirige a `/login`
**Y** me autentico exitosamente
**Entonces** el sistema me redirige de regreso a `/dashboard`

### CA-002.7 — Mostrar/ocultar contraseña

**Dado** que estoy en el formulario de login
**Cuando** hago clic en el ícono de ojo junto al campo de contraseña
**Entonces** el campo alterna entre `type="password"` y `type="text"`

---

## Notas Técnicas

- Endpoint: `POST /api/v1/auth/login`
- Se verifica `is_email_verified = true` y `is_active = true` antes de generar tokens
- Los tokens JWT se generan con `golang-jwt/jwt v5`
- Registro de evento `LOGIN_SUCCESS` o `LOGIN_FAILED` en el audit log
- Los tokens se almacenan temporalmente en el cliente (no en localStorage)
