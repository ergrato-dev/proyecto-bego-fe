// Archivo: audit_log.go (internal/utils)
// Descripción: Utilidades de auditoría para registrar eventos de autenticación.
// ¿Para qué? Mantener un registro de eventos de seguridad (logins exitosos, fallidos,
//            cambios de contraseña) para detección de anomalías y análisis forense.
// ¿Impacto? Sin auditoría, es imposible detectar ataques de fuerza bruta
//            ni investigar incidentes de seguridad a posteriori.

package utils

import (
	"log"
	"time"
)

// AuditEvent representa un evento de seguridad registrado en el sistema.
// ¿Para qué? Estructurar los logs de auditoría de forma consistente.
type AuditEvent struct {
	Timestamp time.Time
	Event     string
	UserID    string
	Email     string
	IPAddress string
	Details   string
}

// LogLoginSuccess registra un inicio de sesión exitoso.
// ¿Para qué? Detectar patrones de acceso sospechosos (ej.: login desde IP inusual).
func LogLoginSuccess(userID, email, ipAddress string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		Event:     "LOGIN_SUCCESS",
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
	}
	logEvent(event)
}

// LogLoginFailed registra un intento de inicio de sesión fallido.
// ¿Para qué? Detectar ataques de fuerza bruta — múltiples intentos fallidos en poco tiempo.
// ¿Impacto? Junto con el rate limiter, estos logs permiten identificar y bloquear atacantes.
func LogLoginFailed(email, ipAddress, reason string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		Event:     "LOGIN_FAILED",
		Email:     email,
		IPAddress: ipAddress,
		Details:   reason,
	}
	logEvent(event)
}

// LogPasswordChanged registra un cambio de contraseña exitoso.
// ¿Para qué? Alertar al usuario si ve un cambio de contraseña que no realizó.
func LogPasswordChanged(userID, email, ipAddress string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		Event:     "PASSWORD_CHANGED",
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
	}
	logEvent(event)
}

// LogPasswordReset registra un restablecimiento de contraseña vía email.
func LogPasswordReset(email, ipAddress string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		Event:     "PASSWORD_RESET",
		Email:     email,
		IPAddress: ipAddress,
	}
	logEvent(event)
}

// LogEmailVerified registra la verificación de un email.
func LogEmailVerified(userID, email string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		Event:     "EMAIL_VERIFIED",
		UserID:    userID,
		Email:     email,
	}
	logEvent(event)
}

// logEvent imprime el evento de auditoría en el log del servidor.
// ¿Para qué? En producción, estos logs deben enviarse a un sistema centralizado
//            (ej.: CloudWatch, Datadog) para análisis y alertas.
// ¿Impacto? NUNCA loggear contraseñas, tokens completos ni información sensible.
//            Solo loggear metadatos del evento (quién, qué, cuándo, desde dónde).
func logEvent(e AuditEvent) {
	if e.Details != "" {
		log.Printf("[AUDIT] %s | event=%s | user_id=%s | email=%s | ip=%s | details=%s",
			e.Timestamp.Format(time.RFC3339),
			e.Event, e.UserID, e.Email, e.IPAddress, e.Details,
		)
		return
	}
	log.Printf("[AUDIT] %s | event=%s | user_id=%s | email=%s | ip=%s",
		e.Timestamp.Format(time.RFC3339),
		e.Event, e.UserID, e.Email, e.IPAddress,
	)
}
