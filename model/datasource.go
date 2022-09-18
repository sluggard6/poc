package model

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sluggard/poc/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

var cfg config.Database

// Init 初始化数据库
func Init() error {
	cfg = config.GetConfig().Database
	// log.Debug(fmt.Sprintf("%s:%s@%s", cfg.Username, cfg.Password, cfg.Url))
	// dsn := fmt.Sprintf("%s:%s@%s", cfg.Username, cfg.Password, cfg.Url)
	// _db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: logger.Default.LogMode(logger.Info)})
	_db, err := makeDb()
	// log.Debug(fmt.Sprintf("%s:%s@%s", cfg.Username, cfg.Password, cfg.Url))
	// log.Debug(cfg.Type)
	// dsn := fmt.Sprintf("%s:%s@%s", cfg.Username, cfg.Password, cfg.Url)
	// _db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: logger.Default.LogMode(logger.Info)})
	// db, err := gorm.Open(sqlite.Open("myfile.db"), &gorm.Config{})
	// db.Logger.LogMode(logger.Info)
	if err != nil {
		log.Error(err)
		return err
	}
	sqlDb, err := _db.DB()
	if err != nil {
		log.Error(err)
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(2)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDb.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(time.Hour)

	// DB = db
	db = _db
	initTable()
	return nil
}

func makeDb() (*gorm.DB, error) {
	switch cfg.Type {
	case config.Mysql:
		return makeMysqlDb()
	default:
		return makeSqlLite()
	}
}

func makeMysqlDb() (*gorm.DB, error) {
	log.Debug(fmt.Sprintf("%s:%s@%s", cfg.Username, cfg.Password, cfg.Url))
	dsn := fmt.Sprintf("%s:%s@%s", cfg.Username, cfg.Password, cfg.Url)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: logger.Default.LogMode(logger.Info)})
}

func makeSqlLite() (*gorm.DB, error) {
	// if cfg.Url ==  {
	cfg.Url = "myfile.db"
	// }
	return gorm.Open(sqlite.Open(cfg.Url), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: logger.Default.LogMode(logger.Info)})
}

func initTable() error {
	// log.Debug(reflect.TypeOf(User{}))
	db.AutoMigrate(&User{})
	// db.AutoMigrate(&Library{})
	// db.AutoMigrate(&Folder{})
	// db.AutoMigrate(&ShareLibrary{})
	// db.AutoMigrate(&File{})
	// db.AutoMigrate(&Policy{})
	return nil
}
