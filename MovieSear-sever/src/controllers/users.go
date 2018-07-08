package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	. "../log"
	. "../models"
	"../utils"
)

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	// GET, 从URL中读取参数, 直接使用mux.Vars(r)
	vars := mux.Vars(r)
	userId := vars["userId"]

	user := User{}
	err := Db["users"].FindId(bson.ObjectIdHex(userId)).One(&user)
	if err != nil {
		Log.Errorf("Get user id: %s failed, %v", userId, err)
		utils.FailureResponse(&w, "获取用户信息失败", "")
		return
	}

	Log.Noticef("Get user successfully: %s", user)
	utils.SuccessResponse(&w, "获取用户信息成功", user)
}

func UserGetAll(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := Db["users"].Find(nil).All(&users)
	if err != nil {
		Log.Errorf("get all users failed, %v", err)
		utils.FailureResponse(&w, "获取用户列表失败", "")
		return
	}
	Log.Notice("get all user successfully")
	utils.SuccessResponse(&w, "获取用户列表成功", users)
}

// post data: {"name":"c","password":"c"}, 首字母大写小写皆可
// password需要在客户端先经过bcrypt加密
func UserAddOne(w http.ResponseWriter, r *http.Request) {
	// 1. load request's body
	newUser := User{}
	ok := utils.LoadRequestBody(r, "insert user", &newUser)
	if !ok {
		utils.FailureResponse(&w, "新建用户失败", "")
		return
	}
	// 2. verify the user existed or not
	existedUser := User{}
	err := Db["users"].Find(bson.M{"name": newUser.Name}).One(&existedUser)
	if err == nil {
		Log.Errorf("insert user failed: user %s is existed", newUser.Name)
		utils.FailureResponse(&w, "用户已存在", "")
		return
	}
	// 3. set a new id
	// 由于将User.Id解释成_id(见user定义), 所以user.Id需要自己指定, 没有这一步会导致插入失败
	newUser.Id = bson.NewObjectId()
	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashedPassword)
	// 4. insert into db
	err = Db["users"].Insert(&newUser)
	if err != nil {
		Log.Error("insert user falied: insert into db failed, ", err)
		utils.FailureResponse(&w, "添加用户失败", "")
		return
	}
	// 5. success
	Log.Notice("add one user successfully")
	utils.SuccessResponse(&w, "添加用户成功", "")
}

func UserUpdateOne(w http.ResponseWriter, r *http.Request) {
	// 1. 获得URL中的参数
	vars := mux.Vars(r)
	userId := vars["userId"]
	// 2. 从request中解析出body数据
	newUser := User{}
	ok := utils.LoadRequestBody(r, "update user", &newUser)
	if !ok {
		utils.FailureResponse(&w, "修改用户信息失败", "")
	}
	newUser.Id = bson.ObjectIdHex(userId)
	// 3. 通过session验证是否有修改权限
	session, _ := SessionGet(w, r, "user")
	sessionUser, _ := session["user"].(User)
	if sessionUser.Id != newUser.Id {
		Log.Error("no privilege to change user detail")
		utils.FailureResponse(&w, "没有权限修改用户信息", "")
		return
	}
	// 4. 修改数据
	// convert structure to bson.M, used to update
	updateData, _ := bson.Marshal(&newUser)
	updateUser := bson.M{}
	_ = bson.Unmarshal(updateData, &updateUser)
	// 此处更新时如果没有"$set",会将整行直接覆盖，而不是按需修改
	err := Db["users"].Update(bson.M{"_id": newUser.Id}, bson.M{"$set": updateUser})
	if err != nil {
		Log.Error("update user falied: failed to update data into db, ", err)
		utils.FailureResponse(&w, "修改用户信息失败", "")
		return
	}
	// 5. 成功返回
	Log.Notice("update user successfully")
	utils.SuccessResponse(&w, "修改用户成功", "")
}

// POST data: {"oldPassword": "", "newPassword": ""}
// oldPassword为明文，newPassword需要在客户端经过bcrypt加密后在传到服务端
func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	// 1. load data from request.body
	passwords := make(map[string]string)
	ok := utils.LoadRequestBody(r, "update password", &passwords)
	Log.Noticef("passwords: %v", passwords)
	if !ok {
		utils.FailureResponse(&w, "修改密码失败", "")
	}
	// 2. get session
	session, _ := SessionGet(w, r, "user")
	sessionUser, _ := session["user"].(User)
	// 3. verify the old password is correct or not
	err := bcrypt.CompareHashAndPassword([]byte(passwords["oldPassword"]), []byte(sessionUser.Password))
	if err != nil {
		Log.Error("update password failed: old password is incorrect")
		utils.FailureResponse(&w, "原密码错误", "")
		return
	}
	// 4. update data in db
	err = Db["users"].Update(bson.M{"_id": sessionUser.Id},
		bson.M{"$set": bson.M{"password": passwords["newPassword"]}})
	if err != nil {
		Log.Error("update password failed: write db failed, ", err)
		utils.FailureResponse(&w, "修改密码失败", "")
		return
	}
	// 5. success
	Log.Notice("update password successfully")
	utils.SuccessResponse(&w, "修改密码成功", "")
}

// 不能提供删除用户的功能，因此此函数只作为示例，其他如电影片则提供删除功能
func UserDeleteOne(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]

	err := Db["users"].Remove(bson.M{"_id": bson.ObjectIdHex(userId)})
	if err != nil {
		Log.Error("delete user from db failed: ", err)
		utils.FailureResponse(&w, "删除用户失败", "")
		return
	}

	Log.Notice("delete user successfully")
	utils.SuccessResponse(&w, "删除用户成功", "")
}

// UserLogin should bind a session on request
func UserLogin(w http.ResponseWriter, r *http.Request) {
	// 1. 解析body数据
	user := User{}
	ok := utils.LoadRequestBody(r, "user login", &user)
	if !ok {
		utils.FailureResponse(&w, "登录失败", "")
	}
	// 2. 从db中取出用户信息
	existedUser := User{}
	err := Db["users"].Find(bson.M{"name": user.Name}).One(&existedUser)
	if err != nil {
		Log.Error("user login failed: user not found, ", err)
		utils.FailureResponse(&w, "登录失败,用户不存在", "")
		return
	}
	// 3. 验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password))
	if err != nil {
		Log.Error("user login failed: password is incorrect, ")
		utils.FailureResponse(&w, "登录失败,密码错误", "")
		return
	}
	// 4. user login successfully, save user into session
	session, _ := SessionGet(w, r, "user")
	session["user"] = existedUser
	_ = SessionSet(w, r, "user", session)
	// 5. response successfully
	Log.Noticef("user %v login successfully", existedUser.Name)
	utils.SuccessResponse(&w, "登录成功", "")
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	UserAddOne(w, r)
}

// UserLogout should delete the session bind on request
func UserLogout(w http.ResponseWriter, r *http.Request) {
	_ = SessionDel(w, r, "user")
	utils.SuccessResponse(&w, "登出成功", "")
}

// 后四个函数涉及到session，需要浏览器进行测试
var UserRoutes Routes = Routes{
	Route{"UserGetOne", "GET", "/user/{userId}", BasicAuth(UserGetOne)},
	Route{"UserGetAll", "GET", "/user/", UserGetAll},
	Route{"UserAddOne", "POST", "/user/", UserAddOne},
	Route{"UserUpdateOne", "PUT", "/user/{userId}", BasicAuth(UserUpdateOne)},
	Route{"UserUpdatePassword", "POST", "/user/updatePassword", BasicAuth(UserUpdatePassword)},
	// Route { "UserDeleteOne",      "DELETE", "/user/{userId}",       UserDeleteOne, },
	Route{"UserLogin", "POST", "/user/login", UserLogin},
	Route{"UserRegister", "POST", "/user/register", UserRegister},
	Route{"UserLogout", "GET", "/user/logout", UserLogout},
}
