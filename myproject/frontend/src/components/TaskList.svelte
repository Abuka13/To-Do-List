<script>
    export let tasks = [];
    export let onTaskUpdated = () => {};

    async function toggleTaskCompletion(task) {
        if (!window.go?.app?.App?.ToggleTaskCompletion) {
            console.error('Wails runtime недоступен');
            alert('Ошибка: соединение с приложением потеряно');
            return;
        }

        try {
            console.log('Переключение статуса задачи:', task.id);
            const updatedTask = await window.go.app.App.ToggleTaskCompletion(task.id);
            console.log('Статус переключен:', updatedTask);

            tasks = tasks.map(t => t.id === updatedTask.id ? updatedTask : t);
            onTaskUpdated();
        } catch (error) {
            console.error('Ошибка при обновлении задачи:', error);
            alert('Ошибка при обновлении задачи: ' + (error.message || error));
        }
    }

    async function deleteTask(taskId) {
        if (!window.go?.app?.App?.DeleteTask) {
            console.error('Wails runtime недоступен');
            alert('Ошибка: соединение с приложением потеряно');
            return;
        }

        if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
            try {
                console.log('Удаление задачи:', taskId);
                await window.go.app.App.DeleteTask(taskId);
                console.log('Задача удалена:', taskId);
                onTaskUpdated();
            } catch (error) {
                console.error('Ошибка при удалении задачи:', error);
                alert('Ошибка при удалении задачи: ' + (error.message || error));
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
            case 'high':
                return 'priority-high';
            case 'medium':
                return 'priority-medium';
            case 'low':
                return 'priority-low';
            default:
                return '';
        }
    }

    function isOverdue(dueDate) {
        if (!dueDate) return false;
        return new Date(dueDate) < new Date();
    }
</script>

<div class="task-list">
    {#if tasks.length === 0}
        <p class="empty-state">Задач пока нет. Создайте свою первую задачу!</p>
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
                                {#if isOverdue(task.dueDate) && !task.isCompleted}
                                    <span class="overdue-label">(Просрочено)</span>
                                {/if}
                            </span>
                        {/if}

                        <span class="task-priority">
                            Приоритет:
                            {#if task.priority === 'high'}Высокий
                            {:else if task.priority === 'medium'}Средний
                            {:else}Низкий
                            {/if}
                        </span>

                        <span class="task-created">
                            Создано: {formatDate(task.createdAt)}
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
        background: #f8f9fa;
        border-radius: 8px;
        border: 2px dashed #dee2e6;
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
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        border-left: 4px solid transparent;
    }

    .task-item:hover {
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
        transform: translateY(-1px);
    }

    .task-item.completed {
        opacity: 0.7;
        background: #f8f9fa;
    }

    /* ЦВЕТА ДЛЯ ПРИОРИТЕТОВ */
    .priority-high {
        border-left-color: #dc3545 !important;
        background: linear-gradient(90deg, rgba(220, 53, 69, 0.05) 0%, white 4%);
    }

    .priority-medium {
        border-left-color: #fd7e14 !important;
        background: linear-gradient(90deg, rgba(253, 126, 20, 0.05) 0%, white 4%);
    }

    .priority-low {
        border-left-color: #28a745 !important;
        background: linear-gradient(90deg, rgba(40, 167, 69, 0.05) 0%, white 4%);
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
        flex: 1;
    }

    .task-title {
        font-size: 16px;
        font-weight: 500;
        margin-bottom: 4px;
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

    .overdue-label {
        color: #dc3545;
        font-weight: bold;
    }

    .task-priority {
        font-size: 12px;
        color: #666;
        font-weight: 500;
    }

    /* Цвет текста для приоритетов */
    .priority-high .task-priority {
        color: #dc3545;
    }

    .priority-medium .task-priority {
        color: #fd7e14;
    }

    .priority-low .task-priority {
        color: #28a745;
    }

    .task-created {
        font-size: 11px;
        color: #999;
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
        background-color: #ffe6e6;
    }

    /* СТИЛИ ДЛЯ ТЕМНОЙ ТЕМЫ */
    :global([data-theme="dark"]) .task-item {
        background: #2d3748;
        border-color: #4a5568;
        color: #e2e8f0;
    }

    :global([data-theme="dark"]) .task-item.completed {
        background: #4a5568;
        opacity: 0.8;
    }

    :global([data-theme="dark"]) .task-item:hover {
        background: #3c4758;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
    }

    :global([data-theme="dark"]) .priority-high {
        background: linear-gradient(90deg, rgba(220, 53, 69, 0.15) 0%, #2d3748 4%) !important;
    }

    :global([data-theme="dark"]) .priority-medium {
        background: linear-gradient(90deg, rgba(253, 126, 20, 0.15) 0%, #2d3748 4%) !important;
    }

    :global([data-theme="dark"]) .priority-low {
        background: linear-gradient(90deg, rgba(40, 167, 69, 0.15) 0%, #2d3748 4%) !important;
    }

    :global([data-theme="dark"]) .empty-state {
        background: #2d3748;
        border-color: #4a5568;
        color: #a0aec0;
    }

    :global([data-theme="dark"]) .completed-text {
        color: #a0aec0;
    }

    :global([data-theme="dark"]) .task-due-date {
        color: #a0aec0;
    }

    :global([data-theme="dark"]) .task-due-date.overdue {
        color: #fc8181;
    }

    :global([data-theme="dark"]) .overdue-label {
        color: #fc8181;
    }

    :global([data-theme="dark"]) .task-priority {
        color: #a0aec0;
    }

    :global([data-theme="dark"]) .priority-high .task-priority {
        color: #fc8181;
    }

    :global([data-theme="dark"]) .priority-medium .task-priority {
        color: #f6ad55;
    }

    :global([data-theme="dark"]) .priority-low .task-priority {
        color: #68d391;
    }

    :global([data-theme="dark"]) .task-created {
        color: #718096;
    }

    :global([data-theme="dark"]) .delete-button:hover {
        background-color: rgba(220, 53, 69, 0.3);
        color: #fc8181;
    }
</style>