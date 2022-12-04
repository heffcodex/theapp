package dep

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type BunConfigPostgres struct {
	DSN            string `mapstructure:"dsn"`
	MaxConnections int    `mapstructure:"max_connections"`
}

func NewBunPostgres(
	cfg BunConfigPostgres,
	onTuneConnector func(conn *pgdriver.Connector),
	onTuneSQLDB func(db *sql.DB),
	onTuneBunDB func(db *bun.DB),
) *D[*bun.DB] {
	resolve := func() (*bun.DB, error) {
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

		return bunDB, nil
	}

	return NewDep(true, resolve)
}
