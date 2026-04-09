/**
 * Archivo: axios.ts (api)
 * Descripción: Instancia de Axios configurada con la URL base y gestión de tokens.
 * ¿Para qué? Centralizar la configuración HTTP — baseURL, headers y el token JWT.
 * ¿Impacto? Sin esta instancia, cada llamada a la API repetiría la URL base
 *            y la lógica de autenticación, violando el principio DRY.
 */

import axios from 'axios'

/**
 * apiClient es la instancia de Axios configurada para la API del backend.
 * ¿Para qué? Todas las peticiones HTTP del frontend pasan por esta instancia.
 */
const apiClient = axios.create({
  // ¿Qué? Leer la URL de la API desde la variable de entorno de Vite.
  // ¿Para qué? Permitir diferentes URLs en desarrollo vs producción sin cambiar código.
  // ¿Impacto? Si VITE_API_URL no está definida, usa el valor por defecto.
  baseURL: import.meta.env.VITE_API_URL ?? 'http://localhost:8000/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

/**
 * _accessToken almacena el JWT en memoria (NO en localStorage).
 * ¿Para qué? Seguridad — en producción, localStorage es vulnerable a ataques XSS.
 *            En memoria, el token desaparece al recargar la página (comportamiento deseado).
 * ¿Impacto? El usuario deberá hacer login nuevamente si recarga la página.
 *            En producción se usarían httpOnly cookies gestionadas por el servidor.
 */
let _accessToken: string | null = null

/** Almacena el access token en memoria. Llamar tras login exitoso. */
export function setAccessToken(token: string | null): void {
  _accessToken = token
}

/** Retorna el access token actual, o null si no hay sesión activa. */
export function getAccessToken(): string | null {
  return _accessToken
}

/**
 * Interceptor de request — inyecta el token JWT en cada petición.
 * ¿Para qué? No repetir la lógica de añadir el header Authorization
 *            en cada llamada individual a la API.
 * ¿Impacto? Si _accessToken es null, el header no se añade y las
 *            rutas protegidas retornan 401.
 */
apiClient.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default apiClient
