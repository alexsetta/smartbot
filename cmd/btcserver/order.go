package main

import (
	"fmt"
	"github.com/alexsetta/smartbot/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//type AssetOrder struct {
//	qty   float64
//	value float64
//}

type OrderResult struct {
	time      string
	from      string
	fromQty   float64
	fromValue float64
	to        string
	toQty     float64
	toValue   float64
}

func (r *OrderResult) String() string {
	res := fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`, r.time, r.from, r.fromQty, r.fromValue, r.to, r.toQty, r.toValue)
	util.AppendFile("../../files/order.log", res)
	return res
}

func (r *OrderResult) Json() string {
	return fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`,
		r.time, r.from, r.fromQty, r.fromValue, r.to, r.toQty, r.toValue)
}

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
	log.Println("Order from", from, "to", to)

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

	dolar := toAsset.data.Preco
	if dolar == 0 {
		dolar = 1
	}
	res := OrderResult{
		time:      util.Now(),
		from:      fromAsset.data.Simbolo,
		fromQty:   fromAsset.data.Quantidade,
		fromValue: fromAsset.data.Preco,
		to:        toAsset.data.Simbolo,
		toQty:     fromAsset.data.Quantidade * fromAsset.data.Preco / dolar,
		toValue:   toAsset.data.Preco,
	}
	fmt.Println(res.String())
	fmt.Fprint(w, res.Json())
}
