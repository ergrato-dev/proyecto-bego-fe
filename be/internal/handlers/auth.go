// Archivo: auth.go (internal/handlers)
// Descripción: Handlers HTTP para todos los endpoints de autenticación.
// ¿Para qué? Procesar requests HTTP: parsear body, validar datos, llamar al service
//            y serializar la respuesta JSON. Los handlers son delgados — la lógica va en services.
// ¿Impacto? Una validación faltante aquí puede permitir datos corruptos llegar
//            a la base de datos o exponer información sensible en los errores.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"nn-auth-system/internal/dto"
	"nn-auth-system/internal/middleware"
	"nn-auth-system/internal/services"
)

// AuthHandler agrupa los handlers de autenticación y sus dependencias.
// ¿Para qué? Inyectar el service y el validator una sola vez en lugar de
//            crear instancias en cada handler.
type AuthHandler struct {
	service  *services.AuthService
	validate *validator.Validate
}

// NewAuthHandler crea un nuevo handler de autenticación con sus dependencias.
func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Register maneja POST /api/v1/auth/register
// ¿Para qué? Crear un nuevo usuario en el sistema con email y contraseña.
// ¿Impacto? Retorna 201 con el usuario creado (sin contraseña) y envía email de verificación.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// ¿Qué? ShouldBindJSON parsea el body JSON y valida los tipos básicos.
	// ¿Para qué? Retornar 400 si el body no es JSON válido o los tipos no coinciden.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	// ¿Qué? Validar el struct con las reglas definidas en las etiquetas `validate`.
	// ¿Para qué? Verificar reglas de negocio: email válido, password mínimo 8 chars, etc.
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	resp, err := h.service.Register(&req)
	if err != nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login maneja POST /api/v1/auth/login
// ¿Para qué? Autenticar al usuario y retornar access_token + refresh_token.
// ¿Impacto? Los errores de credenciales son siempre genéricos para prevenir
//            que el atacante sepa si el email existe o la contraseña es incorrecta.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	resp, err := h.service.Login(&req, c.ClientIP())
	if err != nil {
		// ¿Qué? 401 Unauthorized para credenciales inválidas.
		// ¿Para qué? El código HTTP correcto para "no autenticado" es 401, no 403.
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken maneja POST /api/v1/auth/refresh
// ¿Para qué? Renovar el access token usando el refresh token sin pedir credenciales.
// ¿Impacto? Si el refresh token es inválido o expiró, el usuario debe hacer login.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	resp, err := h.service.RefreshToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ChangePassword maneja POST /api/v1/auth/change-password (requiere auth)
// ¿Para qué? Cambiar la contraseña del usuario autenticado verificando la contraseña actual.
// ¿Impacto? El user_id se obtiene del JWT (middleware RequireAuth) — no del body request.
//            Esto evita que un usuario cambie la contraseña de otro.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	// ¿Qué? Obtener el user_id del contexto de Gin, puesto por el middleware RequireAuth.
	// ¿Para qué? El user_id viene del JWT verificado — no se puede falsificar.
	userID, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "usuario no autenticado"})
		return
	}

	if err := h.service.ChangePassword(userID.(string), &req, c.ClientIP()); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Contraseña cambiada exitosamente"})
}

// ForgotPassword maneja POST /api/v1/auth/forgot-password
// ¿Para qué? Iniciar el flujo de recuperación de contraseña enviando un email.
// ¿Impacto? SIEMPRE retorna 200 con el mismo mensaje, sin importar si el email existe.
//            Esto previene user enumeration attacks.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	// ¿Qué? Ignorar el error — el service ya maneja el caso de email no encontrado internamente.
	// ¿Para qué? Respuesta idéntica para emails existentes y no existentes.
	_ = h.service.ForgotPassword(&req)

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Si el email está registrado, recibirás instrucciones para recuperar tu contraseña.",
	})
}

// ResetPassword maneja POST /api/v1/auth/reset-password
// ¿Para qué? Completar el flujo de recuperación estableciendo una nueva contraseña con el token del email.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	if err := h.service.ResetPassword(&req, c.ClientIP()); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Contraseña restablecida exitosamente"})
}

// VerifyEmail maneja POST /api/v1/auth/verify-email
// ¿Para qué? Activar el email del usuario cuando hace click en el enlace de verificación.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "datos de entrada inválidos"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: formatValidationErrors(err)})
		return
	}

	if err := h.service.VerifyEmail(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Email verificado exitosamente. Ya puedes iniciar sesión."})
}

// formatValidationErrors convierte los errores de validación en un mensaje legible.
// ¿Para qué? El cliente recibe un mensaje claro sobre qué campo falló la validación.
// ¿Impacto? No exponer los mensajes internos de validator para evitar revelar
//            detalles de la implementación al cliente.
func formatValidationErrors(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "error de validación"
	}

	// ¿Qué? Retornar el primer error de validación para mantener el mensaje conciso.
	// ¿Para qué? Un solo error claro es más útil que una lista técnica.
	if len(validationErrors) > 0 {
		field := validationErrors[0].Field()
		tag := validationErrors[0].Tag()
		return "campo '" + field + "' inválido: " + tag
	}
	return "error de validación en los datos enviados"
}
