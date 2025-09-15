package service

import (
	"myproject/pkg/domain"
	"myproject/pkg/repository"

	"strings"
	"time"
)

// TaskService содержит бизнес-логику приложения.
// Работает не с конкретным хранилищем, а с абстракцией TaskRepository.
// Паттерн: Сервис (Use Case Layer).
type TaskService struct {
	repo repository.TaskRepository
}

// NewTaskService - конструктор. Паттерн: Dependency Injection (репозиторий внедряется).
func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// GetAllTasks возвращает все задачи. Делегирует работу репозиторию.
func (s *TaskService) GetAllTasks() ([]*domain.Task, error) {
	return s.repo.FindAll()
}

// СОздание задачи
func (s *TaskService) CreateTask(title string, priority domain.Priority, dueDate *time.Time) (*domain.Task, error) {
	newTask := domain.NewTask(title, priority, dueDate)
	if err := newTask.Validate(); err != nil {
		return nil, err
	}
	err := s.repo.Save(newTask)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

// ищет таски по названию
func (s *TaskService) FindTasksByTitle(title string) ([]*domain.Task, error) {
	if strings.TrimSpace(title) == "" {
		//Если поиск пустой, возвращаем все задачи
		return s.repo.FindAll()
	}
	return s.repo.FindByTitle(title)
}

// Переключение статуса задачи
func (s *TaskService) ToggleTaskCompletion(id int) (*domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task.IsCompleted {
		task.Reopen()
	} else {
		task.Complete()
	}
	err = s.repo.Save(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}

// Бонусный метод, находит все таски где дедлайн сгорает сегодня
func (s *TaskService) GetTasksDueToday() ([]*domain.Task, error) {
	allTasks, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var dueToday []*domain.Task
	today := time.Now().UTC().Truncate(24 * time.Hour) // Начало текущих суток
	for _, task := range allTasks {
		if task.DueDate != nil {
			dueDate := task.DueDate.Truncate(24 * time.Hour)
			if dueDate.Equal(today) {
				dueToday = append(dueToday, task)
			}
		}
	}
	return dueToday, nil
}
