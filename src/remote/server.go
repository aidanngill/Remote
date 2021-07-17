package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func FileGetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	entry, err := findPackageById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	file, err := os.Open(entry.Path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	defer file.Close()

	http.ServeFile(w, r, entry.Path)
}
