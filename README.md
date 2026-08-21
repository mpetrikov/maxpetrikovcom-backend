# maxpetrikovcom-backend

Backend service for maxpetrikov.com, an educational platform focused on
programming, DevOps, and interactive hands-on labs.

## Migrations

This project uses `golang-migrate` for PostgreSQL schema migrations.

Install the migration CLI if it is not already available:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Common commands:

```bash
make migrate-up
make migrate-down
make migrate-version
make migrate-force VERSION=2
```

`migrate-force` does **not** execute migration SQL and does not create or remove tables.

It only changes the migration version stored by `golang-migrate`.

Use it to recover from a broken or `dirty` migration state after manually fixing the database schema.

Example:

```text
migration 2 fails
→ database becomes dirty
→ fix the database manually
→ make migrate-force VERSION=1
→ make migrate-up
```

Use `migrate-force` only when you understand the current database state.

## Local lab run

See [docs/local-lab-run.md](docs/local-lab-run.md) for the reproducible local
flow with Docker Compose, kind, migrations, demo lab seed, Bruno API checks, and
Kubernetes/DB verification.
