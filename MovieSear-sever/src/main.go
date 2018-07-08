package main

import (
	"fmt"
	"net/http"

	"./configs"
	"./controllers"
	"./models"
)

func init() {
	models.DbInit()
}

func main() {
	router := controllers.NewRouter()
	http.ListenAndServe(configs.HOST+":"+configs.PORT, router)
	fmt.Println("start server...")
}
