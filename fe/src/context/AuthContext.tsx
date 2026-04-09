/**
 * Archivo: AuthContext.tsx
 * Descripción: Contexto React que gestiona el estado de autenticación global.
 * ¿Para qué? Proveer a toda la aplicación acceso al usuario autenticado,
 *            tokens y las acciones de login/logout/register.
 * ¿Impacto? Sin este contexto, no habría forma de saber si el usuario está
 *            autenticado ni de proteger rutas que requieren JWT.
 *
 * Nota de seguridad: el access_token se almacena EXCLUSIVAMENTE en memoria
 * (variable de módulo en axios.ts) y NUNCA en localStorage/sessionStorage.
 * Esto mitiga ataques XSS que podrían robar tokens persistidos.
 * En producción, la solución óptima son httpOnly cookies gestionadas por el servidor.
 */

import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from 'react'

import * as authApi from '../api/auth'
import { setAccessToken } from '../api/axios'
import type {
  LoginRequest,
  RegisterRequest,
  UserResponse,
} from '../types/auth'

interface AuthContextType {
  /** Usuario autenticado actual, null si no hay sesión. */
  user: UserResponse | null
  /** true si hay un usuario autenticado en memoria. */
  isAuthenticated: boolean
  /** true mientras se ejecuta una operación de auth. */
  isLoading: boolean
  /** Inicia sesión: llama a la API, guarda el token y obtiene el perfil. */
  login: (data: LoginRequest) => Promise<void>
  /** Cierra sesión: limpia el token en memoria y el estado del usuario. */
  logout: () => void
  /** Registra un nuevo usuario y retorna el mensaje de éxito del servidor. */
  register: (data: RegisterRequest) => Promise<string>
}

const AuthContext = createContext<AuthContextType | null>(null)

/**
 * AuthProvider envuelve la aplicación y gestiona el estado de autenticación.
 * ¿Para qué? Centralizar toda la lógica de sesión: login, logout, registro.
 * ¿Impacto? Todos los componentes dentro pueden acceder al usuario actual
 *            y disparar acciones de auth sin prop drilling.
 */
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserResponse | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  /**
   * login llama a la API de autenticación, guarda el access_token en memoria
   * y obtiene el perfil del usuario.
   * ¿Para qué? Un solo punto de entrada para iniciar sesión.
   * ¿Impacto? Si falla cualquier paso, el error se propaga al componente
   *            que llamó a login para mostrarlo al usuario.
   */
  const login = useCallback(async (data: LoginRequest) => {
    setIsLoading(true)
    try {
      const tokens = await authApi.login(data)
      setAccessToken(tokens.access_token)
      const me = await authApi.getMe()
      setUser(me)
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * logout limpia el token en memoria y resetea el estado del usuario.
   * ¿Para qué? Garantizar que no queden traces de la sesión anterior.
   * ¿Impacto? Después de logout, isAuthenticated será false y ProtectedRoute
   *            redirigirá al usuario al login.
   */
  const logout = useCallback(() => {
    setAccessToken(null)
    setUser(null)
  }, [])

  /**
   * register crea una cuenta nueva y retorna el mensaje del servidor.
   * ¿Para qué? El componente Register puede mostrar el mensaje de éxito
   *            (ej: "revisa tu email") y luego redirigir al login.
   * ¿Impacto? No realiza login automático post-registro porque el email
   *            debe verificarse primero.
   */
  const register = useCallback(async (data: RegisterRequest): Promise<string> => {
    setIsLoading(true)
    try {
      const response = await authApi.register(data)
      return response.message
    } finally {
      setIsLoading(false)
    }
  }, [])

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: user !== null,
        isLoading,
        login,
        logout,
        register,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

/**
 * useAuth hook para acceder al contexto de autenticación.
 * ¿Para qué? Proveer un acceso tipado y con validación al AuthContext.
 * ¿Impacto? Lanza un error descriptivo si se usa fuera del AuthProvider,
 *            facilitando el debug durante el desarrollo.
 */
export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth debe usarse dentro de AuthProvider')
  return ctx
}
