/**
 * Archivo: ResetPasswordPage.tsx
 * Descripción: Página para restablecer la contraseña usando el token del email.
 * ¿Para qué? Completar el flujo de recuperación de contraseña permitiendo
 *            al usuario definir una nueva contraseña con el token enviado por email.
 * ¿Impacto? El token se lee desde la URL (?token=...). Si no existe o expiró,
 *            el backend rechazará la solicitud con un error claro.
 */

import { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import * as authApi from '../api/auth'
import AuthLayout from '../components/layout/AuthLayout'
import InputField from '../components/ui/InputField'
import Button from '../components/ui/Button'
import Alert from '../components/ui/Alert'
import { extractError } from '../utils/error'

/**
 * ResetPasswordPage lee el token de la URL y actualiza la contraseña.
 * ¿Para qué? El token de reset viaja por email en un enlace con ?token=<uuid>.
 *            Esta página lo extrae y lo envía junto con la nueva contraseña.
 * ¿Impacto? Si falta el token en la URL, se muestra error inmediatamente
 *            sin llamar al servidor — validación de inputs en el cliente.
 */
export default function ResetPasswordPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  // ¿Qué? El token viene como query param: /reset-password?token=<uuid>
  const token = searchParams.get('token') ?? ''

  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')

    if (!token) {
      setError(t('errors.invalid_token'))
      return
    }
    if (password.length < 8) {
      setError(t('errors.password_min'))
      return
    }
    if (password !== confirm) {
      setError(t('errors.passwords_not_match'))
      return
    }

    setLoading(true)
    try {
      await authApi.resetPassword({ token, new_password: password })
      setSuccess(t('success.password_reset'))
      // ¿Qué? Redirigir al login después de 3 segundos.
      setTimeout(() => navigate('/login', { replace: true }), 3000)
    } catch (err) {
      setError(extractError(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout
      title={t('auth.reset_title')}
      subtitle={t('auth.reset_subtitle')}
    >
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {success && <Alert type="success" message={success} />}
        {error && <Alert type="error" message={error} onClose={() => setError('')} />}

        <InputField
          id="new-password"
          label={t('auth.new_password')}
          type="password"
          value={password}
          onChange={setPassword}
          placeholder={t('auth.password_placeholder')}
          required
        />

        <InputField
          id="confirm-password"
          label={t('auth.confirm_password')}
          type="password"
          value={confirm}
          onChange={setConfirm}
          required
        />

        <div className="flex justify-end">
          <Button type="submit" loading={loading} disabled={!!success}>
            {t('auth.reset_button')}
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
