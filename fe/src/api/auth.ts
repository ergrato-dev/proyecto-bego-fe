/**
 * Archivo: auth.ts (api)
 * Descripción: Funciones de acceso a la API de autenticación del backend.
 * ¿Para qué? Encapsular las peticiones HTTP en funciones tipadas y reutilizables.
 * ¿Impacto? Si cambia un endpoint del backend, solo hay que actualizar aquí.
 */

import apiClient from './axios'
import type {
  LoginRequest,
  RegisterRequest,
  ChangePasswordRequest,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  VerifyEmailRequest,
  TokenResponse,
  UserResponse,
  RegisterResponse,
  MessageResponse,
} from '../types/auth'

/** login autentica al usuario y retorna los tokens JWT. */
export async function login(data: LoginRequest): Promise<TokenResponse> {
  const response = await apiClient.post<TokenResponse>('/auth/login', data)
  return response.data
}

/** register crea un nuevo usuario y retorna sus datos. */
export async function register(data: RegisterRequest): Promise<RegisterResponse> {
  const response = await apiClient.post<RegisterResponse>('/auth/register', data)
  return response.data
}

/**
 * getMe retorna el perfil del usuario autenticado.
 * ¿Para qué? Obtener los datos del usuario después del login.
 * ¿Impacto? Requiere un access token válido en el header Authorization.
 */
export async function getMe(): Promise<UserResponse> {
  const response = await apiClient.get<UserResponse>('/users/me')
  return response.data
}

/** changePassword cambia la contraseña del usuario autenticado. */
export async function changePassword(data: ChangePasswordRequest): Promise<MessageResponse> {
  const response = await apiClient.post<MessageResponse>('/auth/change-password', data)
  return response.data
}

/**
 * forgotPassword inicia el flujo de recuperación de contraseña.
 * ¿Para qué? El backend envía un email con enlace de reset.
 * ¿Impacto? Siempre retorna 200 — no revela si el email existe (anti-enumeration).
 */
export async function forgotPassword(data: ForgotPasswordRequest): Promise<MessageResponse> {
  const response = await apiClient.post<MessageResponse>('/auth/forgot-password', data)
  return response.data
}

/** resetPassword establece una nueva contraseña usando el token del email. */
export async function resetPassword(data: ResetPasswordRequest): Promise<MessageResponse> {
  const response = await apiClient.post<MessageResponse>('/auth/reset-password', data)
  return response.data
}

/** verifyEmail confirma el email del usuario con el token enviado al registrarse. */
export async function verifyEmail(data: VerifyEmailRequest): Promise<MessageResponse> {
  const response = await apiClient.post<MessageResponse>('/auth/verify-email', data)
  return response.data
}
