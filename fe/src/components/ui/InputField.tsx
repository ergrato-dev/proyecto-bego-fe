/**
 * Archivo: InputField.tsx
 * Descripción: Componente de campo de formulario con label, estado de error
 *              y atributos de accesibilidad ARIA.
 * ¿Para qué? Proveer un input reutilizable que cumpla con WCAG AA:
 *            labels asociados, aria-invalid, aria-describedby para errores.
 * ¿Impacto? Sin labels ni aria-*, los lectores de pantalla no pueden
 *            asociar el campo con su descripción, excluyendo usuarios con
 *            discapacidad visual.
 */

interface InputFieldProps {
  /** Texto del label visible sobre el input. */
  label: string
  type?: string
  value: string
  /** Callback que recibe el nuevo valor como string. */
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  /** Mensaje de error — si existe, el input se marca como inválido. */
  error?: string
  /** id HTML explícito; si se omite se genera desde el label. */
  id?: string
}

/**
 * InputField renderiza un campo de formulario accesible con label y validación.
 * ¿Para qué? Encapsular la lógica de estado de error y atributos aria
 *            para no repetirla en cada formulario.
 * ¿Impacto? Un input sin label accesible no cumple WCAG 2.1 Success Criterion 1.3.1.
 */
export default function InputField({
  label,
  type = 'text',
  value,
  onChange,
  placeholder,
  required,
  error,
  id,
}: InputFieldProps) {
  // ¿Qué? Genera un id único desde el label si no se provee uno explícito.
  // ¿Para qué? El atributo htmlFor del label necesita coincidir con el id del input.
  const inputId = id ?? label.toLowerCase().replace(/\s+/g, '-')

  return (
    <div className="space-y-1">
      <label
        htmlFor={inputId}
        className="block text-sm font-medium text-gray-700 dark:text-gray-300"
      >
        {label}
        {required && (
          <span className="ml-1 text-red-500" aria-hidden="true">
            *
          </span>
        )}
      </label>

      <input
        id={inputId}
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        required={required}
        aria-required={required}
        aria-invalid={!!error}
        aria-describedby={error ? `${inputId}-error` : undefined}
        className={[
          'w-full rounded-md border px-3 py-2 text-sm',
          'bg-white dark:bg-gray-800',
          'text-gray-900 dark:text-gray-100',
          'placeholder-gray-400 dark:placeholder-gray-500',
          'transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-1',
          error
            ? 'border-red-500 dark:border-red-500'
            : 'border-gray-300 dark:border-gray-600',
        ].join(' ')}
      />

      {/* ¿Qué? Mensaje de error con role="alert" para lectores de pantalla. */}
      {error && (
        <p
          id={`${inputId}-error`}
          role="alert"
          className="text-xs text-red-600 dark:text-red-400"
        >
          {error}
        </p>
      )}
    </div>
  )
}
