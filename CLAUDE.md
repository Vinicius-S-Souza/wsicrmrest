# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

WSICRMREST is a REST API service written in Go, converted from WinDev procedures. It provides JWT-based authentication and connects to Oracle databases via TNSNAMES. The service is designed to maintain compatibility with the original WinDev implementation while leveraging Go's performance and concurrency features.

## Build and Run Commands

### Linux/WSL

```bash
# Initial setup (creates dbinit.ini from example, creates directories, downloads dependencies)
make setup

# Build the binary for Linux
make build

# Build for Windows (32 and 64 bits)
make build-windows

# Build for Windows 32 bits only
make build-windows-32

# Build for Windows 64 bits only
make build-windows-64

# Run the server (requires dbinit.ini and Oracle database with ORGANIZADOR table)
make run

# Build and run in one step
make dev

# Format code
make fmt

# Run static analysis
make vet

# Format and verify
make check

# Clean build artifacts and logs
make clean

# Test APIs (requires server running)
make test-api
```

### Windows

```batch
REM Build the application (interactive menu)
scripts\build_windows.bat

REM Run the server (auto-detects 32/64 bit)
scripts\run_windows.bat
```

All compiled binaries are placed in the `build/` directory.

## Critical Architecture Patterns

### 1. Configuration Loading Priority

The system has a **specific loading priority** for configuration values that must be understood:

- **JWT Credentials (hardcoded)**: Values in `internal/config/config.go` are fixed and cannot be configured via `dbinit.ini`:
  - `SecretKey: "CloudI0812IcrMmDB"` (gsKey from WinDev)
  - `Issuer: "WSCloudICrmIntellsys"` (gsIss from WinDev)
  - `KeyDelivery: "Ped2505IcrM"` (gsKeyDelivery from WinDev)
  - `Timezone: 0` (gnFusoHorario from WinDev)

- **Organization Data (database-only)**: All organization configuration comes **exclusively** from the `ORGANIZADOR` table via `database.LeOrganizador()` during startup. The system **will not start** if this table doesn't exist or is empty. Do not attempt to configure organization data in `dbinit.ini`.

- **Application Settings (configurable)**: `dbinit.ini` only controls:
  - Database connection (tns_name, username, password) - **Oracle only**
  - Application settings (port, environment)
  - Logging preferences (ws_grava_log_db, ws_detalhe_log_api)
  - CORS configuration (allowed_origins, allowed_methods, allowed_headers)

- **Version Information (hardcoded)**: Application version is defined as global variables in `internal/config/config.go`:
  - `Version = "Ver 1.26.4.27"` (gsVersao from WinDev)
  - `VersionDate = "2025-10-16T11:55:00"` (gsDataVersao from WinDev)

### 2. Startup Sequence

The main.go initialization follows this **critical order**:

1. Load `dbinit.ini` configuration (including CORS settings)
2. Initialize Zap logger (creates `log/wsicrmrest_YYYY-MM-DD.log`)
3. Connect to Oracle database
4. **Load ORGANIZADOR table** (FATAL if fails - system exits)
5. Log CORS configuration
6. Setup Gin router with CORS middleware
7. Setup routes
8. Start HTTP server

Breaking this sequence or making ORGANIZADOR loading non-fatal will cause production issues.

### 3. Request Context Pattern

Every API request creates a `RequestContext` object that tracks:
- UUID (unique request identifier)
- Start time (for duration calculation)
- Client ID and application name (populated during token validation)
- Log details (if `ws_detalhe_log_api` enabled)

This context is passed through handlers and used for database logging via `database.GravaLogDB()`.

### 4. Scope System (Bitwise)

The `WSAPLSCOPO` field uses bitwise flags (13 defined scopes):
- Bit 1 (1): clientes
- Bit 2 (2): lojas
- Bit 3 (4): ofertas
- Bit 4 (8): produtos
- Bit 5 (16): pontos
- Bit 6 (32): private
- Bit 7 (64): convenio
- Bit 8 (128): giftcard
- Bit 9 (256): cobranca
- Bit 10 (512): basico
- Bit 11 (1024): sistema
- Bit 12 (2048): terceiros
- Bit 13 (4096): totem

Use `utils.Escopo(code)` to convert bitwise codes to space-separated scope strings.

## WinDev to Go Conversion Mapping

When converting WinDev procedures, follow these established patterns:

| WinDev Function | Go Function | Location |
|----------------|-------------|----------|
| `pgGerarToken()` | `handlers.GenerateToken()` | `internal/handlers/token.go` |
| `pgWSRestTeste()` | `handlers.WSTest()` | `internal/handlers/wstest.go` |
| `pgLeOrganizador()` | `database.LeOrganizador()` | `internal/database/organizador.go` |
| `pgGravaLogDB()` | `database.GravaLogDB()` | `internal/database/log.go` |
| `pgScopo()` | `utils.Escopo()` | `internal/utils/helpers.go` |
| `fgEliminaCaracterNulo()` | `utils.EliminaCaracterNulo()` | `internal/utils/helpers.go` |
| `pgStringChange()` | `utils.StringChange()` | `internal/utils/helpers.go` |
| `FcDateTime()` | `utils.FormatDateTimeOracle()` | `internal/utils/helpers.go` |
| `pgCalcTimeStampUnix()` | `utils.CalcTimeStampUnix()` | `internal/utils/helpers.go` |

### Key Conversion Notes:

- WinDev global variables (gsKey, gsIss, etc.) → hardcoded in `config.LoadConfig()`
- WinDev `SQLExec()`/`SQLFetch()` → Go `db.Query()`/`db.QueryRow()`
- WinDev `RESULT` → Go `return`
- WinDev date/time formatting uses MM/DD/YYYY for Oracle compatibility

## Database Schema Requirements

The system requires these Oracle tables (see `docs/DATABASE_SCHEMA.md` for full schema):

1. **ORGANIZADOR** (MANDATORY - system won't start without it)
   - Must have at least one record with `ORGCODIGO > 0`
   - Contains all organization configuration
   - Loaded once during startup

2. **WSAPLICACOES** (application registry)
   - Stores client_id, client_secret, scopes, JWT expiration
   - Queried for every token generation request

3. **WSAPLLOGTOKEN** (token audit log)
   - Written to on successful token generation

4. **WSREQUISICOES** (request audit log)
   - Written to asynchronously via goroutine for every request
   - Can be disabled with `ws_grava_log_db = false`

## Adding New API Endpoints

Follow this pattern (see existing handlers in `internal/handlers/`):

1. Create handler function in `internal/handlers/`:
   ```go
   func NewHandler(cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) gin.HandlerFunc {
       return func(c *gin.Context) {
           reqCtx := reqcontext.NewRequestContext()
           // ... handler logic ...
           go db.GravaLogDB(reqCtx.UUID, method, endpoint, ...) // Log to DB
       }
   }
   ```

2. Add model in `internal/models/models.go` if needed

3. Register route in `internal/routes/routes.go`:
   ```go
   apiGroup.GET("/path", handlers.NewHandler(cfg, db, logger))
   ```

4. All handlers must:
   - Create a RequestContext at the start
   - Calculate duration with `reqCtx.GetDuration()`
   - Log to database via `db.GravaLogDB()` (asynchronously with `go`)
   - Return JSON responses matching WinDev response format

## Logging Architecture

The system has two logging mechanisms:

1. **File Logging** (Zap):
   - Automatic daily rotation: `log/wsicrmrest_YYYY-MM-DD.log`
   - Structured JSON format
   - Used for application events and errors

2. **Database Logging**:
   - Every request logged to `WSREQUISICOES` table
   - Runs asynchronously (goroutine) to not block responses
   - Includes: UUID, timestamps, duration, headers (Authorization removed), parameters, response
   - Controlled by `ws_grava_log_db` setting in dbinit.ini

## JWT Token Generation

Tokens use HS256 (HMAC-SHA256) with the hardcoded `SecretKey`. The generation process:

1. Validate Basic Auth header (base64-encoded client_id:client_secret)
2. Query `WSAPLICACOES` table for application
3. Verify application status and credentials
4. Generate JWT with:
   - Header: `{"typ":"JWT","alg":"HS256"}`
   - Payload: `{"iss", "nbf", "exp", "client_id", "scope", "aplicacao"}`
   - Signature: HMAC-SHA256 with hardcoded key
5. Log token to `WSAPLLOGTOKEN` table
6. Return token response

## Environment-Specific Behavior

- `environment = production` in dbinit.ini sets Gin to ReleaseMode (less verbose logging)
- `environment = development` keeps Gin in debug mode
- Log directory is always created automatically if missing

## CORS Configuration

The system implements custom CORS (Cross-Origin Resource Sharing) middleware to allow controlled access from web applications.

### Configuration Options

In `dbinit.ini`:

```ini
[CORS]
AllowedOrigins=
AllowedMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
AllowedHeaders=Origin,Content-Type,Content-Length,Accept-Encoding,Authorization,Grant_type,X-CSRF-Token
AllowCredentials=true
MaxAge=43200
```

### Two Operating Modes:

1. **Development Mode** (empty `AllowedOrigins`):
   - Allows all origins (`Access-Control-Allow-Origin: *`)
   - Suitable for local development
   - Logged as: "CORS configurado para permitir TODAS as origens (*)"

2. **Production Mode** (specific origins):
   - Only allows listed origins
   - Example: `AllowedOrigins=https://app.example.com,https://admin.example.com`
   - Logged as: "CORS configurado com origens restritas"

### Key Features:

- Custom middleware implementation (no external CORS package)
- Automatic OPTIONS (preflight) request handling
- Configurable via `dbinit.ini`
- Supports `Grant_type` header for WinDev compatibility
- Logs configuration at startup

### Implementation:

- **Middleware:** `internal/middleware/cors.go`
- **Config:** `internal/config/config.go` (CORSConfig struct)
- **Applied in:** `cmd/server/main.go` before routes

See `docs/setup/CONFIGURACAO-CORS.md` for complete documentation.

## Common Troubleshooting

**"Organizador Não Cadastrado" error on startup:**
- The ORGANIZADOR table is missing or has no records with ORGCODIGO > 0
- This is fatal - system will not start
- Insert a record into ORGANIZADOR table

**Oracle connection errors:**
- Verify `LD_LIBRARY_PATH` points to Oracle client libraries
- Ensure `tnsnames.ora` has the configured TNS name
- Test with `sqlplus username/password@tns_name`

**JWT credential changes not applying:**
- JWT credentials are hardcoded in `internal/config/config.go`
- Changes to `dbinit.ini` [jwt] section are ignored (section doesn't exist)
- Must recompile after changing hardcoded values

**CORS errors in browser:**
- Check that origin is in `AllowedOrigins` list
- If using `AllowCredentials=true`, cannot use wildcard origins
- OPTIONS requests should return 204 status
- See `docs/setup/CONFIGURACAO-CORS.md` for troubleshooting

## Documentation References

- `QUICKSTART.md`: 5-minute setup guide
- `docs/DATABASE_SCHEMA.md`: Complete Oracle schema and maintenance queries
- `docs/GLOBAL_VARIABLES.md`: WinDev variable mapping and configuration details
- `docs/USAGE_EXAMPLES.md`: Code examples for common patterns
- `docs/setup/CONFIGURACAO-CORS.md`: Complete CORS configuration guide
- `docs/setup/BUILD_WINDOWS.md`: Windows compilation guide (32 and 64 bits)
- `docs/VERSIONING.md`: Version management and update guide
- `CHANGELOG.md`: Implementation history and conversion notes
