/**
 * Archivo: LoginPage.tsx
 * Descripción: Página de inicio de sesión con formulario email y contraseña.
 * ¿Para qué? Permitir a usuarios registrados obtener un JWT para acceder
 *            a rutas protegidas de la aplicación.
 * ¿Impacto? Sin esta página, no hay forma de iniciar sesión. Los errores
 *            del backend (credenciales incorrectas) se muestran al usuario
 *            sin revelar si el email existe (mensaje genérico del servidor).
 */

import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../hooks/useAuth'
import AuthLayout from '../components/layout/AuthLayout'
import InputField from '../components/ui/InputField'
import Button from '../components/ui/Button'
import Alert from '../components/ui/Alert'
import { extractError } from '../utils/error'

/**
 * LoginPage gestiona el formulario de autenticación.
 * ¿Para qué? Recoger credenciales, llamar a AuthContext.login y
 *            redirigir al dashboard si tiene éxito.
 * ¿Impacto? noValidate en el form desactiva la validación nativa del navegador
 *            para que el feedback de error sea controlado por la app.
 */
export default function LoginPage() {
  const { t } = useTranslation()
  const { login, isLoading } = useAuth()
  const navigate = useNavigate()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    try {
      await login({ email, password })
      // ¿Qué? replace:true evita que el usuario regrese al login con el botón atrás.
      navigate('/dashboard', { replace: true })
    } catch (err) {
      setError(extractError(err))
    }
  }

  return (
    <AuthLayout title={t('auth.login_title')} subtitle={t('auth.login_subtitle')}>
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {error && <Alert type="error" message={error} onClose={() => setError('')} />}

        <InputField
          id="email"
          label={t('auth.email')}
          type="email"
          value={email}
          onChange={setEmail}
          placeholder={t('auth.email_placeholder')}
          required
        />

        <InputField
          id="password"
          label={t('auth.password')}
          type="password"
          value={password}
          onChange={setPassword}
          placeholder={t('auth.password_placeholder')}
          required
        />

        {/* ¿Qué? Enlace de "¿Olvidaste tu contraseña?" alineado a la derecha. */}
        <div className="flex justify-end">
          <Link
            to="/forgot-password"
            className="text-sm text-indigo-600 hover:underline dark:text-indigo-400"
          >
            {t('auth.forgot_password')}
          </Link>
        </div>

        {/* ¿Qué? Botón submit alineado a la derecha — convención del proyecto. */}
        <div className="flex justify-end">
          <Button type="submit" loading={isLoading}>
            {t('auth.login_button')}
          </Button>
        </div>

        <p className="text-center text-sm text-gray-500 dark:text-gray-400">
          {t('auth.no_account')}{' '}
          <Link
            to="/register"
            className="text-indigo-600 hover:underline dark:text-indigo-400"
          >
            {t('nav.register')}
          </Link>
        </p>
      </form>
    </AuthLayout>
  )
}
