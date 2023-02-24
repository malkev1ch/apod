package repository_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/malkev1ch/apod/pkg/postgres"
	"github.com/ory/dockertest/v3"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

const (
	pgPassword = "qwerty"
	dbName     = "postgres"

	dayDuration = time.Hour * 24
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		dockerPool, err := dockertest.NewPool("")
		if err != nil {
			log.Fatalf("failed to connect to docker: %s", err)
		}

		resource, err := dockerPool.Run(
			"postgres",
			"14.1-alpine",
			[]string{
				fmt.Sprintf("POSTGRES_PASSWORD=%v", pgPassword),
				fmt.Sprintf("POSTGRES_DB=%v", dbName),
			},
		)
		if err != nil {
			log.Fatalf("failed to start resource: %s", err)
		}

		defer func() {
			err = resource.Expire(1)
			if err != nil {
				log.Fatalf("failed to expire resource: %s", err)
			}
		}()

		defer func() {
			if err = dockerPool.Purge(resource); err != nil {
				log.Fatalf("failed to purge resource: %s", err)
			}
		}()

		postgresURL := fmt.Sprintf(
			"postgres://postgres:%v@%v/%v?sslmode=disable",
			pgPassword,
			resource.GetHostPort("5432/tcp"),
			dbName,
		)

		err = dockerPool.Retry(func() error {
			pool, err = postgres.NewPool(postgresURL)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			log.Fatalf(`failed to connect to database with URL "%s": %s`, postgresURL, err)
		}

		migr, err := migrate.New("file://../../migrations", postgresURL)
		if err != nil {
			log.Fatalf(`failed to create migrate instance: %s`, err)
		}

		err = migr.Up()
		if err != nil {
			log.Fatalf(`failed to apply migrations: %s`, err)
		}
		os.Exit(m.Run())
	}
}

func clearDatabase(ctx context.Context) error {
	_, err := pool.Exec(ctx, "TRUNCATE pictures")
	return err
}
