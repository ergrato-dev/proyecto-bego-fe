// Archivo: database.go
// Descripción: Inicialización de la conexión a PostgreSQL mediante GORM.
// ¿Para qué? Proveer una instancia de *gorm.DB lista para usar en toda la aplicación.
// ¿Impacto? Sin una conexión correctamente configurada (pool, timeouts),
//            la aplicación bajo carga generaría errores de "too many connections".

package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nn-auth-system/internal/config"
)

// Connect establece la conexión con PostgreSQL usando GORM y configura el pool de conexiones.
// ¿Para qué? Centralizar la inicialización de la BD y sus parámetros de rendimiento.
// ¿Impacto? Un pool mal configurado puede agotar las conexiones disponibles de PostgreSQL
//
//	o mantener conexiones ociosas durante demasiado tiempo.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// ¿Qué? Configurar el nivel de log de GORM según el entorno.
	// ¿Para qué? En desarrollo, loggear todas las queries SQL para depuración.
	//            En producción, solo loggear errores para no saturar los logs.
	logLevel := logger.Error
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("abriendo conexión a la base de datos: %w", err)
	}

	// ¿Qué? Obtener la conexión SQL subyacente para configurar el pool.
	// ¿Para qué? GORM usa database/sql internamente — el pool se configura a ese nivel.
	// ¿Impacto? Sin pool, cada query abriría y cerraría una conexión TCP a PostgreSQL,
	//            lo cual es extremadamente lento bajo carga.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("obteniendo instancia sql.DB: %w", err)
	}

	// ¿Qué? Límite máximo de conexiones abiertas simultáneamente.
	// ¿Para qué? PostgreSQL por defecto acepta 100 conexiones — no saturarlo.
	sqlDB.SetMaxOpenConns(25)

	// ¿Qué? Conexiones que se mantienen abiertas en el pool cuando no se usan.
	// ¿Para qué? Reutilizar conexiones existentes evita la latencia de establecer una nueva.
	sqlDB.SetMaxIdleConns(10)

	// ¿Qué? Tiempo máximo que una conexión puede existir en el pool.
	// ¿Para qué? Evitar conexiones "zombie" que PostgreSQL ya cerró por timeout.
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// ¿Qué? Verificar que la conexión funciona antes de devolver.
	// ¿Para qué? Detectar problemas de configuración al arrancar, no al primer request.
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("verificando conexión a la base de datos (¿está PostgreSQL corriendo?): %w", err)
	}

	return db, nil
}
