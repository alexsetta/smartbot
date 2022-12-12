package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexsetta/smartbot/rsi"
	"strings"
	"time"

	"github.com/alexsetta/smartbot/cfg"
	"github.com/alexsetta/smartbot/cotacao"
	"github.com/alexsetta/smartbot/tipos"
)

var (
	hora   = time.Now().Add(time.Hour * -5)
	alerta = tipos.Alertas{hora, hora, hora, hora, hora, hora}
)

func PrettyJson(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func Cotacao(id string) string {
	if err := cfg.ReadConfig("../carteira.cfg", &carteira); err != nil {
		return fmt.Sprintf("cotacao: read carteira.cfg: %s", err)
	}

	mr := make(map[string]*rsi.RSI)
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo)
	}

	resposta := "["
	outJson := tipos.Result{}
	if id == "all" {
		for _, atv := range carteira.Ativos {
			_, _, out, err := cotacao.Calculo(atv, config, alerta, mr)
			if err != nil {
				return fmt.Sprintf("cotacao: calculo: %s", err)
			}
			outJson = out
			prettyJSON, _ := PrettyJson(out)
			resposta += prettyJSON + ","
		}
		resposta = resposta[:len(resposta)-1] + "]"
		return resposta
	}

	ativo := tipos.Ativo{}
	for _, atv := range carteira.Ativos {
		if strings.ToLower(atv.Simbolo) == id {
			ativo = atv
			break
		}
	}
	var err2 error
	resposta, _, outJson, err2 = cotacao.Calculo(ativo, config, alerta, mr)
	if err2 != nil {
		return fmt.Sprintf("cotacao: calculo[2]: %s", err2)
	}
	// remover as linhas abaixo para mostrar como "string"
	outJson = outJson
	prettyJSON, _ := PrettyJson(outJson)
	resposta = prettyJSON + ","
	return resposta
}

func Total() string {
	if err := cfg.ReadConfig("../carteira.cfg", &carteira); err != nil {
		return fmt.Sprintf("cotacao: read carteira.cfg: %s", err)
	}

	mr := make(map[string]*rsi.RSI)
	atual := 0.0
	for _, atv := range carteira.Ativos {
		if atv.Tipo != "criptomoeda" {
			continue
		}
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo)

		_, _, out, err := cotacao.Calculo(atv, config, alerta, mr)
		if err != nil {
			return fmt.Sprintf("cotacao: calclulo: %s", err)
		}
		atual += out.Atual
	}

	return fmt.Sprintf(`{"hora": "%v","total": %v}`, time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05"), atual)
}
