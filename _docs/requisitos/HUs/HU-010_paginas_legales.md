# HU-010 — Páginas Legales

## Identificación

| Campo         | Valor                   |
| ------------- | ----------------------- |
| **ID**        | HU-010                  |
| **Módulo**    | Frontend / Legal        |
| **Prioridad** | Baja                    |
| **Estado**    | Pendiente               |
| **Sprint**    | Sprint 3                |
| **RFs**       | RF-011                  |
| **RNFs**      | RNF-004 (Accesibilidad) |

---

## Historia

**Como** visitante o usuario del sitio,
**quiero** poder acceder a las políticas de privacidad y términos de uso,
**para** entender cómo se usan mis datos y cuáles son las condiciones del servicio.

---

## Criterios de Aceptación

### CA-010.1 — Política de Privacidad accesible

**Dado** que estoy en cualquier página de la aplicación
**Cuando** hago clic en "Política de Privacidad" en el footer
**Entonces** se abre la página `/privacy-policy` con el contenido de la política

### CA-010.2 — Términos y Condiciones accesibles

**Dado** que estoy en cualquier página de la aplicación
**Cuando** hago clic en "Términos y Condiciones" en el footer
**Entonces** se abre la página `/terms-of-service` con el contenido de los términos

### CA-010.3 — Enlace en el formulario de registro

**Dado** que estoy en el formulario de registro
**Entonces** hay un texto que dice "Al registrarte, aceptas nuestros Términos y Condiciones y Política de Privacidad" con enlaces funcionales

### CA-010.4 — Páginas legales accesibles sin autenticación

**Dado** que un visitante no autenticado intenta ver las páginas legales
**Entonces** las páginas se muestran correctamente sin requerir login

---

## Notas Técnicas

- Componentes: `pages/PrivacyPolicyPage.tsx`, `pages/TermsOfServicePage.tsx`
- El contenido puede ser estático (hardcodeado) para el alcance del bootcamp
- Los enlaces del footer deben estar en el componente `Footer` compartido
- Rutas públicas: `/privacy-policy`, `/terms-of-service`
