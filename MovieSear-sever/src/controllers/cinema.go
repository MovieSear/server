package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	. "../log"
	. "../models"
	"../utils"
)

func CinemaGetOne(w http.ResponseWriter, r *http.Request) {
	// GET, 从URL中读取参数, 直接使用mux.Vars(r)
	vars := mux.Vars(r)
	cinemaId := vars["cinemaId"]

	cinema := Cinema{}
	err := Db["cinemas"].FindId(bson.ObjectIdHex(cinemaId)).One(&cinema)
	if err != nil {
		Log.Errorf("Get cinema id: %s failed, %v", cinemaId, err)
		utils.FailureResponse(&w, "获取电影院信息失败", "")
		return
	}

	Log.Noticef("Get cinema successfully: %s", cinema)
	utils.SuccessResponse(&w, "获取电影院信息成功", cinema)
}

func CinemaGetAll(w http.ResponseWriter, r *http.Request) {
	var cinemas []Cinema
	err := Db["cinemas"].Find(nil).All(&cinemas)
	if err != nil {
		Log.Errorf("get all cinemas failed, %v", err)
		utils.FailureResponse(&w, "获取电影院列表失败", "")
		return
	}
	Log.Notice("get all cinema successfully")
	utils.SuccessResponse(&w, "获取电影院列表成功", cinemas)
}

func CinemaAddOne(w http.ResponseWriter, r *http.Request) {
	// 1. load request's body
	newCinema := Cinema{}
	ok := utils.LoadRequestBody(r, "insert cinema", &newCinema)
	if !ok {
		utils.FailureResponse(&w, "新建电影院失败", "")
		return
	}
	// 2. verify the cinema existed or not
	existedCinema := Cinema{}
	err := Db["cinemas"].Find(bson.M{"name": newCinema.Name}).One(&existedCinema)
	if err == nil {
		Log.Errorf("insert cinema failed: cinema %s is existed", newCinema.Name)
		utils.FailureResponse(&w, "电影院已存在", "")
		return
	}
	// 3. set a new id
	newCinema.Id = bson.NewObjectId()
	// 4. insert into db
	err = Db["cinemas"].Insert(&newCinema)
	if err != nil {
		Log.Error("insert cinema falied: insert into db failed, ", err)
		utils.FailureResponse(&w, "添加电影院失败", "")
		return
	}
	// 5. success
	Log.Notice("add one cinema successfully")
	utils.SuccessResponse(&w, "添加电影院成功", "")
}

func CinemaUpdateOne(w http.ResponseWriter, r *http.Request) {
	// 1. 获得URL中的参数
	vars := mux.Vars(r)
	cinemaId := vars["cinemaId"]
	// 2. 从request中解析出body数据
	newCinema := Cinema{}
	ok := utils.LoadRequestBody(r, "update cinema", &newCinema)
	if !ok {
		utils.FailureResponse(&w, "修改电影院信息失败", "")
	}
	newCinema.Id = bson.ObjectIdHex(cinemaId)

	// 3. 修改数据
	// convert structure to bson.M, used to update
	updateData, _ := bson.Marshal(&newCinema)
	updateCinema := bson.M{}
	_ = bson.Unmarshal(updateData, &updateCinema)
	// 此处更新时如果没有"$set",会将整行直接覆盖，而不是按需修改
	err := Db["cinemas"].Update(bson.M{"_id": newCinema.Id}, bson.M{"$set": updateCinema})
	if err != nil {
		Log.Error("update cinema failed: failed to update data into db, ", err)
		utils.FailureResponse(&w, "修改电影院信息失败", "")
		return
	}

	// 修改其它表中含该电影院的项
	// ...

	// 4. 成功返回
	Log.Notice("update cinema successfully")
	utils.SuccessResponse(&w, "修改电影院成功", "")
}

func CinemaDeleteOne(w http.ResponseWriter, r *http.Request) {
	cinemaId := mux.Vars(r)["cinemaId"]

	err := Db["cinemas"].Remove(bson.M{"_id": bson.ObjectIdHex(cinemaId)})
	if err != nil {
		Log.Error("delete cinema from db failed: ", err)
		utils.FailureResponse(&w, "删除电影院失败", "")
		return
	}

	// 删除其它表中含该电影院的项
	// ...

	Log.Notice("delete cinema successfully")
	utils.SuccessResponse(&w, "删除电影院成功", "")
}

var CinemaRoutes Routes = Routes{
	Route{"CinemaGetOne", "GET", "/cinema/{cinemaId}", CinemaGetOne},
	Route{"CinemaGetAll", "GET", "/cinema/", CinemaGetAll},
	Route{"CinemaAddOne", "POST", "/cinema/", CinemaAddOne},
	Route{"CinemaUpdateOne", "PUT", "/cinema/{cinemaId}", CinemaUpdateOne},
	Route{"CinemaDeleteOne", "DELETE", "/cinema/{cinemaId}", CinemaDeleteOne},
}
