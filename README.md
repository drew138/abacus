<p align="center"><img src="internal/infrastructure/api/rest/web/static/logo.svg" alt="Abacus logo" width="96" height="96"></p>

<h1 align="center">Abacus</h1>

A small personal spending app: every movement by month, by category, and how each
category moves month to month. Go, htmx and rqlite (SQLite as a server), built to run on a
Raspberry Pi behind Caddy, reachable only over Tailscale.

## How data arrives

n8n reads the Nu statements from Paperless, categorises them and posts each month to
`POST /api/import` with `Authorization: Bearer $IMPORT_TOKEN`. Every movement carries a
`ref` built from its statement line, so importing the same statement again changes
nothing. Movements you edit or delete in the app keep your version on later imports, and
merchant rules you set in the app apply to future imports.

## Running it

The image holds one binary, `/abacus`, which serves the web app and the import API. Its
database is [rqlite](https://rqlite.io), a server around SQLite, in its own container
(`rqlite/rqlite:10.6.6` with `ENABLE_FK=1`).

- **Migrations** live in `migration/` (golang-migrate format) and are applied outside the app,
  before it starts, with golang-migrate's official image:
  `migrate -path migration -database "rqlite://host:4001?x-connect-insecure=true&disableClusterDiscovery=true" up`.
- **Your account:** the first visit opens **Create account**. Once an account exists the page
  closes for good, so nobody else can register.

## Configuration

| Variable | Default | |
|---|---|---|
| `IMPORT_TOKEN` | required | token n8n sends to `/api/import`, at least 32 characters |
| `DATABASE_URL` | required | rqlite's address, e.g. `http://abacus-db:4001?disableClusterDiscovery=true` |
| `HTTP_ADDR` | `:8080` | listen address |
| `SECURE_COOKIE` | `true` | set `false` only for plain-HTTP local runs |

## Layout

Ports and adapters:

- `internal/domain`: models and rules, such as the Savings → Food → … → Other order.
- `internal/application`: the use cases. Their one dependency is the `Repository` port.
- `internal/infrastructure`: `api/rest`, every HTTP route (the web app's pages, with their
  templates and static files in `web/`, and the JSON import), and `repository/rqlite`, the
  only outbound adapter.
- `pkg/`: in-process libraries: `password` (bcrypt), `token` (session tokens) and `pesos`
  (peso formatting).

## Tests

```sh
go test ./...                     # unit tests
go test -tags=integration ./...   # adds the repository, against rqlite containers (testcontainers)
```

The integration tests start `rqlite/rqlite:10.6.6` with testcontainers, so they need Docker,
and apply the real migrations from `migration/` to each container.
