package controllers

import (
	. "app/common"
	"app/models"
	"encoding/json"
	conf "github.com/andboson/configlog"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strconv"
)

const BAD_REQUEST_CUSTOM = 422

type InputRequest struct {
	Amount  int `json:"amount"`
	UserId  int `json:"user"`
	FromUid int `json:"from"`
	ToUid   int `json:"to"`
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}

type UserBalanceController struct {
	Request  *http.Request
	Response http.ResponseWriter
	HttpLib
}

///handlers

// http handler GetBalance
func GetBalance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	controller := UserBalanceController{Request: r, Response: w}
	controller.getBalance()
}

// http handler AddAmount
func AddAmount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	controller := UserBalanceController{Request: r, Response: w}
	controller.addAmount()
}

// http handler WithdrawAmount
func WithdrawAmount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	controller := UserBalanceController{Request: r, Response: w}
	controller.withdrawAmount()
}

// http handler TransferAmount
func TransferAmount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	controller := UserBalanceController{Request: r, Response: w}
	controller.transferAmount()
}

/////////////// methods

// get balance bu uid
func (d *UserBalanceController) getBalance() {
	var model = new(models.UserBalance)
	var userId int
	debug, _ := conf.AppConfig.Bool("debug")
	defer LogErr()
	request := d.Request.URL.Query()
	userIdString := request.Get("user")
	userId, _ = strconv.Atoi(userIdString)
	if debug {
		Log.WithField("request", request).Printf("balance", userIdString)
	}
	if userId == 0 {
		d.makeIncorrectRequestParamsReponse("user param is required")
		return
	}

	model = model.GetById(userId)
	if model.Id != 0 {
		resp := &BalanceResponse{
			Balance: model.Balance,
		}
		d.ResponseWriter(d.Request, d.Response, resp)
	} else {
		d.makeNotFound("balance.not_found", "User not found")
		return
	}

	return
}

// add amount to user balance
func (d *UserBalanceController) addAmount() {
	var input InputRequest
	var model = new(models.UserBalance)

	debug, _ := conf.AppConfig.Bool("debug")
	defer LogErr()
	body, error := ioutil.ReadAll(d.Request.Body)
	error = json.Unmarshal(body, &input)
	if debug {
		Log.WithField("request", string(body)).Printf("deposit")
	}
	if error != nil {
		Log.WithField("error", error).Printf("read body or json error")
	}
	if input.UserId == 0 || input.Amount == 0 {
		d.makeIncorrectRequestParamsReponse("User and Amount params is required")
		return
	}

	model = model.GetById(input.UserId)
	if model.Id == 0 {
		model = models.CreateUserBalance(input.UserId, input.Amount, "")
	} else {
		model.AddAmount(input.Amount)
	}

	d.ResponseWriter(d.Request, d.Response, "")

	return
}

// withdraw amount from user balance
func (d *UserBalanceController) withdrawAmount() {
	var input InputRequest
	var model = new(models.UserBalance)

	debug, _ := conf.AppConfig.Bool("debug")
	defer LogErr()
	body, error := ioutil.ReadAll(d.Request.Body)
	error = json.Unmarshal(body, &input)
	if debug {
		Log.WithField("request", string(body)).Printf("deposit")
	}
	if error != nil {
		Log.WithField("error", error).Printf("read body or json error")
	}
	if input.UserId == 0 || input.Amount == 0 {
		d.makeIncorrectRequestParamsReponse("User and Amount params is required")
		return
	}

	model = model.GetById(input.UserId)
	if model.Id == 0 {
		d.makeNotFound("balance.not_found", "User not found")
		return
	}

	model.WithdrawAmount(input.Amount)
	d.ResponseWriter(d.Request, d.Response, "")

	return
}

// transfer amount from between users balance
func (d *UserBalanceController) transferAmount() {
	var input InputRequest
	var from = new(models.UserBalance)
	var to = new(models.UserBalance)
	debug, _ := conf.AppConfig.Bool("debug")
	defer LogErr()
	body, error := ioutil.ReadAll(d.Request.Body)
	error = json.Unmarshal(body, &input)
	if debug {
		Log.WithField("request", string(body)).Printf("transfer - %+v", input)
	}
	if error != nil {
		Log.WithField("error", error).Printf("read body or json error")
	}
	if input.FromUid == 0 || input.Amount == 0 || input.ToUid == 0 {
		d.makeIncorrectRequestParamsReponse("From, To and Amount params is required")
		return
	}

	from = from.GetById(input.FromUid)
	to = to.GetById(input.ToUid)
	if from.Id == 0 || to.Id == 0 {
		d.makeNotFound("balance.not_found", "User not found")
		return
	}
	from.WithdrawAmount(input.Amount)
	to.AddAmount(input.Amount)

	d.ResponseWriter(d.Request, d.Response, "")

	return
}

/////////

// not found
func (d *UserBalanceController) makeNotFound(method, message string) {
	d.Response.WriteHeader(http.StatusNotFound)
	d.Response.Write([]byte(`{"model.` + method + `": "` + message + `"}`))
}

// incorrect params
func (d *UserBalanceController) makeIncorrectRequestParamsReponse(error string) {
	d.Response.WriteHeader(BAD_REQUEST_CUSTOM)
	d.Response.Write([]byte(error))
}
