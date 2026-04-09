/**
 * Archivo: Alert.tsx
 * Descripción: Componente de alerta para mensajes de éxito o error.
 * ¿Para qué? Mostrar feedback visual al usuario después de operaciones
 *            con colores semánticos (verde=éxito, rojo=error).
 * ¿Impacto? role="alert" garantiza que lectores de pantalla anuncien
 *            el mensaje inmediatamente cuando aparece en el DOM.
 */

interface AlertProps {
  type: 'success' | 'error'
  message: string
  /** Callback opcional para mostrar un botón X de cierre. */
  onClose?: () => void
}

/**
 * Alert renderiza un banner de notificación accesible con color semántico.
 * ¿Para qué? Comunicar el resultado de operaciones al usuario de forma clara.
 * ¿Impacto? Sin role="alert", los lectores de pantalla no anuncian el mensaje.
 */
export default function Alert({ type, message, onClose }: AlertProps) {
  // ¿Qué? Clases por tipo — sin gradientes, fondos con baja opacidad en dark mode.
  const styles = {
    success:
      'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800 text-green-800 dark:text-green-300',
    error:
      'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-800 dark:text-red-300',
  }

  return (
    <div
      role="alert"
      className={`flex items-start gap-2 rounded-md border px-3 py-2.5 text-sm ${styles[type]}`}
    >
      <span className="flex-1">{message}</span>

      {onClose && (
        <button
          type="button"
          onClick={onClose}
          aria-label="Cerrar alerta"
          className="shrink-0 rounded hover:opacity-70"
        >
          {/* ¿Qué? Ícono X para cerrar la alerta si onClose fue provisto. */}
          <svg className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22z" />
          </svg>
        </button>
      )}
    </div>
  )
}
