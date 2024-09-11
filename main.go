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
	// Определяем путь к файлу базы данных через переменную окружения
	r := chi.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./web")))

	// обработчики:
	r.HandleFunc("/api/nextdate", handler.NextDateHandler)

	r.Post("/api/task", handler.PostTaskHandler(datab))

	// запускаем сервер
	if err := http.ListenAndServe(":7540", r); err != nil {
		panic(err)

	}

}
