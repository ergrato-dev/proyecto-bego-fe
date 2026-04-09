// Archivo: user.go (internal/handlers)
// Descripción: Handlers HTTP para los endpoints de perfil de usuario.
// ¿Para qué? Proveer acceso al perfil del usuario autenticado.
// ¿Impacto? Toda ruta en este handler requiere el middleware RequireAuth —
//            sin JWT válido, no se retorna información de usuario.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nn-auth-system/internal/middleware"
	"nn-auth-system/internal/services"
)

// UserHandler agrupa los handlers de perfil de usuario.
type UserHandler struct {
	service *services.AuthService
}

// NewUserHandler crea un nuevo handler de usuario con el service inyectado.
func NewUserHandler(service *services.AuthService) *UserHandler {
	return &UserHandler{service: service}
}

// GetMe maneja GET /api/v1/users/me (requiere auth)
// ¿Para qué? Retornar el perfil del usuario actualmente autenticado.
// ¿Impacto? El user_id viene del JWT verificado por RequireAuth — no se acepta
//            ningún parámetro del cliente para determinar qué usuario se retorna.
//            Esto garantiza que cada usuario solo puede ver su propio perfil.
func (h *UserHandler) GetMe(c *gin.Context) {
	// ¿Qué? Obtener el user_id del contexto, puesto por el middleware RequireAuth.
	// ¿Para qué? Identificar al usuario sin depender del body o query params del request.
	userID, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
