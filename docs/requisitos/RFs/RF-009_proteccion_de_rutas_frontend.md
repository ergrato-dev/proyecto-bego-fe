# RF-009 — Protección de Rutas Frontend

## Identificación

| Campo              | Valor                             |
| ------------------ | --------------------------------- |
| **ID**             | RF-009                            |
| **Nombre**         | Protección de rutas (frontend)    |
| **Módulo**         | Enrutamiento / Seguridad frontend |
| **HU relacionada** | HU-008                            |
| **Prioridad**      | Alta                              |

---

## Descripción

El sistema debe impedir que usuarios no autenticados accedan a rutas protegidas (como `/dashboard` o `/change-password`). El componente `ProtectedRoute` actúa como guardia de ruta, verificando el estado de autenticación antes de renderizar el contenido solicitado.

---

## Proceso

1. El usuario intenta navegar a una ruta protegida (ej. `/dashboard`)
2. `ProtectedRoute` consulta `AuthContext` obteniendo `{ user, isLoading }`
3. Si `isLoading === true` → renderizar componente de carga (spinner) sin redirigir
4. Si `isLoading === false && user === null` → redirigir a `/login` con `<Navigate to="/login" replace />`
5. Si `isLoading === false && user !== null` → renderizar el componente hijo de la ruta

---

## Diagrama de Flujo

```
Navegar a /dashboard
        │
        ▼
  AuthContext.isLoading?
        │
   SÍ ─┤─ NO
        │         │
   Spinner    user existe?
               │
          SÍ ─┤─ NO
               │         │
         Renderizar   Redirigir a
         Dashboard    /login
```

---

## Rutas del Sistema

| Ruta               | Tipo      | Componente         |
| ------------------ | --------- | ------------------ |
| `/`                | Pública   | LandingPage        |
| `/login`           | Pública   | LoginPage          |
| `/register`        | Pública   | RegisterPage       |
| `/forgot-password` | Pública   | ForgotPasswordPage |
| `/reset-password`  | Pública   | ResetPasswordPage  |
| `/privacy`         | Pública   | PrivacyPage        |
| `/terms`           | Pública   | TermsPage          |
| `/contact`         | Pública   | ContactPage        |
| `/dashboard`       | Protegida | DashboardPage      |
| `/change-password` | Protegida | ChangePasswordPage |

---

## Reglas de Negocio

| ID     | Regla                                                                       |
| ------ | --------------------------------------------------------------------------- |
| RN-028 | Las rutas protegidas nunca deben mostrar contenido antes de verificar auth  |
| RN-029 | Mientras carga el estado de auth, mostrar spinner (no flash de contenido)   |
| RN-030 | La redirección al login usa `replace` para no dejar historial de navegación |
| RN-031 | Un usuario autenticado que visita `/login` debe redirigirse al `/dashboard` |
