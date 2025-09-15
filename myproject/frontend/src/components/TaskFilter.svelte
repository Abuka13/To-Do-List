<script>
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    export let currentFilter = 'all';

    const filters = [
        { value: 'all', label: 'Все задачи' },
        { value: 'active', label: 'Активные' },
        { value: 'completed', label: 'Выполненные' }
    ];

    async function handleFilterChange(filter) {
        try {
            const filteredTasks = await window.go.app.App.GetTasksByFilter(filter);
            currentFilter = filter;
            // Исправляем вызов события
            dispatch('filterchange', { tasks: filteredTasks, filter });
        } catch (error) {
            console.error('Ошибка при фильтрации задач:', error);
        }
    }
</script>

<div class="filter-buttons">
    {#each filters as filter}
        <button
                class:active={currentFilter === filter.value}
                on:click={() => handleFilterChange(filter.value)}
        >
            {filter.label}
        </button>
    {/each}
</div>

<style>
    .filter-buttons {
        display: flex;
        gap: 10px;
        margin: 20px 0;
    }

    button {
        padding: 8px 16px;
        border: 1px solid #ddd;
        background: white;
        border-radius: 6px;
        cursor: pointer;
    }

    button.active {
        background: #2196F3;
        color: white;
        border-color: #2196F3;
    }
</style>