package main

import (
	"context"
	"embed"
	"fmt"
	app2 "myproject/pkg/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Создаем приложение напрямую
	app := app2.NewApp()

	err := wails.Run(&options.App{
		Title:     "To-Do List",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 248, G: 249, B: 250, A: 1},
		OnStartup: func(ctx context.Context) {
			fmt.Println("Приложение ")
		},
		Bind: []interface{}{
			app,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})

	if err != nil {
		fmt.Printf("Ошибка запуска приложения: %s\n", err.Error())
	}
}
