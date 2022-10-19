package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	l "gorm.io/gorm/logger"

	"github.com/labbs/alfred/pkg/config"
	"github.com/labbs/alfred/pkg/logger"
)

type DbConnection struct {
	DB *gorm.DB
}

var (
	connection DbConnection
	err        error
)

func InitDatabase() {
	newLogger := l.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		l.Config{
			SlowThreshold:             time.Minute, // Slow SQL threshold
			LogLevel:                  l.Silent,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	switch engine := config.Database.Engine; engine {
	case "mysql":
		connection.DB, err = gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{Logger: newLogger})
	case "postgres":
		connection.DB, err = gorm.Open(postgres.Open(config.Database.DSN), &gorm.Config{Logger: newLogger})
	default:
		connection.DB, err = gorm.Open(sqlite.Open(config.Database.DSN), &gorm.Config{Logger: newLogger})
	}

	if err != nil {
		logger.Logger.Error().Err(err).Msg("")
	}
}

func GetDbConnection() DbConnection {
	return connection
}
