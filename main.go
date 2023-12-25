package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"sync"
)

type clipEntry struct {
	Label   string
	Content string
}

var store struct {
	entries []clipEntry
	sync.RWMutex
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Fprintln(os.Stderr, "starting server on "+port)

	startServer(port)
}

func startServer(port string) {
	http.Handle("/simple.css", staticWrite(css, ".css"))
	http.Handle("/", http.HandlerFunc(handle))

	http.ListenAndServe(":"+port, nil)
}

func staticWrite(content []byte, ext string) http.HandlerFunc {
	var mimeType string
	if ext != "" {
		mimeType = mime.TypeByExtension(ext)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if mimeType != "" {
			w.Header().Add("Content-type", mimeType)
		}
		w.Write(content)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "error occured: %v", err)
			fmt.Fprintln(os.Stderr)
		}
	}()

	if r.URL.Path == "/save" && r.Method == http.MethodPost {
		save(r.FormValue("label"), r.FormValue("content"))
	} else if r.URL.Path == "/clear" && r.Method == http.MethodPost {
		clear()
	}

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	t := template.Must(template.New("").Parse(html))
	if err := t.Execute(w, store.entries); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error writing html: %w", err))
	}
}

func save(label, content string) {
	store.Lock()
	defer store.Unlock()

	entry := clipEntry{Label: label, Content: content}
	store.entries = append(store.entries, entry)
}

func clear() {
	store.Lock()
	defer store.Unlock()

	store.entries = nil
}

//go:embed simple.css
var css []byte

//go:embed index.html
var html string
