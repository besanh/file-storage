package data

import (
	"context"
	"database/sql"
	"file/internal/conf"
	db "file/internal/data/db/generated"
	"file/internal/util"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewFileRepo, NewAuthRepo, util.NewPublicPEM, NewTransactionManager, NewPhysicalFileRepo)

// Data .
type Data struct {
	DB    *sql.DB
	Query *db.Queries
}

func (d *Data) Queries(ctx context.Context) *db.Queries {
	if tx, ok := FromContext(ctx); ok {
		return d.Query.WithTx(tx)
	}
	return d.Query
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	dbConn, err := sql.Open(
		c.Database.Driver,
		c.Database.Source,
	)
	if err != nil {
		return nil, nil, err
	}

	dbConn.SetMaxOpenConns(20)
	dbConn.SetMaxIdleConns(10)
	dbConn.SetConnMaxLifetime(time.Hour)

	if err := dbConn.Ping(); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.Info("closing the data resources")
		_ = dbConn.Close()
	}

	return &Data{
		DB:    dbConn,
		Query: db.New(dbConn),
	}, cleanup, nil
}
