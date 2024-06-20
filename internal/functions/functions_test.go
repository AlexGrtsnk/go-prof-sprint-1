package internal

import (
	"bytes"
	"fmt"
	apcfg "go-prof-sprint-1/internal/app_config"
	db "go-prof-sprint-1/internal/db"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			createShortURLPage(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCreateShortURLPage1(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 201,
				//response:    `{"status":"ok"}`,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}
	fmt.Println("gfgdrtgdt")
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
	fmt.Println("gfgdrtgdt")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			_, err = io.WriteString(b, "http://localhost:8080/dasdsdqwe")
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/dasdsdqwe", b)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			createShortURLPage(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestDownloadFullURLPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			downloadFullURLPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestJSONPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 200,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 200,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			jsonPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestUploadBatchFullURLPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 200,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 200,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			uploadBatchFullURLPage(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func TestGetConcreteURLSUser(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 401,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, m[i], nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			getConcreteURLSUser(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			//assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
		i++
	}
}

func BenchmarkGenerateShortKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateShortKey()
	}
}
