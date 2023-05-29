package mysqlx

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	dialer := mysql.New(mysql.Config{
		DSN: dsn,
	})
	//dialer := mysql.Open(dsn)

	gormDB, err := gorm.Open(dialer, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}
	var rawDB *sql.DB
	rawDB, err = gormDB.DB()
	if err != nil {
		return nil, err
	}
	if err = rawDB.Ping(); err != nil {
		return nil, err
	}
	//rawDB.SetMaxIdleConns(100)
	//rawDB.SetMaxOpenConns(1000000)
	//rawDB.SetConnMaxLifetime(60 * time.Second)
	return gormDB, nil
}
