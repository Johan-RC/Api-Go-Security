# Documentación del Proyecto — Sistema de Identidad y Accesos (IAM) JeuxAirel

> Proyecto realizado como desafío técnico: API REST en **Go**, base de datos **PostgreSQL** diseñada por el instructor (entregada como script Liquibase), y aplicación móvil en **React Native (Expo)**.
> Este documento explica, de punta a punta, cómo se construyó todo: desde la base de datos que entregó el instructor, hasta cómo el front consume los endpoints del backend.

---

## Tabla de contenido

1. [Resumen general y arquitectura](#1-resumen-general-y-arquitectura)
2. [Carpetas del proyecto](#2-carpetas-del-proyecto)
3. [La base de datos entregada por el instructor](#3-la-base-de-datos-entregada-por-el-instructor)
4. [Cómo se construyó el backend a partir de la BD](#4-cómo-se-construyó-el-backend-a-partir-de-la-bd)
5. [Arquitectura del backend (capas y flujo de una petición)](#5-arquitectura-del-backend-capas-y-flujo-de-una-petición)
6. [Middlewares: autenticación y autorización RBAC](#6-middlewares-autenticación-y-autorización-rbac)
7. [Endpoints del backend (todos los módulos)](#7-endpoints-del-backend-todos-los-módulos)
8. [Envío de correos](#8-envío-de-correos)
9. [Cómo se implementaron las conexiones](#9-cómo-se-implementaron-las-conexiones)
10. [Cómo se conectó el frontend al backend](#10-cómo-se-conectó-el-frontend-al-backend)
11. [Pantallas y funcionalidad del frontend](#11-pantallas-y-funcionalidad-del-frontend)
12. [Validaciones: frontend vs backend vs base de datos](#12-validaciones-frontend-vs-backend-vs-base-de-datos)
13. [Seguridad implementada](#13-seguridad-implementada)
14. [Cómo ejecutar todo el sistema](#14-cómo-ejecutar-todo-el-sistema)
15. [Instrumentos para la exposición y posibles preguntas](#15-instrumentos-para-la-exposición-y-posibles-preguntas)
16. [Glosario rápido](#16-glosario-rápido)

---

## 1. Resumen general y arquitectura

El sistema es un **módulo de Identidad y Acceso (IAM)** dentro de un ecosistema educativo (SENA). Controla quién puede entrar, qué puede ver y qué puede hacer cada usuario mediante **RBAC** (control de acceso basado en roles).

La arquitectura es de **3 piezas separadas que dialogan por HTTP**:

```
┌──────────────────┐      HTTP/JSON       ┌──────────────────┐       SQL       ┌─────────────────────┐
│   Frontend       │ ───────────────────► │   Backend (Go)   │ ───────────────► │   PostgreSQL 16     │
│  React Native    │   REST /api/v1       │  Gin + GORM      │    (GORM)        │  BD del instructor  │
│  (Expo, TS)      │ ◄─────────────────── │                  │ ◄─────────────── │  (Liquibase)        │
└──────────────────┘  JSON {success,...}  └──────────────────┘                  └─────────────────────┘
```

- **Frontend**: app móvil en **React Native** con **Expo SDK 57**, TypeScript, React Navigation (5 tabs). Corre en el puerto `8081`.
- **Backend**: API REST en **Go**, framework **Gin**, ORM **GORM**, JWT (HS256), bcrypt. Corre en el puerto `8080`.
- **Base de datos**: **PostgreSQL 16** con los esquemas `identity`, `rbac`, `rbac_catalog`, `session`, `identity_audit`. Una sola BD para todos los módulos.

**Tecnologías principales (tabla para la exposición):**

| Capa | Tecnología | Para qué |
|------|-----------|----------|
| Backend | Go 1.25 + Gin | API REST |
| Backend | GORM + driver postgres | ORM y acceso a datos |
| Backend | golang-jwt v5 | Tokens de acceso y refresco |
| Backend | golang.org/x/crypto (bcrypt) | Hash de contraseñas |
| Backend | godotenv | Variables de entorno `.env` |
| Backend | swaggo/swag + gin-swagger | Documentación OpenAPI (Swagger UI) |
| BD | PostgreSQL 16 | Almacenamiento |
| BD | Liquibase | Migraciones versionadas (scripts del instructor) |
| Frontend | React Native + Expo | App móvil / web |
| Frontend | Axios | Cliente HTTP |
| Frontend | AsyncStorage | Guardado seguro de tokens |
| Frontend | React Navigation | Navegación por tabs y stacks |
| Infra | Docker + Docker Compose | Contenedores de BD, API, MailHog |

---

## 2. Carpetas del proyecto

```
Api-Go-Security/
├── API-Go/                       # Backend en Go (la API IAM)
│   ├── main.go                   # Punto de entrada (config, BD, servidor, apagado)
│   ├── router.go                 # Router Gin: middlewares globales + montaje de módulos
│   ├── config/config.go          # Lectura de variables de entorno con defaults
│   ├── database/database.go      # Conexión PostgreSQL vía GORM + pool
│   ├── middleware/               # RequireAuth, RequirePermission, CORS, Logger, Recovery
│   ├── shared/                   # Paquetes: jwt, bcrypt, response, apperror, validation, uuid
│   ├── auth/                     # Register, VerifyEmail, Login, Refresh, Logout, Forgot/Reset, Me
│   ├── user/                     # CRUD de usuarios
│   ├── role/                     # CRUD de roles, features asignadas y scope overrides
│   ├── feature/                  # CRUD de módulos y funcionalidades
│   ├── session/                  # Sesiones activas (refresh tokens) y revocación
│   ├── audit/                    # Auditoría de inicios de sesión
│   └── docs/                     # Swagger generado
│
├── design-software-iam-db-main/  # BD entregada por el instructor (Liquibase)
│   ├── changelog/                # changelog-master.yaml
│   ├── 01_ddl/                   # DDL: esquemas, tablas, índices (versionados)
│   ├── 02_dml/                   # DML: seed de roles, features y usuario admin
│   └── 03_dcl, 04_tcl/           # Permisos y transacciones (si aplica)
│
├── jeuxairel-mobile/             # Frontend React Native (Expo)
│   ├── App.tsx                   # Proveedores + NavigationContainer
│   └── src/
│       ├── api/                  # client (axios), endpoints, tokenStore, errors
│       ├── contexts/AuthContext  # Estado global de sesión + can()/hasRole()
│       ├── navigation/           # RootNavigator (tabs + stacks)
│       ├── screens/              # auth, home, users, roles, catalog, profile
│       ├── hooks/                # useForm, usePaginatedList
│       ├── validation/           # Reglas de validación de formularios
│       ├── types/                # Modelos TypeScript (espejo del backend)
│       └── theme/                # Colores, espaciado, tipografía (verde SENA)
│
└── docker-infra/                 # Docker Compose (BD, API, MailHog) + Dockerfile
    ├── docker-compose.yml
    ├── api-go.Dockerfile
    └── .env.develop              # Variables de entorno de Postgres
```

---

## 3. La base de datos entregada por el instructor

El instructor entregó un proyecto **Liquibase** (`design-software-iam-db-main`) con los cambios DDL versionados en `01_ddl/` y los datos semilla en `02_dml/`. Esa es la **fuente de verdad del modelo de datos**; el backend **no crea sus propias tablas**, respeta exactamente las que define el instructor.

### 3.1 Esquemas creados

| Esquema | Contenido |
|---------|-----------|
| `identity` | Usuarios (`identity.user`) |
| `rbac` | Roles, features por rol, roles por usuario y overrides de scope |
| `rbac_catalog` | Catálogo de módulos y funcionalidades (features) |
| `session` | Refresh tokens, reseteos de contraseña y códigos de verificación de email |
| `identity_audit` | Auditoría de inicios de sesión |

### 3.2 Tablas principales

**`identity.user`** — un usuario puede ser `USER`, `INSTRUCTOR` o `LEARNER` (restricción `CHECK`). Tiene `password_hash`, `is_active`, `failed_attempts`, `locked_until`, `last_login_at`.

**`rbac_catalog.module`** y **`rbac_catalog.feature`** — catálogo de permisos. Un feature tiene una acción base: `READ`, `WRITE`, `DELETE`, `PUBLISH` o `APPROVE` (restricción `CHECK`). Ejemplos: `IDENTITY_USER_VIEW`, `IDENTITY_ROLE_MANAGE`, `AUDIT_LOG_VIEW`.

**`rbac.role`** — roles del sistema (`SYSTEM_ADMIN`, `CENTER_DIRECTOR`, `COORDINATOR`, `AREA_LEADER`, `INSTRUCTOR`, `LEARNER`, `ADMIN_STAFF`).

**`rbac.role_feature`** — asigna funcionalidades a roles con un **scope** (`GLOBAL`, `TRAINING_CENTER`, `AREA`, `OWN_FICHAS`, `OWN_SCHEDULE`, `OWN_PROFILE`, `OWN_FICHA_AS_LEARNER`). Esto define qué tan amplio es cada permiso.

**`rbac.user_role`** — asigna roles a usuarios (con quién lo asignó y a qué training center).

**`rbac.user_scope_override`** — casos especiales: un permiso puede **concederse o denegarse** directamente a un usuario, con vigencia (`expires_at`).

**`session.refresh_token`**, **`session.password_reset_request`**, **`session.email_verification_code`** — sesiones y flujos de recuperación/verificación.

**`identity_audit.audit_login`** — bitácora de accesos con resultado (`SUCCESS`, `INVALID_PASSWORD`, `USER_NOT_FOUND`, `ACCOUNT_LOCKED`, `TOKEN_EXPIRED`).

### 3.3 Datos semilla (seed)

Los archivos `02_dml/00_inserts/` se encargan de precargar:

- **`001_seed_rbac.sql`**: 10 módulos, ~62 funcionalidades y la asignación de features a cada rol (idempotente con `ON CONFLICT DO NOTHING`).
- **`002_seed_demo.sql`**: el usuario administrador `System Admin` con correo `admin@sena.edu.co` y su rol `SYSTEM_ADMIN`.

> Importante para la exposición: la BD **define desde el inicio** qué puede hacer cada rol. Por ejemplo, `SYSTEM_ADMIN` tiene todo lo de identidad (`IDENTITY_USER_*`, `IDENTITY_ROLE_*`), mientras que `INSTRUCTOR` y `LEARNER` solo tienen permisos funcionales de su ámbito (ver propia ficha, propio horario, perfil propio) y **nunca** los de administración de usuarios/roles.

---

## 4. Cómo se construyó el backend a partir de la BD

El proceso fue **dirigido por el esquema (schema-first)**:

1. **Se leyó el esquema Liquibase** y se extrajeron entidades: usuario, rol, feature, role_feature, user_role, scope override, sesiones, auditoría.
2. **Se mapeó cada tabla a una struct de Go** en `model.go` de cada módulo con `TableName()` indicando esquema y tabla exactos (ej. `identity.user`, `rbac.role`).
3. **Se construyeron los DTOs** (request/response) como "contrato" JSON que ve el front.
4. **Se implementaron los repositorios** con GORM usando las rutas reales de la BD (schema.tabla), respetando índices y restricciones.
5. **Se definieron las features del catálogo como constantes de permisos** que el middleware lee desde `rbac_catalog.feature`.
6. **Se reutilizaron los seeds** para probar los endpoints con el usuario admin y los roles precargados.

Ejemplo de mapeo — `user/model.go` (resumen):

```go
type User struct {
    ID            string    `gorm:"column:id;type:uuid;default:gen_random_uuid()"`
    Email         string    `gorm:"column:email;type:varchar(255);uniqueIndex"`
    PasswordHash  string    `gorm:"column:password_hash"`
    ActorType     string    `gorm:"column:actor_type"`
    IsActive      bool      `gorm:"column:is_active"`
    FailedAttempts int16
    LockedUntil   *time.Time
    ...
}
func (User) TableName() string { return "identity.user" }
```

> Clave de diseño: **la API no re-crea la BD**. `database.Connect` solo abre el pool, hace ping con reintento y deja que todo el DDL venga de Liquibase (el comando `liquibase update` se ejecuta con Docker). Así se garantiza que backend y BD siempre coincidan con lo que pidió el instructor.

---

## 5. Arquitectura del backend (capas y flujo de una petición)

Cada módulo (`auth`, `user`, `role`, `feature`, `session`, `audit`) usa **6 archivos con una sola responsabilidad**:

| Capa | Archivo | Responsabilidad |
|------|---------|-----------------|
| Rutas | `routes.go` | Registra endpoints y los protege con middlewares |
| HTTP | `handler.go` | Recibe la petición, valida el body (Gin binding) y responde |
| Negocio | `service.go` | Toda la lógica (flujos, validaciones, transacciones) |
| Datos | `repository.go` | SQL con GORM, sin lógica de negocio |
| Modelo | `model.go` | Struct espejo de la tabla |
| DTO | `dto.go` | Entradas/salidas JSON + tags de validación |

**Flujo completo de una petición protegida:**

```
HTTP Request
  → Router (Gin)
    → Middleware Logger        (método, ruta, status, latencia)
    → Middleware Recovery      (captura pánicos → 500)
    → Middleware CORS          (cabeceras de origen)
    → RequireAuth              (valida JWT Bearer, inyecta user_id)
    → RequirePermission        (RBAC: feature + acción + scope)
      → Handler                (valida body con binding)
        → Service              (lógica de negocio)
          → Repository         (SQL con GORM)
            → PostgreSQL
  ← response JSON con envoltorio estándar
```

**Reglas de oro de la arquitectura:**
- El `Handler` nunca toca la BD; el `Repository` nunca conoce el contexto HTTP.
- Los errores son de **dominio** (`shared/apperror`): cada uno tiene un `code` estable y un status HTTP.
- Todas las respuestas usan el envoltorio `{success, message, data|error}`.

---

## 6. Middlewares: autenticación y autorización RBAC

### 6.1 `RequireAuth` (autenticación) — `middleware/auth.go`

- Lee el header `Authorization: Bearer <token>`.
- Valida firma y expiración del JWT (HS256, `shared/jwt`).
- Si falta → `401 MISSING_TOKEN`; mal formato → `401 INVALID_TOKEN`; expirado → `401 TOKEN_EXPIRED`.
- Deja `user_id`, `email` y `actor_type` en el contexto de Gin.

### 6.2 `RequirePermission` (autorización RBAC) — `middleware/authorization.go`

Es el corazón del control de acceso. Se ejecuta **después** de `RequireAuth` y recibe `(featureCode, action)`:

1. **Resuelve el feature** desde `rbac_catalog.feature` (solo activos).
2. **Revisa overrides de scope** del usuario (`rbac.user_scope_override`):
   - Si hay un override vigente con `is_allowed = false` → **bloquea** (403).
   - Si hay uno con `is_allowed = true` → **permitido** (se salta el resto).
3. **Busca un rol activo del usuario** que tenga ese feature con la acción suficiente: `user_role` → `role_feature` → `feature`.
4. **Compara jerarquía de acciones:** `READ(1) < WRITE(2) < DELETE(3) < PUBLISH(4) < APPROVE(5)`. Tener `WRITE` permite hacer `READ`.

```sql
-- Cómo decide el permiso (resumen de la lógica):
SELECT f.action_level
FROM rbac.user_role ur
JOIN rbac.role_feature rf ON rf.role_id = ur.role_id
JOIN rbac_catalog.feature f  ON f.id = rf.feature_id
WHERE ur.user_id = $1 AND f.id = $2
  AND (ur.expires_at IS NULL OR ur.expires_at > now())
```

> Para la exposición: es fundamental explicar que **frontend y backend aplican la misma regla**. El frontend recibe en `/auth/me` los permisos efectivos y replica la jerarquía de acciones en `AuthContext` (función `can`), pero la **autorización real siempre ocurre en el backend** con este middleware. Si se intenta llamar una ruta sin permiso, la API devuelve `403 FORBIDDEN` aunque la UI trate de ocultarla.

---

## 7. Endpoints del backend (todos los módulos)

Todas las rutas van bajo `/api/v1`. 🔒 = requiere token; entre paréntesis el permiso RBAC.

### 7.1 Auth — `auth/routes.go`

| Método | Ruta | Body | Permiso | Uso |
|--------|------|------|---------|-----|
| POST | `/auth/register` | `{email,password,first_name,last_name,actor_type}` | Público | Registrar cuenta (pendiente de verificación) |
| POST | `/auth/verify-email` | `{email,code}` | Público | Activar cuenta con el código de 6 dígitos |
| POST | `/auth/login` | `{email,password}` | Público | Iniciar sesión → access + refresh token |
| POST | `/auth/refresh` | `{refresh_token}` | Público | Renovar el par de tokens (rotación) |
| POST | `/auth/logout` | `{refresh_token}` | 🔒 | Revocar sesión |
| POST | `/auth/forgot-password` | `{email}` | Público | Enviar código de recuperación |
| POST | `/auth/reset-password` | `{code,new_password}` | Público | Cambiar contraseña |
| GET | `/auth/me` | — | 🔒 | Datos del usuario + roles + permisos efectivos |

### 7.2 Users — `user/routes.go`

| Método | Ruta | Permiso |
|--------|------|---------|
| GET | `/users` | `IDENTITY_USER_VIEW` (read) |
| POST | `/users` | `IDENTITY_USER_MANAGE` (write) |
| GET | `/users/:id` | `IDENTITY_USER_VIEW` (read) |
| PUT | `/users/:id` | `IDENTITY_USER_MANAGE` (write) |
| DELETE | `/users/:id` | `IDENTITY_USER_MANAGE` (write) |

Paginación por query params `page` y `page_size` (clamp: `page >= 1`, `1 <= page_size <= 100`).

### 7.3 Roles y asignaciones — `role/routes.go`

| Método | Ruta | Permiso |
|--------|------|---------|
| GET | `/roles` | `IDENTITY_ROLE_VIEW` (read) |
| POST | `/roles` | `IDENTITY_ROLE_MANAGE` (write) |
| GET | `/roles/:id` | `IDENTITY_ROLE_VIEW` (read) |
| PUT | `/roles/:id` | `IDENTITY_ROLE_MANAGE` (write) |
| DELETE | `/roles/:id` | `IDENTITY_ROLE_MANAGE` (write) |
| GET | `/roles/:id/features` | `IDENTITY_ROLE_VIEW` (read) |
| PUT | `/roles/:id/features` | `IDENTITY_ROLE_MANAGE` (write) |
| POST | `/users/:id/roles` | `IDENTITY_ROLE_ASSIGN` (write) |
| GET | `/users/:id/roles` | `IDENTITY_ROLE_VIEW` (read) |
| DELETE | `/users/:id/roles/:userRoleId` | `IDENTITY_ROLE_ASSIGN` (write) |
| POST | `/scope-overrides` | `IDENTITY_SCOPE_MANAGE` (write) |
| GET | `/scope-overrides/users/:id` | `IDENTITY_SCOPE_MANAGE` (write) |
| DELETE | `/scope-overrides/:id` | `IDENTITY_SCOPE_MANAGE` (write) |

### 7.4 Módulos y funcionalidades — `feature/routes.go`

Todo bajo `IDENTITY_ROLE_MANAGE` (write): CRUD de `/modules`, `/features` y sus relaciones.

### 7.5 Sesiones — `session/routes.go`

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/sessions/active` | Sesiones activas del usuario autenticado |
| DELETE | `/sessions/:id` | Revoca una sesión del usuario autenticado |

### 7.6 Auditoría — `audit/routes.go`

| Método | Ruta | Permiso |
|--------|------|---------|
| GET | `/audit/logins` | `AUDIT_LOG_VIEW` (read) |
| GET | `/audit/logins/users/:userId` | `AUDIT_LOG_VIEW` (read) |

### 7.7 Health — `router.go`

`GET /api/v1/health` → `{"status":"ok"}` (público, para monitorear el servicio).

### 7.8 Formato de respuesta estándar

```json
{ "success": true,  "message": "...", "data": {...} }
{ "success": false, "message": "...", "error": { "code": "INVALID_CREDENTIALS" } }
```

---

## 8. Envío de correos

El módulo `mailer` implementa envío transaccional con la librería estándar de Go (`net/smtp`).

### 8.1 Flujo de verificación de email (registro)

1. `POST /auth/register` crea el usuario con `is_active = false`.
2. El backend genera un **código aleatorio de 6 dígitos** (`crypto/rand`) y guarda solo su **hash SHA-256** en `session.email_verification_code` (caduca en 15 min).
3. El código se envía al correo del usuario.
4. `POST /auth/verify-email` recibe email + código, compara el hash, y en una **transacción** activa la cuenta y marca el código como usado.
5. Hasta verificarse, `login` responde `403 USER_INACTIVE`.

### 8.2 Flujo de recuperación de contraseña

1. `POST /auth/forgot-password` genera un código y guarda su hash en `session.password_reset_request` (caduca en 15 min). **Invalida intentos anteriores** del mismo usuario.
2. Respuesta intencionalmente igual exista o no el usuario (anti enumeración de cuentas).
3. `POST /auth/reset-password`, en una transacción: valida el código, actualiza el hash de la contraseña (bcrypt), revoca todas las sesiones y marca el reset como usado.

### 8.3 Protocolos SMTP soportados

| Caso | Puerto | Método |
|------|--------|--------|
| Gmail (producción/demo) | 587 | STARTTLS + autenticación `smtp.PlainAuth` |
| Gmail (SSL implícito) | 465 | TLS directo |
| Local (MailHog/dev) | 1025 | Sin TLS ni credenciales |

- En **modo desarrollo sin SMTP**, el backend **devuelve el código en la respuesta** (`dev_code`) y lo registra en el log, para poder probar el flujo sin tener correo real.
- El correo es texto plano RFC 5322 con `Content-Type: text/plain; charset=UTF-8`.
- Configuración via `.env`: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM`.

### 8.4 El frontend y los correos

- `VerifyEmailScreen`: pide el código de 6 dígitos; si venía `dev_code` en la respuesta, se precarga.
- `ForgotPasswordScreen`: si el backend devuelve el código (modo dev), navega directo a `ResetPassword` con el código prellenado.
- `ResetPasswordScreen`: pide código + nueva contraseña + confirmación.

---

## 9. Cómo se implementaron las conexiones

### 9.1 Backend → Base de datos

`database/database.go`:

- Arma el DSN: `host, port, user, password, dbname, sslmode`.
- Abre GORM con `postgres.Open`.
- Configura **pool de conexiones**: `DB_MAX_OPEN_CONNS=25`, `DB_MAX_IDLE_CONNS=5`, `DB_CONN_MAX_LIFETIME=5m`.
- **Ping con reintentos** (3 intentos, 2 s de espera) para tolerar el arranque del contenedor de BD.
- Nivel de log SQL según entorno.

### 9.2 Orquestación con Docker

`docker-infra/docker-compose.yml` levanta:

- `postgres` (port 5432): BD `design-software-develop`, usuario `design_software_user`.
- `api` (port 8080): construye `API-Go` con el Dockerfile multi-etapa (`golang:1.25-alpine` → `alpine:3.20`), lee `API-Go/.env`, espera que el healthcheck de la BD esté "healthy".
- `mailhog` (ports 1025/8025): buzón de prueba para validar envíos sin depender de Gmail.

Comandos usados:

```bash
# Levantar BD + API
docker compose --env-file .env.develop up -d postgres api

# Aplicar migraciones de la BD del instructor
docker compose --profile tooling run --rm liquibase-iam update

# Reconstruir la API tras cambios
docker compose --env-file .env.develop up -d --build api
```

> Nota práctica: el archivo de entorno se llama `.env.develop`, por eso se pasa `--env-file .env.develop` para que Compose lea `POSTGRES_DB/USER/PASSWORD`.

### 9.3 Frontend → Backend

El cliente Axios (`src/api/client.ts`) usa:

```ts
export const API_BASE_URL =
  process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1';
```

- Emulador Android: `http://10.0.2.2:8080/api/v1` (se puede sobrescribir con `.env`).
- **Interceptor de request**: adjunta `Authorization: Bearer <access_token>` desde AsyncStorage a cada petición.
- **Interceptor de respuesta**: ante `401`, intenta renovar el token con `POST /auth/refresh` (single-flight para no disparar varios refrescos a la vez) y reenvía la petición original; si el refresh falla, limpia los tokens y cierra sesión.

### 9.4 CORS

`middleware/cors.go` configura los origenes permitidos desde `CORS_ALLOWED_ORIGINS` (default `*` en desarrollo) para que la app web (Expo en `localhost:8081`) pueda llamar a la API.

---

## 10. Cómo se conectó el frontend al backend

### 10.1 Capa de API (`src/api`)

| Archivo | Función |
|---------|---------|
| `client.ts` | Instancia Axios + base URL + interceptores (token y renovación) |
| `endpoints.ts` | Una función por endpoint; desenvuelve la envoltura `{data}` |
| `tokenStore.ts` | Guarda/lee/limpia access y refresh tokens en AsyncStorage |
| `errors.ts` | Traduce códigos de error del backend a mensajes amigables en español |

Ejemplo de cómo el front llama a un endpoint:

```ts
// endpoints.ts
export function login(email: string, password: string): Promise<TokenPair> {
  return unwrap(client.post('/auth/login', { email, password }));
}
export function me(): Promise<MeResponse> {
  return unwrap(client.get('/auth/me'));
}
export function listUsers(page = 1, pageSize = 20) {
  return unwrap(client.get('/users', { params: { page, page_size: pageSize } }));
}
```

### 10.2 Estado global de sesión (`AuthContext`)

`src/contexts/AuthContext.tsx` es el puente entre backend y la navegación:

- `login(email, password)`: llama `POST /auth/login`, guarda tokens, decodifica el JWT (segmento 2) para tener usuario al instante y luego refresca el perfil con `/auth/me`.
- `refreshProfile()`: obtiene `user`, `roles` y `permissions` reales. Guarda siempre **arreglos** (evita el bug de `permissions: null`).
- `can(code, action)`: replica la jerarquía de acciones del backend (`ACTION_RANK`) sobre `permissions`. Quien tiene `WRITE` puede `READ`.
- `hasRole(name)`: comprueba si el usuario tiene un rol por su código.
- `logout()`: llama `POST /auth/logout` con el refresh token, limpia tokens y estado.
- Expiración: si el refresh no puede renovar el token, `setOnSessionExpired` cierra la sesión.

### 10.3 Navegación condicionada por permisos

`src/navigation/RootNavigator.tsx` muestra u oculta **tabs según el rol/permiso** del usuario:

```tsx
const showUsers   = can('IDENTITY_USER_VIEW', 'READ');
const showRoles   = can('IDENTITY_ROLE_VIEW', 'READ');
const showCatalog = can('IDENTITY_ROLE_MANAGE', 'WRITE');
...
{showUsers ? <Tab.Screen name="UsersTab" .../> : null}
{showRoles ? <Tab.Screen name="RolesTab" .../> : null}
{showCatalog ? <Tab.Screen name="CatalogTab" .../> : null}
```

Además, los **botones de acción internos** también se filtran por permiso (FAB de crear usuario/rol, editar, eliminar, asignar roles, etc.), para que un `INSTRUCTOR`/`LEARNER`/`USER` jamás vea opsiones de administración.

| Pantalla | Permiso que activa su contenido |
|----------|--------------------------------|
| Tabs Usuarios | `IDENTITY_USER_VIEW` |
| Tabs Roles | `IDENTITY_ROLE_VIEW` |
| Tab Catálogo | `IDENTITY_ROLE_MANAGE` |
| Tarjetas de estadísticas en Inicio | los mismos READ/WRITE respectivos |
| Auditoría en Perfil e Inicio | `AUDIT_LOG_VIEW` |

---

## 11. Pantallas y funcionalidad del frontend

### 11.1 Flujo de autenticación

| Pantalla | Qué hace |
|----------|----------|
| `LoginScreen` | Formulario email + contraseña, valida en vivo, llama `login()` |
| `RegisterScreen` | Nombre, apellido, email, contraseña (con política visible), confirmación, tipo (`USER`/`INSTRUCTOR`/`LEARNER`) |
| `VerifyEmailScreen` | Ingresa el código de 6 dígitos recibido por correo |
| `ForgotPasswordScreen` | Pide email, invoca `forgot-password` |
| `ResetPasswordScreen` | Código + nueva contraseña + confirmación |

### 11.2 Sección autenticada

- **Inicio** (`HomeScreen`): avatar, bienvenida, tipo de actor, tarjetas de estadísticas (Usuarios, Roles, Funcionalidades, Sesiones activas) solo si hay permiso, accesos rápidos y cerrar sesión.
- **Usuarios** (`UsersScreen` / `UserDetailScreen`): listado con búsqueda y paginación; crear/editar/activar usuarios; asignar o quitar roles (solo con permisos).
- **Roles** (`RolesScreen` / `RoleCreateScreen` / `RoleDetailScreen`): listado de roles, crear roles, editar nombre/descripción, asignar/remover funcionalidades con su scope, eliminar (solo con permisos).
- **Catálogo** (`CatalogScreen`): ver/crear/eliminar módulos y funcionalidades (solo con `IDENTITY_ROLE_MANAGE`).
- **Sesiones** (`SessionsScreen`): sesiones activas del usuario; revocarlas.
- **Perfil** (`ProfileScreen`): datos de la cuenta, atajos de seguridad (Sesiones, Auditoría) y cerrar sesión.
- **Auditoría** (`AuditScreen`): bitácora de accesos con estado (éxito, contraseña incorrecta, etc.).

### 11.3 Hooks de apoyo

- `useForm` (`src/hooks/useForm.ts`): mini motor de formularios. El error aparece al salir del campo (`onBlur`) y se limpia al escribir; `validate()` recalcula todo en el envío.
- `usePaginatedList` (`src/hooks/usePaginatedList.ts`): listas paginadas con refresco y "cargar más" (`onEndReached`).

---

## 12. Validaciones: frontend vs backend vs base de datos

Existen **tres capas de validación** que se complementan (defensa en profundidad).

### 12.1 Frontend (UX inmediata)

En `src/validation/index.ts` se definen reglas reutilizables:

| Regla | Detalle |
|-------|---------|
| `ruleEmail` / `isEmail` | Regex RFC 5322, máx 254, y **descarta partes locales reservadas** (`no-reply`, `postmaster`, etc.) |
| `ruleStrongPassword` | Mín. 8, mayúscula, minúscula, número y símbolo |
| `ruleName` | Solo letras (unicode), espacios, apóstrofe y guion |
| `ruleSixDigitCode` | Exactamente 6 dígitos |
| `ruleConfirm` | Que coincida con el campo de referencia |

La política de contraseña se muestra en vivo en `RegisterScreen` (checks verdes/neutros).

### 12.2 Backend (autoridad)

- **Tags de Gin**: `binding:"required,email"`, `min=8`, `oneof=USER INSTRUCTOR LEARNER`, `len=6`, etc. → si algo falla: `400 VALIDATION_ERROR`.
- **`shared/validation`**: `ValidateEmail` (mismo criterio de correos reservados), `ValidatePassword` (misma política), `ValidateName`.
- **Errores de dominio** (`apperror`): cada caso tiene su propio código. Ejemplos: `INVALID_CREDENTIALS`, `EMAIL_TAKEN`, `WEAK_PASSWORD`, `USER_INACTIVE`, `ACCOUNT_LOCKED`, `INVALID_CODE`, `CODE_EXPIRED`, `INVALID_NAME`.

### 12.3 Base de datos (última barrera)

- `UNIQUE (email)` → un solo usuario por correo.
- `UNIQUE (name)` en roles, `UNIQUE (code)` en módulos y features.
- `CHECK` en `actor_type`, `action_level` y `scope_type` → solo valores permitidos.
- FKs e índices para integridad y rendimiento (por ejemplo, índices sobre `user_id`, `feature_id`, `module_id`).

> Mensaje para la exposición: la validación en el front mejora la experiencia, pero **la seguridad real vive en el backend y la BD**. Aunque alguien desactive la validación del cliente, la API sigue rechazando entradas inválidas y la BD mantiene la integridad.

---

## 13. Seguridad implementada

| Medida | Cómo |
|--------|------|
| Hash de contraseñas | bcrypt (costo por defecto) — nunca en texto plano |
| Sesión | JWT access (15 min) + refresh token (7 días) con **rotación** (cada refresh revoca el anterior) |
| Tokens en BD | El refresh token solo se guarda como **hash SHA-256** |
| Revocación | Logout y reset de contraseña revocan todas las sesiones |
| Bloqueo de cuenta | 5 intentos fallidos → bloqueo temporal de 15 min + revocación de sesiones |
| Anti-enumeración | `forgot-password` responde igual exista o no el usuario |
| Auditoría | Cada login (éxito o fallo) se registra con IP y User-Agent |
| Autorización | RBAC por feature + acción + scope, aplicado por middleware en cada ruta |
| CORS | Oringenes permitidos configurables |
| Panel de auditoría | Solo quienes tienen `AUDIT_LOG_VIEW` |

---

## 14. Cómo ejecutar todo el sistema

### Paso 1 — Arrancar Docker (PostgreSQL está en contenedor)

```bash
# (Windows) iniciar Docker Desktop
docker compose --env-file .env.develop up -d postgres
```

### Paso 2 — Aplicar migraciones de la BD del instructor (la primera vez)

```bash
docker compose --profile tooling run --rm liquibase-iam update
```

### Paso 3 — Levantar la API

```bash
docker compose --env-file .env.develop up -d --build api
# Health check:
curl http://localhost:8080/api/v1/health
```

### Paso 4 — Levantar el frontend

```bash
cd jeuxairel-mobile
npx expo start --port 8081    # abre la app web/emulador
```

### Paso 5 — Credenciales de prueba

| Email | Contraseña | Rol |
|-------|-----------|-----|
| `admin@sena.edu.co` | `Admin#2026` | `SYSTEM_ADMIN` |

### Paso 6 — Swagger (documentación interactiva)

`http://localhost:8080/swagger/index.html` — botón *Authorize* con el `access_token`.

---

## 15. Instrumentos para la exposición y posibles preguntas

### 15.1 Demo sugerida (3-4 min)

1. Abrir Swagger y mostrar `/auth/login` con admin para obtener tokens.
2. Mostrar `/auth/me` → usuario, roles y 45 permisos del admin.
3. En la app: iniciar sesión como admin → ver Usuarios, Roles y Catálogo.
4. Registrar un instructor → verificar correo (modo dev usa `dev_code`)→ iniciar sesión → comprobar que **no ve** las opciones de administración (separación por rol trabajada en front y back).
5. Revocar una sesión desde la pantalla de Sesiones.

### 15.2 Posibles preguntas y respuestas

**¿Cómo supiste qué construir en el backend?**
> Recibí la BD del instructor como proyecto Liquibase. Leí el esquema (`identity`, `rbac`, `rbac_catalog`, `session`, `identity_audit`), identifiqué entidades y relaciones, y mapeé cada tabla a un módulo Go. Las features del catálogo definen los permisos que protegen cada ruta.

**¿Cómo se conecta el front con el backend?**
> Con Axios apuntando a `http://localhost:8080/api/v1`. Los interceptores agregan el Bearer token y renuevan la sesión automáticamente si el access token vence. Los `endpoints.ts` envuelven cada llamada y devuelven la parte `data` de la respuesta.

**¿Dónde se valida el correo y la contraseña?**
> En las 3 capas: el front da feedback inmediato (regex y política visible), el backend valida con tags de Gin + `shared/validation` (autoridad), y la BD impone `UNIQUE`, `CHECK` y no nulos.

**¿Qué pasa si un usuario sin permiso llama a un endpoint directamente?**
> El middleware `RequirePermission` consulta la BD (feacture + roles + scope overrides) y responde `403 FORBIDDEN`. La ocultación de botones en el front es solo estética; la protección real es del backend.

**¿Cómo funciona el envío de correos?**
> El backend genera un código de 6 dígitos, guarda el hash en `session.email_verification_code` (caduca en 15 min) y lo envía por SMTP (Gmail con STARTTLS o MailHog en dev). Sin SMTP configurado, devuelve el código en la respuesta para pruebas.

**¿Cómo se manejan las sesiones?**
> Login entrega `access_token` (15 min) + `refresh_token` (7 días, guardado como hash). Cada refresh revoca el token anterior (rotación). Logout y reset de contraseña revocan todas las sesiones. Hay bloqueo tras 5 intentos fallidos.

**¿Cómo se garantiza que cada rol vea solo lo suyo?**
> En la BD, los roles tienen features con scopes (`GLOBAL`, `OWN_FICHAS`, etc.). El backend filtra por middleware y `/auth/me` expone los permisos efectivos; el front los usa para ocultar/mostrar tabs y botones. Front y backend aplican la misma jerarquía READ→WRITE→DELETE→PUBLISH→APPROVE.

**¿Por qué dos tokens y no uno?**
> El access token corto reduce la ventana de uso si se roba; el refresh permite renovar sin pedir la contraseña otra vez, y su rotación invalida tokens anteriores si se detecta reuso.

**¿Cómo se orquesta el sistema?**
> Docker Compose levanta PostgreSQL, la API (build multi-etapa de Go) y MailHog. El backend espera el healthcheck de la BD antes de conectarse (retries en el pool).

---

## 16. Glosario rápido

- **JWT**: token firmado que transporte la identidad del usuario (`user_id`, `email`, `actor_type`).
- **RBAC**: control de acceso basado en roles (rol → funcionalidades → acciones + scope).
- **Feature**: permiso atómico del catálogo (ej. `IDENTITY_ROLE_MANAGE`).
- **Scope**: alcance del permiso (`GLOBAL`, `TRAINING_CENTER`, `OWN_FICHAS`, ...).
- **Action level**: `READ < WRITE < DELETE < PUBLISH < APPROVE`.
- **Scope override**: permiso general o denegación puntual para un usuario concreto.
- **Refresh token**: credencial duradera que permite renovar la sesión.
- **bCrypt**: algoritmo de hashing de contraseñas (incluye sal automática).
- **Liquibase**: herramienta de migraciones versionadas de base de datos.
- **GORM**: ORM de Go para trabajar con la BD.
- **Gin**: framework HTTP de Go usado para la API.
- **Expo**: kit de desarrollo de React Native.

---

*Documento preparado para la exposición del desafío. Puedes usarlo como guion: capítulos 3–6 explican el backend y la BD, 7–9 los endpoints, correos y conexiones, 10–12 el front y las validaciones, 13–15 seguridad, ejecución y preguntas.*