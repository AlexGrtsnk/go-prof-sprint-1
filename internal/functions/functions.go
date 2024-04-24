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
	ath "go-prof-sprint-1/internal/authentification"
	cks "go-prof-sprint-1/internal/cookies"
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

/*
	func setCookieHandler(w http.ResponseWriter, r *http.Request) {
		// Initialize a new cookie containing the string "Hello world!" and some
		// non-default attributes.
		fmt.Println("We are in cookies1")

		_, err := ath.BuildJWTString()
		if err != nil {
			return
		}
		fmt.Println("We are in cookies2")

		cookie := http.Cookie{
			Name:     "exampleCookie",
			Value:    "Help me",
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
		fmt.Println("We are in cookies")

		// Use the http.SetCookie() function to send the cookie to the client.
		// Behind the scenes this adds a `Set-Cookie` header to the response
		// containing the necessary cookie data.
		http.SetCookie(w, &cookie)
		fmt.Println("We are in cookies")

		// Write a HTTP response as normal.
		w.Write([]byte("cookie set!"))
	}

	func getCookieHandler(w http.ResponseWriter, r *http.Request) {
		// Retrieve the cookie from the request using its name (which in our case is
		// "exampleCookie"). If no matching cookie is found, this will return a
		// http.ErrNoCookie error. We check for this, and return a 400 Bad Request
		// response to the client.

		fmt.Println("We are in cookies8")

		cookie, err := r.Cookie("exampleCookie")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusBadRequest)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}
		fmt.Println("We are in cookies9")

		// Echo out the cookie value in the response body.
		w.Write([]byte(cookie.Value))
	}
*/
func CreateShortURLPage(w http.ResponseWriter, r *http.Request) {
	token, err := cks.GetCookieHandler(w, r)
	kol := 0
	if err != nil {
		kol += 1
	}
	fmt.Println("THIS IS TOKEN : ", token, kol)
	errr := ath.GetUserID(token)
	if errr == -1 {
		token, err = ath.BuildJWTString()
		if err != nil {
			kol += 1
		}
		cks.SetCookieHandler(w, r, token)
	}
	//token = "aaaa"

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
	token, err := cks.GetCookieHandler(res, req)
	kol := 0
	if err != nil {
		kol += 1
	}
	fmt.Println("HERE TEORETICALLY MUST BE COOKIE ", kol)
	errr := ath.GetUserID(token)
	if errr == -1 {
		token, err = ath.BuildJWTString()
		if err != nil {
			kol += 1
		}
		cks.SetCookieHandler(res, req, token)
	}
	//token = "aaaa"
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
		err := db.DataBaseDownloadFullURLPagePost(id, longURL, token)
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

func JSONPage(res http.ResponseWriter, req *http.Request) {
	token, err := cks.GetCookieHandler(res, req)
	kol := 0
	if err != nil {
		kol += 1
	}
	errr := ath.GetUserID(token)
	if errr == -1 {
		token, err = ath.BuildJWTString()
		if err != nil {
			kol += 1
		}
		cks.SetCookieHandler(res, req, token)
	}
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
	token, err := cks.GetCookieHandler(res, req)
	kol := 0
	if err != nil {
		kol += 1
	}
	errr := ath.GetUserID(token)
	if errr == -1 {
		token, err = ath.BuildJWTString()
		if err != nil {
			kol += 1
		}
		cks.SetCookieHandler(res, req, token)
	}
	//token = "aaaa"
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

func GetConcreteURLSUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	/*
		token, err := cks.GetCookieHandler(res, req)
		var a Flw.DeleteList
		var b Flw.DeleteURL
		var c Flw.DeleteURL
		b.CorrelationID = 1
		a = append(a, b)
		c.CorrelationID = 2
		a = append(a, c)
		db.DataBaseDeleteURLs(a, "aaaa")
		if err != nil {
			log.Panic(err)
		}
	*/
	if req.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(res, req)
		fmt.Println("What cookies are send????????????", token)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		kol := 0
		klnm, err := db.DataBaseGetAllURLs(token)
		//token := "aaaa"
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the database side")
			if err != nil {
				log.Fatal(err)
			}
		}
		if klnm == nil {
			res.WriteHeader(http.StatusNoContent)
			return

		} else {
			res.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(res).Encode(klnm); err != nil {
				log.Panic(err)
			}
			fmt.Println("HEEEEEELP   ,", res)
			return
		}
		fmt.Println("HEEEEEELP   ,", klnm)
		errr := ath.GetUserID(token)
		if errr != -1 {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		if errr == -1 {
			token, err = ath.BuildJWTString()
			if err != nil {
				kol += 1
			}
			cks.SetCookieHandler(res, req, token)
		}
		res.WriteHeader(http.StatusAccepted)
	}
	if req.Method == http.MethodDelete {
		reader, err := gzp.GzipFormatHandlerJSON(res, req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(res, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var newProduceItems Flw.DeleteList
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
		//var tempItems []AnswerBatch
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		db.DataBaseDeleteURLs(newProduceItems, "aaaa")
		/*
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
			}*/

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
	//mux1.HandleFunc(`/set`, setCookieHandler)
	//mux1.HandleFunc(`/get`, getCookieHandler)
	mux1.HandleFunc(`/api/user/urls`, lg.WithLogging(authHandler()))
	mux1.HandleFunc(`/api/shorten`, lg.WithLogging(jsonHandler()))
	mux1.HandleFunc(`/ping`, lg.WithLogging(pingHandler()))
	mux1.HandleFunc(`/api/shorten/batch`, lg.WithLogging(batchHandler()))
	mux1.HandleFunc(`/{id}`, lg.WithLogging(apiHandler()))
	mux1.HandleFunc(`/`, lg.WithLogging(mainHandler()))
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

func authHandler() http.Handler {
	fn := GetConcreteURLSUser
	return http.HandlerFunc(fn)
}
