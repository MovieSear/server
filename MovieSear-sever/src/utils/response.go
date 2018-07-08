package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	. "../log"
)

type IResponse struct {
	Code int32
	Msg  string
	Data interface{}
}

func SuccessResponse(wptr *http.ResponseWriter, msg string, data interface{}) {
	w := *wptr
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	res := IResponse{200, msg, data}
	resData, _ := json.Marshal(res)
	fmt.Fprintf(w, "%s", resData)
}

func FailureResponse(wptr *http.ResponseWriter, msg string, data interface{}) {
	w := *wptr
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	res := IResponse{400, msg, data}
	resData, _ := json.Marshal(res)
	fmt.Fprintf(w, "%s", resData)
}

// LoadRequestBody load the body of request as []byte and
// convert it to a structure pointed to by parameter data
// parameter action is a string used for logging
// data is a pointer to a struct that store data
func LoadRequestBody(r *http.Request, action string, data interface{}) bool {
	// use ioutil to read body of http.Request
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 10240))
	if err := r.Body.Close(); err != nil {
		Log.Error(action, "falied: failed to read data from request's body, ", err)
		return false
	}
	if err := json.Unmarshal(body, data); err != nil { // if failed
		Log.Error(action, "falied: failed to decode data", string(body), "from json to structure")
		return false
	}
	return true
}
