package models

import (
	"encoding/gob"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	"../configs"
	. "../log"
)

var Db = make(map[string]*mgo.Collection)
var modelsList = []string{
	"users", "films", "cinemas", "filmShows", "Imgs", "orders", "seats",
}

func gobInit() {
	// 已知bug，session中要保存自定义的结构体
	// 目前只能通过gob.Resister解决
	// see: https://github.com/astaxie/beego/issues/1110
	gob.Register(Cinema{})
	gob.Register(Film{})
	gob.Register(FilmShow{})
	gob.Register(Order{})
	gob.Register(Seat{})
	gob.Register(User{})
}

func DbInit() {
	gobInit()

	session, err := mgo.Dial(configs.MONGODB_URL)
	Log.Notice("connecting to database...")
	if err != nil {
		Log.Error("connecting to database failed...")
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	for _, model := range modelsList {
		Db[model] = session.DB(configs.MONGODB_NAME).C(model)
	}
}
