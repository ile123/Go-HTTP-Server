package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func PrintURLQueryValues(ctx context.Context, r *http.Request) {
	url_query_values := r.URL.Query()
	fmt.Println("\n\nURL Querry values(if any):")
	for k, v := range url_query_values {
		fmt.Println(k, " => ", v)
	}
}

func PrintFormValues(r *http.Request) {
	var err error
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		//Needed to parse multi-part form, 32MB support
		err = r.ParseMultipartForm(32 << 20)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		fmt.Println("Unable to parse form:", err)
		return
	}
	for k, v := range r.Form {
		fmt.Println(k, " => ", v)
	}
}

func PrintURLBodyValues(r io.Reader) {
	body, err := io.ReadAll(r)
	if err != nil {
		fmt.Printf("Could not read request body: %s\n", err)
		return
	}
	fmt.Println("\n\nBody content (if any):")
	fmt.Println(string(body))
}
