package test

import (
	"app/controllers"
	"app/models"
	"bytes"
	"compress/gzip"
	"encoding/json"
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

func TestGetBalance(t *testing.T) {
	AddMockDataModel(101, 100)

	conf.AppConfig.Set("enable_gzip", true)
	b := bytes.NewBufferString(`{"name":"test-name"}`)
	r, _ := http.NewRequest("GET", "/balance?user=101", b)
	r.Header.Add("Accept-Encoding", "deflate, gzip")
	w := httptest.NewRecorder()
	var reader io.ReadCloser

	router := httprouter.New()
	router.GET("/balance", controllers.GetBalance)
	router.ServeHTTP(w, r)

	var response models.UserBalance
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
	json.Unmarshal([]byte(responseText), &response)
	log.Printf("\n response mapped: %+v", response)

	Convey("Subject: Test get balance \n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result `balance` must be equal 100", func() {
			So(response.Balance, ShouldEqual, 100)
		})

	})
}
