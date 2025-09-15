<script>
    import {onMount} from 'svelte';
    import AddTaskForm from './components/AddTaskForm.svelte';
    import TaskList from './components/TaskList.svelte';

    let tasks = [];
    let loading = true;
    let wailsReady = false;
    let errorMessage = '';

    onMount(async () => {
        console.log('App mounted');
        await initializeApp();
    });

    async function initializeApp() {
        console.log('Инициализация приложения...');

        const waitForWails = async () => {
            let attempts = 0;
            const maxAttempts = 100;

            while (attempts < maxAttempts) {
                // ИСПРАВЛЕННЫЙ ПУТЬ!
                if (window.go?.app?.App?.GetAllTasks) {
                    console.log('Wails runtime готов!');
                    wailsReady = true;
                    await loadTasks();
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

    async function loadTasks() {
        try {
            console.log('Загрузка задач...');
            loading = true;
            errorMessage = '';

            // ИСПРАВЛЕННЫЙ ПУТЬ!
            if (!window.go?.app?.App?.GetAllTasks) {
                throw new Error('Wails runtime недоступен');
            }

            const loadedTasks = await window.go.app.App.GetAllTasks();
            tasks = loadedTasks || [];
            console.log('Задачи загружены:', tasks);

        } catch (error) {
            console.error('Ошибка загрузки задач:', error);
            errorMessage = `Ошибка загрузки: ${error.message}`;
            tasks = [];
        } finally {
            loading = false;
        }
    }
</script>

<main>
    <h1>To-Do List</h1>

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
        <AddTaskForm onTaskAdded={loadTasks}/>

        {#if loading}
            <div class="loading">
                <p>Загрузка задач...</p>
            </div>
        {:else}
            <TaskList {tasks} onTaskUpdated={loadTasks}/>
        {/if}
    {/if}
</main>

<style>
    main {
        padding: 20px;
        max-width: 800px;
        margin: 0 auto;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    }

    h1 {
        text-align: center;
        margin-bottom: 30px;
        color: #333;
    }

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
        color: #666;
        padding: 20px;
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
        padding: 8px 16px;
        border-radius: 4px;
        cursor: pointer;
        margin-left: 10px;
        font-size: 14px;
    }

    .retry-btn:hover {
        background: #c82333;
    }
</style>