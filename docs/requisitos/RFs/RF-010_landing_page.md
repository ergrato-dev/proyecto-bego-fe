# RF-010 — Landing Page

## Identificación

| Campo              | Valor               |
| ------------------ | ------------------- |
| **ID**             | RF-010              |
| **Nombre**         | Landing Page        |
| **Módulo**         | Interfaz de Usuario |
| **HU relacionada** | HU-009              |
| **Prioridad**      | Media               |

---

## Descripción

El sistema debe mostrar una página de inicio pública que presente el proyecto, invite al usuario a registrarse o iniciar sesión, y comunique el propósito del sistema de forma clara y atractiva.

---

## Contenido Requerido

| Sección      | Descripción                                                                   |
| ------------ | ----------------------------------------------------------------------------- |
| **Hero**     | Título principal, subtítulo descriptivo, CTA "Registrarse" y "Iniciar sesión" |
| **Features** | Mínimo 3 características del sistema (seguridad, JWT, formularios accesibles) |
| **Footer**   | Links a páginas legales (Privacidad, Términos), formulario de contacto        |

---

## Proceso

1. Usuario navega a `/`
2. React Router renderiza `LandingPage`
3. No se realiza ninguna llamada a la API
4. Si el usuario ya está autenticado (`AuthContext.user !== null`), se puede redirigir a `/dashboard`

---

## Reglas de Diseño

| Aspecto           | Regla                                       |
| ----------------- | ------------------------------------------- |
| Responsividad     | Mobile-first con breakpoints de TailwindCSS |
| Tema              | Soportar dark y light mode con TailwindCSS  |
| Sin degradados    | Colores sólidos y planos                    |
| Accesibilidad     | Contraste WCAG AA, roles y landmarks ARIA   |
| Tipografía        | Fuentes sans-serif (Inter, system-ui)       |
| Botones de acción | Alineados a la derecha (`justify-end`)      |

---

## Reglas de Negocio

| ID     | Regla                                                        |
| ------ | ------------------------------------------------------------ |
| RN-032 | La landing page es pública — no requiere autenticación       |
| RN-033 | Los CTAs deben enlazar directamente a `/register` y `/login` |
| RN-034 | El footer debe contener links legales y de contacto          |
