package controllers

import (
	"net/http"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc 
}

type Routes []Route

/*
 *func errCheck(err Error, logMsg, resMsg string, data interface{}) bool {
 *    if err != nil {
 *        Log.Error("get all users failed")
 *        utils.FailureResponse(&w, "获取用户列表失败", "")
 *        return true
 *    }
 *    return false
 *}
 */
