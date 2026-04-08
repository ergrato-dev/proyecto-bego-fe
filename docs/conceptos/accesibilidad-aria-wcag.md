# Accesibilidad Web — ARIA y WCAG 2.1

<!--
  ¿Qué? Guía de accesibilidad web para el frontend de NN Auth System.
  ¿Para qué? Que los formularios de autenticación sean usables por personas
             con discapacidades visuales, motoras o cognitivas.
  ¿Impacto? Sin accesibilidad, excluimos a millones de usuarios.
             Además, en muchos países, la accesibilidad web es un requisito legal.
-->

---

## ¿Qué es la Accesibilidad Web?

La accesibilidad web (a11y) significa diseñar y desarrollar páginas web que puedan ser usadas por todas las personas, independientemente de sus capacidades.

**Personas beneficiadas:**

- Usuarios de lectores de pantalla (ceguera, baja visión)
- Usuarios que navegan solo con teclado (parálisis, temblores)
- Personas con daltonismo
- Personas con dificultades cognitivas
- Usuarios en conexiones lentas o dispositivos limitados

---

## WCAG 2.1 — Web Content Accessibility Guidelines

La WCAG es el estándar internacional de accesibilidad, publicado por el W3C. Define tres niveles de conformidad:

| Nivel | Descripción                              | Requisito en este proyecto |
| ----- | ---------------------------------------- | -------------------------- |
| A     | Básico — elimina las barreras más graves | Obligatorio                |
| AA    | Intermedio — estándar de la industria    | Obligatorio                |
| AAA   | Avanzado — para contextos especializados | Opcional                   |

**En este proyecto: cumplir nivel WCAG 2.1 AA.**

---

## Los 4 Principios de WCAG (POUR)

### P — Perceptible (Perceivable)

La información debe ser presentable de formas que los usuarios puedan percibir.

**Criterios clave:**

**1.1.1 — Texto alternativo (Nivel A)**

Todas las imágenes deben tener un atributo `alt` descriptivo.

```tsx
{
  /* ✅ Correcto */
}
<img src="/logo.svg" alt="NN Auth System — Logo" />;

{
  /* Para imágenes decorativas, alt vacío */
}
<img src="/decoration.png" alt="" role="presentation" />;

{
  /* ❌ Incorrecto — sin alt */
}
<img src="/logo.svg" />;
```

**1.3.1 — Información y Relaciones (Nivel A)**

La estructura semántica del formulario debe comunicarse al asistente.

```tsx
{
  /* ✅ Correcto — label asociado explícitamente con htmlFor */
}
<div>
  <label htmlFor="email">Correo electrónico</label>
  <input
    id="email"
    type="email"
    name="email"
    aria-required="true"
    aria-describedby="email-error"
  />
  <span id="email-error" role="alert">
    {errors.email}
  </span>
</div>;

{
  /* ❌ Incorrecto — placeholder no es suficiente como label */
}
<input type="email" placeholder="Correo electrónico" />;
```

**1.4.3 — Contraste de Color (Nivel AA)**

Relación de contraste mínima entre texto y fondo:

| Tipo de texto                       | Contraste mínimo |
| ----------------------------------- | ---------------- |
| Texto normal (< 18pt)               | 4.5:1            |
| Texto grande (≥ 18pt o ≥ 14pt bold) | 3:1              |
| Componentes UI (botones, inputs)    | 3:1              |

```tsx
{
  /* ✅ Verificar con herramienta: https://webaim.org/resources/contrastchecker/ */
}
{
  /* TailwindCSS — usar colores con contraste verificado */
}
<button className="bg-blue-700 text-white">
  {" "}
  {/* Contraste: 5.9:1 ✅ */}
  Iniciar sesión
</button>;

{
  /* ❌ Insuficiente contraste */
}
<button className="bg-blue-300 text-white">
  {" "}
  {/* Contraste: 1.8:1 ❌ */}
  Iniciar sesión
</button>;
```

---

### O — Operable (Operable)

Los componentes de la interfaz y la navegación deben ser operables.

**2.1.1 — Teclado (Nivel A)**

Toda funcionalidad debe ser accesible con teclado (sin mouse).

```tsx
{
  /* ✅ Los elementos interactivos reciben foco con Tab */
}
{
  /* El orden del foco debe ser lógico (top → bottom, left → right) */
}

<form>
  <input type="email" /> {/* Tab: 1 */}
  <input type="password" /> {/* Tab: 2 */}
  <button type="submit" /> {/* Tab: 3 */}
</form>;

{
  /* ✅ No suprimir el outline de foco */
}
{
  /* En TailwindCSS: usar focus:ring en lugar de outline-none */
}
<input className="focus:ring-2 focus:ring-blue-500 focus:outline-none" />;

{
  /* ❌ Nunca hacer esto — elimina la visibilidad del foco */
}
<input style={{ outline: "none" }} />;
```

**2.4.7 — Foco Visible (Nivel AA)**

El elemento que tiene el foco del teclado debe tener un indicador visual claro.

```tsx
{
  /* ✅ Estilo de foco consistente con TailwindCSS */
}
const inputClass = `
  border border-gray-300 rounded-md px-3 py-2
  focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent
`;
```

---

### U — Comprensible (Understandable)

La información y el funcionamiento de la UI deben ser comprensibles.

**3.1.1 — Idioma de la Página (Nivel A)**

El idioma principal de la página debe declararse en el HTML.

```html
<!-- index.html -->
<html lang="es"></html>
```

**3.3.1 — Identificación de Errores (Nivel A)**

Si un error de formulario se detecta automáticamente, el elemento erróneo debe ser identificado y el error descrito.

```tsx
{
  /* ✅ Correcto — error asociado al campo y anunciado */
}
function InputField({ id, label, error, ...props }) {
  return (
    <div>
      <label htmlFor={id}>{label}</label>
      <input
        id={id}
        aria-invalid={!!error} // Indica estado inválido al screen reader
        aria-describedby={error ? `${id}-error` : undefined}
        {...props}
      />
      {error && (
        <p
          id={`${id}-error`}
          role="alert"
          className="text-red-600 text-sm mt-1"
        >
          {error}
        </p>
      )}
    </div>
  );
}
```

**3.3.2 — Etiquetas o Instrucciones (Nivel A)**

Los campos de formulario deben tener etiquetas descriptivas.

```tsx
{
  /* ✅ Label clara + requisitos de contraseña visibles */
}
<div>
  <label htmlFor="password">
    Contraseña
    <span className="sr-only"> (requerida)</span>
  </label>
  <input
    id="password"
    type="password"
    aria-required="true"
    aria-describedby="password-hint"
  />
  <p id="password-hint" className="text-sm text-gray-500">
    Mínimo 8 caracteres, con al menos una mayúscula, minúscula y número.
  </p>
</div>;
```

---

### R — Robusto (Robust)

El contenido debe ser suficientemente robusto para ser interpretado por agentes de usuario, incluidas las tecnologías asistivas.

**4.1.2 — Nombre, Rol, Valor (Nivel A)**

Todos los componentes de la interfaz deben tener nombre, rol y valor accesibles.

```tsx
{
  /* ✅ Botón con estado de carga accesible */
}
<button
  type="submit"
  disabled={isLoading}
  aria-busy={isLoading}
  aria-label={
    isLoading ? "Iniciando sesión, por favor espera" : "Iniciar sesión"
  }
>
  {isLoading ? "Iniciando sesión..." : "Iniciar sesión"}
</button>;

{
  /* ✅ Modal con roles correctos */
}
<div
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
  aria-describedby="modal-description"
>
  <h2 id="modal-title">Confirmar acción</h2>
  <p id="modal-description">¿Estás seguro de que deseas cerrar sesión?</p>
</div>;
```

---

## ARIA — Accessible Rich Internet Applications

ARIA es un conjunto de atributos HTML que añaden semántica a elementos que no la tienen nativamente.

### Atributos ARIA más usados en formularios de auth

| Atributo           | Propósito                                        | Ejemplo                            |
| ------------------ | ------------------------------------------------ | ---------------------------------- |
| `aria-label`       | Nombre accesible para el elemento                | `aria-label="Mostrar contraseña"`  |
| `aria-labelledby`  | Referencia al elemento que sirve de label        | `aria-labelledby="titulo-form"`    |
| `aria-describedby` | Referencia al elemento con descripción adicional | `aria-describedby="password-hint"` |
| `aria-required`    | Indica campo obligatorio                         | `aria-required="true"`             |
| `aria-invalid`     | Indica estado de validación                      | `aria-invalid="true"`              |
| `aria-busy`        | Indica que el elemento está cargando             | `aria-busy="true"`                 |
| `aria-live`        | Anuncia cambios dinámicos                        | `aria-live="polite"`               |
| `aria-hidden`      | Oculta del árbol de accesibilidad                | `aria-hidden="true"` (icono dec.)  |
| `role="alert"`     | Anuncia el contenido inmediatamente              | Mensajes de error                  |
| `role="status"`    | Anuncia cambios de estado (no urgente)           | "Guardando..."                     |

### Anuncios dinámicos con `aria-live`

```tsx
{
  /* ✅ Notificar mensajes de éxito/error al screen reader */
}
{
  successMessage && (
    <div
      role="status"
      aria-live="polite" // Sin interrumpir al reader
      className="bg-green-50 text-green-800 p-3 rounded"
    >
      {successMessage}
    </div>
  );
}

{
  errorMessage && (
    <div
      role="alert"
      aria-live="assertive" // Interrumpe e informa inmediatamente
      className="bg-red-50 text-red-800 p-3 rounded"
    >
      {errorMessage}
    </div>
  );
}
```

---

## Checklist de Accesibilidad para Formularios de Auth

### LoginPage

- [ ] `<html lang="es">` en index.html
- [ ] `<title>` descriptivo: "Iniciar Sesión — NN Auth System"
- [ ] Todos los inputs tienen `<label>` con `htmlFor`/`id`
- [ ] El input de email tiene `type="email"`
- [ ] El input de contraseña tiene `type="password"`
- [ ] Botón de mostrar/ocultar contraseña tiene `aria-label`
- [ ] Mensajes de error tienen `role="alert"` y están asociados con `aria-describedby`
- [ ] El botón de submit tiene estado `aria-busy` durante carga
- [ ] Se puede completar todo el formulario solo con teclado
- [ ] El foco es visible en todos los elementos interactivos

### RegisterPage

- [ ] Requisitos de contraseña visibles antes de typing (`aria-describedby`)
- [ ] Confirmación de contraseña tiene label propio
- [ ] Errores de validación anuncian el campo específico con problema
- [ ] Mensaje de éxito tiene `role="status"`

### ForgotPasswordPage / ResetPasswordPage

- [ ] Instrucciones claras antes del formulario
- [ ] Mensaje de respuesta genérico para prevenir enumeración de usuarios
- [ ] Timer de expiración de token visible y comprensible

---

## Dark Mode Accesible

El toggle de tema debe ser accesible:

```tsx
{
  /* ✅ Toggle de tema accesible */
}
<button
  onClick={toggleTheme}
  aria-label={isDark ? "Cambiar a modo claro" : "Cambiar a modo oscuro"}
  aria-pressed={isDark} // Estado actual del toggle
>
  {isDark ? <SunIcon aria-hidden="true" /> : <MoonIcon aria-hidden="true" />}
</button>;
```

Y los colores de dark mode también deben cumplir con el contraste WCAG AA:

```css
/* Variables de color con contraste verificado */
:root {
  --text-primary: #111827; /* sobre blanco: 16:1 ✅ */
  --bg-primary: #ffffff;
}

.dark {
  --text-primary: #f9fafb; /* sobre #111827: 16:1 ✅ */
  --bg-primary: #111827;
}
```

---

## Herramientas de Prueba

| Herramienta                         | Propósito                                     |
| ----------------------------------- | --------------------------------------------- |
| **axe DevTools** (extensión Chrome) | Detecta 57% de los problemas de accesibilidad |
| **NVDA** (Windows, gratuito)        | Lector de pantalla para probar manualmente    |
| **VoiceOver** (macOS/iOS, nativo)   | Lector de pantalla de Apple                   |
| **Lighthouse** (Chrome DevTools)    | Audit de accesibilidad automático             |
| **WebAIM Contrast Checker**         | Verificar ratio de contraste de colores       |
| **keyboard-only navigation**        | Navegar sin mouse para probar teclado         |

```bash
# En el proyecto, se puede instalar axe-core para tests automáticos
cd fe && pnpm add -D @axe-core/react

# En desarrollo, añadir en main.tsx para detectar issues en tiempo real
if (import.meta.env.DEV) {
  import('@axe-core/react').then(axe => axe.default(React, ReactDOM, 1000));
}
```

---

> La accesibilidad no es una característica adicional — es un derecho. Un sistema de autenticación inaccesible excluye a usuarios que tienen tanto derecho de usar la aplicación como cualquier otra persona.
