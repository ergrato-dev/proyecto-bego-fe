# Guía de Instalación — Con Docker

<!--
  ¿Qué? Guía paso a paso para levantar todo el sistema usando Docker Compose.
  ¿Para qué? Que cualquier desarrollador pueda tener el entorno funcionando
             en minutos, sin instalar Go, Node.js ni PostgreSQL manualmente.
  ¿Impacto? Sin Docker, cada desarrollador tendría que configurar su entorno
             manualmente, con versiones distintas y problemas difíciles de reproducir.
-->

---

## Prerrequisitos

Antes de empezar, asegúrate de tener instalado:

| Herramienta    | Versión mínima | Cómo verificar           | Enlace de descarga                             |
| -------------- | -------------- | ------------------------ | ---------------------------------------------- |
| Docker Desktop | 24.0+          | `docker --version`       | https://www.docker.com/products/docker-desktop |
| Docker Compose | 2.20+          | `docker compose version` | (incluido en Docker Desktop)                   |
| Git            | 2.0+           | `git --version`          | https://git-scm.com/                           |

> **Nota:** Docker Desktop incluye Docker Compose automáticamente. Si usas Linux sin Docker Desktop, instala el plugin `docker-compose-plugin` por separado.

---

## Paso 1 — Clonar el Repositorio

```bash
git clone https://github.com/tu-usuario/nn-auth-system.git
cd nn-auth-system
```

---

## Paso 2 — Levantar la Base de Datos y Mailpit

```bash
# Desde la raíz del proyecto (donde está docker-compose.yml)
docker compose up -d
```

Esto levanta dos servicios en segundo plano:

- **PostgreSQL 17** en el puerto `5432`
- **Mailpit** (captura de emails SMTP) en el puerto `8025` (UI) y `1025` (SMTP)

Verificar que los contenedores estén corriendo:

```bash
docker compose ps
```

Salida esperada:

```
NAME               IMAGE                   STATUS
nn_auth_db         postgres:17-alpine      Up
nn_auth_mailpit    axllent/mailpit         Up
```

---

## Paso 3 — Configurar el Backend

```bash
cd be

# Copiar el archivo de ejemplo
cp .env.example .env
```

Editar `be/.env` con los valores de desarrollo (deberían funcionar con los defaults):

```bash
# be/.env
DATABASE_URL=postgresql://nn_user:nn_password@localhost:5432/nn_auth_db
SECRET_KEY=dev-secret-key-change-in-production-32-chars-min
ACCESS_TOKEN_EXPIRE_MINUTES=15
REFRESH_TOKEN_EXPIRE_DAYS=7
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
FROM_EMAIL=noreply@nn-company.com
FRONTEND_URL=http://localhost:5173
ENVIRONMENT=development
```

### Descargar dependencias de Go

```bash
# Desde be/
go mod download    # descarga todas las dependencias listadas en go.mod
```

### Ejecutar las migraciones de la base de datos

```bash
# Desde be/
go run ./cmd/migrate/main.go up
```

Deberías ver:

```
2026/03/15 10:30:00 Applying migration: 000001_create_users.up.sql
2026/03/15 10:30:00 Applying migration: 000002_create_password_reset_tokens.up.sql
2026/03/15 10:30:00 Applying migration: 000003_create_email_verification_tokens.up.sql
2026/03/15 10:30:00 Migrations applied successfully
```

---

## Paso 4 — Configurar el Frontend

```bash
cd ../fe   # (o cd fe desde la raíz)

# Copiar el archivo de ejemplo
cp .env.example .env
```

Editar `fe/.env`:

```bash
# fe/.env
VITE_API_URL=http://localhost:8000/api/v1
```

Instalar dependencias:

```bash
pnpm install
```

---

## Paso 5 — Levantar Backend y Frontend

Necesitas **dos terminales abiertas** (o usar `tmux`/`screen`):

### Terminal 1 — Backend

```bash
cd be
go run ./cmd/api/main.go
```

Deberías ver:

```
2026/03/15 10:30:00 Conectando a la base de datos...
2026/03/15 10:30:00 Base de datos conectada exitosamente
2026/03/15 10:30:00 Servidor iniciado en :8000
[GIN-debug] POST   /api/v1/auth/register   --> handlers.Register
[GIN-debug] POST   /api/v1/auth/login      --> handlers.Login
...
```

### Terminal 2 — Frontend

```bash
cd fe
pnpm dev
```

Deberías ver:

```
  VITE v6.x.x  ready in 300 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

---

## Verificación del Sistema

### URLs del sistema

| Servicio         | URL                          | Descripción                    |
| ---------------- | ---------------------------- | ------------------------------ |
| Frontend         | http://localhost:5173        | Aplicación React               |
| Backend API      | http://localhost:8000        | API Go/Gin                     |
| API Health Check | http://localhost:8000/health | Estado del servidor            |
| Mailpit UI       | http://localhost:8025        | Bandeja de email en desarrollo |
| PostgreSQL       | localhost:5432               | BD (solo con cliente SQL)      |

### Flujo de verificación manual

1. Abre http://localhost:5173 → debería ver la Landing Page
2. Haz clic en "Registrarse" → Completa el formulario
3. Abre http://localhost:8025 → Deberías ver el email de verificación
4. Haz clic en el enlace del email → Email verificado
5. Inicia sesión con las credenciales → Deberías ver el Dashboard

---

## Comandos del Día a Día

```bash
# Iniciar la base de datos y Mailpit (si están detenidos)
docker compose up -d

# Detener la base de datos y Mailpit
docker compose down

# Detener y borrar los volúmenes (borra todos los datos de la BD)
docker compose down -v

# Ver los logs de los contenedores
docker compose logs -f

# Ver logs solo de la BD
docker compose logs -f db

# Conectarse a la BD directamente
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db

# Ver las tablas
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db -c "\dt"
```

---

## Comandos de Testing

```bash
# Backend — ejecutar todos los tests
cd be && go test ./... -v

# Backend — con cobertura
cd be && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# Frontend — ejecutar todos los tests
cd fe && pnpm test

# Frontend — modo watch
cd fe && pnpm test:watch
```

---

## Solución de Problemas

### Error: "Cannot connect to Docker daemon"

```bash
# Verificar que Docker esté corriendo
systemctl status docker        # Linux
# En macOS/Windows: abrir Docker Desktop
```

### Error: "port 5432 already in use"

```bash
# Verificar qué proceso usa el puerto
sudo lsof -i :5432

# Si tienes PostgreSQL local corriendo, detenerlo
sudo systemctl stop postgresql
```

### Error: "go: command not found"

El backend se ejecuta directo con `go run`, no requiere build previo con Docker. Si Go no está instalado, sigue la guía sin Docker en `docs/setup/sin-docker.md`.

### Error: "pnpm: command not found"

```bash
# Instalar pnpm via Node.js
npm install -g pnpm

# O via corepack (Node.js 16.9+)
corepack enable
corepack prepare pnpm@latest --activate
```

### La BD existe pero no tiene las tablas

```bash
# Verificar el estado de las migraciones
cd be && go run ./cmd/migrate/main.go version

# Si es necesario, re-aplicar las migraciones
cd be && go run ./cmd/migrate/main.go up
```

### Error de CORS en el frontend

Verificar que:

1. `FRONTEND_URL=http://localhost:5173` en `be/.env`
2. El backend esté corriendo en el puerto 8000
3. `VITE_API_URL=http://localhost:8000/api/v1` en `fe/.env`
