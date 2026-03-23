# RNF-003 — Usabilidad

## Identificación

| Campo         | Valor           |
| ------------- | --------------- |
| **ID**        | RNF-003         |
| **Nombre**    | Usabilidad      |
| **Categoría** | Usabilidad / UX |
| **Prioridad** | Alta            |

---

## Descripción

El sistema debe ser intuitivo, fácil de usar y proporcionar retroalimentación clara al usuario en todo momento. La interfaz debe adaptarse a dispositivos móviles y de escritorio con igual calidad visual.

---

## Principios UX Aplicados

| Principio              | Implementación en el sistema                                        |
| ---------------------- | ------------------------------------------------------------------- |
| Feedback inmediato     | Mensajes de error inline junto al campo afectado                    |
| Prevención de errores  | Validación en tiempo real antes de enviar formularios               |
| Visibilidad del estado | Spinner/loader durante operaciones asíncronas                       |
| Consistencia           | Mismos componentes (InputField, Button, Alert) en todas las páginas |
| Control del usuario    | Posibilidad de ver/ocultar contraseña con ícono toggle              |

---

## Requisitos de Mensajes

### Mensajes de Error

- Los errores deben estar en **español** y ser descriptivos
- Deben aparecer cerca del campo que los causó
- No exponer información sensible (ej. no indicar si el email existe)

| Error                    | Mensaje al usuario                                           |
| ------------------------ | ------------------------------------------------------------ |
| Email inválido           | "Ingresa un email válido (ej. usuario@ejemplo.com)"          |
| Contraseña muy corta     | "La contraseña debe tener al menos 8 caracteres"             |
| Credenciales incorrectas | "Email o contraseña incorrectos"                             |
| Email no verificado      | "Debes verificar tu email antes de iniciar sesión"           |
| Token expirado           | "Este enlace ha expirado. Solicita uno nuevo."               |
| Error de red             | "Error de conexión. Verifica tu internet e intenta de nuevo" |
| Error del servidor       | "Algo salió mal. Intenta de nuevo en unos minutos."          |

### Estados de Carga

- Botones muestran spinner y texto "Procesando..." durante la petición
- Botones se deshabilitan durante la petición (evitar doble envío)

---

## Responsividad

| Breakpoint TailwindCSS | Ancho mínimo | Comportamiento esperado           |
| ---------------------- | ------------ | --------------------------------- |
| `sm`                   | 640px        | Formularios en columna simple     |
| `md`                   | 768px        | Layout de dos columnas            |
| `lg`                   | 1024px       | Layout completo con sidebar       |
| Mínimo absoluto        | 320px        | Formularios fluidos, sin overflow |

---

## Diseño Visual Obligatorio

| Aspecto           | Regla                                           |
| ----------------- | ----------------------------------------------- |
| Temas             | Dark y Light con toggle visible                 |
| Tipografía        | Fuentes sans-serif (Inter, system-ui)           |
| Colores           | Sólidos y planos — cero degradados (`gradient`) |
| Botones de acción | Alineados a la derecha (`justify-end`)          |
| Espaciado         | Escala consistente de TailwindCSS               |

---

## Criterio de Cumplimiento

- Todos los formularios tienen mensajes de error descriptivos en español
- La interfaz es completamente funcional en pantalla de 375px (iPhone SE)
- Los botones están deshabilitados durante la petición (sin doble envío)
- Dark mode y light mode operan sin recargar la página
