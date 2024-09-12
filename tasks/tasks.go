package cases

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"strconv"
	"strings"

	"log"
	"time"
)

const (
	LimitTasks = 30
	DateFormat = "20060102"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type Datab struct {
	db *sql.DB
}
type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

var ErrorResponses struct {
	Error string `json:"error,omitempty"`
}

// Создаем базу данных
func CreatDb() *sql.DB {

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()
	if install {
		Table := `CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL,
				title TEXT NOT NULL,
				comment TEXT,
				repeat VARCHAR(128) NOT NULL
				);`
		_, err = db.Exec(Table)
		if err != nil {
			log.Fatal(err)
		}

		Index := `CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);`
		_, err = db.Exec(Index)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("База данных создана")
	}
	return db
}
func NewDatab(db *sql.DB) Datab {
	return Datab{db: db}
}
func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("не указана строка")
	}

	nowDate, err := time.Parse(DateFormat, date)

	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	parts := strings.Split(repeat, " ")

	editParts := parts[0]

	switch editParts {
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		moreDays, err := strconv.Atoi(parts[1])
		if err != nil || moreDays < 1 || moreDays > 400 {
			return "", fmt.Errorf("превышен максимально допустимый интервал дней")
		}
		newDate := nowDate.AddDate(0, 0, moreDays)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, moreDays)
		}
		return newDate.Format(DateFormat), nil

	case "y":
		newDate := nowDate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("неверный ввод")

	}
}

func (d *Datab) AddTask(task Task) (string, error) {
	var err error
	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	}

	_, err = time.Parse(DateFormat, task.Date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}
	// Если дата меньше time.Now, то устанавливаем NextDate
	if task.Date < time.Now().Format(DateFormat) {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return "", fmt.Errorf("неверное правило повторения")
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(DateFormat)
		}
	}
	if task.Title == "" {
		return "", fmt.Errorf("не указан заголовок задачи")
	}

	// Добавляем задачу в базу данных
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`
	res, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("задача не добавлена")
	}

	//  Идентификатор созданной задачи
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("id созданной задачи не удалось вернуть")
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

// / Получаем список ближайших задач
func (d *Datab) GetTasks(search string) ([]Task, error) {

	var rows *sql.Rows
	var err error

	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		log.Println("Поиск(текст)")

		rows, err = d.db.Query(
			`SELECT id, date, title, comment, repeat 
			FROM scheduler
			WHERE title LIKE :object OR comment LIKE :object
			ORDER BY date LIMIT :limit`,
			sql.Named("object", "%"+search+"%"),
			sql.Named("limit", LimitTasks),
		)
	} else {
		log.Println("Поиск(дата)")

		object := date.Format("20060102")

		rows, err = d.db.Query(
			`SELECT id, date, title, comment, repeat
			   FROM scheduler
			   WHERE date LIKE :object
			   ORDER BY date LIMIT :limit`,
			sql.Named("object", "%"+object+"%"),
			sql.Named("limit", LimitTasks),
		)
	}
	if err != nil {
		log.Println("задачу не найти", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = []Task{}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tasks, nil
}

// Получение задачи по id
func (d *Datab) GetTask(id string) (Task, error) {
	var task Task
	if id == "" {
		return Task{}, fmt.Errorf("не указан id")
	}
	row := d.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, fmt.Errorf("задача не найдена")
	}
	return task, nil
}

// Редактирование задачи
func (d *Datab) UpdateTask(task Task) error {
	// проверяем поля в базе данных (id, date, title)
	if task.ID == "" {
		return fmt.Errorf("не указан идентификатор")
	}
	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	}

	_, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты")
	}
	// Если дата меньше time.Now, то устанавливаем NextDate
	if task.Date < time.Now().Format(DateFormat) {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("неверное правило повторения")
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(DateFormat)
		}
	}

	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}
	/// что-то не так
	// Обновляем задачу в базе данных
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка в обновлении данных")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось посчитать измененные строки")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача с таким id не найдена")
	}

	log.Println("задача изменена")

	return nil
}

// Выполнение задачи
func (d *Datab) DoneTask(id string) error {
	var task Task

	task, err := d.GetTask(id)
	if err != nil {
		return err
	}
	// Одноразовая задача с пустым полем repeat удаляется
	if task.Repeat == "" {

		err := d.DeleteTask(id)
		if err != nil {
			return err
		}

	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat) // расчет следующей даты
		if err != nil {
			return err
		}
		task.Date = next //изменение у задачи значение Date
		err = d.UpdateTask(task)
		if err != nil {
			return err
		}
	}

	return nil
}

// Удаление задачи из базы данных
func (d *Datab) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("не указан id")
	}
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить задачу")
	}

	rowsAffected, err := res.RowsAffected() // определяем измененные строки
	if err != nil {
		return fmt.Errorf("не удалось посчитать измененные строки")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с таким id не найдена")
	}
	return nil
}
