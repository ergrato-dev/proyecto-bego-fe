// Archivo: email_verification_token.go
// Descripción: Modelo GORM para tokens de verificación de email al registrarse.
// ¿Para qué? Confirmar que el usuario tiene acceso real al email que proporcionó.
// ¿Impacto? Sin verificación de email, cualquiera puede registrarse con emails ajenos
//            o inexistentes, habilitando spam y suplantación de identidad.

package models

import "time"

// EmailVerificationToken representa un token de un solo uso para verificar el email.
// ¿Para qué? Implementar el flujo de verificación post-registro:
//
//	al registrarse, se envía un email con el token; el usuario clickea el enlace
//	y su cuenta queda verificada (is_email_verified = true).
//
// ¿Impacto? Sin verificación, el login estaría abierto a cuentas con emails falsos.
type EmailVerificationToken struct {
	// ¿Qué? UUID como clave primaria del token.
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	// ¿Qué? Referencia al usuario cuyo email se está verificando.
	// ¿Para qué? Al usar el token, saber a qué usuario marcar como verificado.
	UserID string `gorm:"type:uuid;not null;index"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	// ¿Qué? El token enviado en el enlace de verificación por email.
	// ¿Para qué? El índice acelera la búsqueda cuando el usuario hace click en el enlace.
	Token string `gorm:"type:varchar(255);uniqueIndex;not null"`

	// ¿Qué? Fecha y hora de expiración — la verificación de email expira en 24 horas.
	// ¿Para qué? Si el usuario no verifica en 24h, debe solicitar un nuevo enlace.
	// ¿Impacto? Sin expiración, enlaces de verificación antiguos podrían activarse
	//            en cuentas que el usuario ya no controla.
	ExpiresAt time.Time `gorm:"not null"`

	// ¿Qué? Marca si el token ya fue utilizado para verificar el email.
	// ¿Para qué? Garantizar que cada enlace solo verifique una vez.
	Used bool `gorm:"default:false;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
}
