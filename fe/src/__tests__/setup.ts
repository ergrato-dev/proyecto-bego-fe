/**
 * Archivo: setup.ts (__tests__)
 * Descripción: Configuración global del entorno de testing con Vitest.
 * ¿Para qué? Registrar los matchers extendidos de @testing-library/react
 *            para que aserciones como toBeInTheDocument() estén disponibles.
 * ¿Impacto? Sin este setup, los matchers de DOM no están disponibles en los tests
 *            y arrojan errores de "toBeInTheDocument is not a function".
 */

// ¿Qué? Importar matchers extendidos de @testing-library/jest-dom.
// ¿Para qué? Extiende expect() con aserciones como:
//   - expect(element).toBeInTheDocument()
//   - expect(element).toHaveValue('texto')
//   - expect(element).toBeDisabled()
import '@testing-library/react'
