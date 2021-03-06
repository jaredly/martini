package martini

import (
	"log"
	"path/filepath"
	"net/http"
)

// Static returns a middleware handler that serves static files in the given path.
func Static(path string) Handler {
	dir := http.Dir(path)
	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		file := req.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// discard the error?
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return
		}
		if fi.IsDir() {
			file = filepath.Join(file, "index.html")
			f, err = dir.Open(file)
			if err != nil {
				return
			}
			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				return
			}
		}

		log.Println("[Static] Serving " + file)
		http.ServeContent(res, req, file, fi.ModTime(), f)
	}
}
