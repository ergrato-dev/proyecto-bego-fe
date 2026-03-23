# HU-011 — Formulario de Contacto

## Identificación

| Campo         | Valor                   |
| ------------- | ----------------------- |
| **ID**        | HU-011                  |
| **Módulo**    | Frontend / Soporte      |
| **Prioridad** | Baja                    |
| **Estado**    | Pendiente               |
| **Sprint**    | Sprint 3                |
| **RFs**       | RF-012                  |
| **RNFs**      | RNF-004 (Accesibilidad) |

---

## Historia

**Como** visitante o usuario del sitio,
**quiero** poder enviar un mensaje de contacto al equipo de NN,
**para** hacer preguntas, reportar problemas o solicitar soporte.

---

## Criterios de Aceptación

### CA-011.1 — Formulario de contacto funcional

**Dado** que estoy en la página `/contact`
**Y** completo todos los campos requeridos (nombre, email, asunto, mensaje)
**Cuando** envío el formulario
**Entonces** el sistema muestra "Tu mensaje fue enviado. Nos contactaremos pronto."
**Y** limpia el formulario

### CA-011.2 — Validación de campos

**Dado** que intento enviar el formulario con campos vacíos
**Entonces** el sistema muestra mensajes de validación en cada campo requerido

### CA-011.3 — El formulario es accesible

**Dado** que el formulario tiene múltiples campos
**Entonces** cada campo tiene su `<label>` asociada
**Y** los mensajes de error están asociados con `aria-describedby`

### CA-011.4 — Accesible sin autenticación

**Dado** que el formulario de contacto debe ser accesible para cualquier persona
**Entonces** no requiere autenticación para ser usado

---

## Notas Técnicas

- Componente: `pages/ContactPage.tsx`
- Para el alcance del bootcamp, el "envío" puede ser simulado (frontend only, no backend real)
- O puede enviar a un servicio externo de formularios como Formspree
- Ruta pública: `/contact`
- Enlace en el footer
