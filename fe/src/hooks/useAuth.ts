/**
 * Archivo: useAuth.ts
 * Descripción: Re-exportación del hook useAuth desde AuthContext.
 * ¿Para qué? Proveer un punto de importación más semántico:
 *            import { useAuth } from '../hooks/useAuth'
 *            en vez de importar directamente desde el contexto.
 * ¿Impacto? Bajo — solo organización. Facilita refactorizaciones futuras.
 */

export { useAuth } from '../context/AuthContext'
