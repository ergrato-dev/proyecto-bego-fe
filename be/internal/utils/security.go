// Archivo: security.go (internal/utils)
// Descripción: Utilidades de seguridad — hashing de contraseñas y manejo de tokens JWT.
// ¿Para qué? Proveer funciones reutilizables de seguridad usadas en todo el sistema de auth.
// ¿Impacto? Es la base de la seguridad del sistema. Un error aquí compromete
//            toda la autenticación de la aplicación.

package utils

import (
	// ¿Qué? Paquete estándar para manejo de errores con contexto.
	"errors"
	"fmt"
	"time"

	// ¿Qué? golang-jwt/jwt/v5 — biblioteca oficial para JWT en Go.
	// ¿Para qué? Crear y verificar tokens JWT firmados con HS256.
	// ¿Impacto? v4 tenía CVE-2022-21703; v5 lo corrige con validación de audience obligatoria.
	"github.com/golang-jwt/jwt/v5"
	// ¿Qué? bcrypt de la biblioteca estándar extendida de Go.
	// ¿Para qué? Hashear contraseñas con un algoritmo de coste adaptativo — resistente a brute force.
	// ¿Impacto? Sin bcrypt, una filtración de BD expondría todas las contraseñas en texto plano.
	"golang.org/x/crypto/bcrypt"
)

// Claims define la estructura de datos embebida en el token JWT.
// ¿Para qué? Transportar el user_id de forma segura y verificable sin consultar la BD.
// ¿Impacto? Los campos aquí son visibles para quien tenga el token (no cifra, solo firma).
//            NUNCA incluir datos sensibles como contraseñas o roles administrativos.
type Claims struct {
	UserID string `json:"user_id"`
	// ¿Qué? RegisteredClaims incluye exp, iat, iss, sub — campos estándar JWT RFC 7519.
	jwt.RegisteredClaims
}

// HashPassword convierte la contraseña en texto plano a un hash bcrypt seguro.
// ¿Para qué? Almacenar contraseñas de forma segura — nunca en texto plano.
// ¿Impacto? bcrypt.DefaultCost (10) balancea seguridad y rendimiento; aumentar el cost
//            en producción si el hardware lo permite (recomendado: 12).
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword compara una contraseña en texto plano contra su hash bcrypt.
// ¿Para qué? Verificar que la contraseña proporcionada en el login corresponde al hash almacenado.
// ¿Impacto? bcrypt.CompareHashAndPassword es timing-safe — previene timing attacks.
//            Retorna nil si la contraseña es correcta, error si no coincide.
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CreateAccessToken genera un JWT de corta duración para el acceso a recursos protegidos.
// ¿Para qué? Autenticar requests del cliente sin consultar la BD en cada petición.
// ¿Impacto? Si el secret es débil o se filtra, cualquiera puede forjar tokens válidos.
//            La corta duración (15 min) limita la ventana de abuso ante una filtración.
func CreateAccessToken(userID string, secret string, expirationMinutes int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// ¿Qué? ExpiresAt define cuándo el token deja de ser válido.
			// ¿Para qué? Limitar el tiempo de abuso si el token es robado.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signing access token: %w", err)
	}
	return signed, nil
}

// CreateRefreshToken genera un JWT de larga duración solo para renovar el access token.
// ¿Para qué? Permitir sesiones persistentes sin pedir login frecuente,
//            manteniendo el access token de corta duración.
// ¿Impacto? Si el refresh token es robado, el atacante puede generar access tokens
//            por 7 días. Almacenarlo en httpOnly cookie reduce este riesgo.
func CreateRefreshToken(userID string, secret string, expirationDays int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signing refresh token: %w", err)
	}
	return signed, nil
}

// ParseToken valida y decodifica un JWT, retornando sus claims si es válido.
// ¿Para qué? Verificar la autenticidad e integridad del token en cada request protegido.
// ¿Impacto? Si el algoritmo aceptado no se valida explícitamente, es posible el ataque
//            de "alg:none". jwt/v5 ya lo mitiga pero la verificación del método es buena práctica.
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// ¿Qué? Validar explícitamente que el algoritmo sea HMAC (HS256).
		// ¿Para qué? Rechazar tokens con algoritmo "none" o RS256 con clave pública falsa.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de firma inesperado: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expirado")
		}
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims inválidos")
	}

	return claims, nil
}
