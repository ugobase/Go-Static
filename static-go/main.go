package main

import (
	"context"
    "embed"
    "io/fs"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

//go:embed static/*.html
var embeddedFiles embed.FS

// servePage returns a handler that serves the given HTML file
func servePage(name string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := fs.ReadFile(embeddedFiles, "static/"+name+".html")
        if err != nil {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/html")
        w.Write(data)
    }
}

func main() {
    mux := http.NewServeMux()

    // Map endpoints to HTML files
    mux.HandleFunc("/home", servePage("home"))
    mux.HandleFunc("/about", servePage("about"))
    mux.HandleFunc("/contact", servePage("contact"))

    // Health endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

    srv := &http.Server{Addr: ":8080", Handler: mux}

    go func() {
        log.Println("Server listening on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen failed: %v", err)
        }
    }()

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
    log.Println("Server stopped")
}
