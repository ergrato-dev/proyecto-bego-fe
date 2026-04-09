/**
 * Archivo: ProtectedRoute.test.tsx
 * Descripción: Tests del componente ProtectedRoute — acceso y redirección.
 * ¿Para qué? Verificar que las rutas protegidas solo son accesibles con
 *            sesión activa y que redirigen al login en caso contrario.
 * ¿Impacto? Si ProtectedRoute no redirige, cualquier URL de /dashboard
 *            sería accesible sin autenticación, exponiendo datos del usuario.
 */

import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi } from 'vitest'

// ¿Qué? Mockear useAuth para controlar el estado de autenticación en tests.
vi.mock('../../hooks/useAuth')
import { useAuth } from '../../hooks/useAuth'

import ProtectedRoute from '../../components/layout/ProtectedRoute'

/**
 * renderWithRouter monta ProtectedRoute dentro de un MemoryRouter con rutas.
 * ¿Para qué? Simular la navegación del router sin un navegador real.
 * ¿Impacto? isAuthenticated controla si se muestra la ruta protegida o el login.
 */
function renderWithRouter(isAuthenticated: boolean) {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated,
    user: null,
    isLoading: false,
    login: vi.fn(),
    logout: vi.fn(),
    register: vi.fn(),
  })

  render(
    <MemoryRouter initialEntries={['/dashboard']}>
      <Routes>
        <Route path="/login" element={<div>Página de Login</div>} />
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<div>Panel de Control</div>} />
        </Route>
      </Routes>
    </MemoryRouter>
  )
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders the protected child when user is authenticated', () => {
    renderWithRouter(true)
    expect(screen.getByText('Panel de Control')).toBeInTheDocument()
    expect(screen.queryByText('Página de Login')).toBeNull()
  })

  it('redirects to /login when user is NOT authenticated', () => {
    renderWithRouter(false)
    expect(screen.getByText('Página de Login')).toBeInTheDocument()
    expect(screen.queryByText('Panel de Control')).toBeNull()
  })
})
