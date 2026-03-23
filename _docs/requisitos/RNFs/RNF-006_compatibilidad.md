# RNF-006 — Compatibilidad

## Identificación

| Campo         | Valor                         |
| ------------- | ----------------------------- |
| **ID**        | RNF-006                       |
| **Nombre**    | Compatibilidad                |
| **Categoría** | Compatibilidad / Portabilidad |
| **Prioridad** | Media                         |

---

## Descripción

El sistema debe funcionar correctamente en los navegadores y dispositivos más utilizados por los usuarios, garantizando una experiencia consistente independientemente del entorno de acceso.

---

## Navegadores Soportados (Frontend)

| Navegador       | Versión mínima | Plataforma        |
| --------------- | -------------- | ----------------- |
| Google Chrome   | Última estable | Desktop / Android |
| Mozilla Firefox | Última estable | Desktop           |
| Safari          | Última estable | macOS / iOS       |
| Microsoft Edge  | Última estable | Desktop           |

> **Nota:** No se requiere soporte para Internet Explorer.

---

## Dispositivos y Resoluciones

| Categoría | Ancho mínimo | Ejemplo de dispositivo      |
| --------- | ------------ | --------------------------- |
| Mobile S  | 320px        | iPhone SE (2016)            |
| Mobile M  | 375px        | iPhone SE (2022), Galaxy S8 |
| Mobile L  | 425px        | iPhone 14 Pro               |
| Tablet    | 768px        | iPad                        |
| Desktop   | 1024px+      | Laptop / Monitor            |

---

## Entornos de Ejecución (Backend)

| Componente | Compatibilidad                    |
| ---------- | --------------------------------- |
| Go         | 1.22+ en Linux, macOS, Windows    |
| PostgreSQL | 17+ (Docker o instalación nativa) |
| Docker     | 24+ con Docker Compose v2         |

---

## Tecnologías Frontend — Requisitos de Compatibilidad

| Feature                | Soporte requerido                                |
| ---------------------- | ------------------------------------------------ |
| CSS Custom Properties  | Chrome 49+, Firefox 31+, Safari 9.1+             |
| CSS Grid               | Chrome 57+, Firefox 52+, Safari 10.1+            |
| `localStorage`         | Todos los navegadores modernos                   |
| `fetch` API            | Todos los navegadores modernos (usado via Axios) |
| `prefers-color-scheme` | Chrome 76+, Firefox 67+, Safari 12.1+            |

> TailwindCSS 4+ y Vite 6+ configuran Browserslist automáticamente para targets modernos.

---

## Verificación Cruzada

```bash
# Verificar el target de compilación de Vite
cat fe/dist/.vite/manifest.json | grep "browserslist"

# Verificar que el build no tiene errores en modo producción
cd fe && pnpm build
```

---

## Criterio de Cumplimiento

- La aplicación se muestra correctamente en Chrome, Firefox y Safari sin errores de consola
- Los formularios funcionan correctamente en dispositivos touch (375px)
- El backend arranca correctamente en Go 1.22+ en Linux y macOS
- `pnpm build` completa sin errores o advertencias relacionadas con compatibilidad
