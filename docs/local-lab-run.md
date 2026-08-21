# Local lab run

This runbook verifies the local lab flow end to end:

```text
POST /lab-sessions
  -> RabbitMQ lab.create
  -> worker
  -> Kubernetes Pod
  -> lab_sessions.status = running
```

## Prerequisites

Install and verify:

```powershell
docker version
kubectl version --client
kind version
migrate -version
```

Open the Bruno collection:

```text
devtools/bruno/backend
```

Set these Bruno variables:

```text
baseUrl=http://localhost:8080
labSlug=linux-processes
refreshToken=<refresh-token>
```

`POST /lab-sessions` is protected. For local testing, open Google login in the
browser:

```text
http://localhost:8080/auth/google
```

After the callback returns tokens, copy the refresh token into the Bruno
`refreshToken` variable. Then run:

```text
auth / Refresh Token
```

The Bruno request stores the returned access token in the `accessToken`
variable. After that, protected API requests in the collection can run.

## Start local infrastructure

Create the kind cluster and apply local Kubernetes resources:

```powershell
make kind-bootstrap
```

Start PostgreSQL and RabbitMQ:

```powershell
docker compose -f deploy/local/docker-compose.yml up -d
```

Apply database migrations:

```powershell
make migrate-up
```

Seed the local demo lab:

```powershell
make seed-demo-lab
```

## Run backend processes

Use Kubernetes runner mode in `.env`:

```env
LAB_RUNNER_TYPE=kubernetes
KUBECONFIG_PATH=C:\Users\<user-name>\.kube\config
KUBERNETES_LAB_NAMESPACE=maxpetrikov-labs
KUBERNETES_POD_READY_TIMEOUT=60s
```

Run the worker in one terminal:

```powershell
go run ./cmd/worker
```

Run the API in another terminal:

```powershell
go run ./cmd/api
```

Verify that the demo lab is visible with Bruno:

```text
labs / Get Lab by Slug
```

## Start a lab session with Bruno

Run:

```text
lab_session / Start Lab Session
```

The request sends:

```json
{
  "lab_slug": "linux-processes"
}
```

The Bruno request stores the response `id` in the `labSessionId` variable.
The initial API response may still show `pending` because the worker processes
the RabbitMQ job asynchronously.

Poll the session state with:

```text
lab_session / Get Lab Session
```

or list recent sessions with:

```text
lab_session / List My Lab Sessions (1)
```

## Verify Kubernetes

List lab pods:

```powershell
kubectl get pods -n maxpetrikov-labs
```

Expected result:

```text
NAME             READY   STATUS
lab-<session-id> 1/1     Running
```

Inspect the pod:

```powershell
kubectl describe pod lab-<session-id> -n maxpetrikov-labs
```

The pod should include:

- image from `labs.image`
- CPU and memory requests
- CPU and memory limits
- `activeDeadlineSeconds` from `labs.timeout_minutes`
- labels and annotations with `lab-session-id`, `lab-id`, and `user-id`

## Lab runtime command

For the current PoC, the Kubernetes runner overrides the container command:

```text
/bin/sh -c "sleep infinity"
```

This keeps generic Linux images alive while the lab session is running. The current contract requires lab images to include `/bin/sh`.

Before introducing golden images, decide whether the final command ownership belongs to the image entrypoint or to explicit lab template `command`/`args` fields.

## Verify database state

Check the session in PostgreSQL:

```powershell
docker exec maxpetrikov-postgres `
  psql -U max -d maxpetrikov `
  -c "select id, status, namespace, pod_name, started_at, expires_at, failure_reason from lab_sessions order by created_at desc limit 5;"
```

Expected result:

- `status` eventually becomes `running`
- `namespace` is `maxpetrikov-labs`
- `pod_name` is populated
- `started_at` is populated
- `failure_reason` is empty

The expected status transition is:

```text
pending -> provisioning -> running
```

The Bruno `Get Lab Session` response should include populated `namespace`,
`pod_name`, and `started_at` fields after the worker finishes provisioning.

## Useful cleanup

Delete lab pods manually during local testing:

```powershell
kubectl delete pod -n maxpetrikov-labs -l app.kubernetes.io/name=maxpetrikov-lab
```

Delete the kind cluster:

```powershell
make kind-down
```
