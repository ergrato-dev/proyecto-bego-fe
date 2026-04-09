// Archivo: auth_service_test.go (internal/services)
// Descripción: Tests de integración para el servicio de autenticación.
// ¿Para qué? Verificar que la lógica de negocio de auth funciona correctamente
//            contra una base de datos real — incluyendo hashing, JWT y validaciones.
// ¿Impacto? Estos tests necesitan PostgreSQL corriendo (docker compose up -d).
//            Si no hay BD disponible, los tests se omiten automáticamente (t.Skip).

package services_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nn-auth-system/internal/dto"
	"nn-auth-system/internal/models"
	"nn-auth-system/internal/services"
	"nn-auth-system/internal/testhelper"
)

// --- Tests de Register ---

// TestRegister_Success verifica el flujo de registro exitoso.
// ¿Para qué? Confirmar que un nuevo usuario se crea en BD con los datos correctos
//            y que la contraseña queda hasheada (nunca en texto plano).
func TestRegister_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		// ¿Qué? Borrar el usuario de test al finalizar — ON DELETE CASCADE limpia los tokens.
		db.Where("email = ?", email).Delete(&models.User{})
	})

	resp, err := svc.Register(&dto.RegisterRequest{
		Email:    email,
		FullName: "Usuario de Test",
		Password: "Password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.User.ID)
	assert.Equal(t, email, resp.User.Email)
	assert.Equal(t, "Usuario de Test", resp.User.FullName)
	// ¿Qué? El email no debe estar verificado al registrarse.
	// ¿Para qué? Confirmar que el flujo de verificación es necesario antes del login.
	assert.False(t, resp.User.IsEmailVerified)
	assert.True(t, resp.User.IsActive)
}

// TestRegister_DuplicateEmail verifica que no se puede registrar dos veces el mismo email.
// ¿Para qué? El índice UNIQUE en la columna email de la BD lo previene,
//            pero el service debe retornar un error claro (no un panic de DB).
func TestRegister_DuplicateEmail(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	req := &dto.RegisterRequest{
		Email:    email,
		FullName: "Test User",
		Password: "Password123",
	}

	_, err := svc.Register(req)
	require.NoError(t, err, "el primer registro debe ser exitoso")

	_, err = svc.Register(req)
	assert.Error(t, err, "el segundo registro con el mismo email debe fallar")
	assert.Contains(t, err.Error(), "ya está registrado")
}

// --- Tests de Login ---

// TestLogin_UnverifiedEmail verifica que no se puede hacer login sin verificar el email.
// ¿Para qué? Confirmar que el flujo de verificación de email es obligatorio —
//            usuarios con email no verificado deben ser bloqueados.
func TestLogin_UnverifiedEmail(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Test", Password: "Pass1234!",
	})
	require.NoError(t, err)

	_, err = svc.Login(&dto.LoginRequest{Email: email, Password: "Pass1234!"}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verificar tu email")
}

// TestLogin_WrongPassword verifica que una contraseña incorrecta es rechazada.
// ¿Para qué? Confirmar que bcrypt.CompareHashAndPassword detecta contraseñas equivocadas.
// ¿Impacto? El error retornado debe ser genérico ("credenciales inválidas") —
//            sin revelar si el email existe o si la contraseña es la que falla.
func TestLogin_WrongPassword(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Test", Password: "CorrectPass1!",
	})
	require.NoError(t, err)
	// ¿Qué? Verificar email directamente en BD para aislar el test del flujo de email.
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	_, err = svc.Login(&dto.LoginRequest{Email: email, Password: "WrongPass1!"}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credenciales inválidas")
}

// TestLogin_Success verifica el flujo de login exitoso: retorna access y refresh tokens.
// ¿Para qué? Confirmar que los tokens JWT son generados y tienen el formato correcto.
func TestLogin_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Login Test", Password: "MyPass1234!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	resp, err := svc.Login(&dto.LoginRequest{Email: email, Password: "MyPass1234!"}, "127.0.0.1")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "bearer", resp.TokenType)
	// ¿Qué? Los tokens deben tener 3 partes separadas por . (formato JWT header.payload.signature).
	assert.Contains(t, resp.AccessToken, ".", "el access token debe tener formato JWT")
}

// --- Tests de ChangePassword ---

// TestChangePassword_Success verifica que el cambio de contraseña funciona
// y que la nueva contraseña es válida para el siguiente login.
func TestChangePassword_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Change Pass", Password: "OldPassword1!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	err = svc.ChangePassword(regResp.User.ID, &dto.ChangePasswordRequest{
		CurrentPassword: "OldPassword1!",
		NewPassword:     "NewPassword2!",
	}, "127.0.0.1")
	require.NoError(t, err)

	// ¿Qué? Verificar que la nueva contraseña permite el login.
	// ¿Para qué? Confirmar que el hash se actualizó correctamente en la BD.
	_, err = svc.Login(&dto.LoginRequest{Email: email, Password: "NewPassword2!"}, "127.0.0.1")
	require.NoError(t, err, "el login con la nueva contraseña debe funcionar")

	// ¿Qué? Verificar que la contraseña antigua ya no funciona.
	_, err = svc.Login(&dto.LoginRequest{Email: email, Password: "OldPassword1!"}, "127.0.0.1")
	assert.Error(t, err, "la contraseña antigua no debe funcionar después del cambio")
}

// TestChangePassword_WrongCurrentPassword verifica que no se puede cambiar la contraseña
// sin proveer la contraseña actual correcta.
// ¿Para qué? Evitar que un access token robado brevemente permita cambiar
//            la contraseña de la cuenta.
func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Test", Password: "ActualPass1!",
	})
	require.NoError(t, err)

	err = svc.ChangePassword(regResp.User.ID, &dto.ChangePasswordRequest{
		CurrentPassword: "PasswordEquivocada1!",
		NewPassword:     "NuevaPass2!",
	}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contraseña actual es incorrecta")
}

// --- Tests de ForgotPassword ---

// TestForgotPassword_ExistingEmail verifica que ForgotPassword retorna nil
// para un email existente (y genera el token en BD aunque no lo retorne).
// ¿Para qué? Confirmar que el flujo no revela información sobre si el email existe.
func TestForgotPassword_ExistingEmail(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Forgot Test", Password: "Pass1234!",
	})
	require.NoError(t, err)

	// ¿Qué? ForgotPassword siempre retorna nil — tanto si el email existe como si no.
	// ¿Para qué? No revelar al atacante si el email está registrado o no.
	err = svc.ForgotPassword(&dto.ForgotPasswordRequest{Email: email})
	assert.NoError(t, err)
}

// TestForgotPassword_NonExistingEmail verifica que ForgotPassword retorna nil
// incluso cuando el email no existe en la BD.
// ¿Para qué? Garantizar que la respuesta es idéntica para emails existentes y no existentes
//            — previene user enumeration attacks (OWASP A01).
func TestForgotPassword_NonExistingEmail(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	err := svc.ForgotPassword(&dto.ForgotPasswordRequest{
		Email: "email-que-no-existe@test.com",
	})
	assert.NoError(t, err, "ForgotPassword debe retornar nil aunque el email no exista")
}

// --- Tests de VerifyEmail ---

// TestVerifyEmail_Success verifica el flujo de verificación de email.
// ¿Para qué? Confirmar que el token creado al registrarse activa correctamente
//            el campo is_email_verified del usuario en la BD.
// ¿Impacto? Sin este test, podría romperse el flujo de activación de cuenta
//            sin que lo detectemos hasta producción.
func TestVerifyEmail_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Verify Email Test", Password: "Pass1234!",
	})
	require.NoError(t, err)

	// ¿Qué? Recuperar el token de verificación directamente de la BD.
	// ¿Para qué? En tests, no tenemos el email real — accedemos al token vía BD.
	var evToken models.EmailVerificationToken
	err = db.Where("user_id = ?", regResp.User.ID).First(&evToken).Error
	require.NoError(t, err, "el token de verificación debe existir en BD tras el registro")

	err = svc.VerifyEmail(&dto.VerifyEmailRequest{Token: evToken.Token})
	require.NoError(t, err)

	// Verificar que is_email_verified = true en la BD.
	var user models.User
	db.Where("email = ?", email).First(&user)
	assert.True(t, user.IsEmailVerified, "el email debe estar verificado tras usar el token")
}

// TestVerifyEmail_InvalidToken verifica que un token inexistente retorna error.
// ¿Para qué? Confirmar que no se puede verificar un email con un token falso.
func TestVerifyEmail_InvalidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	err := svc.VerifyEmail(&dto.VerifyEmailRequest{Token: "token-que-no-existe-en-bd"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inválido")
}

// --- Tests de RefreshToken ---

// TestRefreshToken_Success verifica que un refresh token válido genera un nuevo access token.
// ¿Para qué? Confirmar el flujo de renovación de sesión sin pedir credenciales de nuevo.
// ¿Impacto? Si este flujo falla, los usuarios son deslogueados cada 15 minutos.
func TestRefreshToken_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Refresh Test", Password: "Pass1234!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	loginResp, err := svc.Login(&dto.LoginRequest{Email: email, Password: "Pass1234!"}, "127.0.0.1")
	require.NoError(t, err)

	newTokens, err := svc.RefreshToken(&dto.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
	assert.Equal(t, "bearer", newTokens.TokenType)
	// ¿Qué? Verificar que los tokens son strings JWT válidos (tres segmentos separados por '.').
	// ¿Para qué? Confirmar formato correcto sin depender del timing (iat/exp con resolución de 1s).
	assert.Equal(t, 3, len(strings.Split(newTokens.AccessToken, ".")),
		"el access token debe tener formato JWT (header.payload.signature)")
}

// TestRefreshToken_InvalidToken verifica que un refresh token inválido retorna error.
// ¿Para qué? Confirmar que el endpoint de refresh rechaza tokens falsos.
func TestRefreshToken_InvalidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	_, err := svc.RefreshToken(&dto.RefreshTokenRequest{
		RefreshToken: "token.invalido.completamente",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token inválido")
}

// --- Tests de ResetPassword ---

// TestResetPassword_Success verifica el flujo completo de recuperación de contraseña.
// ¿Para qué? Confirmar que: forgot-password crea token → reset-password lo consume
//            y actualiza la contraseña → el token queda marcado como used=true.
// ¿Impacto? Si este flujo falla, los usuarios no pueden recuperar sus cuentas.
func TestResetPassword_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Reset Pass Test", Password: "OldPass1!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	// ¿Qué? Triggering forgot-password guarda el token en BD de forma sincrónica.
	// ¿Para qué? El email se envía por goroutine (best-effort), pero el token se guarda antes.
	err = svc.ForgotPassword(&dto.ForgotPasswordRequest{Email: email})
	require.NoError(t, err)

	var resetToken models.PasswordResetToken
	err = db.Where("user_id = ? AND used = ?", regResp.User.ID, false).First(&resetToken).Error
	require.NoError(t, err, "el token de reset debe existir en BD tras forgot-password")

	err = svc.ResetPassword(&dto.ResetPasswordRequest{
		Token:       resetToken.Token,
		NewPassword: "NewPass2!",
	}, "127.0.0.1")
	require.NoError(t, err)

	// Verificar que la nueva contraseña funciona para login.
	_, err = svc.Login(&dto.LoginRequest{Email: email, Password: "NewPass2!"}, "127.0.0.1")
	require.NoError(t, err, "el login debe funcionar con la nueva contraseña")

	// Verificar que el token fue marcado como usado.
	db.Where("id = ?", resetToken.ID).First(&resetToken)
	assert.True(t, resetToken.Used, "el token de reset debe estar marcado como usado")
}

// TestResetPassword_UsedToken verifica que un token ya utilizado no puede usarse de nuevo.
// ¿Para qué? Prevenir que un token interceptado se use múltiples veces.
// ¿Impacto? Sin este control, comprometer el email da acceso permanente para resetear.
func TestResetPassword_UsedToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Used Token Test", Password: "OldPass1!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	err = svc.ForgotPassword(&dto.ForgotPasswordRequest{Email: email})
	require.NoError(t, err)

	var resetToken models.PasswordResetToken
	db.Where("user_id = ? AND used = ?", regResp.User.ID, false).First(&resetToken)

	// Usar el token por primera vez
	err = svc.ResetPassword(&dto.ResetPasswordRequest{
		Token: resetToken.Token, NewPassword: "NewPass1!",
	}, "127.0.0.1")
	require.NoError(t, err)

	// Intentar usar el mismo token de nuevo
	err = svc.ResetPassword(&dto.ResetPasswordRequest{
		Token: resetToken.Token, NewPassword: "AnotherPass1!",
	}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ya fue utilizado")
}

// --- Tests de GetUserByID ---

// TestGetUserByID_Success verifica que se puede obtener un usuario por su ID.
// ¿Para qué? Confirmar que el endpoint GET /users/me puede obtener el perfil
//            del usuario autenticado a partir de su user_id del JWT.
func TestGetUserByID_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "GetByID Test", Password: "Pass1234!",
	})
	require.NoError(t, err)

	userResp, err := svc.GetUserByID(regResp.User.ID)
	require.NoError(t, err)
	assert.Equal(t, regResp.User.ID, userResp.ID)
	assert.Equal(t, email, userResp.Email)
	assert.Equal(t, "GetByID Test", userResp.FullName)
}

// TestGetUserByID_NotFound verifica que buscar un ID inexistente retorna error.
// ¿Para qué? Confirmar que el endpoint maneja correctamente IDs inválidos.
func TestGetUserByID_NotFound(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	_, err := svc.GetUserByID("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}

// --- Tests de rutas de error adicionales ---

// TestLogin_InactiveUser verifica que una cuenta desactivada no puede hacer login.
// ¿Para qué? Confirmar que la desactivación de cuenta bloquea el acceso.
// ¿Impacto? Si no se revisa is_active, cuentas baneadas podrían seguir accediendo.
func TestLogin_InactiveUser(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Inactive User", Password: "Pass1234!",
	})
	require.NoError(t, err)

	// ¿Qué? Marcar el usuario como verificado pero desactivado.
	db.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"is_email_verified": true,
		"is_active":         false,
	})

	_, err = svc.Login(&dto.LoginRequest{Email: email, Password: "Pass1234!"}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "desactivada")
}

// TestVerifyEmail_AlreadyUsedToken verifica que un token de verificación ya utilizado retorna error.
// ¿Para qué? Confirmar que los tokens de verificación son de un solo uso.
// ¿Impacto? Sin este control, interceptar un enlace de verificación permite
//            reactivar una cuenta que fue marcada como no verificada.
func TestVerifyEmail_AlreadyUsedToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Reuse Token Test", Password: "Pass1234!",
	})
	require.NoError(t, err)

	var evToken models.EmailVerificationToken
	db.Where("user_id = ?", regResp.User.ID).First(&evToken)

	// Usar el token por primera vez
	err = svc.VerifyEmail(&dto.VerifyEmailRequest{Token: evToken.Token})
	require.NoError(t, err)

	// Intentar usar el mismo token de nuevo
	err = svc.VerifyEmail(&dto.VerifyEmailRequest{Token: evToken.Token})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ya fue utilizado")
}

// TestResetPassword_ExpiredToken verifica que un token de reset expirado retorna error.
// ¿Para qué? Confirmar que la ventana de 1 hora para el reset se cumple.
// ¿Impacto? Sin expiración, un enlace interceptado podría usarse meses después.
func TestResetPassword_ExpiredToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Expired Token Test", Password: "Pass1234!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	err = svc.ForgotPassword(&dto.ForgotPasswordRequest{Email: email})
	require.NoError(t, err)

	var resetToken models.PasswordResetToken
	db.Where("user_id = ? AND used = ?", regResp.User.ID, false).First(&resetToken)

	// ¿Qué? Forzar la expiración del token directamente en BD.
	// ¿Para qué? Simular que el token fue creado hace más de 1 hora.
	db.Model(&resetToken).Update("expires_at", "2000-01-01 00:00:00")

	err = svc.ResetPassword(&dto.ResetPasswordRequest{
		Token: resetToken.Token, NewPassword: "NewPass1!",
	}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expirado")
}

// TestChangePassword_UserNotFound verifica que cambiar contraseña de un ID inexistente retorna error.
func TestChangePassword_UserNotFound(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	svc := services.NewAuthService(db, testhelper.TestConfig())

	err := svc.ChangePassword("00000000-0000-0000-0000-000000000000", &dto.ChangePasswordRequest{
		CurrentPassword: "Old1!", NewPassword: "New12345!",
	}, "127.0.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no encontrado")
}
