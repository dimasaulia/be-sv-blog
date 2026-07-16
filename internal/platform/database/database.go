package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sv-blog/internal/platform/config"
	"github.com/sv-blog/internal/platform/logger"
)

type Database struct {
	DB  *sqlx.DB
	log *logger.LayerLogger
}

func New(ctx context.Context, cfg config.Config, appLogger *logger.Logger) (*Database, error) {
	log := appLogger.Layer("platform.database")
	end := log.Start(ctx, "New")

	dsn, err := buildDSN(cfg.Database.URL)
	if err != nil {
		end(err)
		return nil, err
	}

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		end(err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(int(cfg.Database.MaxOpenConns))
	sqlDB.SetMaxIdleConns(int(cfg.Database.MinConns))
	sqlDB.SetConnMaxLifetime(cfg.Database.MaxConnLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.MaxConnIdleTime)

	if err := sqlDB.PingContext(ctx); err != nil {
		end(err)
		return nil, err
	}

	end(nil)
	return &Database{
		DB:  sqlx.NewDb(sqlDB, "mysql"),
		log: log,
	}, nil
}

// buildDSN converts a mysql://user:pass@host:port/db?params URL into the
// user:pass@tcp(host:port)/db?params DSN format expected by go-sql-driver/mysql.
func buildDSN(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse database url: %w", err)
	}

	var userInfo string
	if u.User != nil {
		password, _ := u.User.Password()
		userInfo = u.User.Username() + ":" + password
	}

	query := u.Query()

	// Translate JDBC-style params (useSSL, allowPublicKeyRetrieval) to the
	// go-sql-driver/mysql equivalents. The driver requests the server's RSA
	// public key automatically for caching_sha2_password over a plaintext
	// connection, so allowPublicKeyRetrieval has no driver-level equivalent
	// and is simply dropped.
	if useSSL := query.Get("useSSL"); useSSL != "" {
		query.Del("useSSL")
		if _, ok := query["tls"]; !ok {
			query.Set("tls", useSSL)
		}
	}
	query.Del("allowPublicKeyRetrieval")

	if _, ok := query["parseTime"]; !ok {
		query.Set("parseTime", "true")
	}
	if _, ok := query["clientFoundRows"]; !ok {
		query.Set("clientFoundRows", "true")
	}

	dbName := strings.TrimPrefix(u.Path, "/")

	return fmt.Sprintf("%s@tcp(%s)/%s?%s", userInfo, u.Host, dbName, query.Encode()), nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}

func (d *Database) Close() error {
	return d.DB.Close()
}
