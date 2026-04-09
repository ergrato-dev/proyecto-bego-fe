/**
 * Archivo: auth.ts (types)
 * Descripción: Tipos TypeScript para todos los DTOs de autenticación.
 * ¿Para qué? Tener contratos explícitos entre el frontend y la API del backend,
 *            detectando discrepancias en tiempo de compilación (no en runtime).
 * ¿Impacto? Sin tipos, los errores de typo en nombres de campos solo se
 *            detectan en producción con una respuesta 422 inesperada.
 */

// --- Requests (datos que envía el frontend a la API) ---

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  full_name: string
  password: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface ForgotPasswordRequest {
  email: string
}

export interface ResetPasswordRequest {
  token: string
  new_password: string
}

export interface VerifyEmailRequest {
  token: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

// --- Responses (datos que retorna la API al frontend) ---

/** TokenResponse es la respuesta al hacer login o refresh. */
export interface TokenResponse {
  access_token: string
  refresh_token: string
  token_type: string
}

/** UserResponse es la representación pública del usuario (sin contraseña). */
export interface UserResponse {
  id: string
  email: string
  full_name: string
  is_email_verified: boolean
  is_active: boolean
}

/** RegisterResponse incluye un mensaje y los datos del usuario creado. */
export interface RegisterResponse {
  message: string
  user: UserResponse
}

/** MessageResponse es la respuesta genérica sin datos (ej: forgot-password). */
export interface MessageResponse {
  message: string
}

/**
 * ApiError es la estructura de error estandarizada de la API.
 * ¿Para qué? Parsear errores del backend con tipado correcto.
 */
export interface ApiError {
  error: string
}
