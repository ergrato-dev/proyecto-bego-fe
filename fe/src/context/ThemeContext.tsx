/**
 * Archivo: ThemeContext.tsx
 * Descripción: Contexto React que gestiona el tema visual (dark/light mode).
 * ¿Para qué? Proveer a toda la aplicación acceso al tema actual y la función
 *            para alternarlo, con persistencia en localStorage.
 * ¿Impacto? Sin este contexto, cada componente gestionaría su propio estado
 *            de tema, causando inconsistencias visuales.
 */

import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'

type Theme = 'light' | 'dark'

interface ThemeContextType {
  theme: Theme
  toggleTheme: () => void
}

// ¿Qué? Contexto con valor null — se valida en useTheme con el guard.
const ThemeContext = createContext<ThemeContextType | null>(null)

/**
 * getInitialTheme determina el tema inicial de la aplicación.
 * ¿Para qué? Priorizar la preferencia guardada del usuario;
 *            si no existe, usar la preferencia del sistema operativo.
 * ¿Impacto? Sin esto, el tema siempre sería 'light' ignorando la preferencia
 *            del usuario y el modo oscuro configurado en el SO.
 */
function getInitialTheme(): Theme {
  const stored = localStorage.getItem('nn-theme') as Theme | null
  if (stored === 'dark' || stored === 'light') return stored
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

/**
 * ThemeProvider envuelve la aplicación y provee el contexto de tema.
 * ¿Para qué? Centralizar la lógica de dark/light mode en un único lugar.
 * ¿Impacto? Todos los componentes hijos pueden leer/cambiar el tema
 *            sin prop drilling.
 */
export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(getInitialTheme)

  // ¿Qué? Efecto que sincroniza el estado React con el DOM y localStorage.
  // ¿Para qué? La clase 'dark' en <html> activa las variantes dark: de Tailwind.
  // ¿Impacto? Sin esto, cambiar el estado no afectaría los estilos visuales.
  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark')
    localStorage.setItem('nn-theme', theme)
  }, [theme])

  const toggleTheme = () => setTheme((prev) => (prev === 'dark' ? 'light' : 'dark'))

  return <ThemeContext.Provider value={{ theme, toggleTheme }}>{children}</ThemeContext.Provider>
}

/**
 * useTheme hook para acceder al contexto de tema desde cualquier componente.
 * ¿Para qué? Encapsular el useContext con validación — lanzar error claro
 *            si se usa fuera del ThemeProvider.
 * ¿Impacto? Errores de uso incorrecto se detectan en desarrollo inmediatamente.
 */
export function useTheme(): ThemeContextType {
  const ctx = useContext(ThemeContext)
  if (!ctx) throw new Error('useTheme debe usarse dentro de ThemeProvider')
  return ctx
}
