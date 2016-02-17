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


func TestDepositeNew(t *testing.T) {
	DeleteMockDataModel(101)

	conf.AppConfig.Set("enable_gzip", true)
	b := bytes.NewBufferString(`{"user":101, "amount":100}`)
	r, _ := http.NewRequest("POST", "/deposite", b)
	r.Header.Add("Accept-Encoding", "deflate, gzip")
	w := httptest.NewRecorder()
	var reader io.ReadCloser

	router := httprouter.New()
	router.POST("/deposite", controllers.AddAmount)
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
	var model = new(models.UserBalance)
	model = model.GetById(101)

	Convey("Subject: Test Deposite \n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Db record must be exists", func() {
			So(model.Balance, ShouldEqual, 100)
			So(model.Id, ShouldEqual, 101)
		})
	})
}

func TestDepositeExists(t *testing.T) {
	AddMockDataModel(101, 100)

	conf.AppConfig.Set("enable_gzip", true)
	b := bytes.NewBufferString(`{"user":101, "amount":100}`)
	r, _ := http.NewRequest("POST", "/deposite", b)
	r.Header.Add("Accept-Encoding", "deflate, gzip")
	w := httptest.NewRecorder()
	var reader io.ReadCloser

	router := httprouter.New()
	router.POST("/deposite", controllers.AddAmount)
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
	var model = new(models.UserBalance)
	model = model.GetById(101)

	Convey("Subject: Test Deposite \n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Db record balance must equals to 200", func() {
			So(model.Balance, ShouldEqual, 200)
		})
	})
}
