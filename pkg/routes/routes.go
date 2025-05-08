package routes

import (
	"net/http"

	"github.com/ladecadence/EkiAPI/pkg/config"
	"github.com/ladecadence/EkiAPI/pkg/controllers"
	"github.com/ladecadence/EkiAPI/pkg/database"
)

func RegisterRoutes(db database.Database, config config.Config, router *http.ServeMux) {
	router.HandleFunc("GET /api", controllers.ConfMiddleWare(db, config, controllers.ApiRoot))

	// users
	router.HandleFunc("POST /api/signup", controllers.ConfMiddleWare(db, config, controllers.ApiSignup))
	router.HandleFunc("GET /api/login", controllers.ConfMiddleWare(db, config, controllers.ApiLogin))
	router.HandleFunc("GET /api/logout", controllers.ConfMiddleWare(db, config, controllers.ApiLogout))

	// missions
	router.HandleFunc("GET /api/missions", controllers.ConfMiddleWare(db, config, controllers.ApiGetMissions))
	router.HandleFunc("GET /api/mission/{name}", controllers.ConfMiddleWare(db, config, controllers.ApiGetMission))
	router.HandleFunc("POST /api/newmission", controllers.ConfMiddleWare(db, config, controllers.ApiNewMission))

	// data
	router.HandleFunc("POST /api/newdata", controllers.ConfMiddleWare(db, config, controllers.ApiNewDatapoint))
	router.HandleFunc("GET /api/data/{mission}", controllers.ConfMiddleWare(db, config, controllers.ApiGetDataMission))
	router.HandleFunc("GET /api/lastdata/{mission}", controllers.ConfMiddleWare(db, config, controllers.ApiGetLastDataMission))

	// images
	router.HandleFunc("GET /api/images/{mission}", controllers.ConfMiddleWare(db, config, controllers.ApiGetImageListMission))
	router.HandleFunc("GET /api/lastimage/{mission}", controllers.ConfMiddleWare(db, config, controllers.ApiGetLastImageMission))
	router.HandleFunc("POST /api/imgupload", controllers.ConfMiddleWare(db, config, controllers.ApiUploadImage))
	router.HandleFunc("GET /api/imgdownload/{name}", controllers.ConfMiddleWare(db, config, controllers.ApiDownloadImage))

	// updates (SSE)
	router.HandleFunc("GET /api/events", controllers.ConfMiddleWare(db, config, controllers.ApiEvents))
}
