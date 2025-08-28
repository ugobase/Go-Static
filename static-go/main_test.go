package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
)

//go:embed static/*.html
var pages embed.FS

func servePage(name string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := fs.ReadFile(pages, "static/"+name+".html")
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
    mux.HandleFunc("/home", servePage("home"))
    mux.HandleFunc("/about", servePage("about"))
    mux.HandleFunc("/contact", servePage("contact"))
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", mux)
}
