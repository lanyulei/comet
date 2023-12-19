package db

import (
	"database/sql"
	"github.com/lanyulei/comet/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	orm *gorm.DB
)

func Setup() {
	var (
		err   error
		sqlDB *sql.DB
	)
	sqlDB, err = sql.Open(viper.GetString("db.type"), viper.GetString("db.dsn"))
	if err != nil {
		logger.Fatalf("数据库连接失败，错误: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		logger.Fatalf("无法连接到数据库，错误: %v", err)
	}

	err = getDBDriver(sqlDB)
	if err != nil {
		logger.Fatalf("无法连接到数据库，错误: %v", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(viper.GetInt("db.maxIdleConn"))

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(viper.GetInt("db.maxOpenConn"))

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("db.connMaxLifetime")) * time.Minute)
}

func getDBDriver(sqlDB *sql.DB) (err error) {
	var (
		gormDialector gorm.Dialector
		config        *pgx.ConnConfig
	)

	if viper.GetString("db.type") == "mysql" {
		gormDialector = mysql.New(mysql.Config{
			Conn:                     sqlDB,
			DisableDatetimePrecision: true,
		})
	} else {
		config, err = pgx.ParseConfig(viper.GetString("db.dsn"))
		if err != nil {
			return
		}

		sqlDB = stdlib.OpenDB(*config)

		gormDialector = postgres.New(postgres.Config{
			Conn:                 sqlDB,
			PreferSimpleProtocol: true,
		})
	}

	orm, err = gorm.Open(gormDialector, &gorm.Config{})

	return
}

func Orm() *gorm.DB {
	return orm
}
