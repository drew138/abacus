<p align="center"><img src="internal/infrastructure/api/rest/web/static/logo.svg" alt="Abacus logo" width="96" height="96"></p>

<h1 align="center">Abacus</h1>

A small personal spending app: every movement by month, by category, and how each
category moves month to month. Go, htmx and Postgres, built to run on a Raspberry Pi behind
Caddy, reachable only over Tailscale.

## How data arrives

n8n reads the Nu statements from Paperless, categorises them and posts each month to
`POST /api/import` with `Authorization: Bearer $IMPORT_TOKEN`. Every movement carries a
`ref` built from its statement line, so importing the same statement again changes
nothing. Movements you edit or delete in the app keep your version on later imports, and
merchant rules you set in the app apply to future imports.

## Running it

The image holds one binary, `/abacus`, which serves the web app and the import API. Its
database is Postgres 17, in its own container.

- **Migrations** live in `migration/` (golang-migrate format) and are applied outside the app,
  before it starts, with golang-migrate's official image:
  `migrate -path migration -database "postgres://user:password@host:5432/abacus?sslmode=disable" up`.
- **Your account:** the first visit opens **Create account**. Once an account exists the page
  closes for good, so nobody else can register.

## Releases

Pushing a tag such as `v0.1.0` runs the tests and publishes two images for amd64 and arm64:

- `ghcr.io/drew138/abacus:<tag>`: the app
- `ghcr.io/drew138/abacus-migrate:<tag>`: golang-migrate with this version's migrations; run it
  with `-database <postgres url> up` before the app

and creates a GitHub release with generated notes. Every push to `main` and every pull request
runs the unit and integration tests.

## Configuration

| Variable | Default | |
|---|---|---|
| `IMPORT_TOKEN` | required | token n8n sends to `/api/import`, at least 32 characters |
| `DATABASE_URL` | required | Postgres connection string, e.g. `postgres://abacus:…@abacus-db:5432/abacus?sslmode=disable` |
| `HTTP_ADDR` | `:8080` | listen address |
| `SECURE_COOKIE` | `true` | set `false` only for plain-HTTP local runs |

## Layout

Ports and adapters:

- `internal/domain`: models and rules, such as the Savings → Food → … → Other order.
- `internal/application`: the use cases. Their one dependency is the `Repository` port.
- `internal/infrastructure`: `api/rest`, every HTTP route (the web app's pages, with their
  templates and static files in `web/`, and the JSON import), and `repository/postgres`, the
  only outbound adapter.
- `pkg/`: in-process libraries: `password` (bcrypt), `token` (session tokens) and `pesos`
  (peso formatting).

## Tests

```sh
go test ./...                     # unit tests
go test -tags=integration ./...   # adds the repository, against Postgres containers (testcontainers)
```

The integration tests start `postgres:17-alpine` with testcontainers, so they need Docker,
and apply the real migrations from `migration/` to each container.
