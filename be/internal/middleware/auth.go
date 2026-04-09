// Archivo: auth.go (internal/middleware)
// Descripción: Middleware de autenticación JWT para proteger rutas privadas.
// ¿Para qué? Verificar que cada request a rutas protegidas incluya un JWT válido,
//            extrayendo el user_id para que los handlers lo usen sin consultar la BD.
// ¿Impacto? Sin este middleware, cualquiera podría acceder a rutas protegidas
//            como GET /me o POST /change-password sin estar autenticado.

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"nn-auth-system/internal/config"
	"nn-auth-system/internal/utils"
)

// ContextKeyUserID es la clave usada para almacenar el user_id en el contexto de Gin.
// ¿Para qué? Constante tipada para evitar typos al leer el user_id en los handlers.
const ContextKeyUserID = "user_id"

// RequireAuth es el middleware que verifica el JWT en el header Authorization.
// ¿Para qué? Proteger rutas que requieren autenticación.
// ¿Impacto? Si el token es inválido, expirado o ausente, retorna 401 y aborta
//            la cadena de handlers — el handler principal nunca se ejecuta.
func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ¿Qué? Leer el header Authorization con formato "Bearer <token>".
		// ¿Para qué? El estándar OAuth2/JWT usa este header para transportar el token.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "autorización requerida"})
			return
		}

		// ¿Qué? Validar el formato "Bearer <token>".
		// ¿Para qué? Evitar errores en ParseToken si el formato es incorrecto.
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido, usa: Bearer <token>"})
			return
		}

		tokenString := parts[1]

		// ¿Qué? Validar y decodificar el token JWT con la clave secreta del config.
		// ¿Para qué? Verificar firma, expiración y estructura del token.
		claims, err := utils.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// ¿Qué? Almacenar el user_id en el contexto de Gin para los handlers.
		// ¿Para qué? Los handlers pueden obtener el user_id sin re-parsear el token.
		// ¿Impacto? Usar la constante ContextKeyUserID evita errores silenciosos
		//            por inconsistencia en la clave string.
		c.Set(ContextKeyUserID, claims.UserID)

		// ¿Qué? Next() pasa el control al siguiente handler o middleware en la cadena.
		c.Next()
	}
}
