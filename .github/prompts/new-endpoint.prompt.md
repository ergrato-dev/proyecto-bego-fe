---
description: "Crea un endpoint Go/Gin completo: DTO, handler, servicio y test con testify. Usar cuando se necesite agregar una nueva ruta a la API."
name: "Nuevo endpoint Go/Gin"
argument-hint: "Describe el endpoint: método HTTP, ruta, qué hace, qué datos recibe y devuelve, si requiere autenticación"
agent: "agent"
---

# Nuevo endpoint Go/Gin — NN Auth System

Crea un endpoint Gin completo siguiendo las convenciones del proyecto.

## Convenciones obligatorias

- **Idioma del código**: inglés (nombres de funciones, structs, variables, rutas)
- **Idioma de comentarios**: español (godoc style en funciones exportadas)
- **Comentarios pedagógicos**: cada bloque significativo responde ¿Qué? ¿Para qué? ¿Impacto?
- **Tipos**: obligatorios en todos los parámetros y retornos — nunca `interface{}`
- **Errores**: manejar explícitamente — NUNCA ignorar con `_`
- **Cabecera de archivo**: incluir si se crea un archivo nuevo (ver copilot-instructions.md §3.4)

## Lo que debes generar

### 1. DTO (`be/internal/dto/auth.go` o el archivo DTO correspondiente)

- Request struct con tags `json` y `validate` de `go-playground/validator`
- Response struct con solo los campos necesarios (nunca exponer `HashedPassword`)
- Usar `validate:"required,email,max=255"` para emails, `validate:"required,min=8"` para passwords

```go
// ¿Qué? DTO de request para el nuevo endpoint.
// ¿Para qué? Separar la capa de transporte HTTP del modelo de BD.
// ¿Impacto? Permite validar y transformar los datos antes de que lleguen al servicio.
type NuevoRequest struct {
    Campo string `json:"campo" validate:"required,max=255"`
}

type NuevoResponse struct {
    ID    string `json:"id"`
    Campo string `json:"campo"`
}
```

### 2. Función en el servicio (`be/internal/services/auth_service.go`)

- Lógica de negocio pura — sin contexto HTTP (`c *gin.Context`)
- Recibe `db *gorm.DB` y los datos del DTO
- Retorna `(resultado, error)` — nunca `error` como único retorno si hay datos
- Referencia: [be/internal/services/auth_service.go](../../../be/internal/services/auth_service.go)

```go
// NuevoFlujo implementa la lógica de negocio para el nuevo endpoint.
// ¿Para qué? Separar la lógica del transporte HTTP para facilitar los tests.
// ¿Impacto? El handler queda limpio y el servicio es testeable sin HTTP.
func NuevoFlujo(db *gorm.DB, req dto.NuevoRequest) (*dto.NuevoResponse, error) {
    // ...
}
```

### 3. Handler (`be/internal/handlers/auth.go` o `user.go`)

- Usa `c.ShouldBindJSON(&req)` para deserializar
- Valida con `validate.Struct(req)` del `go-playground/validator`
- Llama al servicio y mapea el resultado a la respuesta JSON
- Retorna con `c.JSON(http.StatusOK, response)`
- Referencia: [be/internal/handlers/auth.go](../../../be/internal/handlers/auth.go)

### 4. Registro de la ruta (`be/cmd/api/main.go`)

- Endpoints públicos: directamente en el grupo `/api/v1/auth/`
- Endpoints protegidos: dentro del grupo con `middleware.AuthMiddleware(cfg.SecretKey)`
- Referencia: [be/cmd/api/main.go](../../../be/cmd/api/main.go)

### 5. Test (`be/internal/handlers/` o `be/internal/services/`)

- Usar `net/http/httptest` + `testify/assert`
- Casos mínimos: éxito, input inválido (422), no autenticado si aplica (401), error de negocio
- Referencia: [be/internal/handlers/auth_test.go](../../../be/internal/handlers/auth_test.go)

```go
func TestNuevoEndpoint(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        // ...
        assert.Equal(t, http.StatusOK, w.Code)
    })
    t.Run("invalid input", func(t *testing.T) {
        // ...
        assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
    })
}
```

## Prefijo base de la API

Todos los endpoints se registran bajo `/api/v1/`. Verificar el registro de rutas
en [be/cmd/api/main.go](../../../be/cmd/api/main.go).

## Seguridad (OWASP)

- Validar todos los inputs con `go-playground/validator` — nunca confiar en datos del cliente
- Usar mensajes de error genéricos en auth (no revelar si un email existe)
- Rate limiting ya configurado vía `ulule/limiter` en los endpoints de auth
- Passwords: siempre `utils.HashPassword()` antes de almacenar, nunca texto plano
- Errores de GORM no exponer directamente al cliente (loggear internamente)

## Descripción del endpoint a crear

$input
