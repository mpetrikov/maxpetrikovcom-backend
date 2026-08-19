# Local kind

## Install local kind

Install: `go install sigs.k8s.io/kind@v0.32.0`

Check: `kind --version`

## Local kind cluster

This directory contains the local Kubernetes runtime for hands-on labs.

The local setup uses one kind cluster:

- cluster: `maxpetrikov-labs`
- system namespace: `maxpetrikov-system`
- lab runtime namespace: `maxpetrikov-labs`
- worker service account: `maxpetrikov-system/maxpetrikov-worker`

Create the cluster:

```bash
kind create cluster --config deploy/local/kind/cluster.yaml
```

Apply namespaces and RBAC:

```bash
kubectl apply -f deploy/local/kind/namespace.yaml
kubectl apply -f deploy/local/kind/rbac.yaml
```

Check the worker permissions:

```bash
kubectl auth can-i create pods \
  --as system:serviceaccount:maxpetrikov-system:maxpetrikov-worker \
  --namespace maxpetrikov-labs

kubectl auth can-i watch pods \
  --as system:serviceaccount:maxpetrikov-system:maxpetrikov-worker \
  --namespace maxpetrikov-labs

kubectl auth can-i delete pods \
  --as system:serviceaccount:maxpetrikov-system:maxpetrikov-worker \
  --namespace maxpetrikov-labs
```

Delete the cluster:

```bash
kind delete cluster --name maxpetrikov-labs
```

The first Kubernetes runner iteration should create lab pods in the
`maxpetrikov-labs` namespace. If lab sessions later move to one namespace per
session, replace the namespaced Role with a ClusterRole that can create/delete
namespaces and bind the worker service account to it.
