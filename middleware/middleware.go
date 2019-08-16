package middleware

import (
	"io/ioutil"
	"net/http"
)

// ProcessedRequest adds the validated json object to the incoming request body
type ProcessedRequest struct {
	RawRequest *http.Request

	// TODO: Explore a better solution for this, ideally we would pass an
	// unmarshalled object here so the handler didn't have to unmarshal itself
	JSON []byte
}

// ProcessedHandlerFunc function type for handlers
type ProcessedHandlerFunc func(http.ResponseWriter, *ProcessedRequest)

// ProcessRequest reads in the json body of the request if it exists
func ProcessRequest(handler ProcessedHandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jsonBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Could not parse request body", http.StatusBadRequest)
			return
		}

		handler(w, &ProcessedRequest{
			RawRequest: r,
			JSON:       jsonBytes,
		})
	})
}
