package main

import (
	handler "TESTFIN/handler"
	cases "TESTFIN/tasks"

	"net/http"

	"github.com/go-chi/chi"
	_ "modernc.org/sqlite"
)

func main() {
	db := cases.CreatDb()
	defer db.Close()
	datab := cases.NewDatab(db)

	r := chi.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./web")))

	// обработчики:
	r.HandleFunc("/api/nextdate", handler.NextDateHandler)

	r.Post("/api/task", handler.PostTaskHandler(datab))
	r.Get("/api/tasks", handler.GetTasksHandler(datab))
	r.Get("/api/task", handler.GetTaskHandler(datab))
	r.Put("/api/task", handler.PutTaskHandler(datab))
	r.Post("/api/task/done", handler.DoneTaskHandler(datab))
	r.Delete("/api/task", handler.DeleteTaskHandler(datab))
	// запускаем сервер
	if err := http.ListenAndServe(":7540", r); err != nil {
		panic(err)

	}

}
