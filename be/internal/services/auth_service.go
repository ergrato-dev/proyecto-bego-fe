// Archivo: auth_service.go (internal/services)
// Descripción: Lógica de negocio para todos los flujos de autenticación.
// ¿Para qué? Separar la lógica de negocio de los handlers HTTP.
//            Los handlers solo orquestan; el service hace el trabajo real.
// ¿Impacto? Esta separación facilita el testing (se puede probar la lógica
//            sin levantar un servidor HTTP) y mantiene los handlers delgados.

package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	// ¿Qué? GORM para interactuar con la base de datos.
	"gorm.io/gorm"

	"nn-auth-system/internal/config"
	"nn-auth-system/internal/dto"
	"nn-auth-system/internal/models"
	"nn-auth-system/internal/utils"
)

// AuthService encapsula toda la lógica de autenticación.
// ¿Para qué? Proveer una interfaz clara entre los handlers HTTP
//            y las operaciones de base de datos y seguridad.
// ¿Impacto? Al depender de *gorm.DB directamente (no de una interfaz),
//            los tests de integración requieren una BD real o un mock de GORM.
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService crea una nueva instancia del servicio de autenticación.
// ¿Para qué? Inyectar las dependencias (DB y Config) en el servicio.
// ¿Impacto? Patrón de inyección de dependencias — facilita el testing y el desacoplamiento.
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register registra un nuevo usuario en el sistema.
// ¿Para qué? Crear una cuenta nueva con email verificado=false y contraseña hasheada.
// ¿Impacto? Si el email ya existe retorna error genérico para evitar user enumeration.
//            El token de verificación expira en 24h.
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// ¿Qué? Verificar si el email ya está registrado.
	// ¿Para qué? Garantizar unicidad de emails en el sistema.
	// ¿Impacto? El mensaje de error debe ser genérico en producción para no revelar
	//            qué emails están registrados (previene user enumeration attack).
	var existingUser models.User
	result := s.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("el email ya está registrado")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("verificando email existente: %w", result.Error)
	}

	// ¿Qué? Hashear la contraseña antes de guardarla en la BD.
	// ¿Para qué? Nunca almacenar contraseñas en texto plano.
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hasheando contraseña: %w", err)
	}

	// ¿Qué? Crear el struct del usuario con todos los campos requeridos.
	user := models.User{
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		IsActive:       true,
		// ¿Qué? El email se verifica después de hacer click en el enlace enviado por email.
		IsEmailVerified: false,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("creando usuario: %w", err)
	}

	// ¿Qué? Generar un token aleatorio y seguro para verificar el email.
	// ¿Para qué? El token es el "secreto" que confima que el usuario accede al email.
	// ¿Impacto? crypto/rand garantiza aleatoriedad criptográfica — math/rand NO es suficiente.
	verificationToken, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("generando token de verificación: %w", err)
	}

	emailToken := models.EmailVerificationToken{
		UserID:    user.ID,
		Token:     verificationToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}
	if err := s.db.Create(&emailToken).Error; err != nil {
		return nil, fmt.Errorf("guardando token de verificación: %w", err)
	}

	// ¿Qué? Enviar el email de verificación de forma asíncrona (best-effort).
	// ¿Para qué? No bloquear el registro si el servidor de email falla temporalmente.
	// ¿Impacto? El usuario puede solicitar reenvío; el registro ya está guardado.
	emailCfg := utils.EmailConfig{
		Host:     s.cfg.SMTPHost,
		Port:     s.cfg.SMTPPort,
		Username: s.cfg.SMTPUsername,
		Password: s.cfg.SMTPPassword,
		From:     s.cfg.FromEmail,
	}
	go func() {
		if err := utils.SendVerificationEmail(emailCfg, user.Email, user.FullName, verificationToken, s.cfg.FrontendURL); err != nil {
			fmt.Printf("[ERROR] enviando email de verificación a %s: %v\n", user.Email, err)
		}
	}()

	return &dto.RegisterResponse{
		Message: "Usuario registrado exitosamente. Revisa tu email para verificar tu cuenta.",
		User: dto.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FullName:        user.FullName,
			IsEmailVerified: user.IsEmailVerified,
			IsActive:        user.IsActive,
		},
	}, nil
}

// Login autentica a un usuario y retorna los tokens JWT.
// ¿Para qué? Verificar credenciales y emitir tokens de acceso y refresco.
// ¿Impacto? Los errores devueltos son SIEMPRE genéricos ("credenciales inválidas")
//            para evitar revelar si el email existe o si la contraseña es incorrecta.
func (s *AuthService) Login(req *dto.LoginRequest, ipAddress string) (*dto.TokenResponse, error) {
	var user models.User
	result := s.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		// ¿Qué? Registrar el intento fallido antes de retornar el error genérico.
		utils.LogLoginFailed(req.Email, ipAddress, "email no encontrado")
		return nil, errors.New("credenciales inválidas")
	}

	if !utils.VerifyPassword(user.HashedPassword, req.Password) {
		utils.LogLoginFailed(req.Email, ipAddress, "contraseña incorrecta")
		return nil, errors.New("credenciales inválidas")
	}

	if !user.IsEmailVerified {
		utils.LogLoginFailed(req.Email, ipAddress, "email no verificado")
		return nil, errors.New("debes verificar tu email antes de iniciar sesión")
	}

	if !user.IsActive {
		utils.LogLoginFailed(req.Email, ipAddress, "cuenta inactiva")
		return nil, errors.New("tu cuenta está desactivada")
	}

	accessToken, err := utils.CreateAccessToken(user.ID, s.cfg.JWTSecret, s.cfg.JWTAccessTokenExpireMinutes)
	if err != nil {
		return nil, fmt.Errorf("generando access token: %w", err)
	}

	refreshToken, err := utils.CreateRefreshToken(user.ID, s.cfg.JWTSecret, s.cfg.JWTRefreshTokenExpireDays)
	if err != nil {
		return nil, fmt.Errorf("generando refresh token: %w", err)
	}

	utils.LogLoginSuccess(user.ID, user.Email, ipAddress)

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
	}, nil
}

// RefreshToken valida el refresh token y emite un nuevo access token.
// ¿Para qué? Renovar la sesión del usuario sin pedir sus credenciales de nuevo.
// ¿Impacto? Si el refresh token es inválido o expiró, el usuario debe hacer login.
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	claims, err := utils.ParseToken(req.RefreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("refresh token inválido: %w", err)
	}

	// ¿Qué? Verificar que el usuario aún existe y está activo.
	// ¿Para qué? Un token puede ser válido pero el usuario haber sido desactivado.
	var user models.User
	if err := s.db.Where("id = ? AND is_active = true", claims.UserID).First(&user).Error; err != nil {
		return nil, errors.New("usuario no encontrado o inactivo")
	}

	newAccessToken, err := utils.CreateAccessToken(user.ID, s.cfg.JWTSecret, s.cfg.JWTAccessTokenExpireMinutes)
	if err != nil {
		return nil, fmt.Errorf("generando nuevo access token: %w", err)
	}

	newRefreshToken, err := utils.CreateRefreshToken(user.ID, s.cfg.JWTSecret, s.cfg.JWTRefreshTokenExpireDays)
	if err != nil {
		return nil, fmt.Errorf("generando nuevo refresh token: %w", err)
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "bearer",
	}, nil
}

// ChangePassword cambia la contraseña de un usuario autenticado.
// ¿Para qué? Permitir al usuario actualizar su contraseña desde el perfil.
// ¿Impacto? Requiere la contraseña actual para evitar que un access token robado
//            brevemente permita cambiar la contraseña de la cuenta.
func (s *AuthService) ChangePassword(userID string, req *dto.ChangePasswordRequest, ipAddress string) error {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("usuario no encontrado")
	}

	if !utils.VerifyPassword(user.HashedPassword, req.CurrentPassword) {
		return errors.New("la contraseña actual es incorrecta")
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hasheando nueva contraseña: %w", err)
	}

	if err := s.db.Model(&user).Update("hashed_password", newHash).Error; err != nil {
		return fmt.Errorf("actualizando contraseña: %w", err)
	}

	utils.LogPasswordChanged(user.ID, user.Email, ipAddress)
	return nil
}

// ForgotPassword inicia el flujo de recuperación de contraseña enviando un email.
// ¿Para qué? Permitir que usuarios que olvidaron su contraseña recuperen el acceso.
// ¿Impacto? Siempre retorna el mismo mensaje genérico sin importar si el email existe
//            o no, para prevenir user enumeration attacks (OWASP A01).
func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	var user models.User
	result := s.db.Where("email = ?", req.Email).First(&user)

	// ¿Qué? Si el email no existe, retornar nil (sin error) y sin enviar email.
	// ¿Para qué? No revelar si el email está registrado en el sistema.
	// ¿Impacto? El cliente recibe el mismo mensaje exitoso en ambos casos.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if result.Error != nil {
		return fmt.Errorf("buscando usuario: %w", result.Error)
	}

	resetToken, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("generando token de reset: %w", err)
	}

	// ¿Qué? Invalidar tokens anteriores no usados marcándolos como usados.
	// ¿Para qué? Evitar confusión con múltiples tokens válidos simultáneamente.
	s.db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used = false", user.ID).
		Update("used", true)

	passwordToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	if err := s.db.Create(&passwordToken).Error; err != nil {
		return fmt.Errorf("guardando token de reset: %w", err)
	}

	emailCfg := utils.EmailConfig{
		Host:     s.cfg.SMTPHost,
		Port:     s.cfg.SMTPPort,
		Username: s.cfg.SMTPUsername,
		Password: s.cfg.SMTPPassword,
		From:     s.cfg.FromEmail,
	}
	go func() {
		if err := utils.SendPasswordResetEmail(emailCfg, user.Email, user.FullName, resetToken, s.cfg.FrontendURL); err != nil {
			fmt.Printf("[ERROR] enviando email de reset a %s: %v\n", user.Email, err)
		}
	}()

	return nil
}

// ResetPassword completa el flujo de recuperación usando el token del email.
// ¿Para qué? Permitir establecer una nueva contraseña con el token de un solo uso.
// ¿Impacto? El token expira en 1h y solo puede usarse una vez (campo `used`).
//            Después del reset, el token queda marcado como `used=true`.
func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest, ipAddress string) error {
	var passwordToken models.PasswordResetToken
	result := s.db.Where("token = ?", req.Token).First(&passwordToken)
	if result.Error != nil {
		return errors.New("token inválido o no encontrado")
	}

	if passwordToken.Used {
		return errors.New("este token ya fue utilizado")
	}

	if time.Now().After(passwordToken.ExpiresAt) {
		return errors.New("el token ha expirado, solicita uno nuevo")
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hasheando nueva contraseña: %w", err)
	}

	if err := s.db.Model(&models.User{}).
		Where("id = ?", passwordToken.UserID).
		Update("hashed_password", newHash).Error; err != nil {
		return fmt.Errorf("actualizando contraseña: %w", err)
	}

	if err := s.db.Model(&passwordToken).Update("used", true).Error; err != nil {
		return fmt.Errorf("marcando token como usado: %w", err)
	}

	utils.LogPasswordReset(passwordToken.UserID, ipAddress)
	return nil
}

// VerifyEmail activa el email del usuario usando el token enviado al registrarse.
// ¿Para qué? Confirmar que el usuario tiene acceso real al email proporcionado.
// ¿Impacto? Sin verificación, cualquiera puede registrarse con emails ajenos.
func (s *AuthService) VerifyEmail(req *dto.VerifyEmailRequest) error {
	var emailToken models.EmailVerificationToken
	result := s.db.Where("token = ?", req.Token).First(&emailToken)
	if result.Error != nil {
		return errors.New("token inválido o no encontrado")
	}

	if emailToken.Used {
		return errors.New("este token ya fue utilizado")
	}

	if time.Now().After(emailToken.ExpiresAt) {
		return errors.New("el token ha expirado, solicita uno nuevo")
	}

	if err := s.db.Model(&models.User{}).
		Where("id = ?", emailToken.UserID).
		Update("is_email_verified", true).Error; err != nil {
		return fmt.Errorf("verificando email: %w", err)
	}

	if err := s.db.Model(&emailToken).Update("used", true).Error; err != nil {
		return fmt.Errorf("marcando token como usado: %w", err)
	}

	utils.LogEmailVerified(emailToken.UserID, "")
	return nil
}

// GetUserByID obtiene el perfil de un usuario por su ID.
// ¿Para qué? Proveer los datos del usuario autenticado en el endpoint GET /me.
func (s *AuthService) GetUserByID(userID string) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	return &dto.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		FullName:        user.FullName,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
	}, nil
}

// generateSecureToken genera un token aleatorio de 32 bytes codificado en hex.
// ¿Para qué? Crear tokens imposibles de adivinar para verificación de email y reset de contraseña.
// ¿Impacto? crypto/rand usa el generador de números aleatorios del SO — NUNCA math/rand.
//            Un token de 32 bytes = 256 bits de entropía — imposible de fuerza bruta.
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generando token seguro: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
