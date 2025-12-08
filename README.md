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

---

## Kubernetes Deployment with Helm

This service includes a production-ready Helm chart for Kubernetes deployment with support for:
- Horizontal Pod Autoscaling (HPA)
- AWS Application Load Balancer (ALB) Ingress
- External DNS for automatic Route53 record creation
- Gateway API HTTPRoute support
- Configurable replicas and resources
- Health probes and security contexts

> 📚 **Publishing Your Chart**: To set up GitHub Pages as a Helm chart repository and make your chart publicly available, see the complete guide: [HELM_REPOSITORY_SETUP.md](./HELM_REPOSITORY_SETUP.md)

### Helm Chart Location

The Helm chart is located in: `deploy/helm/bg-products-svc/`

### Quick Start - Add Helm Repository

If the chart is published to a Helm repository:

```bash
# Add the Helm repository (if published)
helm repo add bg-charts https://your-helm-repo.com

# Update repository index
helm repo update

# Search for available versions
helm search repo bg-products-svc
```

### Installation Examples

#### Example 1: Simple Development Deployment

Deploy with minimal configuration - 2 static replicas, no ingress, ClusterIP service only:

```bash
# Install from local chart
helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace

# Or from repository (if published)
helm install bg-products-svc bg-charts/bg-products-svc \
  --namespace bg-products \
  --create-namespace
```

This uses default values:
- 2 replicas (static)
- HPA disabled
- ClusterIP service on port 80 → 8080
- No ingress

**Access the service:**
```bash
# Port forward to local machine
kubectl port-forward -n bg-products svc/bg-products-svc 8080:80

# Test the API
curl http://localhost:8080/api/v1/products
curl http://localhost:8080/api/v1/healthz/ready
```

#### Example 2: Production with ALB Ingress + External DNS

Deploy with AWS Application Load Balancer and automatic Route53 DNS record creation:

```bash
helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace \
  --set replicaCount=2 \
  --set autoscaling.enabled=false \
  --set ingress.enabled=true \
  --set ingress.className=alb \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/scheme"=internet-facing \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/target-type"=ip \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/healthcheck-path"=/api/v1/healthz/ready \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/listen-ports"='[{"HTTP": 80}\, {"HTTPS":443}]' \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/ssl-redirect"='443' \
  --set ingress.annotations."alb\.ingress\.kubernetes\.io/certificate-arn"=arn:aws:acm:eu-west-1:753939038916:certificate/aad20d97-4e4f-4519-8371-e53c11b07948 \
  --set ingress.annotations."external-dns\.alpha\.kubernetes\.io/alias"=true \
  --set ingress.hosts[0].host=api.your-domain.com \
  --set ingress.hosts[0].paths[0].path=/api/v1/products \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

Or use a values file for cleaner configuration:

```bash
# Create values-production.yaml
cat <<EOF > values-production.yaml
replicaCount: 2

ingress:
  enabled: true
  className: alb
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/healthcheck-path: /api/v1/healthz/ready
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:eu-west-1:753939038916:certificate/aad20d97-4e4f-4519-8371-e53c11b07948
    external-dns.alpha.kubernetes.io/alias: "true"
  hosts:
    - host: api.your-domain.com
      paths:
        - path: /api/v1/products
          pathType: Prefix

autoscaling:
  enabled: false
EOF

# Install with values file
helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace \
  -f values-production.yaml
```

**Prerequisites for this setup:**
- AWS Load Balancer Controller installed in cluster
- External DNS installed and configured for Route53
- ACM certificate created and validated
- Route53 hosted zone for your domain

**Access the service:**
```bash
# Wait for DNS propagation
dig api.your-domain.com

# Access via HTTPS
curl https://api.your-domain.com/api/v1/products
curl https://api.your-domain.com/api/v1/healthz/ready
```

#### Example 3: Gateway API with HTTPRoute

Deploy using Gateway API HTTPRoute instead of traditional Ingress:

```bash
helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace \
  --set replicaCount=2 \
  --set autoscaling.enabled=false \
  --set ingress.enabled=false \
  --set httpRoute.enabled=true \
  --set httpRoute.parentRefs[0].name=gateway \
  --set httpRoute.parentRefs[0].sectionName=http \
  --set httpRoute.hostnames[0]=api.your-domain.com \
  --set httpRoute.rules[0].matches[0].path.type=PathPrefix \
  --set httpRoute.rules[0].matches[0].path.value=/api/v1/products
```

Or with a values file:

```bash
cat <<EOF > values-gateway-api.yaml
replicaCount: 2

ingress:
  enabled: false

httpRoute:
  enabled: true
  annotations:
    external-dns.alpha.kubernetes.io/alias: "true"
  parentRefs:
    - name: gateway
      sectionName: http
      namespace: gateway-system
  hostnames:
    - api.your-domain.com
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /api/v1/products
EOF

helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace \
  -f values-gateway-api.yaml
```

**Prerequisites:**
- Gateway API CRDs installed
- Gateway controller installed (e.g., Istio, Envoy Gateway, or AWS VPC Lattice)
- Gateway resource configured

#### Example 4: Autoscaling with HPA

Deploy with Horizontal Pod Autoscaler for automatic scaling based on CPU/Memory:

```bash
helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace \
  --set autoscaling.enabled=true \
  --set autoscaling.minReplicas=2 \
  --set autoscaling.maxReplicas=10 \
  --set autoscaling.targetCPUUtilizationPercentage=70 \
  --set autoscaling.targetMemoryUtilizationPercentage=80 \
  --set resources.requests.cpu=100m \
  --set resources.requests.memory=128Mi \
  --set resources.limits.cpu=500m \
  --set resources.limits.memory=512Mi
```

Or with a values file:

```bash
cat <<EOF > values-hpa.yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi

ingress:
  enabled: true
  className: alb
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/healthcheck-path: /api/v1/healthz/ready
  hosts:
    - host: ""
      paths:
        - path: /api/v1/products
          pathType: Prefix
EOF

helm install bg-products-svc ./deploy/helm/bg-products-svc \
  --namespace bg-products \
  --create-namespace \
  -f values-hpa.yaml
```

**Verify HPA:**
```bash
# Check HPA status
kubectl get hpa -n bg-products

# Watch HPA in action
kubectl get hpa -n bg-products --watch

# Generate load to test scaling
kubectl run -it --rm load-generator --image=busybox -n bg-products -- /bin/sh -c "while true; do wget -q -O- http://bg-products-svc/api/v1/products; done"
```

### Helm Management Commands

#### List Installed Releases
```bash
# List all releases in all namespaces
helm list --all-namespaces

# List releases in specific namespace
helm list -n bg-products
```

#### Get Release Information
```bash
# Get release status
helm status bg-products-svc -n bg-products

# Get release values
helm get values bg-products-svc -n bg-products

# Get all information about release
helm get all bg-products-svc -n bg-products
```

#### Upgrade Release

```bash
# Upgrade with new values
helm upgrade bg-products-svc ./deploy/helm/bg-products-svc \
  -n bg-products \
  --set image.tag=v0.0.2 \
  --set replicaCount=3

# Upgrade with values file
helm upgrade bg-products-svc ./deploy/helm/bg-products-svc \
  -n bg-products \
  -f values-production.yaml

# Upgrade with reuse of existing values
helm upgrade bg-products-svc ./deploy/helm/bg-products-svc \
  -n bg-products \
  --reuse-values \
  --set image.tag=v0.0.2
```

#### Rollback Release

```bash
# View release history
helm history bg-products-svc -n bg-products

# Rollback to previous version
helm rollback bg-products-svc -n bg-products

# Rollback to specific revision
helm rollback bg-products-svc 2 -n bg-products
```

#### Uninstall Release

```bash
# Uninstall the release (keeps history)
helm uninstall bg-products-svc -n bg-products

# Uninstall and purge all history
helm uninstall bg-products-svc -n bg-products --wait

# Delete namespace (optional)
kubectl delete namespace bg-products
```

#### Testing Release

```bash
# Test the chart without installing
helm install bg-products-svc ./deploy/helm/bg-products-svc \
  -n bg-products \
  --dry-run --debug

# Run Helm tests (if configured)
helm test bg-products-svc -n bg-products
```

### Helm Repository Management

#### Add Repository
```bash
# Add your Helm repository
helm repo add bg-charts https://your-helm-repo.com

# Add with authentication
helm repo add bg-charts https://your-helm-repo.com \
  --username your-user \
  --password your-pass

# Verify repository was added
helm repo list
```

#### Update Repository
```bash
# Update all repositories
helm repo update

# Update specific repository
helm repo update bg-charts
```

#### Search Repository
```bash
# Search for charts in repository
helm search repo bg-products-svc

# Search with versions
helm search repo bg-products-svc --versions

# Search all repositories
helm search repo products
```

#### Remove Repository
```bash
# Remove repository after usage
helm repo remove bg-charts

# Verify removal
helm repo list
```

### Configuration Reference

Key configuration values for the Helm chart:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas (when HPA disabled) | `2` |
| `image.repository` | Docker image repository | `savak1990/bg-products-svc` |
| `image.tag` | Docker image tag | `v0.0.1` |
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `service.targetPort` | Container port | `8080` |
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.className` | Ingress class name | `alb` |
| `autoscaling.enabled` | Enable HPA | `false` |
| `autoscaling.minReplicas` | Minimum replicas for HPA | `2` |
| `autoscaling.maxReplicas` | Maximum replicas for HPA | `5` |
| `httpRoute.enabled` | Enable Gateway API HTTPRoute | `false` |

For complete configuration options, see `deploy/helm/bg-products-svc/values.yaml`.

### Debugging

```bash
# Check pod logs
kubectl logs -n bg-products -l app.kubernetes.io/name=bg-products-svc

# Follow logs
kubectl logs -n bg-products -l app.kubernetes.io/name=bg-products-svc -f

# Describe pod
kubectl describe pod -n bg-products -l app.kubernetes.io/name=bg-products-svc

# Check events
kubectl get events -n bg-products --sort-by='.lastTimestamp'

# Check ingress
kubectl describe ingress bg-products-svc -n bg-products

# Get ALB DNS name
kubectl get ingress bg-products-svc -n bg-products -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

---

## License

MIT

