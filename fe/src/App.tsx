/**
 * Archivo: App.tsx
 * Descripción: Componente raíz de la aplicación — placeholder de Fase 5.
 * ¿Para qué? Verificar que el setup (Vite + React + TailwindCSS + i18n) funciona
 *            antes de implementar el router y las páginas en Fase 6.
 * ¿Impacto? Este componente será reemplazado por el router completo en Fase 6.
 */

import { useTranslation } from 'react-i18next'

/**
 * App es el componente raíz de la aplicación.
 * ¿Para qué? Placeholder que verifica: React, TailwindCSS e i18n funcionan.
 */
function App() {
  const { t, i18n } = useTranslation()

  return (
    <div className="min-h-screen bg-white dark:bg-gray-950 flex items-center justify-center">
      <div className="text-center space-y-4 p-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
          {t('app.name')}
        </h1>
        <p className="text-gray-500 dark:text-gray-400">
          {t('app.tagline')}
        </p>
        <div className="flex gap-2 justify-center">
          <button
            onClick={() => i18n.changeLanguage('es')}
            className="px-4 py-2 rounded bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-200 text-sm font-medium"
          >
            ES
          </button>
          <button
            onClick={() => i18n.changeLanguage('en')}
            className="px-4 py-2 rounded bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-200 text-sm font-medium"
          >
            EN
          </button>
        </div>
        <p className="text-xs text-gray-400 dark:text-gray-600">
          Fase 5 — Setup completado ✓ | Tailwind + React {'{18.3.1}'} + i18next
        </p>
      </div>
    </div>
  )
}

export default App
