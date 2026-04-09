/**
 * Archivo: AuthContext.test.tsx
 * Descripción: Tests del contexto de autenticación — login, logout, registro.
 * ¿Para qué? Verificar que la lógica de sesión funciona correctamente:
 *            el token se almacena en memoria y el usuario se actualiza.
 * ¿Impacto? Si el login no guarda el token o el logout no lo limpia,
 *            el sistema de autenticación queda en estado inconsistente.
 */

import { renderHook, act } from '@testing-library/react'
import { vi, type Mock } from 'vitest'
import type { ReactNode } from 'react'

// ¿Qué? Mocks de los módulos de API antes de importar el contexto.
// ¿Para qué? Evitar llamadas HTTP reales en los tests de unidad.
// ¿Impacto? Sin mocks, los tests dependerian de un servidor activo.
vi.mock('../../api/auth')
vi.mock('../../api/axios', () => ({
  setAccessToken: vi.fn(),
  getAccessToken: vi.fn(() => null),
  default: {},
}))

import { AuthProvider, useAuth } from '../../context/AuthContext'
import * as authApi from '../../api/auth'
import { setAccessToken } from '../../api/axios'

/** Usuario de ejemplo que devuelve el mock de getMe. */
const mockUser = {
  id: 'user-uuid-1',
  email: 'test@test.com',
  full_name: 'Test User',
  is_active: true,
  is_email_verified: true,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
}

/** Wrapper que provee AuthProvider a los hooks de test. */
const wrapper = ({ children }: { children: ReactNode }) => <AuthProvider>{children}</AuthProvider>

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ──────────────────── Estado inicial ────────────────────

  it('starts with user=null and isAuthenticated=false', () => {
    const { result } = renderHook(() => useAuth(), { wrapper })
    expect(result.current.user).toBeNull()
    expect(result.current.isAuthenticated).toBe(false)
    expect(result.current.isLoading).toBe(false)
  })

  // ──────────────────── login ────────────────────

  it('login: sets user and isAuthenticated=true on success', async () => {
    vi.mocked(authApi.login as Mock).mockResolvedValue({
      access_token: 'access-tok',
      refresh_token: 'refresh-tok',
      token_type: 'bearer',
    })
    vi.mocked(authApi.getMe as Mock).mockResolvedValue(mockUser)

    const { result } = renderHook(() => useAuth(), { wrapper })

    await act(async () => {
      await result.current.login({ email: 'test@test.com', password: 'pass123' })
    })

    expect(result.current.user).toEqual(mockUser)
    expect(result.current.isAuthenticated).toBe(true)
  })

  it('login: calls setAccessToken with the access_token', async () => {
    vi.mocked(authApi.login as Mock).mockResolvedValue({
      access_token: 'my-token',
      refresh_token: 'ref',
      token_type: 'bearer',
    })
    vi.mocked(authApi.getMe as Mock).mockResolvedValue(mockUser)

    const { result } = renderHook(() => useAuth(), { wrapper })
    await act(async () => {
      await result.current.login({ email: 'test@test.com', password: 'pass' })
    })

    expect(setAccessToken).toHaveBeenCalledWith('my-token')
  })

  it('login: throws and keeps user=null on API failure', async () => {
    vi.mocked(authApi.login as Mock).mockRejectedValue(new Error('Credenciales incorrectas'))

    const { result } = renderHook(() => useAuth(), { wrapper })

    await expect(
      act(async () => {
        await result.current.login({ email: 'bad@test.com', password: 'wrong' })
      })
    ).rejects.toThrow('Credenciales incorrectas')

    expect(result.current.user).toBeNull()
    expect(result.current.isAuthenticated).toBe(false)
  })

  it('login: resets isLoading to false after failure', async () => {
    vi.mocked(authApi.login as Mock).mockRejectedValue(new Error('Error'))

    const { result } = renderHook(() => useAuth(), { wrapper })

    await act(async () => {
      await result.current.login({ email: 'x@x.com', password: 'pass' }).catch(() => {})
    })

    expect(result.current.isLoading).toBe(false)
  })

  // ──────────────────── logout ────────────────────

  it('logout: sets user=null and isAuthenticated=false', async () => {
    vi.mocked(authApi.login as Mock).mockResolvedValue({
      access_token: 'tok',
      refresh_token: 'ref',
      token_type: 'bearer',
    })
    vi.mocked(authApi.getMe as Mock).mockResolvedValue(mockUser)

    const { result } = renderHook(() => useAuth(), { wrapper })
    await act(async () => {
      await result.current.login({ email: 'test@test.com', password: 'pass' })
    })
    expect(result.current.isAuthenticated).toBe(true)

    act(() => {
      result.current.logout()
    })

    expect(result.current.user).toBeNull()
    expect(result.current.isAuthenticated).toBe(false)
  })

  it('logout: calls setAccessToken(null) to clear the in-memory token', async () => {
    vi.mocked(authApi.login as Mock).mockResolvedValue({
      access_token: 'tok',
      refresh_token: 'ref',
      token_type: 'bearer',
    })
    vi.mocked(authApi.getMe as Mock).mockResolvedValue(mockUser)

    const { result } = renderHook(() => useAuth(), { wrapper })
    await act(async () => {
      await result.current.login({ email: 'test@test.com', password: 'pass' })
    })
    act(() => {
      result.current.logout()
    })

    // ¿Qué? Verificar que el token fue limpiado de memoria.
    // ¿Para qué? Si el token queda, el usuario podría seguir haciendo requests.
    expect(setAccessToken).toHaveBeenCalledWith(null)
  })

  // ──────────────────── register ────────────────────

  it('register: returns the message from the API response', async () => {
    vi.mocked(authApi.register as Mock).mockResolvedValue({
      message: 'Cuenta creada. Revisa tu email.',
      user: mockUser,
    })

    const { result } = renderHook(() => useAuth(), { wrapper })
    let msg = ''

    await act(async () => {
      msg = await result.current.register({
        email: 'new@test.com',
        full_name: 'New User',
        password: 'secure123',
      })
    })

    expect(msg).toBe('Cuenta creada. Revisa tu email.')
    // register NO establece sesión — el usuario debe verificar su email primero
    expect(result.current.user).toBeNull()
  })

  // ──────────────────── Guard ────────────────────

  it('useAuth throws descriptive error outside of AuthProvider', () => {
    // ¿Qué? Verificar que el hook lanza un error claro si se usa mal.
    expect(() => renderHook(() => useAuth())).toThrow('useAuth debe usarse dentro de AuthProvider')
  })
})
