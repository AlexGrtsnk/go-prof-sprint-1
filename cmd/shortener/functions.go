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
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
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
	vbn, err := dbMnp()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the database side")
		if err != nil {
			log.Fatal(err)
		}
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

	}
}

func run() error {
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
	err = dbMnCf(flagRunAddr, vbn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running server on", flagRunAddr)
	fmt.Println("Running api on", vbn)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/{id}`, WithLogging(apiHandler()))
	mux1.HandleFunc(`/`, WithLogging(mainHandler()))
	mux1.HandleFunc(`/api/shorten`, WithLogging(jsonHandler()))
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

func apiHandler() http.Handler {
	//fn := a
	fn := apiPage
	return http.HandlerFunc(fn)
}

func mainHandler() http.Handler {
	//fn := a
	fn := mainPage
	return http.HandlerFunc(fn)
}

func jsonPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var ques Ques
		var buf bytes.Buffer
		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
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
		//a, _ := io.ReadAll(req.Body)
		//longURL := string(a)
		//vars := mux.Vars(req)
		//id := vars["url"]
		err = dbAppgPst(shortURL, longURL)
		if err != nil {
			log.Fatal(err)
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
		res.Write(resp)
	}
}

/*
curl -X POST http://localhost:8080/api/shorten -H 'Content-Type: application/json' -d '{"url": "https://practicum.yandex.ru"}'
*/
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

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithLogging(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	sugar := *logger.Sugar()
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// точка, где выполняется хендлер pingHandler
		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)

	}
	// возвращаем функционально расширенный хендлер
	return logFn
}
