/**
 * Archivo: DashboardPage.tsx
 * Descripción: Página principal del usuario autenticado.
 * ¿Para qué? Mostrar el perfil del usuario y accesos rápidos a acciones
 *            disponibles como cambio de contraseña.
 * ¿Impacto? Es la pantalla central post-login. Solo accesible si hay JWT
 *            válido (ProtectedRoute lo garantiza).
 */

import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../hooks/useAuth'
import LanguageSwitcher from '../components/ui/LanguageSwitcher'
import ThemeToggle from '../components/ui/ThemeToggle'
import Button from '../components/ui/Button'

/**
 * DashboardPage muestra el perfil del usuario y acciones disponibles.
 * ¿Para qué? Punto central de la app post-login. Confirma que la
 *            autenticación fue exitosa y muestra datos del usuario.
 * ¿Impacto? user siempre está definido aquí porque ProtectedRoute
 *            ya validó isAuthenticated antes de renderizar esta página.
 */
export default function DashboardPage() {
  const { t } = useTranslation()
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    // ¿Qué? replace:true para que el botón "atrás" no regrese al dashboard.
    navigate('/login', { replace: true })
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      {/* ¿Qué? Barra de navegación superior con controles globales. */}
      <header className="border-b border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900">
        <div className="mx-auto flex h-14 max-w-4xl items-center justify-between px-4 sm:px-6">
          <span className="font-semibold text-gray-900 dark:text-white">{t('app.name')}</span>
          <div className="flex items-center gap-3">
            <LanguageSwitcher />
            <ThemeToggle />
            {/* ¿Qué? Botón logout alineado a la derecha — acción principal del header. */}
            <Button variant="secondary" onClick={handleLogout}>
              {t('nav.logout')}
            </Button>
          </div>
        </div>
      </header>

      {/* ¿Qué? Contenido principal — grid de dos tarjetas. */}
      <main className="mx-auto max-w-4xl px-4 py-10 sm:px-6">
        <h1 className="mb-6 text-2xl font-bold text-gray-900 dark:text-white">
          {t('dashboard.welcome', { name: user?.full_name ?? '' })}
        </h1>

        <div className="grid gap-4 sm:grid-cols-2">
          {/* ¿Qué? Tarjeta de información del perfil del usuario. */}
          <div className="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-gray-900">
            <h2 className="mb-4 text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
              {t('dashboard.profile')}
            </h2>
            <dl className="space-y-3 text-sm">
              <div>
                <dt className="text-gray-500 dark:text-gray-400">{t('auth.full_name')}</dt>
                <dd className="font-medium text-gray-900 dark:text-white">{user?.full_name}</dd>
              </div>
              <div>
                <dt className="text-gray-500 dark:text-gray-400">{t('auth.email')}</dt>
                <dd className="font-medium text-gray-900 dark:text-white">{user?.email}</dd>
              </div>
              <div>
                <dt className="mb-1 text-gray-500 dark:text-gray-400">
                  {t('dashboard.email_status')}
                </dt>
                <dd>
                  {/* ¿Qué? Badge de color para el estado de verificación del email. */}
                  <span
                    className={[
                      'inline-block rounded px-2 py-0.5 text-xs font-medium',
                      user?.is_email_verified
                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                        : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
                    ].join(' ')}
                  >
                    {user?.is_email_verified
                      ? t('dashboard.email_verified')
                      : t('dashboard.email_not_verified')}
                  </span>
                </dd>
              </div>
            </dl>
          </div>

          {/* ¿Qué? Tarjeta de acciones disponibles para el usuario. */}
          <div className="rounded-xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-gray-900">
            <h2 className="mb-4 text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
              {t('dashboard.actions')}
            </h2>
            <div className="space-y-2">
              <Link
                to="/change-password"
                className="flex items-center justify-between rounded-md border border-gray-200 px-4 py-3 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                <span>{t('auth.change_title')}</span>
                <span aria-hidden="true">→</span>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
