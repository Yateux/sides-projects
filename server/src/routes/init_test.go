package routes

import (
	"fmt"
	"log"
	"os"
	"testing"

	"../utils"
	mgo "gopkg.in/mgo.v2"
)

var (
	db        *mgo.Database
	dbSession *mgo.Session
	dbName    string = "TOSET"
)

func dbInit() (*mgo.Session, *mgo.Database, error) {
	dbURL := os.Getenv("TEST_MONGO_DB_URL")
	if dbURL == "" {
		dbURL = "localhost:27017/TOLSET"
	}
	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		return nil, nil, err
	}
	return dbSession, dbSession.DB(utils.DbName), nil
}

func DbClean() {
	if db == nil {
		fmt.Println("Connection to database failed")
		os.Exit(1)
	}
	// Don't forget to adapt to the collections you create in your db here
	// otherwise your tests won't work
	colls := []string{
		"tasks",
		"users",
	}
	for _, coll := range colls {
		if _, err := db.C(coll).RemoveAll(nil); err != nil {
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	var err error
	dbSession, db, err = dbInit()
	if err != nil {
		log.Fatal(err)
	}
	DbClean()
	ret := m.Run()
	DbClean()
	os.Exit(ret)
}
