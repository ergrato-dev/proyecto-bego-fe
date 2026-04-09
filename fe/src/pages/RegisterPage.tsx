/**
 * Archivo: RegisterPage.tsx
 * Descripción: Página de registro de nuevos usuarios.
 * ¿Para qué? Permitir que personas creen una cuenta en el sistema
 *            con nombre, email y contraseña.
 * ¿Impacto? Al registrarse, el usuario NO queda automáticamente logueado.
 *            Debe verificar su email primero. Esta página comunica eso
 *            mostrando el mensaje de éxito del servidor y redirigiendo al login.
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
 * RegisterPage gestiona el formulario de creación de cuenta.
 * ¿Para qué? Recoger los datos del nuevo usuario, validar la contraseña
 *            mínima en el cliente y llamar a AuthContext.register.
 * ¿Impacto? La validación del cliente (min 8 chars) es solo UX.
 *            El backend valida de nuevo — nunca confiar solo en el frontend.
 */
export default function RegisterPage() {
  const { t } = useTranslation()
  const { register, isLoading } = useAuth()
  const navigate = useNavigate()

  const [fullName, setFullName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setSuccess('')

    // ¿Qué? Validación mínima de contraseña en el cliente — UX rápida.
    if (password.length < 8) {
      setError(t('errors.password_min'))
      return
    }

    try {
      const msg = await register({ email, full_name: fullName, password })
      setSuccess(msg || t('success.registered'))
      // ¿Qué? Redirigir al login después de 3 segundos para que el usuario
      //        lea el mensaje antes de ser redirigido.
      setTimeout(() => navigate('/login'), 3000)
    } catch (err) {
      setError(extractError(err))
    }
  }

  return (
    <AuthLayout title={t('auth.register_title')} subtitle={t('auth.register_subtitle')}>
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {error && <Alert type="error" message={error} onClose={() => setError('')} />}
        {success && <Alert type="success" message={success} />}

        <InputField
          id="full-name"
          label={t('auth.full_name')}
          type="text"
          value={fullName}
          onChange={setFullName}
          placeholder={t('auth.full_name_placeholder')}
          required
        />

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

        {/* ¿Qué? Botón deshabilitado post-éxito para evitar envíos duplicados. */}
        <div className="flex justify-end">
          <Button type="submit" loading={isLoading} disabled={!!success}>
            {t('auth.register_button')}
          </Button>
        </div>

        <p className="text-center text-sm text-gray-500 dark:text-gray-400">
          {t('auth.already_have_account')}{' '}
          <Link
            to="/login"
            className="text-indigo-600 hover:underline dark:text-indigo-400"
          >
            {t('nav.login')}
          </Link>
        </p>
      </form>
    </AuthLayout>
  )
}
