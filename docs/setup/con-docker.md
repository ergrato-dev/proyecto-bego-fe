# Guía de Instalación — Con Docker

<!--
  ¿Qué? Guía paso a paso para levantar todo el sistema usando Docker Compose.
  ¿Para qué? Que cualquier desarrollador pueda tener el entorno funcionando
             en minutos, con un solo comando: docker compose up --build
  ¿Impacto? Sin Docker, cada desarrollador tendría que instalar Go, Node.js,
             pnpm y PostgreSQL manualmente, con versiones distintas y problemas
             difíciles de reproducir. Con Docker, basta tener Docker instalado.
-->

> **Enfoque de esta guía:** Docker gestiona **todos** los servicios.
> No necesitas instalar Go, Node.js, pnpm ni PostgreSQL en tu máquina.

---

## Prerrequisitos

Solo necesitas dos herramientas instaladas en tu máquina:

| Herramienta    | Versión mínima | Cómo verificar           | Enlace de descarga                             |
| -------------- | -------------- | ------------------------ | ---------------------------------------------- |
| Docker Desktop | 24.0+          | `docker --version`       | https://www.docker.com/products/docker-desktop |
| Docker Compose | 2.20+          | `docker compose version` | (incluido en Docker Desktop)                   |
| Git            | 2.0+           | `git --version`          | https://git-scm.com/                           |

> **Nota:** Docker Desktop incluye Docker Compose automáticamente. Si usas Linux sin Docker Desktop, instala el plugin `docker-compose-plugin` por separado.

### ¿Qué levanta Docker Compose?

El archivo `docker-compose.yml` en la raíz del proyecto orquesta **4 servicios**:

| Servicio  | Imagen/Build              | Descripción                              |
| --------- | ------------------------- | ---------------------------------------- |
| `db`      | `postgres:17.4-alpine`    | Base de datos PostgreSQL                 |
| `be`      | `./be/Dockerfile`         | API REST Go + Gin (compila desde fuente) |
| `fe`      | `./fe/Dockerfile`         | SPA React compilada con Nginx            |
| `mailpit` | `axllent/mailpit:v1.22.0` | Servidor SMTP de desarrollo              |

El backend aplica automáticamente las migraciones SQL al iniciar.

---

## Paso 1 — Clonar el Repositorio

```bash
git clone https://github.com/tu-usuario/nn-auth-system.git
cd nn-auth-system
```

---

## Paso 2 — Levantar Todos los Servicios

```bash
# Desde la raíz del proyecto (donde está docker-compose.yml)
docker compose up --build
```

> **`--build`** compila las imágenes de `be` y `fe` desde los Dockerfiles.
> Solo es necesario la primera vez o cuando cambies el código.

Docker Compose:

1. Construye la imagen del backend (`be/Dockerfile`) — compila los binarios Go
2. Construye la imagen del frontend (`fe/Dockerfile`) — ejecuta `pnpm build` con Vite
3. Levanta `db` (PostgreSQL) con healthcheck
4. Levanta `mailpit` (SMTP de desarrollo)
5. Inicia `be` después de que la BD esté lista — aplica migraciones y arranca Gin
6. Inicia `fe` después de que `be` esté iniciado — Nginx sirve la SPA

El proceso completo tarda ~2-4 minutos la primera vez (descarga imágenes + compilación).
En builds posteriores es mucho más rápido gracias al caché de Docker.

### Para ejecutar en segundo plano (sin bloquear la terminal)

```bash
docker compose up --build -d

# Ver los logs en tiempo real
docker compose logs -f
```

---

## Verificación del Sistema

### URLs del sistema

| Servicio         | URL                          | Descripción                                |
| ---------------- | ---------------------------- | ------------------------------------------ |
| Frontend         | http://localhost:3000        | Aplicación React (Nginx)                   |
| Backend API      | http://localhost:8000        | API Go/Gin                                 |
| API Health Check | http://localhost:8000/health | Estado del servidor                        |
| Mailpit UI       | http://localhost:8025        | Bandeja de emails capturados en desarrollo |
| PostgreSQL       | localhost:5432               | BD (solo con cliente SQL)                  |

### Estado de los contenedores

```bash
docker compose ps
```

Salida esperada cuando todo está corriendo:

```
NAME               IMAGE               STATUS          PORTS
nn_auth_db         postgres:17.4-...   Up (healthy)    0.0.0.0:5432->5432/tcp
nn_auth_mailpit    axllent/mailpit:... Up              0.0.0.0:1025->1025/tcp, 0.0.0.0:8025->8025/tcp
nn_auth_be         nn-auth-be          Up              0.0.0.0:8000->8000/tcp
nn_auth_fe         nn-auth-fe          Up              0.0.0.0:3000->80/tcp
```

### Flujo de verificación manual

1. Abre http://localhost:3000 → debería ver la Landing Page
2. Haz clic en "Registrarse" → Completa el formulario
3. Abre http://localhost:8025 → Deberías ver el email de verificación capturado por Mailpit
4. Haz clic en el enlace del email → Email verificado
5. Inicia sesión con las credenciales → Deberías ver el Dashboard

---

## Variables de Entorno en Docker

Las variables de entorno del backend están configuradas directamente en `docker-compose.yml`.
No es necesario crear un archivo `be/.env` para el entorno Docker.

> **⚠️ Importante para producción:** El valor de `JWT_SECRET` en `docker-compose.yml`
> es solo para desarrollo. En producción, usa variables de entorno del sistema o un
> secret manager — nunca commitees el secret real al repositorio.

### ¿Por qué VITE_API_URL apunta a `localhost:8000`?

Vite incrusta `VITE_API_URL` en el bundle JavaScript en tiempo de **compilación** (build-time).
Cuando el usuario abre el navegador, el JavaScript hace llamadas al backend desde su máquina,
no desde dentro del contenedor. Por eso `localhost:8000` es la URL correcta en desarrollo.

```
Navegador del usuario → http://localhost:8000 → contenedor be (puerto 8000)
```

---

## Comandos del Día a Día

```bash
# Levantar todos los servicios (primera vez o con cambios de código)
docker compose up --build

# Levantar sin reconstruir (cuando el código no cambió)
docker compose up

# Detener todos los servicios (mantiene los datos en el volumen)
docker compose down

# Detener y borrar los volúmenes (BORRA TODOS LOS DATOS de la BD)
docker compose down -v

# Ver los logs de todos los servicios en tiempo real
docker compose logs -f

# Ver logs solo del backend
docker compose logs -f be

# Ver logs solo de la base de datos
docker compose logs -f db

# Conectarse a la BD directamente con psql
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db

# Ver las tablas creadas por las migraciones
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db -c "\dt"

# Reiniciar un servicio específico (útil al cambiar código)
docker compose restart be
docker compose up --build be   # reconstruye solo el backend

# Ver el estado de los contenedores
docker compose ps
```

---

## Comandos de Testing

Los tests requieren Go y Node.js instalados localmente.
Consulta la guía sin Docker (`docs/setup/sin-docker.md`) para la configuración del entorno de desarrollo nativo.

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

### Error: "port 3000 already in use"

```bash
# Verificar qué proceso usa el puerto 3000
sudo lsof -i :3000

# Alternativa: cambiar el puerto del frontend en docker-compose.yml
#   ports:
#     - "3001:80"   # cambiar 3000 por 3001
```

### El backend no arranca — error de conexión a la BD

El backend depende del healthcheck de PostgreSQL. Si `db` no pasa el healthcheck, `be` no inicia:

```bash
# Ver el estado del healthcheck
docker compose ps

# Ver logs de la BD para diagnosticar
docker compose logs db

# Si la BD está en estado "unhealthy", bajar todo y volver a levantar
docker compose down -v
docker compose up --build
```

### Las migraciones fallan al iniciar

El binario `migrate` se ejecuta automáticamente al iniciar `be`. Si fallan:

```bash
# Ver los logs del backend para leer el error de migración
docker compose logs be

# Si las migraciones están en estado "dirty" (migración parcialmente aplicada):
docker exec -it nn_auth_db psql -U nn_user -d nn_auth_db \
  -c "DELETE FROM schema_migrations WHERE dirty = true;"

# Luego reiniciar el backend
docker compose restart be
```

### El frontend muestra errores de conexión a la API

Verificar que el backend esté corriendo:

```bash
# El endpoint /health debe responder
curl http://localhost:8000/health

# Si no responde, ver los logs del backend
docker compose logs be
```

### Cambié código del backend — ¿cómo aplico los cambios?

```bash
# Reconstruir y reiniciar solo el backend
docker compose up --build be
```

### Cambié código del frontend — ¿cómo aplico los cambios?

```bash
# Reconstruir y reiniciar solo el frontend
docker compose up --build fe
```

> **Nota:** Vite incrusta `VITE_API_URL` en el bundle durante el build.
> Si cambias la URL del backend, debes reconstruir la imagen del frontend.

### Quiero empezar desde cero (base de datos vacía)

```bash
# Detener todos los servicios y eliminar los volúmenes (BORRA TODO)
docker compose down -v

# Volver a levantar — el backend creará las tablas con las migraciones
docker compose up --build
```
