package vertica_prometheus_exporter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	_ "github.com/vertica/vertica-sql-go" // register the Vertica driver
)

// OpenConnection extracts the driver name from the DSN (expected as the URI scheme), adjusts it where necessary (e.g.
// some driver supported DSN formats don't include a scheme), opens a DB handle ensuring early termination if the
// context is closed (this is actually prevented by `database/sql` implementation), sets connection limits and returns
// the handle.

// Vertica
//
// Using the https://github.com/vertica/vertica-sql-go driver, DSN format (passed through to the driver unchanged):
//
//	vertica://user:password@host:port/dbname?param=value
func OpenConnection(ctx context.Context, logContext, dsn string, maxConns, maxIdleConns int, maxConnLifetime time.Duration) (*sql.DB, error) {
	// Extract driver name from DSN.
	idx := strings.Index(dsn, "://")
	if idx == -1 {
		return nil, fmt.Errorf("missing driver in data source name. Expected format `<driver>://<dsn>`")
	}
	driver := dsn[:idx]

	// Open the DB handle in a separate goroutine so we can terminate early if the context closes.
	var (
		conn *sql.DB
		err  error
		ch   = make(chan error)
	)
	go func() {
		conn, err = sql.Open(driver, dsn)
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		if err != nil {
			return nil, err
		}
	}

	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxConns)
	conn.SetConnMaxLifetime(maxConnLifetime)

	if len(logContext) > 0 {
		logContext = fmt.Sprintf("[%s] ", logContext)

		log.Infof("%sDatabase handle successfully opened with '%s' driver", logContext, driver)
	}
	return conn, nil
}

// PingDB is a wrapper around sql.DB.PingContext() that terminates as soon as the context is closed.
//
// sql.DB does not actually pass along the context to the driver when opening a connection (which always happens if the
// database is down) and the driver uses an arbitrary timeout which may well be longer than ours. So we run the ping
// call in a goroutine and terminate immediately if the context is closed.
func PingDB(ctx context.Context, conn *sql.DB) error {
	ch := make(chan error, 1)

	go func() {
		ch <- conn.PingContext(ctx)
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}
