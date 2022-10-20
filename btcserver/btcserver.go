package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/alexsetta/smartbot/cfg"
	"github.com/alexsetta/smartbot/cotacao"
	"github.com/alexsetta/smartbot/tipos"
	"github.com/gorilla/mux"
)

// Generated by https://quicktype.io
type Page struct {
	Result string
}

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	templates *template.Template
	porta     string
	hora      = time.Now().Add(time.Hour * -5)
	alerta    = tipos.Alertas{hora, hora, hora, hora, hora, hora}
	carteira  = tipos.Carteira{}
	config    = tipos.Config{}
	start     = time.Now()
)

func main() {
	porta = "8081"
	if len(os.Args) == 2 {
		porta = os.Args[1]
	}
	templates = template.Must(template.ParseFiles("./templates/result.html"))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/cotacao/{id}", Cotacao)
	router.HandleFunc("/total/", Total)

	fmt.Println("Listen port " + porta)
	log.Fatal(http.ListenAndServe(":"+porta, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "btcserver online\nstart time: ", time.Since(start))
}

func body(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("body: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("body: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("body: %w", err)
	}
	return string(b), nil
}

func formatJSONError(message string) []byte {
	appError := struct {
		Message string `json:"message"`
	}{
		message,
	}
	response, err := json.Marshal(appError)
	if err != nil {
		return []byte(err.Error())
	}
	return response
}

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

func Cotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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

	resposta := "["
	outJson := tipos.Result{}
	if id == "all" {
		for _, atv := range carteira.Ativos {
			_, _, out, err := cotacao.Calculo(atv, config, alerta)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			outJson = out
			prettyJSON, _ := PrettyJson(out)
			resposta += prettyJSON + ","
		}
		resposta = resposta[:len(resposta)-1] + "]"
	} else {
		ativo := tipos.Ativo{}
		for _, atv := range carteira.Ativos {
			if strings.ToLower(atv.Simbolo) == id {
				ativo = atv
				break
			}
		}
		var err2 error
		resposta, _, outJson, err2 = cotacao.Calculo(ativo, config, alerta)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		// remover as linhas abaixo para mostrar como "string"
		outJson = outJson
		prettyJSON, _ := PrettyJson(outJson)
		resposta = prettyJSON + ","
	}

	switch r.Header.Get("Accept") {
	case "*/*":
		saidaJson(w, outJson)
	default:
		w.Write([]byte(resposta))

		//saidaJson(w, resposta)

		//p := Page{Result: resposta}
		//saidaHtml(w, p)
	}
}

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

func saidaJson(w http.ResponseWriter, outJson interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(outJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatJSONError("Erro convertendo em JSON"))
		return
	}
}

func saidaHtml(w http.ResponseWriter, p Page) {
	if err := templates.ExecuteTemplate(w, "result.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
