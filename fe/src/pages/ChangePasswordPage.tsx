/**
 * Archivo: ChangePasswordPage.tsx
 * Descripción: Página para cambiar la contraseña del usuario autenticado.
 * ¿Para qué? Permitir que el usuario actualice su contraseña proporcionando
 *            la contraseña actual como verificación de identidad.
 * ¿Impacto? El backend verifica la contraseña actual antes de actualizar.
 *            Validamos coincidencia y longitud mínima en el cliente por UX.
 */

import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import * as authApi from '../api/auth'
import AuthLayout from '../components/layout/AuthLayout'
import InputField from '../components/ui/InputField'
import Button from '../components/ui/Button'
import Alert from '../components/ui/Alert'
import { extractError } from '../utils/error'

/**
 * ChangePasswordPage gestiona el formulario de cambio de contraseña.
 * ¿Para qué? Ruta protegida — solo accesible con JWT válido.
 * ¿Impacto? Llama directamente a authApi.changePassword (no a AuthContext)
 *            porque esta acción no afecta el estado de sesión.
 */
export default function ChangePasswordPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()

  const [current, setCurrent] = useState('')
  const [next, setNext] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')

    // ¿Qué? Validaciones del cliente — UX inmediata antes de llamar al servidor.
    if (next.length < 8) {
      setError(t('errors.password_min'))
      return
    }
    if (next !== confirm) {
      setError(t('errors.passwords_not_match'))
      return
    }

    setLoading(true)
    try {
      await authApi.changePassword({ current_password: current, new_password: next })
      setSuccess(t('success.password_changed'))
      // ¿Qué? Redirigir al dashboard después de 2 segundos para que lea el mensaje.
      setTimeout(() => navigate('/dashboard'), 2000)
    } catch (err) {
      setError(extractError(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout title={t('auth.change_title')}>
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {error && <Alert type="error" message={error} onClose={() => setError('')} />}
        {success && <Alert type="success" message={success} />}

        <InputField
          id="current-password"
          label={t('auth.current_password')}
          type="password"
          value={current}
          onChange={setCurrent}
          required
        />

        <InputField
          id="new-password"
          label={t('auth.new_password')}
          type="password"
          value={next}
          onChange={setNext}
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

        {/* ¿Qué? Dos botones alineados a la derecha — Cancelar y Guardar. */}
        <div className="flex justify-end gap-2">
          <Button type="button" variant="secondary" onClick={() => navigate('/dashboard')}>
            {t('common.cancel')}
          </Button>
          <Button type="submit" loading={loading} disabled={!!success}>
            {t('auth.change_button')}
          </Button>
        </div>
      </form>
    </AuthLayout>
  )
}
