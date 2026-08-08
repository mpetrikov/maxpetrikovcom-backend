# maxpetrikovcom-backend

Backend service for maxpetrikov.com, an educational platform focused on
programming, DevOps, and interactive hands-on labs.

## Migrations

This project uses `golang-migrate` for PostgreSQL schema migrations.

Install the migration CLI:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Check the installation:

```bash
migrate -version
```

Apply migrations:

```bash
migrate \
  -path migrations \
  -database "postgres://max:123@localhost:5432/maxpetrikov?sslmode=disable" \
  up
```

Downgrade one migration:

```bash
migrate \
  -path migrations \
  -database "postgres://max:123@localhost:5432/maxpetrikov?sslmode=disable" \
  down 1
```

Apply all pending migrations:

```bash
make migrate-up
```

Rollback the latest migration:

```bash
make migrate-down
```

Check the current migration version:

```bash
make migrate-version
```

Force a migration version:

```bash
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
