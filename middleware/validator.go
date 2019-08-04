package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

// ValidatedRequest adds the validated json object to the incoming request body
type ValidatedRequest struct {
	RawRequest *http.Request

	// TODO: Explore a better solution for this, ideally we would pass an
	// unmarshalled object here so the handler didn't have to unmarshal itself
	ValidatedJSON []byte
}

// ValidatedHandlerFunc function type for handlers
type ValidatedHandlerFunc func(http.ResponseWriter, *ValidatedRequest)

// Validator performs a json schema validation on an request body
func Validator(handler ValidatedHandlerFunc, schemaPath string) http.HandlerFunc {

	loader := gojsonschema.NewReferenceLoader(schemaPath)
	schema, err := gojsonschema.NewSchema(loader)

	if err != nil {
		panic(fmt.Errorf("Could not load schema %v, %v", schemaPath, err))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jsonBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		loader := gojsonschema.NewBytesLoader(jsonBytes)
		result, err := schema.Validate(loader)

		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		// if the schema is valid then just run the required handler
		if result.Valid() {

			validated := ValidatedRequest{
				RawRequest:    r,
				ValidatedJSON: jsonBytes,
			}

			handler(w, &validated)
			return
		}

		// else handle the errors back to the user
		errors := result.Errors()

		http.Error(w, fmt.Sprintf("Invalid JSON Body: %v", errors), http.StatusBadRequest)
	})
}
