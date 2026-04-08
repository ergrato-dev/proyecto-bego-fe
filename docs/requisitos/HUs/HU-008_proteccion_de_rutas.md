# HU-008 — Protección de Rutas

## Identificación

| Campo         | Valor                |
| ------------- | -------------------- |
| **ID**        | HU-008               |
| **Módulo**    | Frontend / Seguridad |
| **Prioridad** | Alta                 |
| **Estado**    | Pendiente            |
| **Sprint**    | Sprint 1             |
| **RFs**       | RF-009               |
| **RNFs**      | RNF-001 (Seguridad)  |

---

## Historia

**Como** sistema,
**quiero** proteger automáticamente las rutas que requieren autenticación,
**para** que ningún usuario no autenticado pueda acceder a páginas privadas.

---

## Criterios de Aceptación

### CA-008.1 — Redirección a login desde ruta protegida

**Dado** que un usuario no autenticado intenta acceder a `/dashboard`
**Cuando** el sistema verifica que no hay sesión activa
**Entonces** el sistema redirige inmediatamente a `/login`
**Y** guarda la ruta original en el estado para redirigir de vuelta después del login

### CA-008.2 — Rutas protegidas no renderizan contenido privado

**Dado** que el código fuente de la aplicación carga
**Entonces** el componente `ProtectedRoute` verifica la autenticación ANTES de renderizar el contenido
**Y** un usuario no autenticado NUNCA ve el contenido del Dashboard, ni siquiera brevemente

### CA-008.3 — Rutas públicas accesibles sin autenticación

**Dado** que un usuario no autenticado
**Entonces** puede acceder a: Landing Page (`/`), Login (`/login`), Register (`/register`), Forgot Password (`/forgot-password`), Reset Password (`/reset-password`), páginas legales y contacto

### CA-008.4 — Usuario autenticado redirigido fuera de rutas de auth

**Dado** que soy un usuario autenticado
**Cuando** intento acceder a `/login` o `/register`
**Entonces** el sistema me redirige al Dashboard (no tiene sentido hacer login si ya estoy logueado)

### CA-008.5 — Estado de carga durante verificación

**Dado** que la aplicación está verificando si tengo una sesión activa al cargar
**Entonces** se muestra un spinner o skeleton mientras se verifica
**Y** NO se produce un flash de la página de login antes de mostrar el dashboard

---

## Notas Técnicas

- Implementar `ProtectedRoute` como componente de React Router que usa `<Outlet />`
- El `AuthContext` verifica el token al montar (useEffect inicial)
- La verificación inicial hace `GET /api/v1/users/me` con el token guardado
- Si la verificación falla → limpiar tokens y mostrar login
- `isLoading` en el contexto previene el flash de contenido no autorizado
