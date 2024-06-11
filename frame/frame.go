package frame

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	mydb "my-frame-go/db"
)

// MyGorm gorm
var MyGorm *gorm.DB

// MySqlx sqlx
var MySqlx *sqlx.DB

// MyLog zap
var MyLog *zap.Logger

// InitGorm 初始化gorm
func InitGorm(username string, password string, host string, port string, dbName string, maxIdle int, maxOpen int, initFunc func(db *gorm.DB)) {
	MyGorm = mydb.InitGorm(username, password, host, port, dbName, maxIdle, maxOpen, initFunc)
}

// InitSqlx 初始化sqlx
func InitSqlx(username string, password string, host string, port string, dbName string, maxIdle int, maxOpen int) {
	MySqlx = mydb.InitSqlx(username, password, host, port, dbName, maxIdle, maxOpen)
}
