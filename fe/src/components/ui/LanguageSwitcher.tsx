/**
 * Archivo: LanguageSwitcher.tsx
 * Descripción: Componente de cambio de idioma entre español e inglés.
 * ¿Para qué? Permitir al usuario seleccionar el idioma de la interfaz
 *            con feedback visual del idioma activo.
 * ¿Impacto? Sin este componente, el idioma se fijaría por detección del
 *            navegador sin opción de cambio manual.
 */

import { useTranslation } from 'react-i18next'

const LANGS = ['es', 'en'] as const

/**
 * LanguageSwitcher renderiza dos botones ES/EN para cambiar el idioma.
 * ¿Para qué? Cambiar el idioma de todas las cadenas de texto de la app
 *            usando el sistema i18next.
 * ¿Impacto? i18next persiste la selección en localStorage ('i18nextLng')
 *            para recordarla en la próxima visita.
 */
export default function LanguageSwitcher() {
  const { i18n } = useTranslation()

  // ¿Qué? Normaliza el código de idioma a 2 letras (ej: 'en-US' → 'en').
  const currentLang = i18n.language.startsWith('en') ? 'en' : 'es'

  return (
    <div
      className="flex overflow-hidden rounded-md border border-gray-200 dark:border-gray-700"
      role="group"
      aria-label="Selector de idioma"
    >
      {LANGS.map((lang) => (
        <button
          key={lang}
          onClick={() => i18n.changeLanguage(lang)}
          aria-pressed={currentLang === lang}
          aria-label={lang === 'es' ? 'Español' : 'English'}
          className={[
            'px-2.5 py-1 text-xs font-medium transition-colors',
            currentLang === lang
              ? 'bg-indigo-600 text-white'
              : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800',
          ].join(' ')}
        >
          {lang.toUpperCase()}
        </button>
      ))}
    </div>
  )
}
