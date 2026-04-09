// Archivo: main.go
// Descripción: Punto de entrada del servidor — inicializa config, BD, middleware y rutas.
// ¿Para qué? Orquestar el arranque de todos los componentes en el orden correcto.
// ¿Impacto? Un error aquí impide que el servidor arranque completamente.
//            El orden importa: config → BD → middleware → rutas → escuchar.

package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"nn-auth-system/internal/config"
	"nn-auth-system/internal/database"
	"nn-auth-system/internal/handlers"
	"nn-auth-system/internal/middleware"
	"nn-auth-system/internal/services"
)

func main() {
	// ¿Qué? Cargar la configuración desde variables de entorno.
	// ¿Para qué? Fail fast — si falta alguna variable obligatoria, el servidor no arranca.
	// ¿Impacto? Sin esta validación temprana, errores aparecerían en tiempo de ejecución
	//            y serían difíciles de diagnosticar.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// ¿Qué? Conectar a PostgreSQL con pool de conexiones.
	// ¿Para qué? Verificar la BD antes de recibir requests.
	// ¿Impacto? Si la BD no está disponible y el servidor arranca de todas formas,
	//            cada request fallará con un error confuso en lugar de uno claro al inicio.
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	// ¿Qué? Obtener sql.DB para el defer de cierre.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error obteniendo instancia sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// ¿Qué? Configurar el modo de Gin según el entorno.
	// ¿Para qué? En producción, Gin omite logs de debug y mensajes del banner.
	// ¿Impacto? Loggear en producción genera ruido y puede exponer información interna.
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// ¿Qué? Middleware CORS — define qué orígenes pueden hacer requests al backend.
	// ¿Para qué? El navegador bloquea requests cross-origin sin este header.
	// ¿Impacto? Sin CORS, el frontend en localhost:5173 no puede llamar al backend en :8000.
	//            En producción, NUNCA usar AllowAllOrigins: true.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// ¿Qué? Aplicar cabeceras de seguridad HTTP a todas las respuestas.
	// ¿Para qué? Proteger contra XSS, clickjacking y otros ataques del lado del cliente.
	router.Use(middleware.SecurityHeaders())

	// ¿Qué? Grupo base de rutas de la API con versionamiento.
	// ¿Para qué? Prefixar todas las rutas con /api/v1 permite versionar la API en el futuro
	//            sin romper clientes existentes cuando se lance /api/v2.
	api := router.Group("/api/v1")

	// ¿Qué? Ruta de health check.
	// ¿Para qué? Verificar que el servidor está vivo sin autenticación.
	//            Docker y load balancers usan esto para saber si el servicio está sano.
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"environment": cfg.Environment,
		})
	})

	// TODO Fase 3: registrar rutas de autenticación y usuario aquí
	// api.POST("/auth/register", authHandler.Register)
	// api.POST("/auth/login", authHandler.Login)
	// ...

	// ¿Qué? Instanciar el servicio de autenticación con la BD y la config.
	// ¿Para qué? Inyectar dependencias: el service conoce la BD y el config.
	// ¿Impacto? Si se cambia la BD o el config, solo cambia la instanciación aquí.
	authService := services.NewAuthService(db, cfg)

	// ¿Qué? Instanciar los handlers con el service inyectado.
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	// ¿Qué? Rate limiter para proteger los endpoints de autenticación.
	// ¿Para qué? Limitar 10 requests/min por IP en rutas sensibles — previene brute force.
	authRateLimiter := middleware.NewAuthRateLimiter()

	// ¿Qué? Rutas de autenticación — sin autenticación requerida.
	// ¿Para qué? Registro, login, y recuperación de contraseña son públicos.
	auth := api.Group("/auth")
	{
		// ¿Qué? Rate limiting solo en login y register — los más susceptibles a brute force.
		// ¿Para qué? verify-email y reset-password no necesitan el mismo nivel de protección
		//            ya que dependen de tokens aleatorios de 256 bits — imposibles de adivinar.
		auth.POST("/register", authRateLimiter, authHandler.Register)
		auth.POST("/login", authRateLimiter, authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/forgot-password", authRateLimiter, authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		// ¿Qué? change-password requiere estar autenticado — se registra en el grupo protegido.
	}

	// ¿Qué? Rutas protegidas — requieren JWT válido en el header Authorization.
	// ¿Para qué? Garantizar que solo usuarios autenticados accedan a estos endpoints.
	protected := api.Group("")
	protected.Use(middleware.RequireAuth(cfg))
	{
		protected.POST("/auth/change-password", authHandler.ChangePassword)
		protected.GET("/users/me", userHandler.GetMe)
	}

	addr := ":" + cfg.Port
	log.Printf("Servidor iniciado en http://localhost%s (entorno: %s)", addr, cfg.Environment)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
