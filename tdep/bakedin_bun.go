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
		if onTuneSQLDB != nil {
			onTuneSQLDB(sqlDB)
		}

		bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
		if onTuneBunDB != nil {
			onTuneBunDB(bunDB)
		}

		bunLogger, _ := zap.NewStdLogAt(o.debugLog.Named("bun"), zap.DebugLevel)
		bunDB.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(o.IsDebug()),
				bundebug.WithVerbose(o.IsDebug()),
				bundebug.WithWriter(bunLogger.Writer()),
			),
		)

		return bunDB, nil
	}

	return New(resolve, options...)
}
