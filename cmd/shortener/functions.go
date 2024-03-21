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
	var db *sql.DB
	var vbn string
	db, err := sql.Open("sqlite3", "conf.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	quer := "SELECT vbn FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the database side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	defer rows.Close()
	if rows.Err() != nil {
		log.Fatal(err)
	}
	rows.Next()
	err = rows.Scan(&vbn)
	if err != nil {
		log.Fatal(err)
	}
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
		_, err = io.WriteString(w, "No get method allowed")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func apiPage(res http.ResponseWriter, req *http.Request) {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "shortlongurl.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if req.Method == http.MethodGet {
		vars := mux.Vars(req)
		id, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "bad request")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag := 0
		quer := "SELECT longURL FROM short_longURL WHERE short_url = '" + string(id) + "';"
		rows, err := db.Query(quer)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		defer rows.Close()
		if rows.Err() != nil {
			log.Fatal(err)
		}

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

func run() error {
	var db *sql.DB
	var db1 *sql.DB
	flagRunAddr, vbn := parseFlags()
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
	db, err = sql.Open("sqlite3", "shortlongurl.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sts := `
DROP TABLE IF EXISTS short_longURL;
CREATE TABLE short_longURL(id INTEGER PRIMARY KEY, short_url TEXT, longURL TEXT);`
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	db1, err = sql.Open("sqlite3", "conf.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db1.Close()
	sts1 := `
DROP TABLE IF EXISTS cfg;
CREATE TABLE cfg (id INTEGER PRIMARY KEY, flagRunAddr TEXT, vbn TEXT);
INSERT INTO cfg(flagRunAddr, vbn) VALUES ('` + string(flagRunAddr) + `', '` + string(vbn) + `');`
	_, err = db1.Exec(sts1)

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

func parseFlags() (a string, b string) {
	var flagRunAddr string
	var vbn string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&vbn, "b", "http://localhost:8080", "api page existance url adress")
	flag.Parse()
	if flagRunAddr != "localhost:8080" && vbn == "http://localhost:8080" {
		vbn = "http://" + flagRunAddr
	}
	if flagRunAddr == "localhost:8080" && vbn != "http://localhost:8080" {
		flagRunAddr = vbn[7:]
	}
	return flagRunAddr, vbn
}
