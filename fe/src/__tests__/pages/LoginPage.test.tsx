/**
 * Archivo: LoginPage.test.tsx
 * Descripción: Tests de la página de login — renderizado, éxito y error.
 * ¿Para qué? Verificar el flujo completo: mostrar formulario, enviar
 *            credenciales y reaccionar al resultado del servidor.
 * ¿Impacto? Si el login no navega al dashboard en éxito, o no muestra
 *            el error en fallo, el formulario queda inutilizable para el usuario.
 */

import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'
import { ThemeProvider } from '../../context/ThemeContext'

// ¿Qué? mockNavigate usa vi.hoisted para estar disponible en el factory de vi.mock.
// ¿Para qué? vi.mock se hoistea al inicio del archivo antes de las inicializaciones
//            normales de variables. Vi.hoisted garantiza que mockNavigate esté listo.
const mockNavigate = vi.hoisted(() => vi.fn())
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => mockNavigate }
})

vi.mock('../../hooks/useAuth')
import { useAuth } from '../../hooks/useAuth'

import LoginPage from '../../pages/LoginPage'

/** Mocks reutilizables del AuthContext para los tests de LoginPage. */
const mockLogin = vi.fn()

function setup() {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: false,
    user: null,
    isLoading: false,
    login: mockLogin,
    logout: vi.fn(),
    register: vi.fn(),
  })

  render(
    // ¿Qué? ThemeProvider necesario porque AuthLayout renderiza ThemeToggle → useTheme().
    <MemoryRouter>
      <ThemeProvider>
        <LoginPage />
      </ThemeProvider>
    </MemoryRouter>
  )
}

describe('LoginPage', () => {
  beforeEach(() => {
    mockLogin.mockReset()
    mockNavigate.mockReset()
  })

  // ──────────────────── Renderizado ────────────────────

  it('renders email and password fields', () => {
    setup()
    // Buscar por texto del label que viene de i18n (español por defecto)
    expect(screen.getByLabelText(/correo electrónico/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/contraseña/i)).toBeInTheDocument()
  })

  it('renders submit button', () => {
    setup()
    expect(screen.getByRole('button', { name: /entrar/i })).toBeInTheDocument()
  })

  it('renders forgot password link', () => {
    setup()
    expect(screen.getByRole('link', { name: /olvidaste/i })).toBeInTheDocument()
  })

  // ──────────────────── Login exitoso ────────────────────

  it('calls login with correct credentials on form submit', async () => {
    const user = userEvent.setup()
    mockLogin.mockResolvedValue(undefined)
    setup()

    await user.type(screen.getByLabelText(/correo electrónico/i), 'test@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'password123')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        email: 'test@test.com',
        password: 'password123',
      })
    })
  })

  it('navigates to /dashboard with replace after successful login', async () => {
    const user = userEvent.setup()
    mockLogin.mockResolvedValue(undefined)
    setup()

    await user.type(screen.getByLabelText(/correo electrónico/i), 'test@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'pass')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/dashboard', { replace: true })
    })
  })

  // ──────────────────── Login fallido ────────────────────

  it('shows error alert when login fails', async () => {
    const user = userEvent.setup()
    mockLogin.mockRejectedValue(new Error('Credenciales incorrectas'))
    setup()

    await user.type(screen.getByLabelText(/correo electrónico/i), 'bad@test.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'wrong')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Credenciales incorrectas')
    })
  })

  it('clears error alert when close button is clicked', async () => {
    const user = userEvent.setup()
    mockLogin.mockRejectedValue(new Error('Error de red'))
    setup()

    await user.type(screen.getByLabelText(/correo electrónico/i), 'x@x.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'pass')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument())

    // Cerrar la alerta
    await user.click(screen.getByRole('button', { name: /cerrar/i }))
    // waitFor por si el estado de React no se ha flusheado aún
    await waitFor(() => expect(screen.queryByRole('alert')).toBeNull())
  })

  it('does NOT navigate to dashboard on failed login', async () => {
    const user = userEvent.setup()
    mockLogin.mockRejectedValue(new Error('Error'))
    setup()

    await user.type(screen.getByLabelText(/correo electrónico/i), 'x@x.com')
    await user.type(screen.getByLabelText(/contraseña/i), 'pass')
    await user.click(screen.getByRole('button', { name: /entrar/i }))

    await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument())
    expect(mockNavigate).not.toHaveBeenCalled()
  })
})
