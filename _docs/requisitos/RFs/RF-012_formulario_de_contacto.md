# RF-012 — Formulario de Contacto

## Identificación

| Campo              | Valor                    |
| ------------------ | ------------------------ |
| **ID**             | RF-012                   |
| **Nombre**         | Formulario de contacto   |
| **Módulo**         | Contenido / Comunicación |
| **HU relacionada** | HU-011                   |
| **Prioridad**      | Baja                     |

---

## Descripción

El sistema debe incluir una página de contacto con un formulario que permita a cualquier usuario (autenticado o no) enviar una consulta o mensaje. El formulario es público y no requiere autenticación.

---

## Campos del Formulario

| Campo     | Tipo     | Obligatorio | Validación              |
| --------- | -------- | ----------- | ----------------------- |
| `name`    | text     | Sí          | Mínimo 2 caracteres     |
| `email`   | email    | Sí          | Formato de email válido |
| `subject` | text     | Sí          | Mínimo 5 caracteres     |
| `message` | textarea | Sí          | Mínimo 20 caracteres    |

---

## Proceso

1. Usuario navega a `/contact`
2. Completa el formulario y hace clic en "Enviar"
3. El frontend valida los campos localmente antes de enviar
4. Se muestra mensaje de éxito `"Mensaje enviado. Te responderemos pronto."`

> **Nota de implementación:** Para este proyecto educativo, el envío puede ser simulado (mostrar mensaje de éxito sin backend) o integrarse con un servicio externo como Formspree (`https://formspree.io`). No se requiere endpoint propio en el backend.

---

## Integración con Formspree (opcional)

```
// Configuración en .env
VITE_FORMSPREE_ENDPOINT=https://formspree.io/f/tu-id

// En ContactPage.tsx — acción del formulario
<form action={import.meta.env.VITE_FORMSPREE_ENDPOINT} method="POST">
```

---

## Reglas de Accesibilidad

- Todos los campos deben tener `<label>` asociado
- Mensajes de error deben ser descriptivos y asociados con `aria-describedby`
- El formulario debe ser operablee completamente con teclado

---

## Reglas de Negocio

| ID     | Regla                                                    |
| ------ | -------------------------------------------------------- |
| RN-038 | El formulario de contacto es público (sin autenticación) |
| RN-039 | La validación del lado del cliente es obligatoria        |
| RN-040 | Mostrar mensaje de confirmación tras envío exitoso       |
