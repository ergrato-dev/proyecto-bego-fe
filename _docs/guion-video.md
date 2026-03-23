# Guión de Video — NN Auth System

## Información del Bootcamp

| Campo                           | Valor                                                                    |
| ------------------------------- | ------------------------------------------------------------------------ |
| **Proyecto**                    | NN Auth System                                                           |
| **Duración estimada del video** | 45–60 minutos                                                            |
| **Audiencia**                   | Aprendices SENA — nivel intermedio de programación                       |
| **Objetivo**                    | Construir un sistema de autenticación completo con Go y React desde cero |

---

## Estructura del Video

```
TOTAL: ~60 minutos
├── 0:00 - 3:00   Introducción y contexto (3 min)
├── 3:00 - 10:00  Arquitectura y tech stack (7 min)
├── 10:00 - 18:00 Setup del entorno (8 min)
├── 18:00 - 28:00 Backend — modelos, migraciones, seguridad (10 min)
├── 28:00 - 40:00 Backend — handlers, servicios, middleware (12 min)
├── 40:00 - 52:00 Frontend — componentes y flujos (12 min)
└── 52:00 - 60:00 Demostración final y cierre (8 min)
```

---

## Segmento 1 — Introducción y Contexto (0:00 – 3:00)

### Objetivo del segmento

Ubicar al estudiante en el contexto del proyecto y generar motivación.

### Guión

---

> "Bienvenidos a este bootcamp de autenticación con Go y React."

> "Hoy vamos a construir desde cero un sistema de autenticación completo — el mismo que cualquier aplicación real necesita: registro, login, recuperación de contraseña y protección de rutas."

> "Antes de que existieran los frameworks modernos, los desarrolladores tenían que implementar todo esto manualmente. Hoy vamos a entender exactamente cómo funciona por dentro."

> "Este proyecto tiene dos partes: un backend en Go con el framework Gin, y un frontend en React con TypeScript. Ambos se comunican a través de una API REST usando JWT para la autenticación."

> "Al final de este bootcamp vas a entender no solo CÓMO funciona — sino POR QUÉ está diseñado así, qué problemas de seguridad resuelve y qué pasaría si no lo hiciéramos correctamente."

---

**Mostrar en pantalla:** Diagrama de alto nivel de la arquitectura (frontend ↔ API ↔ PostgreSQL)

---

## Segmento 2 — Arquitectura y Tech Stack (3:00 – 10:00)

### Objetivo del segmento

Explicar las decisiones tecnológicas y la estructura del proyecto.

### Guión

---

> "Veamos el stack tecnológico y por qué elegimos cada herramienta."

**Go + Gin:**

> "Para el backend usamos Go, un lenguaje compilado y estáticamente tipado que ofrece excelente rendimiento. Gin es el framework HTTP más popular de Go — comparable a Express en Node.js o FastAPI en Python, pero más rápido."

**GORM:**

> "Para interactuar con la base de datos usamos GORM, un ORM que nos permite trabajar con structs de Go en lugar de escribir SQL manualmente. Pero ojo — vamos a aprender también las migraciones en SQL puro porque es importante entender qué hay debajo."

**postgresql + golang-migrate:**

> "La base de datos es PostgreSQL 17. Para las migraciones usamos golang-migrate, que aplica archivos `.sql` versionados — no Python ni Go, sino SQL puro. Esto nos da transparencia total sobre los cambios en la base de datos."

**JWT:**

> "La autenticación es stateless y usa JWT — JSON Web Tokens. Emitimos dos tokens: un access token que dura 15 minutos, y un refresh token que dura 7 días. Vamos a ver exactamente por qué esta duración corta es una decisión de seguridad."

**React + Vite + TypeScript:**

> "El frontend es React 18 con TypeScript para tipado estático, y Vite como bundler ultrarrápido. Los estilos son con TailwindCSS."

---

**Mostrar en pantalla:** Árbol de directorios del proyecto completo (`be/` y `fe/`)

---

## Segmento 3 — Setup del Entorno (10:00 – 18:00)

### Objetivo del segmento

Guiar al estudiante en la configuración del entorno de desarrollo.

### Guión

---

> "Vamos a configurar el entorno. Si tienes Docker instalado, esto va a ser muy rápido."

**Con Docker:**

```bash
# Clonar el repositorio
git clone https://github.com/tu-usuario/nn-auth-system.git
cd nn-auth-system

# Levantar la base de datos y el servidor de email
docker compose up -d

# Verificar que los contenedores están corriendo
docker compose ps
```

> "Docker levanta PostgreSQL 17 y Mailpit — un servidor de email de desarrollo que captura todos los emails y los muestra en una interfaz web. Va a ser muy útil para probar el flujo de recuperación de contraseña."

> "Mailpit está disponible en http://localhost:8025 — ábrelo, lo vamos a usar más adelante."

**Configurar el Backend:**

```bash
cd be
cp .env.example .env
# Editar .env con los valores correctos

# Descargar dependencias
go mod download

# Ejecutar migraciones
go run ./cmd/migrate/main.go up

# Iniciar el servidor
go run ./cmd/api/main.go
```

> "Si ves 'Server running on :8000', perfecto. Si hay un error de conexión, verifica que Docker esté corriendo."

**Configurar el Frontend:**

```bash
cd fe
cp .env.example .env
pnpm install
pnpm dev
```

> "Usamos pnpm — no npm, no yarn. Pnpm es más rápido y eficiente con el espacio en disco."

---

## Segmento 4 — Backend: Modelos y Seguridad (18:00 – 28:00)

### Objetivo del segmento

Explicar los modelos de datos, las migraciones y las funciones de seguridad críticas.

### Guión

---

**Modelos GORM:**

> "Empecemos por el modelo de Usuario. En Go, un modelo GORM es simplemente un struct con tags especiales:"

```go
// ¿Qué?  Modelo GORM que mapea la tabla 'users' en PostgreSQL.
// ¿Para qué? Permite a GORM hacer queries y operaciones CRUD sin SQL manual.
// ¿Impacto? Sin este modelo, tendríamos que escribir SQL en cada operación.
type User struct {
    ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email           string    `gorm:"uniqueIndex;not null"`
    FullName        string    `gorm:"not null"`
    HashedPassword  string    `gorm:"not null"`
    IsActive        bool      `gorm:"default:true"`
    IsEmailVerified bool      `gorm:"default:false"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

> "Nótese que el campo se llama `HashedPassword` — no `Password`. Esto no es solo semántico: es un recordatorio constante de que NUNCA almacenamos contraseñas en texto plano."

**Hashing de Contraseñas:**

> "Cuando el usuario se registra, tomamos su contraseña y la procesamos con bcrypt:"

```go
// ¿Qué?  Convierte la contraseña en texto plano a un hash bcrypt irreversible.
// ¿Para qué? Almacenar contraseñas de forma que aunque alguien acceda a la BD,
//            no pueda recuperar las contraseñas originales.
// ¿Impacto? Si omitimos este paso, una filtración de BD expone TODAS las contraseñas.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("hashing password: %w", err)
    }
    return string(bytes), nil
}
```

> "Bcrypt es unidireccional — no se puede 'descifrar'. Para verificar una contraseña, bcrypt aplica el mismo proceso al input del usuario y compara los hashes. Esto puede tomar 100-200ms intencionalmente — ese delay hace que los ataques de fuerza bruta sean imprácticos."

---

## Segmento 5 — Backend: Handlers y Middleware (28:00 – 40:00)

### Objetivo del segmento

Mostrar cómo funcionan los handlers de autenticación y el middleware JWT.

### Guión

---

**Flujo de Login:**

> "El handler de login tiene exactamente 5 responsabilidades. Ni más, ni menos."

```go
// 1. Extraer y validar los datos del request
// 2. Buscar el usuario en la BD por email
// 3. Verificar la contraseña con bcrypt
// 4. Verificar que el email esté confirmado
// 5. Generar y retornar los tokens JWT
```

> "Noten algo importante: si el email no existe O si la contraseña es incorrecta, el mensaje de error es EXACTAMENTE el mismo: 'Email o contraseña incorrectos'. Esto previene un ataque llamado 'enumeración de usuarios' — si diéramos mensajes distintos, un atacante podría descubrir qué emails están registrados."

**Middleware JWT:**

> "El middleware de autenticación intercepta los requests a rutas protegidas:"

```go
// ¿Qué?  Middleware que verifica el JWT en el header Authorization.
// ¿Para qué? Garantizar que solo usuarios autenticados accedan a rutas protegidas.
// ¿Impacto? Sin este middleware, cualquiera podría acceder a GET /users/me
//           o POST /auth/change-password sin autenticarse.
func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractBearerToken(c)
        // ... validar token, extraer user_id, inyectar en contexto
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

---

## Segmento 6 — Frontend: Componentes y Flujos (40:00 – 52:00)

### Objetivo del segmento

Mostrar la arquitectura del frontend y los flujos de autenticación desde React.

### Guión

---

**AuthContext:**

> "En React, usamos Context API para compartir el estado de autenticación con toda la aplicación. Esto evita 'prop drilling' — pasar el usuario por 10 niveles de componentes."

```typescript
// ¿Qué?  Contexto que gestiona el estado global de autenticación.
// ¿Para qué? Proveer a cualquier componente acceso al usuario actual y las
//            acciones de auth (login, logout) sin pasar props manualmente.
// ¿Impacto? Sin este contexto, cada componente tendría que gestionar su propio
//           estado de autenticación — código duplicado e inconsistente.
interface AuthContextType {
  user: UserResponse | null;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
}
```

**Interceptor de Axios:**

> "El interceptor de Axios maneja automáticamente el refresh de tokens. Cuando la API retorna 401, el interceptor intenta renovar el token y reintentar el request original — todo transparentemente para el usuario."

**ProtectedRoute:**

> "Este componente es el guardián de las rutas privadas. Si el usuario no está autenticado, redirige al login. Si está cargando el estado de auth, muestra un spinner para evitar el 'flash' de contenido no autorizado."

---

## Segmento 7 — Demostración Final y Cierre (52:00 – 60:00)

### Objetivo del segmento

Demostrar el flujo completo del sistema funcionando.

### Guión

---

> "Ahora veamos todo el sistema funcionando. Vamos a hacer el flujo completo."

**Flujo a demostrar:**

1. **Registro** → `http://localhost:5173/register`
   - Completar formulario con email, nombre y contraseña
   - Mostrar el email de verificación en Mailpit (`http://localhost:8025`)
   - Hacer clic en el link de verificación

2. **Login** → `http://localhost:5173/login`
   - Iniciar sesión con las credenciales registradas
   - Mostrar el redirect al dashboard
   - Abrir DevTools → Application → LocalStorage (mostrar los tokens)

3. **Ver perfil** → Dashboard
   - Mostrar los datos del usuario (sin contraseña en la respuesta)

4. **Cambio de contraseña** → `/change-password`
   - Actualizar la contraseña con éxito

5. **Recuperación** → `/forgot-password`
   - Solicitar recuperación de contraseña
   - Verificar en Mailpit que llegó el email
   - Seguir el link y establecer nueva contraseña

6. **Logout**
   - Confirmar que las rutas protegidas ya no son accesibles

---

> "Eso es el sistema completo. Hemos cubierto autenticación stateless con JWT, hashing seguro de contraseñas con bcrypt, verificación de email, recuperación de contraseña, y protección de rutas tanto en el backend como en el frontend."

> "Lo más importante que llevas de este bootcamp no son las líneas de código — es entender POR QUÉ están diseñadas así: por qué bcrypt, por qué tokens de corta duración, por qué el mismo mensaje de error para email y contraseña incorrectos, por qué las migraciones son SQL versionado."

> "Ese 'por qué' es lo que diferencia a un desarrollador que copia código de uno que construye sistemas seguros."

> "El código está en el repositorio. Léelo, modifícalo, rómpelo y repáralo. Esa es la mejor forma de aprender."

---

## Notas de Producción

| Aspecto            | Detalle                                                               |
| ------------------ | --------------------------------------------------------------------- |
| Resolución         | 1920x1080 mínimo                                                      |
| Terminal           | Fondo oscuro, fuente monoespaciada grande (≥ 16px)                    |
| Navegador          | Chrome DevTools visible en segmentos de frontend                      |
| Velocidad          | No acelerar el video — las pausas mientras se escribe son pedagógicas |
| Marcadores         | Usar capítulos de YouTube para cada segmento                          |
| Código en pantalla | Mostrar el archivo completo, no solo el fragmento mencionado          |

---

## Recursos Adicionales para el Estudiante

| Recurso                  | URL / Comando                              |
| ------------------------ | ------------------------------------------ |
| Documentación de Gin     | https://gin-gonic.com/docs/                |
| Documentación de GORM    | https://gorm.io/docs/                      |
| golang-jwt               | https://github.com/golang-jwt/jwt          |
| golang.org/x/crypto      | `go doc golang.org/x/crypto/bcrypt`        |
| React Docs               | https://react.dev                          |
| TailwindCSS              | https://tailwindcss.com/docs               |
| OWASP Top 10             | `_docs/conceptos/owasp-top-10.md`          |
| Arquitectura del sistema | `_docs/referencia-tecnica/architecture.md` |
