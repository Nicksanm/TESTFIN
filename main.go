package main

import (
	handler "TESTFIN/handler"
	cases "TESTFIN/tasks"

	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	db := cases.CreatDb()
	defer db.Close()
	datab := cases.NewDatab(db)
	// Определяем путь к файлу базы данных через переменную окружения

	http.Handle("/", http.FileServer(http.Dir("./web")))

	// обработчики:
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)

	http.HandleFunc("/api/task", handler.PostTaskHandler(datab))

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}
