/**
 * Archivo: i18n.ts
 * Descripción: Configuración global de i18next para internacionalización (i18n).
 * ¿Para qué? Centralizar la configuración del motor de traducciones.
 *            Permite cambiar el idioma de toda la app sin recargar la página.
 * ¿Impacto? Sin esta configuración, useTranslation() de react-i18next falla
 *            y todas las cadenas de texto aparecerían como claves vacías.
 */

import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'

// ¿Qué? Importar los archivos JSON de traducción de cada idioma.
// ¿Para qué? Bundle estático — sin petición HTTP al servidor para cargar traducciones.
// ¿Impacto? Añadir un nuevo idioma requiere crear su JSON y agregarlo aquí.
import translationES from './locales/es/translation.json'
import translationEN from './locales/en/translation.json'

/** Recursos de traducción agrupados por código de idioma (ISO 639-1). */
const resources = {
  es: { translation: translationES },
  en: { translation: translationEN },
}

i18n
  // ¿Qué? LanguageDetector detecta automáticamente el idioma del navegador.
  // ¿Para qué? El usuario ve el idioma correcto sin configuración manual.
  // ¿Impacto? Detecta en: querystring (?lng=es), cookie, localStorage,
  //            navigator.language, htmlTag lang attribute.
  .use(LanguageDetector)
  // ¿Qué? Integración con React — hace disponible el hook useTranslation().
  .use(initReactI18next)
  .init({
    resources,
    // ¿Qué? Idioma por defecto si no se detecta el idioma del navegador.
    fallbackLng: 'es',
    // ¿Qué? Idiomas soportados por la aplicación.
    supportedLngs: ['es', 'en'],
    // ¿Qué? Desactivar interpolación por defecto para entornos React.
    // ¿Para qué? React ya escapa los valores — doble escape genera HTML incorrecto.
    interpolation: {
      escapeValue: false,
    },
    // ¿Qué? No autodetectar rutas de idioma (solo navigator y localStorage).
    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage'],
    },
  })

export default i18n
