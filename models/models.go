package models

import (
	"fmt"
	"log"

	"github.com/go-gin-example/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/soft_delete"
)

var (
	db *gorm.DB
)

type Model struct {
	ID         int                   `gorm:"primary_key" json:"id"`
	CreatedOn  int64                 `gorm:"autoCreateTime" json:"created_on"`
	ModifiedOn int64                 `gorm:"autoUpdateTime" json:"modified_on"`
	DeletedOn  soft_delete.DeletedAt `json:"deleted_on"`
}

func Setup() {
	var (
		err                                       error
		dbName, user, password, host, tablePrefix string
	)

	dbName = setting.DatabaseSetting.Name
	user = setting.DatabaseSetting.User
	password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	tablePrefix = setting.DatabaseSetting.TablePrefix

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		log.Println(err)
	}

	db.Set("MaxIdleConns", 10)
	db.Set("MaxOpenConns", 100)
}
