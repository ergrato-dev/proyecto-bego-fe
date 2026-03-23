# RF-011 — Páginas Legales

## Identificación

| Campo              | Valor              |
| ------------------ | ------------------ |
| **ID**             | RF-011             |
| **Nombre**         | Páginas legales    |
| **Módulo**         | Contenido Estático |
| **HU relacionada** | HU-010             |
| **Prioridad**      | Baja               |

---

## Descripción

El sistema debe incluir páginas estáticas de contenido legal: Política de Privacidad y Términos de Servicio. Son páginas informativas, públicas y sin llamadas a la API.

---

## Páginas Requeridas

| Página                 | Ruta       | Componente        |
| ---------------------- | ---------- | ----------------- |
| Política de Privacidad | `/privacy` | `PrivacyPage.tsx` |
| Términos de Servicio   | `/terms`   | `TermsPage.tsx`   |

---

## Contenido Mínimo

### Política de Privacidad (`/privacy`)

- Qué datos se recopilan (email, nombre, contraseña hasheada)
- Cómo se usan los datos (autenticación, no se comparten)
- Cómo protegemos los datos (bcrypt, JWT, HTTPS en producción)
- Derechos del usuario (eliminar cuenta, etc.)
- Fecha de última actualización

### Términos de Servicio (`/terms`)

- Descripción del servicio
- Uso aceptable del sistema
- Responsabilidades del usuario
- Limitación de responsabilidad (proyecto educativo)
- Cambios a los términos
- Fecha de última actualización

---

## Proceso

1. Usuario navega a `/privacy` o `/terms`
2. React Router renderiza el componente correspondiente
3. No se realiza ninguna llamada a la API
4. El contenido es estático (puede ser JSX con texto formateado)

---

## Reglas de Negocio

| ID     | Regla                                                       |
| ------ | ----------------------------------------------------------- |
| RN-035 | Las páginas legales son completamente públicas              |
| RN-036 | Deben incluir fecha de "última actualización"               |
| RN-037 | Deben ser accesibles (contraste, estructura semántica HTML) |
