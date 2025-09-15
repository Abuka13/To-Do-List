<script>
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    export let currentFilter = 'all';
    export let dateFilter = 'all';
    export let sortBy = 'date';

    const statusFilters = [
        {value: 'all', label: 'Все задачи', icon: '📋'},
        {value: 'active', label: 'Активные', icon: '⏳'},
        {value: 'completed', label: 'Выполненные', icon: '✅'}
    ];

    const dateFilters = [
        {value: 'all', label: 'Все даты', icon: '📅'},
        {value: 'today', label: 'Сегодня', icon: '☀️'},
        {value: 'week', label: 'На неделе', icon: '📆'},
        {value: 'overdue', label: 'Просроченные', icon: '⚠️'}
    ];

    const sortOptions = [
        {value: 'date', label: 'По дате', icon: '📊'},
        {value: 'priority', label: 'По приоритету', icon: '⚡'}
    ];

    function handleStatusFilterChange(filter) {
        currentFilter = filter;
        dispatchFilterChange();
    }

    function handleDateFilterChange(filter) {
        dateFilter = filter;
        dispatchFilterChange();
    }

    function handleSortChange(sort) {
        sortBy = sort;
        dispatchFilterChange();
    }

    function dispatchFilterChange() {
        dispatch('filterchange', {
            statusFilter: currentFilter,
            dateFilter: dateFilter,
            sortBy: sortBy
        });
    }
</script>

<div class="sidebar-filters">
    <!-- Статус -->
    <div class="filter-section">
        <h4>📋 Статус</h4>
        <div class="filter-options">
            {#each statusFilters as filter}
                <button
                        class:active={currentFilter === filter.value}
                        on:click={() => handleStatusFilterChange(filter.value)}
                >
                    <span class="icon">{filter.icon}</span>
                    <span class="label">{filter.label}</span>
                </button>
            {/each}
        </div>
    </div>

    <!-- Дата -->
    <div class="filter-section">
        <h4>📅 Дата</h4>
        <div class="filter-options">
            {#each dateFilters as filter}
                <button
                        class:active={dateFilter === filter.value}
                        on:click={() => handleDateFilterChange(filter.value)}
                >
                    <span class="icon">{filter.icon}</span>
                    <span class="label">{filter.label}</span>
                </button>
            {/each}
        </div>
    </div>

    <!-- Сортировка -->
    <div class="filter-section">
        <h4>🔍 Сортировка</h4>
        <div class="filter-options">
            {#each sortOptions as option}
                <button
                        class:active={sortBy === option.value}
                        on:click={() => handleSortChange(option.value)}
                >
                    <span class="icon">{option.icon}</span>
                    <span class="label">{option.label}</span>
                </button>
            {/each}
        </div>
    </div>
</div>

<style>
    .sidebar-filters {
        padding: 20px;
        display: flex;
        flex-direction: column;
        gap: 25px;
    }

    .filter-section {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .filter-section h4 {
        margin: 0;
        font-size: 14px;
        color: var(--text-primary);
        font-weight: 600;
        padding-bottom: 5px;
        border-bottom: 1px solid var(--border-color);
    }

    .filter-options {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    button {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 12px;
        border: 1px solid var(--border-color);
        background: var(--bg-primary);
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s ease;
        text-align: left;
        color: var(--text-primary);
    }

    button:hover {
        border-color: var(--text-secondary);
        background: var(--bg-secondary);
    }

    button.active {
        background: var(--accent-color);
        color: white;
        border-color: var(--accent-color);
    }

    /* Для темной темы - серовый активный цвет */
    [data-theme="dark"] button.active {
        background: #4a5568; /* Серый цвет для темной темы */
        color: white;
        border-color: #4a5568;
    }

    .icon {
        font-size: 16px;
        width: 20px;
        text-align: center;
    }

    .label {
        font-size: 13px;
        font-weight: 500;
    }

    @media (max-width: 768px) {
        .sidebar-filters {
            padding: 15px;
            gap: 20px;
        }

        button {
            padding: 12px 15px;
        }

        .label {
            font-size: 14px;
        }
    }
</style>