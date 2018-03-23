package conn

import (
	"database/sql"
	"net"
	"net/url"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/logrusadapter"
	"github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
)

const (
	// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
	// This error code is seen sometimes when the database is starting.
	psqlCannotConnectNow = "57P03"
)

// New creates a new DB.
// It will block until a connection with the Postgres server has been established,
// or an unexpected error occurs.
func New(logger *logrus.Logger, postgresURL url.URL) (*sql.DB, error) {
	bk := backoff.NewExponentialBackOff()
	bk.MaxElapsedTime = 0 // Ensure we never stop
	// Retry dialling until connection established
	dialFunc := func(network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			KeepAlive: 5 * time.Minute,
			Timeout:   5 * time.Second,
		}
		var conn net.Conn
		connFn := func() error {
			var err error
			conn, err = dialer.Dial(network, addr)
			return err
		}
		// Retry in perpetuity
		_ = backoff.RetryNotify(
			connFn,
			bk,
			func(err error, next time.Duration) {
				logger.WithError(err).Warnf("Failed to connect to postgres, retrying in %v", next.Truncate(time.Millisecond))
			},
		)

		return conn, nil
	}

	driverConfig := stdlib.DriverConfig{
		ConnConfig: pgx.ConnConfig{
			Logger: logrusadapter.NewLogger(logger),
			Dial:   dialFunc,
		},
	}
	stdlib.RegisterDriverConfig(&driverConfig)
	db, err := sql.Open("pgx", driverConfig.ConnectionString(postgresURL.String()))
	if err != nil {
		return nil, err
	}

	// Retry DB ping separately, as it may error with a temporary error
	// even after the dialing has completed.
	err = backoff.RetryNotify(
		func() error {
			// Ping da DB.
			err := db.Ping() // nolint: vetshadow
			if err != nil {
				e, ok := err.(pgx.PgError)
				if ok && e.Code == psqlCannotConnectNow {
					// Retry connection on this one specific error
					return e
				}

				// Other errors are considered permanent
				return backoff.Permanent(err)
			}

			return nil
		},
		bk,
		func(err error, next time.Duration) {
			logger.WithError(err).Warnf("Failed to call postgres client, retrying in %v", next.Truncate(time.Millisecond))
		},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
