package mongo_db

import (
	"context"
	"log"

	"github.com/caarlos0/env"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type mongoConfig struct {
	Username string `env:"MONGO_INITDB_ROOT_USERNAME"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD"`
	DBName   string `env:"MONGO_INITDB_DATABASE"`
}

func Connect() (*mongo.Database, error) {
	cfg := mongoConfig{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	//url := fmt.Sprintf("mongodb://%s:%s",
	//	//cfg.Username,
	//	//cfg.Password,
	//	conf.Address,
	//	conf.Port,
	//)

	url := "mongodb://mongo_db:27017"

	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	log.Printf("connect to mongodb at %s", url)

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(cfg.DBName)

	log.Printf("connected to mongodb at db %s", cfg.DBName)

	return db, nil
}
