package controllers

import (
	"net/http"

	. "../log"
	"../utils"
)

func BasicAuth(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// always return a session, even if the session is empty
		session, _ := SessionGet(w, r, "user")
		if session["user"] == nil {
			Log.Error("user doesn't login")
			utils.FailureResponse(&w, "请先登录", "")
			return
		}

		// auth successfully
		f(w, r)
	}
}
