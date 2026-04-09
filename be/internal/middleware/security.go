// Archivo: security.go (internal/middleware)
// Descripción: Middleware que agrega cabeceras de seguridad HTTP a todas las respuestas.
// ¿Para qué? Proteger al cliente (navegador) contra ataques comunes como XSS,
//            clickjacking y sniffing de contenido sin configurar nada en el frontend.
// ¿Impacto? Estas cabeceras son la primera línea de defensa del navegador.
//            Sin ellas, el sitio es vulnerable a múltiples vectores de ataque (OWASP A05).

package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders agrega las cabeceras de seguridad HTTP recomendadas por OWASP.
// ¿Para qué? Instruir al navegador sobre cómo manejar el contenido de forma segura.
// ¿Impacto? Cada cabecera mitiga un tipo específico de ataque:
//   - X-Content-Type-Options: previene MIME sniffing
//   - X-Frame-Options: previene clickjacking
//   - X-XSS-Protection: capa adicional contra XSS en navegadores legacy
//   - Referrer-Policy: controla qué información de referencia se envía
//   - Content-Security-Policy: la defensa más fuerte contra XSS
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ¿Qué? Previene que el navegador "adivine" el tipo MIME del contenido.
		// ¿Para qué? Evita ataques donde un archivo subido se ejecuta como script.
		c.Header("X-Content-Type-Options", "nosniff")

		// ¿Qué? Prohíbe que la página sea embebida en un iframe.
		// ¿Para qué? Previene clickjacking — el atacante no puede superponer frames engañosos.
		c.Header("X-Frame-Options", "DENY")

		// ¿Qué? Habilita la protección XSS del navegador (legacy, IE/Edge antiguos).
		// ¿Para qué? Capa adicional de protección para navegadores que no soportan CSP.
		c.Header("X-XSS-Protection", "1; mode=block")

		// ¿Qué? Controla qué información del referrer se envía en requests cross-origin.
		// ¿Para qué? Evitar filtrar tokens o paths sensibles en la URL referrer.
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// ¿Qué? Content Security Policy — define qué recursos puede cargar la página.
		// ¿Para qué? La defensa más efectiva contra XSS — restringe scripts a fuentes confiables.
		// ¿Impacto? Esta CSP permite solo resources del mismo origen (self).
		//            En producción, ajustar para incluir CDNs de fuentes (fonts, etc.).
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")

		// ¿Qué? Fuerza conexiones HTTPS por 1 año con preload.
		// ¿Para qué? Previene ataques de downgrade (HTTPS → HTTP) una vez que el browser lo conoce.
		// ¿Impacto? Solo agregar en producción con HTTPS. En desarrollo local no tiene efecto
		//            pero tampoco daña — el navegador no aplica HSTS en localhost.
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// ¿Qué? Controla qué APIs del navegador puede usar la página.
		// ¿Para qué? Deshabilitar por defecto APIs no necesarias (cámara, micrófono, geolocalización).
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}
