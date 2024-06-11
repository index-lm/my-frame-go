package mydb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

// 初始化数据库（sqlx）
func InitSqlx(username string, password string, host string, port string, dbName string, maxIdle int, maxOpen int) *sqlx.DB {
	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbName)
	var err error
	fmt.Println("数据库链接地址：", connInfo)
	fmt.Println()
	db, err := sqlx.Open("mysql", connInfo)
	if err != nil {
		panic(err.Error())
	}
	// 将最大并发空闲链接数设置为 5.
	// 小于或等于 0 表示不保留任何空闲链接.
	db.SetMaxIdleConns(maxIdle)
	//设置同时打开的连接数(使用中+空闲)
	//设为5。将此值设置为小于或等于0表示没有限制
	//最大限制(这也是默认设置)。
	db.SetMaxOpenConns(maxOpen)
	// 设置最大生存时间为1小时
	// 设置为0，表示没有最大生存期，并且连接会被重用
	// forever (这是默认配置).
	db.SetConnMaxLifetime(time.Hour)
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}
