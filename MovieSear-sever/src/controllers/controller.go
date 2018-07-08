package controllers

import (
	"fmt"

	"github.com/gorilla/mux"
)

type RoutesList []Routes

// NewRouter() register all routes defined in RoutesList
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, routes := range routesList {
		for _, route := range routes {
			fmt.Println("register controller: ", route.Name, route.Method, "/api"+route.Pattern)
			router.
				Methods(route.Method).
				Path("/api" + route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
		}
	}
	return router
}

// 新增加的路由，只需将其添加到此数组即可
var routesList RoutesList = RoutesList{
	UserRoutes, FilmRoutes, CinemaRoutes, FilmShowRoutes, ImgRoutes, OrderRoutes, SeatRoutes,
}
