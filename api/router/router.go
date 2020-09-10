package router

import (
	"net/http"

	"github.com/bokultis/evidentor/api/middleware"

	"github.com/gorilla/mux"
)

type RoutePrefix struct {
	Prefix    string
	SubRoutes []Route
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Protected   bool
}

func NewRouter() *mux.Router {

	var AppRoutes []RoutePrefix
	AppRoutes = append(AppRoutes, rootRoutes)
	AppRoutes = append(AppRoutes, userRoutes)
	AppRoutes = append(AppRoutes, studentRoutes)

	//init router
	router := mux.NewRouter()

	for _, route := range AppRoutes {

		//create subroute
		routePrefix := router.PathPrefix(route.Prefix).Subrouter()

		//loop through each sub route
		for _, r := range route.SubRoutes {

			var handler http.Handler
			handler = r.HandlerFunc

			//check to see if route should be protected with jwt
			if r.Protected {
				handler = middleware.JWTMiddleware(r.HandlerFunc)
			}

			//attach sub route
			routePrefix.
				Path(r.Pattern).
				Handler(handler).
				Methods(r.Method).
				Name(r.Name)
		}

	}

	return router
}
