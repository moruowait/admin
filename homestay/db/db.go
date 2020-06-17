package db

import "github.com/jinzhu/gorm"

//orm, _ := gorm.Open("mysql", conn.GetConn().GetDB("default"))

var homeStayDB *HomeStayDB

type HomeStayDB struct {
	DB *gorm.DB
}

func SetHomestayDB(db *gorm.DB) {
	homeStayDB = &HomeStayDB{
		DB: db,
	}
}

func GetHomestayDB() *HomeStayDB {
	return homeStayDB
}
