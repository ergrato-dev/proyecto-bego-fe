# Archivo: Dockerfile (be/)
# Descripción: Build multi-stage para el backend Go + Gin — compila binarios estáticos y corre en Alpine.
# ¿Para qué? Generar una imagen de producción mínima que solo contiene los binarios compilados,
#            sin el toolchain de Go (que pesa ~700 MB).
# ¿Impacto? La imagen final pesa ~30-40 MB en lugar de ~800 MB con el toolchain completo.
#           Los binarios CGO_ENABLED=0 son estáticos — no necesitan libc ni dependencias externas.

# ─────────────────────────────────────────────────────────────
# Stage 1: builder — descarga módulos y compila los binarios
# ─────────────────────────────────────────────────────────────
# ¿Qué? Go 1.25.7 con Alpine — versión pinada exacta que coincide con go.mod.
# ¿Para qué? Garantizar que el compilador usado en Docker es idéntico al de desarrollo local.
FROM golang:1.25.7-alpine AS builder

# ¿Qué? git es requerido por go mod download para algunos módulos con dependencias VCS.
RUN apk --no-cache add git

WORKDIR /build

# ¿Qué? Copiar módulos primero para aprovechar el caché de capas Docker.
# ¿Para qué? Si go.mod y go.sum no cambian, go mod download no se repite en builds posteriores.
# ¿Impacto? Ahorra 1-3 minutos por build cuando solo cambia el código fuente.
COPY go.mod go.sum ./
RUN go mod download

# ¿Qué? Copiar el código fuente completo después de descargar dependencias.
COPY . .

# ¿Qué? Compilar el servidor API como binario estático para Linux.
# ¿Para qué? CGO_ENABLED=0 desactiva CGO — genera un binario sin dependencias de C.
#            GOOS=linux asegura que el binario es para Linux aunque se compile en macOS/Windows.
#            -a -installsuffix cgo fuerza recompilación limpia sin artefactos de CGO.
# ¿Impacto? Sin CGO_ENABLED=0, Alpine puede fallar porque no tiene glibc — solo musl.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/server ./cmd/api

# ¿Qué? Compilar la herramienta de migraciones como binario separado.
# ¿Para qué? Permite ejecutar `./migrate up` al iniciar el contenedor sin el toolchain Go.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/migrate ./cmd/migrate

# ─────────────────────────────────────────────────────────────
# Stage 2: runtime — imagen mínima sin toolchain Go
# ─────────────────────────────────────────────────────────────
# ¿Qué? Alpine 3.21 — imagen base mínima (~5 MB) para ejecutar los binarios compilados.
# ¿Para qué? No necesitamos Go en producción — solo los binarios y sus archivos de datos.
FROM alpine:3.21 AS runtime

# ¿Qué? ca-certificates permite conexiones TLS/HTTPS a servicios externos (SMTP, etc.).
#        tzdata provee zonas horarias para timestamps correctos en logs y BD.
# ¿Para qué? Sin ca-certificates, las llamadas HTTPS fallan con "certificate signed by unknown authority".
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# ¿Qué? Configurar Gin en modo release dentro del contenedor.
# ¿Para qué? Desactiva los logs de debug y advertencias de desarrollo de Gin.
#            Mejora el rendimiento al eliminar chequeos adicionales de desarrollo.
ENV GIN_MODE=release

# ¿Qué? Copiar los binarios compilados desde el stage builder.
# ¿Para qué? Solo los binarios — sin código fuente ni módulos Go.
COPY --from=builder /out/server ./server
COPY --from=builder /out/migrate ./migrate

# ¿Qué? Copiar los archivos SQL de migraciones al directorio de trabajo.
# ¿Para qué? El binario `migrate` usa la ruta relativa "file://migrations" —
#            requiere que la carpeta migrations/ exista en /app/migrations/.
# ¿Impacto? Sin esta carpeta, `./migrate up` falla con "source driver: open migrations: no such file".
COPY migrations/ ./migrations/

EXPOSE 8000

# ¿Qué? Ejecutar las migraciones pendientes y luego iniciar el servidor HTTP.
# ¿Para qué? Garantizar que el esquema de BD está actualizado ANTES de servir tráfico.
#            Si hay una nueva migración, se aplica automáticamente al desplegar.
# ¿Impacto? Si las migraciones fallan (BD no disponible, SQL inválido), el contenedor
#           sale con error — comportamiento correcto para detectar problemas temprano.
CMD ["sh", "-c", "./migrate up && ./server"]
