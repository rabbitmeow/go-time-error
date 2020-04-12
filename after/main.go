package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

var dbConn *sql.DB
var isAppReady bool

type shiftdata struct {
	Name         string    `json:"name"`
	TimeStart    time.Time `json:"-"`
	TimeEnd      time.Time `json:"-"`
	TimeStartStr string    `json:"time_start"`
	TimeEndStr   string    `json:"time_end"`
}

// rawTime is used for scanning mysql time data type
type rawTime []byte

func (t rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type resShift struct {
	ShiftData []shiftdata `json:"shift_data"`
}

type baseResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type requestHandler struct {
	timeoutContext time.Duration
}

func main() {
	var err error
	connURL := viper.GetString("database.user") + ":" + viper.GetString("database.password") + "@tcp(" + viper.GetString("database.host") + ":" + viper.GetString("database.port") + ")/" + viper.GetString("database.dbname") + "?charset=utf8&parseTime=True&loc=Local"
	dbConn, err = sql.Open("mysql", connURL)
	if err != nil {
		panic(err)
	}
	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}
	defer func() {
		dbConn.Close()
	}()
	timeoutContext := time.Duration(viper.GetInt("timeout")) * time.Second
	handleRequests(timeoutContext)
}

func (h *requestHandler) home(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Message string `json:"message"`
	}
	var response baseResponse
	var responseData res
	responseData.Message = "hello, welcome"
	response.Message = "success"
	response.Data = responseData
	response.Status = 200
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *requestHandler) shift(w http.ResponseWriter, r *http.Request) {
	var response baseResponse
	response.Data = new(struct{})
	reqCtx := r.Context()
	if reqCtx == nil {
		reqCtx = context.Background()
	}
	ctx, cancel := context.WithTimeout(reqCtx, h.timeoutContext)
	defer cancel()

	data, err := getShift(ctx)
	response.Data = data
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response.Status = 500
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Status = 200
		response.Message = "success"
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

func getShift(ctx context.Context) (resShift, error) {
	var resData resShift
	resData.ShiftData = make([]shiftdata, 0)
	query := "select name, time_start, time_end from shift order by id desc"
	rows, err := dbConn.QueryContext(ctx, query)
	if err != nil {
		return resData, err
	}
	defer func() {
		rows.Close()
	}()
	for rows.Next() {
		var t shiftdata
		var timeRaw1 []byte
		var timeRaw2 []byte
		err = rows.Scan(&t.Name, &timeRaw1, &timeRaw2)
		if err != nil {
			return resData, err
		}
		loc, _ := time.LoadLocation(viper.GetString("timezone"))
		timeNow := time.Now().In(loc)
		t.TimeStart, _ = time.Parse("15:04:05", string(timeRaw1))
		t.TimeEnd, _ = time.Parse("15:04:05", string(timeRaw2))
		year, month, day := timeNow.Date()
		h1, m1, s1 := t.TimeStart.Clock()
		h2, m2, s2 := t.TimeEnd.Clock()
		t.TimeStartStr = time.Date(year, month, day, h1, m1, s1, 0, loc).Format("02 Jan 2006 15:04:05")
		t.TimeEndStr = time.Date(year, month, day, h2, m2, s2, 0, loc).Format("02 Jan 2006 15:04:05")
		resData.ShiftData = append(resData.ShiftData, t)
	}
	return resData, err
}

func handleRequests(timeout time.Duration) {
	if isAppReady {
		log.Println("app ready on :" + viper.GetString("port"))
	}
	var handler requestHandler
	handler.timeoutContext = timeout
	http.HandleFunc("/", handler.home)
	http.HandleFunc("/shift", handler.shift)
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))
}
