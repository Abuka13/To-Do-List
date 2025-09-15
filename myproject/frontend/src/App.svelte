<script>
    import {onMount} from 'svelte';
    import AddTaskForm from './components/AddTaskForm.svelte';
    import TaskList from './components/TaskList.svelte';
    import SidebarFilters from './components/SidebarFilters.svelte';

    let tasks = [];
    let allTasks = [];
    let loading = true;
    let wailsReady = false;
    let errorMessage = '';
    let currentFilter = 'all';
    let dateFilter = 'all';
    let sortBy = 'date';
    let sidebarOpen = true;
    let darkMode = false;

    function toggleTheme() {
        darkMode = !darkMode;
        const theme = darkMode ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('theme', theme);
        console.log('Тема изменена на:', theme);
    }

    onMount(() => {
        // Загрузка темы из localStorage
        const savedTheme = localStorage.getItem('theme');
        darkMode = savedTheme === 'dark';
        document.documentElement.setAttribute('data-theme', darkMode ? 'dark' : 'light');
        console.log('Тема загружена:', darkMode ? 'dark' : 'light');

        initializeApp();
    });

    function toggleSidebar() {
        sidebarOpen = !sidebarOpen;
        console.log('Сайдбар:', sidebarOpen ? 'открыт' : 'закрыт');
    }

    async function initializeApp() {
        console.log('Инициализация приложения...');

        const waitForWails = async () => {
            let attempts = 0;
            const maxAttempts = 100;

            while (attempts < maxAttempts) {
                if (window.go?.app?.App?.GetAllTasks) {
                    console.log('Wails runtime готов!');
                    wailsReady = true;
                    await loadAllTasks();
                    return true;
                }

                console.log(`Ожидание Wails runtime... попытка ${attempts + 1}/${maxAttempts}`);
                await new Promise(resolve => setTimeout(resolve, 100));
                attempts++;
            }

            console.error('Не удалось подключиться к Wails runtime');
            errorMessage = 'Не удалось подключиться к backend приложения';
            loading = false;
            return false;
        };

        await waitForWails();
    }

    async function loadAllTasks() {
        try {
            console.log('Загрузка всех задач...');
            loading = true;
            errorMessage = '';

            if (!window.go?.app?.App?.GetAllTasks) {
                throw new Error('Wails runtime недоступен');
            }

            allTasks = await window.go.app.App.GetAllTasks();
            applyFilters();
            console.log('Все задачи загружены:', allTasks.length);

        } catch (error) {
            console.error('Ошибка загрузки задач:', error);
            errorMessage = `Ошибка загрузки: ${error.message}`;
            allTasks = [];
            tasks = [];
        } finally {
            loading = false;
        }
    }

    function applyFilters() {
        let filtered = [...allTasks];

        switch (currentFilter) {
            case 'active':
                filtered = filtered.filter(task => !task.isCompleted);
                break;
            case 'completed':
                filtered = filtered.filter(task => task.isCompleted);
                break;
        }

        const now = new Date();
        const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const weekStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 7);

        switch (dateFilter) {
            case 'today':
                filtered = filtered.filter(task => {
                    if (!task.dueDate) return false;
                    const taskDate = new Date(task.dueDate);
                    return taskDate >= todayStart && !task.isCompleted;
                });
                break;
            case 'week':
                filtered = filtered.filter(task => {
                    if (!task.dueDate) return false;
                    const taskDate = new Date(task.dueDate);
                    return taskDate >= weekStart && !task.isCompleted;
                });
                break;
            case 'overdue':
                filtered = filtered.filter(task => {
                    if (!task.dueDate || task.isCompleted) return false;
                    const taskDate = new Date(task.dueDate);
                    return taskDate < todayStart;
                });
                break;
        }

        switch (sortBy) {
            case 'priority':
                filtered.sort((a, b) => {
                    const priorityOrder = {high: 3, medium: 2, low: 1};
                    return (priorityOrder[b.priority] || 0) - (priorityOrder[a.priority] || 0);
                });
                break;
            case 'date':
            default:
                filtered.sort((a, b) => {
                    const dateA = new Date(a.createdAt).getTime();
                    const dateB = new Date(b.createdAt).getTime();
                    return dateB - dateA;
                });
                break;
        }

        tasks = filtered;
    }

    function handleFilterChange(event) {
        currentFilter = event.detail.statusFilter;
        dateFilter = event.detail.dateFilter;
        sortBy = event.detail.sortBy;
        applyFilters();
    }

    async function refreshTasks() {
        await loadAllTasks();
    }
</script>

<main class:sidebar-open={sidebarOpen}>
    {#if !wailsReady && !errorMessage}
        <div class="warning">
            Инициализация приложения...
        </div>
    {/if}

    {#if errorMessage}
        <div class="error">
            {errorMessage}
            <button on:click={initializeApp} class="retry-btn">
                Попробовать снова
            </button>
        </div>
    {/if}

    {#if wailsReady}
        <!-- Кнопка для открытия сайдбара - всегда видимая -->
        {#if !sidebarOpen}
            <button class="sidebar-toggle-btn" on:click={toggleSidebar} title="Открыть фильтры">
                ☰
            </button>
        {/if}

        <aside class:open={sidebarOpen}>
            <div class="sidebar-header">
                <h3>Фильтры</h3>
                <button class="toggle-sidebar" on:click={toggleSidebar}>
                    {sidebarOpen ? '◀' : '▶'}
                </button>
            </div>

            {#if sidebarOpen}
                <SidebarFilters
                        {currentFilter}
                        {dateFilter}
                        {sortBy}
                        on:filterchange={handleFilterChange}
                />
            {/if}
        </aside>

        <div class="content">
            <header>
                <h1>To-Do List</h1>
                <div class="header-actions">
                    <button class="theme-btn" on:click={toggleTheme} title="Сменить тему">
                        {darkMode ? '☀️' : '🌙'}
                    </button>
                </div>
            </header>

            <AddTaskForm onTaskAdded={refreshTasks}/>

            {#if loading}
                <div class="loading">
                    <p>Загрузка задач...</p>
                </div>
            {:else}
                <TaskList {tasks} onTaskUpdated={refreshTasks}/>

                <div class="stats">
                    Всего: {allTasks.length} |
                    Активные: {allTasks.filter(t => !t.isCompleted).length} |
                    Выполненные: {allTasks.filter(t => t.isCompleted).length} |
                    Просроченные: {allTasks.filter(t => {
                    if (!t.dueDate || t.isCompleted) return false;
                    const taskDate = new Date(t.dueDate);
                    const todayStart = new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate());
                    return taskDate < todayStart;
                }).length}
                </div>
            {/if}
        </div>
    {/if}
</main>

<style>
    /* Основной контейнер */
    main {
        display: flex;
        min-height: 100vh;
        background: var(--bg-secondary);
        position: relative;
    }

    /* Плавающая кнопка для открытия сайдбара */
    .sidebar-toggle-btn {
        position: fixed;
        top: 20px;
        left: 20px;
        z-index: 1000;
        background: var(--accent-color);
        color: white;
        border: none;
        border-radius: 50%;
        width: 50px;
        height: 50px;
        font-size: 18px;
        cursor: pointer;
        box-shadow: var(--shadow);
        transition: all 0.3s ease;
    }

    .sidebar-toggle-btn:hover {
        transform: scale(1.1);
        box-shadow: 0 4px 12px rgba(0, 123, 255, 0.3);
    }

    /* Сайдбар */
    aside {
        width: 280px;
        background: var(--bg-primary);
        border-right: 1px solid var(--border-color);
        transition: transform 0.3s ease;
        box-shadow: var(--shadow);
        position: relative;
        z-index: 100;
        flex-shrink: 0;
        height: 100vh;
        overflow-y: auto;
        transform: translateX(0);
    }

    /* Когда сайдбар закрыт */
    aside:not(.open) {
        transform: translateX(-100%);
    }

    /* Контентная область */
    .content {
        flex: 1;
        padding: 20px;
        min-width: 0;
        transition: margin-left 0.3s ease;
        height: 100vh;
        overflow-y: auto;
    }

    /* Заголовок сайдбара */
    .sidebar-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 20px;
        border-bottom: 1px solid var(--border-color);
        background: var(--bg-secondary);
        position: sticky;
        top: 0;
        z-index: 10;
    }

    .sidebar-header h3 {
        margin: 0;
        color: var(--text-primary);
        font-size: 18px;
        font-weight: 600;
    }

    .toggle-sidebar {
        background: none;
        border: 1px solid var(--border-color);
        border-radius: 4px;
        padding: 5px 10px;
        cursor: pointer;
        font-size: 14px;
        color: var(--text-primary);
        transition: all 0.2s ease;
    }

    .toggle-sidebar:hover {
        background: var(--bg-secondary);
        border-color: var(--text-secondary);
    }

    /* Заголовок контента */
    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 30px;
        padding-bottom: 15px;
        border-bottom: 2px solid var(--border-color);
        background: var(--bg-primary);
        position: sticky;
        top: 0;
        z-index: 50;
        padding: 20px;
        margin: -20px -20px 30px -20px;
        box-shadow: var(--shadow);
    }

    h1 {
        margin: 0;
        color: var(--text-primary);
        font-size: 28px;
        font-weight: 600;
    }

    .header-actions {
        display: flex;
        gap: 12px;
        align-items: center;
    }

    .theme-btn {
        background: var(--bg-secondary);
        border: 1px solid var(--border-color);
        border-radius: 8px;
        padding: 10px 14px;
        cursor: pointer;
        font-size: 16px;
        transition: all 0.3s ease;
        color: var(--text-primary);
    }

    .theme-btn:hover {
        background: var(--accent-color);
        color: white;
        border-color: var(--accent-color);
        transform: translateY(-1px);
    }

    /* Сообщения */
    .warning {
        background: #fff3cd;
        border: 1px solid #ffeaa7;
        color: #856404;
        padding: 15px;
        border-radius: 8px;
        margin-bottom: 20px;
        text-align: center;
        font-weight: 500;
    }

    .loading {
        text-align: center;
        color: var(--text-secondary);
        padding: 40px;
        font-style: italic;
    }

    .error {
        background: #f8d7da;
        border: 1px solid #f5c6cb;
        color: #721c24;
        padding: 15px;
        border-radius: 8px;
        text-align: center;
        font-weight: 500;
        margin-bottom: 20px;
    }

    .retry-btn {
        background: #dc3545;
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 6px;
        cursor: pointer;
        margin-top: 10px;
        font-size: 14px;
        font-weight: 500;
        transition: background-color 0.2s ease;
    }

    .retry-btn:hover {
        background: #c82333;
    }

    /* Статистика */
    .stats {
        margin-top: 20px;
        padding: 16px;
        background: var(--bg-primary);
        border-radius: 8px;
        text-align: center;
        font-size: 14px;
        color: var(--text-secondary);
        border: 1px solid var(--border-color);
        font-weight: 500;
    }

    /* Адаптивность для мобильных */
    @media (max-width: 768px) {
        .sidebar-toggle-btn {
            top: 10px;
            left: 10px;
            width: 45px;
            height: 45px;
            font-size: 16px;
        }

        aside {
            position: fixed;
            left: 0;
            top: 0;
            height: 100%;
            z-index: 1001;
        }

        .content {
            padding: 15px;
            width: 100%;
        }

        header {
            padding: 15px;
            margin: -15px -15px 20px -15px;
        }

        h1 {
            font-size: 24px;
        }
    }

    /* Анимации */
    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateX(-20px);
        }
        to {
            opacity: 1;
            transform: translateX(0);
        }
    }

    .content > * {
        animation: slideIn 0.3s ease-out;
    }
</style>