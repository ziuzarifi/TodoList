package models

import (
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Deadline time.Time `json:"deadline"`
}

func ReassignTask(taskID, oldUserID, newUserID int) error {
	_, err := db.Exec("UPDATE tasks SET user_id = $1 WHERE id = $2 AND user_id = $3", newUserID, taskID, oldUserID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTasks() ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Completed, &task.Deadline)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(id string) (Task, error) {
	var task Task

	err := db.QueryRow("SELECT * FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.UserID, &task.Title, &task.Completed, &task.Deadline)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func GetTasksByUserID(userID int) ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Completed, &task.Deadline)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}


func GetOverdueTasksByUserID(userID int) ([]Task, error) {
	var tasks []Task

	currentTime := time.Now()

	rows, err := db.Query("SELECT * FROM tasks WHERE user_id = $1 AND deadline < $2", userID, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Completed, &task.Deadline)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func CreateTask(userID int, title string, completed bool, deadline time.Time) (Task, error) {
	var task Task

	stmt, err := db.Prepare("INSERT INTO tasks (user_id, title, completed, deadline) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return Task{}, err
	}

	err = stmt.QueryRow(userID, title, completed, deadline).Scan(&task.ID)
	if err != nil {
		return Task{}, err
	}

	task.UserID = userID
	task.Title = title
	task.Completed = completed
	task.Deadline = deadline

	return task, nil
}

func UpdateTask(id string, userID int, title string, completed bool, deadline time.Time) error {
	stmt, err := db.Prepare("UPDATE tasks SET user_id = $1, title = $2, completed = $3, deadline = $4 WHERE id = $5")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(userID, title, completed, deadline, id)
	if err != nil {
		return err
	}
	

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Task not found")
	}

	return nil
}

func DeleteTask(id string) error {
	stmt, err := db.Prepare("DELETE FROM tasks WHERE id = $1")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Task not found")
	}

	return nil
}
