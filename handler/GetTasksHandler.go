package handler

import (
	"encoding/json"
	"errors"

	"net/http"

	cases "TESTFIN/tasks"
)

// обработчик GET для "/api/tasks"
func GetTasksHandler(datab cases.Datab) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		taskSearch := req.URL.Query().Get("search")
		tasks, err := datab.GetTasks(taskSearch)

		if err != nil {
			err := errors.New(`ошибка запроса к базе данных`)
			cases.ErrorResponses.Error = err.Error()
			json.NewEncoder(w).Encode(cases.ErrorResponses)
			return
		}

		reply := map[string][]cases.Task{
			"tasks": tasks,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(reply); err != nil {
			http.Error(w, "ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}
	}
}
