package static

import (
	"net/http"
)

type Options struct {
	Path         string
	FallbackPath string
	BasicAuth    string
	MaxAge       int
}

func (s *server) handleStatic() logWriterFunc {
	var fileSystem http.FileSystem = http.Dir(s.staticOptions.Path)

	if len(s.staticOptions.FallbackPath) != 0 {
		fileSystem = fallback{
			defaultPath: s.staticOptions.Path + "/" + s.staticOptions.FallbackPath,
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
