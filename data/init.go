package data

import (
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
)

var Session *mgo.Session

func init() {
	var dbUrl string
	if dbUrl = os.Getenv("MONGODB_URL"); dbUrl == "" {
		dbUrl = "mongodb://localhost:27017/hazel"
	}

	log.Println("Connecting to MongoDB")

	var err error
	Session, err = mgo.Dial(dbUrl)
	if err != nil {
		log.Printf("Couln't connect to MongoDB server: %s", dbUrl)
		log.Fatalf("MongoDB dial error: %s", err.Error())
		os.Exit(1)
	}
}

func Database() *mgo.Database {
	return Session.DB("")
}

// Ensure database indexes are respected for given mongo database
func Index(db *mgo.Database) {
	if err := db.C("users").EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		panic(err)
	}
}
