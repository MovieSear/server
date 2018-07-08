package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Film struct {
	Id     bson.ObjectId `bson:"_id"`
	Name   string        `bson:"name"`
	Info   string        `bson:"info"`
	ImgUrl string        `bson:"imgUrl"`
	Price  int           `bson:"price"`
}
