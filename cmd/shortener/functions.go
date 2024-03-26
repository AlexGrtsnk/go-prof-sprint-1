package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Home          string `env:"HOME"`
	serverAddress string `env:"serverAddress"`
	baseURL       string `env:"baseURL"`
	flnm          string `env:"FILE_STORAGE_PATH"`
}

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
	vbn, err := dbMnp()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	reader, err := xzpjsn(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	if r.Method == http.MethodPost {
		a, err := io.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		longURL := string(a)
		if longURL == "" {
			http.Error(w, "Bad data for url shortener", http.StatusBadRequest)
		}
		shortURL := generateShortKey()
		b := new(bytes.Buffer)
		_, err = io.WriteString(b, longURL)
		if err != nil {
			log.Fatal(err)
		}
		if shortURL != "" {
			resp, err := http.Post(vbn+"/"+string(shortURL), "text/plain", b)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			w.WriteHeader(http.StatusCreated)
			_, err = io.WriteString(w, vbn+"/"+shortURL)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			http.Error(w, "cant create short url", http.StatusBadRequest)
		}
		return
	} else {
		w.Header().Set("Location", "sadasdsadwwq")
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "No get method allowed")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func apiPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		vars := mux.Vars(req)
		id, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
			res.WriteHeader(http.StatusBadRequest)
			_, err := io.WriteString(res, "bad request")
			if err != nil {
				log.Fatal(err)
			}
		}
		longURL, flag, err := dbAppgGt(id)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flag == 0 {
			res.WriteHeader(http.StatusNotFound)
			_, err = io.WriteString(res, "No full url for this address")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			res.Header().Set("Location", longURL)
			res.WriteHeader(http.StatusTemporaryRedirect)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if req.Method == http.MethodPost {
		a, _ := io.ReadAll(req.Body)
		longURL := string(a)
		vars := mux.Vars(req)
		id := vars["id"]
		err := dbAppgPst(id, longURL)
		if err != nil {
			log.Fatal(err)
		}
		err = flpst(id, longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func jsonPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		reader, err := xzpjsn(res, req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}

		var ques Ques
		var buf bytes.Buffer
		// читаем тело запроса
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		// десериализуем JSON в Visitor
		if err = json.Unmarshal(buf.Bytes(), &ques); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		longURL := ques.LongURL

		shortURL := generateShortKey()
		b := new(bytes.Buffer)
		_, err = io.WriteString(b, longURL)
		if err != nil {
			log.Fatal(err)
		}
		err = dbAppgPst(shortURL, longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var answ Answ
		answ.Result = "http://localhost:8080/" + shortURL
		resp, err := json.Marshal(answ)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write(resp)
		if err != nil {
			log.Fatal(err)
		}
		err = flpst(shortURL, longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func run() error {
	var cfg Config
	err := env.Parse(&cfg)
	flagRunAddr, vbn, fileName := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.serverAddress != "" {
		flagRunAddr = "8080"
	}
	if cfg.baseURL != "" {
		vbn = cfg.baseURL
	}
	if cfg.flnm != "" {
		fileName = cfg.flnm
	}
	log.Println(cfg)
	err = dbMnCf(flagRunAddr, vbn, fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = dbins(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running server on DASDASDSDSADSAD", fileName)
	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", vbn)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/{id}`, WithLogging(apiHandler()))
	mux1.HandleFunc(`/`, WithLogging(mainHandler()))
	mux1.HandleFunc(`/api/shorten`, WithLogging(jsonHandler()))
	return http.ListenAndServe(flagRunAddr, gzipHandle(mux1))
}

func parseFlags() (a string, b string, f string) {
	var flagRunAddr string
	var vbn string
	var fileName string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&vbn, "b", "http://localhost:8080", "api page existance url adress")
	flag.StringVar(&fileName, "f", "/tmp/bmq9Ei", "txt file with short and long urls")
	flag.Parse()
	if flagRunAddr != "localhost:8080" && vbn == "http://localhost:8080" {
		vbn = "http://" + flagRunAddr
	}
	if flagRunAddr == "localhost:8080" && vbn != "http://localhost:8080" {
		flagRunAddr = vbn[7:]
	}
	return flagRunAddr, vbn, fileName
}

func apiHandler() http.Handler {
	fn := apiPage
	return http.HandlerFunc(fn)
}

func mainHandler() http.Handler {
	fn := mainPage
	return http.HandlerFunc(fn)
}

type Ques struct {
	LongURL string `json:"url"`
}

type Answ struct {
	Result string `json:"result"`
}

func jsonHandler() http.Handler {
	fn := jsonPage
	return http.HandlerFunc(fn)
}
