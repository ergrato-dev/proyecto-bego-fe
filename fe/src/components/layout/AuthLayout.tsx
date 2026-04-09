/**
 * Archivo: AuthLayout.tsx
 * Descripción: Layout compartido para todas las páginas de autenticación.
 * ¿Para qué? Proveer una estructura consistente (header con controles,
 *            tarjeta centrada, footer) reutilizable en login, registro, etc.
 * ¿Impacto? Sin este layout, cada página de auth repetiría la misma
 *            estructura HTML/CSS, violando el principio DRY.
 */

import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import LanguageSwitcher from '../ui/LanguageSwitcher'
import ThemeToggle from '../ui/ThemeToggle'

interface AuthLayoutProps {
  /** Título principal de la página. */
  title: string
  /** Subtítulo opcional — descripción breve bajo el título. */
  subtitle?: string
  children: ReactNode
}

/**
 * AuthLayout envuelve el contenido de formularios de auth en una tarjeta centrada.
 * ¿Para qué? Reutilizar la estructura visual de todas las páginas de auth.
 * ¿Impacto? LanguageSwitcher y ThemeToggle aparecen en todas las páginas de auth
 *            automáticamente sin necesidad de importarlos en cada página.
 */
export default function AuthLayout({
  title,
  subtitle,
  children,
}: AuthLayoutProps) {
  const { t } = useTranslation()

  return (
    <div className="flex min-h-screen flex-col bg-gray-50 dark:bg-gray-950">
      {/* ¿Qué? Header con controles de accesibilidad global (idioma + tema). */}
      <header className="flex items-center justify-end gap-2 p-4">
        <LanguageSwitcher />
        <ThemeToggle />
      </header>

      {/* ¿Qué? Área principal — tarjeta centrada en la pantalla. */}
      <main className="flex flex-1 items-center justify-center px-4 py-8">
        <div className="w-full max-w-md">
          {/* ¿Qué? Tarjeta blanca/oscura con borde sutil — sin gradientes. */}
          <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm dark:border-gray-800 dark:bg-gray-900">
            {/* ¿Qué? Encabezado de la tarjeta con título y subtítulo. */}
            <div className="mb-6">
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {title}
              </h1>
              {subtitle && (
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  {subtitle}
                </p>
              )}
            </div>

            {children}
          </div>

          {/* ¿Qué? Footer con nombre de la aplicación. */}
          <p className="mt-4 text-center text-xs text-gray-400 dark:text-gray-600">
            {t('app.name')}
          </p>
        </div>
      </main>
    </div>
  )
}
