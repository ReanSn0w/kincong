# Сборка приложения
FROM golang:1.25-alpine AS application
ARG GITHUB_REF=bundle
ADD . /bundle
WORKDIR /bundle

RUN apk --no-cache add ca-certificates

RUN \
    echo "Building service documentation." && \
    GOPRIVATE=https://git.papkovda.ru go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g ./internal/server/api.go --parseVendor -ot yaml -pd

RUN \
    version=${GITHUB_REF} && \
    echo "Building service. Version: ${version}" && \
    go build -ldflags "-X main.build=${version}" -o /srv/app ./cmd/server/main.go

# Финальная сборка образа
FROM scratch

COPY --from=application /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=application /srv /srv
COPY --from=application /bundle/static /srv/static
COPY --from=application /bundle/docs/swagger.yaml /srv/static/swagger.yaml

ENV PORT=8080
ENV BASE_URL=http://localhost:8080

EXPOSE 8080
WORKDIR /srv
ENTRYPOINT ["/srv/app"]
