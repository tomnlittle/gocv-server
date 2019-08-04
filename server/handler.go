package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomnlittle/gocv-server/middleware"
	"gocv.io/x/gocv"
)

// Handler for routes
type Handler struct {
	AwsConfig *AwsConfig
}

// NewHandler returns an initialised handler
func NewHandler() *Handler {

	awsConfig, err := NewAwsConfig()

	if err != nil {
		panic(err.Error())
	}

	return &Handler{
		AwsConfig: awsConfig,
	}
}

func getImage(r *http.Request, h *Handler) ([]byte, error) {
	pathParameters := mux.Vars(r)
	bucket := pathParameters["bucket"]
	id := pathParameters["key"]

	buffer, err := h.AwsConfig.GetObject(bucket, id)
	if err != nil {
		return nil, fmt.Errorf("Could not get %v from bucket %v", id, bucket)
	}

	return buffer, nil
}

// ------------------------------------------------ Simple Handler -----------------------------------------------------

// Simple .
func (h *Handler) Simple(w http.ResponseWriter, r *http.Request) {

	buffer, err := getImage(r, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mat, err := gocv.IMDecode(buffer, gocv.IMReadColor)

	if err != nil {
		http.Error(w, "Server Error when decoding image", http.StatusInternalServerError)
		return
	}

	format := r.FormValue("format")
	quality := r.FormValue("quality")

	rBuf, err := EncodeMatrix(mat, format, quality)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(rBuf)
}

// ------------------------------------------------ Complex Handler ----------------------------------------------------

// ComplexHandlerJSON used to unmarshal
type ComplexHandlerJSON struct {
	Encoding  string `json:"encoding"`
	Quality   int    `json:"quality"`
	Functions []struct {
		ID         string                     `json:"id"`
		Parameters []complexHandlerParameters `json:"parameters"`
	} `json:"functions"`
}

type complexHandlerParameters struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Complex handler takes json objects as input and 'curries' the matrix through
// the desired functions outlined in the json object
func (h *Handler) Complex(w http.ResponseWriter, r *middleware.ValidatedRequest) {

	var parsedBody ComplexHandlerJSON
	err := json.Unmarshal(r.ValidatedJSON, &parsedBody)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	buffer, err := getImage(r.RawRequest, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mat, err := gocv.IMDecode(buffer, gocv.IMReadColor)

	for _, v := range parsedBody.Functions {
		cvFunc, ok := FunctionMappings[v.ID]

		if !ok {
			http.Error(w, fmt.Sprintf("Invalid function id %v", v.ID), http.StatusBadRequest)
			return
		}

		params := buildParameterMapping(v.Parameters)

		rMat, err := cvFunc(mat, params)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mat = *rMat
	}

	// TODO parsed body string
	quality := fmt.Sprintf("%v", parsedBody.Quality)
	rBuf, err := EncodeMatrix(mat, parsedBody.Encoding, quality)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(rBuf)
}

func buildParameterMapping(parameters []complexHandlerParameters) map[string]string {
	mapping := make(map[string]string)
	for _, v := range parameters {
		mapping[v.Key] = v.Value
	}

	return mapping
}
