### Golang APOD API 🚀

#### 👨‍💻 Full list what has been used:
* [echo](https://github.com/labstack/echo) - Web framework.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go.
* [env](https://github.com/caarlos0/env) - Library to parse environment variables.
* [zap](https://github.com/uber-go/zap) - Logger.
* [uuid](https://github.com/google/uuid) - UUID generator.
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [minio-go](https://github.com/minio/minio-go) - AWS S3 MinIO Client SDK for Go.
* [testify](https://github.com/stretchr/testify) - Testing toolkit.
* [dockertest](https://github.com/ory/dockertest) - Tool to run docker container.
* [oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI generator.
* [golangci-lint](https://github.com/golangci/golangci-lint) - Go linter runners.
* [Docker](https://www.docker.com/) - Docker.


#### Tools installation:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest; \
    go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
    go install github.com/golang/mock/mockgen@latest

#### Docker usage:
    make up // run all containers.

#### Local development usage:
    make local-up // run all containers.
    make run // it's easier way to attach debugger or rebuild/rerun project.

#### Docker-compose files:
    docker-compose.yml - run all containers with go application.
    docker-compose.local.yml - run all containers except go application.

#### ENV files:
    .env.example - env for docker environment.
    .env.example.local - env for local environment.

### Go Server:

http://localhost:8080

### SWAGGER UI:

http://localhost:8081

### minio UI:

http://localhost:9000

username=access_key password=secret_key