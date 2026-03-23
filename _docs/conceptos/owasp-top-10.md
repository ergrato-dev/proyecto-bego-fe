# OWASP Top 10 — Seguridad en NN Auth System

<!--
  ¿Qué? Guía que explica las 10 vulnerabilidades de seguridad web más críticas
        de OWASP y cómo las mitigamos en este proyecto con Go y React.
  ¿Para qué? Que cada desarrollador entienda por qué ciertas decisiones de código
             existen y qué problemas concretos evitan.
  ¿Impacto? Sin este conocimiento, un desarrollador podría "mejorar" el código
             de forma que introduzca vulnerabilidades críticas sin saberlo.
-->

> **OWASP** (Open Worldwide Application Security Project) es una fundación sin fines de lucro que publica el "Top 10" — la lista de las vulnerabilidades de seguridad web más críticas y comunes. Referencia: https://owasp.org/www-project-top-ten/

---

## A01 — Control de Acceso Roto (Broken Access Control)

### ¿Qué es?

Ocurre cuando un usuario puede acceder a recursos o realizar acciones para las que no tiene permiso. Por ejemplo: ver el perfil de otro usuario, modificar datos ajenos, acceder a rutas de admin sin serlo.

### ¿Cómo lo mitigamos?

**Middleware JWT en Gin** — Todo endpoint protegido requiere un JWT válido:

```go
// middleware/auth.go
// ¿Qué? Middleware que valida el JWT y extrae el user_id al contexto.
// ¿Para qué? Que NINGÚN handler protegido se ejecute sin autenticación válida.
// ¿Impacto? Sin este middleware, cualquier persona podría llamar a /api/v1/users/me
//            sin autenticarse y obtener datos de usuarios.
func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extraer token del header Authorization: Bearer <token>
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "detail": "Token de autenticación requerido",
            })
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        userID, err := utils.ParseToken(tokenStr, secret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "detail": "Token inválido o expirado",
            })
            return
        }

        // Inyectar user_id al contexto para que los handlers lo usen
        c.Set("user_id", userID)
        c.Next()
    }
}
```

**Uso en rutas:**

```go
// cmd/api/main.go — Solo las rutas protegidas usan AuthMiddleware
protected := r.Group("/api/v1", middleware.AuthMiddleware(cfg.SecretKey))
protected.GET("/users/me", handlers.GetMe)
protected.POST("/auth/change-password", handlers.ChangePassword)
```

---

## A02 — Fallas Criptográficas (Cryptographic Failures)

### ¿Qué es?

Almacenar datos sensibles sin cifrado o con cifrado débil. El caso más común: contraseñas en texto plano o con hash reversible (MD5, SHA1).

### ¿Cómo lo mitigamos?

**bcrypt** via `golang.org/x/crypto/bcrypt`:

```go
// utils/security.go
// ¿Qué? Hashea la contraseña con bcrypt (algoritmo de hashing adaptativo).
// ¿Para qué? Almacenar contraseñas de forma que, incluso con acceso a la BD,
//             no sea viable recuperar la contraseña original.
// ¿Impacto? Si se usara SHA1/MD5, un atacante con acceso a la BD podría
//            descifrar todas las contraseñas en minutos con tablas rainbow.
func HashPassword(password string) (string, error) {
    // DefaultCost = 10 iteraciones — balance entre seguridad y performance
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("hashing password: %w", err)
    }
    return string(bytes), nil
}

// ¿Qué? Compara una contraseña en texto plano con su hash almacenado.
// ¿Para qué? Verificar credenciales sin necesidad de "descifrar" el hash.
// ¿Impacto? Esta es la única forma segura de verificar contraseñas hasheadas.
func VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**¿Por qué bcrypt y no SHA256/SHA512?**

| Algoritmo | Función | Velocidad  | Resistencia a fuerza bruta |
| --------- | ------- | ---------- | -------------------------- |
| MD5       | Hash    | Muy rápido | Nula (obsoleto)            |
| SHA256    | Hash    | Rápido     | Baja (miles de hashes/seg) |
| bcrypt    | KDF     | Lento      | Alta (adaptable, salted)   |

bcrypt es "lento" a propósito: hace que los ataques de fuerza bruta sean computacionalmente inviables.

---

## A03 — Inyección (Injection)

### ¿Qué es?

Enviar datos maliciosos que la aplicación interpreta como comandos. Los tipos más comunes:

- **SQL Injection:** `' OR '1'='1` en un campo de login
- **XSS:** `<script>alert('hack')</script>` en un campo de texto
- **Command Injection:** `; rm -rf /` en un campo enviado a la shell

### ¿Cómo lo mitigamos?

**GORM con parámetros bindeados** — Nunca SQL raw:

```go
// ✅ CORRECTO — GORM parametriza automáticamente
// La variable `email` se pasa como parámetro, no interpolada en el SQL
var user models.User
result := db.Where("email = ?", email).First(&user)

// ❌ INCORRECTO — SQL crudo con interpolación de strings (vulnerable)
// Un email como `' OR '1'='1` permitiría acceder sin contraseña
result := db.Raw("SELECT * FROM users WHERE email = '" + email + "'")
```

**Validación de inputs con go-playground/validator:**

```go
// dto/auth.go
// ¿Qué? Struct con validaciones declarativas.
// ¿Para qué? Rechazar inputs inválidos ANTES de procesarlos.
// ¿Impacto? Sin validación, datos maliciosos llegarían a la lógica de negocio.
type LoginRequest struct {
    Email    string `json:"email"    validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=1,max=128"`
}
```

**XSS en el frontend — React escapa por defecto:**

React renderiza todo texto como texto plano automáticamente. Solo es vulnerable si usas `dangerouslySetInnerHTML` (evitar siempre en este proyecto).

---

## A04 — Diseño Inseguro (Insecure Design)

### ¿Qué es?

Ausencia de controles de seguridad en el diseño de la aplicación. No es un bug de implementación, sino una omisión de los mecanismos de defensa desde el inicio.

### ¿Cómo lo mitigamos?

**Rate limiting** via `ulule/limiter`:

```go
// middleware/ratelimit.go
// ¿Qué? Middleware que limita el número de peticiones por IP en ventana de tiempo.
// ¿Para qué? Prevenir ataques de fuerza bruta (intentos masivos de login)
//             y ataques de denegación de servicio (DoS).
// ¿Impacto? Sin rate limiting, un atacante podría intentar millones de passwords
//            contra una cuenta en minutos.
func RateLimitMiddleware(rate limiter.Rate) gin.HandlerFunc {
    store := memory.NewStore()
    instance := limiter.New(store, rate)
    middleware := gin_limiter.NewMiddleware(instance)
    return middleware
}

// Configuración en main.go — más estricto en endpoints de auth
authLimiter := middleware.RateLimitMiddleware(limiter.Rate{
    Period: 15 * time.Minute,
    Limit:  10, // Max 10 intentos de login cada 15 minutos por IP
})
```

---

## A05 — Configuración de Seguridad Incorrecta (Security Misconfiguration)

### ¿Qué es?

Configuraciones por defecto inseguras, acceso innecesario habilitado, cabeceras HTTP de seguridad ausentes, CORS mal configurado.

### ¿Cómo lo mitigamos?

**Cabeceras de seguridad HTTP:**

```go
// middleware/security.go
// ¿Qué? Añade cabeceras HTTP que indican al navegador cómo manejar el contenido.
// ¿Para qué? Prevenir ataques como clickjacking, MIME sniffing y XSS.
// ¿Impacto? Sin estas cabeceras, el navegador está en modo "permisivo",
//            lo que abre vectores de ataque que se pueden cerrar trivialmente.
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Next()
    }
}
```

**CORS restrictivo:**

```go
// middleware/cors.go
// ¿Qué? Configura CORS para permitir solo el origen del frontend.
// ¿Para qué? Prevenir que cualquier página web en Internet haga requests
//             en nombre de un usuario autenticado (CSRF).
// ¿Impacto? Con CORS permisivo (AllowAllOrigins: true), cualquier sitio
//            malicioso podría hacer fetch a nuestra API con las cookies del usuario.
func CORSMiddleware(frontendURL string) gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{frontendURL}, // Solo el frontend permitido
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

---

## A06 — Componentes Vulnerables y Desactualizados

### ¿Qué es?

Usar librerías o frameworks con vulnerabilidades conocidas y no actualizadas.

### ¿Cómo lo mitigamos?

**Go modules y verificación de dependencias:**

```bash
# Listar dependencias con posibles vulnerabilidades
go list -m all

# Actualizar todas las dependencias a la última versión menor
go get -u ./...

# Auditar vulnerabilidades conocidas (herramienta oficial de Go)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Frontend — auditar vulnerabilidades de npm
cd fe && pnpm audit
```

**`go.sum`** archivo de checksums criptográficos:

- `go.sum` contiene los hashes SHA256 de cada versión de dependencia
- Si alguien modifica una dependencia, la verificación falla
- Esto protege contra ataques a la cadena de suministro (supply chain attacks)

---

## A07 — Fallas de Identificación y Autenticación

### ¿Qué es?

Contraseñas débiles permitidas, sesiones que no expiran, tokens predecibles, reutilización de tokens de reset.

### ¿Cómo lo mitigamos?

**Validación de fortaleza de contraseña:**

```go
// dto/auth.go
type RegisterRequest struct {
    // validate:"min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,
    //           containsany=abcdefghijklmnopqrstuvwxyz,
    //           containsany=0123456789"
    Password string `json:"password" validate:"required,min=8"`
}
```

**JWT con expiración corta:**

```go
// utils/security.go
// ¿Qué? Genera un JWT con expiración de 15 minutos.
// ¿Para qué? Limitar el daño si un access token es interceptado.
// ¿Impacto? Con tokens de larga duración (días), un token robado
//            permite acceso durante mucho tiempo.
func CreateAccessToken(userID string, secret string) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(15 * time.Minute).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

**Tokens de reset de uso único:**

```go
// Al verificar el token de reset, se marca como usado
// ¿Qué? Marcar el token como utilizado después de un reset exitoso.
// ¿Para qué? Prevenir que el mismo enlace de reset se use múltiples veces.
// ¿Impacto? Sin esto, un atacante que intercepte el email podría usar
//            el enlace incluso después de que la víctima ya cambió su contraseña.
db.Model(&token).Update("used", true)
```

---

## A08 — Fallas de Integridad de Software y Datos

### ¿Qué es?

Confiar en actualizaciones o datos sin verificar su integridad. En el contexto de JWT: no verificar la firma del token.

### ¿Cómo lo mitigamos?

**Verificación de firma JWT con `golang-jwt/jwt v5`:**

```go
// utils/security.go
// ¿Qué? Parsea y verifica la firma del JWT.
// ¿Para qué? Garantizar que el token fue emitido por nosotros
//             y no fue manipulado.
// ¿Impacto? Sin verificación de firma, cualquiera podría crear un JWT falso
//            con cualquier user_id y acceder como ese usuario.
func ParseToken(tokenStr string, secret string) (string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Verificar que el algoritmo sea el esperado (HS256)
        // Previene el ataque "alg=none" donde el JWT no tiene firma
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("algoritmo inesperado: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })
    // ...
}
```

---

## A09 — Fallas de Registro y Monitoreo de Seguridad

### ¿Qué es?

No registrar eventos de seguridad relevantes, lo que hace imposible detectar ataques en curso o post-mortem.

### ¿Cómo lo mitigamos?

**Audit log estructurado:**

```go
// utils/audit_log.go
// ¿Qué? Registra eventos de seguridad en formato estructurado.
// ¿Para qué? Detectar patrones de ataque (múltiples logins fallidos),
//             auditar accesos y cumplir con normativas de seguridad.
// ¿Impacto? Sin logs de auditoría, un ataque puede ocurrir sin que
//            nadie lo note hasta que el daño está hecho.
type AuditEvent struct {
    Timestamp time.Time `json:"timestamp"`
    EventType string    `json:"event_type"` // LOGIN_SUCCESS, LOGIN_FAILED, etc.
    UserID    string    `json:"user_id,omitempty"`
    Email     string    `json:"email,omitempty"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
}

func LogLoginFailed(email, ipAddress, userAgent string) {
    event := AuditEvent{
        Timestamp: time.Now(),
        EventType: "LOGIN_FAILED",
        Email:     email,
        IPAddress: ipAddress,
        UserAgent: userAgent,
    }
    // En producción: enviar a sistema de monitoreo (Datadog, CloudWatch, etc.)
    log.Printf("[AUDIT] %+v\n", event)
}
```

**Eventos que se deben registrar:**

| Evento               | Nivel   | Razón                                     |
| -------------------- | ------- | ----------------------------------------- |
| `LOGIN_SUCCESS`      | INFO    | Trazabilidad de accesos                   |
| `LOGIN_FAILED`       | WARNING | Detectar intentos de fuerza bruta         |
| `PASSWORD_CHANGED`   | INFO    | Trazabilidad de cambios de credenciales   |
| `PASSWORD_RESET_REQ` | INFO    | Trazabilidad de resets                    |
| `EMAIL_VERIFIED`     | INFO    | Trazabilidad del ciclo de vida de cuentas |
| `TOKEN_EXPIRED`      | WARNING | Monitoreo de sesiones expiradas           |

---

## A10 — Falsificación de Peticiones del Lado del Servidor (SSRF)

### ¿Qué es?

El servidor hace peticiones a URLs controladas por el atacante, potencialmente accediendo a recursos internos de la red.

### Aplicabilidad en este proyecto

Este proyecto no tiene funcionalidades que impliquen que el servidor haga fetch de URLs externas basadas en input del usuario (por ejemplo, cargar una imagen desde una URL). Si en el futuro se agrega esta funcionalidad (e.g., foto de perfil desde URL), se deben validar las URLs destino con una lista blanca de dominios permitidos.

---

## Resumen

| Vulnerabilidad                  | Herramienta Go                      | Implementación                      |
| ------------------------------- | ----------------------------------- | ----------------------------------- |
| A01 — Broken Access Control     | Middleware JWT (golang-jwt)         | `middleware/auth.go`                |
| A02 — Cryptographic Failures    | bcrypt (golang.org/x/crypto)        | `utils/security.go`                 |
| A03 — Injection                 | GORM ORM + validator v10            | DTOs + queries GORM                 |
| A04 — Insecure Design           | ulule/limiter                       | `middleware/ratelimit.go`           |
| A05 — Security Misconfiguration | CORS config + security headers      | `middleware/cors.go`, `security.go` |
| A06 — Vulnerable Components     | `govulncheck` + `go.sum`            | CI/CD + `pnpm audit`                |
| A07 — Auth Failures             | golang-jwt v5 + validator           | `utils/security.go` + DTOs          |
| A08 — Software Integrity        | JWT signature verification          | `utils/security.go` (ParseToken)    |
| A09 — Logging Failures          | Structured audit log                | `utils/audit_log.go`                |
| A10 — SSRF                      | N/A (no hay fetch de URLs externas) | N/A en la versión actual            |

---

> La seguridad no es un feature adicional que se agrega al final — es una propiedad fundamental que se construye desde el diseño. Cada línea de código relacionada con autenticación tiene implicaciones de seguridad.
