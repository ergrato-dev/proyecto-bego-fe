---
description: "Crea una nueva página React completa: componente, registro de ruta en App.tsx y test. Usar cuando se necesite agregar una nueva vista a la aplicación."
name: "Nueva página React"
argument-hint: "Describe la página: nombre, ruta URL, propósito, si requiere auth, datos que muestra o formularios que contiene"
agent: "agent"
---

# Nueva página React — NN Auth System

Crea una página React completa siguiendo las convenciones del proyecto.

## Convenciones obligatorias

- **Idioma del código**: inglés (nombres de componentes, props, variables, funciones, rutas)
- **Idioma de comentarios y JSDoc**: español
- **Comentarios pedagógicos**: ¿Qué? ¿Para qué? ¿Impacto? en cada bloque significativo
- **TypeScript strict**: todos los props con tipos explícitos, sin `any`
- **TailwindCSS**: sin CSS custom adicional, usar clases de Tailwind
- **Diseño**: dark/light mode, sin degradados (`gradient`), tipografía sans-serif, botones de acción a la derecha
- **Textos en pantalla**: usar siempre `t('clave')` de `useTranslation()` — nunca texto hardcodeado

## Lo que debes generar

### 1. Componente de página (`fe/src/pages/NombreDePagina.tsx`)

- Componente funcional con hooks
- Cabecera de archivo JSDoc (¿Qué? ¿Para qué? ¿Impacto?)
- Importar `useTranslation` de `react-i18next` y usar `t('clave')` para todo texto visible
- Usar componentes UI existentes del proyecto:
  - [fe/src/components/ui/Button.tsx](../../../fe/src/components/ui/Button.tsx)
  - [fe/src/components/ui/InputField.tsx](../../../fe/src/components/ui/InputField.tsx)
  - [fe/src/components/ui/Alert.tsx](../../../fe/src/components/ui/Alert.tsx)
- Si requiere auth: envolver la ruta con `ProtectedRoute` en App.tsx
- Loading states para operaciones asíncronas; manejo de errores con `Alert`
- Página de referencia: [fe/src/pages/LoginPage.tsx](../../../fe/src/pages/LoginPage.tsx)

### 2. Claves de traducción en ambos archivos de locales

Agregar **todas** las claves nuevas en los dos archivos:

- [fe/src/locales/es/translation.json](../../../fe/src/locales/es/translation.json)
- [fe/src/locales/en/translation.json](../../../fe/src/locales/en/translation.json)

Seguir la convención de claves del proyecto:

```json
// es/translation.json — agregar bajo el namespace correspondiente
{
  "nombrePagina": {
    "title": "Título de la página",
    "labelXxx": "...",
    "errorXxx": "...",
    "successXxx": "..."
  }
}
```

### 3. Ruta en App.tsx ([fe/src/App.tsx](../../../fe/src/App.tsx))

- Agregar el `<Route>` correspondiente
- Rutas públicas: `<Route path="/ruta" element={<NombrePagina />} />`
- Rutas protegidas: `<Route path="/ruta" element={<ProtectedRoute><NombrePagina /></ProtectedRoute>} />`

### 4. Test (`fe/src/__tests__/pages/NombreDePagina.test.tsx`)

Casos mínimos:

- Renderiza el título principal
- Si tiene formulario: muestra errores de validación ante inputs inválidos
- Si requiere auth: redirige a login sin token
- Si hace llamada a API: muestra estado de éxito y error

Usar mocks de `vi.mock` para `axios`/módulos de API y el `AuthContext`.
Test de referencia: [fe/src/__tests__/pages/DashboardPage.test.tsx](../../../fe/src/__tests__/pages/DashboardPage.test.tsx)

## Accesibilidad (WCAG AA)

- `<label>` con `htmlFor` en todos los inputs
- `role="alert"` en mensajes de error
- Contraste suficiente en modo claro y oscuro
- Botones con texto descriptivo (no solo iconos)

## Descripción de la página a crear

$input
