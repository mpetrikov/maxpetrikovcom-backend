DATABASE_URL ?= postgres://max:123@localhost:5432/maxpetrikov?sslmode=disable
KIND_CLUSTER ?= maxpetrikov-labs
KIND_CONFIG ?= deploy/local/kind/cluster.yaml

.PHONY: migrate-up migrate-down migrate-version migrate-force
.PHONY: seed-demo-lab
.PHONY: kind-up kind-apply kind-down kind-status kind-check-rbac kind-bootstrap

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

seed-demo-lab:
	docker exec -i maxpetrikov-postgres \
		psql -U max -d maxpetrikov \
		< devtools/sql/seed_demo_lab.sql

kind-up:
	kind create cluster --config $(KIND_CONFIG)

kind-apply:
	kubectl apply -f deploy/local/kind/namespace.yaml
	kubectl apply -f deploy/local/kind/rbac.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl cluster-info --context kind-$(KIND_CLUSTER)
	kubectl get ns

kind-check-rbac:
	kubectl auth can-i create pods --as system:serviceaccount:maxpetrikov-system:maxpetrikov-worker --namespace maxpetrikov-labs
	kubectl auth can-i watch pods --as system:serviceaccount:maxpetrikov-system:maxpetrikov-worker --namespace maxpetrikov-labs
	kubectl auth can-i delete pods --as system:serviceaccount:maxpetrikov-system:maxpetrikov-worker --namespace maxpetrikov-labs

kind-bootstrap: kind-up kind-apply kind-check-rbac
