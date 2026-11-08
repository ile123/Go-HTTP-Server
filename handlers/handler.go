package handlers

import (
	"encoding/json"
	"fmt"
	"ile123/http-server/constants"
	"ile123/http-server/helpers"
	"io"
	"net/http"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: got / request.", ctx.Value(constants.KeyServerAddr))
	helpers.PrintURLQueryValues(ctx, r)
	helpers.PrintURLBodyValues(r.Body)
	sendResponse(w, 200, "This is my website!\n")
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	helpers.PrintURLQueryValues(ctx, r)
	//helpers.PrintURLBodyValues(r.Body)
	helpers.PrintFormValues(r)
	fmt.Printf("%s: got /hello request\n", ctx.Value(constants.KeyServerAddr))
	sendResponse(w, 200, "Hello, HTTP!\n")
}

func HandleJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		sendResponse(w, http.StatusMethodNotAllowed, "method not allowed\n")
		return
	}

	var msg constants.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		sendResponse(w, http.StatusBadRequest, "invalid JSON payload\n",
			constants.ResponseError{Kind: constants.ErrInvalidValue, Value: "body"})
		return
	}

	var errs []constants.ResponseError
	if msg.Name == "" {
		errs = append(errs, constants.ResponseError{Kind: constants.ErrMissingField, Value: "name"})
	}
	if msg.Body == "" {
		errs = append(errs, constants.ResponseError{Kind: constants.ErrMissingField, Value: "body"})
	}
	if len(errs) > 0 {
		sendResponse(w, http.StatusBadRequest, "", errs...)
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Received message from " + msg.Name,
		"body":    msg.Body,
	})
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"healthy"}`))
}

func sendResponse(w http.ResponseWriter, status int, body string, errs ...constants.ResponseError) {
	for _, e := range errs {
		w.Header().Add(string(e.Kind), e.Value)
	}

	w.WriteHeader(status)

	if body != "" {
		io.WriteString(w, body)
	}
}

func sendJSONResponse(w http.ResponseWriter, status int, payload any, errs ...constants.ResponseError) {
	data, err := json.Marshal(payload)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "could not encode response\n")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	sendResponse(w, status, string(data)+"\n", errs...)
}
