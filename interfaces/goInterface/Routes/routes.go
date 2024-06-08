package Routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"goInterface/Routes/LocalHost"
)

/*
Author: Weston Simon
Email: weston@wcloud.com

Creation Date: 2024-04-28 11:44:47

© wcloud
*/

func InitRouters(router *mux.Router) {
	//router.Use(middleware.LogRequest)

	fmt.Println("Main Router Init Complete")

	localRouter := router.PathPrefix("/api/local").Subrouter()

	LocalHost.InitLocalRouts(localRouter)

}
