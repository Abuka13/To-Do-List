<script>
    let newTaskTitle = "";
    let newTaskPriority = "medium";
    let newTaskDueDate = "";

    export let onTaskAdded = () => {};

    async function addTask() {
        console.log('Кнопка добавления нажата');

        if (!window.go?.app?.App?.CreateTask) {
            console.error('Wails runtime недоступен');
            alert('Ошибка: соединение с приложением потеряно');
            return;
        }

        if (!newTaskTitle.trim()) {
            alert('Введите название задачи!');
            return;
        }
        if (newTaskDueDate) {
            const selectedDate = new Date(newTaskDueDate);
            const now = new Date();

            if (selectedDate < now) {
                if (!confirm('Выбранная дата уже прошла. Вы уверены, что хотите создать просроченную задачу?')) {
                    return;
                }
            }
        }
        try {
            console.log('Отправка запроса на создание задачи...');

            const result = await window.go.app.App.CreateTask(
                newTaskTitle.trim(),
                newTaskPriority,
                newTaskDueDate ? new Date(newTaskDueDate).toISOString() : ""
            );

            console.log('Задача создана успешно:', result);

            newTaskTitle = "";
            newTaskPriority = "medium";
            newTaskDueDate = "";

            onTaskAdded();

        } catch (error) {
            console.error('Ошибка создания задачи:', error);
            alert('Ошибка при создании задачи: ' + (error.message || error));
        }
    }

    function handleKeyPress(e) {
        if (e.key === 'Enter') {
            addTask();
        }
    }
</script>

<div class="add-task-form">
    <input
            type="text"
            bind:value={newTaskTitle}
            placeholder="Введите новую задачу..."
            on:keypress={handleKeyPress}
            class="task-input"
    />

    <select bind:value={newTaskPriority} class="priority-select">
        <option value="low">Низкий</option>
        <option value="medium">Средний</option>
        <option value="high">Высокий</option>
    </select>

    <input
            type="datetime-local"
            bind:value={newTaskDueDate}
            class="date-input"
            title="Дата и время выполнения"
    />

    <button on:click={addTask} class="add-button">
        Добавить
    </button>
</div>

<style>
    .add-task-form {
        display: flex;
        gap: 10px;
        margin-bottom: 20px;
        align-items: center;
        flex-wrap: wrap;
    }

    .task-input {
        flex: 1;
        min-width: 200px;
        padding: 12px;
        border: 2px solid var(--border-color);
        border-radius: 8px;
        font-size: 14px;
        transition: border-color 0.2s ease;
        background-color: var(--bg-primary);
        color: var(--text-primary);
    }

    .task-input:focus {
        outline: none;
        border-color: var(--accent-color);
        box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
    }

    .priority-select, .date-input {
        padding: 12px;
        border: 2px solid var(--border-color);
        border-radius: 8px;
        font-size: 14px;
        background-color: var(--bg-primary);
        color: var(--text-primary);
        cursor: pointer;
        transition: border-color 0.2s ease;
    }

    .priority-select:focus, .date-input:focus {
        outline: none;
        border-color: var(--accent-color);
        box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
    }

    .add-button {
        padding: 12px 20px;
        background-color: var(--accent-color);
        color: white;
        border: none;
        border-radius: 8px;
        cursor: pointer;
        font-weight: 600;
        font-size: 14px;
        transition: all 0.2s ease;
    }

    .add-button:hover {
        opacity: 0.9;
        transform: translateY(-1px);
    }

    .add-button:active {
        transform: translateY(0);
    }

    @media (max-width: 768px) {
        .add-task-form {
            flex-direction: column;
            align-items: stretch;
        }

        .task-input, .priority-select, .date-input, .add-button {
            width: 100%;
        }
    }
</style>