package main

import (
	"fmt"
	"github.com/alexsetta/smartbot/util"
	"github.com/gorilla/mux"
	"net/http"
)

//type AssetOrder struct {
//	qty   float64
//	value float64
//}

func Order(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["from"]
	if from == "" {
		http.Error(w, "Argument 'from' not found", http.StatusInternalServerError)
		return
	}

	to := vars["to"]
	if to == "" {
		http.Error(w, "Argument 'to' not found", http.StatusInternalServerError)
		return
	}

	fromAsset, err := NewAsset(from)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toAsset, err := NewAsset(to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resposta := fmt.Sprintf(`{"hora": "%v","from": "%v","to":"%v}`, util.Now(), fromAsset, toAsset)
	_, _ = w.Write([]byte(resposta))
}
