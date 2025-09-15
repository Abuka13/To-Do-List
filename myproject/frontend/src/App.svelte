<script>
  import { onMount } from 'svelte';
  import AddTaskForm from './components/AddTaskForm.svelte';
  import TaskList from './components/TaskList.svelte';

  let tasks = [];
  let loading = true;

  onMount(async () => {
    // Ждем инициализации Wails
    await waitForWails();
    await loadTasks();
  });

  // Функция для ожидания инициализации Wails
  async function waitForWails() {
    const maxAttempts = 50;
    const delay = 100;

    for (let i = 0; i < maxAttempts; i++) {
      if (window.go && window.go.main && window.go.main.App) {
        return true;
      }
      await new Promise(resolve => setTimeout(resolve, delay));
    }

    throw new Error('Wails runtime not available');
  }

  async function loadTasks() {
    try {
      loading = true;
      tasks = await window.go.main.App.GetAllTasks();
    } catch (error) {
      console.error('Ошибка загрузки задач:', error);
    } finally {
      loading = false;
    }
  }
</script>

<main>
  <h1>To-Do List</h1>

  <!-- Добавляем нашу форму -->
  <AddTaskForm onTaskAdded={loadTasks} />

  <!-- Здесь будет список задач -->
  {#if loading}
    <p>Загрузка...</p>
  {:else}
    <TaskList {tasks} onTaskUpdated={loadTasks} />
  {/if}
</main>

<style>
  main {
    padding: 20px;
    max-width: 800px;
    margin: 0 auto;
  }

  h1 {
    text-align: center;
    margin-bottom: 30px;
    color: #333;
  }
</style>