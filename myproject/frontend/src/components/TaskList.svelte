<script>
    export let tasks = [];
    export let onTaskUpdated = () => {};

    async function toggleTaskCompletion(task) {
        if (!window.go?.main?.App) {
            console.error('Wails runtime not available');
            return;
        }

        try {
            const updatedTask = await window.go.main.App.ToggleTaskCompletion(task.id);
            tasks = tasks.map(t => t.id === updatedTask.id ? updatedTask : t);
            onTaskUpdated();
        } catch (error) {
            console.error('Ошибка при обновлении задачи:', error);
        }
    }

    async function deleteTask(taskId) {
        if (!window.go?.main?.App) {
            console.error('Wails runtime not available');
            return;
        }

        if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
            try {
                await window.go.main.App.DeleteTask(taskId);
                onTaskUpdated();
            } catch (error) {
                console.error('Ошибка при удалении задачи:', error);
            }
        }
    }

    function formatDate(dateString) {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function getPriorityClass(priority) {
        switch (priority) {
            case 'high': return 'priority-high';
            case 'medium': return 'priority-medium';
            case 'low': return 'priority-low';
            default: return '';
        }
    }

    function isOverdue(dueDate) {
        if (!dueDate) return false;
        return new Date(dueDate) < new Date();
    }
</script>

<div class="task-list">
    {#if tasks.length === 0}
        <p class="empty-state">Задач пока нет</p>
    {:else}
        {#each tasks as task (task.id)}
            <div class="task-item {task.isCompleted ? 'completed' : ''} {getPriorityClass(task.priority)}">
                <div class="task-content">
                    <input
                            type="checkbox"
                            checked={task.isCompleted}
                            on:change={() => toggleTaskCompletion(task)}
                            class="task-checkbox"
                    />

                    <div class="task-info">
                        <span class="task-title {task.isCompleted ? 'completed-text' : ''}">
                            {task.title}
                        </span>

                        {#if task.dueDate}
                            <span class="task-due-date {isOverdue(task.dueDate) && !task.isCompleted ? 'overdue' : ''}">
                                📅 {formatDate(task.dueDate)}
                                {isOverdue(task.dueDate) && !task.isCompleted ? ' (Просрочено)' : ''}
                            </span>
                        {/if}

                        <span class="task-priority">
                            Приоритет:
                            {#if task.priority === 'high'}Высокий
                            {:else if task.priority === 'medium'}Средний
                            {:else}Низкий
                            {/if}
                        </span>
                    </div>
                </div>

                <button
                        on:click={() => deleteTask(task.id)}
                        class="delete-button"
                        title="Удалить задачу"
                >
                    🗑️
                </button>
            </div>
        {/each}
    {/if}
</div>

<style>
    .task-list {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .empty-state {
        text-align: center;
        color: #666;
        font-style: italic;
        padding: 40px;
    }

    .task-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 15px;
        border: 1px solid #ddd;
        border-radius: 8px;
        background: white;
        transition: all 0.3s ease;
    }

    .task-item:hover {
        box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }

    .task-item.completed {
        opacity: 0.7;
        background: #f8f9fa;
    }

    .task-content {
        display: flex;
        align-items: flex-start;
        gap: 12px;
        flex: 1;
    }

    .task-checkbox {
        margin-top: 2px;
        cursor: pointer;
        width: 18px;
        height: 18px;
    }

    .task-info {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .task-title {
        font-size: 16px;
        font-weight: 500;
    }

    .completed-text {
        text-decoration: line-through;
        color: #666;
    }

    .task-due-date {
        font-size: 12px;
        color: #666;
    }

    .task-due-date.overdue {
        color: #dc3545;
        font-weight: bold;
    }

    .task-priority {
        font-size: 12px;
        color: #666;
    }

    .delete-button {
        background: none;
        border: none;
        cursor: pointer;
        padding: 8px;
        border-radius: 4px;
        font-size: 16px;
        transition: background-color 0.2s ease;
    }

    .delete-button:hover {
        background-color: #f8f9fa;
    }

    /* Стили для приоритетов */
    .priority-high {
        border-left: 4px solid #dc3545;
    }

    .priority-medium {
        border-left: 4px solid #ffc107;
    }

    .priority-low {
        border-left: 4px solid #28a745;
    }
</style>