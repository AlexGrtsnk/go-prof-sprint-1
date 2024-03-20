package main

import (
	"bytes"
	"database/sql"
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
}

var db *sql.DB

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
		_, err := io.WriteString(w, "No get method allowed")
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
		flag := 0
		quer := "SELECT longURL FROM short_longURL WHERE short_url = '" + string(id) + "';"
		rows, err := db.Query(quer)
		if rows.Err() != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {

			var longURL string
			flag = 1
			err = rows.Scan(&longURL)
			res.Header().Set("Location", longURL)
			res.WriteHeader(http.StatusTemporaryRedirect)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s\n", longURL)
		}

		if flag == 0 {
			res.WriteHeader(http.StatusNotFound)
			_, err = io.WriteString(res, "No full url for this address")
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
		quer := "INSERT INTO short_longURL(short_url, longURL) VALUES('" + string(id) + "', '" + string(longURL) + "');"
		_, err := db.Exec(quer)
		if err != nil {
			log.Fatal(err)
		}

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
		flagRunAddr = "8080"
	}
	if cfg.baseURL != "" {
		vbn = cfg.baseURL
	}
	log.Println(cfg)
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sts := `
DROP TABLE IF EXISTS short_longURL;
CREATE TABLE short_longURL(id INTEGER PRIMARY KEY, short_url TEXT, longURL TEXT);
`
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", vbn)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/{id}`, apiPage)
	mux1.HandleFunc(`/`, mainPage)
	return http.ListenAndServe(flagRunAddr, mux1)
}

var flagRunAddr string
var vbn string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&vbn, "b", "http://localhost:8080", "api page existance url adress")
	flag.Parse()
	if flagRunAddr != "localhost:8080" && vbn == "http://localhost:8080" {
		vbn = "http://" + flagRunAddr
	}
	if flagRunAddr == "localhost:8080" && vbn != "http://localhost:8080" {
		flagRunAddr = vbn[7:]
	}
}
