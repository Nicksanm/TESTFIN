package handler

import (
	"encoding/json"

	"net/http"

	cases "TESTFIN/tasks"
)

// обработчик POST "POST /api/task"
func PostTaskHandler(datab cases.Datab) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var t cases.Task
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		id, err := datab.AddTask(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := cases.Response{ID: id}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
