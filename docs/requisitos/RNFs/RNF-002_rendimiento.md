# RNF-002 — Rendimiento

## Identificación

| Campo         | Valor       |
| ------------- | ----------- |
| **ID**        | RNF-002     |
| **Nombre**    | Rendimiento |
| **Categoría** | Rendimiento |
| **Prioridad** | Alta        |

---

## Descripción

El sistema debe responder de manera ágil bajo carga normal, garantizando una experiencia fluida para el usuario tanto en el backend (API) como en el frontend (carga y navegación).

---

## Métricas Requeridas

### Backend (API Go)

| Endpoint                     | Tiempo máximo de respuesta       |
| ---------------------------- | -------------------------------- |
| `POST /auth/login`           | ≤ 500ms (incluye bcrypt)         |
| `POST /auth/register`        | ≤ 500ms (incluye bcrypt)         |
| `GET /users/me`              | ≤ 100ms                          |
| `POST /auth/refresh`         | ≤ 100ms                          |
| `POST /auth/forgot-password` | ≤ 200ms (excluye envío de email) |

> **Nota:** Bcrypt con `DefaultCost` introduce ~100-200ms por diseño — es un costo intencional de seguridad.

### Frontend

| Métrica                   | Objetivo                |
| ------------------------- | ----------------------- |
| Time to Interactive (TTI) | ≤ 3 segundos en 3G      |
| First Contentful Paint    | ≤ 1.5 segundos          |
| Lighthouse Performance    | ≥ 85 en modo producción |
| Bundle inicial (gzip)     | ≤ 200KB                 |

---

## Configuraciones de Rendimiento

### Pool de Conexiones PostgreSQL (GORM)

```go
// database.go — configuración del pool
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)        // máximo conexiones abiertas
sqlDB.SetMaxIdleConns(5)         // conexiones en espera
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

### Gin — Modo Producción

```bash
# Deshabilitar los logs en producción
GIN_MODE=release go run ./cmd/api/main.go
```

### Vite — Optimizaciones de Build

```typescript
// vite.config.ts
build: {
  rollupOptions: {
    output: {
      manualChunks: {
        react: ['react', 'react-dom'],
        router: ['react-router-dom'],
      }
    }
  }
}
```

---

## Herramientas de Medición

| Herramienta                | Propósito                                      |
| -------------------------- | ---------------------------------------------- |
| `curl -w @curl-format.txt` | Medir tiempos de respuesta de la API           |
| Lighthouse                 | Auditoría de performance del frontend          |
| `go test -bench=.`         | Benchmarks de funciones críticas (bcrypt, JWT) |
| Chrome DevTools            | Análisis de red y carga de recursos            |

---

## Criterio de Cumplimiento

- API responde en los tiempos indicados medidos con `curl`
- Lighthouse Score ≥ 85 en la página de Login (producción)
- `go build ./...` completado en ≤ 30 segundos en máquina de desarrollo
