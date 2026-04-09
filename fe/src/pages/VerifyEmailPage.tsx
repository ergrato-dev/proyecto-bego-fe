/**
 * Archivo: VerifyEmailPage.tsx
 * Descripción: Página de verificación de email — se activa automáticamente.
 * ¿Para qué? Completar el flujo de verificación de email cuando el usuario
 *            hace clic en el enlace enviado por el sistema tras el registro.
 * ¿Impacto? Sin esta página, el token de verificación del email no tendría
 *            a dónde apuntar y el flujo de registro quedaría incompleto.
 */

import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import * as authApi from '../api/auth'
import AuthLayout from '../components/layout/AuthLayout'
import Alert from '../components/ui/Alert'
import { extractError } from '../utils/error'

type VerifyStatus = 'loading' | 'success' | 'error'

/**
 * VerifyEmailPage llama a la API de verificación automáticamente al montar.
 * ¿Para qué? El usuario solo hace clic en el enlace del email — la verificación
 *            debe ser automática, sin formularios adicionales.
 * ¿Impacto? El token viene en ?token=<uuid>. Si está ausente o expiró,
 *            se muestra el error del servidor de forma clara.
 */
export default function VerifyEmailPage() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()

  const token = searchParams.get('token') ?? ''
  const [status, setStatus] = useState<VerifyStatus>('loading')
  const [message, setMessage] = useState('')

  // ¿Qué? useEffect vacío con dependencia [token] — se ejecuta una sola vez al montar.
  // ¿Para qué? Verificar el email automáticamente sin interacción del usuario.
  // ¿Impacto? Si falta el token, se evita llamar al servidor y se muestra error local.
  useEffect(() => {
    if (!token) {
      setStatus('error')
      setMessage(t('errors.invalid_token'))
      return
    }

    authApi
      .verifyEmail({ token })
      .then(() => {
        setStatus('success')
        setMessage(t('success.email_verified'))
      })
      .catch((err: unknown) => {
        setStatus('error')
        setMessage(extractError(err))
      })
  }, [token, t])

  return (
    <AuthLayout title={t('auth.verify_email_title')}>
      <div className="space-y-4">
        {/* ¿Qué? Estado de carga — pulsación visual mientras se espera la API. */}
        {status === 'loading' && (
          <p className="animate-pulse text-sm text-gray-500 dark:text-gray-400">
            {t('auth.verify_email_message')}
          </p>
        )}

        {status === 'success' && <Alert type="success" message={message} />}
        {status === 'error' && <Alert type="error" message={message} />}

        {/* ¿Qué? Botón visible solo cuando hay un resultado definitivo. */}
        {status !== 'loading' && (
          <div className="flex justify-end">
            <Link
              to="/login"
              className="inline-flex items-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-700"
            >
              {t('nav.login')}
            </Link>
          </div>
        )}
      </div>
    </AuthLayout>
  )
}
