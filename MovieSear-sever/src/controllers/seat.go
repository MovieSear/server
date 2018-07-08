package controllers

import (
	"net/http"

	. "../log"
	. "../models"
	"../utils"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func SeatAddOne(newOrder Order) bool {
	newSeat := Seat{}
	newSeat.FilmShowId = newOrder.FilmShowId
	newSeat.SeatNum = newOrder.SeatNum

	newSeat.Id = bson.NewObjectId()

	err := Db["seats"].Insert(&newSeat)
	if err != nil {
		Log.Error("insert order falied: insert into db failed, ", err)
		return false
	}

	Log.Notice("add one seat successfully")
	return true
}

func SeatDeleteOne(deleteOrder Order) bool {
	deleteSeat := Seat{}
	deleteSeat.FilmShowId = deleteOrder.FilmShowId
	deleteSeat.SeatNum = deleteOrder.SeatNum

	err := Db["seats"].Remove(bson.M{"filmshowid": deleteSeat.FilmShowId, "seatNum": deleteSeat.SeatNum})
	if err != nil {
		Log.Error("delete seat from db failed: ", err)
		return false
	}

	Log.Notice("delete seat successfully")
	return true
}

func SeatGetFromFilmShowId(w http.ResponseWriter, r *http.Request) {
	var seatNums []int

	vars := mux.Vars(r)
	filmShowId := vars["filmShowId"]

	err := Db["seats"].Find(bson.M{"filmshowid": bson.ObjectIdHex(filmShowId)}).Distinct("seatNum", &seatNums)
	if err != nil {
		Log.Errorf("get seatNums failed, %v", err)
		utils.FailureResponse(&w, "获取座位号列表失败", "")
		return
	}
	Log.Notice("get seatNums successfully")
	utils.SuccessResponse(&w, "获取座位号列表成功", seatNums)
}

func SeatGetAll(w http.ResponseWriter, r *http.Request) {
	var seats []Seat
	err := Db["seats"].Find(nil).All(&seats)
	if err != nil {
		Log.Errorf("get all seats failed, %v", err)
		utils.FailureResponse(&w, "获取座位列表失败", "")
		return
	}
	Log.Notice("get all user successfully")
	utils.SuccessResponse(&w, "获取座位列表成功", seats)
}

var SeatRoutes Routes = Routes{
	Route{"SeatGetFromFilmShowId", "GET", "/seat/{filmShowId}", SeatGetFromFilmShowId},
	Route{"SeatGetAll", "GET", "/seat", SeatGetAll},
}
