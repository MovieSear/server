package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	. "../log"
	. "../models"
	"../utils"
)

func FilmGetOne(w http.ResponseWriter, r *http.Request) {
	// GET, 从URL中读取参数, 直接使用mux.Vars(r)
	vars := mux.Vars(r)
	filmId := vars["filmId"]

	film := Film{}
	err := Db["films"].FindId(bson.ObjectIdHex(filmId)).One(&film)
	if err != nil {
		Log.Errorf("Get film id: %s failed, %v", filmId, err)
		utils.FailureResponse(&w, "获取电影信息失败", "")
		return
	}

	Log.Noticef("Get film successfully: %s", film)
	utils.SuccessResponse(&w, "获取电影信息成功", film)
}

func FilmGetAll(w http.ResponseWriter, r *http.Request) {
	var films []Film
	err := Db["films"].Find(nil).All(&films)
	if err != nil {
		Log.Errorf("get all films failed, %v", err)
		utils.FailureResponse(&w, "获取电影列表失败", "")
		return
	}
	Log.Notice("get all film successfully")
	utils.SuccessResponse(&w, "获取电影列表成功", films)
}

func FilmAddOne(w http.ResponseWriter, r *http.Request) {
	// 1. load request's body
	newFilm := Film{}
	ok := utils.LoadRequestBody(r, "insert film", &newFilm)
	if !ok {
		utils.FailureResponse(&w, "新建电影失败", "")
		return
	}
	// 2. verify the film existed or not
	existedFilm := Film{}
	err := Db["films"].Find(bson.M{"name": newFilm.Name}).One(&existedFilm)
	if err == nil {
		Log.Errorf("insert film failed: film %s is existed", newFilm.Name)
		utils.FailureResponse(&w, "电影已存在", "")
		return
	}
	// 3. set a new id
	newFilm.Id = bson.NewObjectId()
	// 4. insert into db
	err = Db["films"].Insert(&newFilm)
	if err != nil {
		Log.Error("insert film falied: insert into db failed, ", err)
		utils.FailureResponse(&w, "添加电影失败", "")
		return
	}
	// 5. success
	Log.Notice("add one film successfully")
	utils.SuccessResponse(&w, "添加电影成功", "")
}

func FilmUpdateOne(w http.ResponseWriter, r *http.Request) {
	// 1. 获得URL中的参数
	vars := mux.Vars(r)
	filmId := vars["filmId"]
	// 2. 从request中解析出body数据
	newFilm := Film{}
	ok := utils.LoadRequestBody(r, "update film", &newFilm)
	if !ok {
		utils.FailureResponse(&w, "修改电影信息失败", "")
	}
	newFilm.Id = bson.ObjectIdHex(filmId)

	// 3. 修改数据
	// convert structure to bson.M, used to update
	updateData, _ := bson.Marshal(&newFilm)
	updateFilm := bson.M{}
	_ = bson.Unmarshal(updateData, &updateFilm)
	// 此处更新时如果没有"$set",会将整行直接覆盖，而不是按需修改
	err := Db["films"].Update(bson.M{"_id": newFilm.Id}, bson.M{"$set": updateFilm})
	if err != nil {
		Log.Error("update film failed: failed to update data into db, ", err)
		utils.FailureResponse(&w, "修改电影信息失败", "")
		return
	}

	// 修改其它表中含该电影的项
	// ...

	// 4. 成功返回
	Log.Notice("update film successfully")
	utils.SuccessResponse(&w, "修改电影成功", "")
}

func FilmDeleteOne(w http.ResponseWriter, r *http.Request) {
	filmId := mux.Vars(r)["filmId"]

	err := Db["films"].Remove(bson.M{"_id": bson.ObjectIdHex(filmId)})
	if err != nil {
		Log.Error("delete film from db failed: ", err)
		utils.FailureResponse(&w, "删除电影失败", "")
		return
	}

	// 删除其它表中含该电影的项
	// ...

	Log.Notice("delete film successfully")
	utils.SuccessResponse(&w, "删除电影成功", "")
}

var FilmRoutes Routes = Routes{
	Route{"FilmGetOne", "GET", "/film/{filmId}", FilmGetOne},
	Route{"FilmGetAll", "GET", "/film/", FilmGetAll},
	Route{"FilmAddOne", "POST", "/film/", FilmAddOne},
	Route{"FilmUpdateOne", "PUT", "/film/{filmId}", FilmUpdateOne},
	Route{"FilmDeleteOne", "DELETE", "/film/{filmId}", FilmDeleteOne},
}
