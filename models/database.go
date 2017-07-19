package models

import (
	"fmt"

   "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

var database *gorm.DB

// func init() {
// 	database, _ = gorm.Open("postgres", "host=ec2-23-23-225-12.compute-1.amazonaws.com user=xpfggetuhfpsew dbname=d2f30uhjhs6h6b password=4835e89d5d16852397aea1d3cc36be37e37db193738ccd61a12fcd18ded78320")   
// 	// var err error
//  //   	database, err = gorm.Open("postgres", "host=ec2-23-23-225-12.compute-1.amazonaws.com user=xpfggetuhfpsew dbname=d2f30uhjhs6h6b password=4835e89d5d16852397aea1d3cc36be37e37db193738ccd61a12fcd18ded78320")   
//   //  	if err != nil {
// 		// panic("Unable to connect database")
//   //  	}   
// }

func InitializeDB(host, dbName, user, password string) {
	var err error
   	database, err = gorm.Open("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s", host, dbName, user, password))
   	if err != nil {
		panic("Unable to connect database")
   	}
	database.LogMode(true)
}

func AutoMigrate() {
	DropTables()

	database.AutoMigrate(&User{}, &Wallet{}, &Card{})

	database.Model(&Wallet{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	database.Model(&Card{}).AddForeignKey("wallet_id", "wallets(id)", "RESTRICT", "RESTRICT")

	database.Model(&User{}).AddUniqueIndex("idx_user_login", "login")
	database.Model(&Wallet{}).AddUniqueIndex("idx_wallet_name_user", "name", "user_id")
}

func DropTables() {
	database.Model(&User{}).RemoveIndex("idx_user_login")
	database.Model(&Wallet{}).RemoveIndex("idx_wallet_name_user")

	database.DropTableIfExists(&Card{}, &Wallet{}, &User{})
}
