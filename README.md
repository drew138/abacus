# Abacus

A small personal spending app: every Nu movement by month, by category, and how each
category moves month to month. Go, SQLite and htmx, built to run as one container on a
Raspberry Pi behind Caddy, reachable only over Tailscale.

## How data arrives

n8n reads the Nu statements from Paperless, categorises them and posts each month to
`POST /api/import` with `Authorization: Bearer $IMPORT_TOKEN`. Every movement carries a
`ref` built from its statement line, so importing the same statement again changes
nothing. Movements you edit or delete in the app keep your version on later imports, and
merchant rules you set in the app apply to future imports.

## Entry points

| Binary | Does |
|---|---|
| `/rest` | serves the web app and the import API |
| `/migrate` | applies the migrations in `migration/`; run it before `/rest` on every deploy |
| `/cli set-password <username>` | creates the user or changes the password, and signs out every session |

## Configuration

| Variable | Default | |
|---|---|---|
| `IMPORT_TOKEN` | required | token n8n sends to `/api/import`, at least 32 characters |
| `DATABASE_PATH` | `/data/abacus.db` | SQLite file |
| `HTTP_ADDR` | `:8080` | listen address |
| `SECURE_COOKIE` | `true` | set `false` only for plain-HTTP local runs |
| `SESSION_TTL` | `720h` | how long a sign-in lasts |
| `MAX_BODY_BYTES` | `2097152` | largest import body |

## Layout

Ports and adapters: `internal/domain` (models, rules such as the Savings → Food → … →
Other order), `internal/application` (use cases and the ports they need) and
`internal/infrastructure` (SQLite, the htmx web interface, the JSON import API, bcrypt,
random tokens, the clock).

## Tests

```sh
go test ./...                     # unit tests
go test -tags=integration ./...   # adds the SQLite repository against the real migrations
```
