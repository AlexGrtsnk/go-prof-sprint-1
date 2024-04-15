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
	Flw "go-prof-sprint-1/internal/json_parser"
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
		shoortURL, flag, err := db.DataBaseCheckURLExistance(longURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flag == 1 {
			w.WriteHeader(http.StatusConflict)
			_, err = io.WriteString(w, apiRunAddr+"/"+string(shoortURL))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
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
			if id == "ping" {
				return
			}
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
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if db.OldName != "" {
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
}

func JSONPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		apiRunAddr, err := db.DataBaseCreateShortURLPageCfg()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
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
		shoortURL, flag, err := db.DataBaseCheckURLExistance(longURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flag == 1 {
			var answ Answer
			answ.Result = apiRunAddr + "/" + shoortURL
			resp, err := json.Marshal(answ)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusConflict)
			_, err = res.Write(resp)
			if err != nil {
				log.Fatal(err)
			}

			return
		}
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
		answ.Result = apiRunAddr + "/" + shortURL
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

func PingDataBasePage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		err := db.DataBasePingHandler()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(res, "cannot open psql database, using old realization")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			res.WriteHeader(http.StatusOK)
			_, err = io.WriteString(res, "connection went good, using psql databse")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
func UploadBatchFullURLPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		reader, err := gzp.GzipFormatHandlerJSON(res, req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		apiRunAddr, err := db.DataBaseCreateShortURLPageCfg()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}

		var newProduceItems Flw.ProduceList
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &newProduceItems); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		var tempItems []AnswerBatch
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		for idx, produceItem := range newProduceItems {
			if len(produceItem.OriginalURL) <= 0 {
				errMsg := fmt.Sprintf("Item %d: Incorrect produce code sequence or product name. Example code sequence: A12T-4GH7-QPL9-3N4M", idx)
				http.Error(res, errMsg, http.StatusBadRequest)
				return
			} else {
				shortURL := generateShortKey()
				err = db.DataBaseDownloadFullURLPagePost(shortURL, produceItem.OriginalURL)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					_, err = io.WriteString(res, "Error on the database side")
					if err != nil {
						log.Fatal(err)
					}
				}
				var answ AnswerBatch
				answ.CorrelationID = produceItem.CorrelationID
				answ.ShortURL = apiRunAddr + "/" + shortURL
				tempItems = append(tempItems, answ)
				err = db.DataBaseFilePost(shortURL, produceItem.OriginalURL)
				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					_, err = io.WriteString(res, "Error on the database side")
					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}
		res.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(res).Encode(tempItems); err != nil {
			log.Panic(err)
		}
	}
}

func Run() error {
	var cfg apcfg.Config
	err := env.Parse(&cfg)
	flagRunAddr, apiRunAddr, fileName, databaseDSN := apcfg.ParseFlags()
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
	if cfg.DatabaseDSN != "" {
		databaseDSN = cfg.DatabaseDSN
	}
	log.Println(cfg)
	err = db.DataBaseStartConfig(databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	a, b, err := db.DataBaseSelfConfigGet()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dasddasd ", a, b)
	if databaseDSN != "localhost" {
		err = db.DataBasePingHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("dasddasd213123 ", a, b)
	err = db.DataBaseCfg(flagRunAddr, apiRunAddr, fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.DataBaseInsert(fileName)
	if err != nil {
		log.Fatal(err)
	}
	//apiiRunAddr, _ := db.DataBaseCreateShortURLPageCfg()
	//err = db.DataBaseDownloadFullURLPagePost("aaa", "bb")
	//fmt.Println("cdjkdcn:" + string(err.Error()))
	fmt.Println("where postgres is hosted:", databaseDSN)
	fmt.Println("where db is held", fileName)
	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", apiRunAddr)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/{id}`, lg.WithLogging(apiHandler()))
	mux1.HandleFunc(`/`, lg.WithLogging(mainHandler()))
	mux1.HandleFunc(`/api/shorten`, lg.WithLogging(jsonHandler()))
	mux1.HandleFunc(`/ping`, lg.WithLogging(pingHandler()))
	mux1.HandleFunc(`/api/shorten/batch`, lg.WithLogging(batchHandler()))
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

type AnswerBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func jsonHandler() http.Handler {
	fn := JSONPage
	return http.HandlerFunc(fn)
}

func pingHandler() http.Handler {
	fn := PingDataBasePage
	return http.HandlerFunc(fn)
}

func batchHandler() http.Handler {
	fn := UploadBatchFullURLPage
	return http.HandlerFunc(fn)
}
