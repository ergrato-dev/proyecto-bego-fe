# HU-003 — Ver Perfil de Usuario

## Identificación

| Campo         | Valor               |
| ------------- | ------------------- |
| **ID**        | HU-003              |
| **Módulo**    | Usuario             |
| **Prioridad** | Media               |
| **Estado**    | Pendiente           |
| **Sprint**    | Sprint 1            |
| **RFs**       | RF-003              |
| **RNFs**      | RNF-001 (Seguridad) |

---

## Historia

**Como** usuario autenticado,
**quiero** ver los datos de mi perfil en el dashboard,
**para** confirmar que mis datos son correctos y verificar el estado de mi cuenta.

---

## Criterios de Aceptación

### CA-003.1 — Visualización del perfil

**Dado** que soy un usuario autenticado
**Cuando** accedo al Dashboard o a la sección "Mi perfil"
**Entonces** veo mi nombre completo, email y estado de verificación de email

### CA-003.2 — Datos nunca incluyen la contraseña

**Dado** que estoy viendo mi perfil
**Entonces** el sistema NUNCA muestra el `hashed_password` ni ningún campo relacionado con la contraseña en ninguna respuesta

### CA-003.3 — Acceso sin autenticación bloqueado

**Dado** que un usuario no autenticado intenta acceder a `GET /api/v1/users/me`
**Entonces** el sistema devuelve HTTP 401 con "Token de autenticación requerido"
**Y** el frontend lo redirige a la página de login

### CA-003.4 — Token expirado

**Dado** que mi access token ha expirado (15 minutos)
**Cuando** el frontend intenta cargar mi perfil
**Entonces** el interceptor de Axios detecta el 401
**Y** usa el refresh token para obtener un nuevo access token automáticamente
**Y** recarga mi perfil sin que yo lo note

---

## Notas Técnicas

- Endpoint: `GET /api/v1/users/me`
- El middleware `AuthMiddleware` extrae el `user_id` del JWT y lo pasa al handler via contexto de Gin
- El handler usa `user_id` para hacer `db.First(&user, "id = ?", userID)` con GORM
- La respuesta usa `UserResponse` DTO (sin `HashedPassword`)
