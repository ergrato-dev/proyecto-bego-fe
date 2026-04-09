/**
 * Archivo: Button.tsx
 * Descripción: Componente de botón reutilizable con variantes y estado de carga.
 * ¿Para qué? Proveer un botón consistente en toda la app con feedback visual
 *            durante operaciones asíncronas (loading spinner).
 * ¿Impacto? Sin estado de loading ni disabled, el usuario podría hacer clic
 *            múltiples veces enviando requests duplicados a la API.
 */

import type { ReactNode } from 'react'

interface ButtonProps {
  children: ReactNode
  type?: 'button' | 'submit' | 'reset'
  onClick?: () => void
  /** true mientras la operación está en curso — muestra spinner y deshabilita. */
  loading?: boolean
  disabled?: boolean
  /** primary: indigo sólido | secondary: gris neutro */
  variant?: 'primary' | 'secondary'
  className?: string
}

/**
 * Button renderiza un botón accesible con variantes de estilo y spinner de carga.
 * ¿Para qué? Estandarizar la apariencia y comportamiento de botones en toda la app.
 * ¿Impacto? Aria-busy comunica el estado de carga a lectores de pantalla.
 */
export default function Button({
  children,
  type = 'button',
  onClick,
  loading = false,
  disabled = false,
  variant = 'primary',
  className = '',
}: ButtonProps) {
  const baseClasses = [
    'inline-flex items-center gap-2 rounded-md px-4 py-2 text-sm font-medium',
    'transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2',
    'disabled:cursor-not-allowed disabled:opacity-50',
  ].join(' ')

  // ¿Qué? Estilos por variante — sin gradientes, solo colores sólidos.
  const variantClasses = {
    primary: 'bg-indigo-600 text-white hover:bg-indigo-700',
    secondary:
      'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700',
  }

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || loading}
      aria-busy={loading}
      className={`${baseClasses} ${variantClasses[variant]} ${className}`}
    >
      {/* ¿Qué? SVG spinner animado visible solo en estado loading. */}
      {loading && (
        <svg
          className="h-4 w-4 animate-spin"
          viewBox="0 0 24 24"
          fill="none"
          aria-hidden="true"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
          />
        </svg>
      )}
      {children}
    </button>
  )
}
