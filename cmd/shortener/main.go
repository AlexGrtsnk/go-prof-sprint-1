package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
)

type Config struct {
	Home          string `env:"HOME"`
	serverAddress string `env:"serverAddress"`
	baseURL       string `env:"baseURL"`
}

type transform [1000][2]string

var mas transform
var kol int

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		longURL := string(a)
		if longURL == "" {
			http.Error(w, "Bad data for url shortener", http.StatusBadRequest)
		}
		shortURL := generateShortKey()
		b := new(bytes.Buffer)
		io.WriteString(b, longURL)
		if shortURL != "" {
			resp, err := http.Post("http://localhost"+string(flagRunAddr)+vbn+"/"+string(shortURL), "text/plain", b)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, "http://localhost"+string(flagRunAddr)+vbn+"/"+shortURL)
		} else {
			http.Error(w, "cant create short url", http.StatusBadRequest)
		}
		return
	} else {
		io.WriteString(w, "No get method allowed")
	}
}

func apiPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		vars := mux.Vars(req)
		id, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
			res.WriteHeader(http.StatusBadRequest)
			io.WriteString(res, "bad request")
		}
		flag := 0
		for i := 0; i < 1000; i++ {
			if mas[i][0] == string(id) {
				res.Header().Set("Location", mas[i][1])
				res.WriteHeader(http.StatusTemporaryRedirect)
				flag = 1
				break
			}
		}
		//return
		if flag == 0 {
			res.WriteHeader(http.StatusNotFound)
			io.WriteString(res, "No full url for this address")
		}
	}
	if req.Method == http.MethodPost {
		a, _ := io.ReadAll(req.Body)
		longURL := string(a)
		vars := mux.Vars(req)
		id := vars["id"]
		mas[kol][0] = string(id)
		mas[kol][1] = string(longURL)
		kol++
	}
}

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.serverAddress != "" {
		flagRunAddr = cfg.serverAddress
	}
	if cfg.baseURL != "" {
		vbn = cfg.baseURL
	}
	log.Println(cfg)

	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", vbn)
	kol = 0
	mux1 := mux.NewRouter()
	mux1.HandleFunc(vbn+`/{id}`, apiPage)
	mux1.HandleFunc(`/`, mainPage)
	return http.ListenAndServe(flagRunAddr, mux1)
}

var flagRunAddr string
var vbn string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&vbn, "b", "", "api page existance url adress")
	flag.Parse()
}
