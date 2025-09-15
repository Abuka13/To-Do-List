declare interface Window {
    go: {
        app: {
            App: {
                GetAllTasks(): Promise<any[]>;
                CreateTask(title: string, priority: string, dueDate: string): Promise<any>;
                SearchTasks(title: string): Promise<any[]>;
                ToggleTaskCompletion(id: number): Promise<any>;
                DeleteTask(id: number): Promise<void>;
                GetStats(): Promise<Record<string, number>>;
                TestMethod(): Promise<string>;
                GetTasksByFilter(filter: string): Promise<any[]>; // ДОБАВЛЕНО
            }
        }
    }
}