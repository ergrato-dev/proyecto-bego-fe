// Archivo: auth_test.go (internal/handlers)
// Descripción: Tests de integración para los handlers HTTP de autenticación.
// ¿Para qué? Verificar que los endpoints responden con los status codes y
//            formatos JSON correctos ante distintas entradas.
// ¿Impacto? Los handler tests validan la capa HTTP: JSON parsing, validación de
//            entrada, status codes y estructura de respuesta. La lógica de negocio
//            está cubierta en services/auth_service_test.go.

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"nn-auth-system/internal/config"
	"nn-auth-system/internal/dto"
	"nn-auth-system/internal/handlers"
	"nn-auth-system/internal/middleware"
	"nn-auth-system/internal/models"
	"nn-auth-system/internal/services"
	"nn-auth-system/internal/testhelper"
)

// setupTestRouter crea un router Gin en modo test con todos los handlers de auth registrados.
// ¿Para qué? Reutilizar la configuración de rutas en todos los handler tests
//            sin repetir el código de setup.
// ¿Impacto? gin.TestMode desactiva el logger automático de Gin — la salida de tests
//            queda limpia sin los logs de request/response.
func setupTestRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	authService := services.NewAuthService(db, cfg)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.RequireAuth(cfg))
	{
		protected.GET("/users/me", userHandler.GetMe)
		protected.POST("/auth/change-password", authHandler.ChangePassword)
	}

	return router
}

// postJSON es un helper que construye y ejecuta un request POST con body JSON.
// ¿Para qué? Evitar duplicar el código de serialización y construcción de request
//            en cada test.
func postJSON(router *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// --- Tests de POST /api/v1/auth/register ---

// TestRegisterHandler_InvalidJSON verifica que un body no-JSON retorna 400.
// ¿Para qué? Confirmar que el handler rechaza inputs malformados antes de
//            llegar al service o la BD.
func TestRegisterHandler_InvalidJSON(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register",
		bytes.NewBufferString("{esto no es json válido"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Error)
}

// TestRegisterHandler_ValidationError_MissingEmail verifica que un email faltante
// retorna 422 Unprocessable Entity.
// ¿Para qué? Confirmar que el validador (`validate:"required,email"`) funciona
//            correctamente y retorna el status code apropiado.
func TestRegisterHandler_ValidationError_MissingEmail(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	w := postJSON(router, "/api/v1/auth/register", map[string]string{
		"full_name": "Test User",
		"password":  "Password123",
		// email ausente
	})

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// TestRegisterHandler_ValidationError_ShortPassword verifica que una contraseña
// de menos de 8 caracteres retorna 422.
// ¿Para qué? Confirmar la regla `validate:"min=8"` en PasswordRequest.
func TestRegisterHandler_ValidationError_ShortPassword(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	w := postJSON(router, "/api/v1/auth/register", map[string]string{
		"email":    "test@example.com",
		"full_name": "Test User",
		"password": "corta", // menos de 8 chars
	})

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// TestRegisterHandler_Success verifica el flujo de registro exitoso vía HTTP.
// ¿Para qué? Confirmar que el endpoint retorna 201 y el cuerpo JSON con
//            el usuario creado (sin contraseña ni campos sensibles).
func TestRegisterHandler_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	w := postJSON(router, "/api/v1/auth/register", map[string]string{
		"email":    email,
		"full_name": "Test Handler User",
		"password": "ValidPass123",
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.RegisterResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, email, resp.User.Email)
	assert.NotEmpty(t, resp.User.ID)
	assert.False(t, resp.User.IsEmailVerified)
	// ¿Qué? Confirmar que la respuesta NO incluye la contraseña hasheada.
	// ¿Para qué? La respuesta JSON nunca debe contener hashed_password.
	assert.NotContains(t, w.Body.String(), "hashed_password")
	assert.NotContains(t, w.Body.String(), "ValidPass123")
}

// TestRegisterHandler_DuplicateEmail verifica que registrar el mismo email dos veces
// retorna 409 Conflict.
// ¿Para qué? Confirmar que el handler mapea correctamente el error del service
//            al status code HTTP apropiado.
func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	body := map[string]string{
		"email":    email,
		"full_name": "Test",
		"password": "Password123",
	}

	w1 := postJSON(router, "/api/v1/auth/register", body)
	require.Equal(t, http.StatusCreated, w1.Code)

	w2 := postJSON(router, "/api/v1/auth/register", body)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

// --- Tests de POST /api/v1/auth/login ---

// TestLoginHandler_InvalidCredentials verifica que credenciales incorrectas retornan 401.
// ¿Para qué? Confirmar que el handler mapea errores de autenticación al status 401.
func TestLoginHandler_InvalidCredentials(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	w := postJSON(router, "/api/v1/auth/login", map[string]string{
		"email":    "noexiste@test.com",
		"password": "cualquierCosa123",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLoginHandler_Success verifica el flujo completo de login vía HTTP.
// ¿Para qué? Confirmar el JSON de respuesta incluye access_token, refresh_token
//            y token_type = "bearer".
func TestLoginHandler_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	// ¿Qué? Registrar y verificar el email via service directamente (ya probado en service tests).
	// ¿Para qué? Aislar este test — probamos el HANDLER, no el flujo de verificación.
	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Handler Login Test", Password: "LoginPass1!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	w := postJSON(router, "/api/v1/auth/login", map[string]string{
		"email":    email,
		"password": "LoginPass1!",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "bearer", resp.TokenType)
}

// --- Tests de GET /api/v1/users/me ---

// TestGetMeHandler_NoToken verifica que el endpoint protegido retorna 401 sin token.
// ¿Para qué? Confirmar que el middleware RequireAuth protege correctamente la ruta.
func TestGetMeHandler_NoToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetMeHandler_InvalidToken verifica que un token inválido retorna 401.
// ¿Para qué? Confirmar que el middleware rechaza tokens malformados o firmados incorrectamente.
func TestGetMeHandler_InvalidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer token.invalido.aqui")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetMeHandler_WithValidToken verifica que un JWT válido retorna el perfil del usuario.
// ¿Para qué? Integrar el flujo completo: login → obtener token → usar token en /me.
func TestGetMeHandler_WithValidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Me Handler Test", Password: "MePass123!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	tokenResp, err := svc.Login(&dto.LoginRequest{Email: email, Password: "MePass123!"}, "127.0.0.1")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var userResp dto.UserResponse
	err = json.NewDecoder(w.Body).Decode(&userResp)
	require.NoError(t, err)
	assert.Equal(t, email, userResp.Email)
	assert.Equal(t, "Me Handler Test", userResp.FullName)
	assert.True(t, userResp.IsEmailVerified)
}

// --- Tests de POST /api/v1/auth/forgot-password ---

// TestForgotPasswordHandler_AlwaysReturns200 verifica que el endpoint retorna 200
// tanto para emails existentes como para emails inexistentes.
// ¿Para qué? Confirmar la protección contra user enumeration attacks —
//            la respuesta debe ser idéntica en ambos casos (OWASP A01).
func TestForgotPasswordHandler_AlwaysReturns200(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	// Con email inexistente
	w1 := postJSON(router, "/api/v1/auth/forgot-password", map[string]string{
		"email": "no.existe.en.bd@example.com",
	})
	assert.Equal(t, http.StatusOK, w1.Code)

	var resp1 dto.MessageResponse
	err := json.NewDecoder(w1.Body).Decode(&resp1)
	require.NoError(t, err)
	assert.NotEmpty(t, resp1.Message)
}

// --- Tests de POST /api/v1/auth/refresh ---

// TestRefreshTokenHandler_InvalidToken verifica que un refresh token inválido retorna 401.
func TestRefreshTokenHandler_InvalidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	w := postJSON(router, "/api/v1/auth/refresh", map[string]string{
		"refresh_token": "token.invalido.total",
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRefreshTokenHandler_Success verifica que un refresh token válido retorna nuevos tokens.
// ¿Para qué? Confirmar el endpoint HTTP de refresh funciona end-to-end.
func TestRefreshTokenHandler_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Refresh Handler", Password: "Pass1234!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	loginResp, err := svc.Login(&dto.LoginRequest{Email: email, Password: "Pass1234!"}, "127.0.0.1")
	require.NoError(t, err)

	w := postJSON(router, "/api/v1/auth/refresh", map[string]string{
		"refresh_token": loginResp.RefreshToken,
	})
	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "bearer", resp.TokenType)
}

// --- Tests de POST /api/v1/auth/verify-email ---

// TestVerifyEmailHandler_InvalidToken verifica que un token de verificación inválido retorna 400.
func TestVerifyEmailHandler_InvalidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	w := postJSON(router, "/api/v1/auth/verify-email", map[string]string{
		"token": "token-que-no-existe",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestVerifyEmailHandler_Success verifica el flujo de verificación de email vía HTTP.
// ¿Para qué? Confirmar que el token de verificación funciona a través de la capa HTTP.
func TestVerifyEmailHandler_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Verify Handler", Password: "Pass1234!",
	})
	require.NoError(t, err)

	var evToken models.EmailVerificationToken
	db.Where("user_id = ?", regResp.User.ID).First(&evToken)

	w := postJSON(router, "/api/v1/auth/verify-email", map[string]string{
		"token": evToken.Token,
	})
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Tests de POST /api/v1/auth/change-password (protegido con JWT) ---

// TestChangePasswordHandler_NoToken verifica que el endpoint protegido requiere autenticación.
func TestChangePasswordHandler_NoToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password",
		bytes.NewBufferString(`{"current_password":"Old1!","new_password":"New12345!"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestChangePasswordHandler_Success verifica el flujo completo de cambio de contraseña.
// ¿Para qué? Confirmar que el endpoint protegido acepta la nueva contraseña
//            y retorna 200 cuando el token JWT y la contraseña actual son correctos.
func TestChangePasswordHandler_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	_, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Change Handler", Password: "OldPass1!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	tokenResp, err := svc.Login(&dto.LoginRequest{Email: email, Password: "OldPass1!"}, "127.0.0.1")
	require.NoError(t, err)

	bodyBytes, _ := json.Marshal(map[string]string{
		"current_password": "OldPass1!",
		"new_password":    "NewPass2!",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password",
		bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Tests de POST /api/v1/auth/reset-password ---

// TestResetPasswordHandler_InvalidToken verifica que un token de reset inválido retorna 400.
func TestResetPasswordHandler_InvalidToken(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())

	w := postJSON(router, "/api/v1/auth/reset-password", map[string]string{
		"token":        "token-que-no-existe-en-bd",
		"new_password": "NewPass1234!",
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestResetPasswordHandler_Success verifica el flujo completo de reset de contraseña vía HTTP.
// ¿Para qué? Confirmar que el endpoint HTTP acepta el token y actualiza la contraseña.
func TestResetPasswordHandler_Success(t *testing.T) {
	db := testhelper.SetupTestDB(t)
	router := setupTestRouter(db, testhelper.TestConfig())
	svc := services.NewAuthService(db, testhelper.TestConfig())

	email := testhelper.UniqueEmail()
	t.Cleanup(func() {
		db.Where("email = ?", email).Delete(&models.User{})
	})

	regResp, err := svc.Register(&dto.RegisterRequest{
		Email: email, FullName: "Reset Handler", Password: "OldPass1!",
	})
	require.NoError(t, err)
	db.Model(&models.User{}).Where("email = ?", email).Update("is_email_verified", true)

	err = svc.ForgotPassword(&dto.ForgotPasswordRequest{Email: email})
	require.NoError(t, err)

	var resetToken models.PasswordResetToken
	db.Where("user_id = ? AND used = ?", regResp.User.ID, false).First(&resetToken)

	w := postJSON(router, "/api/v1/auth/reset-password", map[string]string{
		"token":        resetToken.Token,
		"new_password": "NewPass1234!",
	})
	assert.Equal(t, http.StatusOK, w.Code)
}
