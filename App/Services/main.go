package main

import (
	"fmt"
	"net/http"
)

func main() {
	// TODO: Change the Service to echo
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Wellcome to the base project")
	})
	http.ListenAndServe(":8080", nil)
}
