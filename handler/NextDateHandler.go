package handler

import (
	"net/http"
	"time"

	cases "TESTFIN/tasks"
)

// обработчик "/api/nextdate"
func NextDateHandler(w http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	nowTime, err := time.Parse(cases.DateFormat, now)
	if err != nil {
		http.Error(w, "неверный формат даты", http.StatusBadRequest)
		return
	}
	nextDate, err := cases.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
