package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"time"
)

var Gorm *gorm.DB

func NewGormConn() *gorm.DB {
	return Gorm.Session(&gorm.Session{})
}

// InitGorm 初始化mysqlOrm
func InitGorm(username string, password string, host string, port string, dbName string, maxIdle int, maxOpen int, initFunc func(db *gorm.DB)) {
	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbName)
	var err error
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connInfo, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置, &gorm.Config{})
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields:                              true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(maxIdle)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(maxOpen)
	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	//initFunc(db)
	Gorm = db
}

func CommitOrRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}
	tx.Commit()
}

func Paginate(currentPage int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if currentPage == 0 {
			currentPage = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (currentPage - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ErrCheckDAO(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	return nil
}
