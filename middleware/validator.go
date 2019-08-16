package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

// Validator performs a json schema validation on an request body
func Validator(handler ProcessedHandlerFunc, schemaPath string) ProcessedHandlerFunc {

	loader := gojsonschema.NewReferenceLoader(schemaPath)
	schema, err := gojsonschema.NewSchema(loader)

	if err != nil {
		panic(fmt.Errorf("Could not load schema %v, %v", schemaPath, err))
	}

	return ProcessedHandlerFunc(func(w http.ResponseWriter, r *ProcessedRequest) {

		loader := gojsonschema.NewBytesLoader(r.JSON)
		result, err := schema.Validate(loader)

		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		// if the schema is valid then just run the required handler
		if result.Valid() {
			handler(w, r)
			return
		}

		// else handle the errors back to the user
		errors := result.Errors()

		http.Error(w, fmt.Sprintf("Invalid JSON Body: %v", errors), http.StatusBadRequest)
	})
}
