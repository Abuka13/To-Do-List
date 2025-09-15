package app

import (
	"fmt"
	"myproject/pkg/domain"
	"myproject/pkg/repository"
	"myproject/pkg/service"
	"time"
)

type Task = domain.Task

type App struct {
	taskService *service.TaskService
}

func NewApp() *App {
	repo := repository.NewJSONTaskRepository("tasks.json")
	taskService := service.NewTaskService(repo)
	return &App{taskService: taskService}
}

func (a *App) GetAllTasks() ([]*Task, error) {
	fmt.Println("App.GetAllTasks() вызван")
	tasks, err := a.taskService.GetAllTasks()
	if err != nil {
		fmt.Printf("Ошибка получения задач: %v\n", err)
		return nil, err
	}
	fmt.Printf("Получено %d задач\n", len(tasks))
	return tasks, nil
}

func (a *App) CreateTask(title string, priorityStr string, dueDateStr string) (*Task, error) {
	fmt.Printf("App.CreateTask() вызван: title='%s', priority='%s', dueDate='%s'\n",
		title, priorityStr, dueDateStr)

	if title == "" {
		fmt.Println("Пустое название задачи")
		return nil, fmt.Errorf("название задачи не может быть пустым")
	}

	priority := domain.Priority(priorityStr)

	if priority != domain.PriorityLow &&
		priority != domain.PriorityMedium &&
		priority != domain.PriorityHigh {
		fmt.Printf("Неверный приоритет '%s', устанавливаем средний\n", priorityStr)
		priority = domain.PriorityMedium
	}

	result, err := a.taskService.CreateTask(title, priority, dueDateStr)
	if err != nil {
		fmt.Printf("Ошибка создания задачи: %v\n", err)
		return nil, err
	}

	fmt.Printf("Задача создана успешно: ID=%d, Title='%s'\n", result.ID, result.Title)
	return result, nil
}

func (a *App) TestMethod() string {
	fmt.Println("TestMethod вызван из фронтенда")
	return "Связь с Go работает! " + time.Now().Format("15:04:05")
}

func (a *App) SearchTasks(title string) ([]*Task, error) {
	fmt.Printf("App.SearchTasks() вызван: title='%s'\n", title)
	tasks, err := a.taskService.FindTasksByTitle(title)
	if err != nil {
		fmt.Printf("Ошибка поиска задач: %v\n", err)
		return nil, err
	}
	fmt.Printf("Найдено %d задач\n", len(tasks))
	return tasks, nil
}

func (a *App) ToggleTaskCompletion(id int) (*Task, error) {
	fmt.Printf("App.ToggleTaskCompletion() вызван: id=%d\n", id)

	task, err := a.taskService.ToggleTaskCompletion(id)
	if err != nil {
		fmt.Printf("Ошибка переключения статуса задачи: %v\n", err)
		return nil, err
	}

	status := "выполнена"
	if !task.IsCompleted {
		status = "активна"
	}
	fmt.Printf("Задача ID=%d теперь %s\n", id, status)

	return task, nil
}

func (a *App) DeleteTask(id int) error {
	fmt.Printf("App.DeleteTask() вызван: id=%d\n", id)

	err := a.taskService.DeleteTask(id)
	if err != nil {
		fmt.Printf("Ошибка удаления задачи: %v\n", err)
		return err
	}

	fmt.Printf("Задача ID=%d удалена\n", id)
	return nil
}

func (a *App) GetStats() (map[string]int, error) {
	fmt.Println("App.GetStats() вызван")

	tasks, err := a.taskService.GetAllTasks()
	if err != nil {
		fmt.Printf("Ошибка получения статистики: %v\n", err)
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

	fmt.Printf("Статистика: всего=%d, выполнено=%d, активных=%d, просрочено=%d\n",
		stats["total"], stats["completed"], stats["active"], stats["overdue"])

	return stats, nil
}
