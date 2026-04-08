# RF-007 — Cierre de Sesión

## Identificación

| Campo              | Valor                     |
| ------------------ | ------------------------- |
| **ID**             | RF-007                    |
| **Nombre**         | Cierre de sesión (logout) |
| **Módulo**         | Autenticación             |
| **HU relacionada** | HU-006                    |
| **Prioridad**      | Alta                      |

---

## Descripción

El sistema debe permitir al usuario cerrar su sesión activa. Dado que la autenticación usa JWT stateless, el logout se implementa en el cliente eliminando los tokens almacenados. No existe un endpoint de logout en el backend.

---

## Proceso (Frontend)

1. El usuario hace clic en "Cerrar sesión"
2. El `AuthContext` ejecuta la función `logout()`
3. Se eliminan `access_token` y `refresh_token` del almacenamiento del cliente (memoria / `localStorage`)
4. Se limpia el estado del contexto: `setUser(null)`
5. Se redirige al usuario a `/login` con `navigate('/login', { replace: true })`

---

## Justificación Técnica

| Pregunta                            | Respuesta                                                |
| ----------------------------------- | -------------------------------------------------------- |
| ¿Por qué no hay endpoint de logout? | JWT es stateless; el backend no guarda sesiones activas  |
| ¿Qué invalida los tokens?           | Su TTL (15 min access, 7 días refresh)                   |
| ¿Es seguro?                         | Sí, porque el access token expira rápido                 |
| ¿Cómo mitigar tokens comprometidos? | Reducir TTL del access token o usar lista negra en Redis |

---

## Salidas

| Acción         | Resultado                                        |
| -------------- | ------------------------------------------------ |
| Logout exitoso | Tokens eliminados, usuario redirigido a `/login` |

---

## Reglas de Negocio

| ID     | Regla                                                                  |
| ------ | ---------------------------------------------------------------------- |
| RN-022 | El logout no requiere llamada al backend                               |
| RN-023 | Tras el logout, toda ruta protegida debe redirigir al login            |
| RN-024 | El estado del AuthContext debe limpiarse completamente al hacer logout |
