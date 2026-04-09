/**
 * Archivo: ForgotPasswordPage.tsx
 * Descripción: Página para solicitar el email de recuperación de contraseña.
 * ¿Para qué? Permitir que usuarios que olvidaron su contraseña inicien
 *            el flujo de recuperación por email.
 * ¿Impacto? El backend siempre responde el mismo mensaje genérico independiente
 *            de si el email existe — esto evita el email enumeration attack.
 */

import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import * as authApi from '../api/auth'
import AuthLayout from '../components/layout/AuthLayout'
import InputField from '../components/ui/InputField'
import Button from '../components/ui/Button'
import Alert from '../components/ui/Alert'
import { extractError } from '../utils/error'

/**
 * ForgotPasswordPage solicita el email de recuperación.
 * ¿Para qué? Disparar el flujo que envía el token de reset al email del usuario.
 * ¿Impacto? Seguridad: el backend no confirma si el email existe.
 *            El botón se deshabilita post-éxito para evitar spam de emails.
 */
export default function ForgotPasswordPage() {
  const { t } = useTranslation()

  const [email, setEmail] = useState('')
  const [success, setSuccess] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await authApi.forgotPassword({ email })
      setSuccess(t('success.forgot_sent'))
    } catch (err) {
      setError(extractError(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout
      title={t('auth.forgot_title')}
      subtitle={t('auth.forgot_subtitle')}
    >
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {success && <Alert type="success" message={success} />}
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

        {/* ¿Qué? Botón deshabilitado post-éxito para evitar múltiples solicitudes. */}
        <div className="flex justify-end">
          <Button type="submit" loading={loading} disabled={!!success}>
            {t('auth.forgot_button')}
          </Button>
        </div>

        <p className="text-center text-sm text-gray-500 dark:text-gray-400">
          <Link
            to="/login"
            className="text-indigo-600 hover:underline dark:text-indigo-400"
          >
            ← {t('nav.back_to_login')}
          </Link>
        </p>
      </form>
    </AuthLayout>
  )
}
