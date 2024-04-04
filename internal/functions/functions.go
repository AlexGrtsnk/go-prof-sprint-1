package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	apcfg "go-prof-sprint-1/internal/app_config"
	db "go-prof-sprint-1/internal/db"
	gzp "go-prof-sprint-1/internal/gzp"
	lg "go-prof-sprint-1/internal/logger"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func CreateShortURLPage(w http.ResponseWriter, r *http.Request) {
	apiRunAddr, err := db.DataBaseCreateShortURLPageCfg()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	reader, err := gzp.GzipFormatHandlerJSON(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	if r.Method == http.MethodPost {
		rReader, err := io.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		longURL := string(rReader)
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
			resp, err := http.Post(apiRunAddr+"/"+string(shortURL), "text/plain", b)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			w.WriteHeader(http.StatusCreated)
			_, err = io.WriteString(w, apiRunAddr+"/"+shortURL)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			http.Error(w, "cant create short url", http.StatusBadRequest)
		}
		return
	}
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "No get method allowed")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DownloadFullURLPage(res http.ResponseWriter, req *http.Request) {
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
		longURL, flag, err := db.DatBaseDownloadFullURLPageGet(id)
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
		err := db.DataBaseDownloadFullURLPagePost(id, longURL)
		if err != nil {
			log.Fatal(err)
		}
		err = db.DataBaseFilePost(id, longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func JSONPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		reader, err := gzp.GzipFormatHandlerJSON(res, req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}

		var ques Question
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
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
		err = db.DataBaseDownloadFullURLPagePost(shortURL, longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var answ Answer
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
		err = db.DataBaseFilePost(shortURL, longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func Run() error {
	var cfg apcfg.Config
	err := env.Parse(&cfg)
	flagRunAddr, apiRunAddr, fileName := apcfg.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress != "" {
		flagRunAddr = "8080"
	}
	if cfg.BaseURL != "" {
		apiRunAddr = cfg.BaseURL
	}
	if cfg.FileStoragePath != "" {
		fileName = cfg.FileStoragePath
	}
	log.Println(cfg)
	err = db.DataBaseCfg(flagRunAddr, apiRunAddr, fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.DataBaseInsert(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("where db is held", fileName)
	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", apiRunAddr)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/{id}`, lg.WithLogging(apiHandler()))
	mux1.HandleFunc(`/`, lg.WithLogging(mainHandler()))
	mux1.HandleFunc(`/api/shorten`, lg.WithLogging(jsonHandler()))
	return http.ListenAndServe(flagRunAddr, gzp.GzipHandle(mux1))
}

func apiHandler() http.Handler {
	fn := DownloadFullURLPage
	return http.HandlerFunc(fn)
}

func mainHandler() http.Handler {
	fn := CreateShortURLPage
	return http.HandlerFunc(fn)
}

type Question struct {
	LongURL string `json:"url"`
}

type Answer struct {
	Result string `json:"result"`
}

func jsonHandler() http.Handler {
	fn := JSONPage
	return http.HandlerFunc(fn)
}
