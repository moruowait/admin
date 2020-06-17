package conn

import (
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/jinzhu/gorm"
	homestaydb "github.com/moruowait/admin/homestay/db"
)

var globalConn db.Connection

func SetConn(conn db.Connection) {
	globalConn = conn
	hsdb, err := gorm.Open("mysql", conn.GetDB("default"))
	if err != nil {
		panic(err)
	}
	homestaydb.SetHomestayDB(hsdb)
}

func GetConn() db.Connection {
	return globalConn
}
