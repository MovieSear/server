package models

import (
	"gopkg.in/mgo.v2/bson"
)

// find from mongodb(users):
// { "_id" : ObjectId("5ae700e02137a87fe6fb1637"), "name" : "lzl", "password" : "new-pass"}
// 说明后面的描述(`bson:"_id"`)表示了储存在数据库内的key
type User struct {
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	Password string `bson:"password"`
}

