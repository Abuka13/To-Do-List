package repository

import (
	"encoding/json"
	"errors"
	"myproject/pkg/domain"
	"os"
	"strings"
	"sync"
)

// JSONTaskRepository - конкретная реализация TaskRepository для хранения в JSON файле.
type JSONTaskRepository struct {
	filePath string
	mu       sync.RWMutex
	tasks    []*domain.Task
	nextID   int
}

// NewJSONTaskRepository - конструктор репозитория.
func NewJSONTaskRepository(filePath string) *JSONTaskRepository {
	repo := &JSONTaskRepository{
		filePath: filePath,
		mu:       sync.RWMutex{},
		tasks:    make([]*domain.Task, 0),
		nextID:   1,
	}

	if err := repo.loadFromFile(); err != nil {
		repo.tasks = make([]*domain.Task, 0)
		repo.nextID = 1
	}

	return repo
}

// FindAll возвращает все задачи из репозитория.
func (r *JSONTaskRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock() // Блокировка для чтения
	defer r.mu.RUnlock()

	// Возвращаем копию среза, чтобы избежать изменений извне
	result := make([]*domain.Task, len(r.tasks))
	copy(result, r.tasks)
	return result, nil
}

// Ищет по ID
func (r *JSONTaskRepository) FindByID(id int) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, task := range r.tasks {
		if task.ID == id {
			foundTask := *task
			return &foundTask, nil
		}
	}
	return nil, errors.New("task not found")
}

// Ищет по title
func (r *JSONTaskRepository) FindByTitle(title string) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	searchTerm := strings.ToLower(title)
	var result []*domain.Task

	for _, task := range r.tasks {
		if strings.Contains(strings.ToLower(task.Title), searchTerm) {
			foundTask := *task
			result = append(result, &foundTask)
		}
	}

	return result, nil
}

// Save сохраняет или обновляет задачу в репозитории.
func (r *JSONTaskRepository) Save(task *domain.Task) error {
	r.mu.Lock() // Полная блокировка для записи
	defer r.mu.Unlock()

	if task.ID == 0 {
		task.ID = r.nextID
		r.nextID++
		r.tasks = append(r.tasks, task)
	} else {
		found := false
		for i, existingTask := range r.tasks {
			if existingTask.ID == task.ID {
				r.tasks[i] = task
				found = true
				break
			}
		}
		if !found {
			return errors.New("Задача не найдена!")
		}
	}

	return r.persistToFile()
}

// Delete удаляет задачу по ID.
func (r *JSONTaskRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			// Удаляем задачу из среза
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return r.persistToFile()
		}
	}
	return errors.New("Задача не найдена!")
}

// загружает задачи из JSON файла.
func (r *JSONTaskRepository) loadFromFile() error {

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var tasks []*domain.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	r.tasks = tasks
	r.nextID = maxID + 1
	return nil
}

func (r *JSONTaskRepository) persistToFile() error {
	// Конвертируем в JSON с форматированием
	data, err := json.MarshalIndent(r.tasks, "", "  ")
	if err != nil {
		return err
	}

	// Создаем файл если не существует, с правами 0644 (rw-r--r--)
	return os.WriteFile(r.filePath, data, 0644)
}

// Функция чтобы посмотреть сколько задач в общем и сколько выполенных
func (r *JSONTaskRepository) GetStats() (int, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.tasks)
	completed := 0
	for _, task := range r.tasks {
		if task.IsCompleted {
			completed++
		}
	}

	return total, completed, nil
}
