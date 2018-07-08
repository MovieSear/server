package models

import (
	"gopkg.in/mgo.v2/bson"
	// "time"
)

type Order struct {
	Id         bson.ObjectId `bson:"_id"`
	UserId     bson.ObjectId `bson:"userId"`
	FilmShowId bson.ObjectId `bson:'filmShowId'`
	SeatNum    int           `bson:"seatNum"`
	// CreateTime time.Time     `bson:"createTime"`
	CreateTime string `bson:"createTime"`
}
