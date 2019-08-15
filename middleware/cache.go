package middleware

import (
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
)

// Cache implements a memcache layer for storing processed images
func Cache(handler func(http.ResponseWriter, *http.Request), active bool, cache *memcache.Client) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if active {
			key := r.URL.String()
			item, err := cache.Get(key)

			if err == memcache.ErrCacheMiss {
				handler(w, r)
				return
			}

			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Error contacting cache", http.StatusInternalServerError)
				return
			}

			w.Write(item.Value)
			return
		}

		handler(w, r)
		return
	})
}
