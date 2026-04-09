// Archivo: testhelper.go (internal/testhelper)
// Descripción: Utilidades compartidas para los tests de integración del backend.
// ¿Para qué? Centralizar la configuración de la BD de test y helpers de datos
//            para que los tests no dupliquen código de setup/teardown.
// ¿Impacto? Sin este helper, cada archivo de test tendría su propio código
//            de conexión a BD — duplicación que dificulta el mantenimiento.

package testhelper

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nn-auth-system/internal/config"
)

// SetupTestDB establece una conexión a PostgreSQL para los tests de integración.
// ¿Para qué? Proveer una *gorm.DB lista para usar en cada test, cargando DATABASE_URL
//            desde el .env del proyecto o del entorno del sistema.
// ¿Impacto? Si DATABASE_URL no está disponible, el test se omite (t.Skip) en lugar
//            de fallar — los tests de integración son opcionales en entornos sin BD.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// ¿Qué? Intentar cargar el .env desde múltiples niveles de directorio.
	// ¿Para qué? Los tests se ejecutan desde su propio directorio de paquete,
	//            no desde la raíz del proyecto. Buscamos be/.env desde internal/*/.
	for _, path := range []string{"../../.env", "../../../.env", ".env"} {
		if err := godotenv.Load(path); err == nil {
			break
		}
		if os.Getenv("DATABASE_URL") != "" {
			break
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL no configurada — omitiendo test de integración (ejecuta: docker compose up -d)")
	}

	// ¿Qué? Abrir conexión GORM con logger silenciado para no contaminar la salida de tests.
	// ¿Para qué? Los logs SQL en los tests dificultan leer los resultados.
	// ¿Impacto? Si se necesita debuggear un test, cambiar LogMode a logger.Info temporalmente.
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("no se pudo conectar a la BD de testing: %v", err)
	}

	// ¿Qué? Verificar conectividad con un ping real a la BD.
	// ¿Para qué? Si la BD está configurada pero inaccesible, Skip en lugar de error críptico.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("error obteniendo sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("BD de testing no accesible: %v — ¿está corriendo 'docker compose up -d'?", err)
	}

	return db
}

// TestConfig retorna una configuración mínima válida para usar en tests.
// ¿Para qué? Evitar que cada test file defina su propia config de prueba.
// ¿Impacto? El JWT_SECRET tiene un valor fijo — solo para testing, nunca en producción.
func TestConfig() *config.Config {
	return &config.Config{
		JWTSecret:                   "test-secret-key-minimum-32-characters-ok!",
		JWTAccessTokenExpireMinutes: 15,
		JWTRefreshTokenExpireDays:   7,
		SMTPHost:                    getEnvWithDefault("SMTP_HOST", "localhost"),
		SMTPPort:                    1025,
		FromEmail:                   "test@nn-company.com",
		FrontendURL:                 "http://localhost:5173",
		Environment:                 "test",
		Port:                        "8001",
	}
}

// UniqueEmail genera un email único basado en el timestamp para evitar colisiones entre tests.
// ¿Para qué? Los tests paralelos o repetidos no pueden usar el mismo email en la BD
//            porque el campo email tiene un índice UNIQUE.
// ¿Impacto? Usar time.Now().UnixNano() garantiza unicidad en tests secuenciales.
//            Para tests paralelos, puede colisionar — no usar t.Parallel() con emails.
func UniqueEmail() string {
	return fmt.Sprintf("test_%d@test.com", time.Now().UnixNano())
}

// getEnvWithDefault retorna el valor de una variable de entorno o un valor por defecto.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
