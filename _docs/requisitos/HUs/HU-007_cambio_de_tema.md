# HU-007 — Cambio de Tema (Dark / Light Mode)

## Identificación

| Campo         | Valor                                         |
| ------------- | --------------------------------------------- |
| **ID**        | HU-007                                        |
| **Módulo**    | UI / Preferencias                             |
| **Prioridad** | Media                                         |
| **Estado**    | Pendiente                                     |
| **Sprint**    | Sprint 3                                      |
| **RFs**       | RF-008                                        |
| **RNFs**      | RNF-003 (Usabilidad), RNF-004 (Accesibilidad) |

---

## Historia

**Como** usuario de la aplicación,
**quiero** poder alternar entre el modo oscuro y el modo claro,
**para** usar la aplicación con el estilo visual que prefiero según mi entorno.

---

## Criterios de Aceptación

### CA-007.1 — Toggle de tema funcional

**Dado** que estoy en cualquier página de la aplicación
**Cuando** hago clic en el botón de toggle de tema (ícono de sol/luna)
**Entonces** el tema de toda la aplicación cambia inmediatamente (sin recarga)
**Y** el ícono del botón refleja el tema activo (sol = light, luna = dark)

### CA-007.2 — Persistencia de preferencia

**Dado** que cambié el tema a dark mode
**Cuando** recargo la página o navego a otra ruta
**Entonces** el tema se mantiene en dark mode (guardado en localStorage)

### CA-007.3 — Respeto de la preferencia del sistema

**Dado** que es la primera vez que visito la aplicación
**Cuando** mi sistema operativo tiene configurado el modo oscuro (`prefers-color-scheme: dark`)
**Entonces** la aplicación carga en modo oscuro por defecto

### CA-007.4 — Contraste WCAG AA en ambos modos

**Dado** que los colores del tema (tanto claro como oscuro) deben cumplir WCAG 2.1 AA
**Entonces** la relación de contraste texto/fondo es mínimo 4.5:1 para texto normal
**Y** 3:1 para texto grande y componentes de UI

### CA-007.5 — Toggle accesible con teclado

**Dado** que el botón de toggle es parte de la barra de navegación
**Entonces** es accesible con Tab y puede activarse con Enter/Space
**Y** tiene `aria-label` descriptivo: "Cambiar a modo claro" o "Cambiar a modo oscuro"

---

## Notas Técnicas

- Usar la clase `dark` de TailwindCSS en el elemento `<html>` o `<body>`
- Guardar preferencia en `localStorage` bajo la clave `theme`
- Leer `window.matchMedia('(prefers-color-scheme: dark)')` para el valor inicial
- El toggle está en el componente `Header` o `Navbar`
