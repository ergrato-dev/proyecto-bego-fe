# HU-009 — Landing Page

## Identificación

| Campo         | Valor                                                                   |
| ------------- | ----------------------------------------------------------------------- |
| **ID**        | HU-009                                                                  |
| **Módulo**    | Frontend / Marketing                                                    |
| **Prioridad** | Media                                                                   |
| **Estado**    | Pendiente                                                               |
| **Sprint**    | Sprint 3                                                                |
| **RFs**       | RF-010                                                                  |
| **RNFs**      | RNF-003 (Usabilidad), RNF-004 (Accesibilidad), RNF-006 (Compatibilidad) |

---

## Historia

**Como** visitante del sitio web,
**quiero** ver una página de inicio que explique el servicio de NN Auth
**y** tenga accesos claros al registro e inicio de sesión,
**para** entender el producto y decidir si quiero crear una cuenta.

---

## Criterios de Aceptación

### CA-009.1 — Contenido de la landing page

**Dado** que soy un visitante en la ruta `/`
**Entonces** veo:

- El nombre y logo de "NN Auth System"
- Una descripción del servicio (qué es y qué ofrece)
- Un botón/enlace prominente hacia `/register`
- Un enlace hacia `/login` para usuarios existentes
- El footer con enlaces a políticas legales

### CA-009.2 — Diseño responsivo

**Dado** que accedo desde un dispositivo móvil (320px de ancho)
**Entonces** el contenido se adapta correctamente sin scroll horizontal
**Y** los botones tienen tamaño suficiente para ser tocados con el dedo (mínimo 44x44px)

### CA-009.3 — Accesibilidad básica

**Dado** que la página carga
**Entonces** tiene una etiqueta `<title>` descriptiva
**Y** la jerarquía de encabezados es lógica (h1 → h2 → h3)
**Y** las imágenes decorativas tienen `alt=""`

### CA-009.4 — Performance

**Dado** que la landing page es la primera impresión del sistema
**Entonces** el First Contentful Paint (FCP) debe ser menor a 2 segundos en conexión 4G
**Y** las imágenes están optimizadas (WebP o SVG)

---

## Notas Técnicas

- Componente: `pages/LandingPage.tsx`
- No requiere autenticación
- Usar TailwindCSS para el diseño
- Sin gradientes (restricción RD-001)
- Colores sólidos, fuentes sans-serif (restricciones de diseño)
