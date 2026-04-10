/**
 * Archivo: App.tsx
 * Descripción: Componente raíz — define el árbol de providers y el router.
 * ¿Para qué? Envolver toda la aplicación con ThemeProvider y AuthProvider
 *            garantiza que cualquier componente puede acceder al tema y
 *            al estado de autenticación sin prop drilling.
 * ¿Impacto? El orden de providers importa: ThemeProvider va primero porque
 *            no depende de nada; AuthProvider después porque ThemeToggle
 *            (dentro de AuthLayout) necesita acceder a useTheme.
 */

import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { ThemeProvider } from './context/ThemeContext'
import ProtectedRoute from './components/layout/ProtectedRoute'
import { LandingPage } from './pages/LandingPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import DashboardPage from './pages/DashboardPage'
import ChangePasswordPage from './pages/ChangePasswordPage'
import ForgotPasswordPage from './pages/ForgotPasswordPage'
import ResetPasswordPage from './pages/ResetPasswordPage'
import VerifyEmailPage from './pages/VerifyEmailPage'
import { TerminosDeUsoPage } from './pages/TerminosDeUsoPage'
import { PoliticaPrivacidadPage } from './pages/PoliticaPrivacidadPage'
import { PoliticaCookiesPage } from './pages/PoliticaCookiesPage'
import { ContactPage } from './pages/ContactPage'

/**
 * App es el componente raíz de la aplicación.
 * ¿Para qué? Centralizar providers y rutas en un único punto de entrada.
 * ¿Impacto? BrowserRouter usa la History API del navegador — las URLs
 *            son limpias (/login, /dashboard) en vez de hash (#/login).
 */
function App() {
  return (
    <BrowserRouter>
      {/* ¿Qué? ThemeProvider inicializa el dark/light mode antes de renderizar. */}
      <ThemeProvider>
        {/* ¿Qué? AuthProvider provee el estado de sesión a todas las rutas. */}
        <AuthProvider>
          <Routes>
            {/* ────────────── Rutas públicas ────────────── */}
            <Route path="/" element={<LandingPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            <Route path="/reset-password" element={<ResetPasswordPage />} />
            <Route path="/verify-email" element={<VerifyEmailPage />} />
            <Route path="/terminos-de-uso" element={<TerminosDeUsoPage />} />
            <Route path="/privacidad" element={<PoliticaPrivacidadPage />} />
            <Route path="/cookies" element={<PoliticaCookiesPage />} />
            <Route path="/contacto" element={<ContactPage />} />

            {/* ── Rutas protegidas — requieren JWT válido en memoria ── */}
            <Route element={<ProtectedRoute />}>
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/change-password" element={<ChangePasswordPage />} />
            </Route>

            {/* ¿Qué? Ruta comodín — redirige cualquier URL desconocida al inicio. */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </AuthProvider>
      </ThemeProvider>
    </BrowserRouter>
  )
}

export default App
