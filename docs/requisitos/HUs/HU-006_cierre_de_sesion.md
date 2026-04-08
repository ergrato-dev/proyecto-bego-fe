# HU-006 — Cierre de Sesión

## Identificación

| Campo         | Valor               |
| ------------- | ------------------- |
| **ID**        | HU-006              |
| **Módulo**    | Autenticación       |
| **Prioridad** | Alta                |
| **Estado**    | Pendiente           |
| **Sprint**    | Sprint 1            |
| **RFs**       | RF-007              |
| **RNFs**      | RNF-001 (Seguridad) |

---

## Historia

**Como** usuario autenticado,
**quiero** poder cerrar mi sesión,
**para** que nadie más pueda acceder a mi cuenta desde este dispositivo.

---

## Criterios de Aceptación

### CA-006.1 — Cierre de sesión exitoso

**Dado** que soy un usuario autenticado
**Cuando** hago clic en el botón "Cerrar sesión" en la barra de navegación
**Entonces** el sistema elimina los tokens del almacenamiento del cliente
**Y** me redirige a la página de inicio o al login
**Y** si intento volver a `/dashboard` con el botón Atrás del navegador, soy redirigido al login

### CA-006.2 — Sesión invalidada en el cliente

**Dado** que cerré sesión
**Cuando** intento acceder a cualquier ruta protegida
**Entonces** veo la página de login, no el contenido protegido

### CA-006.3 — Sin petición al servidor requerida

**Dado** que el sistema es JWT stateless
**Entonces** el cierre de sesión puede implementarse solo en el cliente (eliminar tokens de memoria/sessionStorage)
**Y** no requiere un endpoint de logout en el servidor (en la implementación básica)

---

## Notas Técnicas

- El logout es principalmente del lado del cliente: eliminar el access token y refresh token de la memoria
- JWT es stateless: el servidor no mantiene lista de tokens activos
- El access token expirado naturalmente en 15min es el mecanismo de seguridad principal
- Para mayor seguridad (future scope): implementar token blacklist en BD o Redis
