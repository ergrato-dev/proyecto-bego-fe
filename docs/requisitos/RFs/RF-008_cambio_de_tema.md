# RF-008 — Cambio de Tema Visual

## Identificación

| Campo              | Valor                       |
| ------------------ | --------------------------- |
| **ID**             | RF-008                      |
| **Nombre**         | Cambio de tema (dark/light) |
| **Módulo**         | Interfaz de Usuario         |
| **HU relacionada** | HU-007                      |
| **Prioridad**      | Media                       |

---

## Descripción

El sistema debe ofrecer dos temas visuales (claro y oscuro) y permitir al usuario alternar entre ellos. La preferencia debe persistir entre sesiones usando `localStorage`, respetando también la preferencia del sistema operativo del usuario.

---

## Proceso

1. Al cargar la aplicación, `ThemeContext` (o lógica en `App.tsx`) lee la preferencia en este orden de prioridad:
   - `localStorage.getItem('theme')` → si existe, usar ese valor
   - `window.matchMedia('(prefers-color-scheme: dark)').matches` → usar preferencia del sistema
   - Fallback: tema claro
2. Aplicar la clase `dark` al elemento `<html>` si el tema es oscuro (patrón de TailwindCSS)
3. El usuario hace clic en el botón de toggle
4. Se cambia la clase en `<html>` y se guarda en `localStorage`

---

## Implementación con TailwindCSS

```
// tailwind.config.ts — activar modo dark con clase
darkMode: 'class'

// App.tsx o ThemeContext.tsx
document.documentElement.classList.toggle('dark')
```

---

## Reglas de Diseño Obligatorias

| Regla             | Detalle                                            |
| ----------------- | -------------------------------------------------- |
| Sin degradados    | Colores sólidos y planos en ambos temas            |
| Contraste mínimo  | WCAG AA — 4.5:1 para texto normal, 3:1 para grande |
| Tipografía        | Fuentes sans-serif únicamente (Inter, system-ui)   |
| Botones de acción | Alineados a la derecha (`justify-end`)             |

---

## Salidas

| Acción             | Resultado                                          |
| ------------------ | -------------------------------------------------- |
| Activar dark mode  | Clase `dark` en `<html>`, guardado en localStorage |
| Activar light mode | Clase `dark` removida, guardado en localStorage    |

---

## Reglas de Negocio

| ID     | Regla                                                           |
| ------ | --------------------------------------------------------------- |
| RN-025 | La preferencia de tema persiste entre sesiones                  |
| RN-026 | Si no hay preferencia guardada, respetar `prefers-color-scheme` |
| RN-027 | El cambio de tema es inmediato y sin recarga de página          |
