/**
 * Archivo: Alert.test.tsx
 * Descripción: Tests del componente Alert — mensaje, tipos y cierre.
 * ¿Para qué? Verificar que el feedback visual llega al usuario y que
 *            el botón de cierre funciona correctamente.
 * ¿Impacto? role="alert" es crítico para lectores de pantalla — si falla,
 *            usuarios con discapacidad visual no reciben el feedback.
 */

import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import Alert from '../../components/ui/Alert'

describe('Alert', () => {
  it('renders the message text', () => {
    render(<Alert type="success" message="Operación exitosa" />)
    expect(screen.getByRole('alert')).toHaveTextContent('Operación exitosa')
  })

  it('has role=alert for screen reader accessibility', () => {
    render(<Alert type="error" message="Ha ocurrido un error" />)
    expect(screen.getByRole('alert')).toBeInTheDocument()
  })

  it('renders without close button when onClose is not provided', () => {
    render(<Alert type="success" message="OK" />)
    expect(screen.queryByRole('button')).toBeNull()
  })

  it('renders close button when onClose is provided', () => {
    render(<Alert type="error" message="Error" onClose={vi.fn()} />)
    expect(screen.getByRole('button', { name: /cerrar/i })).toBeInTheDocument()
  })

  it('calls onClose when close button is clicked', async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()
    render(<Alert type="error" message="Error" onClose={onClose} />)
    await user.click(screen.getByRole('button', { name: /cerrar/i }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('success type uses green color classes', () => {
    render(<Alert type="success" message="OK" />)
    expect(screen.getByRole('alert').className).toContain('green')
  })

  it('error type uses red color classes', () => {
    render(<Alert type="error" message="Error" />)
    expect(screen.getByRole('alert').className).toContain('red')
  })
})
