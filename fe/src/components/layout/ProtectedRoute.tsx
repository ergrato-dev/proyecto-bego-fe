/**
 * Archivo: ProtectedRoute.tsx
 * Descripción: Componente de ruta protegida para React Router.
 * ¿Para qué? Redirigir usuarios no autenticados al login antes de
 *            mostrar páginas que requieren sesión activa.
 * ¿Impacto? Sin este componente, cualquier usuario podría acceder al
 *            dashboard o cambio de contraseña sin estar autenticado.
 */

import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

/**
 * ProtectedRoute verifica si el usuario está autenticado.
 * ¿Para qué? Actúa como guardia de ruta — si no hay sesión activa,
 *            redirige al login de forma transparente.
 * ¿Impacto? 'replace' en el Navigate reemplaza la entrada del historial
 *            para que el botón "atrás" no regrese a la ruta protegida.
 *
 * Uso en el router (App.tsx):
 *   <Route element={<ProtectedRoute />}>
 *     <Route path="/dashboard" element={<DashboardPage />} />
 *   </Route>
 */
export default function ProtectedRoute() {
  const { isAuthenticated } = useAuth()

  // ¿Qué? Si no está autenticado, redirige al login.
  // ¿Para qué? Outlet renderiza la ruta hija solo cuando hay sesión.
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}
