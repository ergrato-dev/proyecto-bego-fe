# HU-004 — Cambio de Contraseña

## Identificación

| Campo         | Valor               |
| ------------- | ------------------- |
| **ID**        | HU-004              |
| **Módulo**    | Autenticación       |
| **Prioridad** | Alta                |
| **Estado**    | Pendiente           |
| **Sprint**    | Sprint 2            |
| **RFs**       | RF-004              |
| **RNFs**      | RNF-001 (Seguridad) |

---

## Historia

**Como** usuario autenticado,
**quiero** poder cambiar mi contraseña actual por una nueva,
**para** mantener la seguridad de mi cuenta cuando lo necesite.

---

## Criterios de Aceptación

### CA-004.1 — Cambio exitoso

**Dado** que soy un usuario autenticado en la página de cambio de contraseña
**Y** ingreso mi contraseña actual correcta
**Y** una nueva contraseña que cumple los requisitos
**Y** confirmo la nueva contraseña (coinciden)
**Cuando** envío el formulario
**Entonces** el sistema actualiza mi contraseña en la BD con bcrypt
**Y** muestra el mensaje "Contraseña actualizada exitosamente"
**Y** registra el evento `PASSWORD_CHANGED` en el audit log

### CA-004.2 — Contraseña actual incorrecta

**Dado** que ingreso una contraseña actual incorrecta
**Cuando** envío el formulario
**Entonces** el sistema muestra "La contraseña actual es incorrecta"
**Y** NO cambia la contraseña en la BD

### CA-004.3 — Nueva contraseña no coincide con confirmación

**Dado** que la nueva contraseña y la confirmación no son iguales
**Cuando** intento enviar el formulario
**Entonces** el sistema muestra "Las contraseñas no coinciden" antes de enviar al servidor

### CA-004.4 — Nueva contraseña débil

**Dado** que la nueva contraseña no cumple los requisitos mínimos
**Cuando** intento enviar el formulario
**Entonces** el sistema muestra el error de validación específico

### CA-004.5 — Estado de carga visible

**Dado** que envío el formulario
**Cuando** el servidor está procesando la petición
**Entonces** el botón muestra "Actualizando..." y está deshabilitado para prevenir doble envío

### CA-004.6 — Acceso sin autenticación bloqueado

**Dado** que un usuario no autenticado intenta acceder a la página `/change-password`
**Entonces** el sistema lo redirige a `/login`

---

## Notas Técnicas

- Endpoint: `POST /api/v1/auth/change-password`
- Requiere header `Authorization: Bearer <access_token>`
- La validación "contraseñas coinciden" se hace en el frontend antes de enviar
- El backend re-verifica la contraseña actual con `bcrypt.CompareHashAndPassword`
- Rate limiting: 5 peticiones / 15 minutos por IP
