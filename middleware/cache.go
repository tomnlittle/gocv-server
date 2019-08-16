package middleware

import (
	"log"
	"net/http"

	"github.com/tomnlittle/gocv-server/cache"
)

// Cache implements a memcache layer for storing processed images
func Cache(handler ProcessedHandlerFunc, mc *cache.ImageCache) ProcessedHandlerFunc {

	return ProcessedHandlerFunc(func(w http.ResponseWriter, r *ProcessedRequest) {

		hash := mc.GenerateHash(r.RawRequest.URL.String(), r.RawRequest.Method, string(r.JSON))
		bytes, err := mc.GetBytes(hash)

		if err != nil {
			log.Println(err)
			http.Error(w, "Cache error", http.StatusInternalServerError)
			return
		}

		if bytes == nil || len(bytes) == 0 {
			handler(w, r)
			return
		}

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error contacting cache", http.StatusInternalServerError)
			return
		}

		w.Write(bytes)
		return
	})
}
