package mysql

import (
	"bluebell/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func TestInit(dbCfg *settings.MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DbName,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect test DB failed, err:", zap.Error(err))
		return err
	}

	db.SetMaxOpenConns(dbCfg.MaxOpenConns)
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)

	return nil
}

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		settings.MysqlSettings.User,
		settings.MysqlSettings.Password,
		settings.MysqlSettings.Host,
		settings.MysqlSettings.Port,
		settings.MysqlSettings.DbName,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err:", zap.Error(err))
		return err
	}

	db.SetMaxOpenConns(settings.MysqlSettings.MaxOpenConns)
	db.SetMaxIdleConns(settings.MysqlSettings.MaxIdleConns)
	zap.L().Info("init mysql success")
	return
}

func Close() {
	_ = db.Close()
}
