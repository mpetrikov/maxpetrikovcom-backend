DATABASE_URL ?= postgres://max:123@localhost:5432/maxpetrikov?sslmode=disable

.PHONY: migrate-up migrate-down migrate-version migrate-force

migrate-up:
	migrate \
		-path migrations \
		-database "$(DATABASE_URL)" \
		up

migrate-down:
	migrate \
		-path migrations \
		-database "$(DATABASE_URL)" \
		down 1

migrate-version:
	migrate \
		-path migrations \
		-database "$(DATABASE_URL)" \
		version

migrate-force:
	@test -n "$(VERSION)" || (echo "VERSION is required"; exit 1)
	migrate \
		-path migrations \
		-database "$(DATABASE_URL)" \
		force $(VERSION)
