/**
 * Archivo: main.tsx
 * Descripción: Punto de entrada de la aplicación React.
 * ¿Para qué? Inicializar React, i18n y montar el componente raíz en el DOM.
 * ¿Impacto? El import de './i18n' DEBE ir antes de App para garantizar
 *            que i18next esté inicializado antes de que cualquier componente
 *            llame a useTranslation().
 */
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './i18n'
import './index.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
)
