/**
 * Archivo: setup.ts (__tests__)
 * Descripción: Configuración global del entorno de testing con Vitest.
 * ¿Para qué? Registrar los matchers extendidos de @testing-library/react
 *            para que aserciones como toBeInTheDocument() estén disponibles.
 * ¿Impacto? Sin este setup, los matchers de DOM no están disponibles en los tests
 *            y arrojan errores de "toBeInTheDocument is not a function".
 */

// ¿Qué? Importar matchers extendidos de @testing-library/jest-dom.
// ¿Para qué? Extiende expect() con aserciones como:
//   - expect(element).toBeInTheDocument()
//   - expect(element).toHaveValue('texto')
//   - expect(element).toBeDisabled()
import '@testing-library/jest-dom'

// ¿Qué? Inicializar i18next y forzar español antes de los tests.
// ¿Para qué? Los componentes usan useTranslation(); sin init, t('key') devuelve la clave.
// ¿Impacto? jsdom reporta navigator.language='en-US', así que LanguageDetector elige inglés.
//           Forzamos 'es' para que los tests busquen el texto correcto en español.
import '../i18n'
import i18n from 'i18next'
// changeLanguage es síncrono cuando los recursos ya están en memoria (JSON bundleado).
void i18n.changeLanguage('es')

// ¿Qué? Mock de window.matchMedia — no implementado por defecto en jsdom.
// ¿Para qué? ThemeContext llama matchMedia('(prefers-color-scheme: dark)').
// ¿Impacto? Sin este mock, cualquier test que use ThemeProvider lanzaría
//            TypeError: window.matchMedia is not a function.
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})
