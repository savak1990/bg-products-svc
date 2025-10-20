# --- Build stage ---
FROM golang:1.25-alpine3.22 AS build

WORKDIR /src

COPY go.mod go.sum
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -trimpath -o /out/bg-products-svc ./cmd/server

# --- Runtime stage (distroless static, nonroot) ---
FROM gcr.io/distroless/static:nonroot
WORKDIR /
USER nonroot:nonroot
ENV ADDR=:8080
EXPOSE 8080
COPY --from=build /out/bg-products-svc /bg-products-svc
ENTRYPOINT ["/bg-products-svc"]