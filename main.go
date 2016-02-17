package main

import (
	. "app/common"
	"app/controllers"
	"app/services"
	"github.com/fvbock/endless"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	services.InitDb()
	defer LogErr()
	defer services.DB.Close()

	//routes
	router := httprouter.New()
	router.HandleMethodNotAllowed = false
	router.NotFound = http.HandlerFunc(NotFoundHandler)
	router.GET("/balance", controllers.GetBalance)
	router.POST("/deposite", controllers.AddAmount)
	router.POST("/withdraw", controllers.WithdrawAmount)
	router.POST("/transfer", controllers.TransferAmount)

	//start http
	Log.Printf("started at port: %s", Port)
	Log.Fatal(endless.ListenAndServe(":"+Port, router))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"` + services.ApiName + `.method_not_found": "Method not found"}`))
}
