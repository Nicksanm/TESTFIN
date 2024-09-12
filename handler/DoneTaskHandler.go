package handler

import (
	"encoding/json"

	"net/http"

	cases "TESTFIN/tasks"
)

// обработчик для "/api/task/done"
func DoneTaskHandler(datab cases.Datab) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := datab.DoneTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, "ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}
	}
}
