// Archivo: ratelimit.go (internal/middleware)
// Descripción: Middleware de rate limiting por IP para endpoints de autenticación.
// ¿Para qué? Limitar la cantidad de requests por IP para prevenir ataques de
//            fuerza bruta contra los endpoints de login y registro.
// ¿Impacto? Sin rate limiting, un atacante puede intentar miles de contraseñas
//            por segundo. Con límite de 10 req/min, los ataques automatizados
//            quedan efectivamente frenados (OWASP A07 — Identification and Authentication Failures).

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// ¿Qué? ulule/limiter/v3 — rate limiter en memoria con soporte para múltiples stores.
	// ¿Para qué? Controlar la frecuencia de requests por IP sin necesitar Redis en desarrollo.
	// ¿Impacto? En producción con múltiples instancias, usar el store de Redis
	//            (github.com/ulule/limiter/v3/drivers/store/redis) para compartir el estado.
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// NewAuthRateLimiter crea un middleware de rate limiting para endpoints de autenticación.
// ¿Para qué? Aplicar un límite de requests por IP específico para rutas sensibles de auth.
// ¿Impacto? El límite recomendado para login es 5-10 intentos por minuto por IP.
//            Aumentarlo reduce la seguridad; reducirlo puede bloquear usuarios legítimos en redes NAT.
func NewAuthRateLimiter() gin.HandlerFunc {
	// ¿Qué? Rate: 10 requests por minuto por IP.
	// ¿Para qué? Balancear seguridad y usabilidad — 10 req/min permite uso normal
	//            pero frena ataques automatizados.
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	// ¿Qué? Store en memoria — adecuado para desarrollo o instancias únicas.
	store := memory.NewStore()

	// ¿Qué? Instancia del limiter con la tasa definida.
	rateLimiter := limiter.New(store, rate)

	return func(c *gin.Context) {
		// ¿Qué? Usar la IP del cliente como clave para el rate limiting.
		// ¿Para qué? Cada IP tiene su propio contador independiente.
		// ¿Impacto? c.ClientIP() respeta el header X-Forwarded-For si está configurado,
		//            lo cual es necesario si la app está detrás de un proxy/load balancer.
		ip := c.ClientIP()

		limiterCtx, err := rateLimiter.Get(c.Request.Context(), ip)
		if err != nil {
			// ¿Qué? Si el limiter falla, dejar pasar el request (fail open).
			// ¿Para qué? Preferir disponibilidad sobre restricción en caso de error interno.
			c.Next()
			return
		}

		// ¿Qué? Agregar headers de rate limit a la respuesta (estándar de la industria).
		// ¿Para qué? El cliente puede ver cuántos requests le quedan antes del bloqueo.
		c.Header("X-RateLimit-Limit", "10")
		c.Header("X-RateLimit-Remaining", formatInt(limiterCtx.Remaining))
		c.Header("X-RateLimit-Reset", formatUnix(limiterCtx.Reset))

		if limiterCtx.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "demasiados intentos, espera un minuto antes de volver a intentarlo",
			})
			return
		}

		c.Next()
	}
}

// formatInt convierte int64 a string para los headers HTTP.
func formatInt(n int64) string {
	return fmt.Sprintf("%d", n)
}

// formatUnix convierte el timestamp Unix a string para el header X-RateLimit-Reset.
func formatUnix(n int64) string {
	return fmt.Sprintf("%d", n)
}
