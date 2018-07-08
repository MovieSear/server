package models

import (
	"gopkg.in/mgo.v2/bson"
	// "time"
)

type FilmShow struct {
	Id       bson.ObjectId `bson:"_id"`
	FilmId   bson.ObjectId `bson:"filmId"`
	CinemaId bson.ObjectId `bson:"cinemaId"`
	Time     string        `bson:"time"`
	// Time time.Time     `bson:"time"`
}
