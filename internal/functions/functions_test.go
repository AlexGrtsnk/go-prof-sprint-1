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
			name: "positive test #1",
			want: want{
				code: 200,
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
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "bbb",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			_, err = io.WriteString(b, "http://localhost:8080/multi")

			request := httptest.NewRequest(http.MethodPost, apiRunAddr+"/"+"qwerty", b)

			request.AddCookie(&cookie)
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
	}
}
func TestCreateShortURLPage8(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code: 409,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			_, _ = io.WriteString(b, "http://localhost:8080/multi")

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/"+"qwerty", b)

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
func TestCreateShortURLPage3(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code: 400,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			_, _ = io.WriteString(b, "http://localhost:8080/multi")

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/"+"qwerty", b)

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
	}
}
func TestCreateShortURLPage2(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #2",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			_, _ = io.WriteString(b, "http://localhost:8080/")
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", b)
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
	}
}
func TestCreateShortURLPage4(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #2",
			want: want{
				code: 400,
				//response:    `{"status":"ok"}`,
				//contentType: "text/plain; charset=utf-8",
			},
		},
	}
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "bbb",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			_, _ = io.WriteString(b, "http://localhost:8080/jhjhk")
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/jhjhk", b)
			request.AddCookie(&cookie)
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
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/shorten", nil)
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

func TestJSONPageGoodVibrations(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code: 201,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
		{
			name: "positivetest #2",
			want: want{
				code: 409,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	var m [1]string
	m[0] = "http://localhost:8080/qwerty12"
	//m[1] = "http://localhost:8080/qwerty"
	i := 0
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "bbb",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	//var jsonStr = []byte(`{"url":"Buy cheese and bread for breakfast."}`)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var jsonStr = []byte(`{"url":"http://localhost:8080/qwerty12"}`)
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewBuffer(jsonStr))
			// создаём новый Recorder
			request.AddCookie(&cookie)
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

func TestUploadBatchFullURLPageGoodVibrations(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive #1",
			want: want{
				code: 201,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
		{
			name: "positive test #2",
			want: want{
				code: 201,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	var m [1]string
	//m[0] = "http://localhost:8080/"
	m[0] = "http://localhost:8080/qwerty"
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "bbb",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	i := 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var jsonStr1 = []byte(`[{"correlation_id":"asd", "original_url":"http://localhost:8080/qwerty123"}, {"correlation_id":"asd1", "original_url":"http://localhost:8080/qwerty124"}]`)
			//var jsonStr2 = []byte(`{"correlation_id":"asd1", "original_url":"http://localhost:8080/qwerty124"}`)
			//var mas[2][]byte()
			//mas := append(mas, jsonStr1)
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/batch", bytes.NewBuffer(jsonStr1))
			// создаём новый Recorder
			request.AddCookie(&cookie)
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

func TestUploadBatchFullURLPage1(t *testing.T) {
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
func TestGetConcreteURLSUserbad(t *testing.T) {
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
				code: 403,
				//response:    `{"status":"ok"}`,
				contentType: "text/plain; charset=utf-8",
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
			i += 1
			var request *http.Request
			if i == 1 {
				request = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/urls", nil)
			} else {
				var delst = []byte(`["6qxTVvsy", "RTfd56hn"]`)
				request = httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/urls", bytes.NewBuffer(delst))
			}
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

func TestGetConcreteURLSUserGoodVibrations(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code: 200,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
		{
			name: "positive test #2",
			want: want{
				code: 204,
				//response:    `{"status":"ok"}`,
				contentType: "",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	cookiegood := http.Cookie{
		Name:     "exampleCookie",
		Value:    "bbb",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	cookiebad := http.Cookie{
		Name:     "exampleCookie",
		Value:    "ccc",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i += 1
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/urls", nil)
			// создаём новый Recorder
			if i == 1 {
				request.AddCookie(&cookiegood)
			} else {
				request.AddCookie(&cookiebad)
			}
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

func TestGetConcreteURLSUserGoodVibrations1(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code: 202,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
		{
			name: "positive test #2",
			want: want{
				code: 202,
				//response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	var m [2]string
	m[0] = "http://localhost:8080/"
	m[1] = "http://localhost:8080/qwerty"
	i := 0
	cookiegood := http.Cookie{
		Name:     "exampleCookie",
		Value:    "bbb",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	cookiebad := http.Cookie{
		Name:     "exampleCookie",
		Value:    "ccc",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i += 1
			var delst = []byte(`["8qZzNL", "UE03gB"]`)
			request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/urls", bytes.NewBuffer(delst))
			// создаём новый Recorder
			if i == 1 {
				request.AddCookie(&cookiegood)
			} else {
				request.AddCookie(&cookiebad)
			}
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

func TestBuildRun(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	err := Run()
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
