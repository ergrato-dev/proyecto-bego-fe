/**
 * Archivo: ThemeToggle.tsx
 * Descripción: Botón para alternar entre modo oscuro y modo claro.
 * ¿Para qué? Permitir al usuario cambiar el tema visual en cualquier momento
 *            con feedback visual (ícono sol/luna).
 * ¿Impacto? Sin este componente, solo se usaría la preferencia del sistema
 *            operativo sin opción de personalización manual.
 */

import { useTranslation } from 'react-i18next'
import { useTheme } from '../../hooks/useTheme'

/**
 * ThemeToggle renderiza un botón con ícono de sol (modo oscuro activo)
 * o luna (modo claro activo) para cambiar el tema.
 * ¿Para qué? El ícono muestra el tema AL QUE CAMBIARÁ al hacer clic.
 * ¿Impacto? El aria-label describe la acción futura para lectores de pantalla.
 */
export default function ThemeToggle() {
  const { theme, toggleTheme } = useTheme()
  const { t } = useTranslation()

  const isDark = theme === 'dark'

  return (
    <button
      onClick={toggleTheme}
      aria-label={isDark ? t('theme.toggle_light') : t('theme.toggle_dark')}
      className="rounded-md p-1.5 text-gray-500 transition-colors hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
    >
      {isDark ? (
        // ¿Qué? Ícono sol — se muestra cuando el modo oscuro está activo.
        // ¿Para qué? Indica que al hacer clic se activará el modo claro.
        <svg
          className="h-4 w-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <circle cx="12" cy="12" r="5" />
          <line x1="12" y1="1" x2="12" y2="3" />
          <line x1="12" y1="21" x2="12" y2="23" />
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
          <line x1="1" y1="12" x2="3" y2="12" />
          <line x1="21" y1="12" x2="23" y2="12" />
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
        </svg>
      ) : (
        // ¿Qué? Ícono luna — se muestra cuando el modo claro está activo.
        // ¿Para qué? Indica que al hacer clic se activará el modo oscuro.
        <svg
          className="h-4 w-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        </svg>
      )}
    </button>
  )
}
