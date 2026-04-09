/**
 * Archivo: vite.config.ts
 * Descripción: Configuración de Vite — bundler, plugins de React y TailwindCSS,
 *              y configuración del test runner Vitest.
 * ¿Para qué? Centralizar la configuración de build, dev server y testing.
 * ¿Impacto? Sin esta configuración, TailwindCSS no procesa los estilos y Vitest
 *            no puede ejecutar los tests con soporte JSX/TypeScript.
 */
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    react(),
    // ¿Qué? Plugin oficial de TailwindCSS 4 para Vite.
    // ¿Para qué? Integra TailwindCSS directamente en el pipeline de Vite
    //            sin necesidad de PostCSS ni tailwind.config.js.
    tailwindcss(),
  ],
  test: {
    // ¿Qué? Configuración de Vitest para simular el entorno del navegador.
    // ¿Para qué? Los componentes React usan el DOM — jsdom lo simula en Node.js.
    // ¿Impacto? Sin 'jsdom', render() de @testing-library/react falla porque
    //            no existe document ni window en Node.js puro.
    environment: 'jsdom',
    globals: true,
    setupFiles: './src/__tests__/setup.ts',
  },
})
