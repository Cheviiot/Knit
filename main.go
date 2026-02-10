package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Убираем префикс frontend/dist из пути к файлам
	fsRoot, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	app := application.New(application.Options{
		Name:        "Knit",
		Description: "Movie torrent search application",
		Services: []application.Service{
			application.NewService(NewApp()),
		},
		Assets: application.AssetOptions{
			Handler: http.FileServerFS(fsRoot),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Knit",
		Width:            1200,
		Height:           800,
		MinWidth:         800,
		MinHeight:        600,
		Frameless:        true,
		BackgroundColour: application.NewRGBA(13, 13, 13, 255),
		URL:              "/",
		Linux: application.LinuxWindow{
			WindowIsTranslucent: false,
		},
	})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
