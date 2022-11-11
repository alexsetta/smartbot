package main

import (
	"fmt"
	"github.com/alexsetta/smartbot/cfg"
	"github.com/alexsetta/smartbot/cotacao"
	"net/http"
	"time"
)

func Total(w http.ResponseWriter, r *http.Request) {
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

	atual := 0.0
	for _, atv := range carteira.Ativos {
		if atv.Tipo != "criptomoeda" {
			continue
		}

		_, _, out, err := cotacao.Calculo(atv, config, alerta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		atual += out.Atual
	}

	resposta := fmt.Sprintf(`{"hora": "%v","total": %v}`, time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05"), atual)
	w.Write([]byte(resposta))

}
