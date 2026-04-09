// Archivo: email.go (internal/utils)
// Descripción: Utilidades para el envío de emails transaccionales del sistema de auth.
// ¿Para qué? Proveer funciones reutilizables para enviar emails de verificación
//            y recuperación de contraseña a través del servidor SMTP configurado.
// ¿Impacto? En desarrollo usa Mailpit (localhost:1025) que captura los emails
//            sin enviarlos realmente. En producción debe cambiarse a un proveedor SMTP real.

package utils

import (
	"fmt"
	"net/smtp"

	// ¿Qué? jordan-wright/email — librería para construir y enviar emails con adjuntos y HTML.
	// ¿Para qué? Simplifica la construcción del mensaje MIME con cuerpo HTML y texto plano.
	"github.com/jordan-wright/email"
)

// EmailConfig contiene los parámetros de conexión SMTP para el envío de emails.
// ¿Para qué? Desacoplar las funciones de envío de la configuración global
//            para facilitar el testing y la inyección de dependencias.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SendVerificationEmail envía un email con el enlace de verificación de cuenta al usuario recién registrado.
// ¿Para qué? Confirmar que el usuario tiene acceso real al email proporcionado antes
//            de permitirle iniciar sesión, previniendo registros con emails falsos.
// ¿Impacto? Si este email no llega, el usuario no puede completar el registro.
//            En Mailpit (desarrollo) el email se puede ver en http://localhost:8025.
func SendVerificationEmail(cfg EmailConfig, toEmail, fullName, verificationToken, frontendURL string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, verificationToken)

	// ¿Qué? Construir el cuerpo HTML del email de verificación.
	// ¿Para qué? Proveer una experiencia de usuario clara con los pasos a seguir.
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <h2>Verifica tu cuenta</h2>
  <p>Hola <strong>%s</strong>,</p>
  <p>Gracias por registrarte en NN Auth System. Para activar tu cuenta, haz click en el botón:</p>
  <p style="text-align: center; margin: 30px 0;">
    <a href="%s"
       style="background: #2563eb; color: white; padding: 12px 24px;
              text-decoration: none; border-radius: 6px; font-weight: bold;">
      Verificar mi email
    </a>
  </p>
  <p>O copia y pega este enlace en tu navegador:</p>
  <p style="word-break: break-all; color: #2563eb;">%s</p>
  <p><strong>Este enlace expira en 24 horas.</strong></p>
  <hr>
  <p style="color: #6b7280; font-size: 12px;">
    Si no creaste esta cuenta, ignora este email.
  </p>
</body>
</html>`, fullName, verificationURL, verificationURL)

	return sendEmail(cfg, toEmail, "Verifica tu cuenta — NN Auth System", htmlBody)
}

// SendPasswordResetEmail envía un email con el enlace para restablecer la contraseña.
// ¿Para qué? Permitir la recuperación de acceso a usuarios que olvidaron su contraseña.
// ¿Impacto? El enlace expira en 1 hora para limitar la ventana de abuso.
//            El mensaje de respuesta al endpoint forgot-password es SIEMPRE el mismo
//            para no revelar si el email existe en el sistema (previene user enumeration).
func SendPasswordResetEmail(cfg EmailConfig, toEmail, fullName, resetToken, frontendURL string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <h2>Recuperar contraseña</h2>
  <p>Hola <strong>%s</strong>,</p>
  <p>Recibimos una solicitud para restablecer la contraseña de tu cuenta. Haz click en el botón:</p>
  <p style="text-align: center; margin: 30px 0;">
    <a href="%s"
       style="background: #2563eb; color: white; padding: 12px 24px;
              text-decoration: none; border-radius: 6px; font-weight: bold;">
      Restablecer contraseña
    </a>
  </p>
  <p>O copia y pega este enlace en tu navegador:</p>
  <p style="word-break: break-all; color: #2563eb;">%s</p>
  <p><strong>Este enlace expira en 1 hora.</strong></p>
  <hr>
  <p style="color: #6b7280; font-size: 12px;">
    Si no solicitaste este cambio, ignora este email. Tu contraseña no será modificada.
  </p>
</body>
</html>`, fullName, resetURL, resetURL)

	return sendEmail(cfg, toEmail, "Recuperar contraseña — NN Auth System", htmlBody)
}

// sendEmail es la función interna que construye y envía el email via SMTP.
// ¿Para qué? Centralizar la lógica de envío SMTP para que SendVerificationEmail
//            y SendPasswordResetEmail no dupliquen código.
// ¿Impacto? Si el servidor SMTP no está disponible, retorna error que debe
//            manejarse en el service (loggear pero no bloquear el flujo principal).
func sendEmail(cfg EmailConfig, to, subject, htmlBody string) error {
	e := email.NewEmail()
	e.From = cfg.From
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(htmlBody)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// ¿Qué? En Mailpit (desarrollo, puerto 1025) no se requiere autenticación SMTP.
	// ¿Para qué? Simplificar el setup en desarrollo. En producción se usa smtp.PlainAuth.
	// ¿Impacto? Nunca usar smtp.PlainAuth sin TLS en producción — expone credenciales en red.
	if cfg.Username == "" {
		if err := e.Send(addr, nil); err != nil {
			return fmt.Errorf("enviando email a %s: %w", to, err)
		}
		return nil
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err := e.Send(addr, auth); err != nil {
		return fmt.Errorf("enviando email a %s: %w", to, err)
	}
	return nil
}
