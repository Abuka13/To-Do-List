<script>
    import { onMount } from 'svelte';

    let newTaskTitle = "";
    let newTaskPriority = "medium";
    let newTaskDueDate = "";

    export let onTaskAdded = () => {};

    async function addTask() {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.error('Wails runtime not available');
            return;
        }

        if (newTaskTitle.trim()) {
            try {
                await window.go.main.App.CreateTask(
                    newTaskTitle.trim(),
                    newTaskPriority,
                    newTaskDueDate
                );

                newTaskTitle = "";
                newTaskPriority = "medium";
                newTaskDueDate = "";

                onTaskAdded();

            } catch (error) {
                console.error("Ошибка при добавлении задачи:", error);
            }
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
    }

    .task-input {
        flex: 1;
        padding: 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
    }

    .priority-select, .date-input {
        padding: 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
    }

    .add-button {
        padding: 8px 16px;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    .add-button:hover {
        background-color: #0056b3;
    }
</style>