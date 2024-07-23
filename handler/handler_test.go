package handler 

import (
	"product-info/lib/testutils"
	"product-info/database"

	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M){

	testutils.LoadTestEnv()

	//appName := os.Getenv("APP_NAME")
        //appEnv  := os.Getenv("APP_ENV")
        dbName  := os.Getenv("DBNAME")

        //fmt.Println("DB Name : ", dbName)
        //fmt.Println("APP Name : ", appName)
        //fmt.Println("APP ENV : ", appEnv)

        if err := database.DropAllTables(dbName); err != nil {
                log.Errorln(err)
                return
        }

        if err := database.StartMigration(dbName); err != nil {
                log.Errorln(err)
                return
        }

        code := m.Run()
        os.Exit(code)
}
