/**
 * Archivo: LandingPage.tsx
 * Descripción: Página de inicio — presentación de la aplicación NN Auth System.
 * ¿Para qué? Dar la bienvenida a usuarios no autenticados con un hero
 *            y botones de acción claros para registrarse o iniciar sesión.
 * ¿Impacto? Es la primera impresión del sistema. Debe comunicar el propósito
 *            y guiar al usuario a la acción correcta (registro o login).
 */

import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import LanguageSwitcher from '../components/ui/LanguageSwitcher'
import ThemeToggle from '../components/ui/ThemeToggle'

/**
 * LandingPage muestra el hero principal con controles globales en el header.
 * ¿Para qué? Página pública accesible sin autenticación.
 * ¿Impacto? Es el punto de entrada por defecto (ruta "/").
 */
export default function LandingPage() {
  const { t } = useTranslation()

  return (
    <div className="flex min-h-screen flex-col bg-gray-50 dark:bg-gray-950">
      {/* ¿Qué? Header con logo, controles de idioma/tema y navegación. */}
      <header className="border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900">
        <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4 sm:px-6">
          <span className="font-semibold text-gray-900 dark:text-white">
            {t('app.name')}
          </span>

          <div className="flex items-center gap-3">
            <LanguageSwitcher />
            <ThemeToggle />

            <Link
              to="/login"
              className="text-sm font-medium text-indigo-600 hover:underline dark:text-indigo-400"
            >
              {t('nav.login')}
            </Link>

            {/* ¿Qué? Botón primario de CTA — alineado a la derecha del header. */}
            <Link
              to="/register"
              className="rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-indigo-700"
            >
              {t('nav.register')}
            </Link>
          </div>
        </div>
      </header>

      {/* ¿Qué? Sección hero — texto centrado con dos CTAs. */}
      <main className="flex flex-1 items-center justify-center px-6 py-20">
        <div className="max-w-2xl space-y-8 text-center">
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-5xl">
            {t('landing.hero_title')}
          </h1>

          <p className="text-lg text-gray-500 dark:text-gray-400">
            {t('landing.hero_subtitle')}
          </p>

          {/* ¿Qué? Botones centrados — excepción al patrón right-aligned
               porque el contexto es una sección hero, no un formulario. */}
          <div className="flex flex-wrap justify-center gap-4">
            <Link
              to="/register"
              className="rounded-md bg-indigo-600 px-6 py-3 font-medium text-white transition-colors hover:bg-indigo-700"
            >
              {t('landing.cta_register')}
            </Link>
            <Link
              to="/login"
              className="rounded-md border border-gray-300 px-6 py-3 font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              {t('landing.cta_login')}
            </Link>
          </div>
        </div>
      </main>

      {/* ¿Qué? Footer mínimo con nombre de la app. */}
      <footer className="py-6 text-center text-xs text-gray-400 dark:text-gray-600">
        {t('app.footer')}
      </footer>
    </div>
  )
}
