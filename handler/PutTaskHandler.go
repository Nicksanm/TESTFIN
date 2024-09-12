package handler

import (
	"encoding/json"

	"net/http"

	cases "TESTFIN/tasks"
)

// обработчик PUT для "PUT /api/task"
func PutTaskHandler(datab cases.Datab) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var task cases.Task
		err := json.NewDecoder(req.Body).Decode(&task)
		if err != nil {
			http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
			return
		}
		err = datab.UpdateTask(task)
		if err != nil {
			http.Error(w, "задача не найдена", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, "ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}
	}
}
