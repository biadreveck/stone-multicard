package main

import (
	"flag"

	"multicard/api"
	"multicard/models"
	"multicard/process"
)

var migrateDB bool

func main() {
	flag.BoolVar(&migrateDB, "migrate", false, "use to auto migrate database models when starting application")
	flag.Parse()

	models.InitializeDB("ec2-23-23-225-12.compute-1.amazonaws.com", "d2f30uhjhs6h6b", "xpfggetuhfpsew", "4835e89d5d16852397aea1d3cc36be37e37db193738ccd61a12fcd18ded78320")
	if migrateDB {
		models.AutoMigrate()
	}
	stop := process.Run()
	api.Run()
	stop <- true
}