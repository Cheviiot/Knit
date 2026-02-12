package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

// Image proxy client
var imageProxyClient = &http.Client{Timeout: 30 * time.Second}

func init() {
	// Force WebKit to use system cursor on Linux
	os.Setenv("WEBKIT_DISABLE_COMPOSITING_MODE", "1")
}

func main() {
	// Убираем префикс frontend/dist из пути к файлам
	fsRoot, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	// Create handler that proxies /tmdb-image/* requests
	fileServer := http.FileServerFS(fsRoot)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Proxy TMDB images
		if strings.HasPrefix(r.URL.Path, "/tmdb-image/") {
			imagePath := strings.TrimPrefix(r.URL.Path, "/tmdb-image")
			tmdbURL := "https://imagetmdb.com/t/p" + imagePath

			req, err := http.NewRequest("GET", tmdbURL, nil)
			if err != nil {
				http.Error(w, "Failed to create request", http.StatusInternalServerError)
				return
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

			resp, err := imageProxyClient.Do(req)
			if err != nil {
				http.Error(w, "Failed to fetch image", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			// Copy headers
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	app := application.New(application.Options{
		Name:        "Knit",
		Description: "Movie torrent search application",
		Services: []application.Service{
			application.NewService(NewApp()),
		},
		Assets: application.AssetOptions{
			Handler: handler,
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
