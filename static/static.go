package static

import (
	"net/http"
)

type Options struct {
	Path      string
	Fallback  bool
	BasicAuth string
	MaxAge    int
}

func (s *server) handleStatic() logWriterFunc {
	var fileSystem http.FileSystem = http.Dir(s.staticOptions.Path)

	if s.staticOptions.Fallback {
		fileSystem = fallback{
			defaultPath: "index.html",
			fs:          fileSystem,
		}
	}

	handler := http.FileServer(fileSystem).ServeHTTP

	// Handle headers middleware
	addCacheControlValue(s.staticOptions.MaxAge)
	handler = addHeaders(handler)

	return func(w *logWriter, r *http.Request) {
		handler(w, r)
	}
}
