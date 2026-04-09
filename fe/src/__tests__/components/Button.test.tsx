/**
 * Archivo: Button.test.tsx
 * Descripción: Tests del componente Button — renderizado, loading y click.
 * ¿Para qué? Verificar que el botón comunica correctamente sus estados
 *            y dispara eventos al usuario.
 * ¿Impacto? Si el botón no deshabilita durante loading, el usuario puede
 *            enviar formularios múltiples veces causando requests duplicados.
 */

import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import Button from '../../components/ui/Button'

describe('Button', () => {
  it('renders children text', () => {
    render(<Button>Guardar</Button>)
    expect(screen.getByRole('button', { name: 'Guardar' })).toBeInTheDocument()
  })

  it('applies primary styles by default', () => {
    render(<Button>Entrar</Button>)
    const btn = screen.getByRole('button')
    expect(btn.className).toContain('bg-indigo-600')
  })

  it('applies secondary styles when variant=secondary', () => {
    render(<Button variant="secondary">Cancelar</Button>)
    const btn = screen.getByRole('button')
    expect(btn.className).toContain('bg-gray-100')
  })

  it('is disabled and shows aria-busy when loading', () => {
    render(<Button loading>Cargando</Button>)
    const btn = screen.getByRole('button')
    expect(btn).toBeDisabled()
    expect(btn).toHaveAttribute('aria-busy', 'true')
  })

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Deshabilitado</Button>)
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('calls onClick when clicked', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Click</Button>)
    await user.click(screen.getByRole('button'))
    expect(onClick).toHaveBeenCalledTimes(1)
  })

  it('does not call onClick when loading', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(
      <Button onClick={onClick} loading>
        Click
      </Button>
    )
    await user.click(screen.getByRole('button'))
    expect(onClick).not.toHaveBeenCalled()
  })

  it('renders with type="submit" when specified', () => {
    render(<Button type="submit">Enviar</Button>)
    expect(screen.getByRole('button')).toHaveAttribute('type', 'submit')
  })
})
