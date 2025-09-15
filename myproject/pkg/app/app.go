package app

import (
	"myproject/pkg/domain"
	"myproject/pkg/repository"
	"myproject/pkg/service"
)

type App struct {
	taskService *service.TaskService
}

func NewApp() *App {
	repo := repository.NewJSONTaskRepository("tasks.json")
	taskService := service.NewTaskService(repo)
	return &App{taskService: taskService}
}

// метод для получения всех задач.
func (a *App) GetAllTasks() ([]*domain.Task, error) {
	return a.taskService.GetAllTasks()
}

// метод для создания новой задачи.
func (a *App) CreateTask(title string, priorityStr string, dueDateStr string) (*domain.Task, error) {
	// Парсим приоритет из строки в тип Priority
	priority := domain.Priority(priorityStr)

	// Валидируем приоритет
	if priority != domain.PriorityLow &&
		priority != domain.PriorityMedium &&
		priority != domain.PriorityHigh {
		priority = domain.PriorityMedium // Значение по умолчанию
	}

	return a.taskService.CreateTask(title, priority, dueDateStr)
}

// метод для поиска задач по заголовку.
func (a *App) SearchTasks(title string) ([]*domain.Task, error) {
	return a.taskService.FindTasksByTitle(title)
}

// метод для переключения статуса задачи.
func (a *App) ToggleTaskCompletion(id int) (*domain.Task, error) {
	return a.taskService.ToggleTaskCompletion(id)
}

// метод для удаления задачи.
func (a *App) DeleteTask(id int) error {
	return a.taskService.DeleteTask(id)
}

// метод для получения статистики (бонус).
func (a *App) GetStats() (map[string]int, error) {
	tasks, err := a.taskService.GetAllTasks()
	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total":     len(tasks),
		"completed": 0,
		"active":    0,
		"overdue":   0,
	}

	for _, task := range tasks {
		if task.IsCompleted {
			stats["completed"]++
		} else {
			stats["active"]++
			if task.IsOverdue() {
				stats["overdue"]++
			}
		}
	}

	return stats, nil
}
