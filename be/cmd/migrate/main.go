// Archivo: main.go (cmd/migrate)
// Descripción: Herramienta CLI para ejecutar migraciones de base de datos con golang-migrate.
// ¿Para qué? Versionar los cambios del esquema de BD de forma controlada y reproducible.
// ¿Impacto? Sin migraciones versionadas, los cambios de esquema son manuales
//            y difíciles de replicar en otros entornos (staging, producción).

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	// ¿Qué? Driver pgx/v5 para golang-migrate.
	// ¿Para qué? Reutilizar el driver pgx que ya está en el proyecto
	//            en lugar de agregar lib/pq como dependencia adicional.
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	// ¿Qué? Driver de source "file" para golang-migrate.
	// ¿Para qué? Leer los archivos .sql de la carpeta migrations/ del sistema de archivos.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// ¿Qué? Cargar variables de entorno desde .env antes de cualquier otra operación.
	// ¿Para qué? DATABASE_URL y otras variables sensibles no se hardcodean en el código.
	// ¿Impacto? Si .env no existe, godotenv.Load falla silenciosamente — se depende
	//            de que las variables ya estén definidas en el entorno del sistema.
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: archivo .env no encontrado, usando variables de entorno del sistema")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Error: DATABASE_URL no está definida")
	}

	// ¿Qué? Transformar el esquema de la URL al formato requerido por el driver pgx/v5.
	// ¿Para qué? golang-migrate usa el prefijo de URL para elegir el driver correcto.
	//            El driver pgx/v5 espera "pgx5://", no "postgresql://" ni "postgres://".
	// ¿Impacto? Sin esta transformación, migrate devuelve: "unknown driver pgx5".
	migrateURL := strings.Replace(databaseURL, "postgresql://", "pgx5://", 1)
	migrateURL = strings.Replace(migrateURL, "postgres://", "pgx5://", 1)

	// ¿Qué? Instanciar migrate con la ruta a los archivos SQL y la URL de BD.
	// ¿Para qué? golang-migrate lee los archivos .up.sql y .down.sql de la carpeta migrations/.
	// ¿Impacto? El path "file://migrations" es relativo al directorio de trabajo (be/).
	//            Ejecutar este CLI desde fuera de be/ causará un error de "path not found".
	m, err := migrate.New("file://migrations", migrateURL)
	if err != nil {
		log.Fatalf("Error al inicializar migrate: %v", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Printf("Advertencia al cerrar source: %v", sourceErr)
		}
		if dbErr != nil {
			log.Printf("Advertencia al cerrar BD: %v", dbErr)
		}
	}()

	// ¿Qué? Leer el comando del primer argumento de la línea de comandos.
	// ¿Para qué? Permitir ejecutar "up", "down" o "version" según la necesidad.
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		// ¿Qué? Ejecutar todas las migraciones pendientes en orden ascendente.
		// ¿Para qué? Llevar el esquema de la BD al estado más reciente.
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error al ejecutar migraciones up: %v", err)
		}
		fmt.Println("✅ Migraciones aplicadas correctamente")

	case "down":
		// ¿Qué? Revertir la última migración aplicada.
		// ¿Para qué? Deshacer un cambio de esquema si hay errores en producción.
		// ¿Impacto? Solo revierte UNA migración a la vez — requiere llamar múltiples
		//            veces para revertir varios pasos.
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error al revertir migración: %v", err)
		}
		fmt.Println("✅ Migración revertida correctamente")

	case "version":
		// ¿Qué? Mostrar la versión actual del esquema de la base de datos.
		// ¿Para qué? Saber hasta qué migración se ha ejecutado sin entrar en la BD.
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Error al obtener versión: %v", err)
		}
		if err == migrate.ErrNilVersion {
			fmt.Println("Sin migraciones aplicadas (base de datos vacía)")
		} else {
			fmt.Printf("Versión actual: %d | Dirty: %v\n", version, dirty)
		}

	default:
		fmt.Fprintf(os.Stderr, "Comando desconocido: %q\nUso: go run ./cmd/migrate/main.go [up|down|version]\n", command)
		os.Exit(1)
	}
}
