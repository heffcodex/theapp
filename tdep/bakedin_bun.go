package tdep

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"

	"github.com/heffcodex/theapp/tcfg"
)

func NewBunPostgres(
	cfg tcfg.BunPostgres,
	onTuneConnector func(conn *pgdriver.Connector),
	onTuneSQLDB func(db *sql.DB),
	onTuneBunDB func(db *bun.DB),
	options ...Option,
) *D[*bun.DB] {
	resolve := func(o OptSet) (*bun.DB, error) {
		conn := pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN))
		if onTuneConnector != nil {
			onTuneConnector(conn)
		}

		sqlDB := sql.OpenDB(conn)
		sqlDB.SetMaxOpenConns(cfg.MaxConnections)
		sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTimeSeconds())
		if onTuneSQLDB != nil {
			onTuneSQLDB(sqlDB)
		}

		bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
		if onTuneBunDB != nil {
			onTuneBunDB(bunDB)
		}

		logLevel := zap.ErrorLevel
		if o.IsDebug() {
			logLevel = zap.DebugLevel
		}

		log := o.Log().Named("bun")
		stdLog, _ := zap.NewStdLogAt(log, logLevel)

		bunDB.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithVerbose(o.IsDebug()),
				bundebug.WithWriter(stdLog.Writer()),
			),
		)

		return bunDB, nil
	}

	return New(resolve, options...)
}
