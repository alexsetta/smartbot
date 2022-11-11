package main

import (
	"fmt"
	"github.com/alexsetta/smartbot/cfg"
	"github.com/alexsetta/smartbot/cotacao"
	"github.com/alexsetta/smartbot/util"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type AssetOrder struct {
	qty   float64
	value float64
}

func Order(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["from"]
	if from == "" {
		http.Error(w, "Argument 'from' not found", http.StatusInternalServerError)
		return
	}
	from = strings.ToUpper(from)

	to := vars["to"]
	if to == "" {
		http.Error(w, "Argument 'to' not found", http.StatusInternalServerError)
		return
	}
	to = strings.ToUpper(to)

	if err := cfg.ReadConfig("../smartbot.cfg", &config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	if err := cfg.ReadConfig("../carteira.cfg", &carteira); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fromAsset := AssetOrder{}
	toAsset := AssetOrder{}

	atual := 0.0
	for _, atv := range carteira.Ativos {
		if atv.Tipo != "criptomoeda" {
			continue
		}

		if atv.Simbolo != from && atv.Simbolo != to {
			continue
		}

		_, _, out, err := cotacao.Calculo(atv, config, alerta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if atv.Simbolo == from {
			fromAsset.qty = atv.Quantidade
			fromAsset.value = out.Preco
		} else {
			toAsset.qty = atv.Quantidade
			toAsset.value = out.Preco
		}

		atual += out.Atual
	}

	resposta := fmt.Sprintf(`{"hora": "%v","from": "%v","to":"%v}`, util.Now(), fromAsset, toAsset)
	_, _ = w.Write([]byte(resposta))
}
