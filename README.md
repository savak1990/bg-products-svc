# bg-products-svc

Simple board games products service (HTTP API). Small example service used for local development and demos.

## Features
- In-memory product repository
- HTTP API (Gin)
- Liveness / readiness probes
- Docker multi-stage build (distroless runtime)

## Quick start (local)
Build and run:
```bash
go build -o bg-products-svc ./cmd/server
./bg-products-svc
