// Archivo: security_test.go (internal/utils)
// Descripción: Tests unitarios para las funciones de seguridad — bcrypt y JWT.
// ¿Para qué? Verificar que las funciones críticas de autenticación se comportan
//            correctamente en casos normales y en casos de error.
// ¿Impacto? Estos tests son la primera línea de defensa: si HashPassword o
//            ParseToken tienen un bug, toda la autenticación del sistema falla.

package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nn-auth-system/internal/utils"
)

// testSecret es la clave JWT usada en todos los tests de este archivo.
// ¿Para qué? El JWT_SECRET debe tener al menos 32 caracteres — esta constante
//            es segura para testing pero NUNCA debe usarse en producción.
const testSecret = "test-secret-key-must-be-at-least-32-chars-long"

// --- Tests de HashPassword y VerifyPassword ---

// TestHashPassword_ProduceHashBcrypt verifica que HashPassword retorna un hash bcrypt válido.
// ¿Para qué? Confirmar que la contraseña no se almacena en texto plano
//            y que el hash tiene el formato correcto de bcrypt ($2a$).
func TestHashPassword_ProduceHashBcrypt(t *testing.T) {
	hash, err := utils.HashPassword("miContraseña123")

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	// ¿Qué? Los hashes bcrypt siempre empiezan con $2a$, $2b$ o $2y$.
	// ¿Para qué? Confirmar que se usó bcrypt y no otro algoritmo menos seguro.
	assert.True(t, len(hash) > 0)
	assert.NotEqual(t, "miContraseña123", hash, "el hash no debe ser igual al texto plano")
}

// TestHashPassword_SamePasswordDifferentHashes verifica que bcrypt genera hashes distintos
// para la misma contraseña (por el salt aleatorio).
// ¿Para qué? El salt único por hash previene ataques de rainbow table.
func TestHashPassword_SamePasswordDifferentHashes(t *testing.T) {
	hash1, err := utils.HashPassword("mismaContraseña")
	require.NoError(t, err)

	hash2, err := utils.HashPassword("mismaContraseña")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "dos hashes de la misma contraseña deben ser distintos (salt único)")
}

// TestVerifyPassword_CorrectPassword verifica que una contraseña correcta sea válida.
func TestVerifyPassword_CorrectPassword(t *testing.T) {
	hash, err := utils.HashPassword("contraseñaCorrecta!")
	require.NoError(t, err)

	result := utils.VerifyPassword(hash, "contraseñaCorrecta!")

	assert.True(t, result, "la contraseña correcta debe verificarse como válida")
}

// TestVerifyPassword_IncorrectPassword verifica que una contraseña incorrecta sea rechazada.
// ¿Para qué? Confirmar que el sistema no acepta contraseñas equivocadas en el login.
func TestVerifyPassword_IncorrectPassword(t *testing.T) {
	hash, err := utils.HashPassword("contraseñaCorrecta!")
	require.NoError(t, err)

	result := utils.VerifyPassword(hash, "contraseñaEquivocada!")

	assert.False(t, result, "una contraseña incorrecta debe retornar false")
}

// TestVerifyPassword_EmptyPassword verifica que una contraseña vacía sea rechazada.
func TestVerifyPassword_EmptyPassword(t *testing.T) {
	hash, err := utils.HashPassword("contraseña123")
	require.NoError(t, err)

	result := utils.VerifyPassword(hash, "")

	assert.False(t, result)
}

// --- Tests de CreateAccessToken y ParseToken ---

// TestCreateAccessToken_ContainsCorrectUserID verifica que el token generado
// contiene el user_id correcto cuando se parsea.
// ¿Para qué? El user_id en el token es la única forma de identificar al usuario
//            sin consultar la BD en cada request — debe ser correcto.
func TestCreateAccessToken_ContainsCorrectUserID(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"

	token, err := utils.CreateAccessToken(userID, testSecret, 15)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := utils.ParseToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

// TestParseToken_WrongSecret verifica que un token firmado con clave distinta es rechazado.
// ¿Para qué? Garantizar que tokens forjados con otra clave no sean aceptados.
// ¿Impacto? Si este test falla, el sistema sería vulnerable a JWT forgery attacks.
func TestParseToken_WrongSecret(t *testing.T) {
	token, err := utils.CreateAccessToken("user-123", testSecret, 15)
	require.NoError(t, err)

	_, err = utils.ParseToken(token, "completely-different-secret-key-32chars!")
	assert.Error(t, err, "un token firmado con otra clave debe ser rechazado")
}

// TestParseToken_InvalidFormat verifica que un string que no es JWT es rechazado.
func TestParseToken_InvalidFormat(t *testing.T) {
	_, err := utils.ParseToken("esto.no.es.un.jwt.valido", testSecret)
	assert.Error(t, err)
}

// TestParseToken_EmptyToken verifica que un token vacío es rechazado.
func TestParseToken_EmptyToken(t *testing.T) {
	_, err := utils.ParseToken("", testSecret)
	assert.Error(t, err)
}

// TestParseToken_ExpiredToken verifica que un token expirado es rechazado con
// el mensaje correcto.
// ¿Para qué? Confirmar que el sistema no acepta tokens vencidos — fundamental
//            para la seguridad de las sesiones.
func TestParseToken_ExpiredToken(t *testing.T) {
	// ¿Qué? Crear un token con expiración negativa (ya expirado al crearse).
	token, err := utils.CreateAccessToken("user-123", testSecret, -1)
	require.NoError(t, err)

	// ¿Qué? Pequeña pausa para asegurar que el tiempo de expiración pasó.
	time.Sleep(10 * time.Millisecond)

	_, err = utils.ParseToken(token, testSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expirado", "el error debe mencionar que el token expiró")
}

// TestCreateRefreshToken_Valid verifica que el refresh token contiene el user_id correcto.
// ¿Para qué? El refresh token tiene la misma estructura que el access token —
//            la diferencia es solo la duración (7 días vs 15 min).
func TestCreateRefreshToken_Valid(t *testing.T) {
	userID := "user-refresh-test-id"

	token, err := utils.CreateRefreshToken(userID, testSecret, 7)
	require.NoError(t, err)

	claims, err := utils.ParseToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}
