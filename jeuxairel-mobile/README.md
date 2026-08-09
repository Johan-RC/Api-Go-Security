# JeuxAirel Mobile — Frontend React Native (IdenSecurity)

Aplicación móvil (React Native + TypeScript + Expo) que sirve de **cliente** para la API Go
de identidad y accesos (**JeuxAirel IAM**). Desde esta app puedes: iniciar sesión, gestionar
usuarios, roles, funcionalidades, sesiones activas, y ver la auditoría de accesos.

> **Resumen en una frase:** la app **nunca toca la base de datos**. Se comunica por HTTP/JSON
> con la API Go; la API valida identidad (`JWT`) y permisos (`RBAC`) y es la única que habla
> con PostgreSQL.

---

## 1. ¿Cómo se ven las piezas?

```
┌──────────────────────────┐     HTTP + JSON      ┌──────────────────────────┐      SQL      ┌──────────────────────┐
│   App React Native       │  ──────────────────►  │       API GO (Gin)      │  ───────────► │     PostgreSQL       │
│   (este proyecto)        │                        │ • /api/v1/...            │               │  schemas:            │
│                          │  ◄──────────────────  │ • JWT + RBAC            │  ◄───────────  │   identity, rbac,    │
│   - pxg pantallas        │      respuesta JSON    │ • validación servidor   │               │   session, audit...  │
│   - axios cliente        │                        └──────────────────────────┘               └──────────────────────┘
└──────────────────────────┘
```

- **Frontend (esta app):** muestra pantallas, captura datos, valida formularios y hace peticiones.
- **Backend (API Go):** expone endpoints REST en `http://localhost:8080/api/v1`, valida el token, revisa permisos y ejecuta consultas SQL con GORM.
- **Base de datos (PostgreSQL):** vive dentro de Docker (`docker-compose`). Su esquema se crea con **Liquibase**.

La app solo conoce la **URL base de la API**. Todo lo demás pasa por HTTP.

---

## 2. Cómo levantar todo (flujo completo)

### 2.1 Base de datos (Docker)

```bash
cd docker-infra
cp .env.develop .env        # opcional, ya hay valores por defecto
docker compose up postgres -d
```

Esto crea la base `design-software-develop`. Luego se aplican los cambios del esquema IAM:

```bash
docker compose --profile tooling run --rm liquibase-iam update
```

### 2.2 API Go (backend)

```bash
cd API-Go
cp .env.example .env        # edita JWT_SECRET con un valor fuerte
go run .
```

Verifica que responde:

```bash
curl http://localhost:8080/api/v1/health   # → {"status":"ok"}
```

Swagger de la API (docs interactivas): `http://localhost:8080/swagger/index.html`

### 2.3 App React Native (este proyecto)

```bash
cd jeuxairel-mobile
npm install

# URL de la API según dónde corras la app:
#  - Emulador Android: 10.0.2.2
#  - Simulador iOS/Web: localhost
#  - Teléfono físico : IP de tu PC (ej. 192.168.1.40)
copy .env.example .env      # luego edita EXPO_PUBLIC_API_URL según lo anterior

npm start                   # abre Expo GO / emulador, o usa npm run android / ios
```

> En un **teléfono real**, backend y celular deben estar en la misma red y la API debe
> escuchar en `0.0.0.0` (es el comportamiento por defecto del `.env.example`).

---

## 3. Autenticación: cómo se conecta la app con el backend

Este es el corazón de la comunicación. El backend emite **dos tokens**:

| Token        | Qué es            | Duración por defecto | Cómo viaja             |
|--------------|-------------------|----------------------|------------------------|
| `access_token`| JWT (identidad + permisos) | 15 min      | Header `Authorization: Bearer <token>` |
| `refresh_token`| Sesión de larga duración (guardado SOLO en el servidor en hash) | 7 días | Solo se envía en `auth/refresh` y `auth/logout` |

### 3.1 Inicio de sesión (`POST /auth/login`)
La app envía `{ email, password }`. El backend valida credenciales (bcrypt), registra en
auditoría y responde:

```json
{ "success": true, "data": { "access_token": "...", "refresh_token": "..." } }
```

La app guarda ambos tokens en **AsyncStorage** (almacenamiento local del dispositivo) y
decodifica el JWT para conocer `user_id`, `email` y `actor_type` (el JWT **no se valida** en
la app, solo se lee el payload para mostrar el nombre).

### 3.2 Cada petición autenticada
El archivo `src/api/client.ts` configura **axios** con dos interceptores **globales**:

1. **Interceptor de petición:** antes de salir, lee el access token y agrega
   `Authorization: Bearer <token>` automáticamente. Ninguna pantalla lo hace a mano.
2. **Interceptor de respuesta:** si el servidor responde `401` (token expirado), la app:
   - Intenta renovar con `POST /api/auth/refresh` enviando el refresh token (rotación).
     - → si funciona, guarda el **nuevo par** y **reintenta** la petición original (la que había fallado).
     - → si falla (token revocado), borra todo y manda al usuario a **Login**.

Este flujo evita la pantalla de "sesión expirada" y mantiene al usuario entrando sin fricción.
> El principio `onSessionExpired` (en `AuthContext`) avisa cuándo ya no se puede renovar.

---

## 4. Estructura del proyecto (arquitectura simple)

```
jeuxairel-mobile/
├─ App.tsx                    # Proveedores + navegación raíz
├─ .env.example               # URL de la API (EXPO_PUBLIC_API_URL)
└─ src/
   ├─ api/                    # ★ Así se habla con el backend
   │   ├─ client.ts           #   axios + interceptores (token + refresh)
   │   ├─ endpoints.ts        #   1 función por endpoint (login, listUsers, etc.)
   │   ├─ errors.ts           #   mensajes amigables ante errores HTTP
   │   └─ tokenStore.ts       #   guarda/lee/borra tokens en AsyncStorage
   ├─ types/                  # interfaces TS que reflejan los DTO del backend
   ├─ validation/             # validación de formularios (email, mínimo 8, etc.)
   ├─ contexts/AuthContext.tsx# estado global de la sesión (login/logout)
   ├─ theme/                  # colores, tipografías, espacios, sombras
   ├─ components/             # UI reutilizable (Button, Input, Select, Card...)
   ├─ navigation/             # navigators (auth stack + pestañas)
   ├─ hooks/                  # usePaginatedList (lista + paginación)
   └─ screens/                # pantallas agrupadas por módulo
```

**Regla de oro:** las pantallas **nunca** llaman `axios` directamente. Todo pasa por
`src/api/endpoints.ts`, que además desempaqueta la respuesta del servidor (deja solo `data`).

---

## 5. Envoltorio de respuesta (por qué el front nunca habla con la DB)

Todas las respuestas del backend llevan este "sobre":

```json
{ "success": true, "message": "...", "data": {...} }        // éxito
{ "success": false, "message": "...", "error": { "code": "..." } }  // error
```

El front usa ese `error.code` para mostrar el **mensaje correcto en español** (`errors.ts`),
por ejemplo `INVALID_CREDENTIALS → "Correo o contraseña incorrectos"`. Así la app entiende
el problema sin depender del texto crudo del servidor. Cuando el servidor responde `403`,
significa que el usuario **no tiene el permiso RBAC** necesario → se muestra «No tienes
permisos…».

---

## 6. Tabla de uso de la API desde la app (endpoints implementados)

| Módulo de la app | Endpoint (método)                     | Qué hace la pantalla                    |
|------------------|---------------------------------------|-----------------------------------------|
| Login            | `POST /auth/login`                    | Envía credenciales y recib tokens      |
| Alta usuario     | `POST /users` (protegido)             | Crea el usuario (registro público)     |
| Recuperar acceso | `POST /auth/forgot-password`          | Inicia reseteo de contraseña            |
| Nuevo password   | `POST /auth/reset-password`           | Cambia contraseña con token             |
| Inicio           | `GET /users`, `/roles`, `/features`, `/sessions/active`| Cuenta y pide tarjetas de stats |
| Usuarios         | `GET /users`, `PUT/DELETE /users/:id` | Lista, activa/desactiva, elimina        |
| Detalle usuario  | `GET /users/:id`, `GET /users/:id/roles`, `POST .../roles`, `DELETE .../roles/:userRoleId` | Roles y scope del usuario |
| Roles            | `GET /roles`, `POST /roles`, `GET/PUT/DELETE /roles/:id`, `GET/PUT /roles/:id/features` | CRUD de roles y sus funcionalidades |
| Catálogo         | `GET/POST /modules`, `GET/POST /features`, `DELETE /features/:id` | Módulos y funcionalidades |
| Sesiones         | `GET/DELETE /sessions/active`         | Ver y revocar dispositivos activos      |
| Auditoría        | `GET /audit/logins`                   | Registro de accesos (éxitos/fallos)     |

> Permisos: algunos endpoints piden rol con feature según el código del backend (por ej.
> `IDENTITY_USER_VIEW`); si el usuario no lo tiene, el API responde `403` y la app lo avisa.

---

## 7. Cuántas capas de validación hay? (Frontend + Backend)

1. **Validación del front (React Native)** — para usuario amable y ahorrar llamadas:
   - Correo con formato correcto (`isEmail`).
   - Campos obligatorios (`isRequired`).
   - Contraseña con mínimo 8 caracteres (`minLength`).
   - Confirmación de contraseña que coincida.
   - Tipos permitidos para `actor_type`, `scope_type`, `action_level` (Selects cerrados).
   - La pantalla bloquea el botón mientras envía (`loading`) y muestra errores bajo cada campo.
2. **Validación del backend (Go)** — la fuente de verdad, se repite por seguridad:
   - Gin `binding` en cada DTO (`required`, `email`, `min=8`, `oneof=...`).
   - Errores tipificados con códigos (`VALIDATION_ERROR`, `INVALID_CREDENTIALS`, etc.).
   - Verificación real de contraseña con **bcrypt**; nunca se guarda texto plano.
3. **Base de datos** — restricciones finales (CHECK, UNIQUE, índices) definidas con Liquibase.

> Regla de experiencia: el front valida para **agilidad**, el backend valida para **seguridad**.

---

## 8. Buenas prácticas que ya vienen aplicadas

- **Seguridad:** la app guarda refresh token fuera del alcance del render (AsyncStorage) y el
  nombre nunca se loguea; el JWT se maneja solo en la capa `api/`.
- **No tocar la DB:** la app hace solo HTTP; si alguien "modifica" el front, la API y la DB
  siguen protegidas por tokens y permisos.
- **Centro de errores único:** toda petición pasa por el interceptor y `getErrorMessage`,
  así el usuario ve los errores en español consistentemente.
- **Tipado fuerte:** cada respuesta de la API está tipada en `src/types/` (evita errores).
- **Componentes reutilizables:** `Screen`, `Card`, `Input`, `Select`, `Button`, `Badge`, `Toast`.

---

## 9. Mantenimiento fácil

- **Quiero tocar un flujo:** busca en la pantalla dentro de `src/screens/<modulo>/`.
- **Quiero llamar un endpoint nuevo:** 1) agrega la función en `src/api/endpoints.ts`,
  2) define su tipo de respuesta en `src/types/models.ts`, 3) úsala desde la pantalla.
- **Quiero cambiar colores/tipografía:** centralizados en `src/theme/index.ts`.
- **TypeScript:** corre `npm run typecheck` para comprobar sin errores.

## Documentación relacionada

- Swagger de la API: `http://localhost:8080/swagger/index.html`
- Esquemas DB y changelogs: carpeta `design-software-iam-db-main` (Liquibase).
- Composes para infraestructure: carpeta `docker-infra`.

> Proyecto pensado para desarrollador junior: componentes pequeños, nombres en español,
> un único lugar por responsabilidad y toda la comunicación vía el archivo `endpoints.ts`.