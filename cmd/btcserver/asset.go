package main

import (
	"errors"
	"fmt"
	"github.com/alexsetta/smartbot/cotacao"
	"github.com/alexsetta/smartbot/rsi"
	"github.com/alexsetta/smartbot/tipos"
	"strings"
)

type Asset struct {
	id   string
	data tipos.Result
}

func NewAsset(id string) (*Asset, error) {
	asset := &Asset{id: id, data: tipos.Result{}}

	if err := asset.IsValid(); err != nil {
		return nil, err
	}

	if id == "all" {
		return asset, nil
	}

	if err := asset.Find(); err != nil {
		return nil, err
	}
	return asset, nil
}

func (a *Asset) IsValid() error {
	if a.id == "" {
		return errors.New("id is empty")
	}
	return nil
}

func (a *Asset) Find() error {
	if err := ReadConfig(); err != nil {
		return err
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	// desabilita saveLog
	config.SaveLog = false
	mr := make(map[string]*rsi.RSI)

	ativo := tipos.Ativo{}
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo)
		if strings.ToLower(atv.Simbolo) == a.id {
			ativo = atv
			break
		}
	}

	if ativo == (tipos.Ativo{}) {
		return fmt.Errorf("asset %s not found", a.id)
	}

	_, _, out, err := cotacao.Calculo(ativo, config, alerta, mr)
	if err != nil {
		return err
	}
	a.data = out
	return nil
}

func (a *Asset) GetAll() ([]tipos.Result, error) {
	if err := ReadConfig(); err != nil {
		return []tipos.Result{}, err
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	// desabilita saveLog
	config.SaveLog = false

	mr := make(map[string]*rsi.RSI)

	resposta := ""
	var outJson []tipos.Result
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo)
		_, _, out, err := cotacao.Calculo(atv, config, alerta, mr)
		if err != nil {
			return []tipos.Result{}, err
		}
		outJson = append(outJson, out)
		prettyJSON, _ := PrettyJson(out)
		resposta += prettyJSON + ","
	}

	return outJson, nil
}
