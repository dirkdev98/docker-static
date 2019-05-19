package static

import "log"

func errorLogger(_ int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
