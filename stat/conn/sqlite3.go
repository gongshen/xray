package conn

import (
	"github.com/gongshen/xray/stat/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"os"
	"path"
)

var db *gorm.DB

func initTraffic() error {
	if err := db.AutoMigrate(&models.Traffic{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.TrafficHistory{}); err != nil {
		return err
	}
	var count int64
	if err := db.Model(&models.Traffic{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		user := &models.Traffic{
			Tag:    "admin",
			Enable: true,
		}
		return db.Create(user).Error
	}
	return nil
}

func InitDB(dbPath string) error {
	dir := path.Dir(dbPath)
	err := os.MkdirAll(dir, fs.ModeDir)
	if err != nil {
		return err
	}

	var gormLogger logger.Interface

	gormLogger = logger.Default

	c := &gorm.Config{
		Logger: gormLogger,
	}
	db, err = gorm.Open(sqlite.Open(dbPath), c)
	if err != nil {
		return err
	}

	err = initTraffic()
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db.Debug()
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
