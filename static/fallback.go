package static

import (
	"net/http"
	"os"
)

type fallback struct {
	defaultPath string
	fs          http.FileSystem
}

func (fb fallback) Open(path string) (http.File, error) {
	f, err := fb.fs.Open(path)
	if os.IsNotExist(err) {
		return fb.fs.Open(fb.defaultPath)
	}
	return f, err
}
