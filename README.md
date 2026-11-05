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

## Running it

The image holds one binary, `/abacus`, which serves the web app and the import API.

- **Migrations** live in `migration/` (golang-migrate format) and are applied outside the app,
  before it starts: `migrate -path migration -database sqlite:///data/abacus.db up`, with a
  `migrate` built with `-tags sqlite`. The homelab deploy does this in its own container.
- **Your account:** the first visit opens **Create account**. Once an account exists the page
  closes for good, so nobody else can register.

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
