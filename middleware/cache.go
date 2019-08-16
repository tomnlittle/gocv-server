package middleware

import (
	"net/http"

	"github.com/tomnlittle/gocv-server/cache"
)

// Cache implements a memcache layer for storing processed images
func Cache(handler ProcessedHandlerFunc, mc *cache.ImageCache) ProcessedHandlerFunc {

	// namespace, err := uuid.NewV4()
	// if err != nil {
	// 	panic(err)
	// }

	return ProcessedHandlerFunc(func(w http.ResponseWriter, r *ProcessedRequest) {

		// hash := mc.GenerateHash(namespace, r.RawRequest.URL.String(), r.RawRequest.Method, string(r.JSON))
		// bytes, err := mc.GetBytes(hash)

		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, "Cache error", http.StatusInternalServerError)
		// 	return
		// }

		// if bytes != nil && len(bytes) > 0 {
		// 	w.Write(bytes)
		// 	return
		// }

		handler(w, r)
	})
}
