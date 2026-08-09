# API JeussAirel - IAM en Go

API REST de Identity and Access Management (IAM) construida con **Go + Gin + PostgreSQL + GORM**.

Incluye autenticación JWT con refresh tokens, recuperación de contraseñas, CRUD de usuarios/roles/funcionalidades,
asignación de roles, sobreescritura de permisos por usuario (**scope overrides**) y auditoría de inicios de sesión,
empleando **RBAC** (control de acceso basado en roles).

## Requisitos

- **Go 1.21+** instalado y en el PATH.
- **PostgreSQL** accesible (ver [Configuración](#configuración)).
- Docker (opcional, para levantar la BD con `docker-infra`).

## Cómo ejecutar la API

```bash
# 1. Copiar la configuración de ejemplo
cp .env.example .env

# 2. Ajustar .env (ver tabla de configuración abajo)

# 3. Instalar dependencias (la primera vez)
go mod tidy

# 4. Levantar el servidor
go run .

#   o bien compilar el binario y ejecutarlo:
go build -o api-server .
./api-server
```

Al iniciar se valida la conexión a la base de datos (se crean las tablas automáticamente vía migración
automática de GORM) y se imprime en los logs:

```
conexión a la base de datos establecida
servidor iniciado ... port=8080 env=develop
```

El servidor escucha por defecto en **http://localhost:8080**. Health check:
`GET /api/v1/health` → `{"status":"ok"}`.

Para detenerlo: `Ctrl+C` (apagado ordenado en los segundos configurados en `SERVER_SHUTDOWN_GRACE`).

> Nota para PowerShell (Windows): al leer/escribir JSON con tokens, evita la interpolación de `$`
> escribiendo los cuerpos en archivos y usando `--data-binary "@archivo.json"`.

## Credenciales de ejemplo para probar endpoints

Un usuario administrador viene precargado para probar la API:

| Email                | Password    | Permisos que otorga                                              |
|----------------------|-------------|------------------------------------------------------------------|
| `admin@sena.edu.co`  | `Admin#2026`| Leer/listar usuarios, roles, features, auditoría y sesiones; crear roles, asignar features y roles |

El `access_token` que devuelve `POST /auth/login` se envía en el header `Authorization: Bearer <token>`
en el resto de endpoints. Las que están protegidas por RBAC devuelven `403 FORBIDDEN` si el rol no
posee la feature/acción requerida (por ejemplo, administrar scope overrides requiere
`IDENTITY_SCOPE_MANAGE` y por defecto el admin no la tiene).

## Documentación Swagger (OpenAPI)

La API expone su documentación interactiva (Swagger UI) en **http://localhost:8080/swagger/index.html**
cuando el servidor está corriendo. También en formato crudo en `/swagger/doc.json` y `/swagger/swagger.yaml`.

- Todos los endpoints están documentados con su método, ruta, parámetros, cuerpos de entrada/salida
  y los permisos RBAC requeridos.
- Para probar endpoints protegidos usa el botón **Authorize** de Swagger con el `access_token`
  obtenido en `/auth/login` (formato `Bearer <token>`).
- Los docs se generan desde los comentarios del código con el paquete `github.com/swaggo/swag`.
  Tras cambiar cualquier anotación `@` en los handlers, regénéralos con:

```bash
swag init -g main.go -d ./ -o ./docs
```

## Configuración

Copia `.env.example` a `.env` y ajusta los valores:

| Variable                  | Default                | Descripción |
|---------------------------|------------------------|-------------|
| `APP_ENV`                 | `develop`              | `develop` o `production` (cambia el modo de Gin) |
| `APP_PORT`                | `8080`                 | Puerto del servidor |
| `LOG_LEVEL`               | (vacío)                | Nivel de log (`debug`, `info`, `warn`, `error`) |
| `SERVER_READ_TIMEOUT`     | `10s`                  | Timeout de lectura HTTP |
| `SERVER_WRITE_TIMEOUT`    | `15s`                  | Timeout de escritura HTTP |
| `SERVER_IDLE_TIMEOUT`     | `60s`                  | Timeout de conexiones idle |
| `SERVER_SHUTDOWN_GRACE`   | `10s`                  | Tiempo de apagado ordenado |
| `CORS_ALLOWED_ORIGINS`    | `*`                    | Orígenes permitidos (separados por coma) |
| `DB_HOST` / `DB_PORT`     | `localhost` / `5432`  | Host y puerto de PostgreSQL |
| `DB_USER` / `DB_PASSWORD` | `postgres` / (vacío)   | Credenciales de la BD |
| `DB_NAME`                 | `postgres`              | Nombre de la base de datos |
| `DB_SSLMODE`              | `disable`               | Modo SSL de PostgreSQL |
| `DB_MAX_OPEN_CONNS`       | `25`                    | Pool: conexiones máximas abiertas |
| `DB_MAX_IDLE_CONNS`       | `5`                     | Pool: conexiones idle |
| `DB_CONN_MAX_LIFETIME`    | `5m`                    | Vida máxima de una conexión |
| `JWT_SECRET`              | (vacío)                | **Secreto de firma de tokens JWT (¡crítico en producción!)** |
| `JWT_ACCESS_TTL`          | `15m`                   | Vigencia del token de acceso |
| `JWT_REFRESH_TTL`         | `168h`                  | Vigencia del refresh token (7 días) |

## Estructura del proyecto

```
api-jeussairel/
├── main.go              # Punto de entrada: carga config, logger, conecta BD y arranca el server
├── router.go            # Construye el router Gin: middlewares globales + montaje de rutas por módulo
├── go.mod / go.sum      # Definición del módulo y dependencias
├── .env                 # Variables de entorno (no debe commitearse)
├── .env.example         # Variables de entorno de ejemplo
├── docs/                # Documentación Swagger generada (docs.go, swagger.json, swagger.yaml)
├── config/config.go     # Carga de configuración desde variables de entorno (con defaults)
├── logger/logger.go     # Logger estructurado (slog)
├── database/database.go # Conexión a PostgreSQL con GORM + pool + migración
├── middleware/
│   ├── auth.go           # RequireAuth: valida JWT Bearer y deja claims en el contexto
│   ├── authorization.go  # RequirePermission: autorización RBAC por feature + acción
│   ├── cors.go           # CORS configurable
│   ├── recovery.go       # Recuperación de pánicos (500) con log
│   ├── logger.go         # Logger estructurado de peticiones (método, ruta, status, latencia)
│   └── notfound.go       # Respuestas JSON para 404 / 405
├── shared/
│   ├── apperror/         # Errores de dominio con código y status HTTP
│   ├── jwt/              # Generación, parseo y expiración de tokens JWT
│   ├── bcrypt/           # Hash y verificación de contraseñas
│   ├── response/         # Envoltura estándar de respuestas JSON (Success / Error / AppError)
│   ├── validation/      # Helpers de validación de entrada (email, contraseña)
│   ├── utils/           # Útiles (id de usuario a string, IP del cliente)
│   └── uuid/            # Generación de UUID v4 (PKs)
├── auth/                # Módulo: login, refresh, logout, recuperación de contraseña
├── user/                # Módulo: CRUD de usuarios
├── role/                # Módulo: CRUD de roles, features asignadas y scope overrides
├── feature/             # Módulo: CRUD de módulos y funcionalidades (features)
├── session/             # Módulo: sesiones activas (refresh tokens) y revocación
└── audit/               # Módulo: auditoría de inicios de sesión
```

## Responsabilidades por capa

Cada módulo (auth, user, role, feature, session, audit) sigue la misma arquitectura en 6 archivos:

| Capa    | Archivo         | Responsabilidad |
|---------|-----------------|-----------------|
| Rutas   | `routes.go`     | Registra los endpoints del módulo y los protege con middlewares (auth y permisos) |
| HTTP    | `handler.go`    | Recibe la petición, valida la entrada (Gin binding), llama al service y devuelve la respuesta |
| Negocio | `service.go`    | Toda la lógica de negocio: flujos, validaciones, transacciones |
| Datos   | `repository.go` | Acceso a PostgreSQL **únicamente** con GORM (sin lógica de negocio) |
| Modelo  | `model.go`      | Struct que representa una tabla de la BD (con `TableName()` amarcante) |
| DTO     | `dto.go`        | Estructuras de entrada (request) y salida (response) del API |

Flujo de una petición:

```
HTTP Request → Router → Middlewares (Logger, Recovery, CORS)
            → RequireAuth (JWT) → RequirePermission (RBAC) → Handler → Service → Repository → PostgreSQL
            → response JSON (Body estándar {success, message, data|error})
```

- **No hay** lógica de negocio en los handlers ni SQL fuera de los repositorios.
- Los números entre capas son paquetes `shared/` (respuestas, errores, JWT, bcrypt, uuid).
- El `Handler` nunca habla con la BD ni hace queries; siempre delega en el `Service`.
- `Repository` no conoce el contexto HTTP: recibe datos concretos y devuelve estructuras.

### Middlewares en detalle

| Middleware              | Propósito |
|-------------------------|-----------|
| `RequireAuth(secret)`       | Exige header `Authorization: Bearer <token>`, valida el JWT y deja `user_id`, `email`, `actor_type` en el contexto. Respuestas: 401 `MISSING_TOKEN / INVALID_TOKEN / TOKEN_EXPIRED`. |
| `RequirePermission(db, featureCode, action)` | Autoriza por RBAC: consulta el feature vigente, chequea scope overrides del usuario (pueden **denegar** aunque el rol lo permita) y luego los roles del usuario (`user_role` → `role_feature`). Acciones: `READ < WRITE < DELETE < PUBLISH < APPROVE` (una acción concedida permite las inferiores). Respuesta 403 si no hay permiso. |
| `Logger(log)`           | Log de cada petición: método, ruta, status, latencia, cliente. |
| `Recovery()`         | Captura pánicos y responde 500. |
| `CORS(cfg)`         | Configura `Access-Control-Allow-Origin` según `CORS_ALLOWED_ORIGINS`. |

## Endpoints por módulo

Todas las rutas están bajo `/api/v1`. Las marcadas con 🔒 requieren un token válido
(`Authorization: Bearer <access_token>`); las que indican permiso además exigen la feature/acción del RBAC.

**Auth** (dos públicos) — `auth/routes.go`:

| Método | Ruta                        | Body                    | Permiso |
|--------|-----------------------------|-------------------------|---------|
| POST   | `/auth/login`               | `{email, password}`     | Público |
| POST   | `/auth/refresh`             | `{refresh_token}`       | Público |
| POST   | `/auth/logout`              | `{refresh_token}`       | 🔒 |
| POST   | `/auth/forgot-password`     | `{email}`               | Público |
| POST   | `/auth/reset-password`      | `{token, new_password}` | Público |

**Users** — `user/routes.go`:

| Método | Ruta               | Permiso                            |
|--------|--------------------|------------------------------------|
| GET    | `/users`           | 🔒 `IDENTITY_USER_VIEW` (read)     |
| POST   | `/users`           | 🔒 `IDENTITY_USER_MANAGE` (write)  |
| GET    | `/users/:id`       | 🔒 `IDENTITY_USER_VIEW` (read)     |
| PUT    | `/users/:id`       | 🔒 `IDENTITY_USER_MANAGE` (write)  |
| DELETE | `/users/:id`       | 🔒 `IDENTITY_USER_MANAGE` (write)  |

Se parsean `page` y `page_size` (clamp: `page >= 1`, `1 <= page_size <= 100`).

**Roles** — `role/routes.go`:

| Método | Ruta                                 | Permiso                            |
|--------|--------------------------------------|------------------------------------|
| GET    | `/roles`                             | 🔒 `IDENTITY_ROLE_VIEW` (read)     |
| POST   | `/roles`                             | 🔒 `IDENTITY_ROLE_MANAGE` (write)  |
| GET    | `/roles/:id`                         | 🔒 `IDENTITY_ROLE_VIEW` (read)     |
| PUT    | `/roles/:id`                         | 🔒 `IDENTITY_ROLE_MANAGE` (write)  |
| DELETE | `/roles/:id`                         | 🔒 `IDENTITY_ROLE_MANAGE` (write)  |
| GET    | `/roles/:id/features`                | 🔒 `IDENTITY_ROLE_VIEW` (read)     |
| PUT    | `/roles/:id/features`                | 🔒 `IDENTITY_ROLE_MANAGE` (write)  |
| POST   | `/users/:id/roles`                   | 🔒 `IDENTITY_ROLE_ASSIGN` (write)  |
| GET    | `/users/:id/roles`                   | 🔒 `IDENTITY_ROLE_VIEW` (read)     |
| DELETE | `/users/:id/roles/:userRoleId`       | 🔒 `IDENTITY_ROLE_ASSIGN` (write)  |
| POST   | `/scope-overrides`                   | 🔒 `IDENTITY_SCOPE_MANAGE` (write) |
| GET    | `/scope-overrides/users/:id`         | 🔒 `IDENTITY_SCOPE_MANAGE` (write) |
| DELETE | `/scope-overrides/:id`               | 🔒 `IDENTITY_SCOPE_MANAGE` (write) |

**Modules y Features** — `feature/routes.go` (todo bajo `IDENTITY_ROLE_MANAGE`):

| Método | Ruta                              | Permiso                            |
|--------|-----------------------------------|------------------------------------|
| GET   | `/modules`                        | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| POST   | `/modules`                        | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| GET    | `/modules/:id`                    | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| PUT    | `/modules/:id`                    | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| DELETE | `/modules/:id`                    | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| GET    | `/modules/:id/features`           | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| GET    | `/features`                       | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| POST   | `/features`                       | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| GET    | `/features/:id`                   | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| PUT    | `/features/:id`                   | 🔒 `IDENTITY_ROLE_MANAGE` (write) |
| DELETE | `/features/:id`                   | 🔒 `IDENTITY_ROLE_MANAGE` (write) |

**Sesiones** — `session/routes.go` (solo autenticación):

| Método | Ruta                 | Descripción |
|--------|----------------------|-------------|
| GET    | `/sessions/active`   | Sesiones (refresh tokens) activas del usuario autenticado |
| DELETE | `/sessions/:id`      | Revoca una sesión del usuario autenticado |

**Auditoría** — `audit/routes.go`:

| Método | Ruta                          | Permiso                     |
|--------|-------------------------------|-----------------------------|
| GET   | `/audit/logins`               | 🔒 `AUDIT_LOG_VIEW` (read)  |
| GET   | `/audit/logins/users/:userId` | 🔒 `AUDIT_LOG_VIEW` (read)  |

Todas las respuestas usan el envoltorio estándar:

```json
{ "success": true,  "message": "...", "data": {...} }
{ "success": false, "message": "...", "error": { "code": "..." } }
```

Códigos de error comunes: `INVALID_CREDENTIALS`, `UNAUTHORIZED`, `FORBIDDEN`, `TOKEN_EXPIRED`,
`NOT_FOUND`, `VALIDATION_ERROR`, `INTERNAL`.

## Ejemplos con curl

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Login (usuario demo)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sena.edu.co","password":"Admin#2026"}'
# → data.access_token  y  data.refresh_token

# Listar usuarios (protegido)
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <access_token>"

# Obtener un usuario por id
curl http://localhost:8080/api/v1/users/<uuid> \
  -H "Authorization: Bearer <access_token>"

# Listar roles
curl http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer <access_token>"

# Renovar tokens
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'

# Auditoría de logins
curl http://localhost:8080/api/v1/audit/logins?page=1&page_size=20 \
  -H "Authorization: Bearer <access_token>"

# Swagger UI
curl http://localhost:8080/swagger/index.html
```

En **PowerShell (Windows)** escribe el cuerpo en un archivo y usa `-InFile`:

```powershell
$body = '{"email":"admin@sena.edu.co","password":"Admin#2026"}'
Set-Content -Path login.json -Value $body
Invoke-RestMethod -Uri http://localhost:8080/api/v1/auth/login -Method Post `
  -ContentType "application/json" -InFile login.json
```

## Dependencias

| Dependencia                    | Uso                          |
|--------------------------------|------------------------------|
| `github.com/gin-gonic/gin`     | Framework HTTP y routers     |
| `gorm.io/gorm`                  | ORM                          |
| `gorm.io/driver/postgres`      | Driver de PostgreSQL         |
| `github.com/golang-jwt/jwt/v5` | Emisión y validación de JWT  |
| `golang.org/x/crypto`          | Bcrypt (hash de contraseñas) |
| `github.com/joho/godotenv`     | Lectura de variables `.env`  |
| `github.com/swaggo/swag`      | Generación de docs OpenAPI   |
| `github.com/swaggo/gin-swagger`| Middleware de Swagger UI     |
| `github.com/swaggo/files`     | Archivos estáticos de Swagger UI |

## Pruebas

Los paquetes de `shared/` incluyen tests unitarios (ver `shared/bcrypt/bcrypt_test.go` y `shared/jwt/jwt_test.go`):

```bash
go test ./...       # todos
go test ./shared/...
```

## Notas de seguridad

- **Nunca** subas el archivo `.env` a repositorio.
- Cambia `JWT_SECRET` antes de producción.
- En `production` el router activa el modo release de Gin (`APP_ENV=production`).