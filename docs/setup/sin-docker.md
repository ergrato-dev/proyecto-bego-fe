# Guía de Instalación — Sin Docker

<!--
  ¿Qué? Guía paso a paso para levantar el sistema sin Docker,
        instalando y configurando cada componente manualmente.
  ¿Para qué? Que desarrolladores sin Docker o que quieran entender
             cada componente por separado puedan configurar el entorno.
  ¿Impacto? Más pasos, pero mayor comprensión de cómo funciona
             cada pieza del sistema independientemente.
-->

---

## Prerrequisitos

Instala y verifica cada herramienta antes de continuar:

| Herramienta | Versión mínima | Cómo verificar   | Enlace de descarga                   |
| ----------- | -------------- | ---------------- | ------------------------------------ |
| Go          | 1.25+          | `go version`     | https://go.dev/dl/                   |
| Node.js     | 20 LTS+        | `node --version` | https://nodejs.org/                  |
| pnpm        | 10+            | `pnpm --version` | `npm install -g pnpm`                |
| PostgreSQL  | 17+            | `psql --version` | https://www.postgresql.org/download/ |
| Git         | 2.0+           | `git --version`  | https://git-scm.com/                 |

> **Importante:** En Go no hay "entorno virtual" equivalente a Python's `venv`. Go modules (`go.mod`) maneja automáticamente las dependencias por proyecto. No es necesario activar nada antes de trabajar.

---

## Paso 1 — Clonar el Repositorio

```bash
git clone https://github.com/tu-usuario/nn-auth-system.git
cd nn-auth-system
```

---

## Paso 2 — Configurar PostgreSQL

### 2.1 Crear el usuario y la base de datos

Conéctate a PostgreSQL con el superusuario:

```bash
# En Linux/macOS
sudo -u postgres psql

# En Windows (usando psql)
psql -U postgres
```

Ejecuta los siguientes comandos en la consola de psql:

```sql
-- Crear el usuario de la aplicación
CREATE USER nn_user WITH PASSWORD 'nn_password';

-- Crear la base de datos
CREATE DATABASE nn_auth_db OWNER nn_user;

-- Dar permisos completos al usuario sobre la BD
GRANT ALL PRIVILEGES ON DATABASE nn_auth_db TO nn_user;

-- Salir de psql
\q
```

### 2.2 Verificar la conexión

```bash
psql -h localhost -U nn_user -d nn_auth_db -c "SELECT version();"
```

Si muestra la versión de PostgreSQL, la conexión es correcta.

### 2.3 Habilitar la extensión uuid-ossp (para UUID)

```bash
psql -h localhost -U nn_user -d nn_auth_db -c "CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";"
```

> La extensión `pgcrypto` provee la función `gen_random_uuid()` usada en las migraciones.

---

## Paso 3 — Configurar el Backend (Go)

```bash
cd be

# Copiar el archivo de variables de entorno
cp .env.example .env
```

Editar `be/.env`:

```bash
# be/.env — Configuración para entorno sin Docker
DATABASE_URL=postgresql://nn_user:nn_password@localhost:5432/nn_auth_db
JWT_SECRET=dev-secret-key-change-in-production-min-32-chars-here
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=15
JWT_REFRESH_TOKEN_EXPIRE_DAYS=7
PORT=8000
ENVIRONMENT=development
FRONTEND_URL=http://localhost:5173
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
FROM_EMAIL=noreply@nn-company.com
```

### Descargar dependencias Go

```bash
# Desde be/
go mod download

# Verificar que no haya problemas con las dependencias
go mod verify
```

### Ejecutar las migraciones

```bash
# Desde be/
go run ./cmd/migrate/main.go up
```

Deberías ver:

```
✅ Migraciones aplicadas correctamente
```

> Si todas las migraciones ya estaban aplicadas, el comando no imprime nada (sin error).

---

## Paso 4 — Configurar el Frontend

```bash
cd ../fe   # (o cd fe desde la raíz)

# Copiar variables de entorno
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

## Paso 5 — Configurar el Email (Opciones)

Sin Docker, tienes tres opciones para emails de desarrollo:

### Opción A: Deshabilitar emails (más simple)

Dejar los campos SMTP vacíos en `be/.env`. Los emails no se enviarán pero el sistema funcionará. Para verificar emails manualmente, puedes actualizar la BD directamente:

```sql
UPDATE users SET is_email_verified = TRUE WHERE email = 'test@ejemplo.com';
```

### Opción B: Instalar Mailpit localmente

```bash
# Linux/macOS (requiere Go instalado)
go install github.com/axllent/mailpit@latest

# O descargar binario desde: https://github.com/axllent/mailpit/releases
# Levantar Mailpit
mailpit
```

Mailpit estará disponible en:

- UI: http://localhost:8025
- SMTP: localhost:1025

### Opción C: Usar un servicio de email real (Resend, SMTP de Gmail)

```bash
# be/.env — Usando Resend
SMTP_HOST=smtp.resend.com
SMTP_PORT=587
SMTP_USERNAME=resend
SMTP_PASSWORD=re_xxxxxxxxxxxxxxxxxxxxxxxxx
FROM_EMAIL=noreply@tu-dominio.com
```

---

## Paso 6 — Levantar Backend y Frontend

### Terminal 1 — Backend

```bash
cd be
go run ./cmd/api/main.go
```

Salida esperada:

```
2026/03/15 10:30:00 Servidor iniciado en :8000
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
[GIN-debug] POST   /api/v1/auth/register       --> handlers.Register
[GIN-debug] POST   /api/v1/auth/login          --> handlers.Login
[GIN-debug] POST   /api/v1/auth/refresh        --> handlers.Refresh
[GIN-debug] POST   /api/v1/auth/change-password --> handlers.ChangePassword
[GIN-debug] POST   /api/v1/auth/forgot-password --> handlers.ForgotPassword
[GIN-debug] POST   /api/v1/auth/reset-password  --> handlers.ResetPassword
[GIN-debug] POST   /api/v1/auth/verify-email    --> handlers.VerifyEmail
[GIN-debug] GET    /api/v1/users/me             --> handlers.GetMe
[GIN-debug] Listening and serving HTTP on :8000
```

### Terminal 2 — Frontend

```bash
cd fe
pnpm dev
```

Salida esperada:

```
  VITE v6.x.x  ready in 300 ms

  ➜  Local:   http://localhost:5173/
```

---

## Verificación

| Servicio     | URL                                      | Estado esperado      |
| ------------ | ---------------------------------------- | -------------------- |
| Frontend     | http://localhost:5173                    | Landing page visible |
| Backend      | http://localhost:8000/health             | `{"status":"ok"}`    |
| API Register | POST localhost:8000/api/v1/auth/register | 201 Created          |

### Test rápido con curl

```bash
# Registrar un usuario de prueba
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@ejemplo.com","full_name":"Test User","password":"Password123"}'
```

Respuesta esperada:

```json
{
  "message": "Usuario registrado. Revisa tu email para verificar tu cuenta.",
  "data": { "id": "...", "email": "test@ejemplo.com", ... }
}
```

---

## Compilar el Backend (Opcional)

Para producción o para una ejecución más rápida, compila el binario:

```bash
cd be

# Compilar para el sistema operativo actual
go build -o ./bin/api ./cmd/api/main.go

# Ejecutar el binario compilado
./bin/api

# Cross-compile para Linux desde macOS/Windows
GOOS=linux GOARCH=amd64 go build -o ./bin/api-linux ./cmd/api/main.go
```

---

## Comandos de Testing

```bash
# Backend — todos los tests
cd be && go test ./... -v

# Backend — con cobertura
cd be && go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out   # Abre el reporte en el navegador

# Backend — lint
cd be && golangci-lint run ./...

# Backend — formatear código
cd be && gofmt -w .

# Frontend — tests
cd fe && pnpm test

# Frontend — lint
cd fe && pnpm lint

# Frontend — formatear código
cd fe && pnpm format
```

---

## Comandos de Migraciones

```bash
cd be

# Aplicar todas las migraciones pendientes
go run ./cmd/migrate/main.go up

# Revertir la última migración aplicada (un paso atrás)
go run ./cmd/migrate/main.go down

# Ver en qué versión está la BD
go run ./cmd/migrate/main.go version
```

> **Nota:** El comando `down` revierte **una sola migración** (el último paso aplicado).
> Para revertir varias, ejecutarlo múltiples veces.
> Los argumentos adicionales (ej. `down 1`) son ignorados por el CLI.

---

## Solución de Problemas

### Error: "go: command not found"

Instalar Go desde https://go.dev/dl/ y agregar al PATH:

```bash
# Linux/macOS — añadir al ~/.bashrc o ~/.zshrc
export PATH=$PATH:/usr/local/go/bin

# Recargar el shell
source ~/.bashrc
```

### Error: "pq: role 'nn_user' does not exist"

El usuario de PostgreSQL no fue creado. Ejecutar el Paso 2 nuevamente.

### Error: "pq: password authentication failed"

Verificar que:

1. La contraseña en `DATABASE_URL` coincide con la creada en PostgreSQL
2. PostgreSQL acepta autenticación por contraseña (`md5` o `scram-sha-256` en `pg_hba.conf`)

```bash
# Editar pg_hba.conf para permitir auth con contraseña
sudo nano /etc/postgresql/17/main/pg_hba.conf

# Cambiar 'peer' por 'md5' en líneas de localhost:
# host  all  all  127.0.0.1/32  md5

# Reiniciar PostgreSQL
sudo systemctl restart postgresql
```

### Error: "dial tcp [::1]:5432: connect: connection refused"

PostgreSQL no está corriendo:

```bash
# Linux — iniciar PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql   # Para que inicie automáticamente

# macOS — iniciar con Homebrew
brew services start postgresql@17
```

### Error: "port 8000 already in use"

```bash
# Encontrar el proceso
lsof -i :8000

# Matarlo
kill -9 <PID>
```

### Migraciones fallan: "no such file or directory"

Verificar que se ejecuta desde la carpeta `be/`, no desde la raíz del proyecto:

```bash
cd be
go run ./cmd/migrate/main.go up   # ✅ Correcto
```

### Error de GORM: "column does not exist"

Las migraciones no se aplicaron. Ejecutar:

```bash
cd be && go run ./cmd/migrate/main.go up
```
