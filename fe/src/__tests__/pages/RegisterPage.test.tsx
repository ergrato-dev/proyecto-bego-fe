/**
 * Archivo: RegisterPage.test.tsx
 * Descripción: Tests de la página de registro — validación y flujo de éxito.
 * ¿Para qué? Verificar que las validaciones del cliente funcionan antes de
 *            llamar al servidor y que el mensaje de éxito se muestra.
 * ¿Impacto? Si la validación de contraseña del cliente falla, se harían
 *            requests al servidor con datos inválidos causando errores innecesarios.
 */

import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'
import { ThemeProvider } from '../../context/ThemeContext'

const mockNavigate = vi.hoisted(() => vi.fn())
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

vi.mock('../../hooks/useAuth')
import { useAuth } from '../../hooks/useAuth'

import RegisterPage from '../../pages/RegisterPage'

const mockRegister = vi.fn()

function setup() {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: false,
    user: null,
    isLoading: false,
    login: vi.fn(),
    logout: vi.fn(),
    register: mockRegister,
  })

  render(
    <MemoryRouter>
      <ThemeProvider>
        <RegisterPage />
      </ThemeProvider>
    </MemoryRouter>
  )
}

describe('RegisterPage', () => {
  beforeEach(() => {
    mockRegister.mockReset()
    mockNavigate.mockReset()
  })

  // ──────────────────── Renderizado ────────────────────

  it('renders full name, email and password fields', () => {
    setup()
    expect(screen.getByLabelText(/nombre completo/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/correo electrónico/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/contraseña/i)).toBeInTheDocument()
  })

  it('renders submit button', () => {
    setup()
    expect(screen.getByRole('button', { name: /registrarse/i })).toBeInTheDocument()
  })

  // ──────────────────── Validación del cliente ────────────────────

  it('shows error when password has fewer than 8 characters', async () => {
    const user = userEvent.setup()
    setup()

    await user.type(screen.getByLabelText(/nombre completo/i), 'Test User')
    await user.type(screen.getByLabelText(/correo electrónico/i), 'test@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'short')
    await user.click(screen.getByRole('button', { name: /registrarse/i }))

    // ¿Qué? La validación del cliente muestra el error sin llamar al servidor.
    expect(screen.getByRole('alert')).toHaveTextContent('Mínimo 8 caracteres')
    expect(mockRegister).not.toHaveBeenCalled()
  })

  // ──────────────────── Registro exitoso ────────────────────

  it('calls register with correct data on form submit', async () => {
    const user = userEvent.setup()
    mockRegister.mockResolvedValue('Cuenta creada. Revisa tu email.')
    setup()

    await user.type(screen.getByLabelText(/nombre completo/i), 'Test User')
    await user.type(screen.getByLabelText(/correo electrónico/i), 'test@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'password123')
    await user.click(screen.getByRole('button', { name: /registrarse/i }))

    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith({
        email: 'test@test.com',
        full_name: 'Test User',
        password: 'password123',
      })
    })
  })

  it('shows success message after successful registration', async () => {
    const user = userEvent.setup()
    mockRegister.mockResolvedValue('Cuenta creada. Revisa tu email.')
    setup()

    await user.type(screen.getByLabelText(/nombre completo/i), 'Test User')
    await user.type(screen.getByLabelText(/correo electrónico/i), 'test@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'password123')
    await user.click(screen.getByRole('button', { name: /registrarse/i }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Cuenta creada')
    })
  })

  it('navigates to /login after 3 seconds following successful registration', async () => {
    // ¿Qué? Usa timers reales con waitFor de 4s en lugar de vi.useFakeTimers().
    // ¿Para qué? vi.useFakeTimers() en cualquier configuración bloquea los mecanismos
    //            internos de userEvent en jsdom (gestión de eventos de puntero/teclado
    //            mediante setTimeout) causando timeout del test completo.
    //            Con timers reales, el comportamiento del componente se verifica con
    //            fidelidad máxima: el setTimeout de 3s funciona como en producción.
    // ¿Impacto? Este test tarda ~3s reales. El timeout del test se aumenta a 6s.
    //            Para un proyecto educativo con test suite pequeño, es aceptable.
    mockRegister.mockResolvedValue('Cuenta creada.')
    setup()
    const user = userEvent.setup()

    await user.type(screen.getByLabelText(/nombre completo/i), 'Test User')
    await user.type(screen.getByLabelText(/correo electrónico/i), 'test@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'securepass')
    await user.click(screen.getByRole('button', { name: /registrarse/i }))

    // El componente llama setTimeout(() => navigate('/login'), 3000).
    // waitFor espera hasta 4s a que mockNavigate sea llamado.
    await waitFor(
      () => expect(mockNavigate).toHaveBeenCalledWith('/login'),
      { timeout: 4000 },
    )
  }, 6000)

  // ──────────────────── Error del servidor ────────────────────

  it('shows error alert when registration fails', async () => {
    const user = userEvent.setup()
    mockRegister.mockRejectedValue(new Error('El email ya está registrado'))
    setup()

    await user.type(screen.getByLabelText(/nombre completo/i), 'Test User')
    await user.type(screen.getByLabelText(/correo electrónico/i), 'existing@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'password123')
    await user.click(screen.getByRole('button', { name: /registrarse/i }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('El email ya está registrado')
    })
  })
})
