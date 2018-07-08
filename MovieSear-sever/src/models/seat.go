package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Seat struct {
	Id         bson.ObjectId `bson:"_id"`
	FilmShowId bson.ObjectId `bson:'filmShowId'`
	SeatNum    int           `bson:"seatNum"`
}
