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

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	IsCompleted bool       `json:"isCompleted"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	Priority    Priority   `json:"priority"`
}

// NewTask - конструктор задач. Использовал здесь фабричный паттерн
func NewTask(title string, priority Priority, dueDate *time.Time) *Task {
	return &Task{
		Title:       title,
		IsCompleted: false,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC(),
		CompletedAt: nil,
	}
}

// эта функция отмечает задачу выполненной
func (t *Task) Complete() {
	if !t.IsCompleted {
		t.IsCompleted = true
		now := time.Now().UTC()
		t.CompletedAt = &now
	}
}

// эта функция чтобы выполенную задачу снова вернуть в не выполненную
func (t *Task) Reopen() {
	if t.IsCompleted {
		t.IsCompleted = false
		t.CompletedAt = nil
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
	if t.DueDate == nil || t.IsCompleted {
		return false
	}
	return t.DueDate.Before(time.Now().UTC())
}
