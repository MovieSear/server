package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Cinema struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Info string        `bson:"info"`
}
