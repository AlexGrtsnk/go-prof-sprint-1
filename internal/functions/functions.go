package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	apcfg "go-prof-sprint-1/internal/app_config"
	ath "go-prof-sprint-1/internal/authentification"

	cks "go-prof-sprint-1/internal/cookies"
	db "go-prof-sprint-1/internal/db"
	gzp "go-prof-sprint-1/internal/gzp"
	flw "go-prof-sprint-1/internal/json_parser"
	lg "go-prof-sprint-1/internal/logger"

	"net/http/pprof"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// generateShortKey генерирует сокращенный url
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

// createShortURLPage отправляет пользователю его токен и адрес сокращенного url
func createShortURLPage(w http.ResponseWriter, r *http.Request) {
	reader, err := gzp.GzipFormatHandlerJSON(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if r.Method == http.MethodPost {
		apiRunAddr, err := db.DataBaseCreateShortURLPageCfg()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		var cookiesTmp *http.Cookie
		_, err = cks.GetCookieHandler(w, r)
		if err != nil {
			token, err := ath.BuildJWTString()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
			cookiesTmp = cks.SetCookieHandler(w, r, token)
		} else {
			cookiesTmp, err = r.Cookie("exampleCookie")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		rReader, err := io.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		longURL := string(rReader)
		if longURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}

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
			client := http.Client{}
			request, err := http.NewRequest("POST", apiRunAddr+"/"+string(shortURL), b)

			if err != nil {
				log.Fatal(err)
			}
			request.AddCookie(cookiesTmp)
			resp, err := client.Do(request)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			defer resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
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

// downloadFullURLPage при методе post пищет сокращенный url в базу данных, при методе get возращает полный url
func downloadFullURLPage(res http.ResponseWriter, req *http.Request) {
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
		flag, err := db.DataBaseCheckURLDelition(id)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flag == 1 {
			res.WriteHeader(http.StatusGone)
			return
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
		token, err := cks.GetCookieHandler(res, req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		a, _ := io.ReadAll(req.Body)
		longURL := string(a)
		vars := mux.Vars(req)
		id := vars["id"]
		err = db.DataBaseDownloadFullURLPagePost(id, longURL, token)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		err = db.DataBaseFilePost(id, longURL, token)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// jsonPage позволяет получить сокращенный url из полного в формате json
func jsonPage(res http.ResponseWriter, req *http.Request) {
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

		var token string
		_, err = cks.GetCookieHandler(res, req)
		if err != nil {
			token, err = ath.BuildJWTString()
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(res, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
			_ = cks.SetCookieHandler(res, req, token)
		} else {
			cookiesTmp, err := req.Cookie("exampleCookie")
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(res, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
			token = cookiesTmp.Value
		}
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
		err = db.DataBaseDownloadFullURLPagePost(shortURL, longURL, token)
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
		err = db.DataBaseFilePost(shortURL, longURL, "qew")
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// pingDataBasePage страница, по адресу которой пользователь может узнать состояние базы данных
func pingDataBasePage(res http.ResponseWriter, req *http.Request) {
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

// uploadBatchFullURLPage позволяет получить сразу несколько сокращенных url в формате json
func uploadBatchFullURLPage(res http.ResponseWriter, req *http.Request) {
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
		var token string
		_, err = cks.GetCookieHandler(res, req)
		if err != nil {
			token, err = ath.BuildJWTString()
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(res, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
			_ = cks.SetCookieHandler(res, req, token)
		} else {
			cookiesTmp, err := req.Cookie("exampleCookie")
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(res, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
			token = cookiesTmp.Value
		}

		var newProduceItems flw.ProduceList
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
				err = db.DataBaseDownloadFullURLPagePost(shortURL, produceItem.OriginalURL, token)
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
				err = db.DataBaseFilePost(shortURL, produceItem.OriginalURL, token)
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

// getConcreteURLSUser возврщает все url, загруженные конкретным пользователем, который отпраляет этот запрос
func getConcreteURLSUser(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(res, req)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		var userURLs []db.AnswerBatch
		userURLs, err = db.DataBaseGetAllURLs(token)
		var tmpURLs []db.AnswerBatch
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if userURLs == nil {
			res.WriteHeader(http.StatusNoContent)
			return

		} else {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
			tmpURLs = append(tmpURLs, userURLs[len(userURLs)-1])
			if err := json.NewEncoder(res).Encode(tmpURLs); err != nil {
				log.Panic(err)
			}
			return
		}
	}
	if req.Method == http.MethodDelete {
		reader, err := gzp.GzipFormatHandlerJSON(res, req)
		if err != nil {
			res.WriteHeader(http.StatusConflict)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var newDelitionItems flw.DeleteList
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadGateway)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &newDelitionItems); err != nil {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}
		token, err := cks.GetCookieHandler(res, req)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusAccepted)
		err = db.DataBaseDeleteURLs(newDelitionItems, token)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		}
	}

}

// Run определяет необходимые для работы приложения системные переменные, настраивает базу данных и запускает сам сервер
func Run() *http.Server {
	var cfg apcfg.Config
	var srv = http.Server{}
	err := env.Parse(&cfg)
	flagRunAddr, apiRunAddr, fileName, databaseDSN, enableHttps, config := apcfg.ParseFlags()
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
	if !cfg.EnableHttps {
		enableHttps = cfg.EnableHttps
	}
	if cfg.Config != "" {
		config = cfg.Config
	}
	log.Println(cfg)
	configFile, err := os.Open(config)
	if err != nil {
		fmt.Println("no file was found")
	}
	var setting flw.Setting
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&setting); err != nil {
		fmt.Println("wrong data in file. using old types")
	}
	if flagRunAddr == "8080" && setting.ServerAddress != "" {
		flagRunAddr = setting.ServerAddress
	}
	if apiRunAddr == "http://localhost:8080" && setting.BaseURL != "" {
		apiRunAddr = setting.BaseURL
	}
	if fileName == "text.txt" && setting.FileStoragePath != "" {
		fileName = setting.FileStoragePath
	}
	if databaseDSN == "localhost" && setting.DataBaseDSN != "" {
		databaseDSN = setting.DataBaseDSN
	}
	if !enableHttps && !setting.EnableHttps {
		enableHttps = setting.EnableHttps
	}
	err = db.DataBaseStartConfig(databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if databaseDSN != "localhost" {
		err = db.DataBasePingHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = db.DataBaseCfg(flagRunAddr, apiRunAddr, fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.DataBaseInsert(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("where postgres is hosted:", databaseDSN)
	fmt.Println("where db is held", fileName)
	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", apiRunAddr)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/api/user/urls`, lg.WithLogging(authHandler()))
	mux1.HandleFunc(`/api/shorten`, lg.WithLogging(jsonHandler()))
	mux1.HandleFunc(`/ping`, lg.WithLogging(pingHandler()))
	mux1.HandleFunc(`/api/shorten/batch`, lg.WithLogging(batchHandler()))
	mux1.HandleFunc(`/{id}`, lg.WithLogging(apiHandler()))
	mux1.HandleFunc(`/`, lg.WithLogging(mainHandler()))
	mux1.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux1.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux1.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux1.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux1.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux1.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))

	if !enableHttps {
		srv.Addr = flagRunAddr
		srv.Handler = gzp.GzipHandle(mux1)
		//<-idleConnsClosed
		//fmt.Println(srv)
		return &srv
		//<-idleConnsClosed
	} else {
		//<-idleConnsClosed
		//return http.ListenAndServeTLS(flagRunAddr, "certificate", "key", gzp.GzipHandle(mux1))
		return &srv
	}

}

func apiHandler() http.Handler {
	fn := downloadFullURLPage
	return http.HandlerFunc(fn)
}

func mainHandler() http.Handler {
	fn := createShortURLPage
	return http.HandlerFunc(fn)
}

// Question необходима для чтения запроса в формате json
type Question struct {
	// LongURL изначальный url, отправленный пользователем
	LongURL string `json:"url"`
}

// Answer необходима для отправки ответа пользователю в формате json
type Answer struct {
	// Result сокращенный url, отправляемый пользователю
	Result string `json:"result"`
}

// AnswerBatch  еобходима для чтения batch запроса в формате json
type AnswerBatch struct {
	// CorrelationID изначальный url, отправленный пользователем
	CorrelationID string `json:"correlation_id"`
	// ShortURL сокращенный url, отправляемый пользователем
	ShortURL string `json:"short_url"`
}

// NewAnser еобходима для отправки пользователю batch ответа в формате json
type NewAnser struct {
	// ShortURL сокращенный url, отправляемый пользователю
	ShortURL string `json:"ShortURL"`
	// OriginalURL полный url, отправляемый пользователю
	OriginalURL string `json:"OriginalURL"`
}

func jsonHandler() http.Handler {
	fn := jsonPage
	return http.HandlerFunc(fn)
}

func pingHandler() http.Handler {
	fn := pingDataBasePage
	return http.HandlerFunc(fn)
}

func batchHandler() http.Handler {
	fn := uploadBatchFullURLPage
	return http.HandlerFunc(fn)
}

func authHandler() http.Handler {
	fn := getConcreteURLSUser
	return http.HandlerFunc(fn)
}
