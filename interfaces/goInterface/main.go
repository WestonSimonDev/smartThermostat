package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"goInterface/Connector"
	"goInterface/Routes"
	"net/http"
)

func main() {
	mainRouter := mux.NewRouter()
	Routes.InitRouters(mainRouter)

	Connector.InitDB()

	err := http.ListenAndServe(":8888", mainRouter)
	if err != nil {
		fmt.Println("Bind Error: ", err)
	}
	fmt.Println("hi")

}
