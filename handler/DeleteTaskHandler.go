package handler

import (
	"encoding/json"
	"errors"

	"net/http"

	cases "TESTFIN/tasks"
)

// обработчик DELETE для "DELETE /api/task"
func DeleteTaskHandler(datab cases.Datab) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := datab.DeleteTask(id)
		if err != nil {
			err := errors.New("задача с таким id не найдена")
			cases.ErrorResponses.Error = err.Error()
			json.NewEncoder(w).Encode(cases.ErrorResponses)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, "ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}
	}
}
