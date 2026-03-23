# Patrones Arquitectónicos — NN Auth System

<!--
  ¿Qué? Docuemntación de los patrones de diseño y arquitectura usados en el proyecto.
  ¿Para qué? Que los desarrolladores entiendan el "por qué" de la estructura del código,
             no solo el "qué" hace cada archivo.
  ¿Impacto? Sin entender los patrones, el código se rompe fácilmente:
             se añade lógica de negocio en handlers, se duplica código,
             o se crean dependencias circulares.
-->

---

## 1. Arquitectura en Capas (Layered Architecture)

### ¿Qué es?

Organizar el código en capas con responsabilidades bien definidas. Cada capa solo conoce a la capa inmediatamente inferior.

### ¿Cómo se aplica en el backend Go?

```
┌─────────────────────────────────────────────────────┐
│                  Capa HTTP / Transporte              │
│  handlers/auth.go, handlers/user.go                  │
│  Responsabilidad: Parsear requests, responder JSON   │
│  Solo conoce: DTOs, Services                         │
└──────────────────────────┬──────────────────────────┘
                           │ Llama a
                           ▼
┌─────────────────────────────────────────────────────┐
│               Capa de Lógica de Negocio              │
│  services/auth_service.go                            │
│  Responsabilidad: Orquestar el flujo de negocio      │
│  Solo conoce: Models, DTOs, Utils                    │
└──────────────────────────┬──────────────────────────┘
                           │ Usa
                           ▼
┌─────────────────────────────────────────────────────┐
│               Capa de Acceso a Datos                │
│  models/user.go, GORM                                │
│  Responsabilidad: Queries a PostgreSQL               │
│  Solo conoce: PostgreSQL / GORM                      │
└─────────────────────────────────────────────────────┘
```

### Ejemplo: ¿Por qué NO poner la lógica en el handler?

```go
// ❌ ANTIPATRÓN — Handler con lógica de negocio (difícil de testear y mantener)
func LoginHandler(c *gin.Context) {
    var req dto.LoginRequest
    c.BindJSON(&req)

    var user models.User
    db.Where("email = ?", req.Email).First(&user)

    if !utils.VerifyPassword(req.Password, user.HashedPassword) {
        c.JSON(401, gin.H{"detail": "Credenciales inválidas"})
        return
    }

    token, _ := utils.CreateAccessToken(user.ID.String(), cfg.SecretKey)
    c.JSON(200, gin.H{"access_token": token})
}

// ✅ CORRECTO — Handler delega en Service (testeable, separado)
func LoginHandler(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(422, gin.H{"detail": err.Error()})
        return
    }

    // El handler solo orquesta HTTP; la lógica está en el service
    response, err := authService.Login(c.Request.Context(), req)
    if err != nil {
        c.JSON(401, gin.H{"detail": err.Error()})
        return
    }

    c.JSON(200, response)
}
```

---

## 2. DTO (Data Transfer Object)

### ¿Qué es?

Un objeto que define la forma exacta de los datos que entran y salen de la API, independiente de cómo se almacenan en la base de datos.

### ¿Por qué separar DTOs de Modelos?

```
Cliente HTTP ←──► DTO (dto/auth.go) ←──► Service ←──► Modelo GORM (models/user.go) ←──► BD
```

Si los DTOs y modelos fueran el mismo struct:

- La API quedaría acoplada al esquema de BD (cambiar una columna rompería la API)
- Se podría exponer accidentalmente el `hashed_password` o campos internos
- Las validaciones de API no podrían diferir de las restricciones de BD

### Ejemplo

```go
// models/user.go — Representa la tabla en la BD
type User struct {
    ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
    Email            string    `gorm:"uniqueIndex;not null"`
    FullName         string    `gorm:"not null"`
    HashedPassword   string    `gorm:"not null"`   // ← Campo interno, NUNCA en la API
    IsActive         bool      `gorm:"default:true"`
    IsEmailVerified  bool      `gorm:"default:false"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

// dto/auth.go — Define la forma del request de registro
type RegisterRequest struct {
    Email    string `json:"email"     validate:"required,email,max=255"`
    FullName string `json:"full_name" validate:"required,min=2,max=255"`
    Password string `json:"password"  validate:"required,min=8"` // ← Texto plano, se hashea luego
}

// dto/auth.go — Define la forma del response (sin datos internos)
type UserResponse struct {
    ID              string    `json:"id"`
    Email           string    `json:"email"`
    FullName        string    `json:"full_name"`
    IsActive        bool      `json:"is_active"`
    IsEmailVerified bool      `json:"is_email_verified"`
    CreatedAt       time.Time `json:"created_at"`
    // HashedPassword NO está aquí — nunca se expone al cliente
}
```

---

## 3. Inyección de Dependencias (Dependency Injection)

### ¿Qué es?

Pasar las dependencias de una función/struct desde afuera en lugar de crearlas internamente. Esto hace el código testeable y desacoplado.

### ¿Cómo se aplica en Go?

En Go no hay un framework de DI como FastAPI `Depends()`. La DI se implementa de forma explícita pasando dependencias como parámetros de función o en structs.

```go
// ✅ PATRÓN CORRECTO — Service recibe el *gorm.DB desde afuera
type AuthService struct {
    db     *gorm.DB    // Dependencia de base de datos
    config *config.Config // Dependencia de configuración
}

// Constructor: las dependencias se "inyectan" al crear el service
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
    return &AuthService{db: db, config: cfg}
}

// Los métodos usan las dependencias del struct
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
    // Usa s.db para queries, s.config para configuración
}

// En cmd/api/main.go — Se "inyectan" las dependencias reales
db := database.Connect(cfg)
authSvc := services.NewAuthService(db, cfg)
authHandler := handlers.NewAuthHandler(authSvc)
```

**¿Por qué esto facilita el testing?**

```go
// En tests, se inyecta una BD de prueba en lugar de la real
func TestRegisterUser(t *testing.T) {
    testDB := setupTestDB()  // BD SQLite en memoria o PostgreSQL de test
    authSvc := services.NewAuthService(testDB, testConfig)
    // Los tests no afectan la BD de producción ni dependen de ella
}
```

---

## 4. Patrón Middleware (Chain of Responsibility)

### ¿Qué es?

Una cadena de funciones que procesan la petición HTTP secuencialmente antes de llegar al handler final. Cada middleware decide si continua la cadena (`c.Next()`) o la corta (`c.Abort()`).

### ¿Cómo se aplica en Gin?

```go
// ¿Qué? Cada función gin.HandlerFunc es un eslabón de la cadena.
// ¿Para qué? Separar concerns transversales (auth, logging, rate limiting)
//             del código de negocio en los handlers.
// ¿Impacto? Sin este patrón, cada handler tendría que repetir la lógica
//            de autenticación, rate limiting, etc.

// Middleware 1 — CORS
r.Use(middleware.CORSMiddleware(cfg.FrontendURL))

// Middleware 2 — Security headers
r.Use(middleware.SecurityHeaders())

// Middleware 3 — Rate limiting (solo para rutas de auth)
authRoutes.Use(middleware.RateLimitMiddleware(authRate))

// Middleware 4 — JWT auth (solo para rutas protegidas)
protected.Use(middleware.AuthMiddleware(cfg.SecretKey))

// Flujo de una petición con todos los middlewares:
// Request → CORS → SecurityHeaders → RateLimit → (si necesario) Auth → Handler
```

---

## 5. JWT Stateless (Token-Based Authentication)

### ¿Qué es?

Un esquema de autenticación donde el servidor no guarda el estado de la sesión. Toda la información necesaria está en el token que el cliente envía en cada request.

### Estructura de un JWT

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.   ← Header (algoritmo)
eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQiLCJleHAiOjE3MDk3NTAwMDB9.   ← Payload (datos)
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c   ← Signature (firma)
```

**Payload decodificado:**

```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000", // user_id
  "exp": 1709750000, // expiración (Unix timestamp)
  "iat": 1709749100 // emitido en (issued at)
}
```

### ¿Por qué dos tokens (access + refresh)?

```
Access Token (15 min)           Refresh Token (7 días)
│                               │
├── Corta vida                  ├── Larga vida
├── Se envía en CADA request    ├── Solo se envía para renovar
├── Si se compromete:           ├── Más valioso si se compromete
│   daño limitado a 15 min      │   (requiere rotación inmediata)
└── No se guarda en BD          └── Puede guardarse o revocarse
```

---

## 6. Context/Provider (React)

### ¿Qué es?

Un patrón de React para compartir estado global (como el usuario autenticado) entre componentes sin prop drilling.

### ¿Cómo se aplica?

```tsx
// context/AuthContext.tsx
// ¿Qué? Contexto que centraliza el estado de autenticación.
// ¿Para qué? Que cualquier componente en el árbol acceda al usuario
//             sin pasar props manualmente por cada nivel.
// ¿Impacto? Sin Context, habría que pasar user, token, logout()
//            como prop a través de 5+ niveles de componentes.
interface AuthContextType {
  user: User | null;
  accessToken: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);

  // Al montar: verificar si hay token guardado y obtener el perfil
  useEffect(() => {
    const token = sessionStorage.getItem("access_token");
    if (token) {
      getMeUser(token)
        .then(setUser)
        .catch(() => logout());
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, accessToken, login, logout, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  );
}
```

---

## 7. Custom Hook (React)

### ¿Qué es?

Una función que empieza con `use` y encapsula lógica de estado reutilizable. Permite extraer lógica compleja de los componentes.

### ¿Cómo se aplica?

```tsx
// hooks/useAuth.ts
// ¿Qué? Hook que consume el AuthContext y lo expone.
// ¿Para qué? Simplificar el acceso al contexto de auth y
//             añadir un mensaje de error claro si se usa fuera del Provider.
// ¿Impacto? Sin este hook, cada componente haría useContext(AuthContext)
//            y tendría que manejar el caso null (cuando no hay Provider).
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe usarse dentro de <AuthProvider>");
  }
  return context;
}

// Uso en cualquier componente:
function DashboardPage() {
  const { user, logout } = useAuth(); // Simple, limpio
  return <div>Bienvenido, {user?.fullName}</div>;
}
```

---

## 8. Interceptor (Axios)

### ¿Qué es?

Una función que se ejecuta automáticamente antes o después de cada petición/respuesta HTTP. Permite agregar lógica transversal (como adjuntar tokens) sin repetir código en cada llamada.

### ¿Cómo se aplica?

```typescript
// api/axios.ts
// ¿Qué? Instancia de Axios con interceptors configurados.
// ¿Para qué? Adjuntar el token automáticamente y manejar renovación de tokens
//             de forma transparente para los llamadores.
// ¿Impacto? Sin interceptors, cada función de api/ tendría que:
//            1. Obtener el token del storage
//            2. Añadirlo al header
//            3. Manejar el error 401 y renovar
//            Eso es duplicación masiva de código.
const apiClient = axios.create({ baseURL: import.meta.env.VITE_API_URL });

// Interceptor de REQUEST — añade el token a cada petición
apiClient.interceptors.request.use((config) => {
  const token = getAccessToken(); // Del contexto o sessionStorage
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor de RESPONSE — maneja token expirado automáticamente
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const newToken = await refreshToken();
      if (newToken) {
        // Reintentar la petición original con el nuevo token
        error.config.headers.Authorization = `Bearer ${newToken}`;
        return apiClient(error.config);
      }
      // Si no se pudo renovar → cerrar sesión
      logout();
    }
    return Promise.reject(error);
  },
);
```

---

## 9. Protected Route (Guard Pattern)

### ¿Qué es?

Un componente de React Router que protege rutas: redirige a login si el usuario no está autenticado.

### ¿Cómo se aplica?

```tsx
// components/layout/ProtectedRoute.tsx
// ¿Qué? Componente wrapper que verifica autenticación antes de renderizar.
// ¿Para qué? Centralizar la lógica de protección de rutas en un solo lugar.
// ¿Impacto? Sin ProtectedRoute, cada página protegida tendría que verificar
//            individualmente si el usuario está autenticado,
//            duplicando código y siendo fácil de olvidar.
function ProtectedRoute() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <div>Cargando...</div>; // Esperar a que se cargue el estado inicial
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />; // Renderizar la ruta protegida
}

// App.tsx — Uso del guard:
<Routes>
  <Route path="/" element={<LandingPage />} />
  <Route path="/login" element={<LoginPage />} />
  <Route element={<ProtectedRoute />}>
    {" "}
    {/* ← wrapper */}
    <Route path="/dashboard" element={<DashboardPage />} />
    <Route path="/change-password" element={<ChangePasswordPage />} />
  </Route>
</Routes>;
```

---

## 10. Monorepo

### ¿Qué es?

Un repositorio único que contiene múltiples proyectos relacionados (backend + frontend).

### ¿Cómo se aplica?

```
proyecto/        ← Un solo repositorio git
├── be/          ← Proyecto Go (módulo Go independiente)
└── fe/          ← Proyecto Node.js (package.json independiente)
```

**Ventajas en este bootcamp:**

- Un solo `git clone` para obtener todo el proyecto
- Los cambios de API y frontend van en el mismo commit
- Documentación centralizada en `_docs/`

**Independencia de herramientas:**

- `be/` usa `go.mod` — no afecta a `fe/`
- `fe/` usa `pnpm` — no afecta a `be/`
- Cada uno tiene su propio `.env`

---

## Resumen de Patrones

| Patrón                  | Dónde se usa              | Beneficio principal                   |
| ----------------------- | ------------------------- | ------------------------------------- |
| Layered Architecture    | `be/internal/`            | Separación de concerns, testabilidad  |
| DTO                     | `be/internal/dto/`        | Desacoplamiento API ↔ BD, seguridad   |
| Dependency Injection    | Services en Go            | Testabilidad, flexibilidad            |
| Middleware Chain        | `be/internal/middleware/` | Reutilización de lógica transversal   |
| JWT Stateless           | Auth completo             | Escalabilidad, sin estado en servidor |
| Context/Provider        | `fe/src/context/`         | Estado global sin prop drilling       |
| Custom Hook             | `fe/src/hooks/`           | Reutilización de lógica de estado     |
| Interceptor             | `fe/src/api/axios.ts`     | Manejo automático de tokens           |
| Protected Route (Guard) | `fe/src/components/`      | Protección centralizada de rutas      |
| Monorepo                | Raíz del proyecto         | Cohesión del proyecto completo        |
