package domain

import (
	"errors"
	"time"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type TaskFilter string

const (
	FilterAll       TaskFilter = "all"
	FilterActive    TaskFilter = "active"
	FilterCompleted TaskFilter = "completed"
)

type Task struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	IsCompleted bool     `json:"isCompleted"`
	CreatedAt   string   `json:"createdAt"`
	CompletedAt string   `json:"completedAt,omitempty"`
	DueDate     string   `json:"dueDate,omitempty"`
	Priority    Priority `json:"priority"`
}

// NewTask - конструктор задач. Использовал здесь фабричный паттерн
func NewTask(title string, priority Priority, dueDate string) *Task {
	return &Task{
		Title:       title,
		IsCompleted: false,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		CompletedAt: "",
	}
}

// эта функция отмечает задачу выполненной
func (t *Task) Complete() {
	if !t.IsCompleted {
		t.IsCompleted = true
		t.CompletedAt = time.Now().UTC().Format(time.RFC3339)
	}
}

// эта функция чтобы выполенную задачу снова вернуть в не выполненную
func (t *Task) Reopen() {
	if t.IsCompleted {
		t.IsCompleted = false
		t.CompletedAt = ""
	}
}

// здесь валидация чтобы проверить пустое поле
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	return nil
}

// эта функция чтобы проверить не просрочена ли задача
func (t *Task) IsOverdue() bool {
	if t.DueDate == "" || t.IsCompleted {
		return false
	}

	dueTime, err := time.Parse(time.RFC3339, t.DueDate)
	if err != nil {
		return false
	}

	return dueTime.Before(time.Now().UTC())
}
