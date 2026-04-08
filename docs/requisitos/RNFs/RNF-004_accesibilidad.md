# RNF-004 — Accesibilidad

## Identificación

| Campo         | Valor                |
| ------------- | -------------------- |
| **ID**        | RNF-004              |
| **Nombre**    | Accesibilidad (a11y) |
| **Categoría** | Accesibilidad        |
| **Prioridad** | Alta                 |

---

## Descripción

El sistema debe cumplir con el estándar WCAG 2.1 nivel AA para garantizar que sea usable por personas con discapacidades visuales, motoras o cognitivas.

---

## Estándar Aplicable

| Estándar | Nivel | Descripción                              |
| -------- | ----- | ---------------------------------------- |
| WCAG 2.1 | AA    | Web Content Accessibility Guidelines 2.1 |

---

## Principios POUR

| Principio        | Requisito en este sistema                                   |
| ---------------- | ----------------------------------------------------------- |
| **Perceptible**  | Contraste mínimo 4.5:1 texto normal, 3:1 texto grande       |
| **Operable**     | Toda función accesible con teclado (Tab, Enter, Escape)     |
| **Comprensible** | Etiquetas claras, mensajes de error descriptivos en español |
| **Robusto**      | HTML semántico compatible con lectores de pantalla          |

---

## Requisitos por Área

### Formularios

| Elemento          | Requisito                                              |
| ----------------- | ------------------------------------------------------ |
| `<input>`         | Siempre asociado a un `<label>` con `htmlFor`/`id`     |
| Mensajes de error | Vinculados con `aria-describedby` al campo afectado    |
| Campos inválidos  | Marcados con `aria-invalid="true"`                     |
| Campos requeridos | Marcados con `aria-required="true"` o `required`       |
| Botón submit      | Texto descriptivo ("Iniciar sesión", no solo "Submit") |

### Navegación por Teclado

| Acción            | Comportamiento esperado                                    |
| ----------------- | ---------------------------------------------------------- |
| `Tab`             | Navega a través de todos los elementos interactivos        |
| `Shift + Tab`     | Navegación inversa                                         |
| `Enter` en botón  | Activa el botón                                            |
| `Escape` en modal | Cierra el modal/diálogo                                    |
| Foco visible      | Ring de foco visible en todos los elementos (no se oculta) |

### Contraste de Color

| Contexto              | Mínimo WCAG AA |
| --------------------- | -------------- |
| Texto normal (< 18px) | 4.5 : 1        |
| Texto grande (≥ 18px) | 3.0 : 1        |
| Componentes UI        | 3.0 : 1        |

El contraste debe cumplirse tanto en modo claro como en modo oscuro.

### HTML Semántico y ARIA

```html
<!-- Ejemplo: formulario de login accesible -->
<main>
  <section aria-labelledby="login-title">
    <h1 id="login-title">Iniciar sesión</h1>
    <form novalidate>
      <div role="alert" aria-live="polite">
        <!-- Mensajes de error globales -->
      </div>
      <div>
        <label htmlFor="email">Correo electrónico</label>
        <input
          id="email"
          type="email"
          aria-required="true"
          aria-describedby="email-error"
        />
        <span id="email-error" role="alert">Ingresa un email válido</span>
      </div>
    </form>
  </section>
</main>
```

---

## Herramientas de Verificación

| Herramienta           | Tipo               | Cómo usar                                     |
| --------------------- | ------------------ | --------------------------------------------- |
| axe DevTools          | Extensión Chrome   | Auditoría automática de accesibilidad         |
| Google Lighthouse     | Chrome DevTools    | Pestaña "Accessibility" → score objetivo ≥ 90 |
| NVDA (Windows)        | Lector de pantalla | Probar navegación real con lector             |
| VoiceOver (macOS/iOS) | Lector de pantalla | Probar en Safari                              |
| Keyboard-only test    | Manual             | Navegar toda la app sin ratón                 |

---

## Criterio de Cumplimiento

- Lighthouse Accessibility Score ≥ 90 en todas las páginas
- Sin errores críticos en axe DevTools
- Todos los formularios navegables y usables solo con teclado
- Contraste verificado en modo claro y oscuro con herramienta de contraste
