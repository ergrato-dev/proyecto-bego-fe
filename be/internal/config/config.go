// Archivo: config.go
// Descripción: Configuración central del backend — carga variables de entorno y las expone como struct.
// ¿Para qué? Centralizar toda la configuración en un solo lugar con tipos fuertes.
// ¿Impacto? Sin este archivo, cada paquete leería os.Getenv() directamente,
//            generando código duplicado y difícil de testear.

package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración del backend leída desde variables de entorno.
// ¿Para qué? Proveer un único punto de acceso a la configuración, con tipos correctos.
// ¿Impacto? Facilita el testing (se puede inyectar una Config mock) y
//
//	evita errores de tipo al leer strings como números.
type Config struct {
	// Base de datos
	DatabaseURL string

	// JWT
	JWTSecret                  string
	JWTAccessTokenExpireMinutes int
	JWTRefreshTokenExpireDays  int

	// Servidor
	Port        string
	Environment string
	FrontendURL string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// Load carga las variables de entorno desde el archivo .env y devuelve la configuración.
// ¿Para qué? Inicializar la configuración al arrancar la aplicación.
// ¿Impacto? Si falta alguna variable obligatoria, la aplicación falla en el arranque
//
//	(fail fast) en lugar de fallar silenciosamente en tiempo de ejecución.
func Load() (*Config, error) {
	// ¿Qué? Intentar cargar el archivo .env si existe.
	// ¿Para qué? En desarrollo local, las variables vienen de .env.
	//            En producción, vienen del entorno del sistema (no de .env).
	// ¿Impacto? godotenv.Load no da error si el archivo no existe — es intencional.
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        getEnvWithDefault("PORT", "8000"),
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		FrontendURL: getEnvWithDefault("FRONTEND_URL", "http://localhost:5173"),
		SMTPHost:    getEnvWithDefault("SMTP_HOST", "localhost"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:   getEnvWithDefault("FROM_EMAIL", "noreply@nn-company.com"),
	}

	// ¿Qué? Parsear valores numéricos con validación.
	// ¿Para qué? Detectar configuración inválida al arrancar, no en tiempo de ejecución.
	var err error

	cfg.JWTAccessTokenExpireMinutes, err = strconv.Atoi(getEnvWithDefault("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", "15"))
	if err != nil {
		return nil, fmt.Errorf("JWT_ACCESS_TOKEN_EXPIRE_MINUTES debe ser un número: %w", err)
	}

	cfg.JWTRefreshTokenExpireDays, err = strconv.Atoi(getEnvWithDefault("JWT_REFRESH_TOKEN_EXPIRE_DAYS", "7"))
	if err != nil {
		return nil, fmt.Errorf("JWT_REFRESH_TOKEN_EXPIRE_DAYS debe ser un número: %w", err)
	}

	cfg.SMTPPort, err = strconv.Atoi(getEnvWithDefault("SMTP_PORT", "1025"))
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT debe ser un número: %w", err)
	}

	// ¿Qué? Validar variables obligatorias.
	// ¿Para qué? Fail fast — el servidor no debe arrancar con configuración incompleta.
	// ¿Impacto? Sin esta validación, errores de configuración aparecerían al primer request
	//            y serían difíciles de diagnosticar.
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate verifica que las variables obligatorias estén presentes.
func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL es obligatoria")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET es obligatoria")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET debe tener al menos 32 caracteres")
	}
	return nil
}

// IsDevelopment indica si el entorno es de desarrollo.
// ¿Para qué? Activar logs detallados y modo debug solo en desarrollo.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// getEnvWithDefault devuelve el valor de la variable de entorno o el valor por defecto.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
