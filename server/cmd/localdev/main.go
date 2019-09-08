package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/rhythm/server/internal/api"
)

const (
	port         = "8000"
	staticPrefix = "/static/"
)

var (
	webuiBuildRoot string
)

func init() {
	kingpin.Flag("webui-build-root", "").
		Short('b').Required().StringVar(&webuiBuildRoot)
}

func main() {
	kingpin.Parse()

	cfg := &api.Config{
		Secret: "Spoon!",
	}

	router := localdevRouter()
	router.Mount("/api", api.New(cfg))

	fmt.Printf("Listening to http://localhost:%s/\n", port)
	err := http.ListenAndServe("127.0.0.1:"+port, router)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func localdevRouter() *chi.Mux {
	r := chi.NewRouter()

	fs := http.StripPrefix(staticPrefix,
		http.FileServer(http.Dir(filepath.Join(webuiBuildRoot, "static"))),
	)
	r.Get(staticPrefix+"*", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}),
	)

	index := filepath.Join(webuiBuildRoot, "index.html")
	r.Get("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, index)
		}),
	)

	fs = http.StripPrefix("/",
		http.FileServer(http.Dir(webuiBuildRoot)),
	)
	r.Get("/{:asset-manifest.json|service-worker.js}", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}),
	)

	return r
}
