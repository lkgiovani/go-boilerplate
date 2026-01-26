ARG GO_VERSION=1.25.5

# Build Stage - Compila o binário Go
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /service
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app ./cmd/startServer/main.go

# Docs Build Stage - Compila a documentação OpenAPI com Bun
FROM oven/bun:1-alpine AS docs
WORKDIR /docs
COPY ./docs/package.json ./docs/bun.lock ./
RUN bun install --frozen-lockfile
COPY ./docs/ ./
RUN bun run compile

# Production Stage - Imagem final mínima com apenas o necessário
FROM gcr.io/distroless/static-debian12 AS production
ENV PROFILE=prod
WORKDIR /service
USER nonroot:nonroot

# Copia o binário compilado
COPY --from=build --chown=nonroot:nonroot /app ./app

# Copia os schemas da documentação
COPY --from=docs /docs/schema ./docs/schema

# Copia os viewers HTML
COPY --from=build --chown=nonroot:nonroot /service/resources/viewers/index.html ./resources/viewers/index.html
COPY --from=build --chown=nonroot:nonroot /service/resources/viewers/redoc.html ./resources/viewers/redoc.html
COPY --from=build --chown=nonroot:nonroot /service/resources/viewers/swagger.html ./resources/viewers/swagger.html
COPY --from=build --chown=nonroot:nonroot /service/resources/viewers/scalar.html ./resources/viewers/scalar.html

ENTRYPOINT ["/service/app"]