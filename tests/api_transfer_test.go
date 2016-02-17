package test

import (
	"app/controllers"
	"app/models"
	"bytes"
	"compress/gzip"
	conf "github.com/andboson/configlog"
	"github.com/julienschmidt/httprouter"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestTransfer(t *testing.T) {
	var model = new(models.UserBalance)
	var model2 = new(models.UserBalance)
	AddMockDataModel(101, 100)
	AddMockDataModel(205, 100)

	conf.AppConfig.Set("enable_gzip", true)
	b := bytes.NewBufferString(`{"from":101, "to":205,"amount":25}`)
	r, _ := http.NewRequest("POST", "/transfer", b)
	r.Header.Add("Accept-Encoding", "deflate, gzip")
	w := httptest.NewRecorder()
	var reader io.ReadCloser

	router := httprouter.New()
	router.POST("/transfer", controllers.TransferAmount)
	router.ServeHTTP(w, r)

	var responseText string
	switch w.Header().Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(w.Body)
		var uncompressed []byte
		uncompressed, _ = ioutil.ReadAll(reader)
		responseText = string(uncompressed)
		defer reader.Close()
	default:
		responseText = w.Body.String()
	}

	log.Printf("\n response == ", responseText, w.Header().Get("Content-Encoding"))

	modelFrom := model.GetById(101)
	modelTo := model2.GetById(205)

	Convey("Subject: Test Deposite \n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Db record from balance must equals to 75", func() {
			So(modelFrom.Balance, ShouldEqual, 75)
		})

		Convey("Db record to balance must equals to 125", func() {
			So(modelTo.Balance, ShouldEqual, 125)
		})
	})
}
