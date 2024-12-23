package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// Обработчик для выдачи всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Обработчик для добавления задачи
func postNewTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	// считывем новую запись в буфер
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// преобразуем данные в мапу
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// добавляем новую запись в мапу заданий
	tasks[task.ID] = task

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// в заголовок записываем код ответа, у нас это 201
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// проверяем и, если есть, записываем в переменную данные из мапы заданий
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// сериализуем данные из переменной task
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// в заголовок записываем тип контента
	w.Header().Set("Content-Type", "application/json")
	// в заголовок записываем код ответа, 200
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Обработчик для удаления задачи по ID
func delTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// проверяем и, если есть, записываем в переменную данные из мапы заданий
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// удаляем данные по id из мапы task
	delete(tasks, id)

	// в заголовок записываем тип контента
	w.Header().Set("Content-Type", "application/json")
	// в заголовок записываем код ответа, 200
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks` - все задания
	r.Get("/tasks", getTasks)

	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postNewTask`
	r.Post("/tasks", postNewTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `delTask`
	r.Delete("/tasks/{id}", delTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
