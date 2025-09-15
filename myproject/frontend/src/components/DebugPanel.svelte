<script>
    // Убираем экспорт, так как мы используем bind:
    let messages = [];
    let isVisible = false;
</script>

{#if isVisible}
    <div class="debug-panel">
        <h3>🐛 Debug Panel</h3>
        <div class="messages">
            {#each messages as message}
                <div class="message {message.type}">
                    [{message.timestamp}] {message.text}
                </div>
            {/each}
        </div>
        <button on:click={() => isVisible = false}>Закрыть</button>
    </div>
{/if}

<style>
    .debug-panel {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        background: white;
        border: 2px solid #007bff;
        border-radius: 8px;
        padding: 20px;
        z-index: 1000;
        max-width: 500px;
        max-height: 400px;
        overflow-y: auto;
    }

    .messages {
        margin: 10px 0;
    }

    .message {
        padding: 5px;
        margin: 2px 0;
        border-radius: 4px;
    }

    .message.error {
        background: #ffe6e6;
        color: #d00;
    }

    .message.info {
        background: #e6f7ff;
        color: #007bff;
    }

    .message.success {
        background: #e6ffe6;
        color: #090;
    }
</style>