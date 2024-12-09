package constants

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type DatabaseCfg struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	Database string `env:"DB_NAME" envDefault:"postgres"`
}

var dbCfg *DatabaseCfg

// Database return database config and load it from dbCfg and using singleton DP
func Database() *DatabaseCfg {
	if dbCfg == nil {

		dbCfg = &DatabaseCfg{}
		if err := env.Parse(dbCfg); err != nil {
			panic(err)
		}
	}

	return dbCfg
}

func (d DatabaseCfg) Url() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username, d.Password, d.Host, d.Port, d.Database)
}
