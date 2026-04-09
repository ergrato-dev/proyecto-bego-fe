// Archivo: password_reset_token.go
// Descripción: Modelo GORM para tokens de recuperación de contraseña.
// ¿Para qué? Almacenar tokens temporales que permiten resetear contraseñas vía email.
// ¿Impacto? Un token sin expiración o reutilizable es un vector de ataque —
//            siempre verificar expires_at y used antes de aceptar el token.

package models

import "time"

// PasswordResetToken representa un token de un solo uso para restablecer contraseña.
// ¿Para qué? Implementar el flujo "olvidé mi contraseña" de forma segura:
//
//	el usuario recibe un enlace con el token por email, lo usa una sola vez,
//	y el token queda marcado como usado.
//
// ¿Impacto? Sin este mecanismo, no habría forma de recuperar el acceso a una cuenta.
type PasswordResetToken struct {
	// ¿Qué? UUID como clave primaria del token.
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	// ¿Qué? Referencia al usuario dueño del token.
	// ¿Para qué? Saber a qué usuario se le debe cambiar la contraseña
	//            cuando se use el token.
	// ¿Impacto? Sin esta FK, no se podría asociar el token al usuario correcto.
	UserID string `gorm:"type:uuid;not null;index"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	// ¿Qué? El token en sí — string aleatorio único enviado por email.
	// ¿Para qué? El índice permite buscarlo rápidamente cuando el usuario lo envía.
	Token string `gorm:"type:varchar(255);uniqueIndex;not null"`

	// ¿Qué? Fecha y hora de expiración del token.
	// ¿Para qué? Limitar la ventana de uso — un token no debe ser válido indefinidamente.
	//            El flujo de recuperación expira en 1 hora.
	// ¿Impacto? Sin expiración, un token interceptado puede usarse en cualquier momento.
	ExpiresAt time.Time `gorm:"not null"`

	// ¿Qué? Marca si el token ya fue utilizado.
	// ¿Para qué? Garantizar que cada token solo resetee la contraseña una vez.
	// ¿Impacto? Sin este flag, el mismo enlace de email podría usarse múltiples veces.
	Used bool `gorm:"default:false;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}
