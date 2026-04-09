/**
 * Archivo: InputField.test.tsx
 * Descripción: Tests del componente InputField — label, error y onChange.
 * ¿Para qué? Verificar accesibilidad (label→input vinculados) y que los
 *            estados de error comunican correctamente al usuario y al DOM.
 * ¿Impacto? Un input sin label asociado viola WCAG 1.3.1 (lectores de pantalla).
 */

import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import InputField from '../../components/ui/InputField'

describe('InputField', () => {
  it('renders a label associated with the input via htmlFor', () => {
    render(<InputField label="Correo electrónico" value="" onChange={vi.fn()} />)
    // getByLabelText verifica que el label está vinculado al input (accesibilidad)
    expect(screen.getByLabelText('Correo electrónico')).toBeInTheDocument()
  })

  it('renders the input with the provided value', () => {
    render(<InputField label="Email" value="test@test.com" onChange={vi.fn()} />)
    expect(screen.getByLabelText('Email')).toHaveValue('test@test.com')
  })

  it('calls onChange with the new value when user types', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    // value="" simula un input controlado
    render(<InputField label="Email" value="" onChange={onChange} />)
    await user.type(screen.getByRole('textbox'), 'a')
    expect(onChange).toHaveBeenCalledWith('a')
  })

  it('shows error message with role=alert when error prop is provided', () => {
    render(<InputField label="Email" value="" onChange={vi.fn()} error="Campo requerido" />)
    expect(screen.getByRole('alert')).toHaveTextContent('Campo requerido')
  })

  it('marks input as aria-invalid when error is present', () => {
    render(<InputField label="Email" value="" onChange={vi.fn()} error="Inválido" />)
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'true')
  })

  it('does not show error message when error prop is absent', () => {
    render(<InputField label="Email" value="" onChange={vi.fn()} />)
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('input is aria-invalid=false when no error', () => {
    render(<InputField label="Email" value="" onChange={vi.fn()} />)
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'false')
  })

  it('renders with type="password" when specified', () => {
    render(<InputField label="Contraseña" type="password" value="" onChange={vi.fn()} />)
    // Los inputs de tipo password no tienen rol "textbox"; buscar por label
    const input = screen.getByLabelText('Contraseña')
    expect(input).toHaveAttribute('type', 'password')
  })

  it('marks input as aria-required when required=true', () => {
    render(<InputField label="Email" value="" onChange={vi.fn()} required />)
    // ¿Qué? Con required=true, el label incluye un asterisco (* en span aria-hidden).
    //        getByLabelText falla porque textContent del label es 'Email *',
    //        así que se busca el input por su rol directamente.
    const input = screen.getByRole('textbox')
    expect(input).toHaveAttribute('aria-required', 'true')
    expect(input).toHaveAttribute('required')
  })
})
