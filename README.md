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
go build -o bg-products-svc .
./bg-products-svc
```

## Docker build
To build the Docker image for multiple architectures (recommended):
```bash
docker buildx build --platform linux/amd64,linux/arm64 -t bg-products-svc:v0.0.1 -t bg-products-svc:latest .
```

To build for current architecture only:
```bash
docker build -t bg-products-svc:v0.0.1 -t bg-products-svc:latest .
```

Replace `v0.0.1` with your actual version number (e.g., `v0.0.1`, `v1.2.3`). Using `-t` multiple times tags the image with both a specific version and `latest`.

## Run the Docker image
To run the locally built image:
```bash
docker run -p 8080:8080 bg-products-svc:v0.0.1
```

Or use the latest tag:
```bash
docker run -p 8080:8080 bg-products-svc:latest
```

To run from a remote repository:
```bash
docker run -p 8080:8080 <your-registry>/bg-products-svc:v0.0.1
```

## Push to remote repository
To tag and push the image to your remote Docker registry:
```bash
# Tag the image with your registry
docker tag bg-products-svc:v0.0.1 <your-registry>/bg-products-svc:v0.0.1
docker tag bg-products-svc:latest <your-registry>/bg-products-svc:latest

# Log in to your registry (if needed)
docker login <your-registry>

# Push both tags
docker push <your-registry>/bg-products-svc:v0.0.1
docker push <your-registry>/bg-products-svc:latest
```

Replace `<your-registry>` with your Docker Hub username or private registry address (e.g., `docker.io/username`, `gcr.io/project-id`, `my-registry.azurecr.io`). Replace `v0.0.1` with your version tag.

## Quick install from Docker Hub
To fetch and run the image directly from Docker Hub without building:
```bash
docker run -p 8080:8080 savak1990/bg-products-svc:latest
```

Or use a specific version:
```bash
docker run -p 8080:8080 savak1990/bg-products-svc:v0.0.1
```

The image will be automatically pulled from Docker Hub if not already present locally.

## API Endpoints

### List products
```bash
curl http://localhost:8080/v1/products
```

### Create product
```bash
curl -X POST http://localhost:8080/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Catan", "price_cents": 3999}'
```

### Liveness probe
```bash
curl http://localhost:8080/healthz/live
```

### Readiness probe
```bash
curl http://localhost:8080/healthz/ready
```

### Health check
```bash
curl http://localhost:8080/health
```
