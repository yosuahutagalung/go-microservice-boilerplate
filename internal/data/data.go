package data

import (
	"database/sql"
	"service_boilerplate/internal/conf"
	"service_boilerplate/internal/data/db"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/lib/pq"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewNSQPublisher, NewNSQProducer, NewOutboxRepo)

// Data .
type Data struct {
	db    *sql.DB
	query *db.Queries
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(logger)

	sqlDB, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		log.Errorf("failed opening connection to postgres: %v", err)
		return nil, nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, nil, err
	}

	queries := db.New(sqlDB)

	cleanup := func() {
		log.Info("closing the data resources")
		if err := sqlDB.Close(); err != nil {
			log.Error(err)
		}
	}

	return &Data{
		db:    sqlDB,
		query: queries,
	}, cleanup, nil
}
