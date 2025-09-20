package database

import "github.com/sunvc/NoLet/common"

var DB Database

// Database defines all the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	KeyExists(key string) bool
	Close() error //Close the database
}

func InitDatabase() {
	if dsn := common.LocalConfig.System.DSN; len(dsn) > 10 {
		if database, err := NewMySQL(dsn); err == nil {
			DB = database
			return
		}
	}
	DB = NewBboltdb(common.BaseDir())
}
