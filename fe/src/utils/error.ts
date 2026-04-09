/**
 * Archivo: error.ts (utils)
 * Descripción: Utilidad para extraer mensajes de error de respuestas de Axios.
 * ¿Para qué? Convertir errores genéricos de Axios en mensajes legibles para el usuario.
 * ¿Impacto? Sin esta función, los catch blocks de cada componente repetirían
 *            la misma lógica de parseo de errores.
 */

import axios from 'axios'
import type { ApiError } from '../types/auth'

/**
 * extractError extrae el mensaje de error de una excepción desconocida.
 * ¿Para qué? Normalizar errores de Axios (red, API) y errores de JS
 *            en un string simple para mostrar en la UI.
 * ¿Impacto? Los errores del backend tienen el formato { error: "mensaje" }.
 *            Los errores de red (sin conexión) muestran un mensaje genérico.
 */
export function extractError(err: unknown): string {
  if (axios.isAxiosError(err) && err.response?.data) {
    const data = err.response.data as ApiError
    if (data.error) return data.error
  }
  if (err instanceof Error) return err.message
  return 'Algo salió mal. Intenta de nuevo.'
}
