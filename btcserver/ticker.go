package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Ticker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// desabilita mensagens no Telegram
	config.TelegramID = 0

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
		w.Write([]byte(resposta))
		return
	}

	outJson, err := asset.Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prettyJSON, _ := PrettyJson(outJson)
	resposta = prettyJSON + ","
	resposta = "[" + resposta[:len(resposta)-1] + "]"
	w.Write([]byte(resposta))
}
