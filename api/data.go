package main

import (
	"github.com/maxwellhealth/bongo"
	"labix.org/v2/mgo"
)

type App struct {
	bongo.DocumentBase `bson:",inline"`
	FirstName          string
	LastName           string
	Gender             string
}
