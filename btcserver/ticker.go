package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Ticker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	asset, err := NewAsset(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resposta := ""
	if id == "all" {
		data, err := asset.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resposta, _ = PrettyJson(data)
		_, _ = w.Write([]byte(resposta))
		return
	}

	err = asset.Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prettyJSON, _ := PrettyJson(asset.data)
	resposta = prettyJSON + ","
	resposta = "[" + resposta[:len(resposta)-1] + "]"
	_, _ = w.Write([]byte(resposta))
}
