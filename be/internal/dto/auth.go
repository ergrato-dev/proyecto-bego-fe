// Archivo: auth.go (internal/dto)
// Descripción: Data Transfer Objects para los flujos de autenticación.
// ¿Para qué? Separar los datos que viajan por la API (requests/responses)
//            de los modelos internos de la base de datos (GORM structs).
// ¿Impacto? Sin DTOs, se corre el riesgo de exponer campos sensibles como
//            hashed_password en las respuestas JSON de la API.

package dto

// RegisterRequest contiene los datos necesarios para registrar un nuevo usuario.
// ¿Para qué? Validar y transportar los datos del formulario de registro.
// ¿Impacto? Las etiquetas `validate` garantizan que los datos sean correctos
//            antes de llegar a la capa de negocio o la base de datos.
type RegisterRequest struct {
	Email    string `json:"email"     validate:"required,email,max=255"`
	FullName string `json:"full_name"  validate:"required,min=2,max=255"`
	// ¿Qué? Validación de contraseña: mínimo 8 chars, al menos 1 mayúscula,
	//        1 minúscula y 1 dígito (custom validate "password_strength").
	// ¿Para qué? Garantizar contraseñas seguras antes de hashear con bcrypt.
	Password string `json:"password"   validate:"required,min=8,max=72"`
}

// LoginRequest contiene las credenciales para iniciar sesión.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ChangePasswordRequest contiene las contraseñas para el flujo de cambio cuando el usuario está logueado.
// ¿Para qué? Verificar la contraseña actual antes de permitir el cambio.
// ¿Impacto? Sin verificar la contraseña actual, cualquiera con el access token
//            podría cambiar la contraseña (ej.: si el token fue robado brevemente).
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"      validate:"required,min=8,max=72"`
}

// ForgotPasswordRequest contiene el email para iniciar la recuperación de contraseña.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest contiene el token y la nueva contraseña para completar la recuperación.
type ResetPasswordRequest struct {
	Token       string `json:"token"        validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// VerifyEmailRequest contiene el token de verificación de email.
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// RefreshTokenRequest contiene el refresh token para obtener un nuevo access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse es la respuesta enviada al cliente tras un login o refresh exitoso.
// ¿Para qué? Entregar los tokens al cliente de forma estandarizada.
// ¿Impacto? El campo token_type = "bearer" es estándar OAuth2 — el cliente
//            debe enviar: "Authorization: Bearer <access_token>".
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// UserResponse es la representación pública de un usuario en las respuestas de la API.
// ¿Para qué? Exponer solo los campos necesarios — nunca hashed_password.
// ¿Impacto? Omitir campos sensibles aquí es la última línea de defensa contra
//            filtraciones de datos por serialización incorrecta.
type UserResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	IsEmailVerified bool   `json:"is_email_verified"`
	IsActive        bool   `json:"is_active"`
}

// RegisterResponse es la respuesta tras un registro exitoso.
type RegisterResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

// MessageResponse es una respuesta genérica con solo un mensaje de texto.
// ¿Para qué? Responder a operaciones que no retornan datos (ej.: forgot-password,
//            change-password) sin revelar si el email existe o no.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse es la estructura de error estandarizada de la API.
type ErrorResponse struct {
	Error string `json:"error"`
}
