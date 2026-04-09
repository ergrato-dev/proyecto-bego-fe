<!--
  ¿Qué? Documentación del concepto de internacionalización (i18n) y su implementación en el proyecto.
  ¿Para qué? Explicar qué es i18n, por qué importa y cómo se aplica en React con react-i18next.
  ¿Impacto? Sin esta guía, cada desarrollador implementaría la traducción de forma distinta,
             generando inconsistencias en textos, claves y estructura de archivos.
-->

# Internacionalización (i18n) — NN Auth System

## 1. ¿Qué es i18n y l10n?

| Término  | Nombre completo      | Significado                                                                |
| -------- | -------------------- | -------------------------------------------------------------------------- |
| **i18n** | Internationalization | Preparar el código para soportar múltiples idiomas (sin texto hardcodeado) |
| **l10n** | Localization         | Proveer las traducciones y formatos culturales específicos por idioma      |

> El nombre "i18n" viene de que entre la "i" y la "n" de _internationalization_ hay 18 letras.

**Idiomas soportados en este proyecto:**

| Código | Idioma  | Es el idioma por defecto |
| ------ | ------- | ------------------------ |
| `es`   | Español | ✅ Sí                    |
| `en`   | Inglés  | No                       |

---

## 2. ¿Por qué i18n en un sistema de autenticación?

Aunque el sistema de auth parece "solo para uso interno", hay razones concretas:

- **Audiencia mixta:** usuarios hispanohablantes y anglohablantes pueden usar la misma app.
- **Mensajes de error y validación:** deben mostrarse en el idioma del usuario para comprensión inmediata.
- **Emails transaccionales:** verificación de email, recuperación de contraseña — deben llegar en el idioma del usuario.
- **Accesibilidad:** los lectores de pantalla leen el idioma indicado en el atributo `lang` del HTML.
- **Escalabilidad:** agregar un idioma después es trivial si el código ya está preparado; rehacerlo todo después es costoso.

---

## 3. Librería: `react-i18next`

| Herramienta                        | Versión pinada | Rol                                             |
| ---------------------------------- | -------------- | ----------------------------------------------- |
| `i18next`                          | **25.2.1**     | Motor principal de i18n — gestiona traducciones |
| `react-i18next`                    | **15.5.1**     | Integración de i18next con React (hooks, HOCs)  |
| `i18next-browser-languagedetector` | **8.0.5**      | Detecta el idioma del navegador automáticamente |

**¿Por qué react-i18next?**

- Es el estándar de facto en el ecosistema React.
- API simple: hook `useTranslation()` en cualquier componente.
- Soporta interpolación, plurales, namespaces y contexto.
- Compatible con TypeScript — se puede tipar el namespace de traducciones.
- Mantenimiento activo y amplia comunidad.

---

## 4. Estructura de archivos de traducción

```
fe/src/
└── locales/
    ├── es/
    │   └── translation.json    # Traducciones en español (idioma por defecto)
    └── en/
        └── translation.json    # Traducciones en inglés
```

### 4.1 Ejemplo de archivo de traducción

```json
// fe/src/locales/es/translation.json
{
  "common": {
    "loading": "Cargando...",
    "save": "Guardar",
    "cancel": "Cancelar",
    "back": "Volver",
    "required": "Este campo es obligatorio"
  },
  "auth": {
    "login": {
      "title": "Iniciar sesión",
      "email": "Correo electrónico",
      "password": "Contraseña",
      "submit": "Entrar",
      "forgotPassword": "¿Olvidaste tu contraseña?",
      "noAccount": "¿No tienes cuenta?",
      "register": "Regístrate",
      "errorInvalid": "Correo o contraseña incorrectos",
      "errorInactive": "Tu cuenta está desactivada"
    },
    "register": {
      "title": "Crear cuenta",
      "fullName": "Nombre completo",
      "email": "Correo electrónico",
      "password": "Contraseña",
      "submit": "Registrarse",
      "haveAccount": "¿Ya tienes cuenta?",
      "login": "Inicia sesión",
      "errorEmailExists": "Este correo ya está registrado",
      "successMessage": "Cuenta creada. Revisa tu correo para verificar tu email."
    },
    "changePassword": {
      "title": "Cambiar contraseña",
      "current": "Contraseña actual",
      "new": "Nueva contraseña",
      "submit": "Cambiar contraseña",
      "successMessage": "Contraseña actualizada exitosamente",
      "errorWrong": "La contraseña actual es incorrecta"
    },
    "forgotPassword": {
      "title": "Recuperar contraseña",
      "email": "Correo electrónico",
      "submit": "Enviar instrucciones",
      "successMessage": "Si el correo existe, recibirás instrucciones en breve"
    },
    "resetPassword": {
      "title": "Nueva contraseña",
      "new": "Nueva contraseña",
      "submit": "Guardar contraseña",
      "successMessage": "Contraseña restablecida. Ya puedes iniciar sesión.",
      "errorExpired": "El enlace expiró o ya fue usado"
    }
  },
  "dashboard": {
    "title": "Panel principal",
    "welcome": "Bienvenido, {{name}}",
    "email": "Correo electrónico",
    "logout": "Cerrar sesión"
  },
  "nav": {
    "home": "Inicio",
    "login": "Iniciar sesión",
    "register": "Registrarse"
  },
  "errors": {
    "generic": "Ocurrió un error. Intenta de nuevo.",
    "network": "Error de conexión. Verifica tu internet.",
    "unauthorized": "Tu sesión expiró. Inicia sesión de nuevo."
  }
}
```

```json
// fe/src/locales/en/translation.json
{
  "common": {
    "loading": "Loading...",
    "save": "Save",
    "cancel": "Cancel",
    "back": "Back",
    "required": "This field is required"
  },
  "auth": {
    "login": {
      "title": "Sign in",
      "email": "Email",
      "password": "Password",
      "submit": "Sign in",
      "forgotPassword": "Forgot your password?",
      "noAccount": "Don't have an account?",
      "register": "Sign up",
      "errorInvalid": "Invalid email or password",
      "errorInactive": "Your account is disabled"
    },
    "register": {
      "title": "Create account",
      "fullName": "Full name",
      "email": "Email",
      "password": "Password",
      "submit": "Sign up",
      "haveAccount": "Already have an account?",
      "login": "Sign in",
      "errorEmailExists": "This email is already registered",
      "successMessage": "Account created. Check your email to verify your address."
    },
    "changePassword": {
      "title": "Change password",
      "current": "Current password",
      "new": "New password",
      "submit": "Change password",
      "successMessage": "Password updated successfully",
      "errorWrong": "Current password is incorrect"
    },
    "forgotPassword": {
      "title": "Recover password",
      "email": "Email",
      "submit": "Send instructions",
      "successMessage": "If the email exists, you will receive instructions shortly"
    },
    "resetPassword": {
      "title": "New password",
      "new": "New password",
      "submit": "Save password",
      "successMessage": "Password reset. You can now sign in.",
      "errorExpired": "The link has expired or was already used"
    }
  },
  "dashboard": {
    "title": "Dashboard",
    "welcome": "Welcome, {{name}}",
    "email": "Email",
    "logout": "Sign out"
  },
  "nav": {
    "home": "Home",
    "login": "Sign in",
    "register": "Sign up"
  },
  "errors": {
    "generic": "An error occurred. Please try again.",
    "network": "Connection error. Check your internet.",
    "unauthorized": "Your session expired. Sign in again."
  }
}
```

---

## 5. Configuración de i18next (`fe/src/i18n.ts`)

```typescript
/**
 * Archivo: i18n.ts
 * Descripción: Configuración central de i18next para el proyecto.
 * ¿Para qué? Inicializar el motor de traducciones con los recursos es/en
 *            y configurar la detección automática del idioma del navegador.
 * ¿Impacto? Sin esta configuración, el hook useTranslation() no funcionará
 *           en ningún componente de la aplicación.
 */
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import esTranslation from "./locales/es/translation.json";
import enTranslation from "./locales/en/translation.json";

i18n
  // ¿Qué? LanguageDetector lee navigator.language del navegador.
  // ¿Para qué? Mostrar automáticamente el idioma del usuario sin que tenga que elegirlo.
  // ¿Impacto? Sin esto, siempre mostraría el idioma por defecto (es).
  .use(LanguageDetector)
  // ¿Qué? initReactI18next conecta i18next con React mediante contexto.
  // ¿Para qué? Permite usar useTranslation() en cualquier componente.
  .use(initReactI18next)
  .init({
    resources: {
      es: { translation: esTranslation },
      en: { translation: enTranslation },
    },
    fallbackLng: "es", // Idioma de respaldo si no se detecta el del navegador
    supportedLngs: ["es", "en"],
    interpolation: {
      escapeValue: false, // React ya escapa los valores — evita doble escape
    },
    detection: {
      order: ["localStorage", "navigator"],
      caches: ["localStorage"], // Persistir la selección del usuario
    },
  });

export default i18n;
```

Este archivo se importa una sola vez en `fe/src/main.tsx`:

```typescript
// fe/src/main.tsx
import './i18n'; // debe ir antes de importar App
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
```

---

## 6. Uso en componentes — `useTranslation()`

```typescript
import { useTranslation } from 'react-i18next';

export function LoginPage() {
  // ¿Qué? Hook que provee la función t() para traducir claves.
  // ¿Para qué? Acceder a las traducciones del archivo translation.json
  //            en el idioma activo del usuario.
  // ¿Impacto? Sin este hook, los textos serían hardcodeados en un solo idioma.
  const { t } = useTranslation();

  return (
    <h1>{t('auth.login.title')}</h1>   // → "Iniciar sesión" o "Sign in"
  );
}
```

### 6.1 Interpolación de variables

Para textos dinámicos usar `{{variable}}` en los archivos JSON:

```typescript
// En translation.json: "welcome": "Bienvenido, {{name}}"
t("dashboard.welcome", { name: user.fullName });
// → "Bienvenido, María García"
```

### 6.2 Selector de idioma

```typescript
import { useTranslation } from 'react-i18next';

export function LanguageSwitcher() {
  // ¿Qué? Componente que permite al usuario cambiar entre es/en.
  // ¿Para qué? Dar control al usuario sobre el idioma de la interfaz.
  // ¿Impacto? Sin esto, el idioma se detecta automáticamente pero no se puede cambiar.
  const { i18n } = useTranslation();

  const toggleLanguage = () => {
    const next = i18n.language === 'es' ? 'en' : 'es';
    i18n.changeLanguage(next);
  };

  return (
    <button onClick={toggleLanguage} aria-label="Cambiar idioma">
      {i18n.language === 'es' ? 'EN' : 'ES'}
    </button>
  );
}
```

---

## 7. Estructura de claves — Convenciones

| Convención                          | Ejemplo                          | Descripción                          |
| ----------------------------------- | -------------------------------- | ------------------------------------ |
| Namespaces por sección              | `auth.login.title`               | `sección.subsección.clave`           |
| `camelCase` en claves               | `forgotPassword`, `errorInvalid` | Nunca guión bajo ni guión            |
| Prefijo `error` en errores          | `errorInvalid`, `errorExpired`   | Distingue errores de labels normales |
| Prefijo `success` en confirmaciones | `successMessage`                 | Distingue confirmaciones             |
| `common.*` para reutilizables       | `common.loading`, `common.save`  | Textos compartidos entre páginas     |

**Regla de oro:** si un texto es igual en dos páginas, va en `common`. Si es específico de una sección, va en su namespace.

---

## 8. Atributo `lang` en el HTML

Para que los lectores de pantalla y motores de búsqueda conozcan el idioma activo, actualizar `index.html` dinámicamente:

```typescript
// En fe/src/i18n.ts, después de init():
i18n.on("languageChanged", (lng) => {
  document.documentElement.lang = lng; // <html lang="es"> o <html lang="en">
});
```

---

## 9. Instalación

```bash
cd fe
pnpm add i18next@25.2.1 react-i18next@15.5.1 i18next-browser-languagedetector@8.0.5
```

---

## 10. Checklist de i18n antes de commit

- [ ] ¿Ningún texto visible al usuario está hardcodeado en el componente?
- [ ] ¿Cada clave nueva existe en **ambos** archivos (`es/translation.json` y `en/translation.json`)?
- [ ] ¿Las claves nuevas siguen la convención `namespace.subsección.clave` en `camelCase`?
- [ ] ¿Los textos de variables usan `{{variable}}` en lugar de concatenación de strings?
- [ ] ¿El atributo `lang` del `<html>` se actualiza al cambiar idioma?
