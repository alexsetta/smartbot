package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) != 4 {
		log.Fatal("Digite: BTCBRL <moeda> <quantidade> <cotacao>")
	}

	moeda := os.Args[1]
	qtde, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatal("A quantidade deve ser um número (float) válido. Ex.: 0.12345678")
	}

	cotacao, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		log.Fatal("A cotacao deve ser um número (float) válido. Ex.: 61250.87")
	}

	taxa := 1 - 0.007
	final := 0.00
	res := ""
	if moeda == "btc" {
		final = (qtde * cotacao) * taxa
		res = fmt.Sprintf("R$: %.2f", final)
	} else {
		final = (qtde / cotacao) * taxa
		res = fmt.Sprintf("BTC: %.8f", final)
	}
	fmt.Println(res)
}
