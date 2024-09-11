package handler

import (
	"encoding/json"

	"net/http"

	cases "TESTFIN/tasks"
)

// обработчик POST "POST /api/task"
func PostTaskHandler(datab cases.Datab) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var task cases.Task
		err := json.NewDecoder(req.Body).Decode(&task)
		if err != nil {
			http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
			return
		}
		id, err := datab.AddTask(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reply := map[string]string{
			"id": id,
		}
		res, err := json.Marshal(reply)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
