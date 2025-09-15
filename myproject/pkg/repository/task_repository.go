package repository

import "myproject/pkg/domain"

type TaskRepository interface {
	FindAll() ([]*domain.Task, error)
	FindByID(id int) (*domain.Task, error)
	FindByTitle(title string) ([]*domain.Task, error)
	Save(task *domain.Task) error
	Delete(id int) error
}
