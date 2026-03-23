# RNF-005 — Mantenibilidad

## Identificación

| Campo         | Valor              |
| ------------- | ------------------ |
| **ID**        | RNF-005            |
| **Nombre**    | Mantenibilidad     |
| **Categoría** | Calidad del Código |
| **Prioridad** | Alta               |

---

## Descripción

El sistema debe ser fácil de entender, modificar y extender. El código debe seguir convenciones estrictas, estar bien documentado con comentarios pedagógicos, y contar con suficiente cobertura de tests para detectar regresiones.

---

## Estándares de Código

### Go (Backend)

| Herramienta     | Comando                   | Propósito                                |
| --------------- | ------------------------- | ---------------------------------------- |
| `gofmt`         | `gofmt -w .`              | Formato estándar obligatorio             |
| `goimports`     | `goimports -w .`          | Organiza imports (stdlib→extern→interno) |
| `golangci-lint` | `golangci-lint run ./...` | Linter + analizador estático             |
| `go vet`        | `go vet ./...`            | Verifica errores comunes de Go           |

**Reglas obligatorias:**

- Toda función exportada tiene comentario godoc en español
- Errores siempre manejados explícitamente (nunca `_`)
- Líneas máximo 100 caracteres
- Imports organizados en 3 grupos: stdlib / externo / interno

### TypeScript/React (Frontend)

| Herramienta | Comando             | Propósito                          |
| ----------- | ------------------- | ---------------------------------- |
| ESLint      | `pnpm lint`         | Detectar errores y malos patrones  |
| Prettier    | `pnpm format`       | Formateo consistente del código    |
| TypeScript  | `pnpm tsc --noEmit` | Verificación de tipos sin compilar |

**Reglas obligatorias:**

- `"strict": true` en `tsconfig.json`
- Sin `any` explícito sin justificación
- Componentes funcionales únicamente (sin clases)

---

## Cobertura de Tests

| Capa                         | Cobertura mínima | Herramienta                     |
| ---------------------------- | ---------------- | ------------------------------- |
| Backend — servicios          | 80%              | `go test -cover ./internal/...` |
| Backend — handlers           | 70%              | `go test -cover ./...`          |
| Frontend — componentes clave | 70%              | `pnpm test:coverage`            |

```bash
# Ver reporte de cobertura Go
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out    # abre reporte en navegador

# Ver reporte de cobertura frontend
cd fe && pnpm test:coverage
```

---

## Documentación del Código

Todo bloque significativo debe responder:

```go
// ¿Qué? Middleware que extrae y valida el JWT del header Authorization.
// ¿Para qué? Proteger las rutas que requieren autenticación, inyectando
//            el user_id en el contexto de Gin para los handlers.
// ¿Impacto? Sin este middleware, cualquier usuario podría acceder a rutas
//           protegidas sin credenciales válidas.
func AuthMiddleware(secretKey string) gin.HandlerFunc {
```

---

## Estructura del Proyecto

El proyecto sigue **Separation of Concerns** estricta:

| Capa       | Paquete Go       | Responsabilidad                       |
| ---------- | ---------------- | ------------------------------------- |
| Entrada    | `handlers/`      | HTTP: bind, validate, call service    |
| Negocio    | `services/`      | Lógica de negocio pura                |
| Datos      | `models/` + GORM | Entidades y acceso a la BD            |
| Transporte | `dto/`           | Structs de request/response separados |
| Utilidades | `utils/`         | Funciones reutilizables (bcrypt, JWT) |

---

## Conventional Commits

Todos los commits deben seguir el formato:

```
type(scope): descripción en inglés

What: qué se hizo
For: por qué se hizo
Impact: efecto en el sistema
```

---

## Criterio de Cumplimiento

- `golangci-lint run ./...` sin errores en el backend
- `pnpm lint` sin errores en el frontend
- Cobertura de tests ≥ 80% en paquetes de lógica de negocio
- Toda función exportada de Go tiene comentario godoc
