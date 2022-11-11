package main

import (
	"errors"
	"github.com/alexsetta/smartbot/cotacao"
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
	return asset, nil
}

func (a *Asset) IsValid() error {
	if a.id == "" {
		return errors.New("id is empty")
	}
	return nil
}

func (a *Asset) Find() (tipos.Result, error) {
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	if err := ReadConfig(); err != nil {
		return tipos.Result{}, err
	}

	outJson := tipos.Result{}

	ativo := tipos.Ativo{}
	for _, atv := range carteira.Ativos {
		if strings.ToLower(atv.Simbolo) == a.id {
			ativo = atv
			break
		}
	}

	if ativo == (tipos.Ativo{}) {
		return tipos.Result{}, errors.New("Ativo não encontrado")
	}

	_, _, out, err := cotacao.Calculo(ativo, config, alerta)
	if err != nil {
		return tipos.Result{}, err
	}
	outJson = out

	return outJson, nil
}

func (a *Asset) GetAll() ([]tipos.Result, error) {
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	if err := ReadConfig(); err != nil {
		return []tipos.Result{}, err
	}

	resposta := ""
	outJson := []tipos.Result{}
	for _, atv := range carteira.Ativos {
		_, _, out, err := cotacao.Calculo(atv, config, alerta)
		if err != nil {
			return []tipos.Result{}, err
		}
		outJson = append(outJson, out)
		prettyJSON, _ := PrettyJson(out)
		resposta += prettyJSON + ","
	}

	return outJson, nil
}
